# Stacks

## Stack

> Example

```json
{
  "deployed_datetime": "Mon, 13 Oct 2014 11:01:43 +0000",
  "destroyed_datetime": null,
  "nickname": "deployment stack",
  "name": "dockercloud-app",
  "resource_uri": "/api/app/v1/user_namespace/stack/7fe7ec85-58be-4904-81da-de2219098d7c/",
  "services": [
    "/api/app/v1/user_namespace/service/09cbcf8d-a727-40d9-b420-c8e18b7fa55b/"
  ],
  "state": "Running",
  "synchronized": true,
  "uuid": "09cbcf8d-a727-40d9-b420-c8e18b7fa55b"
}
```

A stack is a logical grouping of closely related services, that may be linked with one another.

This is a [namespaced endpoint](#namespaced-endpoints).

### Attributes

Attribute | Description
--------- | -----------
uuid | A unique identifier for the stack generated automatically on creation
resource_uri | A unique API endpoint that represents the stack
name | A user provided name for the stack.
state | The state of the stack (see table `Stack states` below)
synchronized | Flag indicating if the current stack definition is synchronized with their services.
services | List of service resource URIs belonging to the stack
deployed_datetime | The date and time of the last deployment of the stack (if applicable, `null` otherwise)
destroyed_datetime | The date and time of the `terminate` operation on the stack (if applicable, `null` otherwise)
nickname | A user-friendly name for the stack (`name` by default)


### Stack states

State | Description
----- | -----------
Not Running | The stack has been created and has no deployed services yet. Possible actions in this state: `start`, `terminate`.
Starting | All services for the stack are either starting or already running. No actions allowed in this state.
Running | All services for the service are deployed and running. Possible actions in this state: `redeploy`, `terminate`.
Partly running | One or more services of the stack are deployed and running. Possible actions in this state: `redeploy`, `terminate`.
Stopping | All services for the stack are either stopping or already stopped. No actions allowed in this state.
Stopped | All services for the service are stopped. Possible actions in this state: `start`, `redeploy`, `terminate`.
Redeploying | The stack is redeploying all its services with the updated configuration. No actions allowed in this state.
Terminating | All services for the stack are either being terminated or already terminated. No actions allowed in this state.
Terminated | The stack and all its services have been terminated. No actions allowed in this state.


## List all stacks

```python
import dockercloud

stacks = dockercloud.Stack.list()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stackList, err := dockercloud.ListStacks()

if err != nil {
  log.Println(err)
}

log.Println(stackList)
```

```http
GET /api/app/v1/stack/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack ls
```

Lists all current and recently terminated stacks. Returns a list of `Stack` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]stack/`

### Query Parameters

Parameter | Description
--------- | -----------
uuid | Filter by UUID
name | Filter by stack name


## Create a new stack

```python
import dockercloud

stack = dockercloud.Stack.create(name="my-new-stack", services=[{"name": "hello-word", "image": "tutum/hello-world", "target_num_containers": 2}])
stack.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.CreateStack(dockercloud.StackCreateRequest{Name: "my-new-stack", Services: []dockercloud.ServiceCreateRequest{{Image: "tutum/hello-world", Name: "test", Target_num_containers: 2}}})

if err != nil {
  log.Println(err)
}

log.Println(stack)
```

```http
POST /api/app/v1/stack/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{
    "name": "my-new-stack",
    "services": [
        {
            "name": "hello-word",
            "image": "tutum/hello-world",
            "target_num_containers": 2,
            "linked_to_service": [
                {
                    "to_service": "database",
                    "name": "DB"
                }
            ]
        },
        {
            "name": "database",
            "image": "tutum/mysql"
        }
    ]
}
```

```shell
docker-cloud stack create --name hello-world -f docker-compose.yml
```

Creates a new stack without starting it. The JSON syntax is abstracted to use
[Stack YAML files](/docker-cloud/apps/stack-yaml-reference/) in both
the Docker Cloud CLI and our UI,

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]stack/`

### JSON Parameters

Parameter | Description
--------- | -----------
name | (required) A human-readable name for the stack, such as `my-hello-world-stack`
nickname | (optional) A user-friendly name for the stack (`name` by default)
services | (optional) List of services belonging to the stack. Each service accepts the same parameters as a [Create new service](#create-a-new-service) operation (default: `[]`) plus the ability to refer "links" and "volumes-from" by the name of another service in the stack (see example).


## Export an existing stack

```python
import dockercloud

stack = dockercloud.Stack.fetch("46aca402-2109-4a70-a378-760cfed43816")
stack.export()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

if err = stack.Export(); err != nil {
   log.Println(err)
}
```

```http
GET /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/export/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack export 46aca402
```

Get a JSON representation of the stack following the [Stack YAML representation](/docker-cloud/apps/stacks/).

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]stack/(uuid)/export/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to retrieve



## Get an existing stack

```python
import dockercloud

stack = dockercloud.Stack.fetch("46aca402-2109-4a70-a378-760cfed43816")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

log.Println(stack)
```

```http
GET /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack inspect 46aca402-2109-4a70-a378-760cfed43816
```

Get all the details of an specific stack

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]stack/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to retrieve



## Update an existing stack

```python
import dockercloud

stack = dockercloud.Stack.fetch("46aca402-2109-4a70-a378-760cfed43816")
stack.services = {"services": [{"name": "hello-word", "image": "tutum/hello-world", "target_num_containers": 2}]}
stack.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

if err = stack.Update(dockercloud.StackCreateRequest{Services: []dockercloud.ServiceCreateRequest{{Name: "hello-world", Image: "tutum/hello-world", Target_num_containers: 2}}}); err != nil {
   log.Println(err)
}
```

```http
PATCH /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{
    "services": [
        {
            "name": "hello-word",
            "image": "tutum/hello-world",
            "target_num_containers": 3,
            "linked_to_service": [
                {
                    "to_service": "database",
                    "name": "DB"
                }
            ]
        },
        {
            "name": "database",
            "image": "tutum/mysql"
        }
    ]
}
```

```shell
docker-cloud stack update -f docker-compose.yml 46aca402
```

Updates the details of every service in the stack.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`PATCH /api/app/v1/[optional_namespace/]stack/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to update


### JSON Parameters

Parameter | Description
--------- | -----------
services | (optional) List of services belonging to the stack. Each service accepts the same parameters as a [Update an existing service](#update-an-existing-service) operation (default: `[]`) plus the ability to refer "links" and "volumes-from" by the name of another service in the stack (see example).



## Stop a stack

```python
import dockercloud

stack = dockercloud.Stack.fetch("46aca402-2109-4a70-a378-760cfed43816")
stack.stop()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

if err = stack.Stop(); err != nil {
   log.Println(err)
}
```

```http
POST /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/stop/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack stop 46aca402-2109-4a70-a378-760cfed43816
```

Stops the services in the stack.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]stack/(uuid)/stop/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to stop


## Start a stack

```python
import dockercloud

stack = dockercloud.Stack.fetch()
stack.start()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

if err = stack.Start(); err != nil {
   log.Println(err)
}
```

```http
POST /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/start/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack start 46aca402
```

Starts the services in the stack.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]stack/(uuid)/start/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to start


## Redeploy a stack

```python
import dockercloud

stack = dockercloud.Stack.fetch("46aca402-2109-4a70-a378-760cfed43816")
stack.redeploy()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

//Redeploy(dockercloud.ReuseVolumesOption{Reuse: true}) to reuse the existing volumes
//Redeploy(dockercloud.ReuseVolumesOption{Reuse: false}) to not reuse the existing volumes
if err = stack.Redeploy(dockercloud.ReuseVolumesOption{Reuse: false}); err != nil {
   log.Println(err)
}
```

```http
POST /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/redeploy/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack redeploy 46aca402
```

Redeploys all the services in the stack.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]stack/(uuid)/redeploy/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to redeploy

### Query Parameters

Parameter | Description
--------- | -----------
reuse_volumes | Whether to reuse container volumes for this redeploy operation or not (default: `true`).


## Terminate a stack

```python
import dockercloud

stack = dockercloud.Stack.fetch("46aca402-2109-4a70-a378-760cfed43816")
stack.delete()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

stack, err := dockercloud.GetStack("46aca402-2109-4a70-a378-760cfed43816")

if err != nil {
  log.Println(err)
}

if err = stack.Terminate(); err != nil {
   log.Println(err)
}
```

```http
DELETE /api/app/v1/stack/46aca402-2109-4a70-a378-760cfed43816/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud stack terminate 46aca402
```

Terminate all the services in a the stack and the stack itself.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`DELETE /api/app/v1/[optional_namespace/]stack/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the stack to terminate
