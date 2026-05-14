package tokenize

import (
	"iter"
	"unicode"
	"unicode/utf8"
)

type ScriptTokenizer struct{}

func NewScriptTokenizer() ScriptTokenizer {
	return ScriptTokenizer{}
}

func (t ScriptTokenizer) Tokens(text string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		var script string
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

// currentScript returns the current script's name and length.
// Characters in the Common script are consumed without breaking.
func currentScript(text string) (string, int) {
	c, n := utf8.DecodeRuneInString(text)
	initScript := detectScript(c)

	i := n
	for ; i < len(text); i += n {
		c, n = utf8.DecodeRuneInString(text[i:])
		script := detectScript(c)

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

// detectScript returns the Unicode script of the character.
// If the character has no script or is in the Common script, "" is returned.
func detectScript(c rune) string {
	if !unicode.Is(unicode.Common, c) {
		for name, rt := range unicode.Scripts {
			if unicode.Is(rt, c) {
				return name
			}
		}
	}
	return ""
}
