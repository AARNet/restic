package servicefabric

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// DeployedServicePackageHealthsClient is the client for the DeployedServicePackageHealths methods of the Servicefabric
// service.
type DeployedServicePackageHealthsClient struct {
	ManagementClient
}

// NewDeployedServicePackageHealthsClient creates an instance of the DeployedServicePackageHealthsClient client.
func NewDeployedServicePackageHealthsClient(timeout *int32) DeployedServicePackageHealthsClient {
	return NewDeployedServicePackageHealthsClientWithBaseURI(DefaultBaseURI, timeout)
}

// NewDeployedServicePackageHealthsClientWithBaseURI creates an instance of the DeployedServicePackageHealthsClient
// client.
func NewDeployedServicePackageHealthsClientWithBaseURI(baseURI string, timeout *int32) DeployedServicePackageHealthsClient {
	return DeployedServicePackageHealthsClient{NewWithBaseURI(baseURI, timeout)}
}

// Get get deployed service package healths
//
// nodeName is the name of the node applicationName is the name of the application servicePackageName is the name of
// the service package eventsHealthStateFilter is the filter of the events health state
func (client DeployedServicePackageHealthsClient) Get(nodeName string, applicationName string, servicePackageName string, eventsHealthStateFilter string) (result DeployedServicePackageHealth, err error) {
	req, err := client.GetPreparer(nodeName, applicationName, servicePackageName, eventsHealthStateFilter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabric.DeployedServicePackageHealthsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicefabric.DeployedServicePackageHealthsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabric.DeployedServicePackageHealthsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DeployedServicePackageHealthsClient) GetPreparer(nodeName string, applicationName string, servicePackageName string, eventsHealthStateFilter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationName":    applicationName,
		"nodeName":           autorest.Encode("path", nodeName),
		"servicePackageName": autorest.Encode("path", servicePackageName),
	}

	const APIVersion = "1.0.0"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(eventsHealthStateFilter) > 0 {
		queryParameters["EventsHealthStateFilter"] = autorest.Encode("query", eventsHealthStateFilter)
	}
	if client.Timeout != nil {
		queryParameters["timeout"] = autorest.Encode("query", *client.Timeout)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Nodes/{nodeName}/$/GetApplications/{applicationName}/$/GetServicePackages/{servicePackageName}/$/GetHealth", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DeployedServicePackageHealthsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DeployedServicePackageHealthsClient) GetResponder(resp *http.Response) (result DeployedServicePackageHealth, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Send send deployed service package health
//
// nodeName is the name of the node applicationName is the name of the application serviceManifestName is the name of
// the service manifest deployedServicePackageHealthReport is the report of the deployed service package health
func (client DeployedServicePackageHealthsClient) Send(nodeName string, applicationName string, serviceManifestName string, deployedServicePackageHealthReport DeployedServiceHealthReport) (result String, err error) {
	req, err := client.SendPreparer(nodeName, applicationName, serviceManifestName, deployedServicePackageHealthReport)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabric.DeployedServicePackageHealthsClient", "Send", nil, "Failure preparing request")
		return
	}

	resp, err := client.SendSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicefabric.DeployedServicePackageHealthsClient", "Send", resp, "Failure sending request")
		return
	}

	result, err = client.SendResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabric.DeployedServicePackageHealthsClient", "Send", resp, "Failure responding to request")
	}

	return
}

// SendPreparer prepares the Send request.
func (client DeployedServicePackageHealthsClient) SendPreparer(nodeName string, applicationName string, serviceManifestName string, deployedServicePackageHealthReport DeployedServiceHealthReport) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationName":     applicationName,
		"nodeName":            autorest.Encode("path", nodeName),
		"serviceManifestName": autorest.Encode("path", serviceManifestName),
	}

	const APIVersion = "1.0.0"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if client.Timeout != nil {
		queryParameters["timeout"] = autorest.Encode("query", *client.Timeout)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Nodes/{nodeName}/$/GetApplications/{applicationName}/$/GetServicePackages/{serviceManifestName}/$/ReportHealth", pathParameters),
		autorest.WithJSON(deployedServicePackageHealthReport),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// SendSender sends the Send request. The method will close the
// http.Response Body if it receives an error.
func (client DeployedServicePackageHealthsClient) SendSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// SendResponder handles the response to the Send request. The method always
// closes the http.Response Body.
func (client DeployedServicePackageHealthsClient) SendResponder(resp *http.Response) (result String, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
