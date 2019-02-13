# STL

This is an [STL File](https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FSTL_(file_format) "Wiki") reader and writer written in Go.

###Package use
#####Read
This takes in an `io.Reader` and will return an `stl.STL`.  This read method will automatically determine if the file is binary or ASCII and handle it appropriately.

The io is done sequentially for ASCII files, and concurrently for binary files.  This concurrent reader gives almost a 4x speedup on an E3-1231 v3 @ 3.40GHz reading off a SATA SSD.

The concurrency level can be controlled with `stl.SetConcurrencyLevel(l uint32)`.  1 is the minimum.  Best performance is found using the number of CPU cores returned by `runtime.NumCPU()` and is the default.  Performance tails off after this, but not by much.  Users of this library can easily test by running the benchmark in `stl_benchmark_test.go`.

#####WriteASCII
This will write an `stl.STL` to an `io.Writer` in ASCII format.  The representations of the numbers are minimized to save some space.

#####WriteBinary
This will write an `stl.STL` to an `io.Writer` in binary format.  This format is much more space efficient.

###Examples
#####Read from a file
```go
// Open file
gFile, err := os.Open("/tmp/obj.stl")
defer gFile.Close()
if err != nil {
    t.Errorf("could not open file: %v", err)
}

// Read into STL type
readSTL, err := Read(gFile)
if err != nil {
    t.Errorf("could not read stl: %v", err)
}
```

#####Write to a file (ASCII)
```go
// Open file
dFile, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_RDWR, 0700)
defer dFile.Close()
if err != nil {
    t.Errorf("could not open file: %v", err)
}

// Write to file
err = readSTL.WriteASCII(dFile)
if err != nil {
    t.Errorf("could not write to file: %v", err)
}
```