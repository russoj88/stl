package stl

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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

	return Solid{
		Header:        header,
		TriangleCount: triCount,
		Triangles:     tris,
	}, nil
}
func extractBinaryTriangles(triCount uint32, br *bufio.Reader) ([]*Triangle, error) {
	// Each triangle is 50 bytes.
	// Parsing is done concurrently here depending on concurrencyLevel in config.go.
	triParsed := make(chan *[]*Triangle, concurrencyLevel)

	// Read in binary and send chunks to workers
	binToParse, errChan := sendBinaryToWorkers(br, triCount)

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
	return accumulateTriangles(triCount, triParsed, errChan)
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
func sendBinaryToWorkers(br *bufio.Reader, triCount uint32) (work chan *[]byte, errChan chan error) {
	errChan = make(chan error, 1)
	work = make(chan *[]byte, concurrencyLevel)
	const numTrianglesPerWorkUnit = 1000

	go func() {
		// Close channel when done
		defer close(work)

		// Get bytes for each triangle and send to worker channel
		for i := 0; i < int(triCount); i += numTrianglesPerWorkUnit {
			// Get bytes and put on channel
			bin := make([]byte, 50*numTrianglesPerWorkUnit)
			n, err := io.ReadFull(br, bin)
			if err == io.ErrUnexpectedEOF {
				// This condition is for the last chunk which may not be complete
				bin = bin[:n]
			} else {
				if n < 50*numTrianglesPerWorkUnit {
					errChan <- errors.New("did not read entire contents")
					return
				}
				if err != nil {
					errChan <- fmt.Errorf("could not parse Triangles: %v", err)
					return
				}
			}

			work <- &bin
		}
	}()

	return work, errChan
}
func accumulateTriangles(triCount uint32, in <-chan *[]*Triangle, errChan chan error) ([]*Triangle, error) {
	defer close(errChan)
	tris := make([]*Triangle, 0, triCount)
	for {
		select {
		case p, more := <-in:
			if !more {
				return tris, nil
			}
			tris = append(tris, *p...)
		case err := <-errChan:
			return nil, err
		}
	}
}
func parseChunksOfBinary(in <-chan *[]byte, out chan<- *[]*Triangle, workGroup *sync.WaitGroup) {
	defer workGroup.Done()
	for w := range in {
		t := make([]*Triangle, 0, len(*w)/50)
		for i := 0; i < len(*w); i += 50 {
			t = append(t, triangleFromBinary((*w)[i:i+50]))
		}
		out <- &t
	}
}
func triangleFromBinary(bin []byte) *Triangle {
	return &Triangle{
		Normal: unitVectorFromBinary(bin[0:12]),
		Vertices: [3]*Coordinate{
			coordinateFromBinary(bin[12:24]),
			coordinateFromBinary(bin[24:36]),
			coordinateFromBinary(bin[36:48]),
		},
		AttrByteCnt: uint16(bin[48])<<8 | uint16(bin[49]),
	}
}
func coordinateFromBinary(bin []byte) *Coordinate {
	return &Coordinate{
		X: math.Float32frombits(binary.LittleEndian.Uint32(bin[0:4])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(bin[4:8])),
		Z: math.Float32frombits(binary.LittleEndian.Uint32(bin[8:12])),
	}
}
func unitVectorFromBinary(bin []byte) *UnitVector {
	return &UnitVector{
		Ni: math.Float32frombits(binary.LittleEndian.Uint32(bin[0:4])),
		Nj: math.Float32frombits(binary.LittleEndian.Uint32(bin[4:8])),
		Nk: math.Float32frombits(binary.LittleEndian.Uint32(bin[8:12])),
	}
}
