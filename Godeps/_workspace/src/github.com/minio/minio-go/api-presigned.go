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
	"errors"
	"time"
)

// PresignedGetObject returns a presigned URL to access an object without credentials.
// Expires maximum is 7days - ie. 604800 and minimum is 1.
func (c Client) PresignedGetObject(bucketName, objectName string, expires time.Duration) (string, error) {
	// Input validation.
	if err := isValidBucketName(bucketName); err != nil {
		return "", err
	}
	if err := isValidObjectName(objectName); err != nil {
		return "", err
	}
	if err := isValidExpiry(expires); err != nil {
		return "", err
	}

	expireSeconds := int64(expires / time.Second)
	// Instantiate a new request.
	// Since expires is set newRequest will presign the request.
	req, err := c.newRequest("GET", requestMetadata{
		presignURL: true,
		bucketName: bucketName,
		objectName: objectName,
		expires:    expireSeconds,
	})
	if err != nil {
		return "", err
	}
	return req.URL.String(), nil
}

// PresignedPutObject returns a presigned URL to upload an object without credentials.
// Expires maximum is 7days - ie. 604800 and minimum is 1.
func (c Client) PresignedPutObject(bucketName, objectName string, expires time.Duration) (string, error) {
	// Input validation.
	if err := isValidBucketName(bucketName); err != nil {
		return "", err
	}
	if err := isValidObjectName(objectName); err != nil {
		return "", err
	}
	if err := isValidExpiry(expires); err != nil {
		return "", err
	}

	expireSeconds := int64(expires / time.Second)
	// Instantiate a new request.
	// Since expires is set newRequest will presign the request.
	req, err := c.newRequest("PUT", requestMetadata{
		presignURL: true,
		bucketName: bucketName,
		objectName: objectName,
		expires:    expireSeconds,
	})
	if err != nil {
		return "", err
	}
	return req.URL.String(), nil
}

// PresignedPostPolicy returns POST form data to upload an object at a location.
func (c Client) PresignedPostPolicy(p *PostPolicy) (map[string]string, error) {
	// Validate input arguments.
	if p.expiration.IsZero() {
		return nil, errors.New("Expiration time must be specified")
	}
	if _, ok := p.formData["key"]; !ok {
		return nil, errors.New("object key must be specified")
	}
	if _, ok := p.formData["bucket"]; !ok {
		return nil, errors.New("bucket name must be specified")
	}

	bucketName := p.formData["bucket"]
	// Fetch the location.
	location, err := c.getBucketLocation(bucketName)
	if err != nil {
		return nil, err
	}

	// Keep time.
	t := time.Now().UTC()
	if c.signature.isV2() {
		policyBase64 := p.base64()
		p.formData["policy"] = policyBase64
		// For Google endpoint set this value to be 'GoogleAccessId'.
		if isGoogleEndpoint(c.endpointURL) {
			p.formData["GoogleAccessId"] = c.accessKeyID
		} else {
			// For all other endpoints set this value to be 'AWSAccessKeyId'.
			p.formData["AWSAccessKeyId"] = c.accessKeyID
		}
		// Sign the policy.
		p.formData["signature"] = PostPresignSignatureV2(policyBase64, c.secretAccessKey)
		return p.formData, nil
	}

	// Add date policy.
	p.addNewPolicy(policyCondition{
		matchType: "eq",
		condition: "$x-amz-date",
		value:     t.Format(iso8601DateFormat),
	})
	// Add algorithm policy.
	p.addNewPolicy(policyCondition{
		matchType: "eq",
		condition: "$x-amz-algorithm",
		value:     signV4Algorithm,
	})
	// Add a credential policy.
	credential := getCredential(c.accessKeyID, location, t)
	p.addNewPolicy(policyCondition{
		matchType: "eq",
		condition: "$x-amz-credential",
		value:     credential,
	})
	// get base64 encoded policy.
	policyBase64 := p.base64()
	// Fill in the form data.
	p.formData["policy"] = policyBase64
	p.formData["x-amz-algorithm"] = signV4Algorithm
	p.formData["x-amz-credential"] = credential
	p.formData["x-amz-date"] = t.Format(iso8601DateFormat)
	p.formData["x-amz-signature"] = PostPresignSignatureV4(policyBase64, t, c.secretAccessKey, location)
	return p.formData, nil
}
