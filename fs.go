// Metafile is a file system abstraction that adds persisted user-defined metadata
// to file system objects.
package metafile

import (
	"errors"

	"gopkg.in/src-d/go-billy.v4"
)

// Static errors
var (
	ErrBadBucket   = errors.New("bucket does not exist")
	ErrNeedPointer = errors.New("attempt to get into a non-pointer")
	ErrDecoding    = errors.New("pointer incompatible with the stored value")

	ErrReadOnly        = billy.ErrReadOnly
	ErrNotSupported    = billy.ErrNotSupported
	ErrCrossedBoundary = billy.ErrCrossedBoundary
)

// Filesystem represents a file system abstraction with a key-value store for each object.
type Filesystem interface {
	billy.Filesystem
	Store

	// Close the filesystem and persist the key-value store
	Close() error
}

// Store represents a key-value store associated with file system paths
type Store interface {
	// Put stores a value against a key in a storage bucket.
	//
	// Returns ErrBadBucket if the storage bucket does not exist.
	Put(bucket, key string, value interface{}) error

	// Get retrieves a value for a key in a storage bucket and returns
	// it in pointer.  If the key does not exist ok will be false and
	// the pointer value will be unchanged.
	//
	// Returns ErrBadBucket if the storage bucket does not exist.
	// Returns ErrDecoding if pointer is incompatible with the stored value.
	Get(bucket, key string, pointer *interface{}) (ok bool, err error)

	// Delete deletes a key from a storage bucket.
	//
	// Returns ErrBadBucket if the storage bucket does not exist.
	Delete(bucket, key string) error

	// Empty removes a storage bucket and all nested buckets.  For example
	// if you remove "foo" then "foo/bar" and "foo/baz" will also be
	// removed. However "foobar" will not be affected.
	//
	// Returns ErrBadBucket if the storage bucket does not exist.
	Empty(bucket string) error
}
