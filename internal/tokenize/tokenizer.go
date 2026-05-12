package tokenize

import "iter"

type Tokenizer interface {
	Tokenize() iter.Seq[string]
}
