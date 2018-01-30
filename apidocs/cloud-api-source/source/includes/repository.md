# External Repositories

## External Repository

> Example

```json
{
  "in_use": false,
  "name": "my.registry.com/myrepo",
  "registry": "/api/repo/v1/user_namespace/registry/my.registry.com/",
  "resource_uri": "/api/repo/v1/user_namespace/repository/my.registry.com/myrepo/",
}
```

The `repository` endpoint is used to add and remove existing repositories on third party registries to be used in deployments and builds.

This is a [namespaced endpoint](#namespaced-endpoints).

### Attributes

Attribute | Description
--------- | -----------
resource_uri | A unique API endpoint that represents the repository
name | Name of the repository, such as `my.registry.com/myrepo`
in_use | If the image is being used by any of your services
registry | Resource URI of the registry where this image is hosted


## List all external repositories

```python
import dockercloud

repositories = dockercloud.Repository.list()
```

```http
GET /api/repo/v1/repository/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

repositoriesList, err := dockercloud.ListRepositories()

if err != nil {
    log.Println(err)
}

log.Pringln(repositoriesList)
```

```shell
docker-cloud repository ls
```

Lists all added repositories from third party registries. Returns a list of `Repository` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/repo/v1/[optional_namespace/]repository/`

### Query Parameters

Parameter | Description
--------- | -----------
name | Filter by image name
registry | Filter by resource URI of the target repository registry


## Add a new external repository

```python
import dockercloud

repository = dockercloud.Repository.create(name="registry.local/user1/image1", username=username, password=password)
repository.save()
```

```http
POST /api/repo/v1/repository/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{"name": "registry.local/user1/image1", "username": "username", "password": "password"}
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

image, err := dockercloud.CreateImage(dockercloud.ImageCreateRequest{
  Name: "registry.local/user1/image1",
  Username: "username",
  Password: "password"
})
```

```shell
docker-cloud repository register -u username -p password registry.local/user1/image1
```

Adds an existing repository on a third party registry.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/repo/v1/[optional_namespace/]repository/`

### JSON Parameters

Parameter | Description
--------- | -----------
name | Name of the repository, such as 'my.registry.com/myrepo'
username | Username to authenticate with the third party registry
password | Password to authenticate with the third party registry


## Get an external repository details

```python
import dockercloud

repository = dockercloud.Repository.fetch("registry.local/user1/image1")
```

```http
GET /api/repo/v1/repository/registry.local/user1/image1/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

repository, err = dockercloud.GetRepository("registry.local/user1/image1")

if err != nil {
    log.Println(err)
}

log.Println(repository)
```

```shell
docker-cloud repository inspect registry.local/user1/image1
```

Gets all the details of an specific repository

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/repo/v1/[optional_namespace/]repository/(name)/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the repository to retrieve


## Update credentials of an external repository

```python
import dockercloud

repository = dockercloud.Repository.fetch("registry.local/user1/image1")
repository.username = "new username"
repository.password = "new password"
repository.save()
```

```http
PATCH /api/repo/v1/repository/registry.local/user1/image1/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{"username": "username", "password": "password"}
```

```shell
docker-cloud repository update -n "new username" -p "new password" registry.local/user1/image1
```

Updates the external repository credentials.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`PATCH /api/repo/v1/[optional_namespace/]repository/(name)/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the repository to update


### JSON Parameters

Parameter | Description
--------- | -----------
username | Username to authenticate with the private registry
password | Password to authenticate with the private registry


## Remove an external repository

```python
import dockercloud

repository = dockercloud.Repository.fetch("registry.local/user1/image1")
repository.delete()
```

```http
DELETE /api/repo/v1/repository/registry.local/user1/image1/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

repository, err = dockercloud.GetRepository("registry.local/user1/image1")

if err != nil {
    log.Println(err)
}

repository.Remove()
```

```shell
docker-cloud repository rm registry.local/user1/image1
```

Removes the external repository from Docker Cloud. It doesn't remove the repository from the third party registry where it's stored.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`DELETE /api/repo/v1/[optional_namespace/]repository/`

### Path Parameters

Parameter | Description
--------- | -----------
name | The name of the external repository to remove
