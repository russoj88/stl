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
func extractASCIIHeader(br *bufio.Reader) (string, error) {
	s, _, err := br.ReadLine()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(string(s), "solid")), nil
}
func extractASCIITriangles(br *bufio.Reader) (t []Triangle, err error) {
	// Collect parsed triangles
	doneTris := make(chan Triangle)
	errChan := make(chan error)

	// Read in ASCII data and send to workers
	out := sendASCIIToWorkers(br)

	// Start up workers
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go parseTriangles(out, doneTris, errChan, wg)
	}

	go func() {
		wg.Wait()
		close(doneTris)
		close(errChan)
	}()
	return collectASCIITriangles(doneTris, errChan)
}
func sendASCIIToWorkers(br *bufio.Reader) chan string {
	work := make(chan string)

	go func() {
		defer close(work)

		// Create Scanner with split func for ASCII triangles
		scanner := bufio.NewScanner(br)
		scanner.Split(splitTrianglesASCII)

		// Need to copy each read from the Scanner because it will be overwritten by the next Scan
		for scanner.Scan() {
			bin := make([]byte, len(scanner.Text()))
			copy(bin, scanner.Text())
			work <- string(bin)
		}
	}()
	return work
}
func parseTriangles(out <-chan string, doneTris chan<- Triangle, errChan chan error, wg *sync.WaitGroup) {
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
func collectASCIITriangles(in <-chan Triangle, errChan chan error) ([]Triangle, error) {
	// Read in all triangles
	// Creating space for 1K triangles as even simple designs have a few hundred
	tris := make([]Triangle, 0, 1024)
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
func extractCoords(s string) (Coordinate, error) {
	sl := strings.Split(strings.TrimSpace(s), " ")
	x, err := strconv.ParseFloat(sl[1], 32)
	if err != nil {
		return Coordinate{}, err
	}
	y, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return Coordinate{}, err
	}
	z, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return Coordinate{}, err
	}

	return Coordinate{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
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
