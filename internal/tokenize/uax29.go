package tokenize

import (
	"iter"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// UAX29 is loosely based on UAX #29.
// We require only that the tokenizer yield consistent meaningful word tokens.
// There is no need for characters with little to no meaning contribution.
// This relaxation allows us a lot of flexibility in how we implement UAX29.
type UAX29 struct{}

func NewUAX29() UAX29 {
	return UAX29{}
}

func (t UAX29) Tokens(text string) iter.Seq[string] {
	return func(yield func(string) bool) {
		c, n := utf8.DecodeRuneInString(text)

		for len(text) > 0 {
			for !isWordPart(c) && len(text) > 0 {
				text = text[n:]
				c, n = utf8.DecodeRuneInString(text)
			}

			builder := strings.Builder{}
			for !isBreak(c) && len(text) > 0 {
				if !isSkippable(c) {
					builder.WriteRune(c)
				}
				text = text[n:]
				c, n = utf8.DecodeRuneInString(text)
			}

			token := builder.String()
			if len(token) > 0 && !yield(token) {
				return
			}
		}
	}
}

// isBreak returns true if the rune causes a word break.
// It is the union of `!c.isWordPart()` and some end-of-word runes.
func isBreak(c rune) bool {
	// todo: add end of word runes
	return !isWordPart(c)
}

// isSkippable returns true if the rune should be excluded from tokens.
func isSkippable(c rune) bool {
	if slices.Contains(forceSkips, c) {
		return true
	}
	if slices.Contains(forceNoSkips, c) {
		return false
	}

	// Assume c ∈ Cf ∪ M ∪ Sk because of previous isWordPart call.
	// Within that set, we don't want Cf ∪ Me.
	return unicode.In(c, unicode.Cf, unicode.Me)
}

var forceSkips = []rune{
	'\u00a0',           // nbsp
	'\u00b7',           // middle dot
	'\uff9e', '\uff9f', // katakana sound marks
}

var forceNoSkips = []rune{
	'\u0027', // apostrophe
	'\u005f', // underscore
	'\u059e', // accent gershayim
}

// isWordPart returns true if the rune is part of a word or identifier.
func isWordPart(c rune) bool {
	if unicode.IsLetter(c) || unicode.IsNumber(c) {
		return true
	}
	if unicode.In(c, unicode.M, unicode.Cf) {
		return true
	}
	if _, ok := slices.BinarySearch(extraLetters, c); ok {
		return true
	}
	return false
}

var extraLetters = []rune{
	'\u0027', // apostrophe
	'\u005f', // underscore
	'\u00a0', // nbsp
	'\u00b7', // middle dot
	'\u00b8', // cedilla

	// modifier tone letters (Sk)
	'\u02c2', '\u02c3', '\u02c3', '\u02c4', '\u02c5', '\u02d2', '\u02d3',
	'\u02d4', '\u02d5', '\u02d6', '\u02d7', '\u02de', '\u02df', '\u02e5',
	'\u02e6', '\u02e7', '\u02e8', '\u02e9', '\u02ea', '\u02eb', '\u02ed',
	'\u02ef', '\u02f0', '\u02f1', '\u02f2', '\u02f3', '\u02f4', '\u02f5',
	'\u02f6', '\u02f7', '\u02f8', '\u02f9', '\u02fa', '\u02fb', '\u02fc',
	'\u02fd', '\u02fe', '\u02ff',

	'\u055a', '\u055b', '\u055c', '\u055e', '\u058a', // Armenian
	'\u059e', '\u05f3', // Hebrew
	'\u070f', // Syriac

	// modifier tone letters (Sk)
	'\ua708', '\ua709', '\ua70a', '\ua70b', '\ua70c', '\ua70d', '\ua70e', '\ua70f',
	'\ua710', '\ua711', '\ua712', '\ua713', '\ua714', '\ua715', '\ua716',
	'\ua720', '\ua721', '\ua789', '\ua78a', '\uab5b',

	'\uff9e', '\uff9f', // katakana sound marks
}
