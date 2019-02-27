package stl

import (
	"bufio"
	"strconv"
	"strings"
	"sync"
)

func fromASCII(br *bufio.Reader) (Solid, error) {
	header, err := extractASCIIHeader(br)
	if err != nil {
		return Solid{}, err
	}

	tris, err := extractASCIITriangles(br)
	if err != nil {
		return Solid{}, err
	}

	return Solid{
		Header:        header,
		TriangleCount: uint32(len(tris)),
		Triangles:     tris,
	}, nil
}

func extractASCIITriangles(br *bufio.Reader) (ts []Triangle, err error) {
	scanner := bufio.NewScanner(br)
	scanner.Split(splitTriangles)
	triangles := make(chan Triangle)
	errs := make(chan error)
	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		go extractTriangle(scanner.Text(), triangles, errs, &wg)
	}
	go func() {
		wg.Wait()
		close(triangles)
		close(errs)
	}()
	for {
		select {
		case t, ok := <-triangles:
			if !ok {
				return
			}
			ts = append(ts, t)
		case err, _ = <-errs:
			if err != nil {
				return
			}
		}
	}
}

func extractTriangle(s string, triangles chan Triangle, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sl := strings.Split(s, "\n")

	// Get the normal for a triangle
	norm, err := extractUnitVec(sl[0])
	if err != nil {
		errs <- err
	}

	// Get coordinates
	var v [3]Coordinate
	for i := 0; i < 3; i++ {
		v[i], err = extractCoords(sl[i+2])
		if err != nil {
			errs <- err
		}
	}

	triangles <- Triangle{
		Normal:      norm,
		Vertices:    v,
		AttrByteCnt: 0,
	}
}

func extractCoords(s string) (Coordinate, error) {
	sl := strings.Split(strings.TrimSpace(s), " ")
	f1, err := strconv.ParseFloat(sl[1], 32)
	if err != nil {
		return Coordinate{}, err
	}
	f2, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return Coordinate{}, err
	}
	f3, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return Coordinate{}, err
	}

	return Coordinate{
		X: float32(f1),
		Y: float32(f2),
		Z: float32(f3),
	}, nil
}

func extractUnitVec(s string) (UnitVector, error) {
	sl := strings.Split(strings.TrimSpace(s), " ")
	i, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return UnitVector{}, err
	}
	j, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return UnitVector{}, err
	}
	k, err := strconv.ParseFloat(sl[4], 32)
	if err != nil {
		return UnitVector{}, err
	}

	return UnitVector{
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
