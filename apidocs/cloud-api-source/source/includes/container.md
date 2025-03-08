# Containers

## Container

> Example

```json
{
    "autodestroy": "OFF",
    "autorestart": "OFF",
    "bindings": [
        {
            "volume": "/api/infra/v1/user_namespace/volume/1863e34d-6a7d-4945-aefc-8f27a4ab1a9e/",
            "host_path": null,
            "container_path": "/data",
            "rewritable": true
        },
        {
            "volume": null,
            "host_path": "/etc",
            "container_path": "/etc",
            "rewritable": true
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
            "key": "DB_1_ENV_DEBIAN_FRONTEND",
            "value": "noninteractive"
        },
        {
            "key": "DB_1_ENV_MYSQL_PASS",
            "value": "**Random**"
        },
        {
            "key": "DB_1_ENV_MYSQL_USER",
            "value": "admin"
        },
        {
            "key": "DB_1_ENV_PATH",
            "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
        },
        {
            "key": "DB_1_ENV_REPLICATION_MASTER",
            "value": "**False**"
        },
        {
            "key": "DB_1_ENV_REPLICATION_PASS",
            "value": "replica"
        },
        {
            "key": "DB_1_ENV_REPLICATION_SLAVE",
            "value": "**False**"
        },
        {
            "key": "DB_1_ENV_REPLICATION_USER",
            "value": "replica"
        },
        {
            "key": "DB_1_PORT",
            "value": "tcp://172.16.0.3:3306"
        },
        {
            "key": "DB_1_PORT_3306_TCP",
            "value": "tcp://172.16.0.3:3306"
        },
        {
            "key": "DB_1_PORT_3306_TCP_ADDR",
            "value": "172.16.0.3"
        },
        {
            "key": "DB_1_PORT_3306_TCP_PORT",
            "value": "3306"
        },
        {
            "key": "DB_1_PORT_3306_TCP_PROTO",
            "value": "tcp"
        },
        {
            "key": "DB_ENV_DEBIAN_FRONTEND",
            "value": "noninteractive"
        },
        {
            "key": "DB_ENV_MYSQL_PASS",
            "value": "**Random**"
        },
        {
            "key": "DB_ENV_MYSQL_USER",
            "value": "admin"
        },
        {
            "key": "DB_ENV_PATH",
            "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
        },
        {
            "key": "DB_ENV_REPLICATION_MASTER",
            "value": "**False**"
        },
        {
            "key": "DB_ENV_REPLICATION_PASS",
            "value": "replica"
        },
        {
            "key": "DB_ENV_REPLICATION_SLAVE",
            "value": "**False**"
        },
        {
            "key": "DB_ENV_REPLICATION_USER",
            "value": "replica"
        },
        {
            "key": "DB_PASS",
            "value": "szVaPz925B7I"
        },
        {
            "key": "DB_PORT",
            "value": "tcp://172.16.0.3:3306"
        },
        {
            "key": "DB_PORT_3306_TCP",
            "value": "tcp://172.16.0.3:3306"
        },
        {
            "key": "DB_PORT_3306_TCP_ADDR",
            "value": "172.16.0.3"
        },
        {
            "key": "DB_PORT_3306_TCP_PORT",
            "value": "3306"
        },
        {
            "key": "DB_PORT_3306_TCP_PROTO",
            "value": "tcp"
        },
        {
            "key": "DB_DOCKERCLOUD_API_URL",
            "value": "https://cloud.docker.com/api/app/v1/user_namespace/service/c0fed1dc-c528-40c9-aa4c-dc00672ebcbf/"
        }
    ],
    "container_ports": [
        {
            "endpoint_uri": "http://wordpress-stackable-1.admin.cont.dockerapp.io:49153/",
            "inner_port": 80,
            "outer_port": 49153,
            "port_name": "http",
            "protocol": "tcp",
            "published": true,
            "uri_protocol": "http"
        }
    ],
    "cpu_shares": 100,
    "cpuset": "0,1",
    "cgroup_parent": "m-executor-abcd",
    "deployed_datetime": "Thu, 16 Oct 2014 12:04:08 +0000",
    "destroyed_datetime": null,
    "devices": [
        "/dev/ttyUSB0:/dev/ttyUSB0"
    ],
    "dns": [
        "8.8.8.8"
    ],
    "dns_search": [
        "example.com",
        "c1dd4e1e-1356-411c-8613-e15146633640.local.dockerapp.io"
    ],
    "domainname": "domainname",
    "entrypoint": "",
    "exit_code": null,
    "exit_code_msg": null,
    "extra_hosts": [
        "onehost:50.31.209.229"
    ],
    "hostname": "hostname",
    "image_name": "tutum/wordpress-stackable:latest",
    "labels": {
        "com.example.description": "Accounting webapp",
        "com.example.department": "Finance",
        "com.example.label-with-empty-value": ""
    },
    "linked_to_container": [
    	{
    		"endpoints": {
    			"3306/tcp": "tcp://172.16.0.3:3306"
    		},
    		"from_container": "/api/app/v1/user_namespace/container/c1dd4e1e-1356-411c-8613-e15146633640/",
    		"name": "DB_1",
    		"to_container": "/api/app/v1/user_namespace/container/ba434e1e-1234-411c-8613-e15146633640/"
    	}
    ],
    "link_variables": {
        "WORDPRESS_STACKABLE_1_ENV_DB_HOST": "**LinkMe**",
        "WORDPRESS_STACKABLE_1_ENV_DB_NAME": "wordpress",
        "WORDPRESS_STACKABLE_1_ENV_DB_PASS": "szVaPz925B7I",
        "WORDPRESS_STACKABLE_1_ENV_DB_PORT": "**LinkMe**",
        "WORDPRESS_STACKABLE_1_ENV_DB_USER": "admin",
        "WORDPRESS_STACKABLE_1_ENV_DEBIAN_FRONTEND": "noninteractive",
        "WORDPRESS_STACKABLE_1_ENV_HOME": "/",
        "WORDPRESS_STACKABLE_1_ENV_PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        "WORDPRESS_STACKABLE_1_PORT": "tcp://172.16.0.2:80",
        "WORDPRESS_STACKABLE_1_PORT_80_TCP": "tcp://172.16.0.2:80",
        "WORDPRESS_STACKABLE_1_PORT_80_TCP_ADDR": "172.16.0.2",
        "WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT": "80",
        "WORDPRESS_STACKABLE_1_PORT_80_TCP_PROTO": "tcp",
        "WORDPRESS_STACKABLE_ENV_DB_HOST": "**LinkMe**",
        "WORDPRESS_STACKABLE_ENV_DB_NAME": "wordpress",
        "WORDPRESS_STACKABLE_ENV_DB_PASS": "szVaPz925B7I",
        "WORDPRESS_STACKABLE_ENV_DB_PORT": "**LinkMe**",
        "WORDPRESS_STACKABLE_ENV_DB_USER": "admin",
        "WORDPRESS_STACKABLE_ENV_DEBIAN_FRONTEND": "noninteractive",
        "WORDPRESS_STACKABLE_ENV_HOME": "/",
        "WORDPRESS_STACKABLE_ENV_PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        "WORDPRESS_STACKABLE_PORT": "tcp://172.16.0.2:80",
        "WORDPRESS_STACKABLE_PORT_80_TCP": "tcp://172.16.0.2:80",
        "WORDPRESS_STACKABLE_PORT_80_TCP_ADDR": "172.16.0.2",
        "WORDPRESS_STACKABLE_PORT_80_TCP_PORT": "80",
        "WORDPRESS_STACKABLE_PORT_80_TCP_PROTO": "tcp"
    },
    "mac_address": "02:42:ac:11:65:43",
    "memory": 1024,
    "memory_swap": 4096,
    "name": "wordpress-stackable",
    "net": "bridge",
    "node": "/api/infra/v1/user_namespace/node/9691c44e-3155-4ca2-958d-c9571aac0a14/",
    "pid": "none",
    "private_ip": "10.7.0.1",
    "privileged": false,
    "public_dns": "wordpress-stackable-1.admin.cont.dockerapp.io",
    "read_only": true,
    "resource_uri": "/api/app/v1/user_namespace/container/c1dd4e1e-1356-411c-8613-e15146633640/",
    "roles": ["global"],
    "run_command": "/run-wordpress.sh",
    "security_opt": [
        "label:user:USER",
        "label:role:ROLE"
    ],
    "service": "/api/app/v1/user_namespace/service/adeebc1b-1b81-4af0-b8f2-cefffc69d7fb/",
    "started_datetime": "Thu, 16 Oct 2014 12:04:08 +0000",
    "state": "Running",
    "stdin_open": false,
    "stopped_datetime": null,
    "synchronized": true,
    "tty": false,
    "user": "root",
    "uuid": "c1dd4e1e-1356-411c-8613-e15146633640",
    "working_dir": "/app"
}
```


A container is a representation of a Docker container in a node.

This is a [namespaced endpoint](#namespaced-endpoints).

### Attributes

Attribute | Description
--------- | -----------
uuid | A unique identifier for the container generated automatically on creation
resource_uri | A unique API endpoint that represents the container
image_name | The Docker image name and tag of the container
bindings | A list of volume bindings that the container has mounted (see table `Container Binding attributes` below)
name | A user provided name for the container (inherited from the service)
node | The resource URI of the node where this container is running
service | The resource URI of the service which this container is part of
public_dns | The external FQDN of the container
state | The state of the container (see table `Container states` below)
synchronized | Flag indicating if the container is synchronized with the current service definition.
exit_code | The numeric exit code of the container (if applicable, `null` otherwise)
exit_code_msg | A string representation of the exit code of the container (if applicable, `null` otherwise)
deployed_datetime | The date and time of the last deployment of the container (if applicable, `null` otherwise)
started_datetime | The date and time of the last `start` operation on the container (if applicable, `null` otherwise)
stopped_datetime | The date and time of the last `stop` operation on the container (if applicable, `null` otherwise)
destroyed_datetime | The date and time of the `terminate` operation on the container (if applicable, `null` otherwise)
container_ports | List of published ports of this container (see table `Container Port attributes` below)
container_envvars | List of user-defined environment variables set on the containers of the service, which will override the container environment variables (see table `Container Environment Variable attributes` below)
labels | Container metadata in form of dictionary
working_dir | Working directory for running binaries within a container
user | User used on the container on launch
hostname | Hostname used on the container on launch
domainname | Domainname used on the container on launch
mac_address | Ethernet device's MAC address used on the container on launch
cgroup_name | Optional parent cgroup for the container.
tty | If the container has the tty enable
stdin_open | If the container has stdin opened
dns | Container custom DNS servers
dns_search | Container custom DNS search domain
cap_add | Container added capabilities
cap_drop | Container dropped capabilities
devices | List of container device mappings
extra_hosts | List of container hostname mappings
secuirty_opt | Labeling scheme of this container
entrypoint | Entrypoint used on the container on launch
run_command | Run command used on the container on launch
cpu_shares | The relative CPU priority of the container (see [Runtime Constraints on CPU and Memory](/engine/reference/run/#runtime-constraints-on-cpu-and-memory) for more information)
cpuset | CPUs in which execution is allowed
memory | The memory limit of the container in MB (see [Runtime Constraints on CPU and Memory](/engine/reference/run/#runtime-constraints-on-cpu-and-memory) for more information)
memory_swap | Total memory limit (memory + swap) of the container in MB
autorestart | Whether to restart the container automatically if it stops (see [Crash recovery](/docker-cloud/apps/autorestart/) for more information)
autodestroy | Whether to terminate the container automatically if it stops (see [Autodestroy](/docker-cloud/apps/auto-destroy/) for more information)
roles | List of Docker Cloud roles asigned to this container (see [API roles](/docker-cloud/apps/api-roles/) for more information))
linked_to_container | List of IP addresses of the linked containers (see table `Container Link attributes` below and [Service links](/docker-cloud/apps/service-links/) for more information)
link_variables | List of environment variables that would be exposed in any container that is linked to this one
privileged | Whether the container has Docker's `privileged` flag set or not (see [Runtime privilege](/engine/reference/run/#runtime-privilege-linux-capabilities-and-lxc-configuration) for more information)
read_only | Whether the container filesystem is read-only or not
private_ip | IP address of the container on the overlay network. This IP will be reachable from any other container.
net | Network mode set on the container (see table `Network Modes` below, [more information](/engine/reference/run/#network-settings))
pid | PID (Process) Namespace mode for the container ([more information](/engine/reference/run/#pid-settings-pid))


### Container Binding attributes

Attribute | Description
--------- | -----------
host_path | The host path of the volume
container_path | The container path where the volume is mounted
rewritable | `true` is the volume has writable permissions
volume | The resource URI of the volume


### Container Port attributes

Attribute | Description
--------- | -----------
protocol | The protocol of the port, either `tcp` or `udp`
inner_port | The published port number inside the container
outer_port | The published port number in the node public network interface
port_name | Name of the service associated to this port
uri_protocol | The protocol to be used in the endpoint for this port, such as `http`
endpoint_uri | The URI of the endpoint for this port
published | Whether the port has been published in the host public network interface or not. Non-published ports can only be accessed via links.


### Container Environment Variable attributes

Attribute | Description
--------- | -----------
key | The name of the environment variable
value | The value of the environment variable


### Container States

State | Description
----- | -----------
Starting | The container is being deployed or started (from Stopped). No actions allowed in this state.
Running | The container is deployed and running. Possible actions in this state: `stop`, `terminate`.
Stopping | The container is being stopped. No actions allowed in this state.
Stopped | The container is stopped. Possible actions in this state: `start`, `terminate`.
Terminating | The container is being deleted. No actions allowed in this state.
Terminated | The container has been deleted. No actions allowed in this state.


### Network Modes

Strategy | Description
-------- | -----------
bridge | Creates a new network stack for the container on the docker bridge.
host | Uses the host network stack inside the container.


### Container Link attributes

Attribute | Description
--------- | -----------
name | The name given to the link
from_container | The resource URI of the "client" container
to_container | The resource URI of the "server" container being linked
endpoints | A dictionary with the endpoints (protocol, IP and port) to be used to reach each of the "server" container exposed ports


## List all containers

```python
import dockercloud

containers = dockercloud.Container.list()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

containerList, err := dockercloud.ListContainers()

if err != nil {
  log.Println(err)
}

log.Println(containerList)
```

```http
GET /api/app/v1/container/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud container ps
```

Lists all current and recently terminated containers. Returns a list of `Container` objects.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]container/`

### Query Parameters

Parameter | Description
--------- | -----------
uuid | Filter by UUID
state | Filter by state. Possible values: `Starting`, `Running`, `Stopping`, `Stopped`, `Terminating`, `Terminated`
name | Filter by container name
service | Filter by resource URI of the target service.
node | Filter by resource URI of the target node.

## Get an existing container

```python
import dockercloud

container = dockercloud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockerckoud.GetContainer("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

log.Println(container)
```


```http
GET /api/app/v1/container/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud container inspect 7eaf7fff
```

Get all the details of an specific container

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]container/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container to retrieve


## Get the logs of a container

> Example log line

```json
{
    "type": "log",
    "log": "Log line from the container",
    "streamType": "stdout",
    "timestamp": 1433779324
}
```

```python
import dockercloud

def log_handler(message):
	print message

container = dockercloud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
container.logs(tail=300, follow=True, log_handler=log_handler)
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockercloud.GetContainer("447ecddc-2890-4ea2-849b-99392e0dd7a6")

if err != nil {
	log.Fatal(err)
}
c := make(chan dockercloud.Logs)

go container.Logs(c)
	for {
		s := <-c
		log.Println(s)
	}
```

```http
GET /api/app/v1/container/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/logs/ HTTP/1.1
Host: ws.cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Connection: Upgrade
Upgrade: websocket
```

```shell
docker-cloud container logs 7eaf7fff
```

Get the logs of the specified container.

### Endpoint Type

Available in Docker Cloud's **STREAM API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]container/(uuid)/logs/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container to retrieve logs

### Query Parameters

Parameter | Description
--------- | -----------
tail | Number of lines to show from the end of the logs (default: `300`)
follow | Whether to stream logs or close the connection immediately (default: true)
service | Filter by service (resource URI)


## Start a container

```python
import dockercloud

container = dockercloud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
container.start()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockercloud.GetContainer("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
	log.Println(err)
}

if err = container.Start(); err != nil {
  log.Println(err)
}
```

```http
POST /api/app/v1/container/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/start/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud container start 7eaf7fff
```

Starts a stopped container.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]container/(uuid)/start/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container to start


## Stop a container

```python
import dockercloud

container = dockerlcoud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
container.stop()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockercloud.GetContainer("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
	log.Println(err)
}

if err = container.Stop(); err != nil {
       log.Println(err)
   }
```

```http
POST /api/app/v1/container/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/stop/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud container stop 7eaf7fff
```

Stops a running container.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]container/(uuid)/stop/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container to stop



## Redeploy a container

```python
import dockercloud

container = dockercloud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
container.redeploy()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockercloud.GetContainer("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
	log.Println(err)
}
//Redeploy(dockercloud.ReuseVolumesOption{Reuse: true) to reuse the existing volumes
//Redeploy(dockercloud.ReuseVolumesOption{Reuse: false}) to not reuse the existing volumes
if err = container.Redeploy(dockercloud.ReuseVolumesOption{Reuse: false}); err != nil {
  log.Println(err)
}
```

```http
POST /api/app/v1/container/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/start/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud container redeploy 7eaf7fff
```

Redeploys a container.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`POST /api/app/v1/[optional_namespace/]container/(uuid)/redeploy/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container to redeploy

### Query Parameters

Parameter | Description
--------- | -----------
reuse_volumes | Whether to reuse container volumes for this redeploy operation or not (default: `true`).


## Terminate a container

```python
import dockercloud

container = dockercloud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
container.delete()
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockercloud.GetContainer("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
	log.Println(err)
}

if err = container.Terminate(); err != nil {
       log.Println(err)
   }
```


```http
DELETE /api/app/v1/container/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce/ HTTP/1.1
Host: cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Accept: application/json
```

```shell
docker-cloud container terminate 7eaf7fff
```

Terminates the specified container. This is not reversible. All data stored in the container will be permanently deleted.

### Endpoint Type

Available in Docker Cloud's **REST API**

### HTTP Request

`DELETE /api/app/v1/[optional_namespace/]container/(uuid)/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container to terminate


## Execute command inside a container

```
import dockercloud

def msg_handler(message):
    print message

container = dockercloud.Container.fetch("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
container.execute("ls", handler=msg_handler)
```

```go
import "github.com/docker/go-dockercloud/dockercloud"

container, err := dockercloud.GetContainer("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")

if err != nil {
  log.Println(err)
}

c := make(chan dockercloud.Exec)

container.Exec("ls", c)

```

```http
GET /api/app/v1/container/(uuid)/exec/ HTTP/1.1
Host: ws.cloud.docker.com
Authorization: Basic dXNlcm5hbWU6YXBpa2V5
Connection: Upgrade
Upgrade: websocket
```


```
docker-cloud exec 7eaf7fff ls
```

Executes a command inside the specified running container, creating a bi-directional stream for the process' standard input and output. This endpoint can be connected to using a bi-directional Secure Web Socket `wss://ws.cloud.docker.com/api/app/v1/container/(uuid)/exec/`

### Endpoint Type

Available in Docker Cloud's **STREAM API**

### HTTP Request

`GET /api/app/v1/[optional_namespace/]container/(uuid)/exec/`

### Path Parameters

Parameter | Description
--------- | -----------
uuid | The UUID of the container where the command will be executed

### Query Parameters

Parameter | Description
--------- | -----------
command | Command to be executed (default: `sh`)
