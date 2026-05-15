package rank_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/subtributary/search/internal/rank"
)

func TestBM25F_Rank(t *testing.T) {
	t.Parallel()

	config := map[rank.Field]rank.FieldConfig{
		"title": {Weight: 2.0, B: 0},
		"body":  {Weight: 1.0, B: 0},
	}

	emptyDoc := rank.NewDocument()
	emptyDoc.SetStream("title", []string{})
	emptyDoc.SetStream("body", []string{})

	helloDoc := rank.NewDocument()
	helloDoc.SetStream("title", []string{"hello"})
	helloDoc.SetStream("body", []string{"hello", "blue", "world"})

	natureDoc := rank.NewDocument()
	natureDoc.SetStream("title", []string{"nature"})
	natureDoc.SetStream("body", []string{"blue", "tulip", "blue", "sky", "world"})

	tests := []struct {
		name      string
		documents map[string]rank.Document
		query     string
		want      []rank.Result
	}{
		{
			name:      "no documents",
			documents: map[string]rank.Document{},
			query:     "test",
			want:      []rank.Result{},
		},
		{
			name: "empty fields",
			documents: map[string]rank.Document{
				"empty2": emptyDoc,
				"empty1": emptyDoc,
			},
			query: "test",
			want: []rank.Result{
				// Both are empty and score the same,
				// so they are sorted alphabetically by title.
				{Id: "empty1", Document: emptyDoc},
				{Id: "empty2", Document: emptyDoc},
			},
		},
		{
			name: "single match",
			documents: map[string]rank.Document{
				"empty":  emptyDoc,
				"nature": natureDoc,
			},
			query: "tulip",
			want: []rank.Result{
				// Only natureDoc has the word "tulip", so it is listed first.
				{Id: "nature", Document: natureDoc},
				{Id: "empty", Document: emptyDoc},
			},
		},
		{
			name: "multiple matches",
			documents: map[string]rank.Document{
				"empty":  emptyDoc,
				"nature": natureDoc,
				"hello":  helloDoc,
			},
			query: "world",
			want: []rank.Result{
				// helloDoc and natureDoc both have one "word" in the body,
				// so they are sorted alphabetically by title.
				{Id: "hello", Document: helloDoc},
				{Id: "nature", Document: natureDoc},
				{Id: "empty", Document: emptyDoc},
			},
		},
		{
			name: "overused word",
			documents: map[string]rank.Document{
				"empty":  emptyDoc,
				"hello":  helloDoc,
				"nature": natureDoc,
			},
			query: "blue",
			want: []rank.Result{
				// natureDoc and helloDoc both contain the word "blue",
				// but the word is more frequent in natureDoc.
				{Id: "nature", Document: natureDoc},
				{Id: "hello", Document: helloDoc},
				{Id: "empty", Document: emptyDoc},
			},
		},
		{
			name: "multiword query",
			documents: map[string]rank.Document{
				"empty":  emptyDoc,
				"hello":  helloDoc,
				"nature": natureDoc,
			},
			query: "hello blue",
			want: []rank.Result{
				// natureDoc and helloDoc both contain the world "blue",
				// but helloDoc also contains "hello" in its title.
				{Id: "hello", Document: helloDoc},
				{Id: "nature", Document: natureDoc},
				{Id: "empty", Document: emptyDoc},
			},
		},
	}

	extractNames := func(src []rank.Result) []string {
		results := make([]string, len(src))
		for i, result := range src {
			results[i] = result.Id
		}
		return results
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			corpus := rank.NewCorpus()
			for filename, document := range tt.documents {
				corpus.Upsert(filename, document)
			}

			bm := rank.NewBM25F(1.2, config)
			results := bm.Rank(corpus, strings.Split(tt.query, " "))

			gotNames := extractNames(results)
			wantNames := extractNames(tt.want)
			if !slices.Equal(gotNames, wantNames) {
				t.Errorf("Rank: got %v, want %v", gotNames, wantNames)
			}
		})
	}
}
