package search_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/subtributary/search"
)

var documents = map[string]map[string]string{
	"empty": {
		"title":       "",
		"body":        "",
		"description": "empty document",
	},
	"biography": {
		"title":       "EXID",
		"body":        "EXID is a South Korean girl group formed in 2012. The group consists of Solji (솔지), ELLY, Hani (하니), Hyelin (혜린), and Jeonghwa (정화).",
		"description": "summary of the kpop group EXID",
	},
	"favorites": {
		"title":       "Favorites",
		"body":        "My favorite flavor is vanilla. My favorite color is white. My favorite singer is Solji (솔지). My favorite kpop group is EXID.",
		"description": "my personal favorites",
	},
	"song": {
		"title":       "We Are",
		"body":        "우리를 우리가 될 수 있도록 / 만들어줘서 고마운 마음 뿐이야 내겐 / 앞으로도 소중히 기억할게 / 늘 언제 어디에 있던지 우리",
		"description": "a song that is a promise to EXID's fans",
	},
}

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
}
