package tokenize_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

// echoTokenizer is a simple "tokenizer" for testing purposes.
// It simply outputs its preset tokens in place of any text it parses.
type echoTokenizer struct {
	text string
}

func newEchoTokenizer(text string) echoTokenizer {
	return echoTokenizer{text}
}

func (t echoTokenizer) Tokens(_ string) iter.Seq[string] {
	return slices.Values([]string{t.text})
}

func TestSmartTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		text       string
		tokenizers map[string]tokenize.Tokenizer
		want       []string
	}{
		{
			name: "empty",
			text: "",
			tokenizers: map[string]tokenize.Tokenizer{
				"Latin": newEchoTokenizer("hello"),
			},
			want: []string{},
		},
		{
			name:       "unconfigured",
			text:       "hello world",
			tokenizers: map[string]tokenize.Tokenizer{},
			want:       []string{},
		},
		{
			name:       "empty and unconfigured tokens",
			text:       "",
			tokenizers: map[string]tokenize.Tokenizer{},
			want:       []string{},
		},
		{
			name: "single word",
			text: "hello",
			tokenizers: map[string]tokenize.Tokenizer{
				"Latin": newEchoTokenizer("hello"),
			},
			want: []string{"hello"},
		},
		{
			name: "single script",
			text: "hello world",
			tokenizers: map[string]tokenize.Tokenizer{
				"Latin": newEchoTokenizer("hello"),
			},
			want: []string{"hello"},
		},
		{
			name: "two scripts",
			text: "world 안녕",
			tokenizers: map[string]tokenize.Tokenizer{
				"Latin":  newEchoTokenizer("hello"),
				"Hangul": newEchoTokenizer("world"),
			},
			want: []string{"hello", "world"},
		},
		{
			name: "leading common",
			text: "--hello world",
			tokenizers: map[string]tokenize.Tokenizer{
				"Latin": newEchoTokenizer("hello"),
			},
			want: []string{"hello"},
		},
		{
			name: "intermixed common",
			text: "hello! 안녕!",
			tokenizers: map[string]tokenize.Tokenizer{
				"Latin":  newEchoTokenizer("hello"),
				"Hangul": newEchoTokenizer("world"),
			},
			want: []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokenize.NewSmartTokenizer()
			for script, tokenizer := range tt.tokenizers {
				subject.SetSubTokenizer(script, tokenizer)
			}

			// We only want the text. If it's right, the script is too.
			tokens := subject.Tokens(tt.text)
			var got []string
			for _, text := range tokens {
				got = append(got, text)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens: got %v, want %v", got, tt.want)
			}
		})
	}
}
