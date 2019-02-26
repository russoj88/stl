package stl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// ToBinary writes the Solid out in binary form
func (s *Solid) ToBinary(w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	if _, err := bw.Write(headerBinary(s.Header)); err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	if _, err := bw.Write(triCountBinary(s.TriangleCount)); err != nil {
		return fmt.Errorf("did not write triangle count: %v", err)
	}

	for _, t := range s.Triangles {
		if _, err := bw.Write(triangleBinary(t)); err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	return nil
}
func headerBinary(s string) []byte {
	// Trim header down to 80 bytes
	if len(s) > 80 {
		s = s[:80]
	}

	return append([]byte(s), bytes.Repeat([]byte{0}, 80-len(s))...)
}
func triCountBinary(u uint32) []byte {
	tcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(tcBytes, u)
	return tcBytes
}
func triangleBinary(t Triangle) []byte {
	bin := make([]byte, 0, 50)

	// Convert float32s to binary
	for _, f := range [12]float32{
		t.Normal.Ni, t.Normal.Nj, t.Normal.Nk,
		t.Vertices[0].X, t.Vertices[0].Y, t.Vertices[0].Z,
		t.Vertices[1].X, t.Vertices[1].Y, t.Vertices[1].Z,
		t.Vertices[2].X, t.Vertices[2].Y, t.Vertices[2].Z,
	} {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, math.Float32bits(f))
		bin = append(bin, b...)
	}

	// Attribute byte count binary
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, t.AttrByteCnt)
	bin = append(bin, b...)

	return bin
}
