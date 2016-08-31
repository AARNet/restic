package restic

// Repository stores data in a backend. It provides high-level functions and
// transparently encrypts/decrypts data.
type Repository interface {

	// Backend returns the backend used by the repository
	Backend() Backend

	SetIndex(Index)

	Index() Index
	SaveFullIndex() error

	SaveJSON(BlobType, interface{}) (ID, error)
	SaveUnpacked(FileType, []byte) (ID, error)

	Config() Config

	SaveAndEncrypt(BlobType, []byte, *ID) (ID, error)
	SaveJSONUnpacked(FileType, interface{}) (ID, error)
	SaveIndex() error

	LoadJSONPack(BlobType, ID, interface{}) error
	LoadJSONUnpacked(FileType, ID, interface{}) error
	LoadBlob(ID, BlobType, []byte) ([]byte, error)

	LookupBlobSize(ID, BlobType) (uint, error)

	List(FileType, <-chan struct{}) <-chan ID
	ListPack(ID) ([]Blob, int64, error)

	Flush() error
}

// Deleter removes all data stored in a backend/repo.
type Deleter interface {
	Delete() error
}

// Lister allows listing files in a backend.
type Lister interface {
	List(FileType, <-chan struct{}) <-chan string
}

// Index keeps track of the blobs are stored within files.
type Index interface {
	Has(ID, BlobType) bool
	Lookup(ID, BlobType) ([]PackedBlob, error)
}
