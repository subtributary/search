package rank_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/subtributary/search/internal/rank"
)

func TestBM25F(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		query    string
		want     []string
	}{
		{
			name:     "empty",
			contents: []string{},
			query:    "test",
			want:     []string{},
		},
		{
			name:     "single match",
			contents: []string{"blue", "test", "blue tulips"},
			query:    "test",
			want:     []string{"test", "blue", "blue tulips"},
		},
		{
			name:     "multiple matches",
			contents: []string{"blue", "test", "blue tulips"},
			query:    "blue",
			want:     []string{"blue", "blue tulips", "test"},
		},
		{
			name:     "overused word",
			contents: []string{"blue", "test", "test test"},
			query:    "test",
			want:     []string{"test test", "test", "blue"},
		},
		{
			name:     "multiword query",
			contents: []string{"test", "blue", "blue tulips"},
			query:    "tulips blue",
			want:     []string{"blue tulips", "blue", "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			corpus := rank.NewCorpus()
			for _, content := range tt.contents {
				filename := content
				document := rank.NewDocument()
				document.SetStream("", strings.Split(content, " "))
				corpus.Upsert(filename, document)
			}

			bm := rank.NewBM25F(1.2, map[rank.Field]rank.FieldConfig{
				"": {Weight: 1.0, B: 0.72},
			})
			query := strings.Split(tt.query, " ")
			results := extractNames(bm.Rank(corpus, query))

			if !slices.Equal(results, results) {
				t.Fatalf("Rank: expected %v, got %v", tt.want, results)
			}
		})
	}
}

func extractNames(src []rank.Result) []string {
	results := make([]string, len(src))
	for i, result := range src {
		results[i] = result.Name
	}
	return results
}
