package stl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func (s *STL) WriteBinary(w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	_, err := bw.WriteString(s.header)
	if err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	tcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(tcBytes, s.triangleCount)
	_, err = bw.Write(tcBytes)
	if err != nil {
		return fmt.Errorf("did not write triangle count: %v", err)
	}

	for _, t := range s.triangles {
		err := t.writeBinary(bw)
		if err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	return nil
}
func (t *Triangle) writeBinary(bw *bufio.Writer) error {
	// Collect all float32s that need to be written in order
	float32s := [12]float32{
		t.normal.Ni, t.normal.Nj, t.normal.Nk,
		t.vertices[0].X, t.vertices[0].Y, t.vertices[0].Z,
		t.vertices[1].X, t.vertices[1].Y, t.vertices[1].Z,
		t.vertices[2].X, t.vertices[2].Y, t.vertices[2].Z,
	}

	// Convert them to binary and write
	for _, f := range float32s {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, math.Float32bits(f))
		_, err := bw.Write(b)
		if err != nil {
			return err
		}
	}

	// Write out the attribute byte count
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, t.attrByteCnt)
	_, err := bw.Write(b)

	return err
}
