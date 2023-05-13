package cli

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	apply     = "Apply"
	dontApply = "Don't Apply"
	reprompt  = "Reprompt"
)

func userActionPrompt() (string, error) {
	var (
		result string
		err    error
	)

	if !*requireConfirmation {
		return apply, nil
	}

	items := []string{apply, dontApply}
	label := fmt.Sprintf("Would you like to apply this? [%s/%s/%s]", reprompt, items[0], items[1])

	prompt := promptui.SelectWithAdd{
		Label:    label,
		Items:    items,
		AddLabel: reprompt,
	}
	_, result, err = prompt.Run()

	if err != nil {
		return dontApply, fmt.Errorf("error to run prompt: %w", err)
	}

	return result, nil
}
