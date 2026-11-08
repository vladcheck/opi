package tokenizer

import (
	"log"
	"time"
)

// DebugPrintAllTokens logs all tokens of []Token, each token on a new line and prettified.
func DebugPrintAllTokens(tokens []Token) {
	for _, t := range tokens {
		log.Printf("%s\n", t.String())
	}
}

// DebugPrintTokenizer prints current state of the Tokenizer instance.
func (t *Tokenizer) DebugPrintTokenizer() {
	const sleepTime = 200 * time.Millisecond

	next, _ := t.peekNext()

	log.Printf("i: %d j: %d peek: %q next: %q\n", t.start, t.current, t.peek(), next)

	time.Sleep(sleepTime)
}
