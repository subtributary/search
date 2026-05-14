package search

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
