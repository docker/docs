---
title: Get started with the Docker API
keywords: developing, api, sdk, developers, rest, curl, python, go
redirect_from:
- /engine/api/getting-started/
---

After you
[install Docker](/engine/installation.md), you can try out the Docker API.

Optionally, install an SDK for the language you are using. Official SDKs are
available for Python and Go, and a number of community maintained
libraries for other languages.

[Learn more about installing and using Docker SDKs](sdks.md).

These examples show how to perform the same operation using Python, Go, or
by using `curl` directly.

## Versioned API

The Python and Go examples used here do not specify the API version to use,
because they use features that have been part of Docker for a long time. The
Docker API is fully backward compatible.

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

You can specify the API version to use, in one of the following ways:

- When using `curl` directly, specify the version as the first part of the URL.
  For instance, if the endpoint is `/containers/`, you can use
  `/v1.27/containers/`.

- For the SDKs, or to force the Docker CLI to use a specific version of the API,
  set the environment variable `DOCKER_API_VERSION` to the correct version.
  This works on Linux, Windows, or macOS clients.

  ```bash
  DOCKER_API_VERSION='1.27'
  ```

  While the environment variable is set, that version of the API is used, even
  if the Docker daemon supports a newer version.

- For the SDKs, you can also specify the API version programmatically, as a
  parameter to the `client` object. See the
  [Go constructor](https://github.com/moby/moby/blob/master/client/client.go#L136){: target="_blank" class="_"}
  or the
  [Python SDK documentation for `client`](https://docker-py.readthedocs.io/en/stable/client.html).

## API Examples

### Run a container

This first example shows how to run a container using the Docker API. On the
command line, you would use the `docker run` command, but this is just as easy
to do from your own apps too.

This is the equivalent of typing `docker run alpine echo hello world` at the
command prompt:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-run-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-run-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-run-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-run-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
print client.containers.run("alpine", ["echo", "hello", "world"])
```

</div>

<div id="tab-run-go" class="tab-pane fade" markdown="1">

```go
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

<div id="tab-run-curl" class="tab-pane fade" markdown="1">

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
</div><!-- end tab-content -->

### Run a container in the background

You can also run containers in the background, the equivalent of typing
`docker run -d bfirsh/reticulate-splines`:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-rundetach-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-rundetach-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-rundetach-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-rundetach-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
container = client.containers.run("bfirsh/reticulate-splines", detach=True)
print container.id
```

</div>

<div id="tab-rundetach-go" class="tab-pane fade" markdown="1">

```go
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
```

</div>

<div id="tab-rundetach-curl" class="tab-pane fade" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" \
  -d '{"Image": "bfirsh/reticulate-splines"}' \
  -X POST http:/v1.24/containers/create
{"Id":"1c6594faf5","Warnings":null}

$ curl --unix-socket /var/run/docker.sock -X POST http:/v1.24/containers/1c6594faf5/start
```

</div>
</div><!-- end tab-content -->

### List and manage containers

You can use the API to list containers that are running, just like using
`docker ps`:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-listcontainers-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-listcontainers-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-listcontainers-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-listcontainers-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
for container in client.containers.list():
  print container.id
```

</div>

<div id="tab-listcontainers-go" class="tab-pane fade" markdown="1">

```go
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
```

</div>

<div id="tab-listcontainers-curl" class="tab-pane fade" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock http:/v1.24/containers/json
[{
  "Id":"ae63e8b89a26f01f6b4b2c9a7817c31a1b6196acf560f66586fbc8809ffcd772",
  "Names":["/tender_wing"],
  "Image":"bfirsh/reticulate-splines",
  ...
}]
```

</div>
</div><!-- end tab-content -->

### Stop all running containers

Now that you know what containers exist, you can perform operations on them. For
example, to stop all running containers:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-stopcontainers-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-stopcontainers-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-stopcontainers-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-stopcontainers-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
for container in client.containers.list():
  container.stop()
```

</div>

<div id="tab-stopcontainers-go" class="tab-pane fade in" markdown="1">

```go
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
```

</div>

<div id="tab-stopcontainers-curl" class="tab-pane fade in" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock http:/v1.24/containers/json
[{
  "Id":"ae63e8b89a26f01f6b4b2c9a7817c31a1b6196acf560f66586fbc8809ffcd772",
  "Names":["/tender_wing"],
  "Image":"bfirsh/reticulate-splines",
  ...
}]

$ curl --unix-socket /var/run/docker.sock \
  -X POST http:/v1.24/containers/ae63e8b89a26/stop
```

</div>
</div><!-- end tab-content -->

### Print the logs of a specific container

You can also perform actions on individual containers. This example prints the
logs of a container given its ID:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-containerlogs-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-containerlogs-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-containerlogs-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-containerlogs-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
container = client.containers.get('f1064a8a4c82')
print container.logs()
```

</div>
<div id="tab-containerlogs-go" class="tab-pane fade in" markdown="1">

```go
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
```

</div>

<div id="tab-containerlogs-curl" class="tab-pane fade in" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock "http:/v1.24/containers/ca5f55cdb/logs?stdout=1"
Reticulating spline 1...
Reticulating spline 2...
Reticulating spline 3...
Reticulating spline 4...
Reticulating spline 5...
```
</div>
</div><!-- end tab-content -->

### List all images

List the images on your Engine, similar to `docker images`:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-listimages-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-listimages-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-listimages-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-listimages-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
for image in client.images.list():
  print image.id
```

</div>

<div id="tab-listimages-go" class="tab-pane fade" markdown="1">

```go
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
```

</div>

<div id="tab-listimages-curl" class="tab-pane fade" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock http:/v1.24/images/json
[{
  "Id":"sha256:31d9a31e1dd803470c5a151b8919ef1988ac3efd44281ac59d43ad623f275dcd",
  "ParentId":"sha256:ee4603260daafe1a8c2f3b78fd760922918ab2441cbb2853ed5c439e59c52f96",
  ...
}]
```

</div>
</div><!-- end tab-content -->

### Pull images

Pull images, like `docker pull`:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-pullimages-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-pullimages-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-pullimages-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-pullimages-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
image = client.images.pull("alpine")
print image.id
```

</div>

<div id="tab-pullimages-go" class="tab-pane fade" markdown="1">

```go
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

  defer out.Close()

  io.Copy(os.Stdout, out)
}
```

</div>

<div id="tab-pullimages-curl" class="tab-pane fade" markdown="1">

```bash
$ curl --unix-socket /var/run/docker.sock \
  -X POST "http:/v1.24/images/create?fromImage=alpine"
{"status":"Pulling from library/alpine","id":"3.1"}
{"status":"Pulling fs layer","progressDetail":{},"id":"8f13703509f7"}
{"status":"Downloading","progressDetail":{"current":32768,"total":2244027},"progress":"[\u003e                                                  ] 32.77 kB/2.244 MB","id":"8f13703509f7"}
...
```

</div>
</div> <!-- end tab-content -->


### Pull images with authentication

Pull images, like `docker pull`, with authentication:

> **Note**: Credentials are sent in the clear. Docker's official registries use
> HTTPS. Private registries should also be configured to use HTTPS.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-pullimages-auth-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-pullimages-auth-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-pullimages-auth-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-pullimages-auth-python" class="tab-pane fade in active" markdown="1">

The Python SDK retrieves authentication information from the [credentials
store](/engine/reference/commandline/login/#credentials-store) file and
integrates with [credential
helpers](https://github.com/docker/docker-credential-helpers){: target="_blank"
class="_" }. It is possible to override these credentials, but that is out of
scope for this Getting Started guide. After using `docker login`, the Python SDK
uses these credentials automatically.

```python
import docker
client = docker.from_env()
image = client.images.pull("alpine")
print image.id
```

</div>

<div id="tab-pullimages-auth-go" class="tab-pane fade" markdown="1">

```go
package main

import (
  "io"
  "os"
  "encoding/json"
  "encoding/base64"

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

  authConfig := types.AuthConfig{
      Username: "username",
      Password: "password",
  }
  encodedJSON, err := json.Marshal(authConfig)
  if err != nil {
      panic(err)
  }
  authStr := base64.URLEncoding.EncodeToString(encodedJSON)

  out, err := cli.ImagePull(ctx, "alpine", types.ImagePullOptions{RegistryAuth:
  authStr})
  if err != nil {
    panic(err)
  }

  defer out.Close()
  io.Copy(os.Stdout, out)
}
```

</div>

<div id="tab-pullimages-curl" class="tab-pane fade" markdown="1">

This example will leave the credentials in your shell's history, so consider
this a naive implementation. The credentials are passed as a Base-64-encoded
JSON structure.

```bash
$ JSON=$(echo '{"username": "string", "password": "string", "serveraddress": "string"}' | base64)

$ curl --unix-socket /var/run/docker.sock \
  -H "Content-Type: application/tar" 
  -X POST "http:/v1.24/images/create?fromImage=alpine"
  -H "X-Registry-Auth"
  -d "$JSON"
{"status":"Pulling from library/alpine","id":"3.1"}
{"status":"Pulling fs layer","progressDetail":{},"id":"8f13703509f7"}
{"status":"Downloading","progressDetail":{"current":32768,"total":2244027},"progress":"[\u003e                                                  ] 32.77 kB/2.244 MB","id":"8f13703509f7"}
...
```

</div>
</div> <!-- end tab-content -->
### Commit containers

Commit containers to create images from their contents:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab-commit-python" data-group="python">Python</a></li>
  <li><a data-toggle="tab" data-target="#tab-commit-go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#tab-commit-curl" data-group="curl">Curl</a></li>
</ul>

<div class="tab-content">
<div id="tab-commit-python" class="tab-pane fade in active" markdown="1">

```python
import docker
client = docker.from_env()
container = client.containers.run("alpine", ["touch", "/helloworld"], detach=True)
container.wait()
image = container.commit("helloworld")
print image.id
```

</div>

<div id="tab-commit-go" class="tab-pane fade" markdown="1">

```go
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
```

</div>

<div id="tab-commit-curl" class="tab-pane fade" markdown="1">

```bash
$ docker run -d alpine touch /helloworld
0888269a9d584f0fa8fc96b3c0d8d57969ceea3a64acf47cd34eebb4744dbc52
$ curl --unix-socket /var/run/docker.sock\
  -X POST "http:/v1.24/commit?container=0888269a9d&repo=helloworld"
{"Id":"sha256:6c86a5cd4b87f2771648ce619e319f3e508394b5bfc2cdbd2d60f59d52acda6c"}
```

</div>
</div><!-- end tab-content -->

## Next steps

 - [Full documentation for the Python SDK](https://docker-py.readthedocs.io)
 - [Full documentation for the Go SDK](https://godoc.org/github.com/moby/moby/client)
 - [Full documentation for the HTTP API](/engine/api/v1.27/)
