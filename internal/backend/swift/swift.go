package swift

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/restic/restic/internal/backend"
	"github.com/restic/restic/internal/debug"
	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/restic"

	"github.com/ncw/swift"
)

const connLimit = 10

// beSwift is a backend which stores the data on a swift endpoint.
type beSwift struct {
	conn      *swift.Connection
	sem       *backend.Semaphore
	container string // Container name
	prefix    string // Prefix of object names in the container
	backend.Layout
}

// ensure statically that *beSwift implements restic.Backend.
var _ restic.Backend = &beSwift{}

// Open opens the swift backend at a container in region. The container is
// created if it does not exist yet.
func Open(cfg Config, rt http.RoundTripper) (restic.Backend, error) {
	debug.Log("config %#v", cfg)

	sem, err := backend.NewSemaphore(cfg.Connections)
	if err != nil {
		return nil, err
	}

	be := &beSwift{
		conn: &swift.Connection{
			UserName:       cfg.UserName,
			Domain:         cfg.Domain,
			ApiKey:         cfg.APIKey,
			AuthUrl:        cfg.AuthURL,
			Region:         cfg.Region,
			Tenant:         cfg.Tenant,
			TenantId:       cfg.TenantID,
			TenantDomain:   cfg.TenantDomain,
			TrustId:        cfg.TrustID,
			StorageUrl:     cfg.StorageURL,
			AuthToken:      cfg.AuthToken,
			ConnectTimeout: time.Minute,
			Timeout:        time.Minute,

			Transport: rt,
		},
		sem:       sem,
		container: cfg.Container,
		prefix:    cfg.Prefix,
		Layout: &backend.DefaultLayout{
			Path: cfg.Prefix,
			Join: path.Join,
		},
	}

	// Authenticate if needed
	if !be.conn.Authenticated() {
		if err := be.conn.Authenticate(); err != nil {
			return nil, errors.Wrap(err, "conn.Authenticate")
		}
	}

	// Ensure container exists
	switch _, _, err := be.conn.Container(be.container); err {
	case nil:
		// Container exists

	case swift.ContainerNotFound:
		err = be.createContainer(cfg.DefaultContainerPolicy)
		if err != nil {
			return nil, errors.Wrap(err, "beSwift.createContainer")
		}

	default:
		return nil, errors.Wrap(err, "conn.Container")
	}

	return be, nil
}

func (be *beSwift) createContainer(policy string) error {
	var h swift.Headers
	if policy != "" {
		h = swift.Headers{
			"X-Storage-Policy": policy,
		}
	}

	return be.conn.ContainerCreate(be.container, h)
}

// Location returns this backend's location (the container name).
func (be *beSwift) Location() string {
	return be.container
}

// Load returns a reader that yields the contents of the file at h at the
// given offset. If length is nonzero, only a portion of the file is
// returned. rd must be closed after use.
func (be *beSwift) Load(ctx context.Context, h restic.Handle, length int, offset int64) (io.ReadCloser, error) {
	debug.Log("Load %v, length %v, offset %v", h, length, offset)
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

	headers := swift.Headers{}
	if offset > 0 {
		headers["Range"] = fmt.Sprintf("bytes=%d-", offset)
	}

	if length > 0 {
		headers["Range"] = fmt.Sprintf("bytes=%d-%d", offset, offset+int64(length)-1)
	}

	if _, ok := headers["Range"]; ok {
		debug.Log("Load(%v) send range %v", h, headers["Range"])
	}

	be.sem.GetToken()
	obj, _, err := be.conn.ObjectOpen(be.container, objName, false, headers)
	if err != nil {
		debug.Log("  err %v", err)
		be.sem.ReleaseToken()
		return nil, errors.Wrap(err, "conn.ObjectOpen")
	}

	return be.sem.ReleaseTokenOnClose(obj, nil), nil
}

// Save stores data in the backend at the handle.
func (be *beSwift) Save(ctx context.Context, h restic.Handle, rd io.Reader) (err error) {
	if err = h.Valid(); err != nil {
		return err
	}

	objName := be.Filename(h)

	debug.Log("Save %v at %v", h, objName)

	be.sem.GetToken()
	defer be.sem.ReleaseToken()

	// Check key does not already exist
	switch _, _, err = be.conn.Object(be.container, objName); err {
	case nil:
		debug.Log("%v already exists", h)
		return errors.New("key already exists")

	case swift.ObjectNotFound:
		// Ok, that's what we want

	default:
		return errors.Wrap(err, "conn.Object")
	}

	encoding := "binary/octet-stream"

	debug.Log("PutObject(%v, %v, %v)", be.container, objName, encoding)
	_, err = be.conn.ObjectPut(be.container, objName, rd, true, "", encoding, nil)
	debug.Log("%v, err %#v", objName, err)

	return errors.Wrap(err, "client.PutObject")
}

// Stat returns information about a blob.
func (be *beSwift) Stat(ctx context.Context, h restic.Handle) (bi restic.FileInfo, err error) {
	debug.Log("%v", h)

	objName := be.Filename(h)

	be.sem.GetToken()
	defer be.sem.ReleaseToken()

	obj, _, err := be.conn.Object(be.container, objName)
	if err != nil {
		debug.Log("Object() err %v", err)
		return restic.FileInfo{}, errors.Wrap(err, "conn.Object")
	}

	return restic.FileInfo{Size: obj.Bytes}, nil
}

// Test returns true if a blob of the given type and name exists in the backend.
func (be *beSwift) Test(ctx context.Context, h restic.Handle) (bool, error) {
	objName := be.Filename(h)

	be.sem.GetToken()
	defer be.sem.ReleaseToken()

	switch _, _, err := be.conn.Object(be.container, objName); err {
	case nil:
		return true, nil

	case swift.ObjectNotFound:
		return false, nil

	default:
		return false, errors.Wrap(err, "conn.Object")
	}
}

// Remove removes the blob with the given name and type.
func (be *beSwift) Remove(ctx context.Context, h restic.Handle) error {
	objName := be.Filename(h)

	be.sem.GetToken()
	defer be.sem.ReleaseToken()

	err := be.conn.ObjectDelete(be.container, objName)
	debug.Log("Remove(%v) -> err %v", h, err)
	return errors.Wrap(err, "conn.ObjectDelete")
}

// List returns a channel that yields all names of blobs of type t. A
// goroutine is started for this. If the channel done is closed, sending
// stops.
func (be *beSwift) List(ctx context.Context, t restic.FileType) <-chan string {
	debug.Log("listing %v", t)
	ch := make(chan string)

	prefix, _ := be.Basedir(t)
	prefix += "/"

	go func() {
		defer close(ch)

		err := be.conn.ObjectsWalk(be.container, &swift.ObjectsOpts{Prefix: prefix},
			func(opts *swift.ObjectsOpts) (interface{}, error) {
				be.sem.GetToken()
				newObjects, err := be.conn.ObjectNames(be.container, opts)
				be.sem.ReleaseToken()

				if err != nil {
					return nil, errors.Wrap(err, "conn.ObjectNames")
				}
				for _, obj := range newObjects {
					m := filepath.Base(strings.TrimPrefix(obj, prefix))
					if m == "" {
						continue
					}

					select {
					case ch <- m:
					case <-ctx.Done():
						return nil, io.EOF
					}
				}
				return newObjects, nil
			})

		if err != nil {
			debug.Log("ObjectsWalk returned error: %v", err)
		}
	}()

	return ch
}

// Remove keys for a specified backend type.
func (be *beSwift) removeKeys(ctx context.Context, t restic.FileType) error {
	for key := range be.List(ctx, t) {
		err := be.Remove(ctx, restic.Handle{Type: t, Name: key})
		if err != nil {
			return err
		}
	}

	return nil
}

// IsNotExist returns true if the error is caused by a not existing file.
func (be *beSwift) IsNotExist(err error) bool {
	if e, ok := errors.Cause(err).(*swift.Error); ok {
		return e.StatusCode == http.StatusNotFound
	}

	return false
}

// Delete removes all restic objects in the container.
// It will not remove the container itself.
func (be *beSwift) Delete(ctx context.Context) error {
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

	err := be.Remove(ctx, restic.Handle{Type: restic.ConfigFile})
	if err != nil && !be.IsNotExist(err) {
		return err
	}

	return nil
}

// Close does nothing
func (be *beSwift) Close() error { return nil }
