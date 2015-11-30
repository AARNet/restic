package backend_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/restic/restic/backend"
	. "github.com/restic/restic/test"
)

func testBackendConfig(b backend.Backend, t *testing.T) {
	// create config and read it back
	_, err := b.Get(backend.Config, "")
	Assert(t, err != nil, "did not get expected error for non-existing config")

	blob, err := b.Create()
	OK(t, err)

	_, err = blob.Write([]byte("Config"))
	OK(t, err)
	OK(t, blob.Finalize(backend.Config, ""))

	// try accessing the config with different names, should all return the
	// same config
	for _, name := range []string{"", "foo", "bar", "0000000000000000000000000000000000000000000000000000000000000000"} {
		rd, err := b.Get(backend.Config, name)
		Assert(t, err == nil, "unable to read config")

		buf, err := ioutil.ReadAll(rd)
		OK(t, err)
		OK(t, rd.Close())
		Assert(t, string(buf) == "Config", "wrong data returned for config")
	}
}

func testBackend(b backend.Backend, t *testing.T) {
	testBackendConfig(b, t)

	for _, tpe := range []backend.Type{
		backend.Data, backend.Key, backend.Lock,
		backend.Snapshot, backend.Index,
	} {
		// detect non-existing files
		for _, test := range TestStrings {
			id, err := backend.ParseID(test.id)
			OK(t, err)

			// test if blob is already in repository
			ret, err := b.Test(tpe, id.String())
			OK(t, err)
			Assert(t, !ret, "blob was found to exist before creating")

			// try to open not existing blob
			_, err = b.Get(tpe, id.String())
			Assert(t, err != nil, "blob data could be extracted before creation")

			// try to read not existing blob
			_, err = b.GetReader(tpe, id.String(), 0, 1)
			Assert(t, err != nil, "blob reader could be obtained before creation")

			// try to get string out, should fail
			ret, err = b.Test(tpe, id.String())
			OK(t, err)
			Assert(t, !ret, "id %q was found (but should not have)", test.id)
		}

		// add files
		for _, test := range TestStrings {
			// store string in backend
			blob, err := b.Create()
			OK(t, err)

			_, err = blob.Write([]byte(test.data))
			OK(t, err)
			OK(t, blob.Finalize(tpe, test.id))

			// try to get it out again
			rd, err := b.Get(tpe, test.id)
			OK(t, err)
			Assert(t, rd != nil, "Get() returned nil")

			// try to read it out again
			reader, err := b.GetReader(tpe, test.id, 0, uint(len(test.data)))
			OK(t, err)
			Assert(t, reader != nil, "GetReader() returned nil")
			bytes := make([]byte, len(test.data))
			reader.Read(bytes)
			Assert(t, test.data == string(bytes), "Read() returned different content")

			// try to read it out with an offset and a length
			readerOffLen, err := b.GetReader(tpe, test.id, 1, uint(len(test.data)-2))
			OK(t, err)
			Assert(t, readerOffLen != nil, "GetReader() returned nil")
			bytesOffLen := make([]byte, len(test.data)-2)
			readerOffLen.Read(bytesOffLen)
			Assert(t, test.data[1:len(test.data)-1] == string(bytesOffLen), "Read() with offset and length returned different content")

			buf, err := ioutil.ReadAll(rd)
			OK(t, err)
			Equals(t, test.data, string(buf))

			// compare content
			Equals(t, test.data, string(buf))
		}

		// test adding the first file again
		test := TestStrings[0]

		// create blob
		blob, err := b.Create()
		OK(t, err)

		_, err = blob.Write([]byte(test.data))
		OK(t, err)
		err = blob.Finalize(tpe, test.id)
		Assert(t, err != nil, "expected error, got %v", err)

		// remove and recreate
		err = b.Remove(tpe, test.id)
		OK(t, err)

		// test that the blob is gone
		ok, err := b.Test(tpe, test.id)
		OK(t, err)
		Assert(t, ok == false, "removed blob still present")

		// create blob
		blob, err = b.Create()
		OK(t, err)

		_, err = io.Copy(blob, bytes.NewReader([]byte(test.data)))
		OK(t, err)
		OK(t, blob.Finalize(tpe, test.id))

		// list items
		IDs := backend.IDs{}

		for _, test := range TestStrings {
			id, err := backend.ParseID(test.id)
			OK(t, err)
			IDs = append(IDs, id)
		}

		sort.Sort(IDs)

		i := 0
		for s := range b.List(tpe, nil) {
			Equals(t, IDs[i].String(), s)
			i++
		}

		// remove content if requested
		if TestCleanup {
			for _, test := range TestStrings {
				id, err := backend.ParseID(test.id)
				OK(t, err)

				found, err := b.Test(tpe, id.String())
				OK(t, err)
				Assert(t, found, fmt.Sprintf("id %q was not found before removal", id))

				OK(t, b.Remove(tpe, id.String()))

				found, err = b.Test(tpe, id.String())
				OK(t, err)
				Assert(t, !found, fmt.Sprintf("id %q not found after removal", id))
			}
		}

	}
}
