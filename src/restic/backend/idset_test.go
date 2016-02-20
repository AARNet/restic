package backend_test

import (
	"testing"

	"restic/backend"
	. "restic/test"
)

var idsetTests = []struct {
	id   backend.ID
	seen bool
}{
	{ParseID("7bb086db0d06285d831485da8031281e28336a56baa313539eaea1c73a2a1a40"), false},
	{ParseID("1285b30394f3b74693cc29a758d9624996ae643157776fce8154aabd2f01515f"), false},
	{ParseID("7bb086db0d06285d831485da8031281e28336a56baa313539eaea1c73a2a1a40"), true},
	{ParseID("7bb086db0d06285d831485da8031281e28336a56baa313539eaea1c73a2a1a40"), true},
	{ParseID("1285b30394f3b74693cc29a758d9624996ae643157776fce8154aabd2f01515f"), true},
	{ParseID("f658198b405d7e80db5ace1980d125c8da62f636b586c46bf81dfb856a49d0c8"), false},
	{ParseID("7bb086db0d06285d831485da8031281e28336a56baa313539eaea1c73a2a1a40"), true},
	{ParseID("1285b30394f3b74693cc29a758d9624996ae643157776fce8154aabd2f01515f"), true},
	{ParseID("f658198b405d7e80db5ace1980d125c8da62f636b586c46bf81dfb856a49d0c8"), true},
	{ParseID("7bb086db0d06285d831485da8031281e28336a56baa313539eaea1c73a2a1a40"), true},
}

func TestIDSet(t *testing.T) {
	set := backend.NewIDSet()
	for i, test := range idsetTests {
		seen := set.Has(test.id)
		if seen != test.seen {
			t.Errorf("IDSet test %v failed: wanted %v, got %v", i, test.seen, seen)
		}
		set.Insert(test.id)
	}
}
