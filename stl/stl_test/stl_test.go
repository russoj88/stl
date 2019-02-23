package stl

import (
	"crypto/sha256"
	"fmt"
	"gitlab.com/russoj88/stl/stl"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestSTL_Binary(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Utah_teapot.stl"
	dumpFile := "testdata/dump.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	defer gFile.Close()
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}

	// Read into STL type
	readSTL, err := stl.Read(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Order triangles to make hash comparison between files
	sort.Slice(readSTL.Triangles(), func(i, j int) bool {
		return strings.Compare(hash(readSTL.Triangles()[i]), hash(readSTL.Triangles()[j])) > 0
	})

	// Write back to dump file
	dFile, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_RDWR, 0700)
	defer os.Remove(dumpFile)
	defer dFile.Close()
	if err != nil {
		t.Errorf("could not open dump file: %v", err)
	}
	err = readSTL.WriteBinary(dFile)
	if err != nil {
		t.Errorf("could not write to dump file: %v", err)
	}

	// Set file readers to 0 so the contents of the file are actually read in
	_, _ = gFile.Seek(0, 0)
	_, _ = dFile.Seek(0, 0)

	// Confirm the dumped file matches golden file
	if eq, err := fileAreEqual(gFile, dFile); err != nil {
		t.Errorf("error comparing: %v", err)
	} else if !eq {
		t.Errorf("Files are not equal!")
	}
}
func TestSTL_ASCII(t *testing.T) {
	t.Parallel()
	goldenFile := "testdata/Sphericon.stl"
	dumpFile := "testdata/dump2.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	defer gFile.Close()
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}

	// Read into STL type
	readSTL, err := stl.Read(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

	// Order triangles to make hash comparison between files
	sort.Slice(readSTL.Triangles(), func(i, j int) bool {
		return strings.Compare(hash(readSTL.Triangles()[i]), hash(readSTL.Triangles()[j])) > 0
	})

	// Write back to dump file
	dFile, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_RDWR, 0700)
	defer os.Remove(dumpFile)
	defer dFile.Close()
	if err != nil {
		t.Errorf("could not open dump file: %v", err)
	}
	err = readSTL.WriteASCII(dFile)
	if err != nil {
		t.Errorf("could not write to dump file: %v", err)
	}

	// Set file readers to 0 so the contents of the file are actually read in
	_, _ = gFile.Seek(0, 0)
	_, _ = dFile.Seek(0, 0)

	// Confirm the dumped file matches golden file
	if eq, err := fileAreEqual(gFile, dFile); err != nil {
		t.Errorf("error comparing: %v", err)
	} else if !eq {
		t.Errorf("Files are not equal!")
	}
}
func fileAreEqual(f1 *os.File, f2 *os.File) (bool, error) {
	h1, err := fileHash(f1)
	if err != nil {
		return false, err
	}
	h2, err := fileHash(f2)
	if err != nil {
		return false, err
	}

	return h1 == h2, nil
}
func fileHash(file *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
func hash(t *stl.Triangle) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", *t.Normal())))
	h.Write([]byte(fmt.Sprintf("%v", *t.Vertices()[0])))
	h.Write([]byte(fmt.Sprintf("%v", *t.Vertices()[1])))
	h.Write([]byte(fmt.Sprintf("%v", *t.Vertices()[2])))
	h.Write([]byte(fmt.Sprintf("%v", t.AttrByteCnt())))
	return string(h.Sum(nil))
}
