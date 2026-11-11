//nolint:gochecknoglobals
package tokenizer

import (
	"fmt"
	"log"
	"os"
	"time"
)

var debugLogger = log.New(os.Stdout, "", 0)

// DebugPrintAllTokens logs all tokens of []Token, each token on a new line and prettified.
func DebugPrintAllTokens(tokens []Token) {
	strs := make([]string, 0)

	for _, token := range tokens {
		if token.TokenType == TtNewLine {
			strs = append(strs, token.String())
			debugLogger.Println(strs)
			strs = make([]string, 0)

			continue
		}

		if token.Value == nil {
			strs = append(strs, token.String())
		} else {
			strs = append(strs, fmt.Sprintf("%s(%q)", token.String(), *token.Value))
		}
	}
}

// DebugPrintTokenizer prints current state of the Tokenizer instance.
func (t *Tokenizer) DebugPrintTokenizer() {
	const sleepTime = 200 * time.Millisecond

	next, _ := t.peekNext()

	debugLogger.Printf("i: %d j: %d peek: %q next: %q\n", t.start, t.current, t.peek(), next)

	time.Sleep(sleepTime)
}
