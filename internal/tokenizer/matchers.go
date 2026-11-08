package tokenizer

import "regexp"

// IsWhitespace checks if substring contains only whitespace characters.
func IsWhitespace(substring string) bool {
	re := regexp.MustCompile(`\s`)

	return re.MatchString(substring)
}

// IsAlpha checks if substring is alphabetic.
func IsAlpha(substring string) bool {
	re := regexp.MustCompile(`[_a-zA-Z]`)

	return re.MatchString(substring)
}

// IsAlphanumeric checks if substring is alphanumeric.
func IsAlphanumeric(substring string) bool {
	re := regexp.MustCompile(`[_0-9a-zA-Z]`)

	return re.MatchString(substring)
}

// IsNumeric checks if substring contains only digits.
func IsNumeric(substring string) bool {
	re := regexp.MustCompile("[0-9]")

	return re.MatchString(substring)
}
