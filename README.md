# Search

A BM25F search implementation with per-script tokenizers and normalizers.

Our search implementation is built to handle documents in multiple languages.
The per-script customization lets us choose how best to handle each one.

This package is still in initial development, and there may be breaking changes between versions.
Version 1.0 will be released after its successful full integration into our
[Musings](https://github.com/subtributary/musings) project.

## Features

* Multiple fields per document that can be configured and searched
* Attachment strings that are returned alongside document ids in search results
* Separate normalizer and tokenizer configurations for every Unicode script
* Sane defaults for every script
* Supported normalizers:
    * `LowerCase`, `NFC`, `NFD`, `NFKC`,  `NFKD`
* Supported tokenizers:
    * `UAX29`, `Unigram`, `Bigram`, `Trigram`

## Quick start

```
go get "github.com/subtributary/search"
```

To create an index:

```go
index, err := search.NewIndex(
    search.WithField("title", 5.0),
    search.WithField("content", 1.0),
)
```

To add a document to the index:

```go
err := index.Upsert(unique_id, map[string]string{
    "title":   title,
    "content": content,
})
```

Configured fields are parsed then thrown away.
To attach data to be returned with results:

```go
err := index.Upsert(unique_id, map[string]string{
    "title":         title,
    "content":       content,
	// display_title won't be parsed but will be returned with results.
	"display_title": title
})
```

To search through the corpus and print the results:

```go
results := index.Search(query)
for result := range results {
	id := result.Id
	title := result.Attachments["display_title"]
	fmt.Printf("id: %-*s title: %s", 15, id, title)
}
```

To save the index to a file:

```go
data, err := json.Marshal(index)
err := os.WriteFile(filename, data, fs.ModePerm)
```

To load the index from file:

```go
data, err := os.ReadFile(filename)
var index search.Index
err := json.Unmarshal(data, &index)
```
