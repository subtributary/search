package rank_test

import (
	"maps"
	"strings"
	"testing"

	"github.com/subtributary/search/internal/rank"
)

func TestStream(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tokens     []string
		wantLength int
		wantCounts map[string]int
	}{
		{
			name:       "empty",
			tokens:     []string{},
			wantLength: 0,
			wantCounts: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := rank.NewStream(tt.tokens)
			if s.Length != tt.wantLength {
				t.Errorf("Length: got %d, want %d", s.Length, tt.wantLength)
			}
			if !maps.Equal(s.TermCounts, tt.wantCounts) {
				t.Errorf("TermCounts: got %v, want %v", s.TermCounts, tt.wantCounts)
			}
		})
	}
}

func TestCorpus(t *testing.T) {
	t.Parallel()
	corpus := rank.NewCorpus()

	assertAvgDocSize := func(want float64) {
		t.Helper()
		if got := corpus.AvgStreamLength(""); got != want {
			t.Errorf("AvgStreamLength: got %v, want %v", got, want)
		}
	}

	assertCount := func(term string, want int) {
		t.Helper()
		if got := corpus.DocsWithTerm[term]; got != want {
			t.Errorf("DocsWithTerm[%q]: got %d, want %d", term, got, want)
		}
	}

	assertSize := func(want int) {
		t.Helper()
		if got := len(corpus.Documents); got != want {
			t.Errorf("len(Documents): got %d, want %d", got, want)
		}
	}

	createDocument := func(text string) rank.Document {
		d := rank.NewDocument()
		d.SetStream("", strings.Split(text, " "))
		return d
	}

	// Populate corpus
	corpus.Upsert("one", createDocument("one"))
	corpus.Upsert("two", createDocument("two two"))
	corpus.Upsert("three", createDocument("three three three"))
	assertAvgDocSize(2)
	assertCount("three", 1)
	assertSize(3)
	if t.Failed() {
		t.FailNow()
	}

	// Replace existing document
	corpus.Upsert("one", createDocument("one two three four"))
	assertAvgDocSize(3)
	assertCount("three", 2)
	assertSize(3)
	if t.Failed() {
		t.FailNow()
	}

	// Remove document
	corpus.Remove("one")
	assertAvgDocSize(2.5)
	assertCount("one", 0)
	assertSize(2)
	if t.Failed() {
		t.FailNow()
	}

	// Remove remaining documents
	corpus.Remove("two")
	corpus.Remove("three")
	assertAvgDocSize(0)
	assertCount("two", 0)
	assertSize(0)
	if t.Failed() {
		t.FailNow()
	}

	// Remove nonexistent
	corpus.Remove("missing")
	assertAvgDocSize(0.0)
	assertSize(0)
}
