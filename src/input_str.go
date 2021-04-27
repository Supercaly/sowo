package src

import "strings"

// Represent an input string with methods to
// mange its content.
type InputStr struct {
	value string
}

// Trims all leadind white spaces and new line.
// This works for either CRLF and LF.
func (is *InputStr) TrimSpaceAndNewLine() {
	is.value = strings.TrimLeftFunc(is.value, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n'
	})
}

// Chop characters from the start while the result of a predicate
// function is true.
// The chopped characters are returned as a string.
func (is *InputStr) ChopWhile(predicate func(rune) bool) string {
	var n int
	for n < len(is.value) && predicate(rune(is.value[n])) {
		n++
	}
	return is.ChopOff(n)
}

// Chops the first n characters from the start of the input
// and returns them as a string.
func (is *InputStr) ChopOff(n int) (ret string) {
	ret = is.value[:n]
	is.value = is.value[n:]
	return ret
}

// Returns true if the input string is empty.
func (is InputStr) IsEmpty() bool {
	return len(is.value) == 0
}

// Returns the first character of the input as a rune.
func (is InputStr) First() rune {
	return rune(is.value[0])
}
