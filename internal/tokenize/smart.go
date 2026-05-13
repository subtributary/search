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
// tokenizers and normalizers depending on the script.
//
// Text is segmented by script, a per-script tokenizer further tokenizes it,
// then per-script normalizers are run in order on the subtokens. If no script
// tokenizer is configured, then the script's text is skipped.
type SmartTokenizer struct {
	main  ScriptTokenizer
	norms map[Script][]Normalizer
	subs  map[Script]Tokenizer
}

type SmartOption func(*SmartTokenizer)

func WithScriptTokenizer(script Script, tokenizer Tokenizer) SmartOption {
	return func(t *SmartTokenizer) {
		t.subs[script] = tokenizer
	}
}

func WithScriptNormalizer(script Script, norm Normalizer) SmartOption {
	return func(t *SmartTokenizer) {
		t.norms[script] = append(t.norms[script], norm)
	}
}

func NewSmartTokenizer(opts ...SmartOption) SmartTokenizer {
	t := SmartTokenizer{
		main:  NewScriptTokenizer(),
		norms: make(map[Script][]Normalizer),
		subs:  make(map[Script]Tokenizer),
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t
}

func (t SmartTokenizer) Tokens(text string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for script, token := range t.main.Tokens(text) {
			tokenizer, ok := t.subs[script]
			if !ok {
				// If no tokenizer is configured for the script, it is a no-op.
				// This ensures backwards compatibility for new scripts.
				continue
			}

			normalizers := t.norms[script]

			for subtoken := range tokenizer.Tokens(token) {
				for _, n := range normalizers {
					subtoken = n(subtoken)
				}
				if !yield(subtoken) {
					return
				}
			}
		}
	}
}
