package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

func EndsWithTf(str string) bool {
	return strings.HasSuffix(str, ".tf")
}

func RandomName() string {
	// Initialize a byte slice of desired length
	randomBytes := make([]byte, 5)

	// Read random data into the byte slice
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// Encode the byte slice as a base64 string
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Truncate the string to the desired length

	return fmt.Sprintf("terraform-%s.tf", randomString)
}

func GetName(name string) string {
	name = RemoveBlankLinesFromString(name)
	if EndsWithTf(name) {
		return name
	}

	return RandomName()
}

func TerraformPath() (string, error) {
	cmd := exec.Command("which", "terraform")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running Init: %w", err)
	}

	return strings.TrimRight(string(output), "\n"), nil
}
