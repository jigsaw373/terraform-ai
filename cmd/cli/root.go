package cli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/walles/env"
)

const version = "0.0.1"

var (
	openAIDeploymentName = flag.String("openai-deployment-name", env.GetOr("OPENAI_DEPLOYMENT_NAME", env.String, "text-davinci-003"), "The deployment name used for the model in OpenAI service.")
	maxTokens            = flag.Int("max-tokens", env.GetOr("MAX_TOKENS", strconv.Atoi, 0), "The max token will overwrite the max tokens in the max tokens map.")
	openAIAPIKey         = flag.String("openai-api-key", env.GetOr("OPENAI_API_KEY", env.String, ""), "The API key for the OpenAI service. This is required.")
	azureOpenAIEndpoint  = flag.String("azure-openai-endpoint", env.GetOr("AZURE_OPENAI_ENDPOINT", env.String, ""), "The endpoint for Azure OpenAI service. If provided, Azure OpenAI service will be used instead of OpenAI service.")
	requireConfirmation  = flag.Bool("require-confirmation", env.GetOr("REQUIRE_CONFIRMATION", strconv.ParseBool, true), "Whether to require confirmation before executing the command. Defaults to true.")
	sensitivity          = flag.Float64("sensitivity", env.GetOr("SENSITIVITY", env.WithBitSize(strconv.ParseFloat, 64), 0.0), "The sensitivity to use for the model. Range is between 0 and 1. Set closer to 0 if your want output to be more deterministic but less creative. Defaults to 0.0.")
)

func InitAndExecute() {
	flag.Parse()

	if *openAIAPIKey == "" {
		log.Fatal("Please provide an OpenAI key.")
	}

	if err := RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "terraform-ai",
		Version: version,
		RunE:    runCommand,
	}

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("prompt must be provided")
	}

	return run(args)
}

func run(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	oaiClients, err := newOAIClients()
	if err != nil {
		return err
	}

	completion, err := gptCompletion(ctx, oaiClients, args, *openAIDeploymentName)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("✨ Attempting to apply the following template: %s", completion)
	fmt.Println(text)

	confirmation, err := getUserConfirmation(*requireConfirmation)
	if err != nil {
		return err
	}

	if confirmation {
		if err = applyTemplate(completion); err != nil {
			return err
		}
	}
	return nil
}

func getUserConfirmation(requireConfirmation bool) (bool, error) {
	if !requireConfirmation {
		return true, nil
	}

	prompt := promptui.Select{
		Label: "Would you like to apply this? [Apply/Don't Apply]",
		Items: []string{"Apply", "Don't Apply"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == "Apply", nil
}
