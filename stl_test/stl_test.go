package stl

import (
	"bytes"
	"io"
	"os"
	"sort"
	"testing"

	stl2 "gitlab.com/russoj88/stl"
)

func TestFromFile(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Utah_teapot.stl"

	// Read into Solid type
	solid, err := stl2.FromFile(goldenFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Write solid to a buffer
	buffer := writeToBuffer(solid, err, solid.ToBinary, t)

	// Confirm the buffer matches golden file
	gFile, err := os.Open(goldenFile)
	if err != nil {
		t.Fatalf("could not open golden file %s", goldenFile)
	}
	defer gFile.Close()
	if !contentsAreEqual(gFile, buffer) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func TestToASCIIFile(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/small_ASCII.stl"
	dumpFile := "testdata/dump_ASCII.stl"

	defer os.Remove(dumpFile)

	// Read into Solid type
	solid, err := stl2.FromFile(goldenFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	orderTriangles(&solid)

	// Write solid to a buffer
	err = solid.ToASCIIFile(dumpFile)
	if err != nil {
		t.Errorf("failed to write file: %v", err)
	}

	// Confirm the buffer matches golden file
	gFile, err := os.Open(goldenFile)
	if err != nil {
		t.Fatalf("could not open golden file %s", goldenFile)
	}
	dFile, err := os.Open(dumpFile)
	if err != nil {
		t.Fatalf("could not open dump file %s", dumpFile)
	}
	if !contentsAreEqual(gFile, dFile) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func TestToBinaryFile(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/small_binary.stl"
	dumpFile := "testdata/dump_binary.stl"

	defer os.Remove(dumpFile)

	// Read into Solid type
	solid, err := stl2.FromFile(goldenFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	orderTriangles(&solid)

	// Write solid to a buffer
	err = solid.ToBinaryFile(dumpFile)
	if err != nil {
		t.Errorf("failed to write file: %v", err)
	}

	// Confirm the buffer matches golden file
	gFile, err := os.Open(goldenFile)
	if err != nil {
		t.Fatalf("could not open golden file %s", goldenFile)
	}
	dFile, err := os.Open(dumpFile)
	if err != nil {
		t.Fatalf("could not open dump file %s", dumpFile)
	}
	if !contentsAreEqual(gFile, dFile) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func TestFrom_Binary(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Utah_teapot.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}
	defer gFile.Close()

	// Read into Solid type
	solid, err := stl2.From(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Write solid to a buffer
	buffer := writeToBuffer(solid, err, solid.ToBinary, t)

	// Set the golden file reader to 0 so the contents of the file are actually read in
	_, _ = gFile.Seek(0, 0)

	// Confirm the buffer matches golden file
	if !contentsAreEqual(gFile, buffer) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func TestFrom_BinaryError(t *testing.T) {
	t.Parallel()
	testFile := "testdata/invalid_binary.stl"

	// Read into Solid type
	_, err := stl2.FromFile(testFile)
	if err == nil {
		t.Errorf("expecting error, got none")
	}
}
func TestFrom_ASCII(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Sphericon.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}
	defer gFile.Close()

	// Read into Solid type
	solid, err := stl2.From(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Write solid to a buffer
	buffer := writeToBuffer(solid, err, solid.ToASCII, t)

	// Set the golden file reader to 0 so the contents of the file are actually read in
	_, _ = gFile.Seek(0, 0)

	// Confirm the buffer matches golden file
	if !contentsAreEqual(gFile, buffer) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func TestFrom_ASCIIErrorTriangle(t *testing.T) {
	t.Parallel()
	testFile := "testdata/invalid_ASCII_triangle.stl"

	// Read into Solid type
	_, err := stl2.FromFile(testFile)
	if err == nil {
		t.Errorf("expecting error, got none")
	}
}
func TestFrom_ASCIIErrorLine(t *testing.T) {
	t.Parallel()
	testFile := "testdata/invalid_ASCII_line.stl"

	// Read into Solid type
	_, err := stl2.FromFile(testFile)
	if err == nil {
		t.Errorf("expecting error, got none")
	}
}
func orderTriangles(solid *stl2.Solid) {
	sort.Slice(solid.Triangles, func(i, j int) bool {
		for idx := 0; idx < 3; idx++ {
			l := solid.Triangles[i].Vertices[idx]
			r := solid.Triangles[j].Vertices[idx]
			if l.X == r.X {
				if l.Y == r.Y {
					if l.Z == r.Z {
						continue
					}
					return l.Z < r.Z
				}
				return l.Y < r.Y
			}
			return l.X < r.X
		}

		return solid.Triangles[i].Normal.Ni < solid.Triangles[j].Normal.Ni
	})
}
func writeToBuffer(solid stl2.Solid, err error, To func(io.Writer) error, t *testing.T) *bytes.Buffer {
	orderTriangles(&solid)

	// Write to a binary buffer
	buffer := bytes.NewBuffer([]byte{})
	err = To(buffer)
	if err != nil {
		t.Errorf("could not write to binary buffer: %v", err)
	}
	return buffer
}
func contentsAreEqual(r1 io.Reader, r2 io.Reader) bool {
	for {
		buf1 := make([]byte, 4096)
		buf2 := make([]byte, 4096)

		n1, err1 := r1.Read(buf1)
		n2, err2 := r2.Read(buf2)

		if err1 != err2 || n1 != n2 || !bytes.Equal(buf1, buf2) {
			return false
		}

		if err1 == io.EOF && err2 == io.EOF {
			return true
		}
	}
}
