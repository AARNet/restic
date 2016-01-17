// +build ignore

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

package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-bucketname, my-objectname and
	// my-testfile are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set insecure=true to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API copatibality (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := minio.New("s3.amazonaws.com", "YOUR-ACCESS-KEY-HERE", "YOUR-SECRET-KEY-HERE", false)
	if err != nil {
		log.Fatalln(err)
	}

	reader, stat, err := s3Client.GetObjectPartial("my-bucketname", "my-objectname")
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	localFile, err := os.OpenFile("my-testfile", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer localfile.Close()

	st, err := localFile.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	readAtOffset := st.Size()
	readAtBuffer := make([]byte, 5*1024*1024)

	// For loop.
	for {
		readAtSize, rerr := reader.ReadAt(readAtBuffer, readAtOffset)
		if rerr != nil {
			if rerr != io.EOF {
				log.Fatalln(rerr)
			}
		}
		writeSize, werr := localFile.Write(readAtBuffer[:readAtSize])
		if werr != nil {
			log.Fatalln(werr)
		}
		if readAtSize != writeSize {
			log.Fatalln(errors.New("Something really bad happened here."))
		}
		readAtOffset += int64(writeSize)
		if rerr == io.EOF {
			break
		}
	}

	// totalWritten size.
	totalWritten := readAtOffset

	// If found mismatch error out.
	if totalWritten != stat.Size {
		log.Fatalln(errors.New("Something really bad happened here."))
	}
}
