package walker

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/restic/restic/internal/restic"
)

// TestTree is used to construct a list of trees for testing the walker.
type TestTree map[string]interface{}

// TestNode is used to test the walker.
type TestFile struct{}

func BuildTreeMap(tree TestTree) (m TreeMap, root restic.ID) {
	m = TreeMap{}
	id := buildTreeMap(tree, m)
	return m, id
}

func buildTreeMap(tree TestTree, m TreeMap) restic.ID {
	res := restic.NewTree()

	for name, item := range tree {
		switch elem := item.(type) {
		case TestFile:
			res.Insert(&restic.Node{
				Name: name,
				Type: "file",
			})
		case TestTree:
			id := buildTreeMap(elem, m)
			res.Insert(&restic.Node{
				Name:    name,
				Subtree: &id,
				Type:    "dir",
			})
		default:
			panic(fmt.Sprintf("invalid type %T", elem))
		}
	}

	buf, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	id := restic.Hash(buf)

	if _, ok := m[id]; !ok {
		m[id] = res
	}

	return id
}

// TreeMap returns the trees from the map on LoadTree.
type TreeMap map[restic.ID]*restic.Tree

func (t TreeMap) LoadTree(ctx context.Context, id restic.ID) (*restic.Tree, error) {
	tree, ok := t[id]
	if !ok {
		return nil, errors.New("tree not found")
	}

	return tree, nil
}

// checkFunc returns a function suitable for walking the tree to check
// something, and a function which will check the final result.
type checkFunc func(t testing.TB) (walker WalkFunc, final func(testing.TB))

// checkItemOrder ensures that the order of the 'path' arguments is the one passed in as 'want'.
func checkItemOrder(want []string) checkFunc {
	pos := 0
	return func(t testing.TB) (walker WalkFunc, final func(testing.TB)) {
		walker = func(path string, node *restic.Node, err error) (bool, error) {
			if err != nil {
				t.Errorf("error walking %v: %v", path, err)
				return false, err
			}

			if pos >= len(want) {
				t.Errorf("additional unexpected path found: %v", path)
				return false, nil
			}

			if path != want[pos] {
				t.Errorf("wrong path found, want %q, got %q", want[pos], path)
			}
			pos++
			return false, nil
		}

		final = func(t testing.TB) {
			if pos != len(want) {
				t.Errorf("not enough items returned, want %d, got %d", len(want), pos)
			}
		}

		return walker, final
	}
}

// checkSkipFor returns SkipNode if path is in skipFor, it checks that the
// paths the walk func is called for are exactly the ones in wantPaths.
func checkSkipFor(skipFor map[string]struct{}, wantPaths []string) checkFunc {
	var pos int

	return func(t testing.TB) (walker WalkFunc, final func(testing.TB)) {
		walker = func(path string, node *restic.Node, err error) (bool, error) {
			if err != nil {
				t.Errorf("error walking %v: %v", path, err)
				return false, err
			}

			if pos >= len(wantPaths) {
				t.Errorf("additional unexpected path found: %v", path)
				return false, nil
			}

			if path != wantPaths[pos] {
				t.Errorf("wrong path found, want %q, got %q", wantPaths[pos], path)
			}
			pos++

			if _, ok := skipFor[path]; ok {
				return false, SkipNode
			}

			return false, nil
		}

		final = func(t testing.TB) {
			if pos != len(wantPaths) {
				t.Errorf("wrong number of paths returned, want %d, got %d", len(wantPaths), pos)
			}
		}

		return walker, final
	}
}

// checkIgnore returns SkipNode if path is in skipFor and sets ignore according
// to ignoreFor. It checks that the paths the walk func is called for are exactly
// the ones in wantPaths.
func checkIgnore(skipFor map[string]struct{}, ignoreFor map[string]bool, wantPaths []string) checkFunc {
	var pos int

	return func(t testing.TB) (walker WalkFunc, final func(testing.TB)) {
		walker = func(path string, node *restic.Node, err error) (bool, error) {
			if err != nil {
				t.Errorf("error walking %v: %v", path, err)
				return false, err
			}

			if pos >= len(wantPaths) {
				t.Errorf("additional unexpected path found: %v", path)
				return ignoreFor[path], nil
			}

			if path != wantPaths[pos] {
				t.Errorf("wrong path found, want %q, got %q", wantPaths[pos], path)
			}
			pos++

			if _, ok := skipFor[path]; ok {
				return ignoreFor[path], SkipNode
			}

			return ignoreFor[path], nil
		}

		final = func(t testing.TB) {
			if pos != len(wantPaths) {
				t.Errorf("wrong number of paths returned, want %d, got %d", len(wantPaths), pos)
			}
		}

		return walker, final
	}
}

func TestWalker(t *testing.T) {
	var tests = []struct {
		tree   TestTree
		checks []checkFunc
	}{
		{
			tree: TestTree{
				"foo": TestFile{},
				"subdir": TestTree{
					"subfile": TestFile{},
				},
			},
			checks: []checkFunc{
				checkItemOrder([]string{
					"/",
					"/foo",
					"/subdir",
					"/subdir/subfile",
				}),
				checkSkipFor(
					map[string]struct{}{
						"/subdir": struct{}{},
					}, []string{
						"/",
						"/foo",
						"/subdir",
					},
				),
				checkIgnore(
					map[string]struct{}{}, map[string]bool{
						"/subdir": true,
					}, []string{
						"/",
						"/foo",
						"/subdir",
						"/subdir/subfile",
					},
				),
			},
		},
		{
			tree: TestTree{
				"foo": TestFile{},
				"subdir1": TestTree{
					"subfile1": TestFile{},
				},
				"subdir2": TestTree{
					"subfile2": TestFile{},
					"subsubdir2": TestTree{
						"subsubfile3": TestFile{},
					},
				},
			},
			checks: []checkFunc{
				checkItemOrder([]string{
					"/",
					"/foo",
					"/subdir1",
					"/subdir1/subfile1",
					"/subdir2",
					"/subdir2/subfile2",
					"/subdir2/subsubdir2",
					"/subdir2/subsubdir2/subsubfile3",
				}),
				checkSkipFor(
					map[string]struct{}{
						"/subdir1": struct{}{},
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/subdir2",
						"/subdir2/subfile2",
						"/subdir2/subsubdir2",
						"/subdir2/subsubdir2/subsubfile3",
					},
				),
				checkSkipFor(
					map[string]struct{}{
						"/subdir1":            struct{}{},
						"/subdir2/subsubdir2": struct{}{},
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/subdir2",
						"/subdir2/subfile2",
						"/subdir2/subsubdir2",
					},
				),
				checkSkipFor(
					map[string]struct{}{
						"/foo": struct{}{},
					}, []string{
						"/",
						"/foo",
					},
				),
			},
		},
		{
			tree: TestTree{
				"foo": TestFile{},
				"subdir1": TestTree{
					"subfile1": TestFile{},
					"subfile2": TestFile{},
					"subfile3": TestFile{},
				},
				"subdir2": TestTree{
					"subfile1": TestFile{},
					"subfile2": TestFile{},
					"subfile3": TestFile{},
				},
				"subdir3": TestTree{
					"subfile1": TestFile{},
					"subfile2": TestFile{},
					"subfile3": TestFile{},
				},
				"zzz other": TestFile{},
			},
			checks: []checkFunc{
				checkItemOrder([]string{
					"/",
					"/foo",
					"/subdir1",
					"/subdir1/subfile1",
					"/subdir1/subfile2",
					"/subdir1/subfile3",
					"/subdir2",
					"/subdir2/subfile1",
					"/subdir2/subfile2",
					"/subdir2/subfile3",
					"/subdir3",
					"/subdir3/subfile1",
					"/subdir3/subfile2",
					"/subdir3/subfile3",
					"/zzz other",
				}),
				checkIgnore(
					map[string]struct{}{
						"/subdir1": struct{}{},
					}, map[string]bool{
						"/subdir1": true,
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/zzz other",
					},
				),
				checkIgnore(
					map[string]struct{}{}, map[string]bool{
						"/subdir1": true,
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/subdir1/subfile1",
						"/subdir1/subfile2",
						"/subdir1/subfile3",
						"/zzz other",
					},
				),
				checkIgnore(
					map[string]struct{}{
						"/subdir2": struct{}{},
					}, map[string]bool{
						"/subdir2": true,
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/subdir1/subfile1",
						"/subdir1/subfile2",
						"/subdir1/subfile3",
						"/subdir2",
						"/zzz other",
					},
				),
				checkIgnore(
					map[string]struct{}{}, map[string]bool{
						"/subdir1/subfile1": true,
						"/subdir1/subfile2": true,
						"/subdir1/subfile3": true,
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/subdir1/subfile1",
						"/subdir1/subfile2",
						"/subdir1/subfile3",
						"/zzz other",
					},
				),
				checkIgnore(
					map[string]struct{}{}, map[string]bool{
						"/subdir2/subfile1": true,
						"/subdir2/subfile2": true,
						"/subdir2/subfile3": true,
					}, []string{
						"/",
						"/foo",
						"/subdir1",
						"/subdir1/subfile1",
						"/subdir1/subfile2",
						"/subdir1/subfile3",
						"/subdir2",
						"/subdir2/subfile1",
						"/subdir2/subfile2",
						"/subdir2/subfile3",
						"/zzz other",
					},
				),
			},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			repo, root := BuildTreeMap(test.tree)
			for _, check := range test.checks {
				t.Run("", func(t *testing.T) {
					ctx, cancel := context.WithCancel(context.TODO())
					defer cancel()

					fn, last := check(t)
					err := Walk(ctx, repo, root, restic.NewIDSet(), fn)
					if err != nil {
						t.Error(err)
					}
					last(t)
				})
			}
		})
	}
}
