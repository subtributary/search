package search

import (
	"github.com/subtributary/search/internal/normalize"
	"github.com/subtributary/search/internal/rank"
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

// withDefaults is called automatically before user customizations.
func withDefaults() Option {
	defaults := make([]Option, 0, len(scripts)+3)
	norms := []Normalizer{NFKC, Lowercase}
	for _, script := range scripts {
		defaults = append(defaults, WithNormalizers(script, norms))
		defaults = append(defaults, WithTokenizer(script, UAX29))
	}
	defaults = append(defaults, WithTokenizer("Katakana", BigramTrigram))
	defaults = append(defaults, WithTokenizer("Han", UnigramBigram))
	defaults = append(defaults, WithTokenizer("Hiragana", BigramTrigram))

	return func(idx *Index) error {
		for _, opt := range defaults {
			if err := opt(idx); err != nil {
				return err
			}
		}
		return nil
	}
}
