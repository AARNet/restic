// Copyright 2016, Google
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package b2

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
)

const (
	bucketName    = "b2-tests"
	smallFileName = "TeenyTiny"
	largeFileName = "BigBytes"
)

type testError struct {
	retry    bool
	backoff  time.Duration
	reauth   bool
	reupload bool
}

func (t testError) Error() string {
	return fmt.Sprintf("retry %v; backoff %v; reauth %v; reupload %v", t.retry, t.backoff, t.reauth, t.reupload)
}

type errCont struct {
	errMap map[string]map[int]error
	opMap  map[string]int
}

func (e *errCont) getError(name string) error {
	if e.errMap == nil {
		return nil
	}
	if e.opMap == nil {
		e.opMap = make(map[string]int)
	}
	i := e.opMap[name]
	e.opMap[name]++
	return e.errMap[name][i]
}

type testRoot struct {
	errs      *errCont
	auths     int
	bucketMap map[string]map[string]string
}

func (t *testRoot) authorizeAccount(context.Context, string, string, ...ClientOption) error {
	t.auths++
	return nil
}

func (t *testRoot) backoff(err error) time.Duration {
	e, ok := err.(testError)
	if !ok {
		return 0
	}
	return e.backoff
}

func (t *testRoot) reauth(err error) bool {
	e, ok := err.(testError)
	if !ok {
		return false
	}
	return e.reauth
}

func (t *testRoot) reupload(err error) bool {
	e, ok := err.(testError)
	if !ok {
		return false
	}
	return e.reupload
}

func (t *testRoot) transient(err error) bool {
	e, ok := err.(testError)
	if !ok {
		return false
	}
	return e.retry || e.reupload || e.backoff > 0
}

func (t *testRoot) createBucket(_ context.Context, name, _ string, _ map[string]string, _ []LifecycleRule) (b2BucketInterface, error) {
	if err := t.errs.getError("createBucket"); err != nil {
		return nil, err
	}
	if _, ok := t.bucketMap[name]; ok {
		return nil, fmt.Errorf("%s: bucket exists", name)
	}
	m := make(map[string]string)
	t.bucketMap[name] = m
	return &testBucket{
		n:     name,
		errs:  t.errs,
		files: m,
	}, nil
}

func (t *testRoot) listBuckets(context.Context) ([]b2BucketInterface, error) {
	var b []b2BucketInterface
	for k, v := range t.bucketMap {
		b = append(b, &testBucket{
			n:     k,
			errs:  t.errs,
			files: v,
		})
	}
	return b, nil
}

type testBucket struct {
	n     string
	errs  *errCont
	files map[string]string
}

func (t *testBucket) name() string                                     { return t.n }
func (t *testBucket) btype() string                                    { return "allPrivate" }
func (t *testBucket) attrs() *BucketAttrs                              { return nil }
func (t *testBucket) deleteBucket(context.Context) error               { return nil }
func (t *testBucket) updateBucket(context.Context, *BucketAttrs) error { return nil }

func (t *testBucket) getUploadURL(context.Context) (b2URLInterface, error) {
	if err := t.errs.getError("getUploadURL"); err != nil {
		return nil, err
	}
	return &testURL{
		files: t.files,
	}, nil
}

func (t *testBucket) startLargeFile(_ context.Context, name, _ string, _ map[string]string) (b2LargeFileInterface, error) {
	return &testLargeFile{
		name:  name,
		parts: make(map[int][]byte),
		files: t.files,
		errs:  t.errs,
	}, nil
}

func (t *testBucket) listFileNames(ctx context.Context, count int, cont, pfx, del string) ([]b2FileInterface, string, error) {
	var f []string
	for name := range t.files {
		f = append(f, name)
	}
	sort.Strings(f)
	idx := sort.SearchStrings(f, cont)
	var b []b2FileInterface
	var next string
	for i := idx; i < len(f) && i-idx < count; i++ {
		b = append(b, &testFile{
			n:     f[i],
			s:     int64(len(t.files[f[i]])),
			files: t.files,
		})
		if i+1 < len(f) {
			next = f[i+1]
		}
		if i+1 == len(f) {
			next = ""
		}
	}
	return b, next, nil
}

func (t *testBucket) listFileVersions(ctx context.Context, count int, a, b, c, d string) ([]b2FileInterface, string, string, error) {
	x, y, z := t.listFileNames(ctx, count, a, c, d)
	return x, y, "", z
}

func (t *testBucket) downloadFileByName(_ context.Context, name string, offset, size int64) (b2FileReaderInterface, error) {
	return &testFileReader{
		b: ioutil.NopCloser(bytes.NewBufferString(t.files[name][offset : offset+size])),
	}, nil
}

func (t *testBucket) hideFile(context.Context, string) (b2FileInterface, error) { return nil, nil }
func (t *testBucket) getDownloadAuthorization(context.Context, string, time.Duration) (string, error) {
	return "", nil
}
func (t *testBucket) baseURL() string { return "" }

type testURL struct {
	files map[string]string
}

func (t *testURL) reload(context.Context) error { return nil }

func (t *testURL) uploadFile(_ context.Context, r io.Reader, _ int, name, _, _ string, _ map[string]string) (b2FileInterface, error) {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}
	t.files[name] = buf.String()
	return &testFile{
		n:     name,
		s:     int64(len(t.files[name])),
		files: t.files,
	}, nil
}

type testLargeFile struct {
	name  string
	mux   sync.Mutex
	parts map[int][]byte
	files map[string]string
	errs  *errCont
}

func (t *testLargeFile) finishLargeFile(context.Context) (b2FileInterface, error) {
	var total []byte
	for i := 1; i <= len(t.parts); i++ {
		total = append(total, t.parts[i]...)
	}
	t.files[t.name] = string(total)
	return &testFile{
		n:     t.name,
		s:     int64(len(total)),
		files: t.files,
	}, nil
}

func (t *testLargeFile) getUploadPartURL(context.Context) (b2FileChunkInterface, error) {
	return &testFileChunk{
		parts: t.parts,
		mux:   &t.mux,
		errs:  t.errs,
	}, nil
}

type testFileChunk struct {
	mux   *sync.Mutex
	parts map[int][]byte
	errs  *errCont
}

func (t *testFileChunk) reload(context.Context) error { return nil }

func (t *testFileChunk) uploadPart(_ context.Context, r io.Reader, _ string, _, index int) (int, error) {
	if err := t.errs.getError("uploadPart"); err != nil {
		return 0, err
	}
	buf := &bytes.Buffer{}
	i, err := io.Copy(buf, r)
	if err != nil {
		return int(i), err
	}
	t.mux.Lock()
	t.parts[index] = buf.Bytes()
	t.mux.Unlock()
	return int(i), nil
}

type testFile struct {
	n     string
	s     int64
	t     time.Time
	a     string
	files map[string]string
}

func (t *testFile) name() string         { return t.n }
func (t *testFile) size() int64          { return t.s }
func (t *testFile) timestamp() time.Time { return t.t }
func (t *testFile) status() string       { return t.a }

func (t *testFile) compileParts(int64, map[int]string) b2LargeFileInterface {
	panic("not implemented")
}

func (t *testFile) getFileInfo(context.Context) (b2FileInfoInterface, error) {
	return nil, nil
}

func (t *testFile) listParts(context.Context, int, int) ([]b2FilePartInterface, int, error) {
	return nil, 0, nil
}

func (t *testFile) deleteFileVersion(context.Context) error {
	delete(t.files, t.n)
	return nil
}

type testFileReader struct {
	b io.ReadCloser
	s int64
}

func (t *testFileReader) Read(p []byte) (int, error)                      { return t.b.Read(p) }
func (t *testFileReader) Close() error                                    { return nil }
func (t *testFileReader) stats() (int, string, string, map[string]string) { return 0, "", "", nil }

type zReader struct{}

var pattern = []byte{0x02, 0x80, 0xff, 0x1a, 0xcc, 0x63, 0x22}

func (zReader) Read(p []byte) (int, error) {
	for i := 0; i+len(pattern) < len(p); i += len(pattern) {
		copy(p[i:], pattern)
	}
	return len(p), nil
}

func TestReauth(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	root := &testRoot{
		bucketMap: make(map[string]map[string]string),
		errs: &errCont{
			errMap: map[string]map[int]error{
				"createBucket": {0: testError{reauth: true}},
			},
		},
	}
	client := &Client{
		backend: &beRoot{
			b2i: root,
		},
	}
	auths := root.auths
	if _, err := client.NewBucket(ctx, "fun", &BucketAttrs{Type: Private}); err != nil {
		t.Errorf("bucket should not err, got %v", err)
	}
	if root.auths != auths+1 {
		t.Errorf("client should have re-authenticated; did not")
	}
}

func TestBackoff(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var calls []time.Duration
	ch := make(chan time.Time)
	close(ch)
	after = func(d time.Duration) <-chan time.Time {
		calls = append(calls, d)
		return ch
	}

	table := []struct {
		root *testRoot
		want int
	}{
		{
			root: &testRoot{
				bucketMap: make(map[string]map[string]string),
				errs: &errCont{
					errMap: map[string]map[int]error{
						"createBucket": {
							0: testError{backoff: time.Second},
							1: testError{backoff: 2 * time.Second},
						},
					},
				},
			},
			want: 2,
		},
		{
			root: &testRoot{
				bucketMap: make(map[string]map[string]string),
				errs: &errCont{
					errMap: map[string]map[int]error{
						"getUploadURL": {
							0: testError{retry: true},
						},
					},
				},
			},
			want: 1,
		},
	}

	var total int
	for _, ent := range table {
		client := &Client{
			backend: &beRoot{
				b2i: ent.root,
			},
		}
		b, err := client.NewBucket(ctx, "fun", &BucketAttrs{Type: Private})
		if err != nil {
			t.Fatal(err)
		}
		o := b.Object("foo")
		w := o.NewWriter(ctx)
		if _, err := io.Copy(w, bytes.NewBufferString("foo")); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		total += ent.want
	}
	if len(calls) != total {
		t.Errorf("got %d calls, wanted %d", len(calls), total)
	}
}

func TestBackoffWithoutRetryAfter(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var calls []time.Duration
	ch := make(chan time.Time)
	close(ch)
	after = func(d time.Duration) <-chan time.Time {
		calls = append(calls, d)
		return ch
	}

	root := &testRoot{
		bucketMap: make(map[string]map[string]string),
		errs: &errCont{
			errMap: map[string]map[int]error{
				"createBucket": {
					0: testError{retry: true},
					1: testError{retry: true},
				},
			},
		},
	}
	client := &Client{
		backend: &beRoot{
			b2i: root,
		},
	}
	if _, err := client.NewBucket(ctx, "fun", &BucketAttrs{Type: Private}); err != nil {
		t.Errorf("bucket should not err, got %v", err)
	}
	if len(calls) != 2 {
		t.Errorf("wrong number of backoff calls; got %d, want 2", len(calls))
	}
}

type badTransport struct{}

func (badTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "700 What",
		StatusCode: 700,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		Request:    r,
	}, nil
}

func TestCustomTransport(t *testing.T) {
	ctx := context.Background()
	// Sorta fragile but...
	_, err := NewClient(ctx, "abcd", "efgh", Transport(badTransport{}))
	if err == nil {
		t.Error("NewClient returned successfully, expected an error")
	}
	if !strings.Contains(err.Error(), "700") {
		t.Errorf("Expected nonsense error code 700, got %v", err)
	}
}

func TestReaderDoubleClose(t *testing.T) {
	ctx := context.Background()

	client := &Client{
		backend: &beRoot{
			b2i: &testRoot{
				bucketMap: make(map[string]map[string]string),
				errs:      &errCont{},
			},
		},
	}
	bucket, err := client.NewBucket(ctx, "bucket", &BucketAttrs{Type: Private})
	if err != nil {
		t.Fatal(err)
	}
	o, _, err := writeFile(ctx, bucket, "file", 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	r := o.NewReader(ctx)
	// Read to EOF, and then read some more.
	if _, err := io.Copy(ioutil.Discard, r); err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(ioutil.Discard, r); err != nil {
		t.Fatal(err)
	}
}

func TestReadWrite(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := &Client{
		backend: &beRoot{
			b2i: &testRoot{
				bucketMap: make(map[string]map[string]string),
				errs:      &errCont{},
			},
		},
	}

	bucket, err := client.NewBucket(ctx, bucketName, &BucketAttrs{Type: Private})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := bucket.Delete(ctx); err != nil {
			t.Error(err)
		}
	}()

	sobj, wsha, err := writeFile(ctx, bucket, smallFileName, 1e6+42, 1e8)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := sobj.Delete(ctx); err != nil {
			t.Error(err)
		}
	}()

	if err := readFile(ctx, sobj, wsha, 1e5, 10); err != nil {
		t.Error(err)
	}

	lobj, wshaL, err := writeFile(ctx, bucket, largeFileName, 1e6-1e5, 1e4)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := lobj.Delete(ctx); err != nil {
			t.Error(err)
		}
	}()

	if err := readFile(ctx, lobj, wshaL, 1e7, 10); err != nil {
		t.Error(err)
	}
}

func TestWriterReturnsError(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := &Client{
		backend: &beRoot{
			b2i: &testRoot{
				bucketMap: make(map[string]map[string]string),
				errs: &errCont{
					errMap: map[string]map[int]error{
						"uploadPart": {
							0: testError{},
							1: testError{},
							2: testError{},
							3: testError{},
							4: testError{},
							5: testError{},
							6: testError{},
						},
					},
				},
			},
		},
	}

	bucket, err := client.NewBucket(ctx, bucketName, &BucketAttrs{Type: Private})
	if err != nil {
		t.Fatal(err)
	}
	w := bucket.Object("test").NewWriter(ctx)
	r := io.LimitReader(zReader{}, 1e7)
	w.ChunkSize = 1e4
	w.ConcurrentUploads = 4
	if _, err := io.Copy(w, r); err == nil {
		t.Fatalf("io.Copy: should have returned an error")
	}
}

func TestFileBuffer(t *testing.T) {
	r := io.LimitReader(zReader{}, 1e8)
	w, err := newFileBuffer("")
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		t.Fatal(err)
	}
	bReader, err := w.Reader()
	if err != nil {
		t.Fatal(err)
	}
	hsh := sha1.New()
	if _, err := io.Copy(hsh, bReader); err != nil {
		t.Fatal(err)
	}
	hshText := fmt.Sprintf("%x", hsh.Sum(nil))
	if hshText != w.Hash() {
		t.Errorf("hashes are not equal: bufferWriter is %q, read buffer is %q", w.Hash(), hshText)
	}
}

func writeFile(ctx context.Context, bucket *Bucket, name string, size int64, csize int) (*Object, string, error) {
	r := io.LimitReader(zReader{}, size)
	o := bucket.Object(name)
	f := o.NewWriter(ctx)
	h := sha1.New()
	w := io.MultiWriter(f, h)
	f.ConcurrentUploads = 5
	f.ChunkSize = csize
	if _, err := io.Copy(w, r); err != nil {
		return nil, "", err
	}
	if err := f.Close(); err != nil {
		return nil, "", err
	}
	return o, fmt.Sprintf("%x", h.Sum(nil)), nil
}

func readFile(ctx context.Context, obj *Object, sha string, chunk, concur int) error {
	r := obj.NewReader(ctx)
	r.ChunkSize = chunk
	r.ConcurrentDownloads = concur
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return err
	}
	if err := r.Close(); err != nil {
		return err
	}
	rsha := fmt.Sprintf("%x", h.Sum(nil))
	if sha != rsha {
		return fmt.Errorf("bad hash: got %s, want %s", rsha, sha)
	}
	return nil
}
