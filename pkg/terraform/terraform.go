package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type Terraform struct {
	WorkingDir string
	ExecDir    string
	Exec       *tfexec.Terraform
}

func NewTerraform(workingDir string, execDir string) (*Terraform, error) {
	tf, err := tfexec.NewTerraform(workingDir, execDir)
	if err != nil {
		return nil, fmt.Errorf("error new terraform: %w", err)
	}

	return &Terraform{
		WorkingDir: workingDir,
		ExecDir:    execDir,
		Exec:       tf,
	}, nil
}
