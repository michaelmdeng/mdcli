package wiki

import (
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
