package tokenize_test

import (
	"slices"
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

func TestScriptTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		text     string
		segments []string
		scripts  []string
	}{
		{
			name:     "empty",
			text:     "",
			segments: []string{},
			scripts:  []string{},
		},
		{
			name:     "single",
			text:     "abc",
			segments: []string{"abc"},
			scripts:  []string{"Latin"},
		},
		{
			name:     "multi",
			text:     "a술晴",
			segments: []string{"a", "술", "晴"},
			scripts:  []string{"Latin", "Hangul", "Han"},
		},
		{
			name:     "modifiers",
			text:     "c\u00b8",
			segments: []string{"c\u00b8"},
			scripts:  []string{"Latin"},
		},
		{
			name:     "common",
			text:     "! ?",
			segments: []string{},
			scripts:  []string{},
		},
		{
			name:     "leading common",
			text:     "!abc",
			segments: []string{"abc"},
			scripts:  []string{"Latin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokenize.NewScriptTokenizer()
			tokens := subject.Tokens(tt.text)

			var scripts, segments []string
			for k, v := range tokens {
				scripts = append(scripts, k)
				segments = append(segments, v)
			}

			if !slices.Equal(tt.scripts, scripts) {
				t.Errorf("scripts: got %v, want %v", scripts, tt.scripts)
			}
			if !slices.Equal(tt.segments, segments) {
				t.Errorf("segments: got %v, want %v", segments, tt.segments)
			}
		})
	}
}
