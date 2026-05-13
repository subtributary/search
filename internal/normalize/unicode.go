package normalize

import "golang.org/x/text/unicode/norm"

// UnicodeNormalizer implements the Unicode equivalence algorithms.
type UnicodeNormalizer struct {
	form norm.Form
}

func NewUnicodeNormalizer(form norm.Form) UnicodeNormalizer {
	return UnicodeNormalizer{form}
}

func (n UnicodeNormalizer) Normalize(text string) string {
	return n.form.String(text)
}
