# Docker Cloud Events

## Docker Cloud Event

> Example

```json
{
	"type": "action",
	"action": "update",
	"parents": [
		"/api/app/v1/user_namespace/container/0b0e3538-88df-4f07-9aed-3a3cc4175076/"
	],
	"resource_uri": "/api/app/v1/user_namespace/action/49f0efe8-a704-4a10-b02f-f96344fabadd/",
	"state": "Success",
	"uuid": "093ba3bb-08dd-48f0-8f12-4d3b85ef85b3",
	"datetime": "2016-02-01T16:47:28Z"
}
```

Docker Cloud events are generated every time any of the following objects is created or changes state:

* Stack
* Service
* Container
* Node Cluster
* Node
* Action

This is a [namespaced endpoint](#namespaced-endpoints).

### Attributes

| Attribute    | Description                                                                                                                      |
|:-------------|:---------------------------------------------------------------------------------------------------------------------------------|
| type         | Type of object that was created or updated. For possible values, check the [events types](#event-types) table below.             |
| action       | Type of action that was executed on the object. Possible values: `create`, `update` or `delete`                                  |
| parents      | List of resource URIs (REST API) of the parents of the object, according to the "Parent-child hierarchy" table below             |
| resource_uri | Resource URI (REST API) of the object that was created or updated. You can do a `GET` operation on this URL to fetch its details |
| state        | The current state of the object                                                                                                  |
| uuid         | Unique identifier for the event                                                                                                  |
| datetime     | Date and time of the event in ISO 8601 format                                                                                    |


### Event types

| Type        | Description                                                                                    |
|:------------|:-----------------------------------------------------------------------------------------------|
| stack       | Whenever a `Stack` is created or updated                                                       |
| service     | Whenever a `Service` is created or updated                                                     |
| container   | Whenever a `Container` is created or updated                                                   |
| nodecluster | Whenever a `Node Cluster` is created or updated                                                |
| node        | Whenever a `Node` is created or updated                                                        |
| action      | Whenever a `Action` is created or updated                                                      |
| error       | Sent when an error occurs on the websocket connection or as part of the authentication process |


### Parent-child hierarchy

| Object type  | Parent types                            |
|:-------------|:----------------------------------------|
| Stack        | (None)                                  |
| Service      | Stack                                   |
| Container    | Service, Stack, Node, Node Cluster      |
| Node Cluster | (None)                                  |
| Node         | Node Cluster                            |
| Action       | (object to which the action applies to) |


## Listen to new Docker Cloud Events

```python
import dockercloud

def process_event(event):
    print event

events = dockercloud.Events()
events.on_message(process_event)
events.run_forever()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

// Listens for container events only
myFilter := dockercloud.NewStreamFilter(&dockercloud.EventFilter{Type: "container"})

stream := dockercloud.NewStream(myFilter)

if err := stream.Connect(); err == nil {
	go stream.RunForever()
} else {
	log.Print("Connect err: " + err.Error())
}

for {
	select {
	case event := <-stream.MessageChan:
		log.Println(event)
	case err := <-stream.ErrorChan:
		log.Println(err)
	}
}
```

```http
GET /api/audit/v1/events/ HTTP/1.1
Host: ws.cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Connection: Upgrade
Upgrade: websocket
```

```shell
docker-cloud event
```

Listens for new Docker Cloud Events

### Endpoint Type

Available in Docker Cloud's **STREAM API**

### HTTP Request

`GET /api/audit/v1/[optional_namespace/]events/`

### Query Parameters

Parameter | Description
--------- | -----------
type | Filter by type
object | Filter by object resource URI
parent | Filter by object parents
