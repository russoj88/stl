package stl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func (s *STL) WriteBinary(w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	if _, err := bw.Write(headerBinary(s.header)); err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	if _, err := bw.Write(triCountBinary(s.triangleCount)); err != nil {
		return fmt.Errorf("did not write triangle count: %v", err)
	}

	for _, t := range s.triangles {
		if _, err := bw.Write(triangleBinary(t)); err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	return nil
}
func triCountBinary(u uint32) []byte {
	tcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(tcBytes, u)
	return tcBytes
}
func headerBinary(s string) []byte {
	return append([]byte(s), bytes.Repeat([]byte{0}, 80-len(s))...)
}
func triangleBinary(t *Triangle) []byte {
	bin := make([]byte, 0, 50)

	// Convert float32s to binary
	for _, f := range [12]float32{
		t.normal.Ni, t.normal.Nj, t.normal.Nk,
		t.vertices[0].X, t.vertices[0].Y, t.vertices[0].Z,
		t.vertices[1].X, t.vertices[1].Y, t.vertices[1].Z,
		t.vertices[2].X, t.vertices[2].Y, t.vertices[2].Z,
	} {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, math.Float32bits(f))
		bin = append(bin, b...)
	}

	// Attribute byte count binary
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, t.attrByteCnt)
	bin = append(bin, b...)

	return bin
}
