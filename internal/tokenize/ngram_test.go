package tokenize_test

import (
	"slices"
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

func TestNGram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		minN, maxN int
		text       string
		want       []string
	}{
		{
			name: "empty",
			minN: 1,
			maxN: 1,
			text: "",
			want: []string{},
		},
		{
			name: "unigram",
			minN: 1,
			maxN: 1,
			text: "abcd",
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "bigram",
			minN: 2,
			maxN: 2,
			text: "abcd",
			want: []string{"ab", "bc", "cd"},
		},
		{
			name: "trigram",
			minN: 3,
			maxN: 3,
			text: "abcd",
			want: []string{"abc", "bcd"},
		},
		{
			name: "bigram and trigram",
			minN: 2,
			maxN: 3,
			text: "abcd",
			want: []string{"ab", "bc", "cd", "abc", "bcd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokenize.NewNGram(tt.minN, tt.maxN)
			tokens := slices.Collect(subject.Tokens(tt.text))

			if !slices.Equal(tokens, tt.want) {
				t.Errorf("Tokens: want %#v, got %#v", tt.want, tokens)
			}
		})
	}
}
