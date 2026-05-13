package tokenize_test

import (
	"slices"
	"testing"

	"github.com/subtributary/search/internal/tokenize"
)

func TestUAX29(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "empty",
			text: "",
			want: []string{},
		},
		{
			name: "single",
			text: "single",
			want: []string{"single"},
		},
		{
			name: "multiple words",
			text: "multiple words",
			want: []string{"multiple", "words"},
		},
		{
			name: "letters and numbers",
			text: "88 wb9 10wb o_o",
			want: []string{"88", "wb9", "10wb", "o_o"},
		},
		{
			name: "leading non-word characters",
			text: " -A",
			want: []string{"A"},
		},
		{
			name: "basic punctuation",
			text: "A~!@#$%^&*()+-=/.,?><\\|]}[{B",
			want: []string{"A", "B"},
		},
		{
			name: "emojis",
			text: "melt \U0001fae0 zwj\U0001f468\u200d\U0001f9b0nospace",
			want: []string{"melt", "zwj", "nospace"},
		},
		{
			name: "format and extend characters",
			text: "눈\u302f ac\u0327ai hy\u00adphen z\u200bw",
			want: []string{"눈\u302f", "ac\u0327ai", "hyphen", "zw"},
		},
		{
			name: "spaces that do not break",
			text: "A\u00a0B",
			want: []string{"AB"},
		},
		{
			name: "joiners",
			text: "A\u200dB",
			want: []string{"AB"},
		},
		{
			name: "ambiguous symbols that do not break",
			text: "A\u00b7B C\u0027\u059eD",
			want: []string{"AB", "C\u0027\u059eD"},
		},
		{
			// This is a simplification of UAX #29's WB4 rule for our use case.
			name: "Mc | Me | Mn | FF9E | FF9F | (Cf except ZWSP)",
			text: "A\u302f\u0488\u0327\uff9e\uff9f\u08e2B",
			want: []string{"A\u302f\u0327B"},
		},
		{
			name: "uax29 extra letters",
			text: "A\u00b8\u02c2\u02c3\u02c3\u02c4\u02c5" +
				"\u02d2\u02d3\u02d4\u02d5\u02d6\u02d7" +
				"\u02de\u02df\u02e5\u02e6\u02e7\u02e8\u02e9\u02ea\u02eb\u02ed" +
				"\u02ef\u02f0\u02f1\u02f2\u02f3\u02f4\u02f5\u02f6\u02f7" +
				"\u02f8\u02f9\u02fa\u02fb\u02fc\u02fd\u02fe\u02ff" +
				"\u055a\u055b\u055c\u055e\u058a\u05f3\u070f" +
				"\ua708\ua709\ua70a\ua70b\ua70c\ua70d\ua70e\ua70f" +
				"\ua710\ua711\ua712\ua713\ua714\ua715\ua716" +
				"\ua720\ua721\ua789\ua78a\uab5bB",
			want: []string{
				"A\u00b8\u02c2\u02c3\u02c3\u02c4\u02c5" +
					"\u02d2\u02d3\u02d4\u02d5\u02d6\u02d7" +
					"\u02de\u02df\u02e5\u02e6\u02e7\u02e8\u02e9\u02ea\u02eb\u02ed" +
					"\u02ef\u02f0\u02f1\u02f2\u02f3\u02f4\u02f5\u02f6\u02f7" +
					"\u02f8\u02f9\u02fa\u02fb\u02fc\u02fd\u02fe\u02ff" +
					"\u055a\u055b\u055c\u055e\u058a\u05f3" +
					"\ua708\ua709\ua70a\ua70b\ua70c\ua70d\ua70e\ua70f" +
					"\ua710\ua711\ua712\ua713\ua714\ua715\ua716" +
					"\ua720\ua721\ua789\ua78a\uab5bB",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokenize.NewUAX29()
			tokens := slices.Collect(subject.Tokens(tt.text))

			if !slices.Equal(tokens, tt.want) {
				tr := []rune(tokens[0])
				t.Errorf("got %#x", tr)
				t.Errorf("Tokens: want %#v, got %#v", tt.want, tokens)
			}
		})
	}
}
