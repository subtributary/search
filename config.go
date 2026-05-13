package search

type Config struct {
	Version int                     `json:"version"`
	Scripts map[string]ScriptConfig `json:"scripts"`
}

type ScriptConfig struct {
	Tokenizer   string   `json:"tokenizer"`
	Normalizers []string `json:"normalizers"`
}
