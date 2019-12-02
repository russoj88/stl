package stl

import "testing"

func Test_unitVectorFromBinary(t *testing.T) {
	bin := []byte{0x4d, 0x10, 0x1c, 0x3f, 0x2a, 0xd2, 0xd3, 0xbe, 0x3a, 0x19, 0x2d, 0x3f}
	want := UnitVector{
		Ni: 0.60962373,
		Nj: -0.4137128,
		Nk: 0.6761662,
	}
	got := unitVectorFromBinary(bin)
	if got.Ni != want.Ni {
		t.Errorf("got %g for Ni; want %g", got.Ni, want.Ni)
	}
	if got.Nj != want.Nj {
		t.Errorf("got %g for Nj; want %g", got.Nj, want.Nj)
	}
	if got.Nk != want.Nk {
		t.Errorf("got %g for Nk; want %g", got.Nk, want.Nk)
	}
}
func Test_coordinateFromBinary(t *testing.T) {
	bin := []byte{0x4d, 0xbe, 0x6b, 0x40, 0xa6, 0xe5, 0xd2, 0x3f, 0xc1, 0x40, 0xd8, 0x40}
	want := Coordinate{
		X: 3.68349,
		Y: 1.6476333,
		Z: 6.7579045,
	}
	got := coordinateFromBinary(bin)
	if got.X != want.X {
		t.Errorf("got %g for X; want %g", got.X, want.X)
	}
	if got.Y != want.Y {
		t.Errorf("got %g for Y; want %g", got.Y, want.Y)
	}
	if got.Z != want.Z {
		t.Errorf("got %g for Z; want %g", got.Z, want.Z)
	}
}
func Test_triangleFromBinary(t *testing.T) {
	bin := []byte{0x8a, 0xa5, 0x10, 0x3f, 0xbd, 0x5f, 0x8f, 0x3e, 0x73, 0xae, 0x46, 0x3f, 0x4d, 0xbe, 0x6b, 0x40, 0xa6, 0xe5, 0xd2, 0x3f, 0xc1, 0x40, 0xd8, 0x40, 0xb0, 0x38, 0x5d, 0x40, 0x13, 0xc0, 0x06, 0x40, 0xc1, 0x40, 0xd8, 0x40, 0xb4, 0xad, 0x5a, 0x40, 0xbd, 0x41, 0x05, 0x40, 0xb3, 0x72, 0xd9, 0x40, 0x00, 0x00}
	want := Triangle{
		Normal: UnitVector{
			Ni: 0.5650259,
			Nj: 0.2800273,
			Nk: 0.7760994,
		},
		Vertices: [3]Coordinate{
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
		AttrByteCnt: 0,
	}

	got := triangleFromBinary(bin)
	if got.Normal != want.Normal {
		t.Errorf("got %+v for normal; want %+v", got.Normal, want.Normal)
	}
	for i := 0; i < 3; i++ {
		if got.Vertices[i] != want.Vertices[i] {
			t.Errorf("got %+v for vertex %d; want %+v", want.Vertices[i], i, got.Vertices[i])
		}
	}
	if got.AttrByteCnt != want.AttrByteCnt {
		t.Errorf("got %d for attrByteCnt; want %d", got.AttrByteCnt, want.AttrByteCnt)
	}
}
