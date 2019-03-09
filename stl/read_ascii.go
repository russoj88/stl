package stl

import (
	"bufio"
	"fmt"
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
	s, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(string(s), "solid")), nil
}

// Parsing is done concurrently here depending on concurrencyLevel in config.go.
func extractASCIITriangles(br *bufio.Reader) (t []Triangle, err error) {
	// Read in ASCII data.  Put on work chan raw.
	raw, errChan := sendASCIIToWorkers(br)

	// Start up workers
	triParsed := make(chan Triangle)
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go parseTriangles(raw, triParsed, errChan, wg)
	}

	// When workers are done, close chans
	go func() {
		wg.Wait()
		close(triParsed)
		close(errChan)
	}()

	// Accumulate parsed Triangles until triParsed channel is closed
	return collectASCIITriangles(triParsed, errChan)
}
func sendASCIIToWorkers(br *bufio.Reader) (chan string, chan error) {
	work := make(chan string)
	// errChan needs a space to put error and return
	errChan := make(chan error, concurrencyLevel+1)

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

		if scanner.Err() != nil {
			errChan <- fmt.Errorf("error reading input: %v", scanner.Err())
		}
	}()

	return work, errChan
}
func parseTriangles(raw <-chan string, triParsed chan<- Triangle, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for r := range raw {
		sl := strings.Split(r, "\n")

		// Get the normal for a triangle
		norm, err := extractUnitVector(sl[0])
		if err != nil {
			errChan <- err
			return
		}

		// Get coordinates
		var v [3]Coordinate
		for i := 0; i < 3; i++ {
			v[i], err = extractCoordinate(sl[i+2])
			if err != nil {
				errChan <- err
				return
			}
		}

		triParsed <- Triangle{
			Normal:   norm,
			Vertices: v,
		}
	}
}
func extractCoordinate(s string) (Coordinate, error) {
	sl := strings.Split(strings.TrimSpace(s), " ")
	if len(sl) != 4 {
		return Coordinate{}, fmt.Errorf("invalid input for coordinate: %s", strings.TrimSpace(s))
	}

	x, err := strconv.ParseFloat(sl[1], 32)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid input for coordinate x: %v", err)
	}
	y, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid input for coordinate y: %v", err)
	}
	z, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid input for coordinate z: %v", err)
	}

	return Coordinate{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}, nil
}
func extractUnitVector(s string) (UnitVector, error) {
	sl := strings.Split(strings.TrimSpace(s), " ")
	if len(sl) != 5 {
		return UnitVector{}, fmt.Errorf("invalid input for unit vector: %s", strings.TrimSpace(s))
	}

	i, err := strconv.ParseFloat(sl[2], 32)
	if err != nil {
		return UnitVector{}, fmt.Errorf("invalid input for unit vector i: %v", err)
	}
	j, err := strconv.ParseFloat(sl[3], 32)
	if err != nil {
		return UnitVector{}, fmt.Errorf("invalid input for unit vector j: %v", err)
	}
	k, err := strconv.ParseFloat(sl[4], 32)
	if err != nil {
		return UnitVector{}, fmt.Errorf("invalid input for unit vector k: %v", err)
	}

	return UnitVector{
		Ni: float32(i),
		Nj: float32(j),
		Nk: float32(k),
	}, nil
}
func collectASCIITriangles(triParsed <-chan Triangle, errChan chan error) ([]Triangle, error) {
	// Creating space for 1K triangles as even simple designs have a few hundred
	tris := make([]Triangle, 0, 1024)
	for t := range triParsed {
		tris = append(tris, t)
	}

	err := <-errChan
	if err != nil {
		return nil, err
	}

	return tris, nil
}
