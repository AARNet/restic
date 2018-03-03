package test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/restic/restic/internal/restic"
	"github.com/restic/restic/internal/test"
)

func saveRandomFile(t testing.TB, be restic.Backend, length int) ([]byte, restic.Handle) {
	data := test.Random(23, length)
	id := restic.Hash(data)
	handle := restic.Handle{Type: restic.DataFile, Name: id.String()}
	err := be.Save(context.TODO(), handle, restic.NewByteReader(data))
	if err != nil {
		t.Fatalf("Save() error: %+v", err)
	}
	return data, handle
}

func remove(t testing.TB, be restic.Backend, h restic.Handle) {
	if err := be.Remove(context.TODO(), h); err != nil {
		t.Fatalf("Remove() returned error: %v", err)
	}
}

// BenchmarkLoadFile benchmarks the Load() method of a backend by
// loading a complete file.
func (s *Suite) BenchmarkLoadFile(t *testing.B) {
	be := s.open(t)
	defer s.close(t, be)

	length := 1<<24 + 2123
	data, handle := saveRandomFile(t, be, length)
	defer remove(t, be, handle)

	buf := make([]byte, length)

	t.SetBytes(int64(length))
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		var n int
		err := be.Load(context.TODO(), handle, 0, 0, func(rd io.Reader) (ierr error) {
			n, ierr = io.ReadFull(rd, buf)
			return ierr
		})
		if err != nil {
			t.Fatal(err)
		}

		if n != length {
			t.Fatalf("wrong number of bytes read: want %v, got %v", length, n)
		}

		if !bytes.Equal(data, buf) {
			t.Fatalf("wrong bytes returned")
		}
	}
}

// BenchmarkLoadPartialFile benchmarks the Load() method of a backend by
// loading the remainder of a file starting at a given offset.
func (s *Suite) BenchmarkLoadPartialFile(t *testing.B) {
	be := s.open(t)
	defer s.close(t, be)

	datalength := 1<<24 + 2123
	data, handle := saveRandomFile(t, be, datalength)
	defer remove(t, be, handle)

	testLength := datalength/4 + 555

	buf := make([]byte, testLength)

	t.SetBytes(int64(testLength))
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		var n int
		err := be.Load(context.TODO(), handle, testLength, 0, func(rd io.Reader) (ierr error) {
			n, ierr = io.ReadFull(rd, buf)
			return ierr
		})
		if err != nil {
			t.Fatal(err)
		}

		if n != testLength {
			t.Fatalf("wrong number of bytes read: want %v, got %v", testLength, n)
		}

		if !bytes.Equal(data[:testLength], buf) {
			t.Fatalf("wrong bytes returned")
		}

	}
}

// BenchmarkLoadPartialFileOffset benchmarks the Load() method of a
// backend by loading a number of bytes of a file starting at a given offset.
func (s *Suite) BenchmarkLoadPartialFileOffset(t *testing.B) {
	be := s.open(t)
	defer s.close(t, be)

	datalength := 1<<24 + 2123
	data, handle := saveRandomFile(t, be, datalength)
	defer remove(t, be, handle)

	testLength := datalength/4 + 555
	testOffset := 8273

	buf := make([]byte, testLength)

	t.SetBytes(int64(testLength))
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		var n int
		err := be.Load(context.TODO(), handle, testLength, int64(testOffset), func(rd io.Reader) (ierr error) {
			n, ierr = io.ReadFull(rd, buf)
			return ierr
		})
		if err != nil {
			t.Fatal(err)
		}

		if n != testLength {
			t.Fatalf("wrong number of bytes read: want %v, got %v", testLength, n)
		}

		if !bytes.Equal(data[testOffset:testOffset+testLength], buf) {
			t.Fatalf("wrong bytes returned")
		}

	}
}

// BenchmarkSave benchmarks the Save() method of a backend.
func (s *Suite) BenchmarkSave(t *testing.B) {
	be := s.open(t)
	defer s.close(t, be)

	length := 1<<24 + 2123
	data := test.Random(23, length)
	id := restic.Hash(data)
	handle := restic.Handle{Type: restic.DataFile, Name: id.String()}

	rd := restic.NewByteReader(data)
	t.SetBytes(int64(length))
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		if err := be.Save(context.TODO(), handle, rd); err != nil {
			t.Fatal(err)
		}

		if err := be.Remove(context.TODO(), handle); err != nil {
			t.Fatal(err)
		}
	}
}
