# Registries

## Registry

> Example

```json
{
  "host": "registry-1.docker.io",
  "is_docker_registry": true,
  "is_ssl": true,
  "name": "Docker Hub",
  "port": 443,
  "resource_uri": "/api/repo/v1/user_namespace/registry/registry-1.docker.io/"
}
```

Represents a registry where repositories are hosted.

This is a [namespaced endpoint](#namespaced-endpoints).


### Attributes

Attribute | Description
--------- | -----------
resource_uri | A unique API endpoint that represents the registry
name | Human-readable name of the registry
host | FQDN of the registry, such as `registry-1.docker.io`
is_docker_registry | Whether this registry is run by Docker
is_ssl | Whether this registry has SSL activated or not
port | The port number where the registry is listening to


## List all registries

```http
GET /api/repo/v1/registry/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

Lists all current registries. Returns a list of `Registry` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/repo/v1/[optional_namespace/]registry/`

### Query Parameters

Parameter | Description
--------- | -----------
uuid | Filter by UUID
name | Filter by registry name
host | Filter by registry host
is_docker_registry | Filter by whether the registry is run by Docker or not. Possible values: 'true' or 'false'


## Get an existing registry

```http
GET /api/repo/v1/registry/registry-1.docker.io/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

Gets all the details of an specific registry

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/v1/[optional_namespace/]registry/(host)/`

### Path Parameters

Parameter | Description
--------- | -----------
host | The host of the registry to retrieve
