package utils

import (
	"log"
	"os"
)

func DirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Terraform is not initialized. Run `terraform init` first.")
	} else if err != nil {
		log.Fatalf("Failed to check if Terraform is initialized: %s\n", err.Error())
	}

	return true
}
