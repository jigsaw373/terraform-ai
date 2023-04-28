package terraform

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func (ter *Terraform) Init() error {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Start()

	err := ter.Exec.Init(context.Background())
	if err != nil {
		spin.Stop()

		return fmt.Errorf("error running Init: %w", err)
	}

	spin.Stop()

	return nil
}

func (ter *Terraform) Apply() error {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Start()

	err := ter.Exec.Apply(context.Background())
	if err != nil {
		spin.Stop()

		return fmt.Errorf("error running Apply: %w", err)
	}

	spin.Stop()

	return nil
}
