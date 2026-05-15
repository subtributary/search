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

