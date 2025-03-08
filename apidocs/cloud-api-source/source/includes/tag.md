# Tags

## Tag

> Example

```json
{
    "name": "byon=false",
    "origin": "tutum"
}
```

Tags are used to target the deployment of services to a specific set of nodes. [Learn more](/docker-cloud/apps/deploy-tags/)

### Attributes

Attribute | Description
--------- | -----------
name | Name of the tag
origin | Possible values: `user`, `tutum`


## List all node tags

```http
GET /api/infra/v1/tag/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```
docker-cloud tag ls 7eaf7fff-882c
```

Lists all tags used by all nodes. Returns a list of `Tag` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/tag/`

### Query Parameters

Parameter | Description
--------- | -----------
name | Filter by name
origin | Filter by origin. Possible values: `user`, `tutum`
