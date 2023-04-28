package main

import (
	"log"

	"github.com/ia-ops/terraform-ai/cmd/cli"
	"github.com/ia-ops/terraform-ai/pkg/utils"
)

func main() {
	workingDir, err := utils.CurrenDir()
	if err != nil {
		log.Fatalf("Failed get current dir: %s\n", err.Error())
	}

	execDir, err := utils.TerraformPath()
	if err != nil {
		log.Fatalf("Failed get exec dir: %s\n", err.Error())
	}

	cli.InitAndExecute(workingDir, execDir)
}
