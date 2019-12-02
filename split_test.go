package stl

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func Test_splitTrianglesASCII(t *testing.T) {
	for _, tst := range []struct {
		in   []byte
		eof  bool
		want struct {
			advance int
			token   []byte
			err     error
		}
	}{
		// Empty input and EOF
		{
			in:  []byte{},
			eof: true,
			want: struct {
				advance int
				token   []byte
				err     error
			}{
				advance: 0,
				token:   nil,
				err:     nil,
			},
		},
		// Full token found
		{
			in:  []byte("facet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet\n"),
			eof: false,
			want: struct {
				advance int
				token   []byte
				err     error
			}{
				advance: 136,
				token:   []byte("facet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet"),
				err:     nil,
			},
		},
		// Beginning of token
		{
			in:  []byte("facet normal 0.05082 -0.24321 -0.96864\n"),
			eof: false,
			want: struct {
				advance int
				token   []byte
				err     error
			}{
				advance: 0,
				token:   nil,
				err:     nil,
			},
		},
		// Full token found, but there is more in data
		{
			in:  []byte("facet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet\nfacet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet\n"),
			eof: false,
			want: struct {
				advance int
				token   []byte
				err     error
			}{
				advance: 136,
				token:   []byte("facet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet"),
				err:     nil,
			},
		},
		// End of input, ignoring "endsolid"
		{
			in:  []byte("endsolid ASCII_STL_of_a_sphericon_by_CMG_Lee\n"),
			eof: true,
			want: struct {
				advance int
				token   []byte
				err     error
			}{
				advance: 0,
				token:   nil,
				err:     nil,
			},
		},
	} {
		tst := tst
		t.Run(fmt.Sprintf("splitTrianglesASCII - %q", tst.in), func(t *testing.T) {
			t.Parallel()
			gotAdvance, gotToken, gotError := splitTrianglesASCII(tst.in, tst.eof)
			if gotAdvance != tst.want.advance {
				t.Errorf("got %d; want %d", gotAdvance, tst.want.advance)
			}
			if !bytes.Equal(gotToken, tst.want.token) {
				t.Errorf("got \n%q\nwant\n%q", string(gotToken), string(tst.want.token))
			}
			if gotError == nil && tst.want.err != nil || gotError != nil && tst.want.err == nil {
				t.Errorf("unexpected nil")
				t.FailNow()
			}
			if gotError != nil && tst.want.err != nil && gotError.Error() != tst.want.err.Error() {
				t.Errorf("got %s; want %s", gotError.Error(), tst.want.err.Error())
			}
		})
	}
}
func Test_splitTrianglesScannerASCII(t *testing.T) {
	data := "solid ASCII_STL_of_a_sphericon_by_CMG_Lee\n facet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet\n facet normal -0.05382 -0.80723 0.58777\n  outer loop\n   vertex 0 -1000 0\n   vertex 995 0 105\n   vertex 988 0 156\n  endloop\n endfacet\n facet normal -0.06315 -0.82099 0.56743\n  outer loop\n   vertex 0 -1000 0\n   vertex 999 0 52\n   vertex 995 0 105\n  endloop\n endfacet\nendsolid ASCII_STL_of_a_sphericon_by_CMG_Lee\n"
	tokens := []string{
		" facet normal 0.05082 -0.24321 -0.96864\n  outer loop\n   vertex -1000 0 0\n   vertex 0 -358 -934\n   vertex 0 -407 -914\n  endloop\n endfacet",
		" facet normal -0.05382 -0.80723 0.58777\n  outer loop\n   vertex 0 -1000 0\n   vertex 995 0 105\n   vertex 988 0 156\n  endloop\n endfacet",
		" facet normal -0.06315 -0.82099 0.56743\n  outer loop\n   vertex 0 -1000 0\n   vertex 999 0 52\n   vertex 995 0 105\n  endloop\n endfacet",
	}

	// Create a buffered reader to get past the first line, which is the header
	buf := bufio.NewReader(bytes.NewReader([]byte(data)))
	_, _, _ = buf.ReadLine()

	// Create a scanner that takes in the splitTriangles split func
	scanner := bufio.NewScanner(buf)
	scanner.Split(splitTrianglesASCII)

	// Check that tokens are taken out in order
	for i := 0; scanner.Scan(); i++ {
		if scanner.Text() != tokens[i] {
			t.Errorf("got %q\n; want %q\n", scanner.Text(), tokens[i])
		}
	}
}
