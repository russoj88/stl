package stl

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSTL_Header(t *testing.T) {
	goldenFile := "testdata/Utah_teapot.stl"
	dumpFile := "testdata/dump.stl"

	// Open file
	gFile, err := os.Open(goldenFile)
	defer gFile.Close()
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}

	// Read into STL type
	readSTL, err := Read(gFile)
	if err != nil {
		t.Errorf("could not read stl: %v", err)
	}

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
