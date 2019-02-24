package stl_test

import (
	"fmt"
	"gitlab.com/russoj88/stl/stl"
	"os"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	testFile := "testdata/Utah_teapot.stl"
	for _, testLevel := range []int{
		// Add different levels of concurrency here to see the best performance on a particular machine
		// The number of cores (x2 for hyper-threading) seem to get the best performance
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 32, 40, 48, 56, 64,
	} {
		b.Run(fmt.Sprintf("cl=%02d", testLevel), func(b *testing.B) {
			stl.SetConcurrencyLevel(testLevel)
			for i := 0; i < b.N; i++ {
				runRead(testFile, b)
			}
		})
	}
}
func runRead(testFile string, b *testing.B) {
	// Open file
	gFile, err := os.Open(testFile)
	if err != nil {
		b.Errorf("could not open file: %v", err)
	}
	defer gFile.Close()

	// Read into Solid type
	_, err = stl.From(gFile)
	if err != nil {
		b.Errorf("could not read stl: %v", err)
	}
}
