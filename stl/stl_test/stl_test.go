package stl

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"gitlab.com/russoj88/stl/stl"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
)

func Test_Binary(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Utah_teapot.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	defer gFile.Close()
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}

	// Read into Solid type
	solid, err := stl.From(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Order triangles to make hash comparison between files
	sort.Slice(solid.Triangles, func(i, j int) bool {
		return strings.Compare(hash(solid.Triangles[i]), hash(solid.Triangles[j])) > 0
	})

	// Write to a binary buffer
	buffer := bytes.NewBuffer([]byte{})
	err = solid.ToBinary(buffer)
	if err != nil {
		t.Errorf("could not write to binary buffer: %v", err)
	}

	// Set the golden file reader to 0 so the contents of the file are actually read in
	_, _ = gFile.Seek(0, 0)

	// Confirm the buffer matches golden file
	if !contentsAreEqual(gFile, buffer) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func Test_ASCII(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Sphericon.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	defer gFile.Close()
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}

	// Read into Solid type
	solid, err := stl.From(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Order triangles to make hash comparison between files
	sort.Slice(solid.Triangles, func(i, j int) bool {
		return strings.Compare(hash(solid.Triangles[i]), hash(solid.Triangles[j])) > 0
	})

	// Write to a binary buffer
	buffer := bytes.NewBuffer([]byte{})
	err = solid.ToASCII(buffer)
	if err != nil {
		t.Errorf("could not write to binary buffer: %v", err)
	}

	// Set the golden file reader to 0 so the contents of the file are actually read in
	_, _ = gFile.Seek(0, 0)

	// Confirm the buffer matches golden file
	if !contentsAreEqual(gFile, buffer) {
		t.Errorf("Buffer and golden file are not equal!")
	}
}
func contentsAreEqual(r1, r2 io.Reader) bool {
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
func hash(t stl.Triangle) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", t.Normal)))
	h.Write([]byte(fmt.Sprintf("%v", t.Vertices[0])))
	h.Write([]byte(fmt.Sprintf("%v", t.Vertices[1])))
	h.Write([]byte(fmt.Sprintf("%v", t.Vertices[2])))
	h.Write([]byte(fmt.Sprintf("%v", t.AttrByteCnt)))
	return string(h.Sum(nil))
}
