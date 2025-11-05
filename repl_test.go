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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "chipotle sauce is  good",
			expected: []string{"chipotle", "sauce", "is", "good"},
		}, {
			input:    "  pokemon,  world  ",
			expected: []string{"pokemon,", "world"},
		},
		{
			input:    "                please  please please ",
			expected: []string{"please", "please", "please"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("words are not equal, expected: %v, actual: %v", expectedWord, word)
			}
		}
	}
}
