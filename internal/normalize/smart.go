package normalize

import "github.com/subtributary/search/internal/shared"

// Normalizer types implement a normalization algorithm that can be delegated
// to by SmartNormalizer depending on script configuration
type Normalizer interface {
	Normalize(text string) string
}

type SmartNormalizer struct {
	subs map[shared.Script]Normalizer
}

func NewSmartNormalizer(subs map[shared.Script]Normalizer) SmartNormalizer {
	return SmartNormalizer{subs}
}

func (n SmartNormalizer) Normalize(script shared.Script, text string) string {
	if norm, ok := n.subs[script]; ok {
		return norm.Normalize(text)
	}
	return text
}
