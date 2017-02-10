package archiver

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"restic"
	"restic/checker"
	"restic/repository"
	"testing"
)

func loadBlob(t *testing.T, repo restic.Repository, id restic.ID, buf []byte) int {
	n, err := repo.LoadBlob(restic.DataBlob, id, buf)
	if err != nil {
		t.Fatalf("LoadBlob(%v) returned error %v", id, err)
	}

	return n
}

func checkSavedFile(t *testing.T, repo restic.Repository, treeID restic.ID, name string, rd io.Reader) {
	tree, err := repo.LoadTree(treeID)
	if err != nil {
		t.Fatalf("LoadTree() returned error %v", err)
	}

	if len(tree.Nodes) != 1 {
		t.Fatalf("wrong number of nodes for tree, want %v, got %v", 1, len(tree.Nodes))
	}

	node := tree.Nodes[0]
	if node.Name != "fakefile" {
		t.Fatalf("wrong filename, want %v, got %v", "fakefile", node.Name)
	}

	if len(node.Content) == 0 {
		t.Fatalf("node.Content has length 0")
	}

	// check blobs
	for i, id := range node.Content {
		size, err := repo.LookupBlobSize(id, restic.DataBlob)
		if err != nil {
			t.Fatal(err)
		}

		buf := restic.NewBlobBuffer(int(size))
		n := loadBlob(t, repo, id, buf)
		if n != len(buf) {
			t.Errorf("wrong number of bytes read, want %d, got %d", len(buf), n)
		}

		buf2 := make([]byte, int(size))
		_, err = io.ReadFull(rd, buf2)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(buf, buf2) {
			t.Fatalf("blob %d (%v) is wrong", i, id.Str())
		}
	}
}

// fakeFile returns a reader which yields deterministic pseudo-random data.
func fakeFile(t testing.TB, seed, size int64) io.Reader {
	return io.LimitReader(restic.NewRandReader(rand.New(rand.NewSource(seed))), size)
}

func TestArchiveReader(t *testing.T) {
	repo, cleanup := repository.TestRepository(t)
	defer cleanup()

	seed := rand.Int63()
	size := int64(rand.Intn(50*1024*1024) + 50*1024*1024)
	t.Logf("seed is 0x%016x, size is %v", seed, size)

	f := fakeFile(t, seed, size)

	sn, id, err := ArchiveReader(repo, nil, f, "fakefile", []string{"test"}, "localhost")
	if err != nil {
		t.Fatalf("ArchiveReader() returned error %v", err)
	}

	if id.IsNull() {
		t.Fatalf("ArchiveReader() returned null ID")
	}

	t.Logf("snapshot saved as %v, tree is %v", id.Str(), sn.Tree.Str())

	checkSavedFile(t, repo, *sn.Tree, "fakefile", fakeFile(t, seed, size))

	checker.TestCheckRepo(t, repo)
}

func TestArchiveReaderNull(t *testing.T) {
	repo, cleanup := repository.TestRepository(t)
	defer cleanup()

	sn, id, err := ArchiveReader(repo, nil, bytes.NewReader(nil), "fakefile", nil, "localhost")
	if err != nil {
		t.Fatalf("ArchiveReader() returned error %v", err)
	}

	if id.IsNull() {
		t.Fatalf("ArchiveReader() returned null ID")
	}

	t.Logf("snapshot saved as %v, tree is %v", id.Str(), sn.Tree.Str())

	checker.TestCheckRepo(t, repo)
}

type errReader string

func (e errReader) Read([]byte) (int, error) {
	return 0, errors.New(string(e))
}

func countSnapshots(t testing.TB, repo restic.Repository) int {
	done := make(chan struct{})
	defer close(done)

	snapshots := 0
	for range repo.List(restic.SnapshotFile, done) {
		snapshots++
	}
	return snapshots
}

func TestArchiveReaderError(t *testing.T) {
	repo, cleanup := repository.TestRepository(t)
	defer cleanup()

	sn, id, err := ArchiveReader(repo, nil, errReader("error returned by reading stdin"), "fakefile", nil, "localhost")
	if err == nil {
		t.Errorf("expected error not returned")
	}

	if sn != nil {
		t.Errorf("Snapshot should be nil, but isn't")
	}

	if !id.IsNull() {
		t.Errorf("id should be null, but %v returned", id.Str())
	}

	n := countSnapshots(t, repo)
	if n > 0 {
		t.Errorf("expected zero snapshots, but got %d", n)
	}

	checker.TestCheckRepo(t, repo)
}

func BenchmarkArchiveReader(t *testing.B) {
	repo, cleanup := repository.TestRepository(t)
	defer cleanup()

	const size = 50 * 1024 * 1024

	buf := make([]byte, size)
	_, err := io.ReadFull(fakeFile(t, 23, size), buf)
	if err != nil {
		t.Fatal(err)
	}

	t.SetBytes(size)
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		_, _, err := ArchiveReader(repo, nil, bytes.NewReader(buf), "fakefile", []string{"test"}, "localhost")
		if err != nil {
			t.Fatal(err)
		}
	}
}
