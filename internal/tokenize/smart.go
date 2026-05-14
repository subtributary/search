package tokenize

import (
	"iter"
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
	subs map[string]Tokenizer
}

func NewSmartTokenizer() SmartTokenizer {
	return SmartTokenizer{
		main: NewScriptTokenizer(),
		subs: make(map[string]Tokenizer),
	}
}

// SetSubTokenizer sets the tokenizer to use for a specific script.
// If the id is invalid, an error is returned.
func (t *SmartTokenizer) SetSubTokenizer(script string, tok Tokenizer) {
	t.subs[script] = tok
}

// Tokens returns tokenized text as pairs of <script, token>.
func (t *SmartTokenizer) Tokens(text string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for script, token := range t.main.Tokens(text) {
			tokenizer, ok := t.subs[script]
			if !ok {
				// If no tokenizer is configured for the script, it is a no-op.
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
