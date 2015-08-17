+++
title = "DTR accounts API"
description = "Docker Trusted Registry 1.3 User and Organization Accounts"
keywords = ["API, Docker, index, REST, documentation, Docker Trusted Registry, registry"]
weight = 61
[menu.main]
parent = "smn_dtrapi"
+++

# Docker Trusted Registry 1.3 User and Organization Accounts

## Create an Account

`POST /api/v0/accounts`

### Create a Managed User Account

DTR auth settings must be in "managed" mode.

Example Request:

```http
POST /api/v0/accounts HTTP/1.1
Content-Type: application/json

{
  "type": "user",
  "name": "alice",
  "password": "watchThinkFruitNeighbor"
}
```

Response:

```json
{
  "id": 3,
  "type": "user",
  "name": "alice",
  "isActive": false
}
```

**Authorization**

Anyone may create a managed user account, however managed user accounts start
out inactive and the user cannot authenticate until an admin explicitly
activates the user account using the activate user API endpoint.

**Status Codes**

- *400* invalid input.
- *200* success.

### Create a User Account from LDAP

DTR auth settings must be in "ldap" mode.

Example Request:

```http
POST /api/v0/accounts HTTP/1.1
Content-Type: application/json

{
  "type": "user",
  "name": "robert.smith",
  "ldapLogin": "robert.smith",
  "password": "shakeMeanPlainBaseball"
}
```

The `name` field is the requested username to use in DTR while `ldapLogin`
should be the user's LDAP user login attribute. These need only differ if the
user prefers or if the user's LDAP login name is not compatible with valid DTR
usernames.

Response:

```json
{
  "id": 4,
  "type": "user",
  "name": "robert.smith",
  "ldapLogin": "robert.smith"
}
```

**Authorization**

Anyone may create an LDAP user account, however, the account is only created if
the client provides a valid ldap login and password.

**Status Codes**

- *400* invalid input.
- *200* success.

### Create an Organization Account

DTR auth settings must be in "managed" or "ldap" mode.

Example Request:

```http
POST /api/v0/accounts HTTP/1.1
Content-Type: application/json

{
  "type": "organization",
  "name": "engineering"
}
```

Response:

```json
{
  "id": 5,
  "type": "organization",
  "name": "engineering",
  "ldapLogin": "engineering"
}
```

**Authorization**

Client must be authenticated as a global 'admin' user.

**Status Codes**

- *400* invalid input.
- *401* client must be authenticated.
- *403* client must be an admin.
- *200* success.

## List Accounts

`GET /api/v0/accounts`

Example Response:

```json
{
  "accounts": [
    {
      "id": 1,
      "type": "user",
      "name": "admin",
      "isActive": true
    },
    {
      "id": 2,
      "type": "user",
      "name": "jlhawn",
      "isActive": true
    },
    {
      "id": 3,
      "type": "user",
      "name": "alice",
      "isActive": true
    },
    {
      "id": 4,
      "type": "organization",
      "name": "engineering",
    }
  ]
}
```

**Authorization**

Client must be authenticated as any user in the system.

**Status Codes**

- *401* the client is not authenticated.
- *200* success.

## Account Details

`GET /api/v0/accounts/{name}`

Example Request:

```http
GET /api/v0/accounts/alice HTTP/1.1
```

Response:

```json
{
  "id": 3,
  "type": "user",
  "name": "alice",
  "isActive": true
}
```

**Authorization**

Client must be authenticated as any user in the system.

**Status Codes**

- *404* no such account exists.
- *401* the client is not authenticated.
- *200* success.

## Remove an Account

`DELETE /api/v0/accounts/{name}`

Example Request:

```http
DELETE /api/v0/accounts/alice HTTP/1.1
```

**Authorization**

Client must be authenticated as a system 'admin' user.

**Status Codes**

- *401* client must be authenticated.
- *403* client must be an admin.
- *204* success.

## Change a Managed User's Password

`POST /api/v0/accounts/{name}/changePassword`

Example Request:

```http
POST /api/v0/accounts/alice/changePassword HTTP/1.1
Content-Type: application/json

{
  "oldPassword": "watchThinkFruitNeighbor",
  "newPassword": "pinkCloudBehaviorDozen"
}
```

Response:

```json
{
  "id": 5,
  "type": "user",
  "name": "alice",
  "isActive": true
}
```

**Authorization**

Client must be authenticated as the user in question or as an admin user (in
which case the `oldPassword` field may be omitted from the request body.

**Status Codes**

- *400* invalid input.
- *401* client must be authenticated.
- *200* success.

## Activate a Managed User

`PUT /api/v0/accounts/{name}/activate`

Example Request:

```http
PUT /api/v0/accounts/alice/activate HTTP/1.1
```

Response:

```json
{
  "id": 5,
  "type": "user",
  "name": "alice",
  "isActive": true
}
```

**Authorization**

Client must be authenticated as a system 'admin' user.

**Status Codes**

- *404* no such account.
- *401* client must be authenticated.
- *403* client must be an admin.
- *200* success.

## Deactivate a Managed User

`PUT /api/v0/accounts/{name}/deactivate`

Example Request:

```http
PUT /api/v0/accounts/alice/deactivate HTTP/1.1
```

Examlpe Response:

```json
{
  "id": 5,
  "type": "user",
  "name": "alice",
  "isActive": true
}
```

**Authorization**

Client must be authenticated as a system 'admin' user.

**Status Codes**

- *404* no such account.
- *401* client must be authenticated.
- *403* client must be an admin.
- *200* success.

## List Organizations for a User

`GET /api/v0/accounts/{name}/organizations`

Example Request:

```http
GET /api/v0/accounts/alice/organizations HTTP/1.1
```

Example Response:

```json
{
  "organizations": [
    {
      "id": 4,
      "type": "organization",
      "name": "engineering",
    }
  ]
}
```

**Authorization**

Client must be authenticated as a system 'admin' user or the user in question.

**Status Codes**

- *404* no such account.
- *401* client must be authenticated.
- *403* client must be an admin or target user.
- *200* success.
