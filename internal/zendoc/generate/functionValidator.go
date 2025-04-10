package generate

import "unicode"

func IsPrivateFunction(name string) bool {
	if name == "" {
		return false
	}
	return !unicode.IsUpper(rune(name[0]))
}
