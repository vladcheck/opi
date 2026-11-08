package tokenizer

import "regexp"

func isWhitespace(char string) bool {
	re := regexp.MustCompile(`\s`)

	return re.MatchString(char)
}

func isDigit(char string) bool {
	re := regexp.MustCompile("[0-9]")

	return re.MatchString(char)
}
