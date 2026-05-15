package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/subtributary/search"
)

func createIndex() *search.Index {
	idx, configErr := search.NewIndex(
		search.WithField("name", 2.0),
		search.WithField("content", 1.0),
	)
	if configErr != nil {
		log.Fatalf("new index: %v", configErr)
	}
	return idx
}

func populateIndex(idx *search.Index, root string) {
	// This is a quick and dirty directory walk that isn't prod-worthy.
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, _ error) error {
		if info, err := d.Info(); err != nil {
			log.Fatalf("file info: %v", err)
		} else if !info.Mode().IsRegular() || d.IsDir() {
			return nil
		}

		id, _ := filepath.Rel(root, path)
		name := filepath.Base(path)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read file: %v", err)
		}

		// Add the file data to the index.
		userError := idx.Upsert(id, map[string]string{
			"name":    name,
			"content": string(content),
			// baseName is not configured so will be attached to results.
			"baseName": name,
		})
		if userError != nil {
			log.Fatalf("upsert: %v", userError)
		}

		return nil
	})
}

func saveIndex(idx *search.Index, filename string) {
	if data, err := json.Marshal(idx); err != nil {
		log.Fatalf("marshal: %v", err)
	} else if err := os.WriteFile(filename, data, fs.ModePerm); err != nil {
		log.Fatalf("save: %v", err)
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <directory> <output-file>\n", os.Args[0])
	}
	root := os.Args[1]
	output := os.Args[2]

	idx := createIndex()
	populateIndex(idx, root)
	saveIndex(idx, output)
}
