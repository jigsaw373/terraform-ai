package utils

import (
	"fmt"
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

func StoreFile(name string, contents string) error {
	contents = RemoveBlankLinesFromString(contents)

	err := os.WriteFile(name, []byte(contents), 0o600)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func CurrenDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error current dir: %w", err)
	}

	return currentDir, nil
}
