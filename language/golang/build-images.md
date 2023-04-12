---
title: "Build your Go image"
keywords: containers, images, go, golang, dockerfiles, coding, build, push, run
description: Learn how to build your first Docker image by writing a Dockerfile
redirect_from:
- /get-started/golang/build-images/
---

{% include_relative nav.html selected="1" %}

## Overview

In this section we are going to build a container image. The image includes everything you need 
to run your application – the compiled application binary file, the runtime, the libraries, and
all other resources required by your application.

## Required software

To complete this tutorial, you need the following:

- Go version 1.19 or later. Visit the [download page for Go](https://golang.org/dl/){:target="_blank" rel="noopener" class="_"} first and install the toolchain.
- Docker running locally. Follow the [instructions to download and install Docker](../../desktop/index.md).
- An IDE or a text editor to edit files. [Visual Studio Code](https://code.visualstudio.com/){: target="_blank" rel="noopener" class="_"} is a free and popular choice but you can use anything you feel comfortable with.
- A Git client. We'll use a command-line based `git` client throughout this module, but you are free to use whatever works for you.
- A command-line terminal application. The examples shown in this module are from the Linux shell, but they should work in PowerShell, Windows Command Prompt, or OS X Terminal with minimal, if any, modifications.

## Meet the example application

The example application is a *caricature* of a microservice. It is purposefully trivial to keep focus on learning the basics of containerization for Go applications.

The application offers two HTTP endpoints:

* It responds with a string containing a heart symbol (`<3`) to requests to `/`.
* It responds with `{"Status" : "OK"}` JSON to a request to `/health`.

It responds with HTTP error 404 to any other request.

The application listens on a TCP port defined by the value of environment variable `PORT`. The default value is `8080`.

The application is *stateless*.

The complete source code for the application is on GitHub: [github.com/olliefr/docker-gs-ping](https://github.com/olliefr/docker-gs-ping){: target="_blank" rel="noopener" class="_"}. You are encouraged to fork it and experiment with it as much as you like.

To continue, we clone the application repository to our local machine:

```console
$ git clone https://github.com/olliefr/docker-gs-ping
```

The application's `main.go` file is fairly straightforward, if you are familiar with Go:

{% raw %}
```go
package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Simple implementation of an integer minimum
// Adapted from: https://gobyexample.com/testing-and-benchmarking
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```
{% endraw %}

## Smoke test the application

Let’s start our application and make sure it’s running properly. Open your
terminal and navigate to the directory into which you cloned the project's repo.
From now on, we'll refer to this directory as the **project directory**.

```console
$ go run main.go
```

This should compile and start the server as a foreground application, outputting
the banner, as illustrated in the next figure.

```
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.10.2
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

Let's run a quick _smoke test_ by accessing the application on `http://localhost:8080`. 
You can use your favourite web browser, or even a `curl` command in the terminal:

```console
$ curl http://localhost:8080/
Hello, Docker! <3
```

This verifies that the application builds locally and we can start it without an error. 
That's a milestone to celebrate!

Now we are ready to "containerize" it.

## Create a Dockerfile for the application

To build a container image with Docker, a *Dockerfile* with build instructions is required.

We begin our `Dockerfile` with the (optional) parser directive line that instructs BuildKit to 
interpret our file according to the grammar rules for the specified version of the syntax.

We then tell Docker what *base image* we would like to use for our application:

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.19
```

Docker images can be inherited from other images. Therefore, instead of creating
our own base image from scratch, we can use the official Go image that already has 
all necessary tools and libraries to compile and run a Go application.

> **Note**
>
> If you are curious about creating your own base images, you can check out the following section of this guide: [creating base images](../../build/building/base-images.md).
> Note, however, that this is not necessary to continue with our task at hand.

Now that we have defined the "base" image for our upcoming container image, 
we can begin building on top of it.

To make things easier when running the rest of our commands, let’s create a
directory _inside_ the image that we are building. This also instructs Docker
to use this directory as the default _destination_ for all subsequent commands.
This way we do not have to type out full file paths in the `Dockerfile`, 
the relative paths will be based on this directory.

```dockerfile
WORKDIR /app
```

Usually the very first thing you do once you’ve downloaded a project written in
Go is to install the modules necessary to compile it. Note, that the base image 
has the toolchain already, but our source code is not in it yet.

So before we can run `go mod download` inside our image, we need to get our
`go.mod` and `go.sum` files copied into it. We use the `COPY` command to do this. 

In its simplest form, the `COPY` command takes two parameters. The first
parameter tells Docker what files you want to copy into the image. The last
parameter tells Docker where you want that file to be copied to. 

We’ll copy the `go.mod` and `go.sum` file into our project directory `/app` which,
owing to our use of `WORKDIR`, is the current directory (`./`) inside the image.
Unlike some modern shells that appear to be indifferent to the use of trailing slash (`/`),
and can figure out what the user meant (most of the time), Docker's `COPY` command 
is quite sensitive in its interpretation of the trailing slash.

```dockerfile
COPY go.mod go.sum ./
```

> **Notice**
>
> Please take some time to familiarise yourself with the trailing slash treatment
> by the `COPY` command: [Dockerfile reference](../../engine/reference/builder.md/#copy)
> as it might otherwise trick you up in more ways than you can imagine.

Now that we have the module files inside the Docker image that we are building,
we can use the `RUN` command to execute the command `go mod download` there as
well. This works exactly the same as if we were running `go` locally on our
machine, but this time these Go modules will be installed into a directory
inside the image.

```dockerfile
RUN go mod download
```

At this point, we have a Go toolchain version 1.19.x and all our Go dependencies
installed inside the image.

The next thing we need to do is to copy our source code into the image. We’ll
use the `COPY` command just like we did with our module files before.

```dockerfile
COPY *.go ./
```

This `COPY` command uses a wildcard to copy all files with `.go` extension
located in the current directory on the host (the directory where the `Dockerfile`
is located) into the current directory inside the image. 

Now, we would like to compile our application. To that end, we use the familiar
`RUN` command:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping
```

This should be familiar. The result of that command will be a static application
binary named `docker-gs-ping` and located in the root of the filesystem of the
image that we are building. We could have put the binary into any other place we
desire inside that image, the root directory has no special meaning in this
regard. It's just convenient to use it to keep the file paths short for improved
readability.

Now, all that is left to do is to tell Docker what command to execute when our
image is used to start a container. 

We do this with the `CMD` command:

```dockerfile
CMD ["/docker-gs-ping"]
```

Here's the complete `Dockerfile`:

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]
```

The `Dockerfile` may also contain _comments_. They always begin with a `#` symbol,
and must be at the beginning of a line. Comments are there for your convenience
to allow documenting your `Dockerfile`. 

There is also a concept of Dockerfile _directives_, such as the `syntax` directive we added.
The directives must always be at the very top of the `Dockerfile`, so when adding comments, 
make sure that the comments follow *after* any directives that you may have used:

```dockerfile
# syntax=docker/dockerfile:1
# A sample microservice in Go packaged into a container image.

FROM golang:1.19

# ...
```

## Build the image

Now that we've created our `Dockerfile`, let’s build an image from it. The
`docker build` command creates Docker images from the `Dockerfile` and a "context".
A build _context_ is the set of files located in the specified path or URL. The
Docker build process can access any of the files located in the context.

The build command optionally takes a `--tag` flag. This flag is used to label
the image with a string value, which is easy for humans to read and recognise.
If you do not pass a `--tag`, Docker will use `latest` as the default value.

Let's build our first Docker image!

```console
$ docker build --tag docker-gs-ping .
```

The build process will print some diagnostic messages as it goes through the build steps. 
The following is just an example of what these messages may look like.

```console
[+] Building 2.2s (15/15) FINISHED
 => [internal] load build definition from Dockerfile                                                                                       0.0s
 => => transferring dockerfile: 701B                                                                                                       0.0s
 => [internal] load .dockerignore                                                                                                          0.0s
 => => transferring context: 2B                                                                                                            0.0s
 => resolve image config for docker.io/docker/dockerfile:1                                                                                 1.1s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:39b85bbfa7536a5feceb7372a0817649ecb2724562a38360f4d6a7782a409b14            0.0s
 => [internal] load build definition from Dockerfile                                                                                       0.0s
 => [internal] load .dockerignore                                                                                                          0.0s
 => [internal] load metadata for docker.io/library/golang:1.19                                                                             0.7s
 => [1/6] FROM docker.io/library/golang:1.19@sha256:5d947843dde82ba1df5ac1b2ebb70b203d106f0423bf5183df3dc96f6bc5a705                       0.0s
 => [internal] load build context                                                                                                          0.0s
 => => transferring context: 6.08kB                                                                                                        0.0s
 => CACHED [2/6] WORKDIR /app                                                                                                              0.0s
 => CACHED [3/6] COPY go.mod go.sum ./                                                                                                     0.0s
 => CACHED [4/6] RUN go mod download                                                                                                       0.0s
 => CACHED [5/6] COPY *.go ./                                                                                                              0.0s
 => CACHED [6/6] RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping                                                                  0.0s
 => exporting to image                                                                                                                     0.0s
 => => exporting layers                                                                                                                    0.0s
 => => writing image sha256:ede8ff889a0d9bc33f7a8da0673763c887a258eb53837dd52445cdca7b7df7e3                                               0.0s
 => => naming to docker.io/library/docker-gs-ping                                                                                          0.0s
```

Your exact output will vary, but provided there aren't any errors, you should
see the word `FINISHED` in the first line of output. This means Docker has successfully
built our image named `docker-gs-ping`.

## View local images

To see the list of images we have on our local machine, we have two options. One
is to use the CLI and the other is to use [Docker Desktop](../../desktop/index.md).
Since we are currently working in the terminal, let’s take a look at listing
images with the CLI.

To list images, run the `docker image ls`command (or the `docker images` shorthand):

```console
$ docker image ls

REPOSITORY                       TAG       IMAGE ID       CREATED         SIZE
docker-gs-ping                   latest    7f153fbcc0a8   2 minutes ago   1.11GB
...
```

Your exact output may vary, but you should see the `docker-gs-ping` image with the
`latest` tag. Because we had not specified a custom tag when we built our image, 
Docker assumed that the tag would be `latest`, which is a special value.

## Tag images

An image name is made up of slash-separated name components. Name components may
contain lowercase letters, digits and separators. A separator is defined as a
period, one or two underscores, or one or more dashes. A name component may not
start or end with a separator.

An image is made up of a manifest and a list of layers. In simple terms, a “tag”
points to a combination of these artifacts. You can have multiple tags for the
image and, in fact, most images have multiple tags. Let’s create a second tag
for the image we had built and take a look at its layers.

Use the `docker image tag` (or `docker tag` shorthand) command to create a new
tag for our image. This command takes two arguments; the first argument is the
"source" image, and the second is the new tag to create. The following command
creates a new `docker-gs-ping:v1.0` tag for the `docker-gs-ping:latest` we built
above:

```console
$ docker image tag docker-gs-ping:latest docker-gs-ping:v1.0
```

The Docker `tag` command creates a new tag for the image. It does not create a
new image. The tag points to the same image and is just another way to reference
the image.

Now run the `docker image ls` command again to see the updated list of local
images:

```console
$ docker image ls

REPOSITORY                       TAG       IMAGE ID       CREATED         SIZE
docker-gs-ping                   latest    7f153fbcc0a8   6 minutes ago   1.11GB
docker-gs-ping                   v1.0      7f153fbcc0a8   6 minutes ago   1.11GB
...
```

You can see that we have two images that start with `docker-gs-ping`. We know
they are the same image because if you look at the `IMAGE ID` column, you can
see that the values are the same for the two images. This value is a unique
identifier Docker uses internally to identify the image.

Let’s remove the tag that we had just created. To do this, we’ll use the
`docker image rm` command, or the shorthand `docker rmi` (which stands for
"remove image"):

```console
$ docker image rm docker-gs-ping:v1.0
Untagged: docker-gs-ping:v1.0
```

Notice that the response from Docker tells us that the image has not been removed but only "untagged". 

Verify this by running the following command:

```console
$ docker image ls
```

You will see that the tag `v1.0` is no longer in the list of images kept by your Docker instance.

```
REPOSITORY                       TAG       IMAGE ID       CREATED         SIZE
docker-gs-ping                   latest    7f153fbcc0a8   7 minutes ago   1.11GB
...
```

The tag `v1.0` has been removed but we still have the `docker-gs-ping:latest`
tag available on our machine, so the image is there.

## Multi-stage builds

You may have noticed that our `docker-gs-ping` image weighs in at over a gigabyte (!!!), 
which is *a lot* for a tiny compiled Go application. You may also be wondering what happened
to the full suite of Go tools, including the compiler, after we had built our image.

The answer is that the full toolchain is still there, in the container image. 
Not only this is inconvenient because of the large file size, but it may also
present a security risk when the container is deployed.

These two issues can be solved by using [multi-stage builds](../../build/building/multi-stage.md).

In a nutshell, a multi-stage build can carry over the artifacts from one build stage into another,
and every build stage can be instantiated from a different base image.

Thus, in the following example, we are going to use a full-scale official Go image to build 
our application but then we'll copy the application binary into another image whose base
is very lean and does not include the Go toolchain or other optional components.

The `Dockerfile.multistage` in the sample application's repo has the following
content:

{% raw %}
```dockerfile
# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /docker-gs-ping /docker-gs-ping

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/docker-gs-ping"]
```
{% endraw %}

Since we have two Dockerfiles now, we have to tell Docker what Dockerfile we'd like to use
to build the image. Let's tag the new image with `multistage`. This tag (like any other, 
apart from `latest`) has no special meaning for Docker, it's just something we chose.

```console
$ docker build -t docker-gs-ping:multistage -f Dockerfile.multistage .
```

Comparing the sizes of `docker-gs-ping:multistage` and `docker-gs-ping:latest`
we see a *few orders-of-magnitude* difference! (`docker image ls`)

```
REPOSITORY       TAG          IMAGE ID       CREATED              SIZE
docker-gs-ping   multistage   e3fdde09f172   About a minute ago   28.1MB
docker-gs-ping   latest       336a3f164d0f   About an hour ago    1.11GB
```

This is so because the ["distroless"](https://github.com/GoogleContainerTools/distroless){:target="_blank" rel="noopener" class="_"} 
base image that we have used in the second stage of the build is very barebones and is designed for lean deployments of static binaries.

There's much more to multi-stage builds, including the possibility of multi-architecture builds, 
so please feel free to check out the [multi-stage builds](../../build/building/multi-stage.md) 
section of Docker documentation. This is, however, not essential for our progress here, so we'll
leave it at that.

## Next steps

In this module, we met our example application and built and container image for it. 

In the next module, we’ll take a look at how to:

[Run your image as a container](run-containers.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you
think by creating an issue in the [Docker Docs]({{ site.repo }}/issues/new?title=[Golang%20docs%20feedback]){:target="_blank" rel="noopener" class="_"}
GitHub repository. Alternatively, [create a PR]({{ site.repo }}/pulls){:target="_blank" rel="noopener" class="_"}
to suggest updates.
