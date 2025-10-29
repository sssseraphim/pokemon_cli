package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " Hello World   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " bob    bean bob",
			expected: []string{"bob", "bean", "bob"},
		},
		{
			input:    "  ",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("wrong length, %v != %v", c.expected, actual)
			return
		}
		for i := range c.expected {
			if c.expected[i] != actual[i] {
				t.Errorf("%s != %s", c.expected[i], actual[i])
			}
		}
	}
}
