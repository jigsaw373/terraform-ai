package terraform

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func GetApplyConfirmation(requireConfirmation bool) (bool, error) {
	if !requireConfirmation {
		return true, nil
	}

	prompt := promptui.Select{
		Label: "Would you like to apply this? [Apply/Don't Apply]",
		Items: []string{"Apply", "Don't Apply"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("error prompt run: %w", err)
	}

	return result == "Apply", nil
}
