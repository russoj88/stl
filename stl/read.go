package stl

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Read from any abstracted reader
func Read(r io.Reader) (STL, error) {
	// Use a buffered reader.  Default size is 4096 (4KB).
	br := bufio.NewReader(r)

	// Read first 5 bytes to get file type indicator.
	indicator, err := br.Peek(5)
	if err != nil {
		return STL{}, fmt.Errorf("could not read from file: %v", err)
	}

	// If indicator is "solid" then it is an ASCII file.  Otherwise binary.
	if string(indicator) == "solid" {
		return readASCII(br)
	}
	return readBinary(br)
}

// Concurrent reader for file.  Binary files only.
func ReadFile(f *os.File) (STL, error) {
	// Use a buffered reader.  Default size is 4096 (4KB).
	br := bufio.NewReader(f)

	// Read first 5 bytes to get file type indicator.
	indicator, err := br.Peek(5)
	if err != nil {
		return STL{}, fmt.Errorf("could not read from file: %v", err)
	}

	// If indicator is "solid" then it is an ASCII file.  Otherwise binary.
	if string(indicator) == "solid" {
		return readASCII(br)
	}
	return readBinaryFile(f, br)
}
