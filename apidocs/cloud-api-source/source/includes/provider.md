# Providers

## Provider

> Example

```json
{
    "available": true,
    "label": "Digital Ocean",
    "name": "digitalocean",
    "regions": [
        "/api/infra/v1/region/digitalocean/ams1/",
        "/api/infra/v1/region/digitalocean/ams2/",
        "/api/infra/v1/region/digitalocean/ams3/",
        "/api/infra/v1/region/digitalocean/lon1/",
        "/api/infra/v1/region/digitalocean/nyc1/",
        "/api/infra/v1/region/digitalocean/nyc2/",
        "/api/infra/v1/region/digitalocean/nyc3/",
        "/api/infra/v1/region/digitalocean/sfo1/",
        "/api/infra/v1/region/digitalocean/sgp1/"
    ],
    "resource_uri": "/api/infra/v1/provider/digitalocean/"
}
```

A provider is a representation of a cloud provider supported by Docker Cloud. Providers have one or more regions where nodes are deployed.


### Attributes

Attribute | Description
--------- | -----------
resource_uri | A unique API endpoint that represents the provider
name | A unique identifier for the provider
label | A user-friendly name for the provider
regions | A list of resource URIs of the regions available in this provider
available | Whether the provider is currently available for new node deployments


## List all providers

```python
import dockercloud

providers = dockercloud.Provider.list()
```

```http
GET /api/infra/v1/provider/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

providerList, err := dockercloud.ListProviders()

if err != nil {
  log.Println(err)
}

log.Println(providerList)
```

```shell
docker-cloud nodecluster provider
```

Lists all supported cloud providers. Returns a list of `Provider` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/provider/`

### Query Parameters

Parameter | Description
--------- | -----------
name | Filter by provider name



## Get an individual provider

```python
import dockercloud

provider = dockercloud.Provider.fetch("digitalocean")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

provider, err := dockercloud.GetProvider("digitalocean")

if err != nil {
  log.Println(err)
}

log.Println(provider)
```

```http
GET /api/infra/v1/provider/digitalocean/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```


Get all the details of a specific provider

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/provider/(name)/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the provider to retrieve
