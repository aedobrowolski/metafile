package metafile

import (
	"os"

	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

// filesystem implements Filesystem by wrapping billy.Filesystem and adding metadata.
type filesystem struct {
	billy.Filesystem
	store
}

var _ Filesystem = filesystem{} // ensure interface compliance

// New returns a new metafile filesystem based on the OS filesystem
func New(baseDir string) Filesystem {
	f := osfs.New(baseDir)
	s := newStore(baseDir, nil) // TODO: add fs validation of path
	return filesystem{f, s}
}

// NewMem returns a new metafile filesystem based on an in-memory filesystem
func NewMem() Filesystem {
	f := memfs.New()
	s := newStore("", func(b string) bool {
		stat, err := f.Stat(b)
		if err != nil {
			return false
		}
		mode := stat.Mode()
		return mode == 0 || mode == os.ModeDir
	})
	return filesystem{f, s}
}

// Override the following Basic methods that modify paths in the fs
// - Create (remove metadata)
// - Rename (move metadata)
// - Remove (remove metadata)

// Create creates the named file with mode 0666 (before umask), truncating
// it if it already exists. If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
// Metadata for the file will be removed.
func (fs filesystem) Create(filename string) (billy.File, error) {
	f, err := fs.Filesystem.Create(filename)
	if err != nil {
		return nil, err
	}
	fs.store.emptyOne(filename)
	return f, err
}

// Rename renames (moves) oldpath to newpath. If newpath already exists and
// is not a directory, Rename replaces it. OS-specific restrictions may
// apply when oldpath and newpath are in different directories.
// Metadata for the file will be moved.
func (fs filesystem) Rename(oldpath, newpath string) error {
	return fs.Filesystem.Rename(oldpath, newpath)
}

// Remove removes the named file or directory, which must exist.
// Metadata for the file will be removed.
func (fs filesystem) Remove(filename string) error {
	if err := fs.store.Empty(filename); err != nil {
		return err
	}
	if err := fs.Filesystem.Remove(filename); err != nil {
		return err
	}
	return nil
}

// Close the filesystem and flush the metadata store
func (fs filesystem) Close() error {
	return nil
}
