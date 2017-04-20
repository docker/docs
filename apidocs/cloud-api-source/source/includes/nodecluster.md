# Node Clusters

## Node Cluster

> Example

```json
{
    "current_num_nodes": 1,
    "deployed_datetime": "Tue, 16 Sep 2014 17:01:15 +0000",
    "destroyed_datetime": null,
    "disk": 60,
    "nickname": "my test cluster",
    "name": "TestCluster",
    "node_type": "/api/infra/v1/nodetype/aws/t2.micro/",
    "nodes": [
        "/api/infra/v1/user_namespace/node/75d20367-0948-4f10-8ba4-ffb4d16ed3c6/"
    ],
    "region": "/api/infra/v1/region/aws/us-east-1/",
    "resource_uri": "/api/infra/v1/user_namespace/nodecluster/5516df0b-721e-4470-b350-741ff22e63a0/",
    "state": "Deployed",
    "tags": [
        {"name": "tag_one"},
        {"name": "tag-two"},
        {"name": "tagthree3"}
    ],
    "target_num_nodes": 2,
    "uuid": "5516df0b-721e-4470-b350-741ff22e63a0",
    "provider_options": {
        "vpc": {
            "id": "vpc-aa1c70d4",
            "subnets": ["subnet-aaa7d94f", "subnet-aa15fa64"],
            "security_groups": ["sg-aa1c70d4"]
        },
        "iam": {
            "instance_profile_name": "my_instance_profile"
        }
    }
}
```

A node cluster is a group of nodes that share the same provider, region and/or availability zone, and node type. They are on the same private network.

This is a [namespaced endpoint](#namespaced-endpoints).


### Attributes

Attribute | Description
--------- | -----------
uuid | A unique identifier for the node cluster generated automatically on creation
resource_uri | A unique API endpoint that represents the node cluster
name | A user provided name for the node cluster
state | The state of the node cluster. See the below table for a list of possible states.
node_type | The resource URI of the node type used for the node cluster
disk | The size of the disk where images and containers are stored (in GB)
nodes | A list of resource URIs of the `Node` objects on the node cluster
region | The resource URI of the `Region` object where the node cluster is deployed
target_num_nodes | The desired number of nodes for the node cluster
current_num_nodes | The actual number of nodes in the node cluster. This may differ from `target_num_nodes` if the node cluster is being deployed or scaled
deployed_datetime | The date and time when this node cluster was deployed
destroyed_datetime | The date and time when this node cluster was terminated (if applicable)
tags | List of tags to identify the node cluster nodes when deploying services (see [Tags](/docker-cloud/apps/deploy-tags/) for more information)
provider_options | Provider-specific extra options for the deployment of the node (see `Provider options` table below for more information)
nickname | A user-friendly name for the node cluster (`name` by default)


### Node Cluster states

State | Description
----- | -----------
Init | The node cluster has been created and has no deployed containers yet. Possible actions in this state: `deploy`, `terminate`.
Deploying | All nodes in the cluster are either deployed or being deployed. No actions allowed in this state.
Deployed | All nodes in the cluster are deployed and provisioned. Possible actions in this state: `terminate`.
Partly deployed | One or more nodes of the cluster are deployed and running. Possible actions in this state: `terminate`.
Scaling | The cluster is either deploying new nodes or terminating existing ones responding to a scaling request. No actions allowed in this state.
Terminating | All nodes in the cluster are either being terminated or already terminated. No actions allowed in this state.
Terminated | The node cluster and all its nodes have been terminated. No actions allowed in this state.
Empty cluster | There are no nodes deployed in this cluster. Possible actions in this state: `terminate`.


### Provider options

You can specify the following options when using the Amazon Web Services provider:

* `vpc`: VPC-related options (optional)
    * `id`: AWS VPC identifier of the target VPC where the nodes of the cluster will be deployed (required)
    * `subnets`: a list of target subnet identifiers inside selected VPC. If you specify more than one subnet, Docker Cloud will balance among all of them following a high-availability schema (optional)
    * `security_groups`: the security group that will be applied to every node of the cluster (optional)
* `iam`: IAM-related options (optional)
    * `instance_profile_name`: name of the instance profile (container for instance an IAM role) to attach to every node of the cluster (required)


## List all node clusters

```python
import dockercloud

nodeclusters = dockercloud.NodeCluster.list()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodeclusters, err := dockercloud.ListNodeClusters()

if err != nil {
  log.Println(err)
}

log.Println(nodeclusters)
```

```http
GET /api/infra/v1/nodecluster/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud nodecluster ls
```

Lists all current and recently terminated node clusters. Returns a list of `NodeCluster` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/[optional_namespace/]nodecluster/`

### Query Parameters

Parameter | Description
--------- | -----------
uuid | Filter by UUID
state | Filter by state. Possible values: `Init`, `Deploying`, `Deployed`, `Partly deployed`, `Scaling`, `Terminating`, `Terminated`, `Empty cluster`
name | Filter by node cluster name
region | Filter by resource URI of the target region
node_type | Filter by resource URI of the target node type


## Create a new node cluster

```python
import dockercloud

region = dockercloud.Region.fetch("digitalocean/lon1")
node_type = dockercloud.NodeType.fetch("digitalocean/1gb")
nodecluster = dockercloud.NodeCluster.create(name="my_cluster", node_type=node_type, region=region, disk=60)
nodecluster.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodecluster, err := dockercloud.CreateNodeCluster(dockercloud.NodeCreateRequest{Name: "my_cluster", Region: "/api/infra/v1/region/digitalocean/lon1/", NodeType: "/api/infra/v1/nodetype/digitalocean/1gb/", Target_num_nodes: 2})

if err != nil {
  log.Println(err)
}

log.Println(nodecluster)
```

```http
POST /api/infra/v1/nodecluster/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{"name": "my_cluster", "region": "/api/infra/v1/region/digitalocean/lon1/", "node_type": "/api/infra/v1/nodetype/digitalocean/1gb/", "disk": 60}
```

```shell
docker-cloud nodecluster create my_cluster digitalocean lon1 1gb
```

Creates a new node cluster without deploying it.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/infra/v1/[optional_namespace/]nodecluster/`

### JSON Parameters

Parameter | Description
--------- | -----------
name | (required) A user provided name for the node cluster
node_type | (required) The resource URI of the node type to be used for the node cluster
region | (required) The resource URI of the region where the node cluster is to be deployed
disk | (optional) The size of the volume to create where images and containers will be stored, in GB (default: `60`). Not available for Digital Ocean. To create Softlayer nodes you must select one of the following sizes (in GBs): 10, 20, 25, 30, 40, 50, 75, 100, 125, 150, 175, 200, 250, 300, 350, 400, 500, 750, 1000, 1500 or 2000
nickname | (optional) A user-friendly name for the node cluster (`name` by default)
target_num_nodes | (optional) The desired number of nodes for the node cluster (default: `1`)
tags | (optional) List of tags of the node cluster to be used when deploying services see [Tags](/docker-cloud/apps/deploy-tags/) for more information) (default: `[]`)
provider_options | Provider-specific extra options for the deployment of the node (see table `Provider options` above for more information)


## Get an existing node cluster

```python
import dockercloud

service = dockercloud.NodeCluster.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodecluster, err := dockercloud.GetNodeCluster("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

log.Println(nodecluster)
```

```http
GET /api/infra/v1/nodecluster/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud nodecluster inspect 7eaf7fff
```

Get all the details of an specific node cluster

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/infra/v1/[optional_namespace/]nodecluster/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node cluster to retrieve

## Deploy a node cluster

```python
import dockercloud

nodecluster = dockercloud.NodeCluster.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
nodecluster.deploy()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodecluster, err := dockercloud.GetNodeCluster("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = nodecluster.Deploy(); err != nil {
   log.Println(err)
}
```

```http
POST /api/infra/v1/nodecluster/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/deploy/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

Deploys and provisions a recently created node cluster in the specified region and cloud provider.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/infra/v1/[optional_namespace/]nodecluster/(uuid)/deploy/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node cluster to deploy

## Update an existing node cluster

```python
import dockercloud

nodecluster = dockercloud.NodeCluster.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
nodecluster.target_num_nodes = 3
nodecluster.tags.add("tag-1")
nodecluster.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodecluster, err := dockercloud.GetNodeCluster("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = nodecluster.Update(dockercloud.NodeCreateRequest{Target_num_nodes: 4}); err != nil {
   log.Println(err)
}
```

```http
PATCH /api/infra/v1/nodecluster/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{"target_num_nodes": 3, "tags": [{"name": "tag-1"}]}
```

```shell
docker-cloud nodecluster scale 7eaf7fff 3
docker-cloud tag add -t tag-1 7eaf7fff
docker-cloud tag set -t tag-2 7eaf7fff
```

Updates the node cluster details and applies the changes automatically.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`PATCH /api/infra/v1/[optional_namespace/]nodecluster/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node cluster to update


### JSON Parameters

Parameter | Description
--------- | -----------
target_num_nodes | (optional) The number of nodes to scale this node cluster to
tags | (optional) List of tags the node cluster (and nodes within the node cluster) will have. This operation replaces the user tag list.
## Terminate a node cluster

```python
import dockercloud

nodecluster = dockercloud.NodeCluster.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
nodecluster.delete()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

nodecluster, err := dockercloud.GetNodeCluster("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = nodecluster.Terminate(); err != nil {
   log.Println(err)
}
```

```http
DELETE /api/infra/v1/nodecluster/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
dockercloud nodecluster rm 7eaf7fff
```

Terminates all the nodes in a node cluster and the node cluster itself. This is not reversible.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`DELETE /api/infra/v1/[optional_namespace/]nodecluster/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the node cluster to terminate
