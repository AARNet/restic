package backend_test

import (
	"bytes"
	"math/rand"
	"testing"

	"restic/backend"
	"restic/backend/mem"
	. "restic/test"
)

const KiB = 1 << 10
const MiB = 1 << 20

func TestLoadAll(t *testing.T) {
	b := mem.New()

	for i := 0; i < 20; i++ {
		data := Random(23+i, rand.Intn(MiB)+500*KiB)

		id := backend.Hash(data)
		err := b.Save(backend.Handle{Name: id.String(), Type: backend.Data}, data)
		OK(t, err)

		buf, err := backend.LoadAll(b, backend.Handle{Type: backend.Data, Name: id.String()}, nil)
		OK(t, err)

		if len(buf) != len(data) {
			t.Errorf("length of returned buffer does not match, want %d, got %d", len(data), len(buf))
			continue
		}

		if !bytes.Equal(buf, data) {
			t.Errorf("wrong data returned")
			continue
		}
	}
}

func TestLoadSmallBuffer(t *testing.T) {
	b := mem.New()

	for i := 0; i < 20; i++ {
		data := Random(23+i, rand.Intn(MiB)+500*KiB)

		id := backend.Hash(data)
		err := b.Save(backend.Handle{Name: id.String(), Type: backend.Data}, data)
		OK(t, err)

		buf := make([]byte, len(data)-23)
		buf, err = backend.LoadAll(b, backend.Handle{Type: backend.Data, Name: id.String()}, buf)
		OK(t, err)

		if len(buf) != len(data) {
			t.Errorf("length of returned buffer does not match, want %d, got %d", len(data), len(buf))
			continue
		}

		if !bytes.Equal(buf, data) {
			t.Errorf("wrong data returned")
			continue
		}
	}
}

func TestLoadLargeBuffer(t *testing.T) {
	b := mem.New()

	for i := 0; i < 20; i++ {
		data := Random(23+i, rand.Intn(MiB)+500*KiB)

		id := backend.Hash(data)
		err := b.Save(backend.Handle{Name: id.String(), Type: backend.Data}, data)
		OK(t, err)

		buf := make([]byte, len(data)+100)
		buf, err = backend.LoadAll(b, backend.Handle{Type: backend.Data, Name: id.String()}, buf)
		OK(t, err)

		if len(buf) != len(data) {
			t.Errorf("length of returned buffer does not match, want %d, got %d", len(data), len(buf))
			continue
		}

		if !bytes.Equal(buf, data) {
			t.Errorf("wrong data returned")
			continue
		}
	}
}
