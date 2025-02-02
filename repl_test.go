package main

// need this package, it gives us many useful code to write test
// just name the test with _test suffix
import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello    world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " LEbROn JameS Is THE Goat",
			expected: []string{"lebron", "james", "is", "the", "goat"},
		},
		{
			input:    "     ",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("The actual cleaned input lenght does not match expectation")
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word not match, expected: %s, found: %s", expectedWord, word)
				return
			}
		}
	}
}
