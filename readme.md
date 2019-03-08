# STL

This is an [STL File](https://en.wikipedia.org/wiki/STL_(file_format) "Wiki") reader and writer written in Go.

### Package use
##### From, To
These core methods are to handle reading from an `io.Reader` and writing to an `io.Writer`.  Because most applications use files, these are wrapped in helper functions explained below.

##### FromFile
This takes in a filename and will return an `stl.Solid`.  This read method will automatically determine if the file is binary or ASCII and handle it appropriately.

This reader is concurrent.  For binary files, it gives about a 60% speedup on an E3-1231 v3 @ 3.40GHz reading off a SATA SSD.

##### ToASCIIFile
This will write an `stl.Solid` to a file in ASCII format.  The representations of the numbers are minimized to save some space.

##### ToBinaryFile
This will write an `stl.Solid` to a file in binary format.  This format is much more space efficient.

### Examples
##### Read from a file
```go
solid, err := stl.FromFile("/path/to/file.stl")
if err != nil {
    t.Errorf("could not read stl: %v", err)
}
```

##### Write to a file (ASCII)
```go
err = solid.ToASCIIFile("/path/to/file.stl")
if err != nil {
    t.Errorf("could not write to file: %v", err)
}
```
