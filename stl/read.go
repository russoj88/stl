package stl

import (
	"bufio"
	"fmt"
	"io"
)

func Read(rd io.Reader) (STL, error) {
	// Use a buffered reader.  Default size is 4096 (4KB).
	brd := bufio.NewReader(rd)

	// Read first 5 bytes to get file type indicator.
	indicator, err := brd.Peek(5)
	if err != nil {
		return STL{}, fmt.Errorf("could not read from file: %v", err)
	}

	// If indicator is "solid" then it is an ASCII file.  Otherwise binary.
	if string(indicator) == "solid" {
		return readASCII(brd)
	}
	return readBinary(brd)
}
