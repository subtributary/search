package search

import (
	"fmt"

	"github.com/subtributary/search/internal/tokenize"
)

type Tokenizer string

var (
	// Unigram divides text into single characters.
	// For example: "abcd" becomes {"a", "b", "c", "d"}.
	Unigram Tokenizer = "unigram"

	// Bigram divides text into pairs of characters.
	// For example: "abcd" becomes {"ab", "bc", "cd"}.
	Bigram Tokenizer = "bigram"

	// Trigram divides text into triplets of characters.
	// For example: "abcd" becomes {"abc", "bcd"}
	Trigram Tokenizer = "trigram"

	// UnigramBigram combines the tokens of Unigram and Bigram.
	UnigramBigram Tokenizer = "unigram bigram"

	// BigramTrigram combines the tokens of Bigram and Trigram.
	BigramTrigram Tokenizer = "bigram trigram"

	// UnigramBigramTrigram combines the tokens of UnigramBigram, and Trigram.
	UnigramBigramTrigram Tokenizer = "unigram bigram trigram"

	// UAX29 divides text on detected word breaks. The word break rules are
	// based on Unicode Standard Annex (UAX) #29.
	UAX29 Tokenizer = "uax29"
)

func (t Tokenizer) toInternal() (tokenize.Tokenizer, error) {
	switch t {
	case "unigram":
		return tokenize.NewNGram(1, 1), nil
	case "bigram":
		return tokenize.NewNGram(2, 2), nil
	case "trigram":
		return tokenize.NewNGram(3, 3), nil
	case "unigram bigram":
		return tokenize.NewNGram(1, 2), nil
	case "bigram trigram":
		return tokenize.NewNGram(2, 3), nil
	case "unigram bigram trigram":
		return tokenize.NewNGram(1, 3), nil
	case "uax29":
		return tokenize.NewUAX29(), nil
	default:
		return nil, fmt.Errorf("invalid tokenizer: %s", t)
	}
}
