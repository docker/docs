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
programming languages and by hitting the API directly with `curl`:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#python">Python</a></li>
  <li><a data-toggle="tab" data-target="#go">Go</a></li>
  <li><a data-toggle="tab" data-target="#curl">curl</a></li>
</ul>
<div class="tab-content">
  <div id="python" class="tab-pane fade in active">
  {% highlight python %}
  import docker
  client = docker.from_env()
  print client.containers.run("alpine", ["echo", "hello", "world"])
  {% endhighlight %}
  </div>
  <div id="go" class="tab-pane fade">
  {% highlight go %}
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
  {% endhighlight %}
  </div>
  <div id="curl" class="tab-pane fade">
  {% highlight bash %}
  $ curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" \
    -d '{"Image": "alpine", "Cmd": ["echo", "hello world"]}' \
    -X POST http:/v1.24/containers/create
  {"Id":"1c6594faf5","Warnings":null}

  $ curl --unix-socket /var/run/docker.sock -X POST http:/v1.24/containers/1c6594faf5/start

  $ curl --unix-socket /var/run/docker.sock -X POST http:/v1.24/containers/1c6594faf5/wait
  {"StatusCode":0}

  $ curl --unix-socket /var/run/docker.sock "http:/v1.24/containers/1c6594faf5/logs?stdout=1"
  hello world
  {% endhighlight %}
  </div>
</div>

To learn more, take a look at the [getting started guide](getting-started.md)
