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

const numTrianglesPerWorker = 1000

type readWork struct {
	offset int
	data   []byte
}
type parsedWork struct {
	offset int
	t      []Triangle
}

func readBinary(br *bufio.Reader) (STL, error) {
	header, err := extractBinaryHeader(br)
	if err != nil {
		return STL{}, err
	}
	header = strings.TrimSpace(header)

	triCount, err := extractBinaryTriangleCount(br)
	if err != nil {
		return STL{}, err
	}

	tris, err := extractBinaryTriangles(triCount, br)

	return STL{
		header:        header,
		triangleCount: triCount,
		triangles:     tris,
	}, nil
}
func extractBinaryTriangles(triCount uint32, br *bufio.Reader) ([]*Triangle, error) {
	// Each triangle is 50 bytes.
	// Parsing is done concurrently here depending on concurrencyLevel in config.go.
	binToParse := make(chan readWork)
	triParsed := make(chan parsedWork, triCount/numTrianglesPerWorker+1)

	// Start up workers
	workGroup := sync.WaitGroup{}
	for i := 0; i < int(concurrencyLevel); i++ {
		workGroup.Add(1)
		go parseChunksOfBinary(binToParse, triParsed, &workGroup)
	}

	// Read in binary and send chunks to workers
	err := sendBinaryToWorkers(br, triCount, binToParse)
	if err != nil {
		return nil, err
	}
	close(binToParse)

	// When workers are done, close triParsed
	go func() {
		workGroup.Wait()
		close(triParsed)
	}()

	// Accumulate parsed triangles until triParsed channel is closed
	return accumulateTriangles(triCount, triParsed), nil
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

	return string(hBytes), nil
}
func sendBinaryToWorkers(br *bufio.Reader, triCount uint32, work chan<- readWork) error {
	// Get bytes for each triangle and send to worker channel
	for i := 0; i < int(triCount); i += numTrianglesPerWorker {
		// Get bytes and put on channel
		bin := make([]byte, 50*numTrianglesPerWorker)
		n, err := io.ReadFull(br, bin)
		if err == io.ErrUnexpectedEOF {
			// This condition is for the last chunk which may not be complete
			bin = bin[:n]
		} else {
			if n < 50*numTrianglesPerWorker {
				return errors.New("did not read entire contents")
			}
			if err != nil {
				return fmt.Errorf("could not parse triangles: %v", err)
			}
		}

		work <- readWork{
			offset: i,
			data:   bin,
		}
	}

	return nil
}
func accumulateTriangles(triCount uint32, in <-chan parsedWork) []*Triangle {
	tris := make([]*Triangle, triCount)
	for p := range in {
		for i := 0; i < len(p.t); i++ {
			tris[i+p.offset] = &p.t[i]
		}
	}
	return tris
}
func parseChunksOfBinary(in <-chan readWork, out chan<- parsedWork, workGroup *sync.WaitGroup) {
	defer workGroup.Done()
	for w := range in {
		t := make([]Triangle, 0, len(w.data)/50)
		for i := 0; i < len(w.data); i += 50 {
			t = append(t, triangleFromBinary(w.data[i:i+50]))
		}
		out <- parsedWork{
			offset: w.offset,
			t:      t,
		}
	}
}
func triangleFromBinary(bin []byte) Triangle {
	return Triangle{
		normal: unitVectorFromBinary(bin[0:12]),
		vertices: [3]*Coordinate{
			coordinateFromBinary(bin[12:24]),
			coordinateFromBinary(bin[24:36]),
			coordinateFromBinary(bin[36:48]),
		},
		attrByteCnt: uint16(bin[48])<<8 | uint16(bin[49]),
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
