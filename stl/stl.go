package stl

import (
	"runtime"
)

// Number of worker goroutines
var concurrencyLevel = runtime.NumCPU()

// X, Y, and Z of a Triangle Vertex
type Coordinate struct {
	X float32
	Y float32
	Z float32
}

// i, j, and k of the unit vector for a Triangle's normal
type UnitVector struct {
	Ni float32
	Nj float32
	Nk float32
}

// Triangle
// AttrByteCnt does not get recorded for ASCII type
// It is for binary only, and does not have a standard use
type Triangle struct {
	Normal      UnitVector
	Vertices    [3]Coordinate
	AttrByteCnt uint16
}

// STL object with all triangles
type Solid struct {
	Header        string
	TriangleCount uint32
	Triangles     []Triangle
}
