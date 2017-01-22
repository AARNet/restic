package restic

import "io"

// Backend is used to store and access data.
type Backend interface {
	// Location returns a string that describes the type and location of the
	// repository.
	Location() string

	// Test a boolean value whether a File with the name and type exists.
	Test(t FileType, name string) (bool, error)

	// Remove removes a File with type t and name.
	Remove(t FileType, name string) error

	// Close the backend
	Close() error

	// Load returns the data stored in the backend for h at the given offset
	// and saves it in p. Load has the same semantics as io.ReaderAt, except
	// that a negative offset is also allowed. In this case it references a
	// position relative to the end of the file (similar to Seek()).
	Load(h Handle, p []byte, off int64) (int, error)

	// Save stores the data in the backend under the given handle.
	Save(h Handle, rd io.Reader) error

	// Get returns a reader that yields the contents of the file at h at the
	// given offset. If length is nonzero, only a portion of the file is
	// returned. rd must be closed after use.
	Get(h Handle, length int, offset int64) (io.ReadCloser, error)

	// Stat returns information about the File identified by h.
	Stat(h Handle) (FileInfo, error)

	// List returns a channel that yields all names of files of type t in an
	// arbitrary order. A goroutine is started for this. If the channel done is
	// closed, sending stops.
	List(t FileType, done <-chan struct{}) <-chan string
}

// FileInfo is returned by Stat() and contains information about a file in the
// backend.
type FileInfo struct{ Size int64 }
