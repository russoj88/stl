package stl

import (
	"bufio"
	"strconv"
	"strings"
)

func fromASCII(br *bufio.Reader) (Solid, error) {
	header, err := extractASCIIHeader(br)
	if err != nil {
		return Solid{}, err
	}

	tris, err := extractASCIITriangle(br)
	if err != nil {
		return Solid{}, err
	}

	return Solid{
		Header:        header,
		TriangleCount: uint32(len(tris)),
		Triangles:     tris,
	}, nil
}

func extractASCIITriangle(br *bufio.Reader) (ts []*Triangle, err error) {
	scanner := bufio.NewScanner(br)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "normal") { // start of triangle
			t, err := extractTriangles(scanner)
			if err != nil {
				return nil, err
			}
			ts = append(ts, t)
		}
	}
	return
}

func extractTriangles(scanner *bufio.Scanner) (*Triangle, error) {
	// Get the normal for a triangle
	norm, err := extractUnitVec(scanner.Text())
	if err != nil {
		return nil, err
	}

	// Get past normal
	scanner.Scan()

	// Assumes 3 vertices
	var v [3]*Coordinate
	for i := 0; i < 3; i++ {
		scanner.Scan()
		c, err := extractCoords(scanner.Text())
		if err != nil {
			return nil, err
		}
		v[i] = c
	}

	// Pass "endloop"
	scanner.Scan()

	// Pass "endfacet"
	scanner.Scan()

	return &Triangle{
		Normal:      norm,
		Vertices:    v,
		AttrByteCnt: 0,
	}, nil
}

func extractCoords(s string) (*Coordinate, error) {
	sl := strings.Split(strings.TrimSpace(s), " ")
	f1, err := strconv.ParseFloat(sl[1], 32)
	if err != nil {
		return nil, err
	}
	f2, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return nil, err
	}
	f3, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return nil, err
	}

	return &Coordinate{
		X: float32(f1),
		Y: float32(f2),
		Z: float32(f3),
	}, nil
}

func extractUnitVec(s string) (*UnitVector, error) {
	//should be index 2-4 if spacing remains constant
	sl := strings.Split(strings.TrimSpace(s), " ")
	i, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return nil, err
	}
	j, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return nil, err
	}
	k, err := strconv.ParseFloat(sl[4], 32)
	if err != nil {
		return nil, err
	}

	return &UnitVector{
		Ni: float32(i),
		Nj: float32(j),
		Nk: float32(k),
	}, err
}

func extractASCIIHeader(br *bufio.Reader) (string, error) {
	s, _, err := br.ReadLine()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(string(s), "solid")), nil
}
