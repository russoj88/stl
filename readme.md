# STL

This is an [STL File](https://en.wikipedia.org/wiki/STL_(file_format) "Wiki") reader and writer written in Go.

### Package use
##### Read
This takes in an `io.Reader` and will return an `stl.Solid`.  This read method will automatically determine if the file is binary or ASCII and handle it appropriately.

This reader is concurrent.  For binary files, it gives about a 2x speedup on an E3-1231 v3 @ 3.40GHz reading off a SATA SSD.

The concurrency level can be controlled with `stl.SetConcurrencyLevel(l int)`.  1 is the minimum.  Best performance is found using the number of CPU cores returned by `runtime.NumCPU()` and is the default.  Performance tails off after this, but not by much.  Users of this library can easily test by running the benchmark in `stl_benchmark_test.go`.

##### WriteASCII
This will write an `stl.Solid` to an `io.Writer` in ASCII format.  The representations of the numbers are minimized to save some space.

##### WriteBinary
This will write an `stl.Solid` to an `io.Writer` in binary format.  This format is much more space efficient.

### Examples
##### Read from a file
```go
// Open file
gFile, err := os.Open("/tmp/obj.stl")
defer gFile.Close()
if err != nil {
    t.Errorf("could not open file: %v", err)
}

// Read into STL type
solid, err := From(gFile)
if err != nil {
    t.Errorf("could not read stl: %v", err)
}
```

##### Write to a file (ASCII)
```go
// Open file
dFile, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_RDWR, 0700)
defer dFile.Close()
if err != nil {
    t.Errorf("could not open file: %v", err)
}

// Write to file
err = solid.WriteASCII(dFile)
if err != nil {
    t.Errorf("could not write to file: %v", err)
}
```
