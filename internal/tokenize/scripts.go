package tokenize

import (
	"iter"
	"unicode/utf8"

	"github.com/subtributary/search/internal/shared"
)

type ScriptTokenizer struct{}

func NewScriptTokenizer() ScriptTokenizer {
	return ScriptTokenizer{}
}

func (t ScriptTokenizer) Tokens(text string) iter.Seq2[shared.Script, string] {
	return func(yield func(shared.Script, string) bool) {
		var script shared.Script
		var length int
		for ; len(text) > 0; text = text[length:] {
			script, length = currentScript(text)

			// Common script can be returned for leading Common characters.
			// Skip the leading Common script in this case.
			if script == "" {
				continue
			}

			if !yield(script, text[:length]) {
				return
			}
		}
	}
}

// currentScript returns the current script's length and name.
// Characters in the Common script are consumed without breaking.
func currentScript(text string) (shared.Script, int) {
	c, n := utf8.DecodeRuneInString(text)
	initScript := shared.DetectScript(c)

	i := n
	for ; i < len(text); i += n {
		c, n = utf8.DecodeRuneInString(text[i:])
		script := shared.DetectScript(c)

		if script == "" {
			// Common script doesn't break.
			continue
		}

		if script != initScript {
			break
		}
	}

	return initScript, i
}
