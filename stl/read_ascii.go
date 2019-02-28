package stl

import (
	"bufio"
	"bytes"
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

func extractASCIITriangles(br *bufio.Reader) (t []Triangle, err error) {
	scanner := bufio.NewScanner(br)
	scanner.Split(splitTriangles)
	doneTris := make(chan Triangle)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	out := extractTriangle(scanner)
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go parseTriangle(out, doneTris, errChan, wg)
	}

	go func() {
		wg.Wait()
		close(doneTris)
		close(errChan)
	}()
	return appendTriangles(doneTris, errChan)
}

func parseTriangle(out <-chan string, doneTris chan<- Triangle, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for o := range out {
		var v [3]Coordinate
		sl := strings.Split(o, "\n")

		// Get the normal for a triangle
		norm, err := extractUnitVec(sl[0])
		if err != nil {
			errChan <- err
		}

		// Get coordinates
		for i := 0; i < 3; i++ {
			v[i], err = extractCoords(sl[i+2])
			if err != nil {
				errChan <- err
			}
		}

		doneTris <- Triangle{
			Normal:   norm,
			Vertices: v,
		}
	}
}

func appendTriangles(in <-chan Triangle, errChan chan error) ([]Triangle, error) {
	// Read in all triangles
	tris := make([]Triangle, 0)
	for t := range in {
		tris = append(tris, t)
	}

	// If there is an error on errChan, return it
	err := <-errChan
	if err != nil {
		return nil, err
	}
	return tris, nil
}

func extractTriangle(b *bufio.Scanner) (out chan string) {
	out = make(chan string)
	go func() {
		defer close(out)
		for b.Scan() {
			out <- b.Text()
		}
	}()
	return
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

func splitTriangles(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// End on input
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the data before the 7th newline
	for n := 0; n < 7; n++ {
		idx := bytes.IndexByte(data[advance+1:], '\n')
		if idx < 0 {
			// Request more data
			return 0, nil, nil
		}
		advance += idx + 1
	}

	// Made it to the end of a token
	return advance + 1, data[:advance], nil
}
