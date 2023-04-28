package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ia-ops/terraform-ai/pkg/terraform"
	"github.com/ia-ops/terraform-ai/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const initSubCommand = "You are a Terraform HCL generator, only generate valid provider Terraform HCL templates."

var errLength = errors.New("invalid length")

func addInit() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Run terraform init",
		RunE:  initCommand,
	}

	return initCmd
}

func initCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.Wrap(errLength, "prompt must be provided")
	}

	return initCmd(args)
}

func initCmd(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	oaiClients, err := newOAIClients()
	if err != nil {
		return fmt.Errorf("error new OAI client: %w", err)
	}

	completion, err := completion(ctx, oaiClients, args, *openAIDeploymentName, initSubCommand)
	if err != nil {
		return fmt.Errorf("error completation: %w", err)
	}

	text := fmt.Sprintf("\n😈 Attempting to apply the following template: %s", completion)
	log.Println(text)

	confirmation, err := terraform.GetApplyConfirmation(*requireConfirmation)
	if err != nil {
		return fmt.Errorf("error apply confirmation: %w", err)
	}

	if confirmation {
		if err = terraform.CheckTemplate(completion); err != nil {
			return fmt.Errorf("error check template: %w", err)
		}

		if err = utils.StoreFile("provider.tf", completion); err != nil {
			return fmt.Errorf("error store file: %w", err)
		}

		if err = ops.Init(); err != nil {
			return fmt.Errorf("error run terraform init: %w", err)
		}
	}

	return nil
}
