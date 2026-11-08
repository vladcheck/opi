package tokenizer

import "regexp"

func isWhitespace(char string) bool {
	re := regexp.MustCompile(`\s`)

	return re.MatchString(char)
}

func isAlpha(substring string) bool {
	re := regexp.MustCompile(`[_a-zA-Z]`)

	return re.MatchString(substring)
}

func isAlphanumeric(substring string) bool {
	re := regexp.MustCompile(`[_0-9a-zA-Z]`)

	return re.MatchString(substring)
}

func isDigit(char string) bool {
	re := regexp.MustCompile("[0-9]")

	return re.MatchString(char)
}
