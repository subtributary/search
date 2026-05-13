package search

import "github.com/subtributary/search/internal/shared"

type Script shared.Script

type Config struct {
	Version int                     `json:"version"`
	Scripts map[Script]ScriptConfig `json:"scripts"`
}

type ScriptConfig struct {
	Tokenizer   string   `json:"tokenizer"`
	Normalizers []string `json:"normalizers"`
}
