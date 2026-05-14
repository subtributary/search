package search

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/subtributary/search/internal/normalize"
	"github.com/subtributary/search/internal/tokenize"
)

type Option func(*Index) error

// WithNormalizers sets normalizers to use for a script.
// The script must match one in Go's Unicode library.
func WithNormalizers(script string, ids []Normalizer) Option {
	return func(idx *Index) error {
		casts := stringsFromNorms(ids)
		idx.normalizers[script] = casts
		return idx.normalizer.SetSubNormalizers(script, casts)
	}
}

// WithTokenizer sets the tokenizer to use for a script.
// The script must match one in Go's Unicode library.
func WithTokenizer(script string, id Tokenizer) Option {
	return func(idx *Index) error {
		idx.tokenizers[script] = string(id)
		return idx.tokenizer.SetSubTokenizer(script, string(id))
	}
}

type Index struct {
	version     string
	normalizers map[string][]string // Used for state saving/loading.
	tokenizers  map[string]string   // Used for state saving/loading.
	tokenizer   tokenize.SmartTokenizer
	normalizer  normalize.SmartNormalizer
}

func NewIndex(opts ...Option) (*Index, error) {
	idx := &Index{
		version:     "0.0.0",
		normalizers: make(map[string][]string),
		tokenizers:  make(map[string]string),
		tokenizer:   tokenize.NewSmartTokenizer(),
		normalizer:  normalize.NewSmartNormalizer(),
	}

	// Default normalizers and tokenizers for all scripts
	defaults := make([]Option, 0, len(scripts)+3)
	norms := []Normalizer{NFKC, Lowercase}
	for _, script := range scripts {
		defaults = append(defaults, WithNormalizers(script, norms))
		defaults = append(defaults, WithTokenizer(script, UAX29))
	}
	defaults = append(defaults, WithTokenizer("Katakana", BigramTrigram))
	defaults = append(defaults, WithTokenizer("Han", UnigramBigram))
	defaults = append(defaults, WithTokenizer("Hiragana", BigramTrigram))
	for _, opt := range defaults {
		if err := opt(idx); err != nil {
			return nil, fmt.Errorf("default option: %w", err)
		}
	}

	// Apply customizations
	for _, opt := range opts {
		if err := opt(idx); err != nil {
			return nil, err
		}
	}

	return idx, nil
}

func (idx *Index) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version     string              `json:"version"`
		Normalizers map[string][]string `json:"normalizers"`
		Tokenizers  map[string]string   `json:"tokenizers"`
	}{
		Version:     idx.version,
		Normalizers: idx.normalizers,
		Tokenizers:  idx.tokenizers,
	})
}

func (idx *Index) Upsert(id string, fields map[string]string, attachment any) {
	//
}

type Result struct {
	Id         string
	Attachment any
}

func (idx *Index) Search(query string) iter.Seq[Result] {
	return nil
}

// Version returns the version number of search used to create the index.
func (idx *Index) Version() string {
	return idx.version
}
