---
description: Submit a product to be listed on Docker Store
keywords: Docker, docker, store, purchase images
title: Bring Your Own License (BYOL) products on Store
---

## What is Bring Your Own License (BYOL)?

Bring Your Own License (BYOL) allows customers with existing software licenses
to easily migrate to the containerized version of the software you make
available on Docker Store.

To see and access an ISV's BYOL product listing, customers simply subscribe to
the product with their Docker ID. We call this **Ungated BYOL**.

ISVs can use the Docker Store/Hub as an entitlement and distribution platform.
Using API’s provided by Docker, ISVs can entitle users and distribute their
Dockerized content to many different audiences:

- Existing customers that want their licensed software made available as Docker containers.
- New customers that are only interested in consuming their software as Docker containers.
- Trial or beta customers, where the ISV can distribute feature or time limited software.

Docker provides a fulfillment service so that ISVs can programmatically entitle
users, by creating subscriptions to their content in Docker Store.

## Ungated BYOL

### Prerequisites and setup

To use Docker as your fulfillment service, an ISV must:
- [Apply and be approved as a Docker Store Vendor Partner](https://goto.docker.com/partners)
- Apply and be approved to list an Ungated BYOL product
- Create one or more Ungated BYOL product plans, in the Docker Store Publisher center.

## Creating an ungated BYOL plan

In Plans & Pricing section of the Publisher Center, ensure the following:
- Price/Month should be set to $0
- There should be no free trial associated with the product
- Under the Pull Requirements dropdown, "Subscribed users only" should be selected.

## API reference for ISVs

### Endpoint, authorization, content

All API requests should be made to: <https://store.docker.com/api/fulfillment/v1/...>

For example, the full path for the "Create Order" API POSTs to: https://store.docker.com/api/fulfillment/v1/orders

All API requests to the fulfillment service must have an "Authorization: Bearer"
header with an authorization string provided by Docker. An example
header is:

```none
Authorization: Bearer 9043ea5c-172a-4d4b-b255-a1dab96fb631
```

ISVs should closely protect their authorization token as if it were a password,
and alert Docker if it has been compromised or needs replacement.

All request and response bodies are encoded with JSON using UTF-8.

### Data structures

### OrderCreateRequest (object)

#### Properties

* partner\_id (PartnerID, required) - Business entity creating this order.
* customer (Customer, optional) - Order customer information.
* items (array[OrderItemCreateRequest], required) - One or more items associated with the order.
* token (string, optional) - If supplied, the customer access token for this order.
* metadata (OrderMetadata, optional) - Key/value strings to be stored with order.

### Order (object)

#### Properties

* order\_id: `ord-93b2dba2-79e3-11e6-8b77-86f30ca893d3` (string, required) - The order id.
* token: `DOCKER-TOKEN-234` (string, required) - The access token created for this order by the fulfillment service.
* docker\_id: `a76808b87b6c11e68b7786f30ca893d3` (string, optional) - The docker id of the user that fulfilled the order. This is not set unless the order is in a fulfilled state.
* state: created (enum, required) - The order state.
  * created
  * fulfilled
  * cancelled
* partner\_id (PartnerID, required) - Business entity that created this order.
* customer (Customer, optional) - Order customer information.
* items (array[OrderItem], required) - One or more items associated with the order.
* metadata (OrderMetadata, optional) - Any key/value strings given in the order creation request.
* created: `2016-06-02T05:10:54Z` (string, required) - An ISO-8601 order creation timestamp.
* updated: `2016-06-02T05:10:54Z` (string, required) - An ISO-8601 order updated timestamp.

### OrderItem (object)

#### Properties

* id: `390745e6-faba-11e6-bc64-92361f002671` (string, required) - The order item id.
* product\_id: `bf8f7c15-0c3b-4dc5-b5b3-1595ba9b589e` (string, required) - The Store product id associated with the order item.
* rate\_plan\_id: `85717ec8-6fcf-4fd9-9dbf-051af0ce1eb3` (string, required) - The Store rate plan id associated with the order item.
* subscription\_start\_date: `2016-06-02T05:10:54Z` (string, optional) - An ISO-8601 timestamp representing the subscription start date. If not specified, the subscription starts at order fulfillment time.
* subscription\_end\_date: `2019-06-02T05:10:54Z` (string, optional) - An ISO-8601 timestamp representing the subscription end date. If not specified, the subscription ends based on the plan duration period.
* pricing\_components (array[PricingComponent], required) - One or more pricing components associated with the order item.
* metadata (OrderItemMetadata, optional) - Any key/value strings given for this item in the order creation request.
* product\_keys (array[ProductKey], optional) - Product keys associated with the order item.

### OrderItemCreateRequest (object)

#### Properties

* sku: ZZ456A (string, optional) - The order item SKU.
* product\_id: `bf8f7c15-0c3b-4dc5-b5b3-1595ba9b589e` (string, optional) - The Store product id associated with the order item.
* rate\_plan\_id: `85717ec8-6fcf-4fd9-9dbf-051af0ce1eb3` (string, optional) - The Store rate plan id associated with the order item.
* subscription\_start\_date: `2016-06-02T05:10:54Z` (string, optional) - An ISO-8601 timestamp representing the subscription start date. If not specified, the subscription starts at order fulfillment time.
* subscription\_end\_date: `2019-06-02T05:10:54Z` (string, optional) - An ISO-8601 timestamp representing the subscription end date. If not specified, the subscription ends based on the plan duration period.
* pricing\_components (array[PricingComponent], required) - One or more pricing components associated with the order item.
* metadata (OrderItemMetadata, optional) - Mapping of key/value strings for this order item.
* product\_keys (array[ProductKeyCreateRequest], optional) - Product keys associated with the order item.

### PricingComponent (object)

#### Properties

* name: `Nodes` (string, required) - The pricing component slug. For example Nodes or Engines.
* value: 25 (number, required) - The quantity for the given pricing component. For example 1 support, 15 docker engines, or 5 private repos.

### ProductKeyCreateRequest (object)

#### Properties

* label: `Production` (string, required) - The human-readable label for the given product key that is displayed to the customer.
* media\_type: `text/plain` (enum, required) - An accepted IANA Media Type for this product key. This suggests to the user interface how to display the product key for the customer to use.
  * text/plain
  * application/json
* file\_name: `production-key.txt` (string, optional) - The file name for the downloaded file if the product key is a blob that requires or allows download by the customer.
* value: `AbKe13894Aksel` (string, required) - The contents of the product key as a string. If the value is blob that cannot be represented as a string, the contents are encoded as a Base64 string.

### ProductKey (object)

#### Properties

* id: `390745e6-faba-11e6-bc64-92361f002671` (string, required) - The product key id.
* order\_item\_id: `85717ec8-6fcf-4fd9-9dbf-051af0ce1eb3` (string, required) - The id of the order item that this product key is associated with.
* label: `Production` (string, required) - The human-readable label for the given product key that is displayed to the customer.
* media\_type: `text/plain` (enum, required) - An accepted IANA Media Type for this product key. This suggests to the user interface how to display the product key for the customer to use.
  * text/plain
  * application/json
* file\_name: `production-key.txt` (string, optional) - The file name for the downloaded file if the product key is a blob that requires or allows download by the customer.
* value: `AbKe13894Aksel` (string, required) - The contents of the product key as a string. If the value is blob that cannot be represented as a string, the contents are encoded as a Base64 string.
* created: `2016-06-02T05:10:54Z` (string, required) - An ISO-8601 product key creation timestamp.
* updated: `2016-06-02T05:10:54Z` (string, required) - An ISO-8601 product key updated timestamp.

### OrderFulfillmentRequest (object)

#### Properties

* state: `fulfilled` (string, required) - The desired order state. Only a fulfilled state may be specified.
* docker\_id: `a76808b87b6c11e68b7786f30ca893d3` (string, required) - The docker id (user or org id) to fulfill the order for.
* eusa\_accepted: true (boolean, required) - Whether or not the EUSA associated with the order has been accepted

### OrderCancellationRequest (object)

#### Properties

* state: `cancelled` (string, required) - The desired order state. Only a cancelled state may be specified.

#### Requests

#### Create order [POST]

Create an order.

* Request (application/json)

  * Attributes (OrderCreateRequest)

#### List orders by token [GET /orders{?token}]

Retrieve an order with the given token. An empty array is returned if no orders for the given token are found.

* Parameters

  * token: `DOCKER-TOKEN-2324234` (string, required) - The order token.
* Request (application/json)

#### List orders by partner [GET /orders{?partner\_id}]

List orders associated with the the given partner. An empty array is returned if there are no orders associated with the partner.

* Parameters

  * partner\_id: partner-2342423 (string, required) - The partner identifier.
* Response 200 (application/json)

  * Attributes (array[Order])

### Order [/orders/{order\_id}]

* Parameters
  * order\_id: `ord-93b2dba2-79e3-11e6-8b77-86f30ca893d3` (string, required) - The order id.

#### Get order [GET]

Retrieve an order by id.

* Request (application/json)
* Response 200 (application/json)

  * Attributes (Order)

#### Update order [PATCH]

A number of operations can be performed by `PATCH`ing an order:

**Fulfill** an order. Fulfilling an order puts it in a fulfilled state, and kicks off a process to create subscriptions for each order item associated with the order.

* Request (application/json)

  * Attributes (OrderFulfillmentRequest)
* Response 200

  * Attributes (Order)

**Cancel** an order. Canceling an order puts it in a cancelled state. The order is frozen once cancelled (that is, no further changes may be made to it).

* Request (application/json)
  * Attributes (OrderCancellationRequest)
* Response 200

  * Attributes (Order)

### Order item product keys [/order-items/{order\_item\_id}/product-keys]

* Parameters
  * order\_item\_id: `ord-93b2dba2-79e3-11e6-8b77-86f30ca893d3` (string, required) - The order item id.

#### List product keys for order item [GET]

Retrieve all product keys for an order item by id. An empty array is returned if the order item does not have any product keys.

* Request (application/json)
* Response 200 (application/json)

  * Attributes (array[ProductKey])

#### Create product key for order item [POST]

Create a product key for an existing order item. Adding new product keys does not affect existing product keys.

* Request (application/json)

  * Attributes (ProductKeyCreateRequest)
* Response 201 (application/json)

  * Attributes (ProductKey)

## Group product keys

A product key is a resource attached to an order item that a publisher manages on behalf of their customer. When an order is fulfilled and subscriptions are created for a customer, the product keys associated with that order item can be accessed by the customer.

### Product key [/product-keys/{product\_key\_id}]

* Parameters
  * product\_key\_id: `23443-93b2dba2-79e3-11e6-8b77-86f30ca893d3` (string, required) - The product key id.

#### Get product key [GET]

Retrieve a product key by id.

* Request (application/json)
* Response 200 (application/json)

  * Attributes (ProductKey)

#### Update product key [PUT]

Update a product key by id. All fields shown are required.

* Request (application/json)

  * Attributes (ProductKeyCreateRequest)
* Response 200 (application/json)

  * Attributes (ProductKey)

#### Delete product key [DELETE]

Delete a product key by id.

* Request (application/json)
* Response 204 (application/json)

## What's next?

More information about the publishing flow can be found [here](publish.md).
