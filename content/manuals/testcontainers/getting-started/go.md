---
description: How to install Testcontainers for Go and run your first container
keywords: testcontainers, testcontainers go, testcontainers oss, testcontainers oss go, testcontainers go quickstart
toc_max: 3
title: Go quickstart
linkTitle: Go
aliases:
- /testcontainers/getting-started/go/
weight: 10
---

_Testcontainers for Go_ plays well with the native `go test` framework.

The ideal use case is for integration or end to end tests. It helps you to spin
up and manage the dependencies life cycle via Docker.

## System requirements

### Go version

From the [Go Release Policy](https://go.dev/doc/devel/release#policy):

> Each major Go release is supported until there are two newer major releases. For example, Go 1.5 was supported until the Go 1.7 release, and Go 1.6 was supported until the Go 1.8 release. We fix critical problems, including critical security problems, in supported releases as needed by issuing minor revisions (for example, Go 1.6.1, Go 1.6.2, and so on).

_Testcontainers for Go_ is tested against those two latest Go releases, therefore we recommend using any of them.

## Step 1:Install _Testcontainers for Go_

_Testcontainers for Go_ uses [go mod](https://blog.golang.org/using-go-modules) and you can get it installed via:

```
go get github.com/testcontainers/testcontainers-go
```

## Step 2: Spin up Redis

```go
import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	testcontainers.CleanupContainer(t, redisC)
	require.NoError(t, err)
}
```

The `testcontainers.ContainerRequest` describes how the Docker container will
look.

* `Image` is the Docker image the container starts from.
* `ExposedPorts` lists the ports to be exposed from the container.
* `WaitingFor` is a field you can use to validate when a container is ready. It
  is important to get this set because it helps to know when the container is
  ready to receive any traffic. In this case, we check for the logs we know come
  from Redis, telling us that it is ready to accept requests.

When you use `ExposedPorts` you have to imagine yourself using `docker run -p
<port>`.  When you do so, `dockerd` maps the selected `<port>` from inside the
container to a random one available on your host.

In the previous example, we expose `6379` for `tcp` traffic to the outside. This
allows Redis to be reachable from your code that runs outside the container, but
it also makes parallelization possible because if you add `t.Parallel` to your
tests, and each of them starts a Redis container each of them will be exposed on a
different random port.

`testcontainers.GenericContainer` creates the container. In this example we are
using `Started: true`. It means that the container function will wait for the
container to be up and running. If you set the `Start` value to `false` it won't
start, leaving to you the decision about when to start it.

All the containers must be removed at some point, otherwise they will run until
the host is overloaded. One of the ways we have to clean up is by deferring the
terminated function: `defer testcontainers.TerminateContainer(redisC)` which
automatically handles nil container so is safe to use even in the error case.

> [!TIP]
>
> Look at [features/garbage_collector](/features/garbage_collector/) to know another
> way to clean up resources.

## Step 3: Make your code talk to the container

This is just an example, but usually Go applications that rely on Redis are
using the [redis-go](https://github.com/go-redis/redis) client. This code gets
the endpoint from the container we just started, and it configures the client.

```go
endpoint, err := redisC.Endpoint(ctx, "")
if err != nil {
    t.Error(err)
}

client := redis.NewClient(&redis.Options{
    Addr: endpoint,
})

_ = client
```

We expose only one port, so the `Endpoint` does not need a second argument set.

> [!TIP]
>
> If you expose more than one port you can specify the one you need as a second
> argument.

In this case it returns: `localhost:<mappedportfor-6379>`.

## Step 4: Run the test

You can run the test via `go test ./...`

## Step 5: Want to go deeper with Redis?

You can find a more elaborated Redis example in our examples section. Please check it out [here](https://golang.testcontainers.org/modules/redis/).
