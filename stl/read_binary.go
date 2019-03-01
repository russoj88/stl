package stl

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
)

func fromBinary(br *bufio.Reader) (Solid, error) {
	header, err := extractBinaryHeader(br)
	if err != nil {
		return Solid{}, err
	}

	triCount, err := extractBinaryTriangleCount(br)
	if err != nil {
		return Solid{}, err
	}
	tris, err := extractBinaryTriangles(triCount, br)
	if err != nil {
		return Solid{}, err
	}

	return Solid{
		Header:        header,
		TriangleCount: triCount,
		Triangles:     tris,
	}, nil
}
func extractBinaryHeader(br *bufio.Reader) (string, error) {
	hBytes := make([]byte, 80)
	_, err := br.Read(hBytes)
	if err != nil {
		return "", fmt.Errorf("could not read header: %v", err)
	}

	return strings.TrimSpace(string(hBytes)), nil
}
func extractBinaryTriangleCount(br *bufio.Reader) (uint32, error) {
	cntBytes := make([]byte, 4)
	_, err := br.Read(cntBytes)
	if err != nil {
		return 0, fmt.Errorf("could not read triangle count: %v", err)
	}

	return binary.LittleEndian.Uint32(cntBytes), nil
}
func extractBinaryTriangles(triCount uint32, br *bufio.Reader) ([]Triangle, error) {
	// Each triangle is 50 bytes.
	// Parsing is done concurrently here depending on concurrencyLevel in config.go.
	triParsed := make(chan []Triangle, concurrencyLevel)
	errChan := make(chan error, concurrencyLevel+1)

	// Read in binary and send chunks to workers
	raw := sendBinaryToWorkers(br, errChan)

	// Start up workers
	workGroup := sync.WaitGroup{}
	for i := 0; i < concurrencyLevel; i++ {
		workGroup.Add(1)
		go parseChunksOfBinary(raw, triParsed, errChan, &workGroup)
	}

	// When workers are done, close triParsed
	go func() {
		workGroup.Wait()
		close(triParsed)
		close(errChan)
	}()

	// Accumulate parsed Triangles until triParsed channel is closed
	return collectBinaryTriangles(triCount, triParsed, errChan)
}
func sendBinaryToWorkers(br *bufio.Reader, errChan chan error) chan []byte {
	raw := make(chan []byte)

	go func() {
		defer close(raw)

		// Create Scanner with split func for binary triangle chunks
		scanner := bufio.NewScanner(br)
		scanner.Split(splitTrianglesBinary)

		// Need to copy each read from the Scanner because it will be overwritten by the next Scan
		for scanner.Scan() {
			bin := make([]byte, len(scanner.Bytes()))
			copy(bin, scanner.Bytes())
			raw <- bin
		}

		if scanner.Err() != nil {
			errChan <- scanner.Err()
		}
	}()

	return raw
}
func parseChunksOfBinary(raw <-chan []byte, triParsed chan<- []Triangle, errChan chan error, workGroup *sync.WaitGroup) {
	defer workGroup.Done()
	// Catch panics
	defer func() {
		if r := recover(); r != nil {
			errChan <- errors.New(fmt.Sprintf("unable to parse triangle from input"))
		}
	}()

	for r := range raw {
		t := make([]Triangle, 0, len(r)/50)
		for i := 0; i < len(r); i += 50 {
			t = append(t, triangleFromBinary(r[i:i+50]))
		}
		triParsed <- t
	}
}
func collectBinaryTriangles(triCount uint32, triParsed <-chan []Triangle, errChan <-chan error) ([]Triangle, error) {
	// Read triParsed all triangles
	tris := make([]Triangle, 0, triCount)
	for t := range triParsed {
		tris = append(tris, t...)
	}

	err := <-errChan
	if err != nil {
		return nil, err
	}

	return tris, nil
}
func triangleFromBinary(bin []byte) Triangle {
	return Triangle{
		Normal: unitVectorFromBinary(bin[0:12]),
		Vertices: [3]Coordinate{
			coordinateFromBinary(bin[12:24]),
			coordinateFromBinary(bin[24:36]),
			coordinateFromBinary(bin[36:48]),
		},
		AttrByteCnt: uint16(bin[48])<<8 | uint16(bin[49]),
	}
}
func coordinateFromBinary(bin []byte) Coordinate {
	return Coordinate{
		X: math.Float32frombits(binary.LittleEndian.Uint32(bin[0:4])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(bin[4:8])),
		Z: math.Float32frombits(binary.LittleEndian.Uint32(bin[8:12])),
	}
}
func unitVectorFromBinary(bin []byte) UnitVector {
	return UnitVector{
		Ni: math.Float32frombits(binary.LittleEndian.Uint32(bin[0:4])),
		Nj: math.Float32frombits(binary.LittleEndian.Uint32(bin[4:8])),
		Nk: math.Float32frombits(binary.LittleEndian.Uint32(bin[8:12])),
	}
}
