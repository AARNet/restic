// +build !openbsd
// +build !windows

package fuse

import (
	"sync"

	"github.com/restic/restic"
	"github.com/restic/restic/backend"
	"github.com/restic/restic/pack"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

// Statically ensure that *file implements the given interface
var _ = fs.HandleReader(&file{})
var _ = fs.HandleReleaser(&file{})

// BlobLoader is an abstracted repository with a reduced set of methods used
// for fuse operations.
type BlobLoader interface {
	LookupBlobSize(backend.ID) (uint, error)
	LoadBlob(pack.BlobType, backend.ID, []byte) ([]byte, error)
}

type file struct {
	repo        BlobLoader
	node        *restic.Node
	ownerIsRoot bool

	sizes []uint
	blobs [][]byte
}

const defaultBlobSize = 128 * 1024

var blobPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, defaultBlobSize)
	},
}

func newFile(repo BlobLoader, node *restic.Node, ownerIsRoot bool) (*file, error) {
	sizes := make([]uint, len(node.Content))
	for i, id := range node.Content {
		size, err := repo.LookupBlobSize(id)
		if err != nil {
			return nil, err
		}

		sizes[i] = size
	}

	return &file{
		repo:        repo,
		node:        node,
		sizes:       sizes,
		blobs:       make([][]byte, len(node.Content)),
		ownerIsRoot: ownerIsRoot,
	}, nil
}

func (f *file) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = f.node.Inode
	a.Mode = f.node.Mode
	a.Size = f.node.Size

	if !f.ownerIsRoot {
		a.Uid = f.node.UID
		a.Gid = f.node.GID
	}
	a.Atime = f.node.AccessTime
	a.Ctime = f.node.ChangeTime
	a.Mtime = f.node.ModTime
	return nil
}

func (f *file) getBlobAt(i int) (blob []byte, err error) {
	if f.blobs[i] != nil {
		return f.blobs[i], nil
	}

	buf := blobPool.Get().([]byte)
	buf = buf[:cap(buf)]

	if uint(len(buf)) < f.sizes[i] {
		if len(buf) > defaultBlobSize {
			blobPool.Put(buf)
		}
		buf = make([]byte, f.sizes[i])
	}

	blob, err = f.repo.LoadBlob(pack.Data, f.node.Content[i], buf)
	if err != nil {
		return nil, err
	}
	f.blobs[i] = blob

	return blob, nil
}

func (f *file) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	offset := req.Offset

	// Skip blobs before the offset
	startContent := 0
	for offset > int64(f.sizes[startContent]) {
		offset -= int64(f.sizes[startContent])
		startContent++
	}

	dst := resp.Data[0:req.Size]
	readBytes := 0
	remainingBytes := req.Size
	for i := startContent; remainingBytes > 0 && i < len(f.sizes); i++ {
		blob, err := f.getBlobAt(i)
		if err != nil {
			return err
		}

		if offset > 0 {
			blob = blob[offset:len(blob)]
			offset = 0
		}

		copied := copy(dst, blob)
		remainingBytes -= copied
		readBytes += copied

		dst = dst[copied:]
	}
	resp.Data = resp.Data[:readBytes]

	return nil
}

func (f *file) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	for i := range f.blobs {
		if f.blobs[i] != nil {
			blobPool.Put(f.blobs[i])
			f.blobs[i] = nil
		}
	}
	return nil
}
