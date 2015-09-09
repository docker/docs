+++
title = "Docker Trusted Registry Org Namespace API"
description = "Docker Trusted Registry 1.3 Organization owned Namespace API"
keywords = ["API, Docker, index, REST, documentation, Docker Trusted Registry, registry"]
weight = 69
[menu.main]
parent = "smn_dtrapi"
+++

# Docker Trusted Registry 1.3 Organization-Owned Repository Namespace Access

Teams within an organization account can be granted `read-only`, `read-write`,
or `admin` level access to the entire namespace of repositories owned
by that organization.

**`read-only`**

- members of the team can pull from any repository under the organization's
  namespace

**`read-write`**

- members of the team can pull from any repository under the organization's
  namespace
- members of the team can push to any repository under the organization's
  namespace

**`admin`**

- members of the team can pull from any repository under the organization's
  namespace
- members of the team can push to any repository under the organization's
  namespace
- members of the team can manage other team's access to any repository under
  the organization's namespace
- members of the team can create new repositories under the organization's
  namespace

## List Teams Granted Access to an Organization-Owned Namespace of Repositories

`GET /api/v0/repositoryNamespaces/{namespace}/teamAccess`

```bash
$ curl -v --user admin:password --insecure -X GET https://dtr.domain.com/api/v0/repositoriesNamespaces/engineering/teamAccess
```

Example Response:

```json
{
  "teamAccessList": [
    {
      "accessLevel": "read-write",
      "team": {
        "id": 7,
        "orgID": 22,
        "type": "managed",
        "name": "dev",
        "description": "Developers"
      }
    },
    {
      "accessLevel": "read-only",
      "team": {
        "id": 8,
        "orgID": 22,
        "type": "managed",
        "name": "qa",
        "description": "Quality Assurance"
      }
    }
  ],
  "namespace": {
    "id": 22,
    "type": "organization",
    "name": "engineering"
  }
}
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
namespace.

**Status Codes**

- *400* the namespace is not owned by an organization.
- *403* the client is not authorized.
- *404*
- *200* success.

## Get a Team's Granted Access to an Organization-Owned Namespace of Repositories

`GET /api/v0/repositoryNamespaces/{namespace}/teamAccess/{teamname}`

```bash
$ curl -v --user admin:password --insecure -X GET https://dtr.domain.com/api/v0/repositoryNamespaces/engineering/teamAccess/lead
```

Example Response:

```json
{
  "accessLevel": "read-only",
  "team": {
    "id": 8,
    "orgID": 22,
    "type": "managed",
    "name": "qa",
    "description": "Quality Assurance"
  },
  "namespace": {
    "id": 22,
    "type": "organization",
    "name": "engineering"
  }
}
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
namespace, is a system admin, member of the organization's "owners" team, or
is a member of the team in question.

**Status Codes**

- *400* the namespace is not owned by an organization.
- *403* the client is not authorized.
- *200* success.

## Set a Team's Access to an Organization-Owned Namespace of Repositories

`PUT /api/v0/repositoryNamespaces/{namespace}/teamAccess/{teamname}`

```bash
curl -v --user admin:password --insecure -X PUT --header "Content-type: application/json" --data '{"accessLevel":"admin"}' https://dtr.domain.com/api/v0/repositoryNamespaces/engineering/teamAccess/lead
```

Example Request:

```http
PUT /api/v0/repositoryNamespaces/engineering/teamAccess/8` HTTP/1.1
Content-Type: application/json

{
  "accessLevel": "read-only",
}
```

Example Response:

```json
{
  "accessLevel": "read-only",
  "team": {
    "id": 8,
    "orgID": 22,
    "type": "managed",
    "name": "qa",
    "description": "Quality Assurance"
  },
  "namespace": {
    "id": 22,
    "type": "organization",
    "name": "engineering"
  }
}
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
namespace.

**Status Codes**

- *400* namespace is not owned by an organization.
- *400* the team does not belong to the owning organization.
- *403* the client is unauthorized.
- *200* success.

## Revoke a Team's Access to an Organization-Owned Namespace of Repositories

`DELETE /api/v0/repositoryNamespaces/{namespace}/teamAccess/{teamname}`

```bash
$ curl -v --user admin:password --insecure -X DELETE https://dtr.domain.com/api/v0/repositoryNamespaces/engineering/teamAccess/lead
```

**Authorization**

Client must be authenticated as a user who has 'admin' level access to the
namespace.

**Status Codes**

- *400* the repository is not owned by an organization.
- *403* the client is not authorized.
- *204* (`No Content`) success - or team does not exist in the access list or there is no such team in the organization.
