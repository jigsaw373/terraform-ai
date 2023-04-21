package utils_test

import (
	"testing"

	"github.com/ia-ops/terraform-ai/pkg/utils"
)

func TestRemoveBlankLinesFromString(t *testing.T) {
	input := "\n\n\nHello, world!\n\nHow are you?\n\n\n"
	expectedOutput := "Hello, world!\n\nHow are you?\n\n\n"

	output := utils.RemoveBlankLinesFromString(input)

	if output != expectedOutput {
		t.Errorf("Expected output '%s', but got '%s'", expectedOutput, output)
	}
}
