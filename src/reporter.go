package src

import (
	"fmt"
	"os"
	"strings"
)

// Struct that helps reporting errors
// and failures during the compilation process.
type Reporter struct {
	Input    string
	FileName string
}

// Factory that returns a new Reporter.
func NewReporter(input string, fileName string) Reporter {
	return Reporter{
		Input:    input,
		FileName: fileName}
}

// Represent the location of the report inside the file
// by his row and column.
type Location struct {
	Row, Col int
}

// Reports a warning message to stdin with given offset and args.
func (r Reporter) Warn(offset int, args ...interface{}) {
	r.print(offset, "WARN", args...)
}

// Reports an error message to stdin with given offset and args.
func (r Reporter) Fail(offset int, args ...interface{}) {
	r.print(offset, "ERROR", args...)
	os.Exit(1)
}

func (r Reporter) print(offset int, level string, args ...interface{}) {
	var loc Location = r.offsetToLocation(offset)
	fmt.Print(r.FileName, ":", loc.Col, ":", loc.Row, ": ", level, ": ")
	fmt.Print(args...)
	fmt.Print("\n")
}

func (r Reporter) OffsetFromInput(currentPosition string) int {
	return strings.Index(r.Input, currentPosition)
}

func (r Reporter) offsetToLocation(offset int) Location {
	lines := strings.Split(strings.ReplaceAll(r.Input, "\r\n", "\n"), "\n")

	for lineNum, line := range lines {
		if offset <= len(line) {
			return Location{Col: lineNum + 1, Row: offset + 1}
		}

		offset -= len(line) + 1
	}
	return Location{}
}
