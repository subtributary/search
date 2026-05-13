package tokenize_test

import (
	"iter"
	"slices"
	"strings"
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

type echoTokenizer struct {
	text string
}

func (t echoTokenizer) Tokens(_ string) iter.Seq[string] {
	return slices.Values([]string{t.text})
}

func echo(text string) echoTokenizer {
	return echoTokenizer{text}
}

func uppercase(s string) string {
	return strings.ToUpper(s)
}

func truncate(n int) func(string) string {
	return func(s string) string {
		if len(s) < n {
			return s
		}
		return s[:n]
	}
}

func TestSmartTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		text      string
		tokenizer tokenize.SmartTokenizer
		want      []string
	}{
		{
			name: "empty",
			text: "",
			tokenizer: tokenize.NewSmartTokenizer(
				tokenize.WithScriptTokenizer("Latin", echo("hello")),
			),
			want: []string{},
		},
		{
			name:      "unconfigured",
			text:      "hello world",
			tokenizer: tokenize.NewSmartTokenizer(),
			want:      []string{},
		},
		{
			name:      "empty and unconfigured tokens",
			text:      "",
			tokenizer: tokenize.NewSmartTokenizer(),
			want:      []string{},
		},
		{
			name: "single word",
			text: "hello",
			tokenizer: tokenize.NewSmartTokenizer(
				tokenize.WithScriptTokenizer("Latin", echo("hello")),
			),
			want: []string{"hello"},
		},
		{
			name: "single script",
			text: "hello world",
			tokenizer: tokenize.NewSmartTokenizer(
				tokenize.WithScriptTokenizer("Latin", echo("hello")),
			),
			want: []string{"hello"},
		},
		{
			name: "two scripts",
			text: "world 안녕",
			tokenizer: tokenize.NewSmartTokenizer(
				tokenize.WithScriptTokenizer("Latin", echo("hello")),
				tokenize.WithScriptTokenizer("Hangul", echo("world")),
			),
			want: []string{"hello", "world"},
		},
		{
			name: "single normalizer",
			text: "hello world",
			tokenizer: tokenize.NewSmartTokenizer(
				tokenize.WithScriptTokenizer("Latin", echo("hello")),
				tokenize.WithScriptNormalizer("Latin", uppercase),
			),
			want: []string{"HELLO"},
		},
		{
			name: "double normalizer",
			text: "hello world",
			tokenizer: tokenize.NewSmartTokenizer(
				tokenize.WithScriptTokenizer("Latin", echo("hello")),
				tokenize.WithScriptNormalizer("Latin", uppercase),
				tokenize.WithScriptNormalizer("Latin", truncate(3)),
			),
			want: []string{"HEL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := slices.Collect(tt.tokenizer.Tokens(tt.text))
			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens: got %v, want %v", got, tt.want)
			}
		})
	}
}
