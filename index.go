package search

import (
	"io/fs"
	"iter"

	"github.com/subtributary/search/internal/normalize"
	"github.com/subtributary/search/internal/tokenize"
)

type Index struct {
	tokenizer  tokenize.SmartTokenizer
	normalizer normalize.SmartNormalizer
}

func LoadIndex(root fs.FS, filename string) *Index {
	return nil
}

func NewIndex(config *Config) *Index {
	return nil
}

func (index *Index) Upsert(id string, fields map[string]string, attachment any) {
	//
}

type Result struct {
	Id         string
	Attachment any
}

func (index *Index) Search(query string) iter.Seq[Result] {
	return nil
}
