package search_test

import (
	"slices"
	"testing"

	"github.com/subtributary/search"
)

/*
func TestIndex(t *testing.T) {
	t.Parallel()

	index, err := search.NewIndex(
		search.WithField("title", 5),
		search.WithField("body", 1),
	)
	if err != nil {
		t.Fatalf("new index: %v", err)
	}

	for name, doc := range documents {
		index.Upsert(name, doc)
	}

	// todo: perform some searches

	saved, err := json.Marshal(index)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}

	// todo: load from saved
	log.Printf("saved: %s", saved)

	// todo: perform searches
}*/

func TestIndex_Search(t *testing.T) {
	t.Parallel()

	index, err := search.NewIndex(
		search.WithField("title", 2),
		search.WithField("body", 1),
	)
	if err != nil {
		t.Fatalf("new index: %v", err)
	}

	_ = index.Upsert("empty", map[string]string{
		"title": "",
		"body":  "",
	})
	_ = index.Upsert("EXID", map[string]string{
		"title": "EXID: 솔지, ELLY, 하니, 혜린, and 정화",
		"body":  "EXID is a South Korean girl group formed in 2012. The group consists of Solji (솔지), ELLY, Hani (하니), Hyelin (혜린), and Jeonghwa (정화).",
	})
	_ = index.Upsert("favorites", map[string]string{
		"title": "Favorites",
		"body":  "My favorite flavor is vanilla. My favorite color is white. My favorite singer is Solji (솔지). My favorite kpop group is EXID.",
	})
	_ = index.Upsert("song", map[string]string{
		"title": "We Are",
		"body":  "우리를 우리가 될 수 있도록 / 만들어줘서 고마운 마음 뿐이야 내겐 / 앞으로도 소중히 기억할게 / 늘 언제 어디에 있던지 우리",
	})

	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{
			name:  "empty",
			query: "",
			// By default, everything is sorted alphabetically.
			want: []string{"EXID", "empty", "favorites", "song"},
		},
		{
			name:  "single term",
			query: "솔지",
			want:  []string{"EXID", "favorites", "empty", "song"},
		},
		{
			name:  "multiple terms",
			query: "우리 favorite",
			want:  []string{"favorites", "song", "EXID", "empty"},
		},
	}

	extractIds := func(results []search.Result) []string {
		ids := make([]string, len(results))
		for i, result := range results {
			ids[i] = result.Id
		}
		return ids
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			results := slices.Collect(index.Search(tt.query))
			ids := extractIds(results)

			if !slices.Equal(ids, tt.want) {
				t.Errorf("want %v, got %v", tt.want, ids)
			}
		})
	}
}

func TestIndex_Upsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		content   map[string]string
		wantError bool
	}{
		{
			name:      "no fields",
			content:   map[string]string{},
			wantError: true,
		},
		{
			name:      "only fields",
			content:   map[string]string{"title": ""},
			wantError: false,
		},
		{
			name:      "only attachments",
			content:   map[string]string{"attached": ""},
			wantError: true,
		},
		{
			name:      "fields and attachments",
			content:   map[string]string{"title": "", "attached": ""},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			index, err := search.NewIndex(
				search.WithField("title", 1),
			)
			if err != nil {
				t.Fatalf("new index: %v", err)
			}

			err = index.Upsert(tt.name, tt.content)
			if tt.wantError && err == nil {
				t.Errorf("want error, got none")
			} else if !tt.wantError && err != nil {
				t.Errorf("want no error, got %v", err)
			}
		})
	}
}
