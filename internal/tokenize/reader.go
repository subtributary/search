package tokenize

import (
	"unicode"
	"unicode/utf8"
)

type Script string

type CharReader struct {
	text string

	// Current is the most recently read utf8 rune.
	Current rune

	// Script is the most recently detected script.
	// If punctuation is encountered, for example, the script does not change.
	Script Script
}

func NewCharReader(text string) *CharReader {
	return &CharReader{
		text: text,
	}
}

func (r *CharReader) Read() {
	c, i := utf8.DecodeRuneInString(r.text)
	r.text = r.text[i:]
	r.Current = c
	if s := script(c); s != "" {
		r.Script = script(c)
	}
}

func script(c rune) Script {
	if unicode.Is(unicode.Common, c) {
		return ""
	}

	for name, rt := range unicode.Scripts {
		if unicode.Is(rt, c) {
			return Script(name)
		}
	}

	return ""
}
