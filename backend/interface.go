package backend

// Type is the type of a Blob.
type Type string

// These are the different data types a backend can store.
const (
	Data     Type = "data"
	Key           = "key"
	Lock          = "lock"
	Snapshot      = "snapshot"
	Index         = "index"
	Config        = "config"
)

// Backend is used to store and access data.
type Backend interface {
	// Location returns a string that describes the type and location of the
	// repository.
	Location() string

	// Test a boolean value whether a Blob with the name and type exists.
	Test(t Type, name string) (bool, error)

	// Remove removes a Blob with type t and name.
	Remove(t Type, name string) error

	// Close the backend
	Close() error

	Lister

	// Load returns the data stored in the backend for h at the given offset
	// and saves it in p. Load has the same semantics as io.ReaderAt.
	Load(h Handle, p []byte, off int64) (int, error)

	// Save stores the data in the backend under the given handle.
	Save(h Handle, p []byte) error

	// Stat returns information about the blob identified by h.
	Stat(h Handle) (BlobInfo, error)
}

// Lister implements listing data items stored in a backend.
type Lister interface {
	// List returns a channel that yields all names of blobs of type t in an
	// arbitrary order. A goroutine is started for this. If the channel done is
	// closed, sending stops.
	List(t Type, done <-chan struct{}) <-chan string
}

// Deleter are backends that allow to self-delete all content stored in them.
type Deleter interface {
	// Delete the complete repository.
	Delete() error
}

// BlobInfo is returned by Stat() and contains information about a stored blob.
type BlobInfo struct {
	Size int64
}
