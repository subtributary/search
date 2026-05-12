package tokenize

import "iter"

type ScriptTokenizer struct {
	text       string
	tokenizers map[string]*Tokenizer
}

func NewScriptTokenizer(text string, tokenizers map[string]*Tokenizer) *ScriptTokenizer {
	return &ScriptTokenizer{text, tokenizers}
}

func (t *ScriptTokenizer) Tokenize() iter.Seq[string] {
	return nil
}
