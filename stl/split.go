package stl

import (
	"bytes"
	"fmt"
	"math"
)

func splitTrianglesASCII(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// End on input
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the data before the 7th newline
	for n := 0; n < 7; n++ {
		idx := bytes.IndexByte(data[advance+1:], '\n')
		if idx < 0 {
			if atEOF && (len(data) < 8 || string(data[:8]) != "endsolid") {
				return 0, nil, fmt.Errorf("invalid input data")
			}
			// Request more data
			return 0, nil, nil
		}
		advance += idx + 1
	}

	// Made it to the end of a token
	return advance + 1, data[:advance], nil
}
func splitTrianglesBinary(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const chunkSize = 50000

	// Return the next chunk, or ask for more data
	if len(data) >= chunkSize {
		return chunkSize, data[:chunkSize], nil
	}

	// Invalid data
	if atEOF && math.Mod(float64(len(data)), 50) != 0 {
		return 0, nil, fmt.Errorf("invalid input data")
	}

	// Last chunk of data
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}

	// Request more data
	return 0, nil, nil
}
