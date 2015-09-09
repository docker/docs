+++
title = "Docker Trusted Registry User and Org API"
description = "Docker Trusted Registry 1.3 User and Organization Accounts API"
keywords = ["API, Docker, index, REST, documentation, Docker Trusted Registry, registry"]
weight = 62
[menu.main]
parent = "smn_dtrapi"
+++

# Docker Trusted Registry 1.3 Teams

As with Docker HUb, Docker Trusted Registry teams of users can only exist within an Organization.

## List Teams in an Organization

`GET /api/v0/accounts/{name}/teams`

```bash
$ curl --insecure -v --user admin:password https://dtr.domain.com/api/v0/accounts/engineering/teams
```

Example Response:

```json
{
  "teams": [
    {
      "id": 4,
      "orgID": 4,
      "type": "managed",
      "name": "owners",
      "description": ""
    },
    {
      "id": 5,
      "orgID": 4,
      "type": "managed",
      "name": "testers",
      "description": "i have altered the description, pray that i do not alter it any further"
    }
  ]
}
```

**Authorization**

Client must be authenticated as a member of the organization containing the team(s).

**Status Codes**

- *403* the client is not authorized.
- *404* the Organization has no teams.
- *200* success.

## View Details of a Team

`GET /api/v0/accounts/{name}/teams/{teamname}`

```bash
$ curl --insecure -v --user admin:password https://dtr.domain.com/api/v0/accounts/engineering/teams/qa
```

Example Response:

```json
{
  "id": 5,
  "orgID": 4,
  "type": "managed",
  "name": "testers",
  "description": "I have altered the description, pray that I do not alter it any further"
}
```

**Authorization**

Client must be authenticated as a member of the organization containing the team(s).

**Status Codes**

- *403* the client is not authorized.
- *404* no such team exists.
- *200* success.

## List a Team's Members

`GET /api/v0/accounts/{name}/teams/{teamname}/members`

```bash
$ curl --insecure -v --user admin:password https://dtr.domain.com/api/v0/accounts/engineering/teams/qa/members
```

Example Response:

```json
{
  "members": [
    {
      "id": 8,
      "type": "user",
      "name": "midei",
      "isActive": true
    },
    {
      "id": 10,
      "type": "user",
      "name": "rajat",
      "isActive": true
    },
    {
      "id": 12,
      "type": "user",
      "name": "banjot",
      "isActive": true
    },
    {
      "id": 15,
      "type": "user",
      "name": "jon",
      "isActive": true
    }
  ]
}
```

**Authorization**

Client must be authenticated as a system admin, a member of the "owners" team
in the organization, or a member of the team in question.

**Status Codes**

- *403* the client is not authorized.
- *404* no such team exists.
- *200* success.

## Check if a User is a Member of a Team

`GET /api/v0/accounts/{name}/teams/{teamname}/members/{member}`

```bash
$ curl --insecure -v --user admin:password -X GET https://dtr.domain.com/api/v0/accounts/engineering/teams/qa/members/test
```

**Authorization**

Client must be authenticated as a user who has visibility into the team
(i.e., a member of the team or an owner of the organization).

**Status Codes**

- *403* the client is not authorized.
- *404* no such teams exists or user is not a member.
- *204* success (user is a member).

## Create a Team in an Organization

`POST /api/v0/accounts/{name}/teams`

```bash
$ curl --insecure -v --user admin:password -X POST --data '{"name": "qa", "type": "managed"}' --header "Content-type: application/json" https://dtr.domain.com/api/v0/accounts/engineering/teams
```

Example Request:

```http
POST /api/v0/accounts/engineering/teams HTTP/1.1
Content-Type: application/json

{
  "name": "qa",
  "description": "QA Engineering Team",
  "type": "ldap",
  "ldapDN": "cn=qatesters,ou=engineering,ou=groups,dc=example,dc=com",
  "ldapGroupMemberAttribute": "member"
}
```

Example Response:

```json
{
  "id": 5,
  "orgID": 4,
  "type": "ldap",
  "name": "qa",
  "description": "QA Engineering Team",
  "ldapDN": "cn=qatesters,ou=engineering,ou=groups,dc=example,dc=com",
  "ldapGroupMemberAttribute": "member"
}
```

**Authorization**

Client must be authenticated as a system admin or a member of the "owners"
team in the organization.

**Status Codes**

- *403* the client is not authorized.
- *400* invalid team name or LDAP filter.
- *201* success.

## Update a Teams's Details

`PATCH /api/v0/accounts/{name}/teams/{teamname}`

```bash
$ curl --insecure -v --user admin:password -X PATCH --data '{"description":"add one"}' --header "Content-type: application/json" https://dtr.domain.com/api/v0/accounts/engineering/teams/qa
```

Example Request:

```http
POST /api/v0/accounts/engineering/teams/5 HTTP/1.1
Content-Type: application/json

{
  "name": "qualityassurance",
  "description": "Quality Assurance Engineers"
}
```

Example Response:

```json
{
  "id": 5,
  "orgID": 4,
  "type": "ldap",
  "name": "qualityassurance",
  "description": "Quality Assurance Engineers",
  "ldapDN": "cn=qatesters,ou=engineering,ou=groups,dc=example,dc=com",
  "ldapGroupMemberAttribute": "member"
}
```

**Authorization**

Client must be authenticated as a system admin or a member of the "owners"
team in the organization.

**Status Codes**

- *404* no such team exists.
- *403* the client is not authorized.
- *400* invalid updated detail values.
- *200* success.


## Add a User to a Team (if not LDAP synced).

`PUT /api/v0/accounts/{name}/teams/{teamname}/members/{member}`

```bash
$ curl --insecure -v --user admin:password -X PUT https://dtr.domain.com/api/v0/accounts/engineering/teams/qa/members/alice
```

**Authorization**

Client must be authenticated as a system admin or a member of the "owners"
team in the organization.

**Status Codes**

- *403* the client is not authorized.
- *404* no such team or user.
- *200* success.

## Remove a User from a Team (if not LDAP synced).

`DELETE /api/v0/accounts/{name}/teams/{teamname}/members/{member}`

```bash
$ curl --insecure -v --user admin:password -X DELETE https://dtr.domain.com/api/v0/accounts/engineering/teams/qa/members/alice
```

**Authorization**

Client must be authenticated as a system admin or a member of the "owners"
team in the organization.

**Status Codes**

- *403* the client is not authorized.
- *404* no such team exists.
- *204* (`No Content`) success - or user is not in the team.

## Remove a Team.

`DELETE /api/v0/accounts/{name}/teams/{teamname}`

```bash
$ curl --insecure -v --user admin:password -X DELETE https://dtr.domain.com/api/v0/accounts/engineering/teams/qa
```

**Authorization**

Client must be authenticated as a system admin or a member of the "owners"
team in the organization.

**Status Codes**

- *403* the client is not authorized.
- *204* (`No Content`) success - or team does not exist.
