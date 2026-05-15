package search

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/subtributary/search/internal/normalize"
	"github.com/subtributary/search/internal/rank"
	"github.com/subtributary/search/internal/tokenize"
)

const version = "0.1.0"

type savedState struct {
	Version     string                          `json:"version"`
	Normalizers map[string][]Normalizer         `json:"normalizers"`
	Tokenizers  map[string]Tokenizer            `json:"tokenizers"`
	Fields      map[rank.Field]rank.FieldConfig `json:"fields"`
	Corpus      rank.Corpus                     `json:"corpus"`
}

func (s *savedState) apply(idx *Index) error {
	if currVer, err := parseSemVer(version); err != nil {
		return fmt.Errorf("invalid current version: %s", version)
	} else if srcVer, err := parseSemVer(s.Version); err != nil {
		return fmt.Errorf("invalid version: %s", s.Version)
	} else if !currVer.canLoad(srcVer) {
		return fmt.Errorf("unable to load version: %s", s.Version)
	}

	// These are fully serialized so can just be set as-is.
	idx.fieldConfigs = s.Fields
	idx.corpus = s.Corpus

	// Normalizers only save their configurations, so they need rebuilt.
	idx.normalizer = normalize.NewSmartNormalizer()
	idx.normalizers = map[string][]Normalizer{}
	for script, norms := range s.Normalizers {
		if err := WithNormalizers(script, norms)(idx); err != nil {
			return fmt.Errorf("instantiating normalizer: %w", err)
		}
	}

	// Tokenizers only save their configurations, so they need rebuilt.
	idx.tokenizer = tokenize.NewSmartTokenizer()
	idx.tokenizers = map[string]Tokenizer{}
	for script, tok := range s.Tokenizers {
		if err := WithTokenizer(script, tok)(idx); err != nil {
			return fmt.Errorf("instantiating tokenizer: %w", err)
		}
	}

	return nil
}

type Index struct {
	normalizers map[string][]Normalizer // Used for state saving normalizer
	tokenizers  map[string]Tokenizer    // Used for state saving tokenizer

	fieldConfigs map[rank.Field]rank.FieldConfig
	corpus       rank.Corpus
	normalizer   normalize.SmartNormalizer
	tokenizer    tokenize.SmartTokenizer
}

// NewIndex creates a new search indexer with the provided options.
// An error is returned if the configuration is invalid.
func NewIndex(opts ...Option) (*Index, error) {
	idx := &Index{
		normalizers:  make(map[string][]Normalizer),
		tokenizers:   make(map[string]Tokenizer),
		fieldConfigs: make(map[rank.Field]rank.FieldConfig),
		corpus:       rank.NewCorpus(),
		normalizer:   normalize.NewSmartNormalizer(),
		tokenizer:    tokenize.NewSmartTokenizer(),
	}

	if err := withDefaults()(idx); err != nil {
		return nil, fmt.Errorf("error setting default options: %w", err)
	}

	for _, opt := range opts {
		if err := opt(idx); err != nil {
			return nil, err
		}
	}

	return idx, nil
}

func (idx *Index) MarshalJSON() ([]byte, error) {
	return json.Marshal(&savedState{
		Version:     version,
		Normalizers: idx.normalizers,
		Tokenizers:  idx.tokenizers,
		Fields:      idx.fieldConfigs,
		Corpus:      idx.corpus,
	})
}

func (idx *Index) UnmarshalJSON(data []byte) error {
	var state savedState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	return state.apply(idx)
}

// Upsert parses document fields and upserts the document into the corpus.
//
// Fields that have a configuration are parsed into internal metadata,
// but fields that do not have a configuration are treated as attachments.
// Attachments are returned unaltered alongside search results.
//
// An error is returned if any configured fields are missing from the document.
func (idx *Index) Upsert(id string, fields map[string]string) error {
	document := rank.NewDocument()
	for field, text := range fields {
		if _, ok := idx.fieldConfigs[rank.Field(field)]; !ok {
			document.SetAttachment(field, text)
		} else {
			document.SetStream(rank.Field(field), idx.tokenize(text))
		}
	}

	// Validation
	if len(document.Streams) != len(idx.fieldConfigs) {
		var missing []string
		for field := range idx.fieldConfigs {
			if _, ok := document.Streams[field]; !ok {
				missing = append(missing, string(field))
			}
		}
		return fmt.Errorf("document is missing fields: %v", missing)
	}

	idx.corpus.Upsert(id, document)
	return nil
}

type Result struct {
	Id          string
	Attachments map[string]string
}

// Search returns all documents sorted by how well they match the query.
// The best match is returned first.
// Equal matches are returned in alphabetical order by their id.
func (idx *Index) Search(query string) iter.Seq[Result] {
	queryTokens := idx.tokenize(query)
	bm := rank.NewBM25F(1.2, idx.fieldConfigs)
	results := bm.Rank(idx.corpus, queryTokens)

	return func(yield func(Result) bool) {
		for _, result := range results {
			converted := Result{
				Id:          result.Id,
				Attachments: result.Document.Attachments,
			}
			if !yield(converted) {
				return
			}
		}
	}
}

func (idx *Index) tokenize(text string) []string {
	tokens := make([]string, 0)
	for script, token := range idx.tokenizer.Tokens(text) {
		token = idx.normalizer.Normalize(script, token)
		tokens = append(tokens, token)
	}
	return tokens
}
