# Metafile

Metafile is a file system abstraction that adds persisted user-defined metadata
to file system objects.

## Introduction

Metadata can add value to data in files and directories. For example a text
document has an author, a subject, a title, a version, keywords, and so on. Some
document formats provide ad hoc support for select metadata (e.g. Microsoft
Office and Open Office file formats). Document management systems have built in
support for metadata. But the operating system provides only the most basic
metadata for native files (owner, creation time, last modification time, access
control). Metafile can bridge the gap by providing an API to manage user defined
metadata for ordinary operating system files in a platform independent manner.

## Metadata Storage

The metadata is segregated into storage buckets that correspond to the file
system objects.  In this way multiple objects may reuse the same key with no
collision.  These buckets are named by relative paths to the objects from the
file system root. The paths are normalized following normal go conventions.

The methods `Put`, `Get`,  and `Delete` are used to store, retrieve and delete
the values associated with keys.  Values can be any go data type.  The `Get`
client must supply a pointer to a variable that matches the stored data.

Storage buckets are not created explicitly.  However they can be removed
explicitly with the `Empty` method or implicitly when the associated object is
created or removed.  They also get renamed when their object is renamed. Buckets
will not be copied when an object is copied.

One caveat: if the file system is changed without the use of the metafile API,
then the metadata will not reflect those changes.

It is an error to try to access a storage bucket for a file system object that
does not exist.  To check if a bucket can be accessed stat the corresponding
file system object.

## Example

The following sample shows the use of the metafile API.  Start by creating
or opening an existing filesystem and adding some files to it.

```go
// Create a file system
games metafile.Filesystem
games, _ = metafile.New(os.TempDir()) // for in memory: metafile.MemNew()
defer games.Close()

// Add a subdirectory to the root of the file system
games.MkdirAll("classic", 0777)

// Create an io.Writer, io.Reader, io.ReaderAt, io.Seeker, io.Closer
f, _ := games.Create("classic/board.txt")
f.Write([]byte("chess\nbackgammon\ncheckers\n"))
f.Close()
```

Associate some metadata with one or more of the file system objects. Close the
file system to test the persistence feature.

```go
// Add some metadata values, including complex structures
type Author struct { Name, Title, Org string }
games.Put("classic/board.txt", "author", Author{"Andrew", "TF", "example.com"})
games.Put("classic/board.txt", "version", "1.0")

// Move a file... this also moves the metadata
games.Rename("classic/board.txt", "twoPerson/board.txt")

// Delete some metadata (ok if the key does not exist)
games.Delete("twoPerson/board.txt", "version")
games.Delete("twoPerson/board.txt", "draft")

fsRoot := games.Root()
games.Close()
```

Reopen the file system from the same root. The metadata is still there.

```go
games, _ = metafile.New(fsRoot)
defer games.Close()

// Retrieve some metadata
var auth Author
games.Get("twoPerson/board.txt", "author", &auth)
fmt.Println(auth) // Author{"Andrew", "TF", "PTC"}
```

## Beneath the Hood

Metafile wraps [go-billy](https://github.com/src-d/go-billy), an abstract file
system that uses the same interfaces for files and directories as the go
libraries. Metafile adds a persisted key-value store for storing metadata with
every file system object. Adding metadata is done through the file system
abstraction. When persisted to disk the metadata storage is put in
a file at the root of the file system called "`.metadata`".

The format of the metadata file is not defined and is subject to change, but it
will remain backward compatible with all newer releases of Metafile.

### Why Gob format

Gob is an efficient persistance format native to go. It can encode and decode
virtually any go data types. It is version tolerant, meaning that if there is a
mismatch between the type used to encode and the one used to decode there is a
chance that the decoding will partially succeed. This is a feature that we want.

The problem with using Gob in this context is that a small encoded object, in
isolation, includes both the type definition and the value encoded as bytes.
Decoding needs to read the type before decoding the value.  Encoding the values
independently is inefficient but the only way in which the values can be read in
random order.

While it may be possible to reuse the Encoder or Decoder, they will need to be
reset before each value written or read.  Creating new ones is cleaner.