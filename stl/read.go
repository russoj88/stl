package stl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// From creates a Solid from the input.
// It handles both ASCII and binary formats.
func From(r io.Reader) (s Solid, err error) {
	// Catch panics
	defer func() {
		if r := recover(); r != nil {
			s = Solid{}
			err = fmt.Errorf("unable to parse input")
		}
	}()

	// Use a buffered reader.  Default size is 4096 (4KB).
	br := bufio.NewReader(r)

	// Read first 6 bytes to get file type indicator.
	indicator, err := br.Peek(6)
	if err != nil {
		return Solid{}, fmt.Errorf("input has no content")
	}

	// If indicator is "solid " then it is an ASCII file.  Otherwise binary.
	if string(indicator) == "solid " {
		return fromASCII(br)
	}

	return fromBinary(br)
}

// FromFile creates a Solid from a file
// See stl.From for more info
func FromFile(filename string) (Solid, error) {
	// Open file for reading
	file, err := os.Open(strings.TrimSpace(filename))
	if err != nil {
		return Solid{}, err
	}
	defer file.Close()

	return From(file)
}
