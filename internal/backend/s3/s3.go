package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/restic/restic/internal/backend"
	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/restic"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"

	"github.com/restic/restic/internal/debug"
)

// Backend stores data on an S3 endpoint.
type Backend struct {
	client *minio.Client
	sem    *backend.Semaphore
	cfg    Config
	backend.Layout
}

// make sure that *Backend implements backend.Backend
var _ restic.Backend = &Backend{}

const defaultLayout = "default"

func open(cfg Config) (*Backend, error) {
	debug.Log("open, config %#v", cfg)

	if cfg.MaxRetries > 0 {
		minio.MaxRetry = int(cfg.MaxRetries)
	}

	var client *minio.Client
	var err error

	if cfg.KeyID == "" || cfg.Secret == "" {
		debug.Log("key/secret not found, trying to get them from IAM")
		creds := credentials.NewIAM("")
		client, err = minio.NewWithCredentials(cfg.Endpoint, creds, !cfg.UseHTTP, "")

		if err != nil {
			return nil, errors.Wrap(err, "minio.NewWithCredentials")
		}
	} else {
		debug.Log("key/secret found")
		client, err = minio.New(cfg.Endpoint, cfg.KeyID, cfg.Secret, !cfg.UseHTTP)

		if err != nil {
			return nil, errors.Wrap(err, "minio.New")
		}
	}

	sem, err := backend.NewSemaphore(cfg.Connections)
	if err != nil {
		return nil, err
	}

	be := &Backend{
		client: client,
		sem:    sem,
		cfg:    cfg,
	}

	client.SetCustomTransport(backend.Transport())

	l, err := backend.ParseLayout(be, cfg.Layout, defaultLayout, cfg.Prefix)
	if err != nil {
		return nil, err
	}

	be.Layout = l

	return be, nil
}

// Open opens the S3 backend at bucket and region. The bucket is created if it
// does not exist yet.
func Open(cfg Config) (restic.Backend, error) {
	return open(cfg)
}

// Create opens the S3 backend at bucket and region and creates the bucket if
// it does not exist yet.
func Create(cfg Config) (restic.Backend, error) {
	be, err := open(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}
	found, err := be.client.BucketExists(cfg.Bucket)
	if err != nil {
		debug.Log("BucketExists(%v) returned err %v", cfg.Bucket, err)
		return nil, errors.Wrap(err, "client.BucketExists")
	}

	if !found {
		// create new bucket with default ACL in default region
		err = be.client.MakeBucket(cfg.Bucket, "")
		if err != nil {
			return nil, errors.Wrap(err, "client.MakeBucket")
		}
	}

	return be, nil
}

// IsNotExist returns true if the error is caused by a not existing file.
func (be *Backend) IsNotExist(err error) bool {
	debug.Log("IsNotExist(%T, %#v)", err, err)
	if os.IsNotExist(errors.Cause(err)) {
		return true
	}

	if e, ok := errors.Cause(err).(minio.ErrorResponse); ok && e.Code == "NoSuchKey" {
		return true
	}

	return false
}

// Join combines path components with slashes.
func (be *Backend) Join(p ...string) string {
	return path.Join(p...)
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fi fileInfo) Name() string       { return fi.name }    // base name of the file
func (fi fileInfo) Size() int64        { return fi.size }    // length in bytes for regular files; system-dependent for others
func (fi fileInfo) Mode() os.FileMode  { return fi.mode }    // file mode bits
func (fi fileInfo) ModTime() time.Time { return fi.modTime } // modification time
func (fi fileInfo) IsDir() bool        { return fi.isDir }   // abbreviation for Mode().IsDir()
func (fi fileInfo) Sys() interface{}   { return nil }        // underlying data source (can return nil)

// ReadDir returns the entries for a directory.
func (be *Backend) ReadDir(dir string) (list []os.FileInfo, err error) {
	debug.Log("ReadDir(%v)", dir)

	// make sure dir ends with a slash
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}

	done := make(chan struct{})
	defer close(done)

	for obj := range be.client.ListObjects(be.cfg.Bucket, dir, false, done) {
		if obj.Key == "" {
			continue
		}

		name := strings.TrimPrefix(obj.Key, dir)
		if name == "" {
			return nil, errors.Errorf("invalid key name %v, removing prefix %v yielded empty string", obj.Key, dir)
		}
		entry := fileInfo{
			name:    name,
			size:    obj.Size,
			modTime: obj.LastModified,
		}

		if name[len(name)-1] == '/' {
			entry.isDir = true
			entry.mode = os.ModeDir | 0755
			entry.name = name[:len(name)-1]
		} else {
			entry.mode = 0644
		}

		list = append(list, entry)
	}

	return list, nil
}

// Location returns this backend's location (the bucket name).
func (be *Backend) Location() string {
	return be.Join(be.cfg.Bucket, be.cfg.Prefix)
}

// Path returns the path in the bucket that is used for this backend.
func (be *Backend) Path() string {
	return be.cfg.Prefix
}

// nopCloserFile wraps *os.File and overwrites the Close() method with method
// that does nothing. In addition, the method Len() is implemented, which
// returns the size of the file (filesize - current offset).
type nopCloserFile struct {
	*os.File
}

func (f nopCloserFile) Close() error {
	debug.Log("prevented Close()")
	return nil
}

// Len returns the remaining length of the file (filesize - current offset).
func (f nopCloserFile) Len() int {
	debug.Log("Len() called")
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	pos, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	size := fi.Size() - pos
	debug.Log("returning file size %v", size)
	return int(size)
}

type lenner interface {
	Len() int
	io.Reader
}

// nopCloserLenner wraps a lenner and overwrites the Close() method with method
// that does nothing. In addition, the method Size() is implemented, which
// returns the size of the file (filesize - current offset).
type nopCloserLenner struct {
	lenner
}

func (f *nopCloserLenner) Close() error {
	debug.Log("prevented Close()")
	return nil
}

// Save stores data in the backend at the handle.
func (be *Backend) Save(ctx context.Context, h restic.Handle, rd io.Reader) (err error) {
	debug.Log("Save %v", h)

	if err := h.Valid(); err != nil {
		return err
	}

	objName := be.Filename(h)

	// Check key does not already exist
	_, err = be.client.StatObject(be.cfg.Bucket, objName)
	if err == nil {
		debug.Log("%v already exists", h)
		return errors.New("key already exists")
	}

	// FIXME: This is a workaround once we move to minio-go 4.0.x this can be
	// removed and size can be directly provided.
	if f, ok := rd.(*os.File); ok {
		debug.Log("reader is %#T, using nopCloserFile{}", rd)
		rd = nopCloserFile{f}
	} else if l, ok := rd.(lenner); ok {
		debug.Log("reader is %#T, using nopCloserLenner{}", rd)
		rd = nopCloserLenner{l}
	} else {
		debug.Log("reader is %#T, no specific workaround enabled", rd)
	}

	be.sem.GetToken()
	debug.Log("PutObject(%v, %v)", be.cfg.Bucket, objName)
	n, err := be.client.PutObject(be.cfg.Bucket, objName, rd, "application/octet-stream")
	be.sem.ReleaseToken()

	debug.Log("%v -> %v bytes, err %#v: %v", objName, n, err, err)

	return errors.Wrap(err, "client.PutObject")
}

// wrapReader wraps an io.ReadCloser to run an additional function on Close.
type wrapReader struct {
	io.ReadCloser
	f func()
}

func (wr wrapReader) Close() error {
	err := wr.ReadCloser.Close()
	wr.f()
	return err
}

// Load returns a reader that yields the contents of the file at h at the
// given offset. If length is nonzero, only a portion of the file is
// returned. rd must be closed after use.
func (be *Backend) Load(ctx context.Context, h restic.Handle, length int, offset int64) (io.ReadCloser, error) {
	debug.Log("Load %v, length %v, offset %v from %v", h, length, offset, be.Filename(h))
	if err := h.Valid(); err != nil {
		return nil, err
	}

	if offset < 0 {
		return nil, errors.New("offset is negative")
	}

	if length < 0 {
		return nil, errors.Errorf("invalid length %d", length)
	}

	objName := be.Filename(h)

	byteRange := fmt.Sprintf("bytes=%d-", offset)
	if length > 0 {
		byteRange = fmt.Sprintf("bytes=%d-%d", offset, offset+int64(length)-1)
	}
	headers := minio.NewGetReqHeaders()
	headers.Add("Range", byteRange)

	be.sem.GetToken()
	debug.Log("Load(%v) send range %v", h, byteRange)

	coreClient := minio.Core{Client: be.client}
	rd, _, err := coreClient.GetObject(be.cfg.Bucket, objName, headers)
	if err != nil {
		be.sem.ReleaseToken()
		return nil, err
	}

	closeRd := wrapReader{
		ReadCloser: rd,
		f: func() {
			debug.Log("Close()")
			be.sem.ReleaseToken()
		},
	}

	return closeRd, err
}

// Stat returns information about a blob.
func (be *Backend) Stat(ctx context.Context, h restic.Handle) (bi restic.FileInfo, err error) {
	debug.Log("%v", h)

	objName := be.Filename(h)
	var obj *minio.Object

	obj, err = be.client.GetObject(be.cfg.Bucket, objName)
	if err != nil {
		debug.Log("GetObject() err %v", err)
		return restic.FileInfo{}, errors.Wrap(err, "client.GetObject")
	}

	// make sure that the object is closed properly.
	defer func() {
		e := obj.Close()
		if err == nil {
			err = errors.Wrap(e, "Close")
		}
	}()

	fi, err := obj.Stat()
	if err != nil {
		debug.Log("Stat() err %v", err)
		return restic.FileInfo{}, errors.Wrap(err, "Stat")
	}

	return restic.FileInfo{Size: fi.Size}, nil
}

// Test returns true if a blob of the given type and name exists in the backend.
func (be *Backend) Test(ctx context.Context, h restic.Handle) (bool, error) {
	found := false
	objName := be.Filename(h)
	_, err := be.client.StatObject(be.cfg.Bucket, objName)
	if err == nil {
		found = true
	}

	// If error, then not found
	return found, nil
}

// Remove removes the blob with the given name and type.
func (be *Backend) Remove(ctx context.Context, h restic.Handle) error {
	objName := be.Filename(h)
	err := be.client.RemoveObject(be.cfg.Bucket, objName)
	debug.Log("Remove(%v) at %v -> err %v", h, objName, err)

	if be.IsNotExist(err) {
		err = nil
	}

	return errors.Wrap(err, "client.RemoveObject")
}

// List returns a channel that yields all names of blobs of type t. A
// goroutine is started for this. If the channel done is closed, sending
// stops.
func (be *Backend) List(ctx context.Context, t restic.FileType) <-chan string {
	debug.Log("listing %v", t)
	ch := make(chan string)

	prefix := be.Dirname(restic.Handle{Type: t})

	// make sure prefix ends with a slash
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	listresp := be.client.ListObjects(be.cfg.Bucket, prefix, true, ctx.Done())

	go func() {
		defer close(ch)
		for obj := range listresp {
			m := strings.TrimPrefix(obj.Key, prefix)
			if m == "" {
				continue
			}

			select {
			case ch <- path.Base(m):
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}

// Remove keys for a specified backend type.
func (be *Backend) removeKeys(ctx context.Context, t restic.FileType) error {
	for key := range be.List(ctx, restic.DataFile) {
		err := be.Remove(ctx, restic.Handle{Type: restic.DataFile, Name: key})
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes all restic keys in the bucket. It will not remove the bucket itself.
func (be *Backend) Delete(ctx context.Context) error {
	alltypes := []restic.FileType{
		restic.DataFile,
		restic.KeyFile,
		restic.LockFile,
		restic.SnapshotFile,
		restic.IndexFile}

	for _, t := range alltypes {
		err := be.removeKeys(ctx, t)
		if err != nil {
			return nil
		}
	}

	return be.Remove(ctx, restic.Handle{Type: restic.ConfigFile})
}

// Close does nothing
func (be *Backend) Close() error { return nil }

// Rename moves a file based on the new layout l.
func (be *Backend) Rename(h restic.Handle, l backend.Layout) error {
	debug.Log("Rename %v to %v", h, l)
	oldname := be.Filename(h)
	newname := l.Filename(h)

	if oldname == newname {
		debug.Log("  %v is already renamed", newname)
		return nil
	}

	debug.Log("  %v -> %v", oldname, newname)

	src := minio.NewSourceInfo(be.cfg.Bucket, oldname, nil)

	dst, err := minio.NewDestinationInfo(be.cfg.Bucket, newname, nil, nil)
	if err != nil {
		return errors.Wrap(err, "NewDestinationInfo")
	}

	err = be.client.CopyObject(dst, src)
	if err != nil && be.IsNotExist(err) {
		debug.Log("copy failed: %v, seems to already have been renamed", err)
		return nil
	}

	if err != nil {
		debug.Log("copy failed: %v", err)
		return err
	}

	return be.client.RemoveObject(be.cfg.Bucket, oldname)
}
