package tokenize_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/subtributary/search/internal/shared"
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
				map[shared.Script]tokenize.Tokenizer{
					"Latin": echo("hello"),
				},
			),
			want: []string{},
		},
		{
			name: "unconfigured",
			text: "hello world",
			tokenizer: tokenize.NewSmartTokenizer(
				map[shared.Script]tokenize.Tokenizer{},
			),
			want: []string{},
		},
		{
			name: "empty and unconfigured tokens",
			text: "",
			tokenizer: tokenize.NewSmartTokenizer(
				map[shared.Script]tokenize.Tokenizer{},
			),
			want: []string{},
		},
		{
			name: "single word",
			text: "hello",
			tokenizer: tokenize.NewSmartTokenizer(
				map[shared.Script]tokenize.Tokenizer{
					"Latin": echo("hello"),
				},
			),
			want: []string{"hello"},
		},
		{
			name: "single script",
			text: "hello world",
			tokenizer: tokenize.NewSmartTokenizer(
				map[shared.Script]tokenize.Tokenizer{
					"Latin": echo("hello"),
				},
			),
			want: []string{"hello"},
		},
		{
			name: "two scripts",
			text: "world 안녕",
			tokenizer: tokenize.NewSmartTokenizer(
				map[shared.Script]tokenize.Tokenizer{
					"Latin":  echo("hello"),
					"Hangul": echo("world"),
				},
			),
			want: []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokens := tt.tokenizer.Tokens(tt.text)

			// We only want the text. If it's right, the script is too.
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
