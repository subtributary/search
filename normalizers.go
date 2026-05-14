package search

import (
	"fmt"

	"github.com/subtributary/search/internal/normalize"
	"golang.org/x/text/unicode/norm"
)

type Normalizer string

// I don't see a reason to use anything other than `text -> Lowercase -> NFKC`,
// but I do need this configurable anyway. If I change it in a future version,
// we will need to still be able to load old configurations. I may as well make
// it configurable from the beginning.
var (
	// Lowercase normalizes text to lowercase.
	//
	// This is okay even for Turkish and Azerbaijani because the inaccuracies
	// are consistent and do not change word meanings for search indexing.
	Lowercase Normalizer = "lowercase"

	// NFC converts characters into their canonical composed forms.
	// Character meanings and appearances are preserved.
	NFC Normalizer = "nfc"

	// NFD converts characters into their canonical decomposed forms.
	// Character meanings and appearances are preserved.
	NFD Normalizer = "nfd"

	// NFKC converts characters into their compatibility composed forms.
	// Similar characters are represented by the same code point, which is a
	// lossy conversion, but it enables fuzzy searching.
	NFKC Normalizer = "nfkc"

	// NFKD converts characters into their compatibility decomposed forms.
	// Similar characters are represented by the same code point, which is a
	// lossy conversion, but it enables fuzzy searching.
	NFKD Normalizer = "nfkd"
)

func (n Normalizer) toInternal() (normalize.Normalizer, error) {
	switch n {
	case "lowercase":
		return normalize.NewCaseNormalizer(), nil
	case "nfc":
		return normalize.NewUnicodeNormalizer(norm.NFC), nil
	case "nfd":
		return normalize.NewUnicodeNormalizer(norm.NFD), nil
	case "nfkc":
		return normalize.NewUnicodeNormalizer(norm.NFKC), nil
	case "nfkd":
		return normalize.NewUnicodeNormalizer(norm.NFKD), nil
	default:
		return nil, fmt.Errorf("invalid normalizer: %s", n)
	}
}
