package stl

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

// ToASCII writes the Solid out in ASCII form
func (s *Solid) ToASCII(w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	_, err := bw.WriteString("solid " + s.Header + "\n")
	if err != nil {
		return fmt.Errorf("did not write header: %v", err)
	}

	for _, t := range s.Triangles {
		if _, err := bw.WriteString(triangleASCII(t)); err != nil {
			return fmt.Errorf("did not write triangle: %v", err)
		}
	}

	_, err = bw.WriteString("endsolid " + s.Header + "\n")
	if err != nil {
		return fmt.Errorf("did not write footer: %v", err)
	}

	return nil
}

// Helper func to write ASCII directly to a file
func (s *Solid) ToASCIIFile(filename string) error {
	file, err := os.OpenFile(strings.TrimSpace(filename), os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.ToASCII(file)
}
func triangleASCII(t Triangle) string {
	return fmt.Sprintf(" facet normal %s %s %s\n", shortFloat(t.Normal.Ni), shortFloat(t.Normal.Nj), shortFloat(t.Normal.Nk)) +
		"  outer loop\n" +
		fmt.Sprintf("   vertex %s %s %s\n", shortFloat(t.Vertices[0].X), shortFloat(t.Vertices[0].Y), shortFloat(t.Vertices[0].Z)) +
		fmt.Sprintf("   vertex %s %s %s\n", shortFloat(t.Vertices[1].X), shortFloat(t.Vertices[1].Y), shortFloat(t.Vertices[1].Z)) +
		fmt.Sprintf("   vertex %s %s %s\n", shortFloat(t.Vertices[2].X), shortFloat(t.Vertices[2].Y), shortFloat(t.Vertices[2].Z)) +
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
