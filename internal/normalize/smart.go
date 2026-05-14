package normalize

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
func (n *SmartNormalizer) SetSubNormalizers(script string, norms []Normalizer) {
	n.subs[script] = norms
}
