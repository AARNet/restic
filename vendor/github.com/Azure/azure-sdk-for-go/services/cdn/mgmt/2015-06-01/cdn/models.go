package cdn

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
)

// CustomDomainResourceState enumerates the values for custom domain resource state.
type CustomDomainResourceState string

const (
	// Active specifies the active state for custom domain resource state.
	Active CustomDomainResourceState = "Active"
	// Creating specifies the creating state for custom domain resource state.
	Creating CustomDomainResourceState = "Creating"
	// Deleting specifies the deleting state for custom domain resource state.
	Deleting CustomDomainResourceState = "Deleting"
)

// EndpointResourceState enumerates the values for endpoint resource state.
type EndpointResourceState string

const (
	// EndpointResourceStateCreating specifies the endpoint resource state creating state for endpoint resource state.
	EndpointResourceStateCreating EndpointResourceState = "Creating"
	// EndpointResourceStateDeleting specifies the endpoint resource state deleting state for endpoint resource state.
	EndpointResourceStateDeleting EndpointResourceState = "Deleting"
	// EndpointResourceStateRunning specifies the endpoint resource state running state for endpoint resource state.
	EndpointResourceStateRunning EndpointResourceState = "Running"
	// EndpointResourceStateStarting specifies the endpoint resource state starting state for endpoint resource state.
	EndpointResourceStateStarting EndpointResourceState = "Starting"
	// EndpointResourceStateStopped specifies the endpoint resource state stopped state for endpoint resource state.
	EndpointResourceStateStopped EndpointResourceState = "Stopped"
	// EndpointResourceStateStopping specifies the endpoint resource state stopping state for endpoint resource state.
	EndpointResourceStateStopping EndpointResourceState = "Stopping"
)

// OriginResourceState enumerates the values for origin resource state.
type OriginResourceState string

const (
	// OriginResourceStateActive specifies the origin resource state active state for origin resource state.
	OriginResourceStateActive OriginResourceState = "Active"
	// OriginResourceStateCreating specifies the origin resource state creating state for origin resource state.
	OriginResourceStateCreating OriginResourceState = "Creating"
	// OriginResourceStateDeleting specifies the origin resource state deleting state for origin resource state.
	OriginResourceStateDeleting OriginResourceState = "Deleting"
)

// ProfileResourceState enumerates the values for profile resource state.
type ProfileResourceState string

const (
	// ProfileResourceStateActive specifies the profile resource state active state for profile resource state.
	ProfileResourceStateActive ProfileResourceState = "Active"
	// ProfileResourceStateCreating specifies the profile resource state creating state for profile resource state.
	ProfileResourceStateCreating ProfileResourceState = "Creating"
	// ProfileResourceStateDeleting specifies the profile resource state deleting state for profile resource state.
	ProfileResourceStateDeleting ProfileResourceState = "Deleting"
	// ProfileResourceStateDisabled specifies the profile resource state disabled state for profile resource state.
	ProfileResourceStateDisabled ProfileResourceState = "Disabled"
)

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// ProvisioningStateCreating specifies the provisioning state creating state for provisioning state.
	ProvisioningStateCreating ProvisioningState = "Creating"
	// ProvisioningStateFailed specifies the provisioning state failed state for provisioning state.
	ProvisioningStateFailed ProvisioningState = "Failed"
	// ProvisioningStateSucceeded specifies the provisioning state succeeded state for provisioning state.
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

// QueryStringCachingBehavior enumerates the values for query string caching behavior.
type QueryStringCachingBehavior string

const (
	// BypassCaching specifies the bypass caching state for query string caching behavior.
	BypassCaching QueryStringCachingBehavior = "BypassCaching"
	// IgnoreQueryString specifies the ignore query string state for query string caching behavior.
	IgnoreQueryString QueryStringCachingBehavior = "IgnoreQueryString"
	// NotSet specifies the not set state for query string caching behavior.
	NotSet QueryStringCachingBehavior = "NotSet"
	// UseQueryString specifies the use query string state for query string caching behavior.
	UseQueryString QueryStringCachingBehavior = "UseQueryString"
)

// ResourceType enumerates the values for resource type.
type ResourceType string

const (
	// MicrosoftCdnProfilesEndpoints specifies the microsoft cdn profiles endpoints state for resource type.
	MicrosoftCdnProfilesEndpoints ResourceType = "Microsoft.Cdn/Profiles/Endpoints"
)

// SkuName enumerates the values for sku name.
type SkuName string

const (
	// Premium specifies the premium state for sku name.
	Premium SkuName = "Premium"
	// Standard specifies the standard state for sku name.
	Standard SkuName = "Standard"
)

// CheckNameAvailabilityInput is input of CheckNameAvailability API.
type CheckNameAvailabilityInput struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// CheckNameAvailabilityOutput is output of check name availability API.
type CheckNameAvailabilityOutput struct {
	autorest.Response `json:"-"`
	NameAvailable     *bool   `json:"NameAvailable,omitempty"`
	Reason            *string `json:"Reason,omitempty"`
	Message           *string `json:"Message,omitempty"`
}

// CustomDomain is CDN CustomDomain represents a mapping between a user specified domain name and a CDN endpoint. This
// is to use custom domain names to represent the URLs for branding purposes.
type CustomDomain struct {
	autorest.Response       `json:"-"`
	ID                      *string `json:"id,omitempty"`
	Name                    *string `json:"name,omitempty"`
	Type                    *string `json:"type,omitempty"`
	*CustomDomainProperties `json:"properties,omitempty"`
}

// CustomDomainListResult is
type CustomDomainListResult struct {
	autorest.Response `json:"-"`
	Value             *[]CustomDomain `json:"value,omitempty"`
}

// CustomDomainParameters is customDomain properties required for custom domain creation or update.
type CustomDomainParameters struct {
	*CustomDomainPropertiesParameters `json:"properties,omitempty"`
}

// CustomDomainProperties is
type CustomDomainProperties struct {
	HostName          *string                   `json:"hostName,omitempty"`
	ResourceState     CustomDomainResourceState `json:"resourceState,omitempty"`
	ProvisioningState ProvisioningState         `json:"provisioningState,omitempty"`
}

// CustomDomainPropertiesParameters is
type CustomDomainPropertiesParameters struct {
	HostName *string `json:"hostName,omitempty"`
}

// DeepCreatedOrigin is deep created origins within a CDN endpoint.
type DeepCreatedOrigin struct {
	Name                         *string `json:"name,omitempty"`
	*DeepCreatedOriginProperties `json:"properties,omitempty"`
}

// DeepCreatedOriginProperties is properties of deep created origin on a CDN endpoint.
type DeepCreatedOriginProperties struct {
	HostName  *string `json:"hostName,omitempty"`
	HTTPPort  *int32  `json:"httpPort,omitempty"`
	HTTPSPort *int32  `json:"httpsPort,omitempty"`
}

// Endpoint is CDN endpoint is the entity within a CDN profile containing configuration information regarding caching
// behaviors and origins. The CDN endpoint is exposed using the URL format <endpointname>.azureedge.net by default, but
// custom domains can also be created.
type Endpoint struct {
	autorest.Response   `json:"-"`
	ID                  *string             `json:"id,omitempty"`
	Name                *string             `json:"name,omitempty"`
	Type                *string             `json:"type,omitempty"`
	Location            *string             `json:"location,omitempty"`
	Tags                *map[string]*string `json:"tags,omitempty"`
	*EndpointProperties `json:"properties,omitempty"`
}

// EndpointCreateParameters is endpoint properties required for new endpoint creation.
type EndpointCreateParameters struct {
	Location                            *string             `json:"location,omitempty"`
	Tags                                *map[string]*string `json:"tags,omitempty"`
	*EndpointPropertiesCreateParameters `json:"properties,omitempty"`
}

// EndpointListResult is
type EndpointListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Endpoint `json:"value,omitempty"`
}

// EndpointProperties is
type EndpointProperties struct {
	HostName                   *string                    `json:"hostName,omitempty"`
	OriginHostHeader           *string                    `json:"originHostHeader,omitempty"`
	OriginPath                 *string                    `json:"originPath,omitempty"`
	ContentTypesToCompress     *[]string                  `json:"contentTypesToCompress,omitempty"`
	IsCompressionEnabled       *bool                      `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed              *bool                      `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed             *bool                      `json:"isHttpsAllowed,omitempty"`
	QueryStringCachingBehavior QueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
	Origins                    *[]DeepCreatedOrigin       `json:"origins,omitempty"`
	ResourceState              EndpointResourceState      `json:"resourceState,omitempty"`
	ProvisioningState          ProvisioningState          `json:"provisioningState,omitempty"`
}

// EndpointPropertiesCreateParameters is
type EndpointPropertiesCreateParameters struct {
	OriginHostHeader           *string                    `json:"originHostHeader,omitempty"`
	OriginPath                 *string                    `json:"originPath,omitempty"`
	ContentTypesToCompress     *[]string                  `json:"contentTypesToCompress,omitempty"`
	IsCompressionEnabled       *bool                      `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed              *bool                      `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed             *bool                      `json:"isHttpsAllowed,omitempty"`
	QueryStringCachingBehavior QueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
	Origins                    *[]DeepCreatedOrigin       `json:"origins,omitempty"`
}

// EndpointPropertiesUpdateParameters is
type EndpointPropertiesUpdateParameters struct {
	OriginHostHeader           *string                    `json:"originHostHeader,omitempty"`
	OriginPath                 *string                    `json:"originPath,omitempty"`
	ContentTypesToCompress     *[]string                  `json:"contentTypesToCompress,omitempty"`
	IsCompressionEnabled       *bool                      `json:"isCompressionEnabled,omitempty"`
	IsHTTPAllowed              *bool                      `json:"isHttpAllowed,omitempty"`
	IsHTTPSAllowed             *bool                      `json:"isHttpsAllowed,omitempty"`
	QueryStringCachingBehavior QueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
}

// EndpointUpdateParameters is endpoint properties required for new endpoint creation.
type EndpointUpdateParameters struct {
	Tags                                *map[string]*string `json:"tags,omitempty"`
	*EndpointPropertiesUpdateParameters `json:"properties,omitempty"`
}

// ErrorResponse is
type ErrorResponse struct {
	autorest.Response `json:"-"`
	Code              *string `json:"code,omitempty"`
	Message           *string `json:"message,omitempty"`
}

// LoadParameters is parameters required for endpoint load.
type LoadParameters struct {
	ContentPaths *[]string `json:"contentPaths,omitempty"`
}

// Operation is CDN REST API operation
type Operation struct {
	Name    *string           `json:"name,omitempty"`
	Display *OperationDisplay `json:"display,omitempty"`
}

// OperationDisplay is
type OperationDisplay struct {
	Provider  *string `json:"provider,omitempty"`
	Resource  *string `json:"resource,omitempty"`
	Operation *string `json:"operation,omitempty"`
}

// OperationListResult is
type OperationListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Operation `json:"value,omitempty"`
}

// Origin is CDN origin is the source of the content being delivered via CDN. When the edge nodes represented by an
// endpoint do not have the requested content cached, they attempt to fetch it from one or more of the configured
// origins.
type Origin struct {
	autorest.Response `json:"-"`
	ID                *string `json:"id,omitempty"`
	Name              *string `json:"name,omitempty"`
	Type              *string `json:"type,omitempty"`
	*OriginProperties `json:"properties,omitempty"`
}

// OriginListResult is
type OriginListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Origin `json:"value,omitempty"`
}

// OriginParameters is origin properties needed for origin creation or update.
type OriginParameters struct {
	*OriginPropertiesParameters `json:"properties,omitempty"`
}

// OriginProperties is
type OriginProperties struct {
	HostName          *string             `json:"hostName,omitempty"`
	HTTPPort          *int32              `json:"httpPort,omitempty"`
	HTTPSPort         *int32              `json:"httpsPort,omitempty"`
	ResourceState     OriginResourceState `json:"resourceState,omitempty"`
	ProvisioningState ProvisioningState   `json:"provisioningState,omitempty"`
}

// OriginPropertiesParameters is
type OriginPropertiesParameters struct {
	HostName  *string `json:"hostName,omitempty"`
	HTTPPort  *int32  `json:"httpPort,omitempty"`
	HTTPSPort *int32  `json:"httpsPort,omitempty"`
}

// Profile is CDN profile represents the top level resource and the entry point into the CDN API. This allows users to
// set up a logical grouping of endpoints in addition to creating shared configuration settings and selecting pricing
// tiers and providers.
type Profile struct {
	autorest.Response  `json:"-"`
	ID                 *string             `json:"id,omitempty"`
	Name               *string             `json:"name,omitempty"`
	Type               *string             `json:"type,omitempty"`
	Location           *string             `json:"location,omitempty"`
	Tags               *map[string]*string `json:"tags,omitempty"`
	*ProfileProperties `json:"properties,omitempty"`
}

// ProfileCreateParameters is profile properties required for profile creation.
type ProfileCreateParameters struct {
	Location                           *string             `json:"location,omitempty"`
	Tags                               *map[string]*string `json:"tags,omitempty"`
	*ProfilePropertiesCreateParameters `json:"properties,omitempty"`
}

// ProfileListResult is
type ProfileListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Profile `json:"value,omitempty"`
}

// ProfileProperties is
type ProfileProperties struct {
	Sku               *Sku                 `json:"sku,omitempty"`
	ResourceState     ProfileResourceState `json:"resourceState,omitempty"`
	ProvisioningState ProvisioningState    `json:"provisioningState,omitempty"`
}

// ProfilePropertiesCreateParameters is
type ProfilePropertiesCreateParameters struct {
	Sku *Sku `json:"sku,omitempty"`
}

// ProfileUpdateParameters is profile properties required for profile update.
type ProfileUpdateParameters struct {
	Tags *map[string]*string `json:"tags,omitempty"`
}

// PurgeParameters is parameters required for endpoint purge.
type PurgeParameters struct {
	ContentPaths *[]string `json:"contentPaths,omitempty"`
}

// Resource is
type Resource struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// Sku is the SKU (pricing tier) of the CDN profile.
type Sku struct {
	Name SkuName `json:"name,omitempty"`
}

// SsoURI is SSO URI required to login to third party web portal.
type SsoURI struct {
	autorest.Response `json:"-"`
	SsoURIValue       *string `json:"ssoUriValue,omitempty"`
}

// TrackedResource is ARM tracked resource
type TrackedResource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// ValidateCustomDomainInput is input of the custom domain to be validated.
type ValidateCustomDomainInput struct {
	HostName *string `json:"hostName,omitempty"`
}

// ValidateCustomDomainOutput is output of custom domain validation.
type ValidateCustomDomainOutput struct {
	autorest.Response     `json:"-"`
	CustomDomainValidated *bool   `json:"customDomainValidated,omitempty"`
	Reason                *string `json:"reason,omitempty"`
	Message               *string `json:"message,omitempty"`
}
