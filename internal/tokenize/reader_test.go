package tokenize_test

import (
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

func TestCharReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		runes   []rune
		scripts []tokenize.Script
	}{
		{
			name:    "empty",
			text:    "",
			runes:   []rune{},
			scripts: []tokenize.Script{},
		},
		{
			name:    "single",
			text:    "abc",
			runes:   []rune{'a', 'b', 'c'},
			scripts: []tokenize.Script{"Latin", "Latin", "Latin"},
		},
		{
			name:    "multi",
			text:    "a술晴",
			runes:   []rune{'a', '술', '晴'},
			scripts: []tokenize.Script{"Latin", "Hangul", "Han"},
		},
		{
			name:    "modifiers",
			text:    "c\u00b8",
			runes:   []rune{'c', '\u00b8'},
			scripts: []tokenize.Script{"Latin", ""},
		},
		{
			name:    "common",
			text:    "! ?",
			runes:   []rune{'!', ' ', '?'},
			scripts: []tokenize.Script{"", "", ""},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			subject := tokenize.NewCharReader(tt.text)
			for i := 0; i != len(tt.runes); i++ {
				subject.Read()
				if subject.Current != tt.runes[i] {
					t.Errorf("Runes[%d]: got %q, want %q", i, subject.Current, tt.runes[i])
				}
				if subject.Script != tt.scripts[i] {
					t.Errorf("Scripts[%d]: got %q, want %q", i, subject.Script, tt.scripts[i])
				}
				if t.Failed() {
					t.FailNow()
				}
			}
		})
	}
}
