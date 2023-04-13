package wiki

import (
	"errors"
	"testing"
)

func Test_GenerateLink(t *testing.T) {
	cases := map[string]string{
		"foo":                 "[foo](foo)",
		"hello world":         "[hello world](hello-world)",
		"Bar":                 "[Bar](bar)",
		"camelCase":           "[camelCase](camelcase)",
		"This is a sentence.": "[This is a sentence.](this-is-a-sentence.)",
		"multi\nline":         "[multi\nline](multi-line)",
	}

	for input, expected := range cases {
		actual := GenerateLink(input)
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func Test_HtmlOutputPath(t *testing.T) {
	type testCase struct {
		input string
		expected string
		expectedErr error
	}
	cases := []testCase{
		{
			input: "/home/mdeng/MyDrive/vimwiki/wiki/index.md",
			expected: "/home/mdeng/MyDrive/vimwiki/html/index.html",
		},
		{
			input: "/home/mdeng/MyDrive/vimwiki/wiki/foo/bar/baz.md",
			expected: "/home/mdeng/MyDrive/vimwiki/html/foo/bar/baz.html",
		},
		{
			input: "/foo/bar/baz.md",
			expectedErr: errors.New("not a wiki path"),
		},
	}

	for _, c := range cases {
		actual, err := HtmlOutputPath(c.input)
		if c.expectedErr != nil && err == nil {
			t.Errorf("expected error, got none")
		}
		if actual != c.expected {
			t.Errorf("expected %s, got %s", c.expected, actual)
		}
	}
}
