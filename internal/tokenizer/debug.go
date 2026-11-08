package tokenizer

import (
	"log"
	"os"
	"time"
)

//nolint:unused
func (t *Tokenizer) debugPrint() {
	debugLogger := log.New(os.Stdout, "DEBUG", 0)

	const sleepTime = 200 * time.Millisecond

	next, _ := t.peekNext()

	debugLogger.Printf("i: %d j: %d peek: %q next: %q\n", t.start, t.current, t.peek(), next)

	time.Sleep(sleepTime)
}
