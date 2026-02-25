package utf16

import "unicode/utf16"

func Index(s string, target rune) (int, bool) {
	runes := []rune(s)

	for i, r := range runes {
		if r == target {
			return len(utf16.Encode(runes[:i])), true
		}
	}
	return 0, false
}
