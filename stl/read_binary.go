package stl

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
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

func readBinary(rd *bufio.Reader) (STL, error) {
	// Header is contained in the first 80 bytes.
	hBytes := make([]byte, 80)
	_, err := rd.Read(hBytes)
	if err != nil {
		return STL{}, fmt.Errorf("could not read header: %v", err)
	}

	// Triangle count is next 4 bytes.
	cntBytes := make([]byte, 4)
	_, err = rd.Read(cntBytes)
	if err != nil {
		return STL{}, fmt.Errorf("could not read triangle count: %v", err)
	}
	triCount := binary.LittleEndian.Uint32(cntBytes)

	// A list of triangles completes the document.  Each triangle is 50 bytes.
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
	err = sendBinaryToWorkers(rd, triCount, binToParse)
	if err != nil {
		return STL{}, err
	}
	close(binToParse)

	// When workGroup is done, close triParsed
	go func() {
		workGroup.Wait()
		close(triParsed)
	}()

	// Accumulate parsed triangles until triParsed channel is closed
	tris := accumulateTriangles(triCount, triParsed)

	return STL{
		header:        string(hBytes),
		triangleCount: triCount,
		triangles:     tris,
	}, nil
}
func sendBinaryToWorkers(rd *bufio.Reader, triCount uint32, work chan<- readWork) error {
	// Get bytes for each triangle and send to worker channel
	for i := 0; i < int(triCount); i += numTrianglesPerWorker {
		// Get bytes and put on channel
		bin := make([]byte, 50*numTrianglesPerWorker)
		n, err := io.ReadFull(rd, bin)
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
func accumulateTriangles(triCount uint32, in <-chan parsedWork) []Triangle {
	tris := make([]Triangle, triCount)
	for p := range in {
		for i := 0; i < len(p.t); i++ {
			tris[i+p.offset] = p.t[i]
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
		Normal: unitVectorFromBinary(bin[0:12]),
		vertices: [3]Coordinate{
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
