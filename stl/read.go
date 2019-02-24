package stl

import (
	"bufio"
	"fmt"
	"io"
)

// From creates a Solid from the input
func From(r io.Reader) (Solid, error) {
	// Use a buffered reader.  Default size is 4096 (4KB).
	br := bufio.NewReader(r)

	// Read first 5 bytes to get file type indicator.
	indicator, err := br.Peek(5)
	if err != nil {
		return Solid{}, fmt.Errorf("could not read from file: %v", err)
	}

	// If indicator is "solid" then it is an ASCII file.  Otherwise binary.
	if string(indicator) == "solid" {
		return fromASCII(br)
	}
	return fromBinary(br)
}
