package servicemap

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
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// ClientGroupsClient is the service Map API Reference
type ClientGroupsClient struct {
	ManagementClient
}

// NewClientGroupsClient creates an instance of the ClientGroupsClient client.
func NewClientGroupsClient(subscriptionID string) ClientGroupsClient {
	return NewClientGroupsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewClientGroupsClientWithBaseURI creates an instance of the ClientGroupsClient client.
func NewClientGroupsClientWithBaseURI(baseURI string, subscriptionID string) ClientGroupsClient {
	return ClientGroupsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get retrieves the specified client group
//
// resourceGroupName is resource group name within the specified subscriptionId. workspaceName is OMS workspace
// containing the resources of interest. clientGroupName is client Group resource name. startTime is UTC date and time
// specifying the start time of an interval. When not specified the service uses DateTime.UtcNow - 10m endTime is UTC
// date and time specifying the end time of an interval. When not specified the service uses DateTime.UtcNow
func (client ClientGroupsClient) Get(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time) (result ClientGroup, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9_-]+`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `[a-zA-Z0-9_][a-zA-Z0-9_-]+[a-zA-Z0-9_]`, Chain: nil}}},
		{TargetValue: clientGroupName,
			Constraints: []validation.Constraint{{Target: "clientGroupName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "clientGroupName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "servicemap.ClientGroupsClient", "Get")
	}

	req, err := client.GetPreparer(resourceGroupName, workspaceName, clientGroupName, startTime, endTime)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ClientGroupsClient) GetPreparer(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clientGroupName":   autorest.Encode("path", clientGroupName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2015-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if startTime != nil {
		queryParameters["startTime"] = autorest.Encode("query", *startTime)
	}
	if endTime != nil {
		queryParameters["endTime"] = autorest.Encode("query", *endTime)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/features/serviceMap/clientGroups/{clientGroupName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ClientGroupsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ClientGroupsClient) GetResponder(resp *http.Response) (result ClientGroup, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetMembersCount returns the approximate number of members in the client group.
//
// resourceGroupName is resource group name within the specified subscriptionId. workspaceName is OMS workspace
// containing the resources of interest. clientGroupName is client Group resource name. startTime is UTC date and time
// specifying the start time of an interval. When not specified the service uses DateTime.UtcNow - 10m endTime is UTC
// date and time specifying the end time of an interval. When not specified the service uses DateTime.UtcNow
func (client ClientGroupsClient) GetMembersCount(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time) (result ClientGroupMembersCount, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9_-]+`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `[a-zA-Z0-9_][a-zA-Z0-9_-]+[a-zA-Z0-9_]`, Chain: nil}}},
		{TargetValue: clientGroupName,
			Constraints: []validation.Constraint{{Target: "clientGroupName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "clientGroupName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "servicemap.ClientGroupsClient", "GetMembersCount")
	}

	req, err := client.GetMembersCountPreparer(resourceGroupName, workspaceName, clientGroupName, startTime, endTime)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "GetMembersCount", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetMembersCountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "GetMembersCount", resp, "Failure sending request")
		return
	}

	result, err = client.GetMembersCountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "GetMembersCount", resp, "Failure responding to request")
	}

	return
}

// GetMembersCountPreparer prepares the GetMembersCount request.
func (client ClientGroupsClient) GetMembersCountPreparer(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clientGroupName":   autorest.Encode("path", clientGroupName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2015-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if startTime != nil {
		queryParameters["startTime"] = autorest.Encode("query", *startTime)
	}
	if endTime != nil {
		queryParameters["endTime"] = autorest.Encode("query", *endTime)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/features/serviceMap/clientGroups/{clientGroupName}/membersCount", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetMembersCountSender sends the GetMembersCount request. The method will close the
// http.Response Body if it receives an error.
func (client ClientGroupsClient) GetMembersCountSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetMembersCountResponder handles the response to the GetMembersCount request. The method always
// closes the http.Response Body.
func (client ClientGroupsClient) GetMembersCountResponder(resp *http.Response) (result ClientGroupMembersCount, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListMembers returns the members of the client group during the specified time interval.
//
// resourceGroupName is resource group name within the specified subscriptionId. workspaceName is OMS workspace
// containing the resources of interest. clientGroupName is client Group resource name. startTime is UTC date and time
// specifying the start time of an interval. When not specified the service uses DateTime.UtcNow - 10m endTime is UTC
// date and time specifying the end time of an interval. When not specified the service uses DateTime.UtcNow top is
// page size to use. When not specified, the default page size is 100 records.
func (client ClientGroupsClient) ListMembers(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time, top *int32) (result ClientGroupMembersCollection, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `[a-zA-Z0-9_-]+`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "workspaceName", Name: validation.Pattern, Rule: `[a-zA-Z0-9_][a-zA-Z0-9_-]+[a-zA-Z0-9_]`, Chain: nil}}},
		{TargetValue: clientGroupName,
			Constraints: []validation.Constraint{{Target: "clientGroupName", Name: validation.MaxLength, Rule: 256, Chain: nil},
				{Target: "clientGroupName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMaximum, Rule: 200, Chain: nil},
					{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "servicemap.ClientGroupsClient", "ListMembers")
	}

	req, err := client.ListMembersPreparer(resourceGroupName, workspaceName, clientGroupName, startTime, endTime, top)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "ListMembers", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListMembersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "ListMembers", resp, "Failure sending request")
		return
	}

	result, err = client.ListMembersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "ListMembers", resp, "Failure responding to request")
	}

	return
}

// ListMembersPreparer prepares the ListMembers request.
func (client ClientGroupsClient) ListMembersPreparer(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time, top *int32) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clientGroupName":   autorest.Encode("path", clientGroupName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2015-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if startTime != nil {
		queryParameters["startTime"] = autorest.Encode("query", *startTime)
	}
	if endTime != nil {
		queryParameters["endTime"] = autorest.Encode("query", *endTime)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/features/serviceMap/clientGroups/{clientGroupName}/members", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListMembersSender sends the ListMembers request. The method will close the
// http.Response Body if it receives an error.
func (client ClientGroupsClient) ListMembersSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListMembersResponder handles the response to the ListMembers request. The method always
// closes the http.Response Body.
func (client ClientGroupsClient) ListMembersResponder(resp *http.Response) (result ClientGroupMembersCollection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListMembersNextResults retrieves the next set of results, if any.
func (client ClientGroupsClient) ListMembersNextResults(lastResults ClientGroupMembersCollection) (result ClientGroupMembersCollection, err error) {
	req, err := lastResults.ClientGroupMembersCollectionPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "ListMembers", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListMembersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "ListMembers", resp, "Failure sending next results request")
	}

	result, err = client.ListMembersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicemap.ClientGroupsClient", "ListMembers", resp, "Failure responding to next results request")
	}

	return
}

// ListMembersComplete gets all elements from the list without paging.
func (client ClientGroupsClient) ListMembersComplete(resourceGroupName string, workspaceName string, clientGroupName string, startTime *date.Time, endTime *date.Time, top *int32, cancel <-chan struct{}) (<-chan ClientGroupMember, <-chan error) {
	resultChan := make(chan ClientGroupMember)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			close(resultChan)
			close(errChan)
		}()
		list, err := client.ListMembers(resourceGroupName, workspaceName, clientGroupName, startTime, endTime, top)
		if err != nil {
			errChan <- err
			return
		}
		if list.Value != nil {
			for _, item := range *list.Value {
				select {
				case <-cancel:
					return
				case resultChan <- item:
					// Intentionally left blank
				}
			}
		}
		for list.NextLink != nil {
			list, err = client.ListMembersNextResults(list)
			if err != nil {
				errChan <- err
				return
			}
			if list.Value != nil {
				for _, item := range *list.Value {
					select {
					case <-cancel:
						return
					case resultChan <- item:
						// Intentionally left blank
					}
				}
			}
		}
	}()
	return resultChan, errChan
}
