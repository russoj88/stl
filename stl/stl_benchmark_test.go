package stl

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	testFile := "testdata/Utah_teapot.stl"
	for _, testLevel := range []uint32{
		// Add different levels of concurrency here to see the best performance on a particular machine
		// The number of cores (x2 for hyper-threading) seem to get the best performance
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 32, 40, 48, 56, 64,
	} {
		b.Run(fmt.Sprintf("cl=%02d", testLevel), func(b *testing.B) {
			SetConcurrencyLevel(testLevel)
			for i := 0; i < b.N; i++ {
				runForFile(testFile, b)
			}
		})
	}
}

func runForFile(testFile string, b *testing.B) {
	// Open file
	gFile, err := os.Open(testFile)
	if err != nil {
		b.Errorf("could not open file: %v", err)
	}
	defer gFile.Close()

	// Read into STL type
	_, err = Read(gFile)
	if err != nil {
		b.Errorf("could not read stl: %v", err)
	}
}
