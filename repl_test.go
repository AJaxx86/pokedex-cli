package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "   hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input: "this   is cleanedup text.",
			expected: []string{"this", "is", "cleanedup", "text."},
		},
		{
			input: "allonesentenceplease",
			expected: []string{"allonesentenceplease"},
		},
		{
			input: "This Sentence has some CAPITAL letters",
			expected: []string{"This", "Sentence", "has", "some", "CAPITAL", "letters"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) == 0 {
			t.Errorf("no strings returned")
			continue
		}
		
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("word '%s' does not match expected '%s'", word, expectedWord)
				continue
			}
		}
	}
}