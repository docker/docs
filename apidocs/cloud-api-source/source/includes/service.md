# Services

## Service

> Example

```json
{
  "autodestroy": "OFF",
  "autoredeploy": false,
  "autorestart": "ON_FAILURE",
  "bindings": [
    {
        "host_path": null,
        "container_path": "/tmp",
        "rewritable": true,
        "volumes_from": null
    },
    {
        "host_path": "/etc",
        "container_path": "/etc",
        "rewritable": true,
        "volumes_from": null
    },
    {
        "host_path": null,
        "container_path": null,
        "rewritable": true,
        "volumes_from": "/api/app/v1/user_namespace/service/2f4f54e5-9d3b-4ac1-85ad-a2d4ff25a179/"
    }
  ],
  "cap_add": [
    "ALL"
  ],
  "cap_drop": [
    "NET_ADMIN",
    "SYS_ADMIN"
  ],
  "container_envvars": [
    {
      "key": "DB_PASS",
      "value": "test"
    }
  ],
  "container_ports": [
    {
      "endpoint_uri": "http://wordpress-stackable.admin.srv.dockerapp.io:80/",
      "inner_port": 80,
      "outer_port": 80,
      "port_name": "http",
      "protocol": "tcp",
      "published": true
    }
  ],
  "containers": [
    "/api/app/v1/user_namespace/container/6f8ee454-9dc3-4387-80c3-57aac1be3cc6/",
    "/api/app/v1/user_namespace/container/fdf9c116-7c08-4a60-b0ce-c54ca72c2f25/"
  ],
  "cpu_shares": 100,
  "cpuset": "0,1",
  "cgroup_parent": "m-executor-abcd",
  "current_num_containers": 2,
  "deployed_datetime": "Mon, 13 Oct 2014 11:01:43 +0000",
  "deployment_strategy": "EMPTIEST_NODE",
  "destroyed_datetime": null,
  "devices": [
    "/dev/ttyUSB0:/dev/ttyUSB0"
  ],
  "dns": [
    "8.8.8.8"
  ],
  "dns_search": [
    "example.com"
  ],
  "domainname": "domainname",
  "entrypoint": "",
  "extra_hosts": [
    "onehost:50.31.209.229"
  ],
  "hostname": "hostname",
  "image_name": "tutum/wordpress-stackable:latest",
  "nickname": "wordpress-stackable",
  "labels": {
    "com.example.description": "Accounting webapp",
    "com.example.department": "Finance",
    "com.example.label-with-empty-value": ""
  },
  "link_variables": {
    "WORDPRESS_STACKABLE_1_ENV_DB_HOST": "**LinkMe**",
    "WORDPRESS_STACKABLE_1_ENV_DB_NAME": "wordpress",
    "WORDPRESS_STACKABLE_1_ENV_DB_PASS": "szVaPz925B7I",
    "WORDPRESS_STACKABLE_1_ENV_DB_PORT": "**LinkMe**",
    "WORDPRESS_STACKABLE_1_ENV_DB_USER": "admin",
    "WORDPRESS_STACKABLE_1_ENV_DEBIAN_FRONTEND": "noninteractive",
    "WORDPRESS_STACKABLE_1_ENV_HOME": "/",
    "WORDPRESS_STACKABLE_1_ENV_PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "WORDPRESS_STACKABLE_1_PORT": "tcp://wordpress-stackable-1.admin.cont.dockerapp.io:49153",
    "WORDPRESS_STACKABLE_1_PORT_80_TCP": "tcp://wordpress-stackable-1.admin.cont.dockerapp.io:49153",
    "WORDPRESS_STACKABLE_1_PORT_80_TCP_ADDR": "wordpress-stackable-1.admin.cont.dockerapp.io",
    "WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT": "49153",
    "WORDPRESS_STACKABLE_1_PORT_80_TCP_PROTO": "tcp",
    "WORDPRESS_STACKABLE_2_ENV_DB_HOST": "**LinkMe**",
    "WORDPRESS_STACKABLE_2_ENV_DB_NAME": "wordpress",
    "WORDPRESS_STACKABLE_2_ENV_DB_PASS": "szVaPz925B7I",
    "WORDPRESS_STACKABLE_2_ENV_DB_PORT": "**LinkMe**",
    "WORDPRESS_STACKABLE_2_ENV_DB_USER": "admin",
    "WORDPRESS_STACKABLE_2_ENV_DEBIAN_FRONTEND": "noninteractive",
    "WORDPRESS_STACKABLE_2_ENV_HOME": "/",
    "WORDPRESS_STACKABLE_2_ENV_PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "WORDPRESS_STACKABLE_2_PORT": "tcp://wordpress-stackable-2.admin.cont.dockerapp.io:49154",
    "WORDPRESS_STACKABLE_2_PORT_80_TCP": "tcp://wordpress-stackable-2.admin.cont.dockerapp.io:49154",
    "WORDPRESS_STACKABLE_2_PORT_80_TCP_ADDR": "wordpress-stackable-2.admin.cont.dockerapp.io",
    "WORDPRESS_STACKABLE_2_PORT_80_TCP_PORT": "49154",
    "WORDPRESS_STACKABLE_2_PORT_80_TCP_PROTO": "tcp",
    "WORDPRESS_STACKABLE_ENV_DB_HOST": "**LinkMe**",
    "WORDPRESS_STACKABLE_ENV_DB_NAME": "wordpress",
    "WORDPRESS_STACKABLE_ENV_DB_PASS": "szVaPz925B7I",
    "WORDPRESS_STACKABLE_ENV_DB_PORT": "**LinkMe**",
    "WORDPRESS_STACKABLE_ENV_DB_USER": "admin",
    "WORDPRESS_STACKABLE_ENV_DEBIAN_FRONTEND": "noninteractive",
    "WORDPRESS_STACKABLE_ENV_HOME": "/",
    "WORDPRESS_STACKABLE_ENV_PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "WORDPRESS_STACKABLE_PORT": "tcp://wordpress-stackable-1.admin.cont.dockerapp.io:49153",
    "WORDPRESS_STACKABLE_PORT_80_TCP": "tcp://wordpress-stackable-1.admin.cont.dockerapp.io:49153",
    "WORDPRESS_STACKABLE_PORT_80_TCP_ADDR": "wordpress-stackable-1.admin.cont.dockerapp.io",
    "WORDPRESS_STACKABLE_PORT_80_TCP_PORT": "49153",
    "WORDPRESS_STACKABLE_PORT_80_TCP_PROTO": "tcp",
    "WORDPRESS_STACKABLE_DOCKERCLOUD_API_URL": "https://cloud.docker.com/api/app/v1/user_namespace/service/adeebc1b-1b81-4af0-b8f2-cefffc69d7fb/"
  },
  "linked_from_service": [],
  "linked_to_service": [
    {
      "from_service": "/api/app/v1/user_namespace/service/09cbcf8d-a727-40d9-b420-c8e18b7fa55b/",
      "name": "DB",
      "to_service": "/api/app/v1/user_namespace/service/72f175bd-390b-46e3-9463-830aca32ce3e/"
    }
  ],
  "mac_address": "02:42:ac:11:65:43",
  "memory": 2048,
  "memory_swap": 8192,
  "name": "wordpress-stackable",
  "net": "bridge",
  "privileged": false,
  "public_dns": "wordpress-stackable.admin.svc.dockerapp.io",
  "read_only": true,
  "resource_uri": "/api/app/v1/user_namespace/service/09cbcf8d-a727-40d9-b420-c8e18b7fa55b/",
  "roles": ["global"],
  "run_command": "/run-wordpress.sh",
  "running_num_containers": 1,
  "security_opt": [
  ],
  "sequential_deployment": false,
  "started_datetime": "Mon, 13 Oct 2014 11:01:43 +0000",
  "state": "Partly running",
  "stack": "/api/app/v1/user_namespace/stack/46aca402-2109-4a70-a378-760cfed43816/",
  "stdin_open": false,
  "stopped_datetime": null,
  "stopped_num_containers": 0,
  "synchronized": true,
  "tags": [
        {"name": "tag_one"},
        {"name": "tag-two"},
        {"name": "tagthree3"}
  ],
  "target_num_containers": 2,
  "tty": false,
  "user": "root",
  "uuid": "09cbcf8d-a727-40d9-b420-c8e18b7fa55b",
  "working_dir": "/app"
}
```


A service is a template used to deploy one or more containers.

This is a [namespaced endpoint](#namespaced-endpoints).

### Attributes

Attribute | Description
--------- | -----------
uuid | A unique identifier for the service generated automatically on creation
resource_uri | A unique API endpoint that represents the service
image_name | The Docker image name and tag used for the service containers
name | A user provided name for the service. This name will be inherited by the service containers and will be used in endpoint URLs, environment variable names, etc.
public_dns | An external FQDN that resolves to all IPs of the nodes where the service containers are running on (as an `A` record with multiple IP entries which will be used by clients in a [round-robin fashion](http://en.wikipedia.org/wiki/Round-robin_DNS)). If the service is not publishing any ports, this FQDN will fail to resolve.
state | The state of the service (see table `Service states` below)
net | Network mode to set on the containers (see table `Network Modes` below, more information /docker-cloud/apps/service-links/)
pid | Set the PID (Process) Namespace mode for the containers ([more information](/engine/reference/run/#pid-settings-pid))
synchronized | Flag indicating if the current service definition is synchronized with the current containers.
deployed_datetime | The date and time of the last deployment of the service (if applicable, `null` otherwise)
started_datetime | The date and time of the last `start` operation on the service (if applicable, `null` otherwise)
stopped_datetime | The date and time of the last `stop` operation on the service (if applicable, `null` otherwise)
destroyed_datetime | The date and time of the `terminate` operation on the service (if applicable, `null` otherwise)
target_num_containers | The requested number of containers to deploy for the service
current_num_containers | The actual number of containers deployed for the service
running_num_containers | The actual number of containers deployed for the service in `Running` state
stopped_num_containers | The actual number of containers deployed for the service in `Stopped` state
stack | Resource URIs of the stack that the service belongs to
containers | List of resource URIs of the containers launched as part of the service
container_ports | List of ports to be published on the containers of this service (see table `Service Port attributes` below)
container_envvars | List of user-defined environment variables to set on the containers of the service, which will override the image environment variables (see table `Service Environment Variable attributes` below)
labels | Metadata in form of dictionary used for every container of this service
working_dir | Working directory for running binaries within a container of this service
user | Set the user used on containers of this service (`root` by default)
hostname | Set the hostname used on containers of this service
domainname | Set the domainname used on containers of this service
mac_address | Ethernet device's MAC address used on containers of this service
cgroup_name | Optional parent cgroup used on containers of this service.
tty | If the containers of this service have the tty enable (`false` by default)
stdin_open | If the containers of this service have stdin opened (`false` by default)
dns | Custom DNS servers for containers of this service
dns_search | Custom DNS search domain for containers of this service
cap_add | Added capabilities for containers of this service
cap_drop | Dropped capabilities for containers of this service
devices | List of device mappings for containers of this service
extra_hosts | List of hostname mappings for containers of this service
secuirty_opt | Labeling scheme for containers of this service
entrypoint | Entrypoint to be set on the containers launched as part of the service, which will override the image entrypoint
run_command | Run command to be set on the containers launched as part of the service, which will override the image run command
sequential_deployment | Whether the containers for this service should be deployed in sequence, linking each of them to the previous containers (see [Service scaling](/docker-cloud/apps/service-scaling/) for more information)
cpu_shares | The relative CPU priority of the containers of the service (see [Runtime Constraints on CPU and Memory](/engine/reference/run/#runtime-constraints-on-cpu-and-memory) for more information)
cpuset | CPUs in which to allow execution
memory | The memory limit of the containers of the service in MB (see [Runtime Constraints on CPU and Memory](/engine/reference/run/#runtime-constraints-on-cpu-and-memory) for more information)
memory_swap | Total memory limit (memory + swap) of the containers of the service in MB
linked_from_service | A list of services that are linked to this one (see table `Related services attributes` below)
linked_to_service | A list of services that the service is linked to (see table `Related services attributes` below)
bindings | A list of volume bindings that the service has mounted (see table `Service binding attributes` below)
autorestart | Whether to restart the containers of the service automatically if they stop (see [Crash recovery](/docker-cloud/apps/autorestart/) for more information)
autodestroy | Whether to terminate the containers of the service automatically if they stop (see [Autodestroy](/docker-cloud/apps/auto-destroy/) for more information)
roles | List of Docker Cloud roles assigned to this service (see [Service links](/docker-cloud/apps/service-links/) for more information)
link_variables | List of environment variables that would be exposed in the containers if they are linked to this service
privileged | Whether to start the containers with Docker's `privileged` flag set or not, which allows containers to access all devices on the host among other things (see [Runtime privilege](/engine/reference/run/#runtime-privilege-linux-capabilities-and-lxc-configuration) for more information)
read_only | Whether the filesystem of every service container is read-only or not (`false` by default)
deployment_strategy | Container distribution among nodes (see table `Deployment strategies` below and [Deployment strategies](/docker-cloud/infrastructure/deployment-strategies/) for more information)
tags | List of tags to be used to deploy the service (see [Tags](/docker-cloud/apps/deploy-tags/) for more information)
autoredeploy | Whether to redeploy the containers of the service when its image is updated in Docker Cloud registry (see [Docker Cloud's private registry](/docker-cloud/apps/auto-redeploy/) for more information)
nickname | A user-friendly name for the service (`name` by default)


### Service binding attributes

Attribute | Description
--------- | -----------
host_path | The host path of the volume
container_path | The container path where the volume is mounted
rewritable | `true` is the volume has writable permissions
volumes_from | The resource URI of the service


### Service Port attributes

Attribute | Description
--------- | -----------
protocol | The protocol of the port, either `tcp` or `udp`
inner_port | The published port number inside the container
outer_port | The published port number in the node public network interface
port_name | Name of the service associated to this port
endpoint_uri | The URI of the service endpoint for this port
published | Whether the port has been published in the host public network interface or not. Non-published ports can only be accessed via links.


### Service Environment Variable attributes

Attribute | Description
--------- | -----------
key | The name of the environment variable
value | The value of the environment variable


### Related services attributes

Attribute | Description
--------- | -----------
name | The link name
from_service | The resource URI of the origin of the link
to_service | The resource URI of the target of the link


### Service states

State | Description
----- | -----------
Not running | The service has been created and has no deployed containers yet. Possible actions in this state: `start`, `terminate`.
Starting | All containers for the service are either starting or already running. No actions allowed in this state.
Running | All containers for the service are deployed and running. Possible actions in this state: `stop`, `redeploy`, `terminate`.
Partly running | One or more containers of the service are deployed and running. Possible actions in this state: `stop`, `redeploy`, `terminate`.
Scaling | The service is either deploying new containers or destroying existing ones responding to a scaling request. No actions allowed in this state.
Redeploying | The service is redeploying all its containers with the updated configuration. No actions allowed in this state.
Stopping | All containers for the service are either stopping or already stopped. No actions allowed in this state.
Stopped | All containers for the service are stopped. Possible actions in this state: `start`, `redeploy`, `terminate`.
Terminating | All containers for the service are either being terminated or already terminated. No actions allowed in this state.
Terminated | The service and all its containers have been terminated. No actions allowed in this state.


### Deployment strategies

Strategy | Description
-------- | -----------
EMPTIEST_NODE | It will deploy containers to the node with the lower total amount of running containers (default).
HIGH_AVAILABILITY | It will deploy containers to the node with the lower amount of running containers of the same service.
EVERY_NODE | It will deploy one container on every node. The service can't scale manually. New containers will be deployed to new nodes automatically.


### Network Modes

Strategy | Description
-------- | -----------
bridge | Creates a new network stack for the container on the docker bridge.
host | Uses the host network stack inside the container.

## List all services

```python
import dockercloud

services = dockercloud.Service.list()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

serviceList, err := dockercloud.ListServices()

if err != nil {
  log.Println(err)
}

log.Println(serviceList)
```

```http
GET /api/app/v1/service/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service ps
```

Lists all current and recently terminated services. Returns a list of `Service` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]service/`

### Query Parameters

Parameter | Description
--------- | -----------
uuid | Filter by UUID
state | Filter by state. Possible values: `Not running`, `Starting`, `Running`, `Partly running`, `Scaling`, `Redeploying`, `Stopping`, `Stopped`, `Terminating`, `Terminated`
name | Filter by service name
stack | Filter by resource URI of the target stack.

## Create a new service

```python
import dockercloud

service = dockercloud.Service.create(image="tutum/hello-world", name="my-new-app", target_num_containers=2)
service.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.CreateService(dockercloud.ServiceCreateRequest{Image: "tutum/hello-world",  Name: "my-new-app", Target_num_containers: 2})

if err != nil {
  log.Println(err)
}

log.Println(service)
```

```http
POST /api/app/v1/service/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{"image": "tutum/hello-world", "name": "my-new-app", "target_num_containers": 2}
```

```shell
docker-cloud service create -t 2 --name my-new-app tutum/hello-world
```

Creates a new service without starting it.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]service/`

### JSON Parameters

Parameter | Description
--------- | -----------
image | (required) The image used to deploy this service in docker format, such as `tutum/hello-world`
name | (optional) A human-readable name for the service, such as `my-hello-world-app` (default: `image` without namespace)
target_num_containers | (optional) The number of containers to run for this service initially (default: 1)
run_command | (optional) The command used to start the containers of this service, overriding the value specified in the image, i.e. `/run.sh` (default: `null`)
entrypoint | (optional) The command prefix used to start the containers of this service, overriding the value specified in the image, i.e. `/usr/sbin/sshd` (default: `null`)
container_ports | (optional) An array of objects with port information to be published in the containers for this service, which will be added to the image port information, i.e. `[{"protocol": "tcp", "inner_port": 80, "outer_port": 80}]` (default: `[]`) (See table `Service Port attributes` below)
container_envvars | (optional) An array of objects with environment variables to be added in the service containers on launch (overriding any image-defined environment variables), i.e. `[{"key": "DB_PASSWORD", "value": "mypass"}]` (default: `[]`) (See table `Service Environment Variable attributes` below)
linked_to_service | (optional) An array of service resource URIs to link this service to, including the link name, i.e. `[{"to_service": "/api/app/v1/service/80ff1635-2d56-478d-a97f-9b59c720e513/", "name": "db"}]` (default: `[]`) (See table `Related services attributes` below)
bindings | (optional) An array of bindings this service mounts, i.e. `[{"volumes_from": "/api/app/v1/service/80ff1635-2d56-478d-a97f-9b59c720e513/", "rewritable": true}]` (default: `[]`) (See table `Related bindings attributes` below)
autorestart | (optional) Whether the containers for this service should be restarted if they stop, i.e. `ALWAYS` (default: `OFF`, possible values: `OFF`, `ON_FAILURE`, `ALWAYS`) (see [Crash recovery](/docker-cloud/apps/autorestart/) for more information)
autodestroy | (optional) Whether the containers should be terminated if they stop, i.e. `OFF` (default: `OFF`, possible values: `OFF`, `ON_SUCCESS`, `ALWAYS`) (see [Autodestroy](/docker-cloud/apps/auto-destroy/) for more information)
sequential_deployment | (optional) Whether the containers should be launched and scaled in sequence, i.e. `true` (default: `false`) (see [Service scaling](/docker-cloud/apps/service-scaling/) for more information)
roles | (optional) A list of Docker Cloud API roles to grant the service, i.e. `["global"]` (default: `[]`, possible values: `global`) (see [Service links](/docker-cloud/apps/service-links/) for more information)
privileged | (optional) Whether to start the containers with Docker's `privileged` flag set or not, i.e. `false` (default: `false`) (see [Runtime privilege](/engine/reference/run/#runtime-privilege-linux-capabilities-and-lxc-configuration) for more information)
deployment_strategy | (optional) Container distribution among nodes (default: `EMPTIEST_NODE`, see table `Deployment strategies` above and [Deployment strategies](/docker-cloud/infrastructure/deployment-strategies/) for more information)
tags | (optional) A list of tags to be used to deploy the service (see [Tags](/docker-cloud/apps/deploy-tags/) for more information) (default: `[]`)
autoredeploy | (optional) Whether to redeploy the containers of the service when its image is updated in Docker Cloud registry (default: `false`) (see [Docker Cloud's private registry](/docker-cloud/apps/auto-redeploy/) for more information)
net | (optional) Set the network mode to the containers (default: `bridge`, possible values: `bridge`, `host`)
pid | (optional) Set the PID (Process) Namespace mode for the containers (default: `none` value, possible values: `none`, `host`)
working_dir | (optional) Working directory for running binaries within a container of this service (default: `/`)
nickname | (optional) A user-friendly name for the service (`name` by default)


### Related bindings attributes

Attribute | Description
--------- | -----------
host_path | (optional) The host path of the volume
container_path | (required if `volumes_from` is omitted) The container path where the volume is mounted
rewritable | (optional) `true` is the volume has writable permissions (default: `true`)
volumes_from | (required if `container_path` is omitted) The resource URI of the service


### Service Port attributes

Attribute | Description
--------- | -----------
protocol | (required) The protocol of the port, either `tcp` or `udp`
inner_port | (required) The port number inside the container to be published
outer_port | (optional) The port number in the node public network interface to be published (default: dynamic allocation if `published` is `true`)
published | (optional) Whether to publish the port in the host public network interface or not. Non-published ports can only be accessed via links. (default: `false`)


### Service Environment Variable attributes

Attribute | Description
--------- | -----------
key | (required) The name of the environment variable
value | (required) The value of the environment variable


### Related services attributes

Attribute | Description
--------- | -----------
to_service | (required) The resource URI of the target of the link
name | (optional) The link name


## Get an existing service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

log.Println(service)
```

```http
GET /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service inspect 7eaf7fff
```

Get all the details of an specific service

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]service/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to retrieve


## Get the logs of a service

> Example log line

```json
{
    "type": "log",
    "source": "wordpress-stackable-1",
    "log": "Log line from the container indicated by 'source'",
    "streamType": "stdout",
    "timestamp": 1433779324
}
```

```python
import dockercloud

def log_handler(message):
	print message

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.logs(tail=300, follow=True, log_handler=log_handler)
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

c := make(chan Logs)

go service.Logs(c)
	for {
		s := <-c
		log.Println(s)
	}
```

```http
GET /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/logs/ HTTP/1.1
Host: ws.cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Connection: Upgrade
Upgrade: websocket
```

```shell
docker-cloud service logs 7eaf7fff
```

Get the aggregated logs of all the containers of the service.

### Endpoint Type

Available in Docker Cloud's **STREAM API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]service/(uuid)/logs/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to retrieve logs

### Query Parameters

Parameter | Description
--------- | -----------
tail | Number of lines to show from the end of the logs (default: `300`)
follow | Whether to stream logs or close the connection immediately (default: true)


## Update an existing service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.target_num_containers = 3
service.tags.append({"name":"tag-1"})
service.save()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = service.Update(dockercloud.ServiceCreateRequest{Target_num_containers: 3}); err != nil {
   log.Println(err)
}
```

```http
PATCH /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
Content-Type: application/json

{"autorestart": "ON_FAILURE", "autodestroy": "OFF", "container_envvars": [{"key": "DB_PASSWORD", "value": "mypass"}],
"container_ports": [{"protocol": "tcp", "inner_port": 80, "outer_port": 80}], "cpu_shares": 512,
"entrypoint": "/usr/sbin/sshd", "image": "tutum/hello-world",
"linked_to_service": [{"to_service": "/api/app/v1/service/80ff1635-2d56-478d-a97f-9b59c720e513/", "name": "db"}],
"memory": 2048, "privileged": True, "roles": ["global"], "run_command": "/run.sh", "sequential_deployment": False,
"tags": [{"name": "tag-1"}], "target_num_containers": 3, "autoredeploy": False}

```

```shell
docker-cloud service scale 7eaf7fff 3
docker-cloud tag add -t tag-1 7eaf7fff
docker-cloud tag set -t tag-2 7eaf7fff
```

Updates the service details.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`PATCH /api/app/v1/[optional_namespace/]service/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to update


### JSON Parameters

Parameter | Description
--------- | -----------
autorestart | (optional) Whether the containers for this service should be restarted if they stop, i.e. `ALWAYS` (possible values: `OFF`, `ON_FAILURE`, `ALWAYS`) (see [Crash recovery](/docker-cloud/apps/autorestart/) for more information)
autodestroy | (optional) Whether the containers should be terminated if they stop, i.e. `OFF` (possible values: `OFF`, `ON_SUCCESS`, `ALWAYS`) (see [Autodestroy](/docker-cloud/apps/auto-destroy/) for more information)
container_envvars | (optional) An array of objects with environment variables to be added in the service containers on launch (overriding any image-defined environment variables), i.e. `[{"key": "DB_PASSWORD", "value": "mypass"}]` (See table `Service Environment Variable attributes`)
container_ports | (optional) An array of objects with port information to be published in the containers for this service, which will be added to the image port information, i.e. `[{"protocol": "tcp", "inner_port": 80, "outer_port": 80}]` (See table `Service Port attributes`)
cpu_shares | (optional) The relative CPU priority of the containers the service describes (see [Runtime Constraints on CPU and Memory](/engine/reference/run/#runtime-constraints-on-cpu-and-memory) for more information)
entrypoint | (optional) The command prefix used to start the containers of this service, overriding the value specified in the image, i.e. `/usr/sbin/sshd`
image | (optional) The image used to deploy this service in docker format, i.e. `tutum/hello-world`, `tutum/ubuntu:5.6`. If no tag is indicated, it will be set to `latest` by default
linked_to_service | (optional) An array of service resource URIs to link this service to, including the link name, i.e. `[{"to_service": "/api/app/v1/service/80ff1635-2d56-478d-a97f-9b59c720e513/", "name": "db"}]` (See table `Related services attributes` below)
memory | (optional) The memory limit of the containers of the service in MB (see [Runtime Constraints on CPU and Memory](/engine/reference/run/#runtime-constraints-on-cpu-and-memory) for more information)
privileged | (optional) Whether to start the containers with Docker's `privileged` flag set or not, i.e. `false` (see [Runtime privilege](/engine/reference/run/#runtime-privilege-linux-capabilities-and-lxc-configuration) for more information)
roles | (optional) A list of Docker Cloud API roles to grant the service, i.e. `["global"]` (possible values: `global`) (see [Service links](/docker-cloud/apps/service-links/) for more information)
run_command | (optional) The command used to start the containers of this service, overriding the value specified in the image, i.e. `/run.sh`
sequential_deployment | (optional) Whether the containers should be launched and scaled in sequence, i.e. `true` (see [Service scaling](/docker-cloud/apps/service-scaling/) for more information)
tags | (optional) List of new tags the service will have. This operation replaces the tag list
target_num_containers | (optional) The number of containers to scale this service to
deployment_strategy | (optional) Container distribution among nodes. A service cannot be updated to or from a deployment strategy of `EVERY_NODE`. (See table `Deployment strategies` above and [Deployment strategies](/docker-cloud/infrastructure/deployment-strategies/) for more information)
autoredeploy | Whether to redeploy the containers of the service when its image is updated in Docker Cloud registry (see [Docker Cloud's private registry](/docker-cloud/apps/auto-redeploy/) for more information)
net | (optional) Set the network mode to the containers (default: `bridge`, possible values: `bridge`, `host`)
pid | (optional) Set the PID (Process) Namespace mode for the containers (default: `none` value, possible values: `none`, `host`)
working_dir | (optional) Working directory for running binaries within a container of this service (default: `/`)
nickname | (optional) A user-friendly name for the service (`name` by default)


## Start a service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.start()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = service.Start(); err != nil {
   log.Println(err)
}
```

```http
POST /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/start/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service start 7eaf7fff
```

Starts all containers in a stopped or partly running service.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]service/(uuid)/start/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to start


## Stop a service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.stop()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = service.Stop(); err != nil {
   log.Println(err)
}
```

```http
POST /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/stop/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service stop 7eaf7fff
```

Stops all containers in a running or partly running service.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]service/(uuid)/stop/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to stop


## Scale a service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.target_num_containers = 3
service.save()
service.scale()
```

```http
POST /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/scale/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service scale 7eaf7fff-882c-4f3d-9a8f-a22317ac00ce 3
```

Scales the service to its current `target_num_containers` field.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]service/(uuid)/scale/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to scale


## Redeploy a service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.redeploy()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

//Redeploy(dockercloud.ReuseVolumesOption{Reuse: true}) to reuse the existing volumes
//Redeploy(dockercloud.ReuseVolumesOption{Reuse: false}) to not reuse the existing volumes
if err = service.Redeploy(dockercloud.ReuseVolumesOption{Reuse: false}); err != nil {
   log.Println(err)
}
```

```http
POST /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/redeploy/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service redeploy 7eaf7fff
```

Redeploys all containers in the service with the current service configuration.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]service/(uuid)/redeploy/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to redeploy

### Query Parameters

Parameter | Description
--------- | -----------
reuse_volumes | Whether to reuse container volumes for this redeploy operation or not (default: `true`).


## Terminate a service

```python
import dockercloud

service = dockercloud.Service.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
service.delete()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

if err = service.Terminate(); err != nil {
   log.Println(err)
}
```

```http
DELETE /api/app/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud service terminate 7eaf7fff
```

Terminate all the containers in a service and the service itself. This is not reversible. All the data stored in all containers of the service will be permanently deleted.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`DELETE /api/app/v1/[optional_namespace/]service/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the service to terminate
