package stl

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
	normal      UnitVector
	vertices    [3]Coordinate
	attrByteCnt uint16
}
type STL struct {
	header        string
	triangleCount uint32
	triangles     []Triangle
}

func (s *STL) Header() string {
	return s.header
}
func (s *STL) TriangleCount() uint32 {
	return s.triangleCount
}
func (s *STL) Triangles() *[]Triangle {
	return &s.triangles
}
func (t *Triangle) Normal() UnitVector {
	return t.normal
}
func (t *Triangle) Vertices() [3]Coordinate {
	return t.vertices
}
func (t *Triangle) AttrByteCnt() uint16 {
	return t.attrByteCnt
}
func (c *Coordinate) SetX(x float32) {
	c.X = x
}
func (c *Coordinate) SetY(y float32) {
	c.Y = y
}
func (c *Coordinate) SetZ(z float32) {
	c.Z = z
}
