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
		err := writeTriangleASCII(bw, t)
		if err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	_, err = bw.WriteString("endsolid " + s.header + "\n")
	if err != nil {
		return fmt.Errorf("did not write footer: %v", err)
	}

	return nil
}

func writeTriangleASCII(bw *bufio.Writer, t *Triangle) (err error) {
	_, err = bw.WriteString(fmt.Sprintf(" facet normal %s %s %s\n", shortFloat(t.normal.Ni), shortFloat(t.normal.Nj), shortFloat(t.normal.Nk)))
	_, err = bw.WriteString("  outer loop\n")
	for _, v := range t.vertices {
		_, err = bw.WriteString(fmt.Sprintf("   vertex %s %s %s\n", shortFloat(v.X), shortFloat(v.Y), shortFloat(v.Z)))
	}
	_, err = bw.WriteString("  endloop\n")
	_, err = bw.WriteString(" endfacet\n")

	return err
}
func shortFloat(f float32) string {
	// If number is an integer, return such
	if float64(f) == math.Floor(float64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}

	// Return the shortest scientific representation of this number
	return fmt.Sprintf("%g", f)
}
