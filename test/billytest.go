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
		// open fails for memfs only: cannot open directory: \billytests
		// fs.Stat("dirname") works nicely for both
		// what happens when you stat a non-existing file or dir?
		dstat, err := StatDirectory(fs)
		if err != nil || !dstat.IsDir() {
			panic(err)
		}

		fmt.Printf("  Created dir in `%s`.\n", fs.Root())
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

// StatDirectory Test the ability to open a directory as a file
func StatDirectory(fs billy.Filesystem) (os.FileInfo, error) {
	// Create a directory
	err := fs.MkdirAll("billytests", 0777)
	if err != nil {
		return nil, err
	}
	// Open the directory
	dstat, err := fs.Stat("billytests") // tried with trailing slash as well
	if err != nil {
		return nil, err
	}

	return dstat, nil
}
