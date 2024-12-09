package translator

import (
	"fmt"
	"testing"
)

func TestNormalizeLabel(t *testing.T) {
	tests := []struct {
		label    string
		expected string
	}{
		{"", ""},
		{"label:with:colons", "label_with_colons"},
		{"LabelWithCapitalLetters", "LabelWithCapitalLetters"},
		{"label!with&special$chars)", "label_with_special_chars_"},
		{"label_with_foreign_characters_字符", "label_with_foreign_characters___"},
		{"label.with.dots", "label_with_dots"},
		{"123label", "key_123label"},
		{"_label_starting_with_underscore", "key_label_starting_with_underscore"},
		{"__label_starting_with_2underscores", "__label_starting_with_2underscores"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			result := NormalizeLabel(test.label)
			if test.expected != result {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
