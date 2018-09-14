package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

func main() {
	// test both os and memory based filesystems
	filesystems := []billy.Filesystem{osfs.New(os.TempDir()), memfs.New()}

	fmt.Println("\n=== File creation ===")
	for i, fs := range filesystems {
		fmt.Println("  Using fs ", i)
		f, err := CreateWrite(fs)
		if err != nil {
			panic(err)
		}

		fmt.Printf("  Created file `%s` in `%s`.\n", f.Name(), fs.Root())
		if err := f.Close(); err != nil {
			panic(err)
		}
	}

	fmt.Println("\n=== Directory creation ===")
	for i, fs := range filesystems {
		fmt.Println("  Using fs ", i)
		d, err := OpenDirectory(fs)
		if err != nil {
			panic(err) // fails for memfs: cannot open directory: \billytests
		}
		s, _ := fs.Stat(d.Name())
		if err != nil || !s.IsDir() {
			panic(err)
		}

		fmt.Printf("  Created dir `%s` in `%s`.\n", d.Name(), fs.Root())
		if err := d.Close(); err != nil {
			panic(err)
		}
	}
}

// CreateWrite tests the behavior of a filesystem FS.
func CreateWrite(fs billy.Filesystem) (billy.File, error) {
	// Create a directory
	err := fs.MkdirAll("billytests", 0777)
	if err != nil {
		return nil, err
	}
	// Create a new file (or truncate existing)
	f, err := fs.Create("billytests/mytestfile.txt")
	if err != nil {
		return nil, err
	}
	// Write to the file
	buf := []byte("Hello World!\n")
	n, err := f.Write(buf)
	if n != len(buf) || err != nil {
		return nil, err
	}

	return f, nil
}

// OpenDirectory Test the ability to open a directory as a file
func OpenDirectory(fs billy.Filesystem) (billy.File, error) {
	// Create a directory
	err := fs.MkdirAll("billytests", 0777)
	if err != nil {
		return nil, err
	}
	// Open the directory
	f, err := fs.Open("billytests/")
	if err != nil {
		return nil, err
	}

	return f, nil
}
