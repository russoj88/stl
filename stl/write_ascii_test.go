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
