package migrate

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
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// GroupsClient is the move your workloads to Azure.
type GroupsClient struct {
	BaseClient
}

// NewGroupsClient creates an instance of the GroupsClient client.
func NewGroupsClient(subscriptionID string, acceptLanguage string) GroupsClient {
	return NewGroupsClientWithBaseURI(DefaultBaseURI, subscriptionID, acceptLanguage)
}

// NewGroupsClientWithBaseURI creates an instance of the GroupsClient client.
func NewGroupsClientWithBaseURI(baseURI string, subscriptionID string, acceptLanguage string) GroupsClient {
	return GroupsClient{NewWithBaseURI(baseURI, subscriptionID, acceptLanguage)}
}

// Create create a new group by sending a json object of type 'group' as given in Models section as part of the Request
// Body. The group name in a project is unique. Labels can be applied on a group as part of creation.
//
// If a group with the groupName specified in the URL already exists, then this call acts as an update.
//
// This operation is Idempotent.
//
// resourceGroupName is name of the Azure Resource Group that project is part of. projectName is name of the Azure
// Migrate project. groupName is unique name of a group within a project. group is new or Updated Group object.
func (client GroupsClient) Create(ctx context.Context, resourceGroupName string, projectName string, groupName string, group *Group) (result Group, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: group,
			Constraints: []validation.Constraint{{Target: "group", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "group.GroupProperties", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "group.GroupProperties.Machines", Name: validation.Null, Rule: true, Chain: nil}}},
				}}}}}); err != nil {
		return result, validation.NewError("migrate.GroupsClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, resourceGroupName, projectName, groupName, group)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client GroupsClient) CreatePreparer(ctx context.Context, resourceGroupName string, projectName string, groupName string, group *Group) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"groupName":         autorest.Encode("path", groupName),
		"projectName":       autorest.Encode("path", projectName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-02-02"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/projects/{projectName}/groups/{groupName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if group != nil {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithJSON(group))
	}
	if len(client.AcceptLanguage) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("Accept-Language", autorest.String(client.AcceptLanguage)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client GroupsClient) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client GroupsClient) CreateResponder(resp *http.Response) (result Group, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError, http.StatusServiceUnavailable),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete the group from the project. The machines remain in the project. Deleting a non-existent group results
// in a no-operation.
//
// A group is an aggregation mechanism for machines in a project. Therefore, deleting group does not delete machines in
// it.
//
// resourceGroupName is name of the Azure Resource Group that project is part of. projectName is name of the Azure
// Migrate project. groupName is unique name of a group within a project.
func (client GroupsClient) Delete(ctx context.Context, resourceGroupName string, projectName string, groupName string) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(ctx, resourceGroupName, projectName, groupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client GroupsClient) DeletePreparer(ctx context.Context, resourceGroupName string, projectName string, groupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"groupName":         autorest.Encode("path", groupName),
		"projectName":       autorest.Encode("path", projectName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-02-02"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/projects/{projectName}/groups/{groupName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(client.AcceptLanguage) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("Accept-Language", autorest.String(client.AcceptLanguage)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client GroupsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client GroupsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent, http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusServiceUnavailable),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get get information related to a specific group in the project. Returns a json object of type 'group' as specified
// in the models section.
//
// resourceGroupName is name of the Azure Resource Group that project is part of. projectName is name of the Azure
// Migrate project. groupName is unique name of a group within a project.
func (client GroupsClient) Get(ctx context.Context, resourceGroupName string, projectName string, groupName string) (result Group, err error) {
	req, err := client.GetPreparer(ctx, resourceGroupName, projectName, groupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client GroupsClient) GetPreparer(ctx context.Context, resourceGroupName string, projectName string, groupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"groupName":         autorest.Encode("path", groupName),
		"projectName":       autorest.Encode("path", projectName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-02-02"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/projects/{projectName}/groups/{groupName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(client.AcceptLanguage) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("Accept-Language", autorest.String(client.AcceptLanguage)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client GroupsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client GroupsClient) GetResponder(resp *http.Response) (result Group, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusServiceUnavailable),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByProject get all groups created in the project. Returns a json array of objects of type 'group' as specified in
// the Models section.
//
// resourceGroupName is name of the Azure Resource Group that project is part of. projectName is name of the Azure
// Migrate project.
func (client GroupsClient) ListByProject(ctx context.Context, resourceGroupName string, projectName string) (result GroupResultList, err error) {
	req, err := client.ListByProjectPreparer(ctx, resourceGroupName, projectName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "ListByProject", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByProjectSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "ListByProject", resp, "Failure sending request")
		return
	}

	result, err = client.ListByProjectResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.GroupsClient", "ListByProject", resp, "Failure responding to request")
	}

	return
}

// ListByProjectPreparer prepares the ListByProject request.
func (client GroupsClient) ListByProjectPreparer(ctx context.Context, resourceGroupName string, projectName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"projectName":       autorest.Encode("path", projectName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-02-02"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/projects/{projectName}/groups", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(client.AcceptLanguage) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("Accept-Language", autorest.String(client.AcceptLanguage)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByProjectSender sends the ListByProject request. The method will close the
// http.Response Body if it receives an error.
func (client GroupsClient) ListByProjectSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByProjectResponder handles the response to the ListByProject request. The method always
// closes the http.Response Body.
func (client GroupsClient) ListByProjectResponder(resp *http.Response) (result GroupResultList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError, http.StatusServiceUnavailable),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
