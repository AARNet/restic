package local

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"restic/backend"
	"restic/debug"
)

// Local is a backend in a local directory.
type Local struct {
	p string
}

func paths(dir string) []string {
	return []string{
		dir,
		filepath.Join(dir, backend.Paths.Data),
		filepath.Join(dir, backend.Paths.Snapshots),
		filepath.Join(dir, backend.Paths.Index),
		filepath.Join(dir, backend.Paths.Locks),
		filepath.Join(dir, backend.Paths.Keys),
		filepath.Join(dir, backend.Paths.Temp),
	}
}

// Open opens the local backend as specified by config.
func Open(dir string) (*Local, error) {
	// test if all necessary dirs are there
	for _, d := range paths(dir) {
		if _, err := os.Stat(d); err != nil {
			return nil, fmt.Errorf("%s does not exist", d)
		}
	}

	return &Local{p: dir}, nil
}

// Create creates all the necessary files and directories for a new local
// backend at dir. Afterwards a new config blob should be created.
func Create(dir string) (*Local, error) {
	// test if config file already exists
	_, err := os.Lstat(filepath.Join(dir, backend.Paths.Config))
	if err == nil {
		return nil, errors.New("config file already exists")
	}

	// create paths for data, refs and temp
	for _, d := range paths(dir) {
		err := os.MkdirAll(d, backend.Modes.Dir)
		if err != nil {
			return nil, err
		}
	}

	// open backend
	return Open(dir)
}

// Location returns this backend's location (the directory name).
func (b *Local) Location() string {
	return b.p
}

// Construct path for given Type and name.
func filename(base string, t backend.Type, name string) string {
	if t == backend.Config {
		return filepath.Join(base, "config")
	}

	return filepath.Join(dirname(base, t, name), name)
}

// Construct directory for given Type.
func dirname(base string, t backend.Type, name string) string {
	var n string
	switch t {
	case backend.Data:
		n = backend.Paths.Data
		if len(name) > 2 {
			n = filepath.Join(n, name[:2])
		}
	case backend.Snapshot:
		n = backend.Paths.Snapshots
	case backend.Index:
		n = backend.Paths.Index
	case backend.Lock:
		n = backend.Paths.Locks
	case backend.Key:
		n = backend.Paths.Keys
	}
	return filepath.Join(base, n)
}

// Load returns the data stored in the backend for h at the given offset
// and saves it in p. Load has the same semantics as io.ReaderAt.
func (b *Local) Load(h backend.Handle, p []byte, off int64) (n int, err error) {
	if err := h.Valid(); err != nil {
		return 0, err
	}

	f, err := os.Open(filename(b.p, h.Type, h.Name))
	if err != nil {
		return 0, err
	}

	defer func() {
		e := f.Close()
		if err == nil && e != nil {
			err = e
		}
	}()

	if off > 0 {
		_, err = f.Seek(off, 0)
		if err != nil {
			return 0, err
		}
	}

	return io.ReadFull(f, p)
}

// writeToTempfile saves p into a tempfile in tempdir.
func writeToTempfile(tempdir string, p []byte) (filename string, err error) {
	tmpfile, err := ioutil.TempFile(tempdir, "temp-")
	if err != nil {
		return "", err
	}

	n, err := tmpfile.Write(p)
	if err != nil {
		return "", err
	}

	if n != len(p) {
		return "", errors.New("not all bytes writen")
	}

	if err = tmpfile.Sync(); err != nil {
		return "", err
	}

	err = tmpfile.Close()
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

// Save stores data in the backend at the handle.
func (b *Local) Save(h backend.Handle, p []byte) (err error) {
	if err := h.Valid(); err != nil {
		return err
	}

	tmpfile, err := writeToTempfile(filepath.Join(b.p, backend.Paths.Temp), p)
	debug.Log("local.Save", "saved %v (%d bytes) to %v", h, len(p), tmpfile)

	filename := filename(b.p, h.Type, h.Name)

	// test if new path already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("Rename(): file %v already exists", filename)
	}

	// create directories if necessary, ignore errors
	if h.Type == backend.Data {
		err = os.MkdirAll(filepath.Dir(filename), backend.Modes.Dir)
		if err != nil {
			return err
		}
	}

	err = os.Rename(tmpfile, filename)
	debug.Log("local.Save", "save %v: rename %v -> %v: %v",
		h, filepath.Base(tmpfile), filepath.Base(filename), err)

	if err != nil {
		return err
	}

	// set mode to read-only
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

	return setNewFileMode(filename, fi)
}

// Stat returns information about a blob.
func (b *Local) Stat(h backend.Handle) (backend.BlobInfo, error) {
	if err := h.Valid(); err != nil {
		return backend.BlobInfo{}, err
	}

	fi, err := os.Stat(filename(b.p, h.Type, h.Name))
	if err != nil {
		return backend.BlobInfo{}, err
	}

	return backend.BlobInfo{Size: fi.Size()}, nil
}

// Test returns true if a blob of the given type and name exists in the backend.
func (b *Local) Test(t backend.Type, name string) (bool, error) {
	_, err := os.Stat(filename(b.p, t, name))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Remove removes the blob with the given name and type.
func (b *Local) Remove(t backend.Type, name string) error {
	fn := filename(b.p, t, name)

	// reset read-only flag
	err := os.Chmod(fn, 0666)
	if err != nil {
		return err
	}

	return os.Remove(fn)
}

// List returns a channel that yields all names of blobs of type t. A
// goroutine is started for this. If the channel done is closed, sending
// stops.
func (b *Local) List(t backend.Type, done <-chan struct{}) <-chan string {
	var pattern string
	if t == backend.Data {
		pattern = filepath.Join(dirname(b.p, t, ""), "*", "*")
	} else {
		pattern = filepath.Join(dirname(b.p, t, ""), "*")
	}

	ch := make(chan string)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		close(ch)
		return ch
	}

	for i := range matches {
		matches[i] = filepath.Base(matches[i])
	}

	go func() {
		defer close(ch)
		for _, m := range matches {
			if m == "" {
				continue
			}

			select {
			case ch <- m:
			case <-done:
				return
			}
		}
	}()

	return ch
}

// Delete removes the repository and all files.
func (b *Local) Delete() error {
	return os.RemoveAll(b.p)
}

// Close closes all open files.
func (b *Local) Close() error {
	// this does not need to do anything, all open files are closed within the
	// same function.
	return nil
}
