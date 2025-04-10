package generate

import "unicode"

/*
@description Determine if a function is private based on its name (starts with a lowercase letter)
@param name string - The name of the function to check
@return bool - true if the function is private, false otherwise
@example IsPrivateFunction("myFunc") => true
*/

func IsPrivateFunction(name string) bool {
	if name == "" {
		return false
	}
	return !unicode.IsUpper(rune(name[0]))
}
