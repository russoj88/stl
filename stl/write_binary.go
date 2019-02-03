package stl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func (s *STL) WriteBinary(w io.Writer) error {
	brw := bufio.NewWriter(w)
	defer brw.Flush()

	_, err := brw.WriteString(s.header)
	if err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	tcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(tcBytes, s.triangleCount)
	_, err = brw.Write(tcBytes)
	if err != nil {
		return fmt.Errorf("did not write triangle count: %v", err)
	}

	for _, t := range s.triangles {
		err := t.writeBinary(brw)
		if err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	return nil
}
func (t *Triangle) writeBinary(brw *bufio.Writer) error {
	// Collect all float32s that need to be written in order
	float32s := [12]float32{
		t.Normal.Ni, t.Normal.Nj, t.Normal.Nk,
		t.Vertices[0].X, t.Vertices[0].Y, t.Vertices[0].Z,
		t.Vertices[1].X, t.Vertices[1].Y, t.Vertices[1].Z,
		t.Vertices[2].X, t.Vertices[2].Y, t.Vertices[2].Z,
	}

	// Convert them to binary and write
	for _, f := range float32s {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, math.Float32bits(f))
		_, err := brw.Write(b)
		if err != nil {
			return err
		}
	}

	// Write out the attribute byte count
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, t.AttrByteCnt)
	_, err := brw.Write(b)

	return err
}
