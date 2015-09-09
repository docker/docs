+++
title = "Docker Trusted Registry Org Repository API"
description = "Docker Trusted Registry 1.3 Organization owned Repository API"
keywords = ["API, Docker, index, REST, documentation, Docker Trusted Registry, registry"]
weight = 68
[menu.main]
parent = "smn_dtrapi"
+++

# Docker Trusted Registry 1.3 Organization-Owned Repository Access

Teams within an organization account can be granted `read-only`, `read-write`,
or `admin` level access to any repository owned by that organization.

**`read-only`**

- members of the team can pull from the repository

**`read-write`**

- members of the team can pull from the repository
- members of the team can push to the repository

**`admin`**

- members of the team can pull from the repository
- members of the team can push to the repository
- members of the team can manage other team's access to the repository

## List Teams Granted Access to an Organization-Owned Repository

`GET /api/v0/repositories/{namespace}/{reponame}/teamAccess`

```bash
$ curl -v --user admin:password --insecure -X GET https://dtr.domain.com/api/v0/repositories/engineering/public-app/teamAccess
```

Example Response:

```json
{
  "teamAccessList": [
    {
      "accessLevel": "read-write",
      "team": {
        "id": 13,
        "orgID": 22,
        "type": "managed",
        "name": "docs",
        "description": "Documentation Team"
      }
    },
    {
      "accessLevel": "read-only",
      "team": {
        "id": 10,
        "orgID": 22,
        "type": "managed",
        "name": "product",
        "description": "Product Managers"
      }
    }
  ],
  "repository": {
    "id": 51,
    "namespace": "dtr",
    "name": "private-app",
    "shortDescription": "A Private App for Docker Trusted Registry",
    "longDescription": "We're building the next big thing - Tools of Mass Innovation",
    "visibility": "private",
    "status": "ok"
  }
}
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
repository.

**Status Codes**

- *404* the repository is not visible to the client.
- *400* the repository is not owned by an organization.
- *403* the client is not authorized.
- *200* success.

## List Repository Access Grants for a Team

`GET /api/v0/accounts/{name}/teams/{teamname}/repositoryAccess`

```bash
$ curl -v --user admin:password --insecure -X GET https://dtr.domain.com/api/v0/accounts/engineering/teams/dev/repositoryAccess
```

Example Response:

```json
{
  "team": {
    "id": 13,
    "orgID": 22,
    "type": "managed",
    "name": "docs",
    "description": "Documentation Team"
  },
  "repositoryAccessList": [
    {
      "accessLevel": "read-write",
      "repository": {
        "id": 51,
        "namespace": "dtr",
        "name": "private-app",
        "shortDescription": "A Private App for Docker Trusted Registry",
        "longDescription": "We're building the next big thing - Tools of Mass Innovation",
        "visibility": "private",
        "status": "ok"
      }
    },
    {
      "accessLevel": "read-only",
      "repository": {
        "id": 52,
        "namespace": "dtr",
        "name": "private-app-2",
        "shortDescription": "Another Private App for Docker Trusted Registry",
        "longDescription": "We're still building the next big thing - Tools of Mass Innovation",
        "visibility": "private",
        "status": "ok"
      }
    }
  ]
}
```

**Authorization**

Client must be authenticated as a user who owns the organization the team is
in or be a member of that team.

**Status Codes**

- *404* the repository is not visible to the client.
- *400* the repository is not owned by an organization.
- *400* the team does not belong to the organization.
- *403* the client is not authorized.
- *200* success.

## Set a Team's Access to an Organization-Owned Repository

`PUT /api/v0/repositories/{namespace}/{reponame}/teamAccess/{teamname}`

```bash
$ curl -v --user admin:password --insecure -X PUT --header "Content-type: application/json" --data '{"accessLevel":"read-write"}' https://dtr.domain.com/api/v0/repositories/engineering/public-app/teamAccess/dev
```

Example Request:

```http
PUT /api/v0/repositories/dtr/private-app/teamAccess/13 HTTP/1.1
Content-Type: application/json

{
  "accessLevel": "read-write",
}
```

Example Response:

```json
{
  "accessLevel": "read-write",
  "team": {
    "id": 13,
    "orgID": 22,
    "type": "managed",
    "name": "docs",
    "description": "Documentation Team"
  },
  "repository": {
    "id": 51,
    "namespace": "dtr",
    "name": "private-app",
    "shortDescription": "A Private App for Docker Trusted Registry",
    "longDescription": "We're building the next big thing - Tools of Mass Innovation",
    "visibility": "private",
    "status": "ok"
  }
}
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
repository.

**Status Codes**

- *404* the repository is not visible to the client.
- *400* repository is not owned by an organization.
- *400* the team does not belong to the owning organization.
- *403* the client is unauthorized.
- *200* success.

## Revoke a Team's Access to an Organization-Owned Repository

`DELETE /api/v0/repositories/{namespace}/{reponame}/teamAccess/{teamname}`

```bash
$ curl -v --user admin:password --insecure -X DELETE https://dtr.domain.com/api/v0/repositories/engineering/public-app/teamAccess/dev
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
repository.

**Status Codes**

- *404* the repository is not visible to the client.
- *400* the repository is not owned by an organiztion.
- *403* the client is not authorized.
- *204* (`No Content`) success - or team is not in the access list or there is no such team in the organization.
