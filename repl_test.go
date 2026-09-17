package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "Book ah flight go Dubai KILLAS on board",
			expected: []string{"book", "ah", "flight", "go", "dubai", "killas", "on", "board"},
		},
		{
			input:    "ONEDRIVE Sucks",
			expected: []string{"onedrive", "sucks"},
		},
	}

	for _, testCase := range testCases {
		actual := cleanInput(testCase.input)
		if len(actual) != len(testCase.expected) {
			t.Errorf("input and output should be equal.\n input: %v \n output: %v", actual, testCase.expected)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := testCase.expected[i]

			if word != expectedWord {
				t.Errorf("current word: %v does not match the expected word: %v", word, expectedWord)
			}
		}
	}

}
