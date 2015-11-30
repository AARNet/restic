package backend

import (
	"errors"
	"io"
)

// MockBackend implements a backend whose functions can be specified. This
// should only be used for tests.
type MockBackend struct {
	CloseFn     func() error
	CreateFn    func() (Blob, error)
	GetFn       func(Type, string) (io.ReadCloser, error)
	GetReaderFn func(Type, string, uint, uint) (io.ReadCloser, error)
	ListFn      func(Type, <-chan struct{}) <-chan string
	RemoveFn    func(Type, string) error
	TestFn      func(Type, string) (bool, error)
	DeleteFn    func() error
	LocationFn  func() string
}

func (m *MockBackend) Close() error {
	if m.CloseFn == nil {
		return nil
	}

	return m.CloseFn()
}

func (m *MockBackend) Location() string {
	if m.LocationFn == nil {
		return ""
	}

	return m.LocationFn()
}

func (m *MockBackend) Create() (Blob, error) {
	if m.CreateFn == nil {
		return nil, errors.New("not implemented")
	}

	return m.CreateFn()
}

func (m *MockBackend) Get(t Type, name string) (io.ReadCloser, error) {
	if m.GetFn == nil {
		return nil, errors.New("not implemented")
	}

	return m.GetFn(t, name)
}

func (m *MockBackend) GetReader(t Type, name string, offset, len uint) (io.ReadCloser, error) {
	if m.GetReaderFn == nil {
		return nil, errors.New("not implemented")
	}

	return m.GetReaderFn(t, name, offset, len)
}

func (m *MockBackend) List(t Type, done <-chan struct{}) <-chan string {
	if m.ListFn == nil {
		ch := make(chan string)
		close(ch)
		return ch
	}

	return m.ListFn(t, done)
}

func (m *MockBackend) Remove(t Type, name string) error {
	if m.RemoveFn == nil {
		return errors.New("not implemented")
	}

	return m.RemoveFn(t, name)
}

func (m *MockBackend) Test(t Type, name string) (bool, error) {
	if m.TestFn == nil {
		return false, errors.New("not implemented")
	}

	return m.TestFn(t, name)
}

func (m *MockBackend) Delete() error {
	if m.DeleteFn == nil {
		return errors.New("not implemented")
	}

	return m.DeleteFn()
}

// Make sure that MockBackend implements the backend interface.
var _ Backend = &MockBackend{}
