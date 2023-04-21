package utils

import "strings"

// ToPtr converts type T to a *T as a convenience.
func ToPtr[T any](i T) *T {
	return &i
}

func RemoveBlankLinesFromString(input string) string {
	return strings.TrimLeft(input, "\n\r \t")
}
