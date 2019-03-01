package stl

import (
	"bufio"
	"encoding/binary"
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

	return Solid{
		Header:        header,
		TriangleCount: triCount,
		Triangles:     extractBinaryTriangles(triCount, br),
	}, nil
}
func extractBinaryTriangles(triCount uint32, br *bufio.Reader) []Triangle {
	// Each triangle is 50 bytes.
	// Parsing is done concurrently here depending on concurrencyLevel in config.go.
	triParsed := make(chan []Triangle, concurrencyLevel)

	// Read in binary and send chunks to workers
	binToParse := sendBinaryToWorkers(br)

	// Start up workers
	workGroup := sync.WaitGroup{}
	for i := 0; i < concurrencyLevel; i++ {
		workGroup.Add(1)
		go parseChunksOfBinary(binToParse, triParsed, &workGroup)
	}

	// When workers are done, close triParsed
	go func() {
		workGroup.Wait()
		close(triParsed)
	}()

	// Accumulate parsed Triangles until triParsed channel is closed
	return accumulateTriangles(triCount, triParsed)
}
func extractBinaryTriangleCount(br *bufio.Reader) (uint32, error) {
	cntBytes := make([]byte, 4)
	_, err := br.Read(cntBytes)
	if err != nil {
		return 0, fmt.Errorf("could not read triangle count: %v", err)
	}

	return binary.LittleEndian.Uint32(cntBytes), nil
}
func extractBinaryHeader(br *bufio.Reader) (string, error) {
	hBytes := make([]byte, 80)
	_, err := br.Read(hBytes)
	if err != nil {
		return "", fmt.Errorf("could not read header: %v", err)
	}

	return strings.TrimSpace(string(hBytes)), nil
}
func sendBinaryToWorkers(br *bufio.Reader) chan []byte {
	work := make(chan []byte)

	go func() {
		// Close channel when done
		defer close(work)

		// Create scanner
		scanner := bufio.NewScanner(br)
		scanner.Split(splitTrianglesBinary)
		for scanner.Scan() {
			bin := make([]byte, len(scanner.Bytes()))
			copy(bin, scanner.Bytes())
			work <- bin
		}
	}()

	return work
}
func accumulateTriangles(triCount uint32, in <-chan []Triangle) []Triangle {
	// Read in all triangles
	tris := make([]Triangle, 0, triCount)
	for t := range in {
		tris = append(tris, t...)
	}

	return tris
}
func parseChunksOfBinary(in <-chan []byte, out chan<- []Triangle, workGroup *sync.WaitGroup) {
	defer workGroup.Done()
	for w := range in {
		t := make([]Triangle, 0, len(w)/50)
		for i := 0; i < len(w); i += 50 {
			t = append(t, triangleFromBinary(w[i:i+50]))
		}
		out <- t
	}
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
