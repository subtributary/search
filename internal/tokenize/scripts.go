package tokenize

import (
	"iter"
	"unicode"
	"unicode/utf8"
)

// Script is the name of a Unicode script or empty for the Common script.
type Script string

type ScriptTokenizer struct{}

func NewScriptTokenizer() ScriptTokenizer {
	return ScriptTokenizer{}
}

func (t ScriptTokenizer) Tokens(text string) iter.Seq2[Script, string] {
	return func(yield func(Script, string) bool) {
		var script Script
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
func currentScript(text string) (Script, int) {
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

func detectScript(c rune) Script {
	if !unicode.Is(unicode.Common, c) {
		for name, rt := range unicode.Scripts {
			if unicode.Is(rt, c) {
				return Script(name)
			}
		}
	}
	return ""
}
