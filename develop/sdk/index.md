---
title: Develop with Docker Engine SDKs and API
description: Using Docker SDKs and APIs to automate Docker tasks in your language of choice
keywords: developing, api, sdk
redirect_from:
- /engine/api/
- /engine/reference/api/
- /engine/reference/api/docker_remote_api/
- /reference/api/
- /reference/api/docker_remote_api/
- /engine/api/sdks/
---

Docker provides an API for interacting with the Docker daemon (called the Docker
Engine API), as well as SDKs for Go and Python. The SDKs allow you to build and
scale Docker apps and solutions quickly and easily. If Go or Python don't work
for you, you can use the Docker Engine API directly.

The Docker Engine API is a RESTful API accessed by an HTTP client such as `wget` or
`curl`, or the HTTP library which is part of most modern programming languages.

## Install the SDKs

Use the following commands to install the Go or Python SDK. Both SDKs can be
installed and coexist together.

### Go SDK

```bash
go get github.com/docker/docker/client
```

[Read the full Docker Engine Go SDK reference](https://godoc.org/github.com/docker/docker/client).

### Python SDK

- **Recommended**: Run `pip install docker`.

- **If you can't use `pip`**:

  1.  [Download the package directly](https://pypi.python.org/pypi/docker/).
  2.  Extract it and change to the extracted directory,
  3.  Run `python setup.py install`.

[Read the full Docker Engine Python SDK reference](https://docker-py.readthedocs.io/).

## View the API reference

You can
[view the reference for the latest version of the API](/engine/api/latest/)
or [choose a specific version](/engine/api/version-history/).

## Versioned API and SDK

The version of the Docker Engine API you should use depends upon the version of
your Docker daemon and Docker client.

A given version of the Docker Engine SDK supports a specific version of the
Docker Engine API, as well as all earlier versions. If breaking changes occur,
they are documented prominently.

> Daemon and client API mismatches
>
> The Docker daemon and client do not necessarily need to be the same version
> at all times. However, keep the following in mind.
>
> - If the daemon is newer than the client, the client does not know about new
>   features or deprecated API endpoints in the daemon.
>
> - If the client is newer than the daemon, the client can request API
>   endpoints that the daemon does not know about.

A new version of the API is released when new features are added. The Docker API
is backward-compatible, so you do not need to update code that uses the API
unless you need to take advantage of new features.

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

- When using the SDK, use the latest version you can, but at least the version
  that incorporates the API version with the features you need.

- When using `curl` directly, specify the version as the first part of the URL.
  For instance, if the endpoint is `/containers/`, you can use
  `/v1.27/containers/`.

- To force the Docker CLI or the Docker Engine SDKs to use an old version
  version of the API than the version reported by `docker version`, set the
  environment variable `DOCKER_API_VERSION` to the correct version. This works
  on Linux, Windows, or macOS clients.

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

### Docker EE and CE API mismatch

If you use Docker EE in production, we recommend using Docker EE in development
too. If you can't, such as when your developers use Docker for Mac or Docker for
Windows and manually build and push images, then your developers need to configure
their Docker clients to use the same version of the API reported by their Docker
daemon. This prevents the developer from using a feature that is not yet supported
on the daemon where the workload runs in production. You can do this one of two ways:

- Configure the Docker client to connect to an external daemon running Docker EE.
  You can use the `-H` flag on the `docker` command or set the `DOCKER_HOST`
  environment variable. The client uses the daemon's latest supported API version.
- Configure the Docker client to use a specific API by setting the `DOCKER_API_VERSION`
  environment variable to the API version to use, such as `1.30`.

### API version matrix

Docker does not recommend running versions prior to 1.12, which means you
are encouraged to use an API version of 1.24 or higher.

{% include api-version-matrix.md %}

### Choose the SDK or API version to use

Use the following guidelines to choose the SDK or API version to use in your
code:

- If you're starting a new project, use the
  [latest version](/engine/api/latest/), but do specify the version you are
  using. This helps prevent surprises.
- If you need a new feature, update your code to use at least the minimum version
  that supports the feature, and prefer the latest version you can use.
- Otherwise, continue to use the version that your code is already using.

## SDK and API quickstart

As an example, the `docker run` command can be easily implemented using the
Docker API directly, or using the Python or Go SDK.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python">Python</a></li>
  <li><a data-toggle="tab" data-target="#curl">HTTP</a></li>
</ul>
<div class="tab-content">

  <div id="go" class="tab-pane fade in active" markdown="1">

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

    statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
            panic(err)
        }
    case <-statusCh:
    }

    out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
    if err != nil {
        panic(err)
    }

    io.Copy(os.Stdout, out)
}
```

  </div>
  <div id="python" class="tab-pane fade" markdown="1">

```python
import docker
client = docker.from_env()
print client.containers.run("alpine", ["echo", "hello", "world"])
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

For more examples, take a look at the
[getting started guide](examples.md).

## Unofficial libraries

There are a number of community supported libraries available for other
languages. They have not been tested by Docker, so if you run into any issues,
file them with the library maintainers.

| Language              | Library                                                                     |
|:----------------------|:----------------------------------------------------------------------------|
| C                     | [libdocker](https://github.com/danielsuo/libdocker)                         |
| C#                    | [Docker.DotNet](https://github.com/ahmetalpbalkan/Docker.DotNet)            |
| C++                   | [lasote/docker_client](https://github.com/lasote/docker_client)             |
| Dart                  | [bwu_docker](https://github.com/bwu-dart/bwu_docker)                        |
| Erlang                | [erldocker](https://github.com/proger/erldocker)                            |
| Gradle                | [gradle-docker-plugin](https://github.com/gesellix/gradle-docker-plugin)    |
| Groovy                | [docker-client](https://github.com/gesellix/docker-client)                  |
| Haskell               | [docker-hs](https://github.com/denibertovic/docker-hs)                      |
| HTML (Web Components) | [docker-elements](https://github.com/kapalhq/docker-elements)               |
| Java                  | [docker-client](https://github.com/spotify/docker-client)                   |
| Java                  | [docker-java](https://github.com/docker-java/docker-java)                   |
| NodeJS                | [dockerode](https://github.com/apocas/dockerode)                            |
| NodeJS                | [harbor-master](https://github.com/arhea/harbor-master)                     |
| Perl                  | [Eixo::Docker](https://github.com/alambike/eixo-docker)                     |
| PHP                   | [Docker-PHP](https://github.com/docker-php/docker-php)                      |
| Ruby                  | [docker-api](https://github.com/swipely/docker-api)                         |
| Rust                  | [docker-rust](https://github.com/abh1nav/docker-rust)                       |
| Rust                  | [shiplift](https://github.com/softprops/shiplift)                           |
| Scala                 | [tugboat](https://github.com/softprops/tugboat)                             |
| Scala                 | [reactive-docker](https://github.com/almoehi/reactive-docker)               |
| Swift                 | [docker-client-swift](https://github.com/valeriomazzeo/docker-client-swift) |
