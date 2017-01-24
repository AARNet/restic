// +build !openbsd
// +build !windows

package fuse

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/net/context"

	"restic/repository"

	"bazil.org/fuse"

	"restic"
	. "restic/test"
)

func testRead(t testing.TB, f *file, offset, length int, data []byte) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &fuse.ReadRequest{
		Offset: int64(offset),
		Size:   length,
	}
	resp := &fuse.ReadResponse{
		Data: data,
	}
	OK(t, f.Read(ctx, req, resp))
}

func firstSnapshotID(t testing.TB, repo restic.Repository) (first restic.ID) {
	done := make(chan struct{})
	defer close(done)
	for id := range repo.List(restic.SnapshotFile, done) {
		if first.IsNull() {
			first = id
		}
	}
	return first
}

func loadFirstSnapshot(t testing.TB, repo restic.Repository) *restic.Snapshot {
	id := firstSnapshotID(t, repo)
	sn, err := restic.LoadSnapshot(repo, id)
	OK(t, err)
	return sn
}

func loadTree(t testing.TB, repo restic.Repository, id restic.ID) *restic.Tree {
	tree, err := repo.LoadTree(id)
	OK(t, err)
	return tree
}

func TestFuseFile(t *testing.T) {
	repo, cleanup := repository.TestRepository(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timestamp, err := time.Parse(time.RFC3339, "2017-01-24T10:42:56+01:00")
	OK(t, err)
	restic.TestCreateSnapshot(t, repo, timestamp, 2, 0.1)

	sn := loadFirstSnapshot(t, repo)
	tree := loadTree(t, repo, *sn.Tree)

	var content restic.IDs
	for _, node := range tree.Nodes {
		content = append(content, node.Content...)
	}
	t.Logf("tree loaded, content: %v", content)

	var (
		filesize uint64
		memfile  []byte
	)
	for _, id := range content {
		size, err := repo.LookupBlobSize(id, restic.DataBlob)
		OK(t, err)
		filesize += uint64(size)

		buf := restic.NewBlobBuffer(int(size))
		n, err := repo.LoadBlob(restic.DataBlob, id, buf)
		OK(t, err)

		if uint(n) != size {
			t.Fatalf("not enough bytes read for id %v: want %v, got %v", id.Str(), size, n)
		}

		if uint(len(buf)) != size {
			t.Fatalf("buffer has wrong length for id %v: want %v, got %v", id.Str(), size, len(buf))
		}

		memfile = append(memfile, buf...)
	}

	t.Logf("filesize is %v, memfile has size %v", filesize, len(memfile))

	node := &restic.Node{
		Name:    "foo",
		Inode:   23,
		Mode:    0742,
		Size:    filesize,
		Content: content,
	}
	f, err := newFile(repo, node, false)
	OK(t, err)

	attr := fuse.Attr{}
	OK(t, f.Attr(ctx, &attr))

	Equals(t, node.Inode, attr.Inode)
	Equals(t, node.Mode, attr.Mode)
	Equals(t, node.Size, attr.Size)
	Equals(t, (node.Size/uint64(attr.BlockSize))+1, attr.Blocks)

	for i := 0; i < 200; i++ {
		offset := rand.Intn(int(filesize))
		length := rand.Intn(int(filesize)-offset) + 100

		b := memfile[offset : offset+length]

		buf := make([]byte, length)

		testRead(t, f, offset, length, buf)
		if !bytes.Equal(b, buf) {
			t.Errorf("test %d failed, wrong data returned (offset %v, length %v)", i, offset, length)
		}
	}

	OK(t, f.Release(ctx, nil))
}
