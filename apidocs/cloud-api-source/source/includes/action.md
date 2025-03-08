# Actions

## Action

> Example

```json
{
    "action": "Cluster Create",
    "end_date": "Wed, 17 Sep 2014 08:26:22 +0000",
    "ip": "56.78.90.12",
    "is_user_action": true,
    "can_be_canceled": false,
    "location": "New York, USA",
    "method": "POST",
    "object": "/api/infra/v1/user_namespace/cluster/eea638f4-b77a-4183-b241-22dbd7866f22/",
    "path": "/api/infra/v1/user_namespace/cluster/",
    "resource_uri": "/api/audit/v1/action/6246c558-976c-4df6-ba60-eb1a344a17af/",
    "start_date": "Wed, 17 Sep 2014 08:26:22 +0000",
    "state": "Success",
    "user": "user_namespace",
    "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.78.2 (KHTML, like Gecko) Version/7.0.6 Safari/537.78.2",
    "uuid": "6246c558-976c-4df6-ba60-eb1a344a17af"
}
```

An action represents an API call by a user. Details of the API call such as timestamp, origin IP address, and user agent are logged in the action object.

Simple API calls that do not require asynchronous execution will return immediately with the appropriate HTTP error code and an action object will be created either in `Success` or `Failed` states. API calls that do require asynchronous execution will return HTTP code `202 Accepted` immediately and create an action object in `In progress` state, which will change to `Success` or `Failed` state depending on the outcome of the operation being performed. In both cases the response will include a `X-DockerCloud-Action-URI` header with the resource URI of the created action.


### Attributes

| Attribute       | Description                                                                        |
|:----------------|:-----------------------------------------------------------------------------------|
| resource_uri    | A unique API endpoint that represents the action                                   |
| uuid            | A unique identifier for the action generated automatically on creation             |
| object          | The API object (resource URI) to which the action applies to                       |
| action          | Name of the operation performed/being performed                                    |
| method          | HTTP method used to access the API                                                 |
| path            | HTTP path of the API accessed                                                      |
| user            | The user authenticated in the request that created the action                      |
| user_agent      | The user agent provided by the client when accessing the API endpoint              |
| start_date      | Date and time when the API call was performed and the operation started processing |
| end_date        | Date and time when the API call finished processing                                |
| state           | State of the operation (see table below)                                           |
| ip              | IP address of the user that performed the API call                                 |
| location        | Geographic location of the IP address of the user that performed the API call      |
| is_user_action  | If the action has been triggered by the user                                       |
| can_be_canceled | If the action can be canceled by the user in the middle of its execution           |
| can_be_retried  | If the action can be retried by the user                                           |


### Action states

| State       | Description                                                                                  |
|:------------|:---------------------------------------------------------------------------------------------|
| Pending     | The action needed asynchronous execution and it is waiting for an in progress action         |
| In progress | The action needed asynchronous execution and is being performed                              |
| Canceling   | The action is being canceled by user request                                                 |
| Canceled    | The action has been canceled                                                                 |
| Success     | The action was executed successfully                                                         |
| Failed      | There was an issue when the action was being performed. Check the logs for more information. |


## List all actions

```python
import dockercloud

actions = dockercloud.Action.list()
```
```go
import "github.com/docker/go-dockercloud/dockercloud"

actionList, err := dockercloud.ListActions()

if err != nil {
  log.Println(err)
}

log.Println(actionList)
```

```http
GET /api/audit/v1/action/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud action ls
```

Lists all actions in chronological order. Returns a list of `Action` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/audit/v1/action/`

### Query Parameters

| Parameter       | Description                                                                                                                                                                                                                                                                                    |
|:----------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| uuid            | Filter by UUID.                                                                                                                                                                                                                                                                                |
| state           | Filter by state. Possible values: `In progress`, `Success`, `Failed`                                                                                                                                                                                                                           |
| start_date      | Filter by start date. Valid filtering values are `start_date__gte` (after or on the date supplied) and `start_date__lte` (before or on the date supplied)                                                                                                                                      |
| end_date        | Filter by end date. Valid filtering values are `end_date__gte` (after or on the date supplied) and `end_date__lte` (before or on the date supplied)                                                                                                                                            |
| object          | Filter by resource URI of the related object. This filter can only be combined with 'include_related' filter                                                                                                                                                                                   |
| include_related | There is a parent-child relationship between Docker Cloud objects, described in table `Relationships between Docker Cloud objects`. If set to 'true', will include the actions of the related objects to the object specified in "object" filter parameter. Possible values: 'true' or 'false' |


## Relationships between Docker Cloud objects

| Object       | Relationships                                                                  |
|:-------------|:-------------------------------------------------------------------------------|
| Container    | Container, service, stack (if any)                                             |
| Service      | All containers in the service, service, stack (if any)                         |
| Stack        | All services in the stack, all containers in every service in the stack, stack |
| Node         | Node, node cluster (if any)                                                    |
| Node cluster | All nodes in the cluster, node cluster                                         |


## Get an action by UUID

```python
import dockercloud

action = dockercloud.Action.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

action, err := dockercloud.GetAction("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

log.Println(action)
```

```http
GET /api/audit/v1/action/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud action inspect 7eaf7fff
```


Get all the details of an specific action

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/audit/v1/action/(uuid)/`

### Path Parameters

| Parameter | Description                        |
|:----------|:-----------------------------------|
| uuid      | The UUID of the action to retrieve |


## Get the logs of an action

> Example log line

```json
{
    "type": "log",
    "log": "Log line from the action",
    "timestamp": 1433779324
}
```

```python
import dockercloud

def log_handler(message):
    print message

action = dockercloud.Action.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
action.logs(tail=300, follow=True, log_handler=log_handler)
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

c := make(chan dockercloud.Logs)
action, err := dockercloud.GetAction("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
    log.Println(err)
}

go action.GetLogs(c)

for {
	log.Println(<-c)
}
```

```http
GET /api/audit/v1/action/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/logs/ HTTP/1.1
Host: ws.cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Connection: Upgrade
Upgrade: websocket
```

```shell
docker-cloud action logs 7eaf7fff-882c-4f3d-9a8f-a22317ac00ce
```


Get the logs of the specified action.


### Endpoint Type

Available in Docker Cloud's **STREAM API**

### HTTP Request

`GET /api/audit/v1/action/(uuid)/logs/`

### Path Parameters

| Parameter | Description                             |
|:----------|:----------------------------------------|
| uuid      | The UUID of the action to retrieve logs |

### Query Parameters

| Parameter | Description                                                                |
|:----------|:---------------------------------------------------------------------------|
| tail      | Number of lines to show from the end of the logs (default: `300`)          |
| follow    | Whether to stream logs or close the connection immediately (default: true) |

## Cancel an action

```http
POST /api/audit/v1/action/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/cancel/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

action, err := dockercloud.GetAction("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

action, err = action.Cancel()

if err != nil {
  log.Println(err)
}

log.Println(action)
```

Cancels an action in Pending or In progress state.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/audit/v1/action/(uuid)/cancel/`

### Path Parameters

| Parameter | Description                      |
|:----------|:---------------------------------|
| uuid      | The UUID of the action to cancel |


## Retry an action

```python
import dockercloud

def log_handler(message):
  print message

action = dockercloud.Action.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
action.logs(tail=300, follow=True, log_handler=log_handler)
```

```http
POST /api/audit/v1/action/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/retry/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

action, err := dockercloud.GetAction("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

action, err = action.Retry()

if err != nil {
  log.Println(err)
}

log.Println(action)
```

```shell
docker-cloud action logs 7eaf7fff-882c-4f3d-9a8f-a22317ac00ce
```

Retries an action in Success, Failed or Canceled state.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/audit/v1/action/(uuid)/retry/`

### Path Parameters

| Parameter | Description                     |
|:----------|:--------------------------------|
| uuid      | The UUID of the action to retry |
