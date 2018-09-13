package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

func main() {
	for _, fs := range []billy.Filesystem{osfs.New(os.TempDir()), memfs.New()} {
		f, err := CreateWrite(fs)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Created file `%s` in `%s`.\n", f.Name(), fs.Root())
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}

// CreateWrite tests the behavior of a filesystem FS.
func CreateWrite(fs billy.Filesystem) (billy.File, error) {
	err := fs.MkdirAll("billytests", 0777)
	if err != nil {
		return nil, err
	}
	f, err := fs.Create("billytests/mytestfile.txt")
	if err != nil {
		return nil, err
	}
	buf := []byte("Hello World!\n")
	n, err := f.Write(buf)
	if n != len(buf) || err != nil {
		return nil, err
	}

	return f, nil
}
