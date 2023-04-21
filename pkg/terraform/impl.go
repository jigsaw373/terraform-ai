package terraform

import (
	"context"
	"fmt"
)

func (ter *Terraform) Init() error {
	err := ter.Exec.Init(context.Background())
	if err != nil {
		return fmt.Errorf("error running Init: %w", err)
	}

	return nil
}

func (ter *Terraform) Apply() error {
	err := ter.Exec.Apply(context.Background())
	if err != nil {
		return fmt.Errorf("error running Apply: %w", err)
	}

	return nil
}
