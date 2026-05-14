package normalize

import (
	"fmt"

	"golang.org/x/text/unicode/norm"
)

// Normalizer types implement a normalization algorithm that can be delegated
// to by SmartNormalizer depending on script configuration
type Normalizer interface {
	Normalize(text string) string
}

type SmartNormalizer struct {
	subs map[string][]Normalizer
}

func NewSmartNormalizer() SmartNormalizer {
	return SmartNormalizer{subs: make(map[string][]Normalizer)}
}

func (n *SmartNormalizer) Normalize(script string, text string) string {
	if subs, ok := n.subs[script]; ok {
		for _, sub := range subs {
			text = sub.Normalize(text)
		}
	}
	return text
}

// SetSubNormalizers sets the normalizers to use for a specific script.
// If any of the ids is invalid, an error is returned.
func (n *SmartNormalizer) SetSubNormalizers(script string, ids []string) error {
	subs := make([]Normalizer, 0, len(ids))
	for _, id := range ids {
		switch id {
		case "lowercase":
			subs = append(subs, NewCaseNormalizer())
		case "nfc":
			subs = append(subs, NewUnicodeNormalizer(norm.NFC))
		case "nfd":
			subs = append(subs, NewUnicodeNormalizer(norm.NFD))
		case "nfkc":
			subs = append(subs, NewUnicodeNormalizer(norm.NFKC))
		case "nfkd":
			subs = append(subs, NewUnicodeNormalizer(norm.NFKD))
		default:
			return fmt.Errorf("invalid normalizer id: %s", id)
		}
	}
	return nil
}
