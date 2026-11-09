package tokenizer_test

import (
	"testing"
	pkg "vladcheck/opi/internal/tokenizer"
)

func TestTokenizeTreatsKeywordAsKeyword(t *testing.T) {
	t.Parallel()

	input := "for"
	tokenizer := pkg.NewTokenizer(input)

	got, err := tokenizer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize(%q) returned unexpected error: %v", input, err)
	}

	if len(got) == 0 {
		t.Fatalf("Tokenize(%q) returned no tokens, want one", input)
	}

	want := pkg.TtFor
	if got[0].TokenType != want {
		t.Errorf("Tokenize(%q)[0].TokenType = %v, want %v", input, got[0].TokenType, want)
	}
}
