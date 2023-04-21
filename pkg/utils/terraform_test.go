package utils_test

import (
	"testing"

	"github.com/ia-ops/terraform-ai/pkg/utils"
)

func TestEndsWithTf(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"main.tf", true},
		{"provider.txt", false},
		{"provider.tf", true},
		{"subnet.doc", false},
	}

	for _, c := range cases {
		result := utils.EndsWithTf(c.input)

		if result != c.expected {
			t.Errorf("EndsWithTf(%q) == %v, expected %v", c.input, result, c.expected)
		}
	}
}

func TestRandomName(t *testing.T) {
	name1 := utils.RandomName()
	name2 := utils.RandomName()

	if name1 == name2 {
		t.Errorf("Expected unique names, but got: %s, %s", name1, name2)
	}
}
