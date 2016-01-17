/*
 * Minio Go Library for Amazon S3 Compatible Cloud Storage (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package minio

import (
	"net/url"
	"testing"
)

func TestSignature(t *testing.T) {
	clnt := Client{}
	if !clnt.signature.isV4() {
		t.Fatal("Error")
	}
	clnt.signature = SignatureV2
	if !clnt.signature.isV2() {
		t.Fatal("Error")
	}
	if clnt.signature.isV4() {
		t.Fatal("Error")
	}
	clnt.signature = SignatureV4
	if !clnt.signature.isV4() {
		t.Fatal("Error")
	}
}

func TestACLTypes(t *testing.T) {
	want := map[string]bool{
		"private":            true,
		"public-read":        true,
		"public-read-write":  true,
		"authenticated-read": true,
		"invalid":            false,
	}
	for acl, ok := range want {
		if BucketACL(acl).isValidBucketACL() != ok {
			t.Fatal("Error")
		}
	}
}

func TestPartSize(t *testing.T) {
	var maxPartSize int64 = 1024 * 1024 * 1024 * 5
	partSize := optimalPartSize(5000000000000000000)
	if partSize > minimumPartSize {
		if partSize > maxPartSize {
			t.Fatal("invalid result, cannot be bigger than maxPartSize 5GiB")
		}
	}
	partSize = optimalPartSize(50000000000)
	if partSize > minimumPartSize {
		t.Fatal("invalid result, cannot be bigger than minimumPartSize 5MiB")
	}
}

func TestURLEncoding(t *testing.T) {
	type urlStrings struct {
		name        string
		encodedName string
	}

	want := []urlStrings{
		{
			name:        "bigfile-1._%",
			encodedName: "bigfile-1._%25",
		},
		{
			name:        "本語",
			encodedName: "%E6%9C%AC%E8%AA%9E",
		},
		{
			name:        "本語.1",
			encodedName: "%E6%9C%AC%E8%AA%9E.1",
		},
		{
			name:        ">123>3123123",
			encodedName: "%3E123%3E3123123",
		},
		{
			name:        "test 1 2.txt",
			encodedName: "test%201%202.txt",
		},
		{
			name:        "test++ 1.txt",
			encodedName: "test%2B%2B%201.txt",
		},
	}

	for _, u := range want {
		if u.encodedName != urlEncodePath(u.name) {
			t.Fatal("Error")
		}
	}
}

func TestGetEndpointURL(t *testing.T) {
	if _, err := getEndpointURL("s3.amazonaws.com", false); err != nil {
		t.Fatal("Error:", err)
	}
	if _, err := getEndpointURL("192.168.1.1", false); err != nil {
		t.Fatal("Error:", err)
	}
	if _, err := getEndpointURL("13333.123123.-", false); err == nil {
		t.Fatal("Error")
	}
	if _, err := getEndpointURL("s3.aamzza.-", false); err == nil {
		t.Fatal("Error")
	}
	if _, err := getEndpointURL("s3.amazonaws.com:443", false); err == nil {
		t.Fatal("Error")
	}
}

func TestValidIP(t *testing.T) {
	type validIP struct {
		ip    string
		valid bool
	}

	want := []validIP{
		{
			ip:    "192.168.1.1",
			valid: true,
		},
		{
			ip:    "192.1.8",
			valid: false,
		},
		{
			ip:    "..192.",
			valid: false,
		},
		{
			ip:    "192.168.1.1.1",
			valid: false,
		},
	}
	for _, w := range want {
		valid := isValidIP(w.ip)
		if valid != w.valid {
			t.Fatal("Error")
		}
	}
}

func TestValidEndpointDomain(t *testing.T) {
	type validEndpoint struct {
		endpointDomain string
		valid          bool
	}

	want := []validEndpoint{
		{
			endpointDomain: "s3.amazonaws.com",
			valid:          true,
		},
		{
			endpointDomain: "s3.amazonaws.com_",
			valid:          false,
		},
		{
			endpointDomain: "%$$$",
			valid:          false,
		},
		{
			endpointDomain: "s3.amz.test.com",
			valid:          true,
		},
		{
			endpointDomain: "s3.%%",
			valid:          false,
		},
		{
			endpointDomain: "localhost",
			valid:          true,
		},
		{
			endpointDomain: "-localhost",
			valid:          false,
		},
		{
			endpointDomain: "",
			valid:          false,
		},
		{
			endpointDomain: "\n \t",
			valid:          false,
		},
		{
			endpointDomain: "    ",
			valid:          false,
		},
	}
	for _, w := range want {
		valid := isValidDomain(w.endpointDomain)
		if valid != w.valid {
			t.Fatal("Error:", w.endpointDomain)
		}
	}
}

func TestValidEndpointURL(t *testing.T) {
	type validURL struct {
		url   string
		valid bool
	}
	want := []validURL{
		{
			url:   "https://s3.amazonaws.com",
			valid: true,
		},
		{
			url:   "https://s3.amazonaws.com/bucket/object",
			valid: false,
		},
		{
			url:   "192.168.1.1",
			valid: false,
		},
	}
	for _, w := range want {
		u, err := url.Parse(w.url)
		if err != nil {
			t.Fatal("Error:", err)
		}
		valid := false
		if err := isValidEndpointURL(u); err == nil {
			valid = true
		}
		if valid != w.valid {
			t.Fatal("Error")
		}
	}
}
