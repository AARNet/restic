package mock

import (
	"io"
	"restic"

	"restic/errors"
)

// Backend implements a mock backend.
type Backend struct {
	CloseFn    func() error
	SaveFn     func(h restic.Handle, rd io.Reader) error
	LoadFn     func(h restic.Handle, length int, offset int64) (io.ReadCloser, error)
	StatFn     func(h restic.Handle) (restic.FileInfo, error)
	ListFn     func(restic.FileType, <-chan struct{}) <-chan string
	RemoveFn   func(h restic.Handle) error
	TestFn     func(h restic.Handle) (bool, error)
	DeleteFn   func() error
	LocationFn func() string
}

// Close the backend.
func (m *Backend) Close() error {
	if m.CloseFn == nil {
		return nil
	}

	return m.CloseFn()
}

// Location returns a location string.
func (m *Backend) Location() string {
	if m.LocationFn == nil {
		return ""
	}

	return m.LocationFn()
}

// Save data in the backend.
func (m *Backend) Save(h restic.Handle, rd io.Reader) error {
	if m.SaveFn == nil {
		return errors.New("not implemented")
	}

	return m.SaveFn(h, rd)
}

// Load loads data from the backend.
func (m *Backend) Load(h restic.Handle, length int, offset int64) (io.ReadCloser, error) {
	if m.LoadFn == nil {
		return nil, errors.New("not implemented")
	}

	return m.LoadFn(h, length, offset)
}

// Stat an object in the backend.
func (m *Backend) Stat(h restic.Handle) (restic.FileInfo, error) {
	if m.StatFn == nil {
		return restic.FileInfo{}, errors.New("not implemented")
	}

	return m.StatFn(h)
}

// List items of type t.
func (m *Backend) List(t restic.FileType, done <-chan struct{}) <-chan string {
	if m.ListFn == nil {
		ch := make(chan string)
		close(ch)
		return ch
	}

	return m.ListFn(t, done)
}

// Remove data from the backend.
func (m *Backend) Remove(h restic.Handle) error {
	if m.RemoveFn == nil {
		return errors.New("not implemented")
	}

	return m.RemoveFn(h)
}

// Test for the existence of a specific item.
func (m *Backend) Test(h restic.Handle) (bool, error) {
	if m.TestFn == nil {
		return false, errors.New("not implemented")
	}

	return m.TestFn(h)
}

// Delete all data.
func (m *Backend) Delete() error {
	if m.DeleteFn == nil {
		return errors.New("not implemented")
	}

	return m.DeleteFn()
}

// Make sure that Backend implements the backend interface.
var _ restic.Backend = &Backend{}
