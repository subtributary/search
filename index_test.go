package search_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/subtributary/search"
)

// todo: create a map[version]jsonresult
// todo: create standard documents

func TestIndex_Save(t *testing.T) {
	t.Parallel()

	index, err := search.NewIndex()
	if err != nil {
		t.Fatalf("new index: %v", err)
	}

	data, err := json.Marshal(index)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// todo: finish the test
	fmt.Printf("%v\n", string(data))
}
