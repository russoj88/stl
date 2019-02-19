package stl

import "testing"

func Test_unitVectorFromBinary(t *testing.T) {
	bin := []byte{0x4d, 0x10, 0x1c, 0x3f, 0x2a, 0xd2, 0xd3, 0xbe, 0x3a, 0x19, 0x2d, 0x3f}
	expected := UnitVector{
		Ni: 0.60962373,
		Nj: -0.4137128,
		Nk: 0.6761662,
	}
	got := unitVectorFromBinary(bin)
	if got.Ni != expected.Ni {
		t.Errorf("Expecting %g for Ni, got %g", expected.Ni, got.Ni)
	}
	if got.Nj != expected.Nj {
		t.Errorf("Expecting %g for Nj, got %g", expected.Nj, got.Nj)
	}
	if got.Nk != expected.Nk {
		t.Errorf("Expecting %g for Nk, got %g", expected.Nk, got.Nk)
	}
}
func Test_coordinateFromBinary(t *testing.T) {
	bin := []byte{0x4d, 0xbe, 0x6b, 0x40, 0xa6, 0xe5, 0xd2, 0x3f, 0xc1, 0x40, 0xd8, 0x40}
	expected := Coordinate{
		X: 3.68349,
		Y: 1.6476333,
		Z: 6.7579045,
	}
	got := coordinateFromBinary(bin)
	if got.X != expected.X {
		t.Errorf("Expecting %g for X, got %g", expected.X, got.X)
	}
	if got.Y != expected.Y {
		t.Errorf("Expecting %g for Y, got %g", expected.Y, got.Y)
	}
	if got.Z != expected.Z {
		t.Errorf("Expecting %g for Z, got %g", expected.Z, got.Z)
	}
}
func Test_triangleFromBinary(t *testing.T) {
	bin := []byte{0x8a, 0xa5, 0x10, 0x3f, 0xbd, 0x5f, 0x8f, 0x3e, 0x73, 0xae, 0x46, 0x3f, 0x4d, 0xbe, 0x6b, 0x40, 0xa6, 0xe5, 0xd2, 0x3f, 0xc1, 0x40, 0xd8, 0x40, 0xb0, 0x38, 0x5d, 0x40, 0x13, 0xc0, 0x06, 0x40, 0xc1, 0x40, 0xd8, 0x40, 0xb4, 0xad, 0x5a, 0x40, 0xbd, 0x41, 0x05, 0x40, 0xb3, 0x72, 0xd9, 0x40, 0x00, 0x00}
	expected := Triangle{
		normal: &UnitVector{
			Ni: 0.5650259,
			Nj: 0.2800273,
			Nk: 0.7760994,
		},
		vertices: [3]*Coordinate{
			{
				X: 3.68349,
				Y: 1.6476333,
				Z: 6.7579045,
			},
			{
				X: 3.456585,
				Y: 2.1054733,
				Z: 6.7579045,
			},
			{
				X: 3.416852,
				Y: 2.0821373,
				Z: 6.7952514,
			},
		},
		attrByteCnt: 0,
	}

	got := triangleFromBinary(bin)
	if *got.normal != *expected.normal {
		t.Errorf("Expecting %+v, got %+v for normal", *expected.normal, *got.normal)
	}
	for i := 0; i < 3; i++ {
		if *got.vertices[i] != *expected.vertices[i] {
			t.Errorf("Expecting %+v, got %+v for vertex %d", *got.vertices[i], *expected.vertices[i], i)
		}
	}
	if got.attrByteCnt != expected.attrByteCnt {
		t.Errorf("Expecting %d, got %d for attrByteCnt", expected.attrByteCnt, got.attrByteCnt)
	}
}
