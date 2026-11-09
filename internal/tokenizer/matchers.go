package tokenizer

import "regexp"

var (
	alphaRE          = regexp.MustCompile(`[_a-zA-Z]+`)
	isAlphanumericRE = regexp.MustCompile(`[_0-9a-zA-Z]+`)
	isNumericRE      = regexp.MustCompile("[0-9]+")
)

// IsHorizontalWhitespace checks if substring contains only whitespace characters except \n.
func IsHorizontalWhitespace(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r != ' ' && r != '\t' {
			return false
		}
	}

	return true
}

// IsAlpha checks if substring is alphabetic.
func IsAlpha(substring string) bool {
	return alphaRE.MatchString(substring)
}

// IsAlphanumeric checks if substring is alphanumeric.
func IsAlphanumeric(substring string) bool {
	return isAlphanumericRE.MatchString(substring)
}

// IsNumeric checks if substring contains only digits.
func IsNumeric(substring string) bool {
	return isNumericRE.MatchString(substring)
}
