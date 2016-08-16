package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"restic/backend"
	"restic/crypto"
)

// BlobType specifies what a blob stored in a pack is.
type BlobType uint8

// These are the blob types that can be stored in a pack.
const (
	Invalid BlobType = iota
	Data
	Tree
)

func (t BlobType) String() string {
	switch t {
	case Data:
		return "data"
	case Tree:
		return "tree"
	}

	return fmt.Sprintf("<BlobType %d>", t)
}

// MarshalJSON encodes the BlobType into JSON.
func (t BlobType) MarshalJSON() ([]byte, error) {
	switch t {
	case Data:
		return []byte(`"data"`), nil
	case Tree:
		return []byte(`"tree"`), nil
	}

	return nil, errors.New("unknown blob type")
}

// UnmarshalJSON decodes the BlobType from JSON.
func (t *BlobType) UnmarshalJSON(buf []byte) error {
	switch string(buf) {
	case `"data"`:
		*t = Data
	case `"tree"`:
		*t = Tree
	default:
		return errors.New("unknown blob type")
	}

	return nil
}

// Blob is a blob within a pack.
type Blob struct {
	Type   BlobType
	Length uint
	ID     backend.ID
	Offset uint
}

// Packer is used to create a new Pack.
type Packer struct {
	blobs []Blob

	bytes uint
	k     *crypto.Key
	wr    io.Writer

	m sync.Mutex
}

// NewPacker returns a new Packer that can be used to pack blobs
// together. If wr is nil, a bytes.Buffer is used.
func NewPacker(k *crypto.Key, wr io.Writer) *Packer {
	if wr == nil {
		wr = bytes.NewBuffer(nil)
	}
	return &Packer{k: k, wr: wr}
}

// Add saves the data read from rd as a new blob to the packer. Returned is the
// number of bytes written to the pack.
func (p *Packer) Add(t BlobType, id backend.ID, data []byte) (int, error) {
	p.m.Lock()
	defer p.m.Unlock()

	c := Blob{Type: t, ID: id}

	n, err := p.wr.Write(data)
	c.Length = uint(n)
	c.Offset = p.bytes
	p.bytes += uint(n)
	p.blobs = append(p.blobs, c)

	return n, err
}

var entrySize = uint(binary.Size(BlobType(0)) + binary.Size(uint32(0)) + backend.IDSize)

// headerEntry is used with encoding/binary to read and write header entries
type headerEntry struct {
	Type   uint8
	Length uint32
	ID     [backend.IDSize]byte
}

// Finalize writes the header for all added blobs and finalizes the pack.
// Returned are the number of bytes written, including the header. If the
// underlying writer implements io.Closer, it is closed.
func (p *Packer) Finalize() (uint, error) {
	p.m.Lock()
	defer p.m.Unlock()

	bytesWritten := p.bytes

	hdrBuf := bytes.NewBuffer(nil)
	bytesHeader, err := p.writeHeader(hdrBuf)
	if err != nil {
		return 0, err
	}

	encryptedHeader, err := crypto.Encrypt(p.k, nil, hdrBuf.Bytes())
	if err != nil {
		return 0, err
	}

	// append the header
	n, err := p.wr.Write(encryptedHeader)
	if err != nil {
		return 0, err
	}

	hdrBytes := bytesHeader + crypto.Extension
	if uint(n) != hdrBytes {
		return 0, errors.New("wrong number of bytes written")
	}

	bytesWritten += hdrBytes

	// write length
	err = binary.Write(p.wr, binary.LittleEndian, uint32(uint(len(p.blobs))*entrySize+crypto.Extension))
	if err != nil {
		return 0, err
	}
	bytesWritten += uint(binary.Size(uint32(0)))

	p.bytes = uint(bytesWritten)

	if w, ok := p.wr.(io.Closer); ok {
		return bytesWritten, w.Close()
	}

	return bytesWritten, nil
}

// writeHeader constructs and writes the header to wr.
func (p *Packer) writeHeader(wr io.Writer) (bytesWritten uint, err error) {
	for _, b := range p.blobs {
		entry := headerEntry{
			Length: uint32(b.Length),
			ID:     b.ID,
		}

		switch b.Type {
		case Data:
			entry.Type = 0
		case Tree:
			entry.Type = 1
		default:
			return 0, fmt.Errorf("invalid blob type %v", b.Type)
		}

		err := binary.Write(wr, binary.LittleEndian, entry)
		if err != nil {
			return bytesWritten, err
		}

		bytesWritten += entrySize
	}

	return
}

// Size returns the number of bytes written so far.
func (p *Packer) Size() uint {
	p.m.Lock()
	defer p.m.Unlock()

	return p.bytes
}

// Count returns the number of blobs in this packer.
func (p *Packer) Count() int {
	p.m.Lock()
	defer p.m.Unlock()

	return len(p.blobs)
}

// Blobs returns the slice of blobs that have been written.
func (p *Packer) Blobs() []Blob {
	p.m.Lock()
	defer p.m.Unlock()

	return p.blobs
}

// Writer return the underlying writer.
func (p *Packer) Writer() io.Writer {
	return p.wr
}

func (p *Packer) String() string {
	return fmt.Sprintf("<Packer %d blobs, %d bytes>", len(p.blobs), p.bytes)
}

// Unpacker is used to read individual blobs from a pack.
type Unpacker struct {
	rd      io.ReadSeeker
	Entries []Blob
	k       *crypto.Key
}

const preloadHeaderSize = 2048

// NewUnpacker returns a pointer to Unpacker which can be used to read
// individual Blobs from a pack.
func NewUnpacker(k *crypto.Key, ldr Loader) (*Unpacker, error) {
	var err error

	// read the last 2048 byte, this will mostly be enough for the header, so
	// we do not need another round trip.
	buf := make([]byte, preloadHeaderSize)
	n, err := ldr.Load(buf, -int64(len(buf)))

	if err == io.ErrUnexpectedEOF {
		err = nil
		buf = buf[:n]
	}

	if err != nil {
		return nil, fmt.Errorf("Load at -%d failed: %v", len(buf), err)
	}
	buf = buf[:n]

	bs := binary.Size(uint32(0))
	p := len(buf) - bs

	// read the length from the end of the buffer
	length := int(binary.LittleEndian.Uint32(buf[p : p+bs]))
	buf = buf[:p]

	// if the header is longer than the preloaded buffer, call the loader again.
	if length > len(buf) {
		buf = make([]byte, length)
		n, err := ldr.Load(buf, -int64(len(buf)+bs))
		if err != nil {
			return nil, fmt.Errorf("Load at -%d failed: %v", len(buf), err)
		}
		buf = buf[:n]
	}

	buf = buf[len(buf)-length:]

	// read header
	hdr, err := crypto.Decrypt(k, buf, buf)
	if err != nil {
		return nil, err
	}

	rd := bytes.NewReader(hdr)

	var entries []Blob

	pos := uint(0)
	for {
		e := headerEntry{}
		err = binary.Read(rd, binary.LittleEndian, &e)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		entry := Blob{
			Length: uint(e.Length),
			ID:     e.ID,
			Offset: pos,
		}

		switch e.Type {
		case 0:
			entry.Type = Data
		case 1:
			entry.Type = Tree
		default:
			return nil, fmt.Errorf("invalid type %d", e.Type)
		}

		entries = append(entries, entry)

		pos += uint(e.Length)
	}

	up := &Unpacker{
		rd:      rd,
		k:       k,
		Entries: entries,
	}

	return up, nil
}
