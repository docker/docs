---
title: Docker Engine API and SDKs
redirect_from:
  - /engine/reference/api/
  - /engine/reference/api/docker_remote_api/
  - /reference/api/
  - /reference/api/docker_remote_api/
---

The Engine API is the API served by Docker Engine. It allows you to control
every aspect of Docker from within your own applications, build tools to manage
and monitor applications running on Docker, and even use it to build apps on
Docker itself.

It is the API the Docker client uses to communicate with the Engine, so
everything the Docker client can do can be done with the API. For example:

* Running and managing containers
* Managing Swarm nodes and services
* Reading logs and metrics
* Creating and managing Swarms
* Pulling and managing images
* Managing networks and volumes

The API can be accessed with any HTTP client, but we also provide
Python and Go [SDKs](sdks.md).

## Open schema

Docker’s remote API uses an **open schema model**. In this model, unknown
properties in incoming messages are ignored. Client applications need to take
this behavior into account to ensure they do not break when talking to newer
Docker daemons.

## Authentication

Authentication configuration is handled on the client. The client sends
the `authConfig` as a `POST` in the `/images/(name)/push` endpoint. The
`authConfig`, set as the `X-Registry-Auth` header, is implemented as a Base64
encoded (JSON) string with the following structure:

```json
{"username": "string",
 "password": "string",
 "email": "string",
 "serveraddress" : "string",
 "auth": ""
}
```

Callers should leave the `auth` field empty. The `serveraddress` is a domain or
IP address without protocol specifier. Throughout this structure, double quotes
are **required**.

## Docker events

The following diagram depicts the container states accessible through the API:

![Container states which can be accessed via the API](/engine/api/images/event_state.png)

Some container-related events are not affected by container state, so they are
not included in this diagram. These events are:

- **export** emitted by `docker export`
- **exec_create** emitted by `docker exec`
- **exec_start** emitted by `docker exec` after **exec_create**
- **detach** emitted when client is detached from container process
- **exec_detach** emitted when client is detached from exec process

Running `docker rmi` (or `docker image remove`) emits an **untag** event when
removing an image name. The `docker rmi` command may also emit **delete** events
when images are deleted by ID directly or by deleting the last tag referring to
the image.

> **Acknowledgment**: This diagram and the accompanying text were used with the
> permission of Matt Good and Gilder Labs. See Matt’s original blog post
> [Docker Events Explained](https://gliderlabs.com/blog/2015/04/14/docker-events-explained/).

## Versioned API

The version of the API you should use depends upon the version of your Docker
daemon. A new version of the API is released when new features are added. The
Docker API is backward-compatible, so you do not need to update code that uses
the API unless you need to take advantage of new features.

To see the highest version of the API your Docker daemon and client support, use
`docker version`:

```bash
$ docker version

Client:
 Version:      17.04.0-ce
 API version:  1.28
 Go version:   go1.7.5
 Git commit:   4845c56
 Built:        Wed Apr  5 06:06:36 2017
 OS/Arch:      darwin/amd64

Server:
 Version:      17.04.0-ce
 API version:  1.28 (minimum version 1.12)
 Go version:   go1.7.5
 Git commit:   4845c56
 Built:        Tue Apr  4 00:37:25 2017
 OS/Arch:      linux/amd64
 Experimental: true
```

## API example

As an example, the `docker run` command can be easily implemented in various
programming languages and by hitting the API directly with `curl`.

> **`curl` compatibility**: cURL 7.5x introduced a backward incompatible change.
> Because of that change, you may need to adapt the `curl` commands listed in
> this reference when communicating with a Docker daemon running on the same
> system as the client.
>
> - cURL 7.50 and up requires a valid URL that includes a hostname and two
>   slashes after the protocol specifier, such as
>   `http://localhost/v1.25/containers/json`. The actual value of the hostname
>   does not matter.
> - cURL 7.40 and below require an "invalid" URL that does **not** include a
>   hostname portion, such as `http:/v1.25/containers/json`. Note the single
>   slash after the protocol specifier.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#python">Python</a></li>
  <li><a data-toggle="tab" data-target="#go">Go</a></li>
  <li><a data-toggle="tab" data-target="#curl">curl</a></li>
</ul>
<div class="tab-content">
  <div id="python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
print client.containers.run("alpine", ["echo", "hello", "world"])
```
</div>
<div id="go" class="tab-pane fade" markdown="1">

```go
package main

import (
	"io"
	"os"

	"github.com/moby/moby/client"
	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
```
</div>
<div id="curl" class="tab-pane fade" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" \
  -d '{"Image": "alpine", "Cmd": ["echo", "hello world"]}' \
  -X POST http:/v1.24/containers/create
{"Id":"1c6594faf5","Warnings":null}

$ curl --unix-socket /var/run/docker.sock -X POST http:/v1.24/containers/1c6594faf5/start

$ curl --unix-socket /var/run/docker.sock -X POST http:/v1.24/containers/1c6594faf5/wait
{"StatusCode":0}

$ curl --unix-socket /var/run/docker.sock "http:/v1.24/containers/1c6594faf5/logs?stdout=1"
hello world
```
  </div>
</div>

To learn more, take a look at the [getting started guide](getting-started.md)
