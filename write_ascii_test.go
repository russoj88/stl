package stl

import (
	"fmt"
	"testing"
)

func Test_shortFloat(t *testing.T) {
	for _, tst := range []struct {
		in   float32
		want string
	}{
		{
			in:   1,
			want: "1",
		},
		{
			in:   100000,
			want: "100000",
		},
		{
			in:   1000000,
			want: "1e+06",
		},
		{
			in:   1234567,
			want: "1234567",
		},
		{
			in:   1.000,
			want: "1",
		},
		{
			in:   45.754,
			want: "45.754",
		},
		{
			in:   .0000000005,
			want: "5e-10",
		},
		{
			in:   1000.00006,
			want: "1000.00006",
		},
	} {
		tst := tst
		t.Run(fmt.Sprintf("shortFloat - %f", tst.in), func(t *testing.T) {
			t.Parallel()
			if got := shortFloat(tst.in); got != tst.want {
				t.Errorf("got %s; want %s", got, tst.want)
			}
		})
	}
}
func Test_triangleASCII(t *testing.T) {
	tri := Triangle{
		Normal: UnitVector{
			Ni: 1,
			Nj: 2,
			Nk: 3,
		},
		Vertices: [3]Coordinate{
			{
				X: 10000000,
				Y: 7,
				Z: 234.67,
			},
			{
				X: 1234.34,
				Y: 8.231,
				Z: 1.345,
			},
			{
				X: 8,
				Y: 1123,
				Z: 5,
			},
		},
		AttrByteCnt: 0,
	}
	want := " facet normal 1 2 3\n  outer loop\n   vertex 1e+07 7 234.67\n   vertex 1234.34 8.231 1.345\n   vertex 8 1123 5\n  endloop\n endfacet\n"
	if got := triangleASCII(tri); got != want {
		t.Errorf("got \n%s, \nwant \n%s", got, want)
	}
}
