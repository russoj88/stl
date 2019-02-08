package stl

import (
	"bufio"
	"fmt"
	"io"
)

func (s *STL) WriteASCII(w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	_, err := bw.WriteString("solid " + s.header + "\n")
	if err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	for _, t := range s.triangles {
		err := writeTriangleAscii(bw, t)
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

func writeTriangleAscii(bw *bufio.Writer, t *Triangle) (err error) {
	_, err = bw.WriteString(fmt.Sprintf(" facet normal %e %e %e\n", t.normal.Ni, t.normal.Nj, t.normal.Nk))
	_, err = bw.WriteString("  outer loop\n")
	for _, v := range t.vertices {
		_, err = bw.WriteString(fmt.Sprintf("   vertex %e %e %e\n", v.X, v.Y, v.Z))
	}
	_, err = bw.WriteString("  endloop\n")
	_, err = bw.WriteString(" endfacet\n")

	return err
}
