package search

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/subtributary/search/internal/normalize"
	"github.com/subtributary/search/internal/rank"
	"github.com/subtributary/search/internal/tokenize"
)

type Option func(*Index) error

// WithField configures a document field for searching.
//
// weight is the relative weight of the field.
// A value of 1 should be used for the main part of the document.
func WithField(field string, weight float64) Option {
	const b = 0.72 // document length normalization strength
	return func(idx *Index) error {
		field := rank.Field(field)
		idx.fieldConfigs[field] = rank.FieldConfig{
			Weight: weight,
			B:      b,
		}
		return nil
	}
}

// WithNormalizers sets normalizers to use for a script.
// The script must match one in Go's Unicode library.
func WithNormalizers(script string, ids []Normalizer) Option {
	return func(idx *Index) error {
		norms := make([]normalize.Normalizer, len(ids))
		for i, id := range ids {
			if norm, err := id.toInternal(); err != nil {
				return err
			} else {
				norms[i] = norm
			}
		}
		idx.normalizers[script] = ids
		idx.normalizer.SetSubNormalizers(script, norms)
		return nil
	}
}

// WithTokenizer sets the tokenizer to use for a script.
// The script must match one in Go's Unicode library.
func WithTokenizer(script string, id Tokenizer) Option {
	return func(idx *Index) error {
		idx.tokenizers[script] = id
		if tok, err := id.toInternal(); err != nil {
			return err
		} else {
			idx.tokenizer.SetSubTokenizer(script, tok)
		}
		return nil
	}
}

type Index struct {
	version     string
	normalizers map[string][]Normalizer // Used for state saving normalizer
	tokenizers  map[string]Tokenizer    // Used for state saving tokenizer

	fieldConfigs map[rank.Field]rank.FieldConfig
	corpus       rank.Corpus
	normalizer   normalize.SmartNormalizer
	tokenizer    tokenize.SmartTokenizer
}

func NewIndex(opts ...Option) (*Index, error) {
	idx := &Index{
		version:      "0.0.0",
		normalizers:  make(map[string][]Normalizer),
		tokenizers:   make(map[string]Tokenizer),
		fieldConfigs: make(map[rank.Field]rank.FieldConfig),
		corpus:       rank.NewCorpus(),
		normalizer:   normalize.NewSmartNormalizer(),
		tokenizer:    tokenize.NewSmartTokenizer(),
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
		Version     string                          `json:"version"`
		Normalizers map[string][]Normalizer         `json:"normalizers"`
		Tokenizers  map[string]Tokenizer            `json:"tokenizers"`
		Fields      map[rank.Field]rank.FieldConfig `json:"fields"`
		Corpus      rank.Corpus                     `json:"corpus"`
	}{
		Version:     idx.version,
		Normalizers: idx.normalizers,
		Tokenizers:  idx.tokenizers,
		Fields:      idx.fieldConfigs,
		Corpus:      idx.corpus,
	})
}

// Upsert parses document fields and upserts the document into the corpus.
//
// Unconfigured fields will not be parsed, but they will be attached unchanged
// to search results. Configured fields are the opposite: they are parsed but
// not attached to search results. To have both, separate fields are needed.
//
// todo: ensure all configured fields are set.
func (idx *Index) Upsert(id string, fields map[string]string) {
	document := rank.NewDocument()
	for field, text := range fields {
		field := rank.Field(field)
		if _, ok := idx.fieldConfigs[field]; !ok {
			document.SetAttachment(field, text)
		} else {
			document.SetStream(field, idx.tokenize(text))
		}
	}
	idx.corpus.Upsert(id, document)
}

func (idx *Index) tokenize(text string) []string {
	tokens := make([]string, 0)
	for script, token := range idx.tokenizer.Tokens(text) {
		token = idx.normalizer.Normalize(script, token)
		tokens = append(tokens, token)
	}
	return tokens
}

type Result struct {
	Id          string
	Attachments map[string]string
}

func (idx *Index) Search(query string) iter.Seq[Result] {
	return nil
}

// Version returns the version number of search used to create the index.
func (idx *Index) Version() string {
	return idx.version
}
