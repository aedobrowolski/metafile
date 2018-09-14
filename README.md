# Metafile

A file system abstraction that adds persisted user metadata to file system
objects.

## Introduction

Metafile is a file system whose objects have user defined metadata.  Metafile
wraps [go-billy](https://github.com/src-d/go-billy), an abstract file system
that uses the same interfaces as the go libraries. In addition, every file
system object in Metafile implements the Metadata interface, allowing client
code to treat each file system object as a key-value map.

The file system objects are mostly files and directories. By default they have
no metadata. Adding metdata is easy (see example below). When persisted to disk
the metadata is stored in a hidden file at the root of the file system.

```go
// Create a file system
games metafile.Filesystem // can be either physical or in memory
games, _ = metafile.Osfs(os.TempDir())
games, _ = metafile.Memfs()

// Add a subdirectory to the root of the file system
games.MkdirAll("classic", 0777)

// Create an io.Writer, io.Reader, io.ReaderAt, io.Seeker, io.Closer
f, _ := games.Create("classic/board.txt")
f.Write([]byte("chess\nbackgammon\ncheckers\n"))

// Add some metadata
type Author struc { Name, Title, Org string }
f.Put("version", "1.0")
f.Put("author", Author{"Andrew", "TF", "PTC"})

// Move a file... does not affect metadata
game.Rename("classic/board.txt", "classic/twoPerson.txt")

// Retrieve some metadata
var auth Author
f, _ = game.Open("classic/twoPerson.txt")
f.Get("author", &auth)
fmt.Println(auth)

// Delete some metadata (ok if it does not exist)
f.Delete("version")
f.Delete("draft")
```

## Notes

Billy does not handle directories symmetrically in the os and in memory file
systems. You can open an OS directory but you cannot open an memfs directory. As
a result we will not open the directory to add metadata.
