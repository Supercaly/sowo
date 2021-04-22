package src

import "strings"

type InputStr struct {
	Value string
}

func Input(s string) InputStr {
	return InputStr{s}
}

func (is *InputStr) TrimSpaceAndNewLine() {
	is.Value = strings.TrimLeftFunc(is.Value, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n'
	})
}

func (is *InputStr) ChopWhile(predicate func(rune) bool) string {
	var n int
	for n < len(is.Value) && predicate(rune(is.Value[n])) {
		n++
	}
	return is.ChopOff(n)
}

func (is *InputStr) ChopOff(n int) (ret string) {
	ret = is.Value[:n]
	is.Value = is.Value[n:]
	return ret
}

func (is InputStr) IsEmpty() bool {
	return len(is.Value) == 0
}

func (is InputStr) First() rune {
	return rune(is.Value[0])
}
