package stl

import "runtime"

var concurrencyLevel = uint32(runtime.NumCPU())

// Allow users of stl package to override default concurrency level
func SetConcurrencyLevel(l uint32) {
	concurrencyLevel = l
}
