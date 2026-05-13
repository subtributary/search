package shared

import "unicode"

// Script is the name of a Unicode script or empty for the Common script.
type Script string

// DetectScript returns the Unicode script of the character.
// If the character has no script or is in the Common script, "" is returned.
func DetectScript(c rune) Script {
	if !unicode.Is(unicode.Common, c) {
		for name, rt := range unicode.Scripts {
			if unicode.Is(rt, c) {
				return Script(name)
			}
		}
	}
	return ""
}
