package utils_test

import (
	"os"
	"testing"

	"github.com/ia-ops/terraform-ai/pkg/utils"
)

func TestStoreFile(t *testing.T) {
	// Create a temporary file to store the test data
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}

	defer os.Remove(tmpfile.Name())

	// Define the test input and expected output
	name := tmpfile.Name()
	contents := "test data"

	// Call the function being tested
	err = utils.StoreFile(name, contents)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Read the contents of the file and compare with the expected output
	actual, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("failed to read file: %s", err)
	}

	if string(actual) != contents {
		t.Errorf("unexpected file contents: got %s, want %s", actual, contents)
	}
}
