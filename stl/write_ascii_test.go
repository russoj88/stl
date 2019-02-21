package stl

import (
	"fmt"
	"testing"
)

func Test_shortFloat(t *testing.T) {
	for _, tst := range []struct {
		in       float32
		expected string
	}{
		{
			in:       1,
			expected: "1",
		},
		{
			in:       1.000,
			expected: "1",
		},
		{
			in:       45.754,
			expected: "45.754",
		},
		{
			in:       100000000,
			expected: "1e+08",
		},
		{
			in:       .0000000005,
			expected: "5e-10",
		},
		{
			in:       1000.00006,
			expected: "1000.00006",
		},
	} {
		tst := tst
		t.Run(fmt.Sprintf("shortFloat - %f", tst.in), func(t *testing.T) {
			t.Parallel()
			got := shortFloat(tst.in)
			if got != tst.expected {
				t.Errorf("Expecting %s, got %s", tst.expected, got)
			}
		})
	}
}
func Test_triangleASCII(t *testing.T) {
	tri := Triangle{
		normal: &UnitVector{
			Ni: 1,
			Nj: 2,
			Nk: 3,
		},
		vertices: [3]*Coordinate{
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
		attrByteCnt: 0,
	}
	expected := " facet normal 1 2 3\n  outer loop\n   vertex 1e+07 7 234.67\n   vertex 1234.34 8.231 1.345\n   vertex 8 1123 5\n  endloop\n endfacet\n"
	got := triangleASCII(&tri)
	if got != expected {
		t.Errorf("Expecting \n%s, \ngot \n%s", expected, got)
	}
}
