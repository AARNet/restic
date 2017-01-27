package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"restic"
	"strings"

	"restic/debug"
	"restic/errors"

	"restic/backend"
)

const connLimit = 10

// make sure the rest backend implements restic.Backend
var _ restic.Backend = &restBackend{}

// restPath returns the path to the given resource.
func restPath(url *url.URL, h restic.Handle) string {
	u := *url

	var dir string

	switch h.Type {
	case restic.ConfigFile:
		dir = ""
		h.Name = "config"
	case restic.DataFile:
		dir = backend.Paths.Data
	case restic.SnapshotFile:
		dir = backend.Paths.Snapshots
	case restic.IndexFile:
		dir = backend.Paths.Index
	case restic.LockFile:
		dir = backend.Paths.Locks
	case restic.KeyFile:
		dir = backend.Paths.Keys
	default:
		dir = string(h.Type)
	}

	u.Path = path.Join(url.Path, dir, h.Name)

	return u.String()
}

type restBackend struct {
	url      *url.URL
	connChan chan struct{}
	client   http.Client
}

// Open opens the REST backend with the given config.
func Open(cfg Config) (restic.Backend, error) {
	connChan := make(chan struct{}, connLimit)
	for i := 0; i < connLimit; i++ {
		connChan <- struct{}{}
	}
	tr := &http.Transport{MaxIdleConnsPerHost: connLimit}
	client := http.Client{Transport: tr}

	return &restBackend{url: cfg.URL, connChan: connChan, client: client}, nil
}

// Location returns this backend's location (the server's URL).
func (b *restBackend) Location() string {
	return b.url.String()
}

// Save stores data in the backend at the handle.
func (b *restBackend) Save(h restic.Handle, rd io.Reader) (err error) {
	if err := h.Valid(); err != nil {
		return err
	}

	// make sure that client.Post() cannot close the reader by wrapping it in
	// backend.Closer, which has a noop method.
	rd = backend.Closer{Reader: rd}

	<-b.connChan
	resp, err := b.client.Post(restPath(b.url, h), "binary/octet-stream", rd)
	b.connChan <- struct{}{}

	if resp != nil {
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			e := resp.Body.Close()

			if err == nil {
				err = errors.Wrap(e, "Close")
			}
		}()
	}

	if err != nil {
		return errors.Wrap(err, "client.Post")
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("unexpected HTTP response code %v", resp.StatusCode)
	}

	return nil
}

// Load returns a reader that yields the contents of the file at h at the
// given offset. If length is nonzero, only a portion of the file is
// returned. rd must be closed after use.
func (b *restBackend) Load(h restic.Handle, length int, offset int64) (io.ReadCloser, error) {
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

	req, err := http.NewRequest("GET", restPath(b.url, h), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	byteRange := fmt.Sprintf("bytes=%d-", offset)
	if length > 0 {
		byteRange = fmt.Sprintf("bytes=%d-%d", offset, offset+int64(length)-1)
	}
	req.Header.Add("Range", byteRange)
	debug.Log("Load(%v) send range %v", h, byteRange)

	<-b.connChan
	resp, err := b.client.Do(req)
	b.connChan <- struct{}{}

	if err != nil {
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
		return nil, errors.Wrap(err, "client.Do")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		return nil, errors.Errorf("unexpected HTTP response code %v", resp.StatusCode)
	}

	return resp.Body, nil
}

// Stat returns information about a blob.
func (b *restBackend) Stat(h restic.Handle) (restic.FileInfo, error) {
	if err := h.Valid(); err != nil {
		return restic.FileInfo{}, err
	}

	<-b.connChan
	resp, err := b.client.Head(restPath(b.url, h))
	b.connChan <- struct{}{}
	if err != nil {
		return restic.FileInfo{}, errors.Wrap(err, "client.Head")
	}

	io.Copy(ioutil.Discard, resp.Body)
	if err = resp.Body.Close(); err != nil {
		return restic.FileInfo{}, errors.Wrap(err, "Close")
	}

	if resp.StatusCode != 200 {
		return restic.FileInfo{}, errors.Errorf("unexpected HTTP response code %v", resp.StatusCode)
	}

	if resp.ContentLength < 0 {
		return restic.FileInfo{}, errors.New("negative content length")
	}

	bi := restic.FileInfo{
		Size: resp.ContentLength,
	}

	return bi, nil
}

// Test returns true if a blob of the given type and name exists in the backend.
func (b *restBackend) Test(h restic.Handle) (bool, error) {
	_, err := b.Stat(h)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// Remove removes the blob with the given name and type.
func (b *restBackend) Remove(h restic.Handle) error {
	if err := h.Valid(); err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", restPath(b.url, h), nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}
	<-b.connChan
	resp, err := b.client.Do(req)
	b.connChan <- struct{}{}

	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	if resp.StatusCode != 200 {
		return errors.New("blob not removed")
	}

	io.Copy(ioutil.Discard, resp.Body)
	return resp.Body.Close()
}

// List returns a channel that yields all names of blobs of type t. A
// goroutine is started for this. If the channel done is closed, sending
// stops.
func (b *restBackend) List(t restic.FileType, done <-chan struct{}) <-chan string {
	ch := make(chan string)

	url := restPath(b.url, restic.Handle{Type: t})
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	<-b.connChan
	resp, err := b.client.Get(url)
	b.connChan <- struct{}{}

	if resp != nil {
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			e := resp.Body.Close()

			if err == nil {
				err = errors.Wrap(e, "Close")
			}
		}()
	}

	if err != nil {
		close(ch)
		return ch
	}

	dec := json.NewDecoder(resp.Body)
	var list []string
	if err = dec.Decode(&list); err != nil {
		close(ch)
		return ch
	}

	go func() {
		defer close(ch)
		for _, m := range list {
			select {
			case ch <- m:
			case <-done:
				return
			}
		}
	}()

	return ch
}

// Close closes all open files.
func (b *restBackend) Close() error {
	// this does not need to do anything, all open files are closed within the
	// same function.
	return nil
}
