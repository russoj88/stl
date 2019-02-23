package stl

import (
	"math"
	"runtime"
)

var concurrencyLevel = runtime.NumCPU()

// Allow users of stl package to override default concurrency level
func SetConcurrencyLevel(l int) {
	concurrencyLevel = int(math.Max(1, float64(l)))
}
