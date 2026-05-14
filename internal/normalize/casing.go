package normalize

import "strings"

// CaseNormalizer changes everything to lowercase.
// It's wrong sometimes, such as for the Turkish 'I', but it's consistent.
type CaseNormalizer struct{}

func NewCaseNormalizer() CaseNormalizer {
	return CaseNormalizer{}
}

func (CaseNormalizer) Normalize(text string) string {
	return strings.ToLower(text)
}
