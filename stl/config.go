package stl

import (
	"math"
	"runtime"
)

var concurrencyLevel = uint32(runtime.NumCPU())

// Allow users of stl package to override default concurrency level
func SetConcurrencyLevel(l uint32) {
	concurrencyLevel = uint32(math.Max(1, float64(l)))
}
