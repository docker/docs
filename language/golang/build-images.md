---
title: "Build your Go image"
keywords: containers, images, go, golang, dockerfiles, coding, build, push, run
description: Learn how to build your first Docker image by writing a Dockerfile
redirect_from:
- /get-started/golang/build-images/
---

{% include_relative nav.html selected="1" %}

## Prerequisites

* Some understanding of Go and its toolchain. This is not a tutorial on Go. If you are new to the language, the [Go website](https://golang.org/){: target="_blank" rel="noopener" class="_"} is a good starting point, so go (pun intended) check it out.
* Some awareness of basic Docker concepts. If unsure, work through the orientation and setup in Get started [Part 1](/get-started/).

## Overview

Now that we have a good overview of containers and the Docker platform, let’s take a look at building our first image. An image includes everything you need to run an application - the code or binary, runtime, dependencies, and any other file system objects required.

To complete this tutorial, you need the following:

- Go version 1.16 or later. 
- Docker running locally. Follow the instructions to [download and install Docker](https://docs.docker.com/desktop/).
- An IDE or a text editor to edit files. We recommend using [Visual Studio Code](https://code.visualstudio.com/){: target="_blank" rel="noopener" class="_"}.

## Introducting the sample application

To avoid losing focus on Docker's features, the sample application is a minimal HTTP server that has only three features:

* It responds with a text message containing a heart symbol ("<3") on requests to `/`.
* It responds with `{"Status" : "OK"}` to the health check request on requests to `/ping`.
* The port it listens on is configurable using the environment variable `HTTP_PORT`. The default value is `8080`.

Thus, it somewhat mimics enough basic properties of a REST microservice to be useful for our learning of Docker.

The source code for the application is in the [olliefr/docker-gs-ping](https://github.com/olliefr/docker-gs-ping){: target="_blank" rel="noopener" class="_"} GitHub repo. Please feel free to clone or fork it.

For our present study, we clone it to our local machine:

```shell
$ git clone github.com/olliefr/docker-gs-ping
```

The application's `main.go` file is fairly straightforward, if you are familiar with Go:

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

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
```

## Quickly smoke test the application

Let’s start our application and make sure it’s running properly. Open your terminal and navigate to the directory into which you cloned the project's repo.

```shell
$ go run main.go
```

This should compile and start the server as a foreground application, outputting the banner and some logging information, as illustrated in the next figure.

```shell
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.2.2
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

Let's run a quick _smoke test_ on the application. **In a new terminal**, run a request using `curl`. Alternatively, you can use your favourite web browser as well.

```shell
$ curl http://localhost:8080/
Hello, Docker! <3
```

So, the application responds with a greeting, just as the first business requirement says it should. Great.

Having established that the server is running and is accessible, let's proceed to "dockerising" it.

## Create a Dockerfile for the application

A `Dockerfile` is a text document that contains all the commands a user could call on the command line to assemble an image. When we tell Docker to build our image by executing the `docker build` command, Docker reads these instructions and executes them one by one and creates a Docker image as a result.

Let’s walk through the process of creating a `Dockerfile` for our application. In the root of your working directory, create a file named `Dockerfile` and open this file in your text editor.

> **Note**
>
> The name of the file is not _that_ important but the default filename for many commands is simply `Dockerfile`. So, we’ll use that as our filename throughout this series.

The first thing we need to do is to add a line in our Dockerfile that tells Docker what base image we would like to use for our application.

```dockerfile
FROM golang:1.16-alpine
```

Docker images can be inherited from other images. Therefore, instead of creating our own base image, we’ll use the official Go image that already has all the tools and packages to compile and run a Go application. You can think of this in the same way you would think about class inheritance in object oriented programming or functional composition in functional programming.

When we have used that `FROM` command, we told Docker to include in our image all the functionality from the `golang:1.16-alpine` image. All of our consequent commands would build on top of that "base" image.

> **Note**
>
> If you want to learn more about creating your own base images, see [creating base images](https://docs.docker.com/develop/develop-images/baseimages/) section of the guide.

To make things easier when running the rest of our commands, let’s create a working directory _inside_ the image that we are building. This instructs Docker to use this path as the default _destination path_ for all subsequent commands. This way we do not have to type out full file paths but can use relative paths based on this working directory.

```dockerfile
WORKDIR /app
```

Usually the very first thing you do once you’ve downloaded a project written in Go is to install the modules necessary to compile it.

Before we can run `go mod download`, we need to get our `go.mod` and `go.sum` files into our images. We use the `COPY` command to do this. 

In its simplest form, the `COPY` command takes two parameters. The first parameter tells Docker what file you would like to copy into the image. The second parameter tells Docker where you want that file to be copied to. 

We’ll copy the `go.mod` and `go.sum` file into our working directory `/app` which, owing to our use of `WORKDIR`, is the current directory (`.`)

```dockerfile
COPY go.mod .
COPY go.sum .
```

Now that we have the module files inside the Docker image that we are building, we can use the `RUN` command to execute the command `go mod download` there as well. This works exactly the same as if we were running `go` locally on our machine, but this time these Go modules will be installed into the a directory inside our image.

```dockerfile
RUN go mod download
```

At this point, we have an image that is based on Go environment version 1.16 (or a later minor version, since we had specified `1.16` as our tag in the `FROM` command) and we have installed our dependencies. 

The next thing we need to do is to add our source code into the image. We’ll use the `COPY` command just like we did with our module files before.

```dockerfile
COPY *.go .
```

This `COPY` command uses a wildcard to copy all files with `.go` extension located in the current directory into the image. 

Now, we would like to compile our application. To that end, we use the familiar `RUN` command.

```dockerfile
RUN go build -o /docker-gs-ping
```

This should be familiar. The result of that command will be a static application binary named `docker-gs-ping` and located in the root of the filesystem of the image that we are building. We could have put the binary into any other place we desire inside that image, there is no special significance for placing it in the root.

Now, all that is left to do is to tell Docker what command to execute when our image is used to start a container. 

We do this with the `CMD` command.

```dockerfile
CMD [ "/docker-gs-ping" ]
```

Here's the complete `Dockerfile`. The file in the repository that you cloned may also contain comments. They always begin with a `#` symbol and make no difference to Docker. The comments are there for the convenience of humans tasked to maintain the `Dockerfile`.

```dockerfile
FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .

RUN go build -o /docker-gs-ping

EXPOSE 8080

CMD [ "/docker-gs-ping" ]
```

You can also find this `Dockerfile` for this part in [go-docker](https://github.com/olliefr/go-docker) repo.

## Build the image

Now that we've created our `Dockerfile`, let’s build an image from it. The `docker build` command creates Docker images from a Dockerfile and a "context". A build's _context_ is the set of files located in the specified `PATH` or `URL`. The Docker build process can access any of the files located in the context.

The build command optionally takes a `--tag` flag. This flag is used to label the image with a string value, which is easy for humans to read and recognise. If you do not pass a `--tag`, Docker will use `latest` as the default value.

Let's build our first Docker image!

```shell
$ docker build --tag docker-gs-ping .
```

```shell
[+] Building 3.6s (12/12) FINISHED
 => [internal] load build definition from Dockerfile                                      0.1s
 => => transferring dockerfile: 38B                                                       0.0s
 => [internal] load .dockerignore                                                         0.1s
 => => transferring context: 2B                                                           0.0s
 => [internal] load metadata for docker.io/library/golang:1.16-alpine                     3.0s
 => [1/7] FROM docker.io/library/golang:1.16-alpine@sha256:49c07aa83790aca732250c2258b59  0.0s
 => => resolve docker.io/library/golang:1.16-alpine@sha256:49c07aa83790aca732250c2258b59  0.0s
 => [internal] load build context                                                         0.1s
 => => transferring context: 114B                                                         0.0s
 => CACHED [2/7] WORKDIR /app                                                             0.0s
 => CACHED [3/7] COPY go.mod .                                                            0.0s
 => CACHED [4/7] COPY go.sum .                                                            0.0s
 => CACHED [5/7] RUN go mod download                                                      0.0s
 => CACHED [6/7] COPY *.go .                                                              0.0s
 => CACHED [7/7] RUN go build -o /docker-gs-ping                                          0.0s
 => exporting to image                                                                    0.1s
 => => exporting layers                                                                   0.0s
 => => writing image sha256:336a3f164d0f079f2e42cd1d38f24ab9110d47d481f1db7f2a0b0d2859ec  0.0s
 => => naming to docker.io/library/docker-gs-ping                                         0.0s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
```

Your exact output will vary, but provided there aren't any errors, you should see the `FINISHED` line in the build output.

## Viewing local images

To see a list of images we have on our local machine, we have two options. One is to use the CLI and the other is to use [Docker Desktop](/products/docker-desktop). Since we are currently working in the terminal, let’s take a look at listing images with the CLI.

To list images, simply run the `images` command.

```shell
$ docker images
```

```
REPOSITORY       TAG       IMAGE ID       CREATED          SIZE
docker-gs-ping   latest    336a3f164d0f   39 minutes ago   540MB
postgres         13.2      c5ec7353d87d   7 weeks ago      314MB
```

Your exact output may vary, but you should see `docker-gs-ping` image with the `latest` tag.

## Tag images

An image name is made up of slash-separated name components. Name components may contain lowercase letters, digits and separators. A separator is defined as a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.

An image is made up of a manifest and a list of layers. In simple terms, a “tag” points to a combination of these artifacts. You can have multiple tags for an image. Let’s create a second tag for the image we built and take a look at its layers.

To create a new tag for the image we built above, run the following command.

```shell
$ docker tag docker-gs-ping:latest docker-gs-ping:v1.0
```

The Docker tag command creates a new tag for an image. It does not create a new image. The tag points to the same image and is just another way to reference the image.

Now run the `docker images` command to see a list of our local images.

```shell
$ docker images
```

```
REPOSITORY       TAG       IMAGE ID       CREATED          SIZE
docker-gs-ping   latest    336a3f164d0f   43 minutes ago   540MB
docker-gs-ping   v1.0      336a3f164d0f   43 minutes ago   540MB
postgres         13.2      c5ec7353d87d   7 weeks ago      314MB
```

You can see that we have two images that start with `docker-gs-ping`. We know they are the same image because if you look at the `IMAGE ID` column, you can see that the values are the same for the two images. This value is a checksum Docker uses internally to identify the image.

Let’s remove the tag that we just created. To do this, we’ll use the `rmi` command. The `rmi` command stands for "remove image".

```shell
$ docker rmi docker-gs-ping:v1.0
Untagged: docker-gs-ping:v1.0
```

Notice that the response from Docker tells us that the image has not been removed but only "untagged". Verify this by running the images command.

```shell
$ docker images
```

```
REPOSITORY       TAG       IMAGE ID       CREATED          SIZE
docker-gs-ping   latest    336a3f164d0f   45 minutes ago   540MB
postgres         13.2      c5ec7353d87d   7 weeks ago      314MB
```

Our image that was tagged with `v1.0` has been removed but we still have the `docker-gs-ping:latest` tag available on our machine.

## Multi-stage builds

You may have noticed that our `docker-gs-ping` image stands at 540MB, which you may think is a lot. You may also be wondering whether our dockerised application still needs the full suite of Go tools, including the compiler, after the application binary had been compiled.

These are legit concerns. Both can be solved by using _multi-stage builds_. The following example is provided with little explanation because this would derail us from our current concerns, but please feel free to explore on your own later.

The `Dockerfile.multistage` in the sample application's repo has the following content.

```dockerfile
##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .

RUN go build -o /docker-gs-ping

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-gs-ping /docker-gs-ping

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/docker-gs-ping"]
```

The application can be built and tagged with the following.

```shell
docker build -t docker-gs-ping:multistage -f Dockerfile.multistage .
```

Compare the sizes of `docker-gs-ping:multistage` and `docker-gs-ping:latest` we see an order-of-magnitude difference!

```
REPOSITORY       TAG          IMAGE ID       CREATED              SIZE
docker-gs-ping   multistage   e3fdde09f172   About a minute ago   27.1MB
docker-gs-ping   latest       336a3f164d0f   About an hour ago    540MB
```

This is due to the fact that ["distroless" base images](https://github.com/GoogleContainerTools/distroless) are very barebones and are meant for lean deployments of static binaries, such as our Go server application.

For more information on multi-stage builds, please feel free to check out [other parts](/develop/develop-images/multistage-build/) of Docker documentation. This is, however, not essential for our progress here, so we'll leave it at that.

## Next steps

In this module, we took a look at setting up our example Go application that we will use for the rest of the tutorial. We also created a Dockerfile that we used to build our Docker image. Then, we took a look at tagging our images and removing images. In the next module, we’ll take a look at how to:

[Run your image as a container](run-containers.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs ](https://github.com/docker/docker.github.io/issues/new?title=[Golang%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
