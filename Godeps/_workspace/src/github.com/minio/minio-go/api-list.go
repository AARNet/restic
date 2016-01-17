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
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ListBuckets list all buckets owned by this authenticated user.
//
// This call requires explicit authentication, no anonymous requests are
// allowed for listing buckets.
//
//   api := client.New(....)
//   for message := range api.ListBuckets() {
//       fmt.Println(message)
//   }
//
func (c Client) ListBuckets() ([]BucketStat, error) {
	// Instantiate a new request.
	req, err := c.newRequest("GET", requestMetadata{})
	if err != nil {
		return nil, err
	}
	// Initiate the request.
	resp, err := c.do(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return nil, HTTPRespToErrorResponse(resp, "", "")
		}
	}
	listAllMyBucketsResult := listAllMyBucketsResult{}
	err = xmlDecoder(resp.Body, &listAllMyBucketsResult)
	if err != nil {
		return nil, err
	}
	return listAllMyBucketsResult.Buckets.Bucket, nil
}

// ListObjects - (List Objects) - List some objects or all recursively.
//
// ListObjects lists all objects matching the objectPrefix from
// the specified bucket. If recursion is enabled it would list
// all subdirectories and all its contents.
//
// Your input paramters are just bucketName, objectPrefix and recursive. If you
// enable recursive as 'true' this function will return back all the
// objects in a given bucket name and object prefix.
//
//   api := client.New(....)
//   recursive := true
//   for message := range api.ListObjects("mytestbucket", "starthere", recursive) {
//       fmt.Println(message)
//   }
//
func (c Client) ListObjects(bucketName, objectPrefix string, recursive bool, doneCh <-chan struct{}) <-chan ObjectStat {
	// Allocate new list objects channel.
	objectStatCh := make(chan ObjectStat, 1000)
	// Default listing is delimited at "/"
	delimiter := "/"
	if recursive {
		// If recursive we do not delimit.
		delimiter = ""
	}
	// Validate bucket name.
	if err := isValidBucketName(bucketName); err != nil {
		defer close(objectStatCh)
		objectStatCh <- ObjectStat{
			Err: err,
		}
		return objectStatCh
	}
	// Validate incoming object prefix.
	if err := isValidObjectPrefix(objectPrefix); err != nil {
		defer close(objectStatCh)
		objectStatCh <- ObjectStat{
			Err: err,
		}
		return objectStatCh
	}

	// Initiate list objects goroutine here.
	go func(objectStatCh chan<- ObjectStat) {
		defer close(objectStatCh)
		// Save marker for next request.
		var marker string
		for {
			// Get list of objects a maximum of 1000 per request.
			result, err := c.listObjectsQuery(bucketName, objectPrefix, marker, delimiter, 1000)
			if err != nil {
				objectStatCh <- ObjectStat{
					Err: err,
				}
				return
			}

			// If contents are available loop through and send over channel.
			for _, object := range result.Contents {
				// Save the marker.
				marker = object.Key
				select {
				// Send object content.
				case objectStatCh <- object:
				// If receives done from the caller, return here.
				case <-doneCh:
					return
				}
			}

			// Send all common prefixes if any.
			// NOTE: prefixes are only present if the request is delimited.
			for _, obj := range result.CommonPrefixes {
				object := ObjectStat{}
				object.Key = obj.Prefix
				object.Size = 0
				select {
				// Send object prefixes.
				case objectStatCh <- object:
				// If receives done from the caller, return here.
				case <-doneCh:
					return
				}
			}

			// If next marker present, save it for next request.
			if result.NextMarker != "" {
				marker = result.NextMarker
			}

			// Listing ends result is not truncated, return right here.
			if !result.IsTruncated {
				return
			}
		}
	}(objectStatCh)
	return objectStatCh
}

/// Bucket Read Operations.

// listObjects - (List Objects) - List some or all (up to 1000) of the objects in a bucket.
//
// You can use the request parameters as selection criteria to return a subset of the objects in a bucket.
// request paramters :-
// ---------
// ?marker - Specifies the key to start with when listing objects in a bucket.
// ?delimiter - A delimiter is a character you use to group keys.
// ?prefix - Limits the response to keys that begin with the specified prefix.
// ?max-keys - Sets the maximum number of keys returned in the response body.
func (c Client) listObjectsQuery(bucketName, objectPrefix, objectMarker, delimiter string, maxkeys int) (listBucketResult, error) {
	// Validate bucket name.
	if err := isValidBucketName(bucketName); err != nil {
		return listBucketResult{}, err
	}
	// Validate object prefix.
	if err := isValidObjectPrefix(objectPrefix); err != nil {
		return listBucketResult{}, err
	}
	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	// Set object prefix.
	urlValues.Set("prefix", urlEncodePath(objectPrefix))
	// Set object marker.
	urlValues.Set("marker", urlEncodePath(objectMarker))
	// Set delimiter.
	urlValues.Set("delimiter", delimiter)
	// Set max keys.
	urlValues.Set("max-keys", fmt.Sprintf("%d", maxkeys))

	// Initialize a new request.
	req, err := c.newRequest("GET", requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	})
	if err != nil {
		return listBucketResult{}, err
	}
	// Execute list buckets.
	resp, err := c.do(req)
	defer closeResponse(resp)
	if err != nil {
		return listBucketResult{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return listBucketResult{}, HTTPRespToErrorResponse(resp, bucketName, "")
		}
	}
	// Decode listBuckets XML.
	listBucketResult := listBucketResult{}
	err = xmlDecoder(resp.Body, &listBucketResult)
	if err != nil {
		return listBucketResult, err
	}
	return listBucketResult, nil
}

// ListIncompleteUploads - List incompletely uploaded multipart objects.
//
// ListIncompleteUploads lists all incompleted objects matching the
// objectPrefix from the specified bucket. If recursion is enabled
// it would list all subdirectories and all its contents.
//
// Your input paramters are just bucketName, objectPrefix and recursive.
// If you enable recursive as 'true' this function will return back all
// the multipart objects in a given bucket name.
//
//   api := client.New(....)
//   recursive := true
//   for message := range api.ListIncompleteUploads("mytestbucket", "starthere", recursive) {
//       fmt.Println(message)
//   }
//
func (c Client) ListIncompleteUploads(bucketName, objectPrefix string, recursive bool, doneCh <-chan struct{}) <-chan ObjectMultipartStat {
	// Turn on size aggregation of individual parts.
	isAggregateSize := true
	return c.listIncompleteUploads(bucketName, objectPrefix, recursive, isAggregateSize, doneCh)
}

// listIncompleteUploads lists all incomplete uploads.
func (c Client) listIncompleteUploads(bucketName, objectPrefix string, recursive, aggregateSize bool, doneCh <-chan struct{}) <-chan ObjectMultipartStat {
	// Allocate channel for multipart uploads.
	objectMultipartStatCh := make(chan ObjectMultipartStat, 1000)
	// Delimiter is set to "/" by default.
	delimiter := "/"
	if recursive {
		// If recursive do not delimit.
		delimiter = ""
	}
	// Validate bucket name.
	if err := isValidBucketName(bucketName); err != nil {
		defer close(objectMultipartStatCh)
		objectMultipartStatCh <- ObjectMultipartStat{
			Err: err,
		}
		return objectMultipartStatCh
	}
	// Validate incoming object prefix.
	if err := isValidObjectPrefix(objectPrefix); err != nil {
		defer close(objectMultipartStatCh)
		objectMultipartStatCh <- ObjectMultipartStat{
			Err: err,
		}
		return objectMultipartStatCh
	}
	go func(objectMultipartStatCh chan<- ObjectMultipartStat) {
		defer close(objectMultipartStatCh)
		// object and upload ID marker for future requests.
		var objectMarker string
		var uploadIDMarker string
		for {
			// list all multipart uploads.
			result, err := c.listMultipartUploadsQuery(bucketName, objectMarker, uploadIDMarker, objectPrefix, delimiter, 1000)
			if err != nil {
				objectMultipartStatCh <- ObjectMultipartStat{
					Err: err,
				}
				return
			}
			// Save objectMarker and uploadIDMarker for next request.
			objectMarker = result.NextKeyMarker
			uploadIDMarker = result.NextUploadIDMarker
			// Send all multipart uploads.
			for _, obj := range result.Uploads {
				// Calculate total size of the uploaded parts if 'aggregateSize' is enabled.
				if aggregateSize {
					// Get total multipart size.
					obj.Size, err = c.getTotalMultipartSize(bucketName, obj.Key, obj.UploadID)
					if err != nil {
						objectMultipartStatCh <- ObjectMultipartStat{
							Err: err,
						}
					}
				}
				select {
				// Send individual uploads here.
				case objectMultipartStatCh <- obj:
				// If done channel return here.
				case <-doneCh:
					return
				}
			}
			// Send all common prefixes if any.
			// NOTE: prefixes are only present if the request is delimited.
			for _, obj := range result.CommonPrefixes {
				object := ObjectMultipartStat{}
				object.Key = obj.Prefix
				object.Size = 0
				select {
				// Send delimited prefixes here.
				case objectMultipartStatCh <- object:
				// If done channel return here.
				case <-doneCh:
					return
				}
			}
			// Listing ends if result not truncated, return right here.
			if !result.IsTruncated {
				return
			}
		}
	}(objectMultipartStatCh)
	// return.
	return objectMultipartStatCh
}

// listMultipartUploads - (List Multipart Uploads).
//   - Lists some or all (up to 1000) in-progress multipart uploads in a bucket.
//
// You can use the request parameters as selection criteria to return a subset of the uploads in a bucket.
// request paramters. :-
// ---------
// ?key-marker - Specifies the multipart upload after which listing should begin.
// ?upload-id-marker - Together with key-marker specifies the multipart upload after which listing should begin.
// ?delimiter - A delimiter is a character you use to group keys.
// ?prefix - Limits the response to keys that begin with the specified prefix.
// ?max-uploads - Sets the maximum number of multipart uploads returned in the response body.
func (c Client) listMultipartUploadsQuery(bucketName, keyMarker, uploadIDMarker, prefix, delimiter string, maxUploads int) (listMultipartUploadsResult, error) {
	// Get resources properly escaped and lined up before using them in http request.
	urlValues := make(url.Values)
	// Set uploads.
	urlValues.Set("uploads", "")
	// Set object key marker.
	urlValues.Set("key-marker", urlEncodePath(keyMarker))
	// Set upload id marker.
	urlValues.Set("upload-id-marker", uploadIDMarker)
	// Set prefix marker.
	urlValues.Set("prefix", urlEncodePath(prefix))
	// Set delimiter.
	urlValues.Set("delimiter", delimiter)
	// Set max-uploads.
	urlValues.Set("max-uploads", fmt.Sprintf("%d", maxUploads))

	// Instantiate a new request.
	req, err := c.newRequest("GET", requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	})
	if err != nil {
		return listMultipartUploadsResult{}, err
	}
	// Execute list multipart uploads request.
	resp, err := c.do(req)
	defer closeResponse(resp)
	if err != nil {
		return listMultipartUploadsResult{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return listMultipartUploadsResult{}, HTTPRespToErrorResponse(resp, bucketName, "")
		}
	}
	// Decode response body.
	listMultipartUploadsResult := listMultipartUploadsResult{}
	err = xmlDecoder(resp.Body, &listMultipartUploadsResult)
	if err != nil {
		return listMultipartUploadsResult, err
	}
	return listMultipartUploadsResult, nil
}

// listObjectParts list all object parts recursively.
func (c Client) listObjectParts(bucketName, objectName, uploadID string) (partsInfo map[int]objectPart, err error) {
	// Part number marker for the next batch of request.
	var nextPartNumberMarker int
	partsInfo = make(map[int]objectPart)
	for {
		// Get list of uploaded parts a maximum of 1000 per request.
		listObjPartsResult, err := c.listObjectPartsQuery(bucketName, objectName, uploadID, nextPartNumberMarker, 1000)
		if err != nil {
			return nil, err
		}
		// Append to parts info.
		for _, part := range listObjPartsResult.ObjectParts {
			// Trim off the odd double quotes from ETag in the beginning and end.
			part.ETag = strings.TrimPrefix(part.ETag, "\"")
			part.ETag = strings.TrimSuffix(part.ETag, "\"")
			partsInfo[part.PartNumber] = part
		}
		// Keep part number marker, for the next iteration.
		nextPartNumberMarker = listObjPartsResult.NextPartNumberMarker
		// Listing ends result is not truncated, return right here.
		if !listObjPartsResult.IsTruncated {
			break
		}
	}

	// Return all the parts.
	return partsInfo, nil
}

// findUploadID lists all incomplete uploads and finds the uploadID of the matching object name.
func (c Client) findUploadID(bucketName, objectName string) (string, error) {
	// Make list incomplete uploads recursive.
	isRecursive := true
	// Turn off size aggregation of individual parts, in this request.
	isAggregateSize := false
	// NOTE: done Channel is set to 'nil, this will drain go routine until exhaustion.
	for mpUpload := range c.listIncompleteUploads(bucketName, objectName, isRecursive, isAggregateSize, nil) {
		if mpUpload.Err != nil {
			return "", mpUpload.Err
		}
		// if object name found, return the upload id.
		if objectName == mpUpload.Key {
			return mpUpload.UploadID, nil
		}
	}
	// No upload id was found, return success and empty upload id.
	return "", nil
}

// getTotalMultipartSize - calculate total uploaded size for the a given multipart object.
func (c Client) getTotalMultipartSize(bucketName, objectName, uploadID string) (size int64, err error) {
	// Iterate over all parts and aggregate the size.
	partsInfo, err := c.listObjectParts(bucketName, objectName, uploadID)
	if err != nil {
		return 0, err
	}
	for _, partInfo := range partsInfo {
		size += partInfo.Size
	}
	return size, nil
}

// listObjectPartsQuery (List Parts query)
//     - lists some or all (up to 1000) parts that have been uploaded for a specific multipart upload
//
// You can use the request parameters as selection criteria to return a subset of the uploads in a bucket.
// request paramters :-
// ---------
// ?part-number-marker - Specifies the part after which listing should begin.
func (c Client) listObjectPartsQuery(bucketName, objectName, uploadID string, partNumberMarker, maxParts int) (listObjectPartsResult, error) {
	// Get resources properly escaped and lined up before using them in http request.
	urlValues := make(url.Values)
	// Set part number marker.
	urlValues.Set("part-number-marker", fmt.Sprintf("%d", partNumberMarker))
	// Set upload id.
	urlValues.Set("uploadId", uploadID)
	// Set max parts.
	urlValues.Set("max-parts", fmt.Sprintf("%d", maxParts))

	req, err := c.newRequest("GET", requestMetadata{
		bucketName:  bucketName,
		objectName:  objectName,
		queryValues: urlValues,
	})
	if err != nil {
		return listObjectPartsResult{}, err
	}
	// Exectue list object parts.
	resp, err := c.do(req)
	defer closeResponse(resp)
	if err != nil {
		return listObjectPartsResult{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return listObjectPartsResult{}, HTTPRespToErrorResponse(resp, bucketName, objectName)
		}
	}
	// Decode list object parts XML.
	listObjectPartsResult := listObjectPartsResult{}
	err = xmlDecoder(resp.Body, &listObjectPartsResult)
	if err != nil {
		return listObjectPartsResult, err
	}
	return listObjectPartsResult, nil
}
