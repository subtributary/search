package tokenize_test

import (
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

func TestScriptTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		text     string
		segments []string
		scripts  []tokenize.Script
	}{
		{
			name:     "empty",
			text:     "",
			segments: []string{},
			scripts:  []tokenize.Script{},
		},
		{
			name:     "single",
			text:     "abc",
			segments: []string{"abc"},
			scripts:  []tokenize.Script{"Latin"},
		},
		{
			name:     "multi",
			text:     "a술晴",
			segments: []string{"a", "술", "晴"},
			scripts:  []tokenize.Script{"Latin", "Hangul", "Han"},
		},
		{
			name:     "modifiers",
			text:     "c\u00b8",
			segments: []string{"c\u00b8"},
			scripts:  []tokenize.Script{"Latin"},
		},
		{
			name:     "common",
			text:     "! ?",
			segments: []string{},
			scripts:  []tokenize.Script{},
		},
		{
			name:     "leading common",
			text:     "!abc",
			segments: []string{"abc"},
			scripts:  []tokenize.Script{"Latin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokenize.NewScriptTokenizer()
			tokens := subject.Tokens(tt.text)

			i := 0
			for script, text := range tokens {
				if i >= len(tt.segments) {
					t.Errorf("len(tokens): got %d, want %d", i, len(tt.segments))
				}
				if text != tt.segments[i] {
					t.Errorf("Text: got %q, want %q", text, tt.segments[i])
				}
				if script != tt.scripts[i] {
					t.Errorf("Script: got %q, want %q", script, tt.scripts[i])
				}
				i++
			}

			if t.Failed() {
				t.FailNow()
			}
		})
	}
}
