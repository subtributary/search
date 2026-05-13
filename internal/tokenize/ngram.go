package tokenize

import (
	"iter"
	"slices"
)

type NGram struct {
	minN, maxN int
}

func NewNGram(minN, maxN int) NGram {
	return NGram{minN, maxN}
}

func (t NGram) Tokens(text string) iter.Seq[string] {
	runes := []rune(text)

	// Text shorter than the minimum is just returned.
	if len(runes) == 0 {
		return slices.Values([]string{})
	}
	if len(runes) < t.minN {
		return slices.Values([]string{text})
	}

	// maxN can't be longer than the text.
	maxN := min(t.maxN, len(runes))

	// Return overlapping chunks of sizes from minN to maxN.
	return func(yield func(string) bool) {
		for n := t.minN; n <= maxN; n++ {
			for i := 0; i+n <= len(runes); i++ {
				if !yield(string(runes[i : i+n])) {
					return
				}
			}
		}
	}
}
