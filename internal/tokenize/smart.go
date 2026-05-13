package tokenize

import (
	"iter"

	"github.com/subtributary/search/internal/shared"
)

// Tokenizer types implement a tokenization algorithm that can be delegated to
// by SmartTokenizer depending on script configuration.
type Tokenizer interface {
	Tokens(text string) iter.Seq[string]
}

// SmartTokenizer is the top-level tokenizer that delegates work to other
// tokenizers depending on the script. If no script is configured for a script,
// then its text is skipped.
type SmartTokenizer struct {
	main ScriptTokenizer
	subs map[shared.Script]Tokenizer
}

func NewSmartTokenizer(subs map[shared.Script]Tokenizer) SmartTokenizer {
	return SmartTokenizer{
		main: NewScriptTokenizer(),
		subs: subs,
	}
}

func (t SmartTokenizer) Tokens(text string) iter.Seq2[shared.Script, string] {
	return func(yield func(shared.Script, string) bool) {
		for script, token := range t.main.Tokens(text) {
			tokenizer, ok := t.subs[script]
			if !ok {
				// If no tokenizer is configured for the script, it is a no-op.
				// This ensures backwards compatibility for new scripts.
				continue
			}

			for subtoken := range tokenizer.Tokens(token) {
				if !yield(script, subtoken) {
					return
				}
			}
		}
	}
}
