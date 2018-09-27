package metafile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aedobrowolski/metafile/gob"
)

// Implement the Store methods.  Store does not handle its own persistence.

type buckets map[string]map[string][]byte

// store is the structure that holds metadata for a filesystem.
// This module is decoupled from the file system which may provide
// a validate method to check that paths (relative to root) exist.
type store struct {
	base  string            // fixed prefix used for bucket names (meta dir to root path)
	dirty bool              // something changed
	meta  buckets           // map: bucket -> map: keyword -> value
	valid func(string) bool // _, err := os.Stat("/path/to/whatever"); return !os.IsNotExist(err)
	codec Codec
}

// Codec is an encoder/decoder of arbitrary values into byte slices.
type Codec interface {
	// Encode returns a byte slice representing the value passed as an argument
	Encode(interface{}) ([]byte, error)
	// Decode converts the byte slice to data that it stores in the provided pointer
	Decode([]byte, interface{}) error
}

// newStore creates and initializes a metafile Store structure
func newStore(base string, valid func(string) bool) store {
	return store{base: base, dirty: false, valid: valid, meta: make(buckets), codec: &gob.Codec{}}
}

func prepare(s store, path string) (map[string][]byte, error) {
	path = Clean(path)
	if !valid(s, path) {
		return nil, ErrBadBucket
	}
	if b, ok := s.meta[path]; ok {
		return b, nil
	}
	b := make(map[string][]byte)
	s.meta[path] = b
	return b, nil
}

// valid test if a storage bucket is valid.  Any existing bucket is automatically valid.
// Otherwise the (typically expensive) function supplied by the file system must be used.
func valid(s store, path string) bool {
	if _, ok := s.meta[path]; ok {
		return true
	}
	if s.valid != nil {
		return s.valid(path)
	}
	return true
}

// Clean normalizes a bucket path. It does this in an OS dependent way by
// calling filepath.Clean and trimming slashes.
func Clean(path string) string {
	path = filepath.Clean(path) // output contains OS specific separator
	path = strings.Trim(path, "/\\")
	return path
}

func (s store) Put(path, key string, value interface{}) error {
	b, err := prepare(s, path)
	if err != nil {
		return err
	}

	b[key], err = s.codec.Encode(value)
	if err != nil {
		return err // very rare that value cannot be encoded
	}
	s.dirty = true
	return nil
}

func (s store) Get(path, key string, valuePtr interface{}) (ok bool, err error) {
	b, err := prepare(s, path)
	if err != nil {
		return false, err
	}
	value, ok := b[key] // get the encoded value
	if !ok {
		return false, nil
	}

	err = s.codec.Decode(value, valuePtr)
	if err != nil {
		// common error if the client supplies an incompatible pointer
		return false, fmt.Errorf("decoding `%s[%s]`: %s", path, key, err)
	}
	return true, nil
}

func (s store) Delete(path, key string) error {
	b, err := prepare(s, path)
	if err != nil {
		return err
	}

	delete(b, key)
	s.dirty = true
	return nil
}

// Methods that operate on buckets

func (s store) Empty(path string) error {
	path = Clean(path)
	if !valid(s, path) {
		return ErrBadBucket
	}
	delete(s.meta, path)
	path = Clean(path + "/")
	for k := range s.meta {
		if strings.HasPrefix(k, path) {
			delete(s.meta, k)
		}
	}
	s.dirty = true
	return nil
}

// emptyOne removes one bucket from storage - it is not recursive like Empty.
func (s store) emptyOne(path string) error {
	path = Clean(path)
	delete(s.meta, path)
	s.dirty = true
	return nil
}

// Move moves all metadata from PATH to TO including children of PATH.
// These must be the same type of object (e.g. two directories or two files)
//
// Returns ErrBadBucket if PATH is not a valid path.
func (s store) Move(path, to string) error {
	path = Clean(path)
	if !valid(s, path) {
		return ErrBadBucket
	}

	to = Clean(to)
	s.meta[to] = s.meta[path]
	delete(s.meta, path)
	s.dirty = true

	// handle child metadata
	path += "/"
	for k := range s.meta {
		if strings.HasPrefix(k, path) {
			s.meta[to+k[len(path)-1:]] = s.meta[k]
			delete(s.meta, k)
		}
	}

	return nil
}

// moveOne moves all metadata from PATH to TO
func (s store) moveOne(path, to string) error {
	path = Clean(path)
	to = Clean(to)
	s.meta[to] = s.meta[path]
	delete(s.meta, path)

	return nil
}
