package stl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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

// ToBinaryFile writes the Solid to a file in binary format
// See stl.ToBinary for more info
func (s *Solid) ToBinaryFile(filename string) error {
	file, err := os.OpenFile(strings.TrimSpace(filename), os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.ToBinary(file)
}
func headerBinary(s string) []byte {
	// Trim header down to 80 bytes
	if len(s) > 80 {
		s = s[:80]
	}

	// Pad header with zeroes
	return append([]byte(s), bytes.Repeat([]byte{0}, 80-len(s))...)
}
func triCountBinary(u uint32) []byte {
	tcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(tcBytes, u)
	return tcBytes
}
func triangleBinary(t Triangle) []byte {
	bin := make([]byte, 50)

	// Normal
	binary.LittleEndian.PutUint32(bin[0:4], math.Float32bits(t.Normal.Ni))
	binary.LittleEndian.PutUint32(bin[4:8], math.Float32bits(t.Normal.Nj))
	binary.LittleEndian.PutUint32(bin[8:12], math.Float32bits(t.Normal.Nk))

	// Vertex 1
	binary.LittleEndian.PutUint32(bin[12:16], math.Float32bits(t.Vertices[0].X))
	binary.LittleEndian.PutUint32(bin[16:20], math.Float32bits(t.Vertices[0].Y))
	binary.LittleEndian.PutUint32(bin[20:24], math.Float32bits(t.Vertices[0].Z))

	// Vertex 2
	binary.LittleEndian.PutUint32(bin[24:28], math.Float32bits(t.Vertices[1].X))
	binary.LittleEndian.PutUint32(bin[28:32], math.Float32bits(t.Vertices[1].Y))
	binary.LittleEndian.PutUint32(bin[32:40], math.Float32bits(t.Vertices[1].Z))

	// Vertex 3
	binary.LittleEndian.PutUint32(bin[36:40], math.Float32bits(t.Vertices[2].X))
	binary.LittleEndian.PutUint32(bin[40:44], math.Float32bits(t.Vertices[2].Y))
	binary.LittleEndian.PutUint32(bin[44:48], math.Float32bits(t.Vertices[2].Z))

	// Attribute byte count binary
	binary.LittleEndian.PutUint16(bin[48:50], t.AttrByteCnt)

	return bin
}
