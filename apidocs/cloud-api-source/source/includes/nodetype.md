# Node Types

## Node Type

> Example

```json
{
	"availability_zones": [],
	"available": true,
	"label": "1GB",
	"name": "1gb",
	"provider": "/api/infra/v1/provider/digitalocean/",
	"regions": [
		"/api/infra/v1/region/digitalocean/ams1/",
		"/api/infra/v1/region/digitalocean/sfo1/",
		"/api/infra/v1/region/digitalocean/nyc2/",
		"/api/infra/v1/region/digitalocean/ams2/",
		"/api/infra/v1/region/digitalocean/sgp1/",
		"/api/infra/v1/region/digitalocean/lon1/",
		"/api/infra/v1/region/digitalocean/nyc3/",
		"/api/infra/v1/region/digitalocean/nyc1/"
	],
	"resource_uri": "/api/infra/v1/nodetype/digitalocean/1gb/"
}
```

A node type is a representation of an instance size supported by a certain cloud provider in a certain region and/or availability zone.


### Attributes

Attribute | Description
--------- | -----------
resource_uri | A unique API endpoint that represents the node type
name | An identifier for the node type
label | A user-friendly name for the node type
regions | A list of resource URIs of the regions to which this node type can be deployed to
availability_zones | A list of resource URIs of the availability zones to which this node type can be deployed to
provider | The resource URI of the provider of the node type
available | Whether the node type is currently available for new node deployments


## List all node types

```python
import dockercloud

nodetypes = dockercloud.NodeType.list()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodetypeList, err := dockercloud.ListNodeTypes()

if err != nil {
  log.Println(err)
}

log.Println(nodetypeList)
```

```http
GET /api/infra/v1/nodetype/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud nodecluster nodetype
```

Lists all node types of all supported cloud providers. Returns a list of `NodeType` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/nodetype/`

### Query Parameters

Parameter | Description
--------- | -----------
name | Filter by node type name
regions | Filter by resource URI of the target regions
availability_zones | Filter by resource URI of the target availability zones


## Get an individual node type

```python
import dockercloud

nodetype = dockercloud.NodeType.fetch("digitalocean/1gb")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodetype, err := dockercloud.GetNodeType("digitalocean","1gb")

if err != nil {
  log.Println(err)
}

log.Println(nodetype)
```

```http
GET /api/infra/v1/nodetype/digitalocean/1gb/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```


Get all the details of a specific node type

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/nodetype/(provider.name)/(name)/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the node type to retrieve
provider.name | The name of the provider of the node type
