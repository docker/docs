---
title: SCIM integration
description: Learn how to configure SCIM with UCP.
keywords: authentication, authorization, SSO, SAML, SCIM, management
---

>{% include enterprise_label_shortform.md %}

# SCIM integration

Simple-Cloud-Identity-Management/System-for-Cross-domain-Identity-Management (SCIM) provides an LDAP alternative for provisioning and managing users and groups, as well as syncing users and groups with an upstream identity provider. Using SCIM schema and API, you can utilize Single Sign-on services (SSO) across various tools.

Prior to Docker Enterprise 3.0, when deactivating a user or changing a user’s group membership association 
in the identity provider, these events were not synchronized with UCP (the service provider). 
You were required to manually change the status and group membership of the user, and possibly revoke the client bundle. 
SCIM implementation allows proactive synchronization with UCP and eliminates this manual intervention.

## Supported identity providers

- Okta 3.2.0

    - References:

        - [System for Cross-domain Identity Management: Core Schema](https://tools.ietf.org/html/rfc7643) 
        - [System for Cross-domain Identity Management: Protocol](https://tools.ietf.org/html/rfc7644) 

## Typical steps involved in SCIM integration:
1. [Configure SCIM for UCP](#configure-scim-for-ucp)
2. [Configure SCIM authentication and access](#scim-authentication-and-access)
3. [Specify user attributes](#specify-user-access)

Other information in this topic includes:

- [Supported SCIM API endpoints](#supported-scim-api-endpoints)
    - [User operations](#user-operations)
    - [Group operations](#group-operations)
    - [Service provider configuration operations](#service-provider-configuration-operations)

### Configure SCIM for UCP
Docker's SCIM implementation utilizes SCIM version 2.0. 

Navigate to **Admin Settings** -> **Authentication and Authorization**. By default, `docker-datacenter` is the organization to which the SCIM team belongs. 
Enter the API token in the UI or have UCP generate a UUID for you.

### Configure SCIM authentication and access
The base URL for all SCIM API calls is `https://<Host IP>/enzi/v0/scim/v2/`. All SCIM methods are accessible API endpoints of this base URL. 

[Bearer Auth](https://swagger.io/docs/specification/authentication/bearer-authentication/) is the API authentication method. When configured, SCIM API endpoints are accessed via the following HTTP header Authorization: 
`Bearer <token>`

Note: 

- SCIM API endpoints are not accessible by any other user (or their token), including the UCP administrator and UCP admin Bearer token.
- As of UCP 3.2.0, an HTTP authentication request header that contains a Bearer token is the only method supported. 

### Specify user attributes 
The following table maps SCIM and SAML attributes to user attribute fields that Docker uses.

  | Docker UCP                | SAML | SCIM |
  | :-----------------------|:-----------------:|:-------------------------:|
  | Account name            | `nameID` in response                          |  userName    |
  | Account full name       | attribute value in `fullname` assertion       |   user's `name.formatted` |
  | Team group link name    | attribute value in `member-of` assertion      |  group's `displayName`    |
  | Team name               |  N/A               | when creating a team, use group's `displayName` +`_SCIM`|
  
  
### Supported SCIM API endpoints

- User operations
    - Retrieve user information
    - Create a new user
    - Update user information
- Group operations
    - Create a new user group
    - Retrieve group information
    - Update user group membership (add/replace/remove users)
- Service provider configuration operations
    - Retrieve service provider resource type metadata
    - Retrieve schema for service provider and SCIM resources
    - Retrieve schema for service provider configuration

#### User operations
For user GET and POST operations:

- Filtering is only supported using the `userName` attribute and `eq` operator. For example, `filter=userName Eq "john"`. 
- Attribute name and attribute operator are case insensitive. For example, the following two expressions evaluate to the same logical value:
    - `filter=userName Eq "john"`
    - `filter=Username eq "john"`
- Pagination is fully supported. 
- Sorting is not supported.

##### GET /Users
Returns a list of SCIM users, 200 users per page by default. Use the `startIndex` and `count` query parameters to paginate long lists of users.

For example, to retrieve the first 20 Users, set `startIndex` to 1 and `count` to 20, provide the following json request:
   ```
   GET Host IP/enzi/v0/scim/v2/Users?startIndex=1&count=20
   Host: example.com
   Accept: application/scim+json
   Authorization: Bearer h480djs93hd8
   ```
   The response to the previous query returns metadata regarding paging that is
   similar to the following example:
   ```
   {
     "totalResults":100,
     "itemsPerPage":20,
     "startIndex":1,
     "schemas":["urn:ietf:params:scim:api:messages:2.0:ListResponse"],
     "Resources":[{
       ...
     }]
   }
   ```
   
##### GET /Users/{id}
Retrieves a single user resource. The value of the `{id}` should be the user's ID. You can also use the `userName` attribute to filter the results.
   ```
   GET {Host IP}/enzi/v0/scim/v2/Users?{user ID}
   Host: example.com
   Accept: application/scim+json
   Authorization: Bearer h480djs93hd8
   ```

##### POST /Users
Creates a user. Must include the `userName` attribute and at least one email address.
   ```
   POST {Host IP}/enzi/v0/scim/v2/Users
   Host: example.com
   Accept: application/scim+json
   Authorization: Bearer h480djs93hd8
   ```
   
##### PATCH /Users/{id}
Updates a user’s `active` status. Inactive users can be reactivated by specifying `"active": true`. Active users can be deactivated by specifying `"active": false`. The value of the `{id}` should be the user's ID. 
   ```
   PATCH {Host IP}/enzi/v0/scim/v2/Users?{user ID}
   Host: example.com
   Accept: application/scim+json
   Authorization: Bearer h480djs93hd8
   ```
   
##### PUT /Users/{id}
Updates existing user information. All attribute values are overwritten, including attributes for 
which empty values or no values were provided. If a previously set attribute value is left blank 
during a `PUT` operation, the value is updated with a blank value in accordance with the attribute 
data type and storage provider. The value of the `{id}` should be the user's ID.

#### Group operations
For group `GET` and `POST` operations:

- Pagination is fully supported.
- Sorting is not supported.

##### GET /Groups/{id}
Retrieves information for a single group.
  ``` 
  GET /scim/v1/Groups?{Group ID}
  Host: example.com
  Accept: application/scim+json
  Authorization: Bearer h480djs93hd8
  ```
  
##### GET /Groups
Returns a paginated list of groups, ten groups per page by default. Use the `startIndex` and `count` query parameters to paginate long lists of groups.
   ```
   GET /scim/v1/Groups?startIndex=4&count=500 HTTP/1.1
   Host: example.com
   Accept: application/scim+json
   Authorization: Bearer h480djs93hd8
   ```
   
##### POST /Groups
Creates a new group. Users can be added to the group during group creation by supplying user ID values in the `members` array.

##### PATCH /Groups/{id}
Updates an existing group resource, allowing individual (or groups of) users to be added or removed from 
the group with a single operation. `Add` is the default operation.

Setting the operation attribute of a member object to `delete` removes members from a group.

##### PUT /Groups/{id}
Updates an existing group resource, overwriting all values for a group even if an attribute is empty or not provided.
`PUT` replaces all members of a group with members provided via the `members` attribute. If a previously set attribute 
is left blank during a `PUT` operation, the new value is set to blank in accordance with the data type of the 
attribute and the storage provider. 

#### Service provider configuration operations
SCIM defines three endpoints to facilitate discovery of SCIM service provider features and schema that can be retrieved using HTTP GET:

##### GET /ResourceTypes
Discovers the resource types available on a SCIM service provider, for example, Users and Groups. Each 
resource type defines the endpoints, the core schema URI that defines the resource, and any supported schema extensions.

##### GET /Schemas
Retrieves information about all supported resource schemas supported by a SCIM service provider.

##### GET /ServiceProviderConfig
Returns a JSON structure that describes the SCIM specification features available on a service provider using a `schemas` attribute of `urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig`.
