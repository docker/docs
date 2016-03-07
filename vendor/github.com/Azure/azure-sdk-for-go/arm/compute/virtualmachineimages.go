package compute

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
// Code generated by Microsoft (R) AutoRest Code Generator 0.14.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
	"net/url"
)

// VirtualMachineImagesClient is the the Compute Management Client.
type VirtualMachineImagesClient struct {
	ManagementClient
}

// NewVirtualMachineImagesClient creates an instance of the
// VirtualMachineImagesClient client.
func NewVirtualMachineImagesClient(subscriptionID string) VirtualMachineImagesClient {
	return NewVirtualMachineImagesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVirtualMachineImagesClientWithBaseURI creates an instance of the
// VirtualMachineImagesClient client.
func NewVirtualMachineImagesClientWithBaseURI(baseURI string, subscriptionID string) VirtualMachineImagesClient {
	return VirtualMachineImagesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets a virtual machine image.
//
func (client VirtualMachineImagesClient) Get(location string, publisherName string, offer string, skus string, version string) (result VirtualMachineImage, err error) {
	req, err := client.GetPreparer(location, publisherName, offer, skus, version)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client VirtualMachineImagesClient) GetPreparer(location string, publisherName string, offer string, skus string, version string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       url.QueryEscape(location),
		"offer":          url.QueryEscape(offer),
		"publisherName":  url.QueryEscape(publisherName),
		"skus":           url.QueryEscape(skus),
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
		"version":        url.QueryEscape(version),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus/{skus}/versions/{version}"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineImagesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client VirtualMachineImagesClient) GetResponder(resp *http.Response) (result VirtualMachineImage, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of virtual machine images.
//
// filter is the filter to apply on the operation.
func (client VirtualMachineImagesClient) List(location string, publisherName string, offer string, skus string, filter string, top *int32, orderby string) (result VirtualMachineImageResourceList, err error) {
	req, err := client.ListPreparer(location, publisherName, offer, skus, filter, top, orderby)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client VirtualMachineImagesClient) ListPreparer(location string, publisherName string, offer string, skus string, filter string, top *int32, orderby string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       url.QueryEscape(location),
		"offer":          url.QueryEscape(offer),
		"publisherName":  url.QueryEscape(publisherName),
		"skus":           url.QueryEscape(skus),
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = filter
	}
	if top != nil {
		queryParameters["$top"] = top
	}
	if len(orderby) > 0 {
		queryParameters["$orderby"] = orderby
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus/{skus}/versions"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineImagesClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client VirtualMachineImagesClient) ListResponder(resp *http.Response) (result VirtualMachineImageResourceList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListOffers gets a list of virtual machine image offers.
//
func (client VirtualMachineImagesClient) ListOffers(location string, publisherName string) (result VirtualMachineImageResourceList, err error) {
	req, err := client.ListOffersPreparer(location, publisherName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListOffers", nil, "Failure preparing request")
	}

	resp, err := client.ListOffersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListOffers", resp, "Failure sending request")
	}

	result, err = client.ListOffersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListOffers", resp, "Failure responding to request")
	}

	return
}

// ListOffersPreparer prepares the ListOffers request.
func (client VirtualMachineImagesClient) ListOffersPreparer(location string, publisherName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       url.QueryEscape(location),
		"publisherName":  url.QueryEscape(publisherName),
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// ListOffersSender sends the ListOffers request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineImagesClient) ListOffersSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListOffersResponder handles the response to the ListOffers request. The method always
// closes the http.Response Body.
func (client VirtualMachineImagesClient) ListOffersResponder(resp *http.Response) (result VirtualMachineImageResourceList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListPublishers gets a list of virtual machine image publishers.
//
func (client VirtualMachineImagesClient) ListPublishers(location string) (result VirtualMachineImageResourceList, err error) {
	req, err := client.ListPublishersPreparer(location)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListPublishers", nil, "Failure preparing request")
	}

	resp, err := client.ListPublishersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListPublishers", resp, "Failure sending request")
	}

	result, err = client.ListPublishersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListPublishers", resp, "Failure responding to request")
	}

	return
}

// ListPublishersPreparer prepares the ListPublishers request.
func (client VirtualMachineImagesClient) ListPublishersPreparer(location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       url.QueryEscape(location),
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// ListPublishersSender sends the ListPublishers request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineImagesClient) ListPublishersSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListPublishersResponder handles the response to the ListPublishers request. The method always
// closes the http.Response Body.
func (client VirtualMachineImagesClient) ListPublishersResponder(resp *http.Response) (result VirtualMachineImageResourceList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListSkus gets a list of virtual machine image skus.
//
func (client VirtualMachineImagesClient) ListSkus(location string, publisherName string, offer string) (result VirtualMachineImageResourceList, err error) {
	req, err := client.ListSkusPreparer(location, publisherName, offer)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListSkus", nil, "Failure preparing request")
	}

	resp, err := client.ListSkusSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListSkus", resp, "Failure sending request")
	}

	result, err = client.ListSkusResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute/VirtualMachineImagesClient", "ListSkus", resp, "Failure responding to request")
	}

	return
}

// ListSkusPreparer prepares the ListSkus request.
func (client VirtualMachineImagesClient) ListSkusPreparer(location string, publisherName string, offer string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       url.QueryEscape(location),
		"offer":          url.QueryEscape(offer),
		"publisherName":  url.QueryEscape(publisherName),
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmimage/offers/{offer}/skus"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

// ListSkusSender sends the ListSkus request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineImagesClient) ListSkusSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListSkusResponder handles the response to the ListSkus request. The method always
// closes the http.Response Body.
func (client VirtualMachineImagesClient) ListSkusResponder(resp *http.Response) (result VirtualMachineImageResourceList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
