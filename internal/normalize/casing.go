package normalize

import "strings"

// CasingNormalizer changes everything to lowercase.
// It's wrong sometimes, such as for the Turkish 'I', but it's consistent.
type CasingNormalizer struct{}

func NewCasingNormalizer() CasingNormalizer {
	return CasingNormalizer{}
}

func (CasingNormalizer) Normalize(text string) string {
	return strings.ToLower(text)
}
