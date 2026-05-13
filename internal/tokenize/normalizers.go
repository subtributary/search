package tokenize

// Normalizer types implement a normalization algorithm that can be delegated
// to by SmartTokenizer depending on script configuration.
type Normalizer func(string) string

// nfkd
// uppercase
// dictionary split?
// small words limit
// large words limit
