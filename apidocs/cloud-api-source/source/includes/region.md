# Regions

## Region

> Example

```json
{
    "availability_zones": [],
    "available": true,
    "label": "Amsterdam 2",
    "name": "ams2",
    "node_types": [
        "/api/infra/v1/nodetype/digitalocean/1gb/",
        "/api/infra/v1/nodetype/digitalocean/2gb/",
        "/api/infra/v1/nodetype/digitalocean/4gb/",
        "/api/infra/v1/nodetype/digitalocean/8gb/",
        "/api/infra/v1/nodetype/digitalocean/16gb/",
        "/api/infra/v1/nodetype/digitalocean/32gb/",
        "/api/infra/v1/nodetype/digitalocean/48gb/",
        "/api/infra/v1/nodetype/digitalocean/64gb/"
    ],
    "provider": "/api/infra/v1/provider/digitalocean/",
    "resource_uri": "/api/infra/v1/region/digitalocean/ams2/"
}
```

A region is a representation of an entire or a subset of a data center of a cloud provider. It can contain availability zones (depending on the provider) and one or more node types.


### Attributes

Attribute | Description
--------- | -----------
resource_uri | A unique API endpoint that represents the region
name | An identifier for the region
label | A user-friendly name for the region
node_types | A list of resource URIs of the node types available in the region
availability_zones | A list of resource URIs of the availability zones available in the region
provider | The resource URI of the provider of the region
available | Whether the region is currently available for new node deployments


## List all regions

```python
import dockercloud

regions = dockercloud.Region.list()
```

```http
GET /api/infra/v1/region/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

regionList, err := dockercloud.ListRegions()

if err != nil {
  log.Println(err)
}

log.Println(regionList)
```

```shell
docker-cloud nodecluster region
```

Lists all regions of all supported cloud providers. Returns a list of `Region` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/region/`

### Query Parameters

Parameter | Description
--------- | -----------
name | Filter by region name
provider | Filter by resource URI of the target provider



## Get an individual region

```python
import dockercloud

region = dockercloud.Region.fetch("digitalocean/lon1")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

region, err := dockercloud.GetRegion("digitalocean","lon1")

if err != nil {
  log.Println(err)
}

log.Println(region)
```

```http
GET /api/infra/v1/region/digitalocean/lon1/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```


Get all the details of a specific region

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/region/(provider.name)/(name)/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the region to retrieve
provider.name | The name of the provider of the region
