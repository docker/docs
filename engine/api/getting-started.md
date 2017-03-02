---
title: Getting started with the Engine API
---

To try out the Docker Engine API in development, [you first need to install Docker](/engine/installation.md).

Next, you need to install an SDK for the language you are using. There are official ones available for Python and Go, and a number of community maintained libraries for other languages. [Head to the SDKs page to find and install them.](sdks.md)

## Running a container

The most basic thing you can do with Docker is running a container. On the command line, you would use the `docker run` command, but this is just as easy to do from your own apps too.

This is the equivalent of doing `docker run alpine echo hello world`:

<dl class="horizontal tabs" data-tab >
  <dd class="active"><a href="#tab-run-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-run-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-run-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-run-python">
{% highlight python %}
import docker
client = docker.from_env()
print client.containers.run("alpine", ["echo", "hello", "world"])
{% endhighlight %}
</section>
<section class="content" id="tab-run-go">
{% highlight go %}
package main

import (
	"io"
	"os"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
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
<section class="content" id="tab-run-curl">
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

You can also run containers in the background, the equivalent of `docker run -d bfirsh/reticulate-splines`:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-rundetach-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-rundetach-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-rundetach-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-rundetach-python">
{% highlight python %}
import docker
client = docker.from_env()
container = client.containers.run("bfirsh/reticulate-splines", detach=True)
print container.id
{% endhighlight %}
</section>
<section class="content" id="tab-rundetach-go">
{% highlight go %}
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}
{% endhighlight %}
</section>
<section class="content" id="tab-rundetach-curl">
{% highlight bash %}
$ curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" \
  -d '{"Image": "bfirsh/reticulate-splines"}' \
  -X POST http:/v1.24/containers/create
{"Id":"1c6594faf5","Warnings":null}

$ curl --unix-socket /var/run/docker.sock -X POST http:/v1.24/containers/1c6594faf5/start
{% endhighlight %}
</section>
</div>

## Listing and managing containers

Like `docker ps`, we can use the API to list containers that are running:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-listcontainers-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-listcontainers-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-listcontainers-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-listcontainers-python">
{% highlight python %}
import docker
client = docker.from_env()
for container in client.containers.list():
  print container.id
{% endhighlight %}
</section>
<section class="content" id="tab-listcontainers-go">
{% highlight go %}
package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}
{% endhighlight %}
</section>
<section class="content" id="tab-listcontainers-curl">
{% highlight bash %}
$ curl --unix-socket /var/run/docker.sock http:/v1.24/containers/json
[{
  "Id":"ae63e8b89a26f01f6b4b2c9a7817c31a1b6196acf560f66586fbc8809ffcd772",
  "Names":["/tender_wing"],
  "Image":"bfirsh/reticulate-splines",
  ...
}]

{% endhighlight %}
</section>
</div>

Now we know what containers exist, we can perform operations on them. For example, we can stop all running containers:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-stopcontainers-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-stopcontainers-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-stopcontainers-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-stopcontainers-python">
{% highlight python %}
import docker
client = docker.from_env()
for container in client.containers.list():
  container.stop()
{% endhighlight %}
</section>
<section class="content" id="tab-stopcontainers-go">
{% highlight go %}
package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
	}
}
{% endhighlight %}
</section>
<section class="content" id="tab-stopcontainers-curl">
{% highlight bash %}
$ curl --unix-socket /var/run/docker.sock http:/v1.24/containers/json
[{
  "Id":"ae63e8b89a26f01f6b4b2c9a7817c31a1b6196acf560f66586fbc8809ffcd772",
  "Names":["/tender_wing"],
  "Image":"bfirsh/reticulate-splines",
  ...
}]

$ curl --unix-socket /var/run/docker.sock \
  -X POST http:/v1.24/containers/ae63e8b89a26/stop

{% endhighlight %}
</section>
</div>

We can also perform actions on individual containers. For example, to print the logs of a container given its ID:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-containerlogs-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-containerlogs-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-containerlogs-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-containerlogs-python">
{% highlight python %}
import docker
client = docker.from_env()
container = client.containers.get('f1064a8a4c82')
print container.logs()
{% endhighlight %}
</section>
<section class="content" id="tab-containerlogs-go">
{% highlight go %}
package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}
	out, err := cli.ContainerLogs(ctx, "f1064a8a4c82", options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
{% endhighlight %}
</section>
<section class="content" id="tab-containerlogs-curl">
{% highlight bash %}
$ curl --unix-socket /var/run/docker.sock "http:/v1.24/containers/ca5f55cdb/logs?stdout=1"
Reticulating spline 1...
Reticulating spline 2...
Reticulating spline 3...
Reticulating spline 4...
Reticulating spline 5...
{% endhighlight %}
</section>
</div>

## Managing images

Images are the basis of containers, and can be managed in a similar way. You can list the images on your Engine, similar to `docker images`:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-listimages-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-listimages-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-listimages-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-listimages-python">
{% highlight python %}
import docker
client = docker.from_env()
for image in client.images.list():
  print image.id
{% endhighlight %}
</section>
<section class="content" id="tab-listimages-go">
{% highlight go %}
package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID)
	}
}
{% endhighlight %}
</section>
<section class="content" id="tab-listimages-curl">
{% highlight bash %}
$ curl --unix-socket /var/run/docker.sock http:/v1.24/images/json
[{
  "Id":"sha256:31d9a31e1dd803470c5a151b8919ef1988ac3efd44281ac59d43ad623f275dcd",
  "ParentId":"sha256:ee4603260daafe1a8c2f3b78fd760922918ab2441cbb2853ed5c439e59c52f96",
  ...
}]
{% endhighlight %}
</section>
</div>

You can pull images, like `docker pull`:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-pullimages-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-pullimages-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-pullimages-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-pullimages-python">
{% highlight python %}
import docker
client = docker.from_env()
image = client.images.pull("alpine")
print image.id
{% endhighlight %}
</section>
<section class="content" id="tab-pullimages-go">
{% highlight go %}
package main

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, "alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
{% endhighlight %}
</section>
<section class="content" id="tab-pullimages-curl">
{% highlight bash %}
$ curl --unix-socket /var/run/docker.sock \
  -X POST "http:/v1.24/images/create?fromImage=alpine"
{"status":"Pulling from library/alpine","id":"3.1"}
{"status":"Pulling fs layer","progressDetail":{},"id":"8f13703509f7"}
{"status":"Downloading","progressDetail":{"current":32768,"total":2244027},"progress":"[\u003e                                                  ] 32.77 kB/2.244 MB","id":"8f13703509f7"}
...
{% endhighlight %}
</section>
</div>

And commit containers to create images from their contents:

<dl class="horizontal tabs" data-tab>
  <dd class="active"><a href="#tab-commit-python" class="noanchor">Python</a></dd>
  <dd><a href="#tab-commit-go" class="noanchor">Go</a></dd>
  <dd><a href="#tab-commit-curl" class="noanchor">curl</a></dd>
</dl>
<div class="tabs-content">
<section class="content active" id="tab-commit-python">
{% highlight python %}
import docker
client = docker.from_env()
container = client.run("alpine", ["touch", "/helloworld"], detached=True)
container.wait()
image = container.commit("helloworld")
print image.id
{% endhighlight %}
</section>
<section class="content" id="tab-commit-go">
{% highlight go %}
package main

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	createResp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"touch", "/helloworld"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	if _, err = cli.ContainerWait(ctx, createResp.ID); err != nil {
		panic(err)
	}

	commitResp, err := cli.ContainerCommit(ctx, createResp.ID, types.ContainerCommitOptions{Reference: "helloworld"})
	if err != nil {
		panic(err)
	}

	fmt.Println(commitResp.ID)
}
{% endhighlight %}
</section>
<section class="content" id="tab-commit-curl">
{% highlight bash %}
$ docker run -d alpine touch /helloworld
0888269a9d584f0fa8fc96b3c0d8d57969ceea3a64acf47cd34eebb4744dbc52
$ curl --unix-socket /var/run/docker.sock\
  -X POST "http:/v1.24/commit?container=0888269a9d&repo=helloworld"
{"Id":"sha256:6c86a5cd4b87f2771648ce619e319f3e508394b5bfc2cdbd2d60f59d52acda6c"}
{% endhighlight %}
</section>
</div>

## Next steps

 - [Full documentation for the Python SDK.](https://docker-py.readthedocs.io)
 - [Full documentation for the Go SDK.](https://godoc.org/github.com/docker/docker/client)
 - [Full documentation for the HTTP API.](/engine/api/v1.26/)
