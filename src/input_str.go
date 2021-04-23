package src

import "strings"

// Represent an input string with methods to
// mange its content.
type InputStr struct {
	Value string
}

// Creates a new InputStr from a string.
func Input(s string) InputStr {
	return InputStr{s}
}

// Trims all leadind white spaces and new line.
// This works for either CRLF and LF.
func (is *InputStr) TrimSpaceAndNewLine() {
	is.Value = strings.TrimLeftFunc(is.Value, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n'
	})
}

// Chop characters from the start while the result of a predicate
// function is true.
// The chopped characters are returned as a string.
func (is *InputStr) ChopWhile(predicate func(rune) bool) string {
	var n int
	for n < len(is.Value) && predicate(rune(is.Value[n])) {
		n++
	}
	return is.ChopOff(n)
}

// Chops the first n characters from the start of the input
// and returns them as a string.
func (is *InputStr) ChopOff(n int) (ret string) {
	ret = is.Value[:n]
	is.Value = is.Value[n:]
	return ret
}

// Returns true if the input string is empty.
func (is InputStr) IsEmpty() bool {
	return len(is.Value) == 0
}

// Returns the first character of the input as a rune.
func (is InputStr) First() rune {
	return rune(is.Value[0])
}
