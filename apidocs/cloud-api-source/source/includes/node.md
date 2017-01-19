# Nodes

## Node

> Example

```json
{
    "availability_zone": "/api/infra/v1/az/testing-provider/testing-region/testing-az/",
	"cpu": 1,
	"current_num_containers": 4,
	"deployed_datetime": "Tue, 16 Sep 2014 17:01:15 +0000",
	"destroyed_datetime": null,
	"disk": 60,
	"docker_execdriver": "native-0.2",
	"docker_graphdriver": "aufs",
	"docker_version": "1.5.0",
	"external_fqdn": "fc1a5bb9-user.node.dockerapp.io",
	"last_seen": "Thu, 25 Sep 2014 13:14:44 +0000",
	"memory": 1792,
	"nickname": "fc1a5bb9-user.node.dockerapp.io",
	"node_cluster": "/api/infra/v1/user_namespace/nodecluster/d787a4b7-d525-4061-97a0-f423e8f1d229/",
	"node_type": "/api/infra/v1/user_namespace/nodetype/testing-provider/testing-type/",
	"public_ip": "10.45.2.11",
	"region": "/api/infra/v1/region/testing-provider/testing-region/",
	"resource_uri": "/api/infra/v1/user_namespace/node/fc1a5bb9-17f5-4819-b667-8c7cd819e949/",
	"state": "Deployed",
	"tags": [
		{"name": "tag_one"},
		{"name": "tag-two"}
	],
	"tunnel": "https://tunnel01.cloud.docker.com:12345",
	"uuid": "fc1a5bb9-17f5-4819-b667-8c7cd819e949"
}
```

A node is a virtual machine provided by a cloud provider where containers can be deployed.

This is a [namespaced endpoint](#namespaced-endpoints).

### Attributes

Attribute | Description
--------- | -----------
availability_zone | The resource URI of the availability zone where the node is deployed, if any
uuid | A unique identifier for the node generated automatically on creation
resource_uri | A unique API endpoint that represents the node
external_fqdn | An automatically generated FQDN for the node. Containers deployed on this node will inherit this FQDN.
state | The state of the node. See the below table for a list of possible states.
node_cluster | The resource URI of the node cluster to which this node belongs to (if applicable)
node_type | The resource URI of the node type used for the node
region | The resource URI of the region where the node is deployed
docker_execdriver | Docker's execution driver used in the node
docker_graphdriver | Docker's storage driver used in the node
docker_version | Docker's version used in the node
cpu | Node number of CPUs
disk | Node storage size in GB
memory | Node memory in MB
current_num_containers | The actual number of containers deployed in this node
last_seen | Date and time of the last time the node was contacted by Docker Cloud
public_ip | The public IP allocated to the node
tunnel | If the node does not accept incoming connections to port 2375, the address of the reverse tunnel to access the docker daemon, or `null` otherwise
deployed_datetime | The date and time when this node cluster was deployed
destroyed_datetime | The date and time when this node cluster was terminated (if applicable)
tags | List of tags to identify the node when deploying services (see [Tags](/docker-cloud/apps/deploy-tags/) for more information)
nickname | A user-friendly name for the node (`external_fqdn` by default)


### Node states

State | Description
----- | -----------
Deploying | The node is being deployed in the cloud provider. No actions allowed in this state.
Deployed | The node is deployed and provisioned and is ready to deploy containers. Possible actions in this state: `terminate`, `docker-upgrade`.
Unreachable | The node is deployed but Docker Cloud cannot connect to the docker daemon. Possible actions in this state: `health-check` and `terminate`.
Upgrading | The node docker daemon is being upgraded. No actions allowed in this state.
Terminating | The node is being terminated in the cloud provider. No actions allowed in this state.
Terminated | The node has been terminated and is no longer present in the cloud provider. No actions allowed in this state.


## List all nodes

```python
import dockercloud

nodes = dockercloud.Node.list()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodeList, err := dockercloud.ListNodes()

if err != nil {
  log.Println(err)
}

log.Println(nodeList)
```

```http
GET /api/infra/v1/node/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud node ls
```

Lists all current and recently terminated nodes. Returns a list of `Node` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/[optional_namespace/]node/`

### Query Parameters

Parameter | Description
--------- | -----------
uuid | Filter by UUID
state | Filter by state. Possible values: `Deploying`, `Deployed`, `Unreachable`, `Upgrading`, `Terminating`, `Terminated`
node_cluster | Filter by resource URI of the target node cluster
node_type | Filter by resource URI of the target node type
region | Filter by resource URI of the target region
docker_version | Filter by Docker engine version running in the nodes



## Get an existing node

```python
import dockercloud

node = dockercloud.Node.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

node, err := dockercloud.GetNode("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

log.Println(node)
```

```http
GET /api/infra/v1/node/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud node inspect 7eaf7fff
```

Get all the details of an specific node

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/[optional_namespace/]node/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node to retrieve


## Update a node

```python
import dockercloud

node = dockercloud.Node.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
node.tags.add(["tag-1"])
node.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

node, err := dockercloud.GetNode("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
	log.Println(err)
}

if err = node.Update(dockercloud.Node{Tags: []string{{Name: "tag-1"}}}); err != nil {
			log.Println(err)
}
```

```http
PATCH /api/infra/v1/node/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json

{"tags": [{"name": "tag-1"}], "nickname": "dev node"}
```

```shell
docker-cloud tag add -t tag-1 7eaf7fff
docker-cloud tag set -t tag-2 7eaf7fff
```

Names the node with a user-friendly name and/or replaces the old tags for the new list provided.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`PATCH /api/infra/v1/[optional_namespace/]node/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node to retrieve

### JSON Parameters

Parameter | Description
--------- | -----------
nickname | (optional) A user-friendly name for the node (`external_fqdn` by default)
tags | (optional) List of tags the node will have. This operation replaces the user tag list.


## Upgrade Docker Daemon

```python
import dockercloud

node = dockercloud.Node.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
node.upgrade_docker()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

node, err := dockercloud.GetNode("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = node.Upgrade(); err != nil {
       log.Println(err)
   }
```

```http
POST /api/infra/v1/node/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/docker-upgrade/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud node upgrade 7eaf7fff
```

Upgrades the docker daemon of the node. This will restart your containers on that node. See [Docker upgrade](/docker-cloud/infrastructure/docker-upgrade/) for more information.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/infra/v1/[optional_namespace/]node/(uuid)/docker-upgrade/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node to upgrade


## Perform a health check of a node

```http
POST /api/infra/v1/node/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/health-check/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

Tests connectivity between Docker Cloud and the node. Updates the node status to `Deployed` if the check was successful, or to `Unreachable` otherwise.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/infra/v1/[optional_namespace/]node/(uuid)/health-check/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node to perform the health check to


## Terminate a node

```python
import dockercloud

node = dockercloud.Node.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
node.delete()
```

```http
DELETE /api/infra/v1/node/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

node, err := dockercloud.GetNode("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = node.Terminate(); err != nil {
   log.Println(err)
}
```

```shell
docker-cloud node rm 7eaf7fff
```

Terminates the specified node.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`DELETE /api/infra/v1/[optional_namespace/]node/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node to terminate
