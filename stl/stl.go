package stl

import (
	"runtime"
)

// Number of worker goroutines
var concurrencyLevel = runtime.NumCPU()

type Coordinate struct {
	X float32
	Y float32
	Z float32
}
type UnitVector struct {
	Ni float32
	Nj float32
	Nk float32
}
type Triangle struct {
	Normal      UnitVector
	Vertices    [3]Coordinate
	AttrByteCnt uint16
}
type Solid struct {
	Header        string
	TriangleCount uint32
	Triangles     []Triangle
}
