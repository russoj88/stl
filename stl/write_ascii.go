package stl

import (
	"bufio"
	"fmt"
	"io"
	"math"
)

func (s *STL) WriteASCII(w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	_, err := bw.WriteString("solid " + s.header + "\n")
	if err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	for _, t := range s.triangles {
		if _, err := bw.WriteString(triangleASCII(t)); err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	_, err = bw.WriteString("endsolid " + s.header + "\n")
	if err != nil {
		return fmt.Errorf("did not write footer: %v", err)
	}

	return nil
}

func triangleASCII(t *Triangle) string {
	return fmt.Sprintf(" facet normal %s %s %s\n", shortFloat(t.normal.Ni), shortFloat(t.normal.Nj), shortFloat(t.normal.Nk)) +
		"  outer loop\n" +
		fmt.Sprintf("   vertex %s %s %s\n", shortFloat(t.vertices[0].X), shortFloat(t.vertices[0].Y), shortFloat(t.vertices[0].Z)) +
		fmt.Sprintf("   vertex %s %s %s\n", shortFloat(t.vertices[1].X), shortFloat(t.vertices[1].Y), shortFloat(t.vertices[1].Z)) +
		fmt.Sprintf("   vertex %s %s %s\n", shortFloat(t.vertices[2].X), shortFloat(t.vertices[2].Y), shortFloat(t.vertices[2].Z)) +
		"  endloop\n" +
		" endfacet\n"
}
func shortFloat(f float32) string {
	// Scientific notation
	sn := fmt.Sprintf("%g", f)

	// If f is an integer, and its shorter than scientific notation form, return an integer
	if float64(f) == math.Floor(float64(f)) {
		in := fmt.Sprintf("%d", int64(f))
		if len(sn) > len(in) {
			return in
		}
	}

	return sn
}
