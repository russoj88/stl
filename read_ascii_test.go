package stl

import "testing"

func Test_extractUnitVector(t *testing.T) {
	for _, tst := range []struct {
		in   string
		want UnitVector
	}{
		{
			in: " facet normal 0.01388 -0.69223 -0.72154",
			want: UnitVector{
				Ni: 0.01388,
				Nj: -0.69223,
				Nk: -0.72154,
			},
		},
	} {
		t.Run("extractUnitVector", func(t *testing.T) {
			if got, _ := extractUnitVector(tst.in); got != tst.want {
				t.Errorf("got %v; want %v", got, tst.want)
			}
		})
	}
}
func Test_extractUnitVectorError(t *testing.T) {
	for _, tst := range []struct {
		in   string
		want string
	}{
		{
			in:   " facet normal 0.01388 -0.69223",
			want: `invalid input for unit vector: facet normal 0.01388 -0.69223`,
		},
		{
			in:   " facet normal 0.01388 -0.69223 a",
			want: `invalid input for unit vector k: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
		{
			in:   " facet normal 0.01388 a -0.72154",
			want: `invalid input for unit vector j: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
		{
			in:   " facet normal a -0.69223 a",
			want: `invalid input for unit vector i: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
	} {
		t.Run("extractUnitVector", func(t *testing.T) {
			if _, got := extractUnitVector(tst.in); got == nil || got.Error() != tst.want {
				t.Errorf("got %v; want %v", got, tst.want)
			}
		})
	}
}
func Test_extractCoordinate(t *testing.T) {
	for _, tst := range []struct {
		in   string
		want Coordinate
	}{
		{
			in: "   vertex -1000 0 0",
			want: Coordinate{
				X: -1000,
				Y: 0,
				Z: 0,
			},
		},
	} {
		t.Run("extractCoordinate", func(t *testing.T) {
			if got, _ := extractCoordinate(tst.in); got != tst.want {
				t.Errorf("got %v; want %v", got, tst.want)
			}
		})
	}
}
func Test_extractCoordinateError(t *testing.T) {
	for _, tst := range []struct {
		in   string
		want string
	}{
		{
			in:   "   vertex -1000 0 ",
			want: `invalid input for coordinate: vertex -1000 0`,
		},
		{
			in:   "   vertex -1000 0 a",
			want: `invalid input for coordinate z: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
		{
			in:   "   vertex -1000 a 0",
			want: `invalid input for coordinate y: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
		{
			in:   "   vertex a 0 0",
			want: `invalid input for coordinate x: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
	} {
		t.Run("extractCoordinate", func(t *testing.T) {
			if _, got := extractCoordinate(tst.in); got == nil || got.Error() != tst.want {
				t.Errorf("got %v; want %v", got, tst.want)
			}
		})
	}
}
