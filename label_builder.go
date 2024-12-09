package translator

import (
	"regexp"
	"strings"
	"unicode"
)

var invalidLabelCharRE = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// Normalizes the specified label to follow Prometheus label names standard.
//
// See rules at https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels.
//
// Labels that start with non-letter rune will be prefixed with "key_".
// An exception is made for double-underscores which are allowed.
func NormalizeLabel(label string) string {
	// Trivial case.
	if len(label) == 0 {
		return label
	}

	label = SanitizeLabelName(label)

	// If label starts with a number, prepend with "key_".
	if unicode.IsDigit(rune(label[0])) {
		label = "key_" + label
	} else if strings.HasPrefix(label, "_") && !strings.HasPrefix(label, "__") {
		label = "key" + label
	}

	return label
}

// SanitizeLabelName replaces anything that doesn't match
// client_label.LabelNameRE with an underscore.
// Note: this does not handle all Prometheus label name restrictions (such as
// not starting with a digit 0-9), and hence should only be used if the label
// name is prefixed with a known valid string.
func SanitizeLabelName(name string) string {
	return invalidLabelCharRE.ReplaceAllString(name, "_")
}
