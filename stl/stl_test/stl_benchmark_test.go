package stl_test

import (
	"fmt"
	"gitlab.com/russoj88/stl/stl"
	"runtime"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	testFile := "testdata/Utah_teapot.stl"
	for _, testLevel := range []int{
		// Threads allowed for execution. This will then get used by concurrencyLevel to determine how many worker goroutines are made
		// The number of cores (x2 for hyper-threading) seem to get the best performance
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 32, 40, 48, 56, 64,
	} {
		b.Run(fmt.Sprintf("cl=%02d", testLevel), func(b *testing.B) {
			runtime.GOMAXPROCS(testLevel)
			for i := 0; i < b.N; i++ {
				runRead(testFile, b)
			}
		})
	}
}
func runRead(testFile string, b *testing.B) {
	// Read into Solid type
	_, err := stl.FromFile(testFile)
	if err != nil {
		b.Errorf("could not read stl: %v", err)
	}
}
