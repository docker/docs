# Availability Zones

## Availability Zone

> Example

```json
{
    "available": true,
    "name": "ap-northeast-1a",
    "region": "/api/infra/v1/region/az/ap-northeast-1/",
    "resource_uri": "/api/infra/v1/az/aws/ap-northeast-1/ap-northeast-1a/"
}
```

An Availability Zone is an isolated location inside a region. Providers that support availability zones: AWS


### Attributes

Attribute | Description
--------- | -----------
available | Whether the availability zone is currently available for new node deployments
name | An identifier for the availability zone
region | The resource URI of the region where the availability zone is allocated
resource_uri | A unique API endpoint that represents the zone


## List all availability zones

```python
import dockercloud

az = dockercloud.AZ.list()
```

```http
GET /api/infra/v1/az/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

az, err := dockercloud.ListAZ()
if err != nil {
	log.Println(err)
}

log.Println(az)
```

```shell
docker-cloud nodecluster az
```

Lists all availability zones from all regions of all supported cloud providers. Returns a list of `Availability Zone` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/az/`

### Query Parameters

Parameter | Description
--------- | -----------
name | Filter by availability zone name
region | Filter by resource URI of the target region



## Get an individual availability zone

```python
import dockercloud

az = dockercloud.AZ.fetch("aws/sa-east-1/sa-east-1a")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

az, err := dockercloud.GetAZ("aws/sa-east-1/sa-east-1a")
if err != nil {
	log.Println(err)
}

log.Println(az)
```

```http
GET /api/infra/v1/az/aws/sa-east-1/sa-east-1a/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

Get all the details of a specific availability zone

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/az/(provider.name)/(region.name)/(name)/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the availability zone to retrieve
provider | The name of the provider
region | The name of the region
