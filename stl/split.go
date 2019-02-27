package stl

import "bytes"

func splitTriangles(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// End on input
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the data before the 7th newline
	for n := 0; n < 7; n++ {
		idx := bytes.IndexByte(data[advance+1:], '\n')
		if idx < 0 {
			// Request more data
			return 0, nil, nil
		}
		advance += idx + 1
	}

	// Made it to the end of a token
	return advance + 1, data[:advance], nil
}
