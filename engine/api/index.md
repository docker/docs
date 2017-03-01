---
title: Docker Engine API and SDKs
redirect_from:
  - /engine/reference/api/
  - /engine/reference/api/docker_remote_api/
  - /reference/api/
  - /reference/api/docker_remote_api/
---

The Engine API is the API served by Docker Engine. It allows you to control every aspect of Docker from within your own applications, build tools to manage and monitor applications running on Docker, and even use it to build apps on Docker itself.

It is the API the Docker client uses to communicate with the Engine, so everything the Docker client can do can be done with the API. For example:

* Running and managing containers
* Managing Swarm nodes and services
* Reading logs and metrics
* Creating and managing Swarms
* Pulling and managing images
* Managing networks and volumes

The API can be accessed with any HTTP client, but we also provide [SDKs](sdks.md) in Python and Go to make it easier to use from programming languages.

As an example, the `docker run` command can be easily implemented in various programming languages and by hitting the API directly with `curl`:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-python">
{% highlight python %}
import docker
client = docker.from_env()
print client.containers.run("alpine", ["echo", "hello", "world"])
{% endhighlight %}
</section>
<section class="content" id="tab-go">
{% highlight go %}
package main

import (
	"io"
	"os"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, "alpine", types.ImagePullOptions{})
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
</section>
<section class="content" id="tab-curl">
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
</section>
</div>

To learn more, take a look at the [getting started guide](getting-started.md)
