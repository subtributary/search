package main

import (
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"os"

	"github.com/subtributary/search"
)

func loadIndex(filename string) (idx *search.Index) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}
	if err := json.Unmarshal(data, &idx); err != nil {
		log.Fatalf("unmarshal: %v", err)
	}
	return
}

func printResults(results iter.Seq[search.Result], count int) {
	for result := range results {
		id := result.Id
		baseName := result.Attachments["baseName"]
		fmt.Printf("baseName: %-*s id: %s\n", 24, baseName, id)

		count--
		if count == 0 {
			break
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <index-file> <query>\n", os.Args[0])
	}
	indexFile := os.Args[1]
	query := os.Args[2]

	idx := loadIndex(indexFile)
	results := idx.Search(query)
	printResults(results, 10)
}
