---
title: Go language-specific guide
linkTitle: Go
description: Containerize Go apps using Docker
keywords: docker, getting started, go, golang, language, dockerfile
summary: |
  This guide teaches you how to containerize Go applications using Docker.
aliases:
  - /language/golang/
  - /guides/language/golang/
  - /get-started/golang/build-images/
  - /language/golang/build-images/
  - /get-started/golang/run-containers/
  - /language/golang/run-containers/
  - /get-started/golang/develop/
  - /language/golang/develop/
  - /get-started/golang/run-tests/
  - /language/golang/run-tests/
  - /language/golang/configure-ci-cd/
  - /language/golang/deploy/
  - /guides/golang/build-images/
  - /guides/golang/configure-ci-cd/
  - /guides/golang/deploy/
  - /guides/golang/develop/
  - /guides/golang/run-containers/
  - /guides/golang/run-tests/
params:
  tags: [languages]
  time: 30 minutes
---

This guide will show you how to create, test, and deploy containerized Go applications using Docker.

> **Acknowledgment**
>
> Docker would like to thank [Oliver Frolovs](https://www.linkedin.com/in/ofr/) for his contribution to this guide.

## What will you learn?

In this guide, you’ll learn how to:

- Create a `Dockerfile` which contains the instructions for building a container image for a program written in Go.
- Run the image as a container in your local Docker instance and manage the container's lifecycle.
- Use multi-stage builds for building small images efficiently while keeping your Dockerfiles easy to read and maintain.
- Use Docker Compose to orchestrate running of multiple related containers together in a development environment.

## Prerequisites

Some basic understanding of Go and its toolchain is assumed. This isn't a Go tutorial. If you are new to the : languages:,
the [Go website](https://golang.org/) is a great place to explore,
so _go_ (pun intended) check it out!

You also must know some basic [Docker concepts](/get-started/docker-concepts/the-basics/what-is-a-container.md) as well as to
be at least vaguely familiar with the [Dockerfile format](/manuals/build/concepts/dockerfile.md).

Your Docker set-up must have BuildKit enabled. BuildKit is enabled by default for all users on [Docker Desktop](/manuals/desktop/_index.md).
If you have installed Docker Desktop, you don’t have to manually enable BuildKit. If you are running Docker on Linux,
please check out BuildKit [getting started](/manuals/build/buildkit/_index.md#getting-started) page.

Some familiarity with the command line is also expected.

## What's next?

The aim of this guide is to provide enough examples and instructions for you to containerize your own Go application and deploy it into the Cloud.

Start by building your first Go image.

## Build your Go image

### Overview

In this section you're going to build a container image. The image includes
everything you need to run your application – the compiled application binary
file, the runtime, the libraries, and all other resources required by your
application.

### Required software

To complete this tutorial, you need the following:

- Docker running locally. Follow the [instructions to download and install Docker](/manuals/desktop/_index.md).
- An IDE or a text editor to edit files. [Visual Studio Code](https://code.visualstudio.com/) is a free and popular choice but you can use anything you feel comfortable with.
- A Git client. This guide uses a command-line based `git` client, but you are free to use whatever works for you.
- A command-line terminal application. The examples shown in this module are from the Linux shell, but they should work in PowerShell, Windows Command Prompt, or OS X Terminal with minimal, if any, modifications.

### Meet the example application

The example application is a caricature of a microservice. It is purposefully trivial to keep focus on learning the basics of containerization for Go applications.

The application offers two HTTP endpoints:

- It responds with a string containing a heart symbol (`<3`) to requests to `/`.
- It responds with `{"Status" : "OK"}` JSON to a request to `/health`.

It responds with HTTP error 404 to any other request.

The application listens on a TCP port defined by the value of environment variable `PORT`. The default value is `8080`.

The application is stateless.

The complete source code for the application is on GitHub: [github.com/docker/docker-gs-ping](https://github.com/docker/docker-gs-ping). You are encouraged to fork it and experiment with it as much as you like.

To continue, clone the application repository to your local machine:

```console
$ git clone https://github.com/docker/docker-gs-ping
```

The application's `main.go` file is straightforward, if you are familiar with Go:

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

### Create a Dockerfile for the application

To build a container image with Docker, a `Dockerfile` with build instructions is required.

Begin your `Dockerfile` with the (optional) parser directive line that instructs BuildKit to
interpret your file according to the grammar rules for the specified version of the syntax.

You then tell Docker what base image you would like to use for your application:

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.19
```

Docker images can be inherited from other images. Therefore, instead of creating
your own base image from scratch, you can use the official Go image that already
has all necessary tools and libraries to compile and run a Go application.

> [!NOTE]
>
> If you are curious about creating your own base images, you can check out the following section of this guide: [creating base images](/manuals/build/building/base-images.md#create-a-base-image).
> Note, however, that this isn't necessary to continue with your task at hand.

Now that you have defined the base image for your upcoming container image, you
can begin building on top of it.

To make things easier when running the rest of your commands, create a directory
inside the image that you're building. This also instructs Docker to use this
directory as the default destination for all subsequent commands. This way you
don't have to type out full file paths in the `Dockerfile`, the relative paths
will be based on this directory.

```dockerfile
WORKDIR /app
```

Usually the very first thing you do once you’ve downloaded a project written in
Go is to install the modules necessary to compile it. Note, that the base image
has the toolchain already, but your source code isn't in it yet.

So before you can run `go mod download` inside your image, you need to get your
`go.mod` and `go.sum` files copied into it. Use the `COPY` command to do this.

In its simplest form, the `COPY` command takes two parameters. The first
parameter tells Docker what files you want to copy into the image. The last
parameter tells Docker where you want that file to be copied to.

Copy the `go.mod` and `go.sum` file into your project directory `/app` which,
owing to your use of `WORKDIR`, is the current directory (`./`) inside the
image. Unlike some modern shells that appear to be indifferent to the use of
trailing slash (`/`), and can figure out what the user meant (most of the time),
Docker's `COPY` command is quite sensitive in its interpretation of the trailing
slash.

```dockerfile
COPY go.mod go.sum ./
```

> [!NOTE]
>
> If you'd like to familiarize yourself with the trailing slash treatment by the
> `COPY` command, see [Dockerfile
> reference](/reference/dockerfile.md#copy). This trailing slash can
> cause issues in more ways than you can imagine.

Now that you have the module files inside the Docker image that you are
building, you can use the `RUN` command to run the command `go mod download`
there as well. This works exactly the same as if you were running `go` locally
on your machine, but this time these Go modules will be installed into a
directory inside the image.

```dockerfile
RUN go mod download
```

At this point, you have a Go toolchain version 1.19.x and all your Go
dependencies installed inside the image.

The next thing you need to do is to copy your source code into the image. You’ll
use the `COPY` command just like you did with your module files before.

```dockerfile
COPY *.go ./
```

This `COPY` command uses a wildcard to copy all files with `.go` extension
located in the current directory on the host (the directory where the `Dockerfile`
is located) into the current directory inside the image.

Now, to compile your application, use the familiar `RUN` command:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping
```

This should be familiar. The result of that command will be a static application
binary named `docker-gs-ping` and located in the root of the filesystem of the
image that you are building. You could have put the binary into any other place
you desire inside that image, the root directory has no special meaning in this
regard. It's convenient to use it to keep the file paths short for improved
readability.

Now, all that is left to do is to tell Docker what command to run when your
image is used to start a container.

You do this with the `CMD` command:

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
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]
```

The `Dockerfile` may also contain comments. They always begin with a `#` symbol,
and must be at the beginning of a line. Comments are there for your convenience
to allow documenting your `Dockerfile`.

There is also a concept of Dockerfile directives, such as the `syntax` directive
you added. The directives must always be at the very top of the `Dockerfile`, so
when adding comments, make sure that the comments follow after any directives
that you may have used:

```dockerfile
# syntax=docker/dockerfile:1
# A sample microservice in Go packaged into a container image.

FROM golang:1.19

# ...
```

### Build the image

Now that you've created your `Dockerfile`, build an image from it. The `docker
build` command creates Docker images from the `Dockerfile` and a context. A
build context is the set of files located in the specified path or URL. The
Docker build process can access any of the files located in the context.

The build command optionally takes a `--tag` flag. This flag is used to label
the image with a string value, which is easy for humans to read and recognize.
If you don't pass a `--tag`, Docker will use `latest` as the default value.

Build your first Docker image.

```console
$ docker build --tag docker-gs-ping .
```

The build process will print some diagnostic messages as it goes through the build steps.
The following is an example of what these messages may look like.

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
see the word `FINISHED` in the first line of output. This means Docker has
successfully built your image named `docker-gs-ping`.

### View local images

To see the list of images you have on your local machine, you have two options.
One is to use the CLI and the other is to use [Docker
Desktop](/manuals/desktop/_index.md). Since you're working in the
terminal, take a look at listing images with the CLI.

To list images, run the `docker image ls`command (or the `docker images` shorthand):

```console
$ docker image ls

REPOSITORY                       TAG       IMAGE ID       CREATED         SIZE
docker-gs-ping                   latest    7f153fbcc0a8   2 minutes ago   1.11GB
...
```

Your exact output may vary, but you should see the `docker-gs-ping` image with
the `latest` tag. Because you didn't specify a custom tag when you built your
image, Docker assumed that the tag would be `latest`, which is a special value.

### Tag images

An image name is made up of slash-separated name components. Name components may
contain lowercase letters, digits, and separators. A separator is defined as a
period, one or two underscores, or one or more dashes. A name component may not
start or end with a separator.

An image is made up of a manifest and a list of layers. In simple terms, a tag
points to a combination of these artifacts. You can have multiple tags for the
image and, in fact, most images have multiple tags. Create a second tag
for the image you built and take a look at its layers.

Use the `docker image tag` (or `docker tag` shorthand) command to create a new
tag for your image. This command takes two arguments; the first argument is the
source image, and the second is the new tag to create. The following command
creates a new `docker-gs-ping:v1.0` tag for the `docker-gs-ping:latest` you
built:

```console
$ docker image tag docker-gs-ping:latest docker-gs-ping:v1.0
```

The Docker `tag` command creates a new tag for the image. It doesn't create a
new image. The tag points to the same image and is another way to reference
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

You can see that you have two images that start with `docker-gs-ping`. You know
they're the same image because if you look at the `IMAGE ID` column, you can
see that the values are the same for the two images. This value is a unique
identifier Docker uses internally to identify the image.

Remove the tag that you just created. To do this, you’ll use the
`docker image rm` command, or the shorthand `docker rmi` (which stands for
"remove image"):

```console
$ docker image rm docker-gs-ping:v1.0
Untagged: docker-gs-ping:v1.0
```

Notice that the response from Docker tells you that the image hasn't been
removed but only untagged.

Verify this by running the following command:

```console
$ docker image ls
```

You will see that the tag `v1.0` is no longer in the list of images kept by your Docker instance.

```text
REPOSITORY                       TAG       IMAGE ID       CREATED         SIZE
docker-gs-ping                   latest    7f153fbcc0a8   7 minutes ago   1.11GB
...
```

The tag `v1.0` has been removed but you still have the `docker-gs-ping:latest`
tag available on your machine, so the image is there.

### Multi-stage builds

You may have noticed that your `docker-gs-ping` image weighs in at over a
gigabyte, which is a lot for a tiny compiled Go application. You may also be
wondering what happened to the full suite of Go tools, including the compiler,
after you had built your image.

The answer is that the full toolchain is still there, in the container image.
Not only this is inconvenient because of the large file size, but it may also
present a security risk when the container is deployed.

These two issues can be solved by using [multi-stage builds](/manuals/build/building/multi-stage.md).

In a nutshell, a multi-stage build can carry over the artifacts from one build stage into another,
and every build stage can be instantiated from a different base image.

Thus, in the following example, you are going to use a full-scale official Go
image to build your application. Then you'll copy the application binary into
another image whose base is very lean and doesn't include the Go toolchain or
other optional components.

The `Dockerfile.multistage` in the sample application's repository has the
following content:

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

Since you have two Dockerfiles now, you have to tell Docker what Dockerfile
you'd like to use to build the image. Tag the new image with `multistage`. This
tag (like any other, apart from `latest`) has no special meaning for Docker,
it's something you chose.

```console
$ docker build -t docker-gs-ping:multistage -f Dockerfile.multistage .
```

Comparing the sizes of `docker-gs-ping:multistage` and `docker-gs-ping:latest`
you see a few orders-of-magnitude difference.

```console
$ docker image ls
REPOSITORY       TAG          IMAGE ID       CREATED              SIZE
docker-gs-ping   multistage   e3fdde09f172   About a minute ago   28.1MB
docker-gs-ping   latest       336a3f164d0f   About an hour ago    1.11GB
```

This is so because the ["distroless"](https://github.com/GoogleContainerTools/distroless)
base image that you have used in the second stage of the build is very barebones and is designed for lean deployments of static binaries.

There's much more to multi-stage builds, including the possibility of multi-architecture builds,
so feel free to check out [multi-stage builds](/manuals/build/building/multi-stage.md). This is, however, not essential for your progress here.

### Next steps

In this module, you met your example application and built and container image
for it.

In the next module, you’ll take a look at how to run your image as a container.

## Run your Go image as a container

### Prerequisites

Work through the steps to containerize a Go application in [Build your Go image](./).

### Overview

In the previous module you created a `Dockerfile` for your example application and then you created your Docker image using the command `docker build`. Now that you have the image, you can run that image and see if your application is running correctly.

A container is a normal operating system process except that this process is isolated and has its own file system, its own networking, and its own isolated process tree separate from the host.

To run an image inside of a container, you use the `docker run` command. It requires one parameter and that's the image name. Start your image and make sure it's running correctly. Run the following command in your terminal.

```console
$ docker run docker-gs-ping
```

```text
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

When you run this command, you’ll notice that you weren't returned to the command prompt. This is because your application is a REST server and will run in a loop waiting for incoming requests without returning control back to the OS until you stop the container.

Make a GET request to the server using the curl command.

```console
$ curl http://localhost:8080/
curl: (7) Failed to connect to localhost port 8080: Connection refused
```

Your curl command failed because the connection to your server was refused.
Meaning that you weren't able to connect to localhost on port 8080. This is
expected because your container is running in isolation which includes
networking. Stop the container and restart with port 8080 published on your
local network.

To stop the container, press ctrl-c. This will return you to the terminal prompt.

To publish a port for your container, you’ll use the `--publish` flag (`-p` for short) on the `docker run` command. The format of the `--publish` command is `[host_port]:[container_port]`. So if you wanted to expose port `8080` inside the container to port `3000` outside the container, you would pass `3000:8080` to the `--publish` flag.

Start the container and expose port `8080` to port `8080` on the host.

```console
$ docker run --publish 8080:8080 docker-gs-ping
```

Now, rerun the curl command.

```console
$ curl http://localhost:8080/
Hello, Docker! <3
```

Success! You were able to connect to the application running inside of your container on port 8080. Switch back to the terminal where your container is running and you should see the `GET` request logged to the console.

Press `ctrl-c` to stop the container.

### Run in detached mode

This is great so far, but your sample application is a web server and you
shouldn't have to have your terminal connected to the container. Docker can run
your container in detached mode in the background. To do this, you can use the
`--detach` or `-d` for short. Docker will start your container the same as
before but this time will detach from the container and return you to the
terminal prompt.

```console
$ docker run -d -p 8080:8080 docker-gs-ping
d75e61fcad1e0c0eca69a3f767be6ba28a66625ce4dc42201a8a323e8313c14e
```

Docker started your container in the background and printed the container ID on the terminal.

Again, make sure that your container is running. Run the same `curl` command:

```console
$ curl http://localhost:8080/
Hello, Docker! <3
```

### List containers

Since you ran your container in the background, how do you know if your container is running or what other containers are running on your machine? Well, to see a list of containers running on your machine, run `docker ps`. This is similar to how the ps command is used to see a list of processes on a Linux machine.

```console
$ docker ps

CONTAINER ID   IMAGE            COMMAND             CREATED          STATUS          PORTS                    NAMES
d75e61fcad1e   docker-gs-ping   "/docker-gs-ping"   41 seconds ago   Up 40 seconds   0.0.0.0:8080->8080/tcp   inspiring_ishizaka
```

The `ps` command tells you a bunch of stuff about your running containers. You can see the container ID, the image running inside the container, the command that was used to start the container, when it was created, the status, ports that are exposed, and the names of the container.

You are probably wondering where the name of your container is coming from. Since you didn’t provide a name for the container when you started it, Docker generated a random name. You'll fix this in a minute but first you need to stop the container. To stop the container, run the `docker stop` command, passing the container's name or ID.

```console
$ docker stop inspiring_ishizaka
inspiring_ishizaka
```

Now rerun the `docker ps` command to see a list of running containers.

```console
$ docker ps

CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
```

### Stop, start, and name containers

Docker containers can be started, stopped and restarted. When you stop a container, it's not removed but the status is changed to stopped and the process inside of the container is stopped. When you ran the `docker ps` command, the default output is to only show running containers. If you pass the `--all` or `-a` for short, you will see all containers on your system, including stopped containers and running containers.

```console
$ docker ps --all

CONTAINER ID   IMAGE            COMMAND                  CREATED              STATUS                      PORTS     NAMES
d75e61fcad1e   docker-gs-ping   "/docker-gs-ping"        About a minute ago   Exited (2) 23 seconds ago             inspiring_ishizaka
f65dbbb9a548   docker-gs-ping   "/docker-gs-ping"        3 minutes ago        Exited (2) 2 minutes ago              wizardly_joliot
aade1bf3d330   docker-gs-ping   "/docker-gs-ping"        3 minutes ago        Exited (2) 3 minutes ago              magical_carson
52d5ce3c15f0   docker-gs-ping   "/docker-gs-ping"        9 minutes ago        Exited (2) 3 minutes ago              gifted_mestorf
```

If you’ve been following along, you should see several containers listed. These are containers that you started and stopped but haven't removed yet.

Restart the container that you have just stopped. Locate the name of the container and replace the name of the container in the following `restart` command:

```console
$ docker restart inspiring_ishizaka
```

Now, list all the containers again using the `ps` command:

```console
$ docker ps -a

CONTAINER ID   IMAGE            COMMAND                  CREATED          STATUS                     PORTS                    NAMES
d75e61fcad1e   docker-gs-ping   "/docker-gs-ping"        2 minutes ago    Up 5 seconds               0.0.0.0:8080->8080/tcp   inspiring_ishizaka
f65dbbb9a548   docker-gs-ping   "/docker-gs-ping"        4 minutes ago    Exited (2) 2 minutes ago                            wizardly_joliot
aade1bf3d330   docker-gs-ping   "/docker-gs-ping"        4 minutes ago    Exited (2) 4 minutes ago                            magical_carson
52d5ce3c15f0   docker-gs-ping   "/docker-gs-ping"        10 minutes ago   Exited (2) 4 minutes ago                            gifted_mestorf
```

Notice that the container you just restarted has been started in detached mode and has port `8080` exposed. Also, note that the status of the container is `Up X seconds`. When you restart a container, it will be started with the same flags or commands that it was originally started with.

Stop and remove all of your containers and take a look at fixing the random naming issue.

Stop the container you just started. Find the name of your running container and replace the name in the following command with the name of the container on your system:

```console
$ docker stop inspiring_ishizaka
inspiring_ishizaka
```

Now that all of your containers are stopped, remove them. When a container is removed, it's no longer running nor is it in the stopped state. Instead, the process inside the container is terminated and the metadata for the container is removed.

To remove a container, run the `docker rm` command passing the container name. You can pass multiple container names to the command in one command.

Again, make sure you replace the containers names in the following command with the container names from your system:

```console
$ docker rm inspiring_ishizaka wizardly_joliot magical_carson gifted_mestorf

inspiring_ishizaka
wizardly_joliot
magical_carson
gifted_mestorf
```

Run the `docker ps --all` command again to verify that all containers are gone.

Now, address the pesky random name issue. Standard practice is to name your containers for the simple reason that it's easier to identify what's running in the container and what application or service it's associated with. Just like good naming conventions for variables in your code makes it simpler to read. So goes naming your containers.

To name a container, you must pass the `--name` flag to the `run` command:

```console
$ docker run -d -p 8080:8080 --name rest-server docker-gs-ping
3bbc6a3102ea368c8b966e1878a5ea9b1fc61187afaac1276c41db22e4b7f48f
```

```console
$ docker ps

CONTAINER ID   IMAGE            COMMAND             CREATED          STATUS          PORTS                    NAMES
3bbc6a3102ea   docker-gs-ping   "/docker-gs-ping"   25 seconds ago   Up 24 seconds   0.0.0.0:8080->8080/tcp   rest-server
```

Now, you can easily identify your container based on the name.

### Next steps

In this module, you learned how to run containers and publish ports. You also learned to manage the lifecycle of containers. You then learned the importance of naming your containers so that they're more easily identifiable. In the next module, you’ll learn how to run a database in a container and connect it to your application.

## Use containers for Go development

### Prerequisites

Work through the steps of the [run your image as a container](./) module to learn how to manage the lifecycle of your containers.

### Introduction

In this module, you'll take a look at running a database engine in a container and connecting it to the extended version of the example application. You are going to see some options for keeping persistent data and for wiring up the containers to talk to one another. Finally, you'll learn how to use Docker Compose to manage such multi-container local development environments effectively.

### Local database and containers

The database engine you are going to use is called [CockroachDB](https://www.cockroachlabs.com/product/). It is a modern, Cloud-native, distributed SQL database.

Instead of compiling CockroachDB from the source code or using the operating system's native package manager to install CockroachDB, you are going to use the [Docker image for CockroachDB](https://hub.docker.com/r/cockroachdb/cockroach) and run it in a container.

CockroachDB is compatible with PostgreSQL to a significant extent, and shares many conventions with the latter, particularly the default names for the environment variables. So, if you are familiar with Postgres, don't be surprised if you see some familiar environment variable names. The Go modules that work with Postgres, such as [pgx](https://pkg.go.dev/github.com/jackc/pgx), [pq](https://pkg.go.dev/github.com/lib/pq), [GORM](https://gorm.io/index.html), and [upper/db](https://upper.io/v4/) also work with CockroachDB.

For more information on the relation between Go and CockroachDB, refer to the [CockroachDB documentation](https://www.cockroachlabs.com/docs/v20.2/build-a-go-app-with-cockroachdb.html), although this isn't necessary to continue with the present guide.

#### Storage

The point of a database is to have a persistent store of data. [Volumes](/manuals/engine/storage/volumes.md) are the preferred mechanism for persisting data generated by and used by Docker containers. Thus, before you start CockroachDB, create the volume for it.

To create a managed volume, run :

```console
$ docker volume create roach
roach
```

You can view the list of all managed volumes in your Docker instance with the following command:

```console
$ docker volume list
DRIVER    VOLUME NAME
local     roach
```

#### Networking

The example application and the database engine are going to talk to one another over the network. There are different kinds of network configuration possible, and you're going to use what's called a user-defined bridge network. It is going to provide you with a DNS lookup service so that you can refer to your database engine container by its hostname.

The following command creates a new bridge network named `mynet`:

```console
$ docker network create -d bridge mynet
51344edd6430b5acd121822cacc99f8bc39be63dd125a3b3cd517b6485ab7709
```

As it was the case with the managed volumes, there is a command to list all networks set up in your Docker instance:

```console
$ docker network list
NETWORK ID     NAME          DRIVER    SCOPE
0ac2b1819fa4   bridge        bridge    local
51344edd6430   mynet         bridge    local
daed20bbecce   host          host      local
6aee44f40a39   none          null      local
```

Your bridge network `mynet` has been created successfully. The other three networks, named `bridge`, `host`, and `none` are the default networks and they had been created by the Docker itself. While it's not relevant to this guide, you can learn more about Docker networking in the [networking overview](/manuals/engine/network/_index.md) section.

#### Choose good names for volumes and networks

As the saying goes, there are only two hard things in Computer Science: cache invalidation and naming things. And off-by-one errors.

When choosing a name for a network or a managed volume, it's best to choose a name that's indicative of the intended purpose. This guide aims for brevity, so it used short, generic names.

#### Start the database engine

Now that the housekeeping chores are done, you can run CockroachDB in a container and attach it to the volume and network you had just created. When you run the following command, Docker will pull the image from Docker Hub and run it for you locally:

```console
$ docker run -d \
  --name roach \
  --hostname db \
  --network mynet \
  -p 26257:26257 \
  -p 8080:8080 \
  -v roach:/cockroach/cockroach-data \
  cockroachdb/cockroach:latest-v25.4 start-single-node \
  --insecure

# ... output omitted ...
```

Notice a clever use of the tag `latest-v25.4` to make sure that you're pulling the latest patch version of 25.4. The diversity of available tags depend on the image maintainer. Here, your intent was to have the latest patched version of CockroachDB while not straying too far away from the known working version as the time goes by. To see the tags available for the CockroachDB image, you can go to the [CockroachDB page on Docker Hub](https://hub.docker.com/r/cockroachdb/cockroach/tags).

#### Configure the database engine

Now that the database engine is live, there is some configuration to do before your application can begin using it. Fortunately, it's not a lot. You must:

1. Create a blank database.
2. Register a new user account with the database engine.
3. Grant that new user access rights to the database.

You can do that with the help of CockroachDB built-in SQL shell. To start the SQL shell in the same container where the database engine is running, type:

```console
$ docker exec -it roach ./cockroach sql --insecure
```

1. In the SQL shell, create the database that the example application is going to use:

   ```sql
   CREATE DATABASE mydb;
   ```

2. Register a new SQL user account with the database engine. Use the username `totoro`.

   ```sql
   CREATE USER totoro;
   ```

3. Give the new user the necessary permissions:

   ```sql
   GRANT ALL ON DATABASE mydb TO totoro;
   ```

4. Type `quit` to exit the shell.

The following is an example of interaction with the SQL shell.

```console
$ sudo docker exec -it roach ./cockroach sql --insecure
#
# Welcome to the CockroachDB SQL shell.
# All statements must be terminated by a semicolon.
# To exit, type: \q.
#
# Server version: CockroachDB CCL v20.1.15 (x86_64-unknown-linux-gnu, built 2021/04/26 16:11:58, go1.13.9) (same version as client)
# Cluster ID: 7f43a490-ccd6-4c2a-9534-21f393ca80ce
#
# Enter \? for a brief introduction.
#
root@:26257/defaultdb> CREATE DATABASE mydb;
CREATE DATABASE

Time: 22.985478ms

root@:26257/defaultdb> CREATE USER totoro;
CREATE ROLE

Time: 13.921659ms

root@:26257/defaultdb> GRANT ALL ON DATABASE mydb TO totoro;
GRANT

Time: 14.217559ms

root@:26257/defaultdb> quit
oliver@hki:~$
```

#### Meet the example application

Now that you have started and configured the database engine, you can switch your attention to the application.

The example application for this module is an extended version of `docker-gs-ping` application you've used in the previous modules. You have two options:

- You can update your local copy of `docker-gs-ping` to match the new extended version presented in this chapter; or
- You can clone the [docker/docker-gs-ping-dev](https://github.com/docker/docker-gs-ping-dev) repository. This latter approach is recommended.

To checkout the example application, run:

```console
$ git clone https://github.com/docker/docker-gs-ping-dev.git
# ... output omitted ...
```

The application's `main.go` now includes database initialization code, as well as the code to implement a new business requirement:

- An HTTP `POST` request to `/send` containing a `{ "value" : string }` JSON must save the value to the database.

You also have an update for another business requirement. The requirement was:

- The application responds with a text message containing a heart symbol ("`<3`") on requests to `/`.

And now it's going to be:

- The application responds with the string containing the count of messages stored in the database, enclosed in the parentheses.

  Example output: `Hello, Docker! (7)`

The full source code listing of `main.go` follows.

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := initStore()
	if err != nil {
		log.Fatalf("failed to initialize the store: %s", err)
	}
	defer db.Close()

	e.GET("/", func(c echo.Context) error {
		return rootHandler(db, c)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/send", func(c echo.Context) error {
		return sendHandler(db, c)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

type Message struct {
	Value string `json:"value"`
}

func initStore() (*sql.DB, error) {

	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	var (
		db  *sql.DB
		err error
	)
	openDB := func() error {
		db, err = sql.Open("postgres", pgConnString)
		return err
	}

	err = backoff.Retry(openDB, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS message (value TEXT PRIMARY KEY)"); err != nil {
		return nil, err
	}

	return db, nil
}

func rootHandler(db *sql.DB, c echo.Context) error {
	r, err := countRecords(db)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("Hello, Docker! (%d)\n", r))
}

func sendHandler(db *sql.DB, c echo.Context) error {

	m := &Message{}

	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := crdb.ExecuteTx(context.Background(), db, nil,
		func(tx *sql.Tx) error {
			_, err := tx.Exec(
				"INSERT INTO message (value) VALUES ($1) ON CONFLICT (value) DO UPDATE SET value = excluded.value",
				m.Value,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			return nil
		})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}

func countRecords(db *sql.DB) (int, error) {

	rows, err := db.Query("SELECT COUNT(*) FROM message")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		rows.Close()
	}

	return count, nil
}
```

The repository also includes the `Dockerfile`, which is almost exactly the same as the multi-stage `Dockerfile` introduced in the previous modules. It uses the official Docker Go image to build the application and then builds the final image by placing the compiled binary into the much slimmer, distroless image.

Regardless of whether you had updated the old example application, or checked out the new one, this new Docker image has to be built to reflect the changes to the application source code.

#### Build the application

You can build the image with the familiar `build` command:

```console
$ docker build --tag docker-gs-ping-roach .
```

#### Run the application

Now, run your container. This time you'll need to set some environment variables so that your application knows how to access the database. For now, you’ll do this right in the `docker run` command. Later you'll see a more convenient method with Docker Compose.

> [!NOTE]
>
> Since you're running your CockroachDB cluster in insecure mode, the value for the password can be anything.
>
> In production, don't run in insecure mode.

```console
$ docker run -it --rm -d \
  --network mynet \
  --name rest-server \
  -p 80:8080 \
  -e PGUSER=totoro \
  -e PGPASSWORD=myfriend \
  -e PGHOST=db \
  -e PGPORT=26257 \
  -e PGDATABASE=mydb \
  docker-gs-ping-roach
```

There are a few points to note about this command.

- You map container port `8080` to host port `80` this time. Thus, for `GET` requests you can get away with literally `curl localhost`:

  ```console
  $ curl localhost
  Hello, Docker! (0)
  ```

  Or, if you prefer, a proper URL would work just as well:

  ```console
  $ curl http://localhost/
  Hello, Docker! (0)
  ```

- The total number of stored messages is `0` for now. This is fine, because you haven't posted anything to your application yet.
- You refer to the database container by its hostname, which is `db`. This is why you had `--hostname db` when you started the database container.

- The actual password doesn't matter, but it must be set to something to avoid confusing the example application.
- The container you've just run is named `rest-server`. These names are useful for managing the container lifecycle:

  ```console
  # Don't do this just yet, it's only an example:
  $ docker container rm --force rest-server
  ```

#### Test the application

In the previous section, you've already tested querying your application with `GET` and it returned zero for the stored message counter. Now, post some messages to it:

```console
$ curl --request POST \
  --url http://localhost/send \
  --header 'content-type: application/json' \
  --data '{"value": "Hello, Docker!"}'
```

The application responds with the contents of the message, which means it has been saved in the database:

```json
{ "value": "Hello, Docker!" }
```

Send another message:

```console
$ curl --request POST \
  --url http://localhost/send \
  --header 'content-type: application/json' \
  --data '{"value": "Hello, Oliver!"}'
```

And again, you get the value of the message back:

```json
{ "value": "Hello, Oliver!" }
```

Run curl and see what the message counter says:

```console
$ curl localhost
Hello, Docker! (2)
```

In this example, you sent two messages and the database kept them. Or has it? Stop and remove all your containers, but not the volumes, and try again.

First, stop the containers:

```console
$ docker container stop rest-server roach
rest-server
roach
```

Then, remove them:

```console
$ docker container rm rest-server roach
rest-server
roach
```

Verify that they're gone:

```console
$ docker container list --all
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
```

And start them again, database first:

```console
$ docker run -d \
  --name roach \
  --hostname db \
  --network mynet \
  -p 26257:26257 \
  -p 8080:8080 \
  -v roach:/cockroach/cockroach-data \
  cockroachdb/cockroach:latest-v25.4 start-single-node \
  --insecure
```

And the service next:

```console
$ docker run -it --rm -d \
  --network mynet \
  --name rest-server \
  -p 80:8080 \
  -e PGUSER=totoro \
  -e PGPASSWORD=myfriend \
  -e PGHOST=db \
  -e PGPORT=26257 \
  -e PGDATABASE=mydb \
  docker-gs-ping-roach
```

Lastly, query your service:

```console
$ curl localhost
Hello, Docker! (2)
```

Great! The count of records from the database is correct although you haven't only stopped the containers, but you've also removed them before starting new instances. The difference is in the managed volume for CockroachDB, which you reused. The new CockroachDB container has read the database files from the disk, just as it normally would if it were running outside the container.

#### Wind down everything

Remember, that you're running CockroachDB in insecure mode. Now that you've built and tested your application, it's time to wind everything down before moving on. You can list the containers that you are running with the `list` command:

```console
$ docker container list
```

Now that you know the container IDs, you can use `docker container stop` and `docker container rm`, as demonstrated in the previous modules.

Stop the CockroachDB and `docker-gs-ping-roach` containers before moving on.

### Better productivity with Docker Compose

At this point, you might be wondering if there is a way to avoid having to deal with long lists of arguments to the `docker` command. The toy example you used in this series requires five environment variables to define the connection to the database. A real application might need many, many more. Then there is also a question of dependencies. Ideally, you want to make sure that the database is started before your application is run. And spinning up the database instance may require another Docker command with many options. But there is a better way to orchestrate these deployments for local development purposes.

In this section, you'll create a Docker Compose file to start your `docker-gs-ping-roach` application and CockroachDB database engine with a single command.

#### Configure Docker Compose

In your application's directory, create a new text file named `compose.yaml` with the following content.

```yaml
services:
  docker-gs-ping-roach:
    depends_on:
      - roach
    build:
      context: .
    container_name: rest-server
    hostname: rest-server
    networks:
      - mynet
    ports:
      - 80:8080
    environment:
      - PGUSER=${PGUSER:-totoro}
      - PGPASSWORD=${PGPASSWORD:?database password not set}
      - PGHOST=${PGHOST:-db}
      - PGPORT=${PGPORT:-26257}
      - PGDATABASE=${PGDATABASE:-mydb}
    deploy:
      restart_policy:
        condition: on-failure
  roach:
    image: cockroachdb/cockroach:latest-v25.4
    container_name: roach
    hostname: db
    networks:
      - mynet
    ports:
      - 26257:26257
      - 8080:8080
    volumes:
      - roach:/cockroach/cockroach-data
    command: start-single-node --insecure

volumes:
  roach:

networks:
  mynet:
    driver: bridge
```

This Docker Compose configuration is super convenient as you don't have to type all the parameters to pass to the `docker run` command. You can declaratively do that in the Docker Compose file. The [Docker Compose documentation pages](/manuals/compose/_index.md) are quite extensive and include a full reference for the Docker Compose file format.

#### The `.env` file

Docker Compose will automatically read environment variables from a `.env` file if it's available. Since your Compose file requires `PGPASSWORD` to be set, add the following content to the `.env` file:

```bash
PGPASSWORD=whatever
```

The exact value doesn't really matter for this example, because you run CockroachDB in insecure mode. Make sure you set the variable to some value to avoid getting an error.

#### Merging Compose files

The filename `compose.yaml` is the default filename which `docker compose` command recognizes if no `-f` flag is provided. This means you can have multiple Docker Compose files if your environment has such requirements. Furthermore, Docker Compose files are composable, so multiple files can be specified on the command line to merge parts of the configuration together. The following list shows a few examples of scenarios where such a feature would be useful:

- Using a bind mount for the source code for local development but not when running the CI tests;
- Switching between using a pre-built image for the frontend for some API application vs creating a bind mount for source code;
- Adding additional services for integration testing;
- And many more...

You aren't going to cover any of these advanced use cases here.

#### Variable substitution in Docker Compose

One of the really cool features of Docker Compose is [variable substitution](/reference/compose-file/interpolation.md). You can see some examples in the Compose file, `environment` section. By means of an example:

- `PGUSER=${PGUSER:-totoro}` means that inside the container, the environment variable `PGUSER` shall be set to the same value as it has on the host machine where Docker Compose is run. If there is no environment variable with this name on the host machine, the variable inside the container gets the default value of `totoro`.
- `PGPASSWORD=${PGPASSWORD:?database password not set}` means that if the environment variable `PGPASSWORD` isn't set on the host, Docker Compose will display an error. This is OK, because you don't want to hard-code default values for the password. You set the password value in the `.env` file, which is local to your machine. It is always a good idea to add `.env` to `.gitignore` to prevent the secrets being checked into the version control.

Other ways of dealing with undefined or empty values exist, as documented in the [variable substitution](/reference/compose-file/interpolation.md) section of the Docker documentation.

#### Validating Docker Compose configuration

Before you apply changes made to a Compose configuration file, there is an opportunity to validate the content of the configuration file with the following command:

```console
$ docker compose config
```

When this command is run, Docker Compose reads the file `compose.yaml`, parses it into a data structure in memory, validates where possible, and prints back the reconstruction of that configuration file from its internal representation. If this isn't possible due to errors, Docker prints an error message instead.

#### Build and run the application using Docker Compose

Start your application and confirm that it's running.

```console
$ docker compose up --build
```

You passed the `--build` flag so Docker will compile your image and then start it.

> [!NOTE]
>
> Docker Compose is a useful tool, but it has its own quirks. For example, no rebuild is triggered on the update to the source code unless the `--build` flag is provided. It is a very common pitfall to edit one's source code, and forget to use the `--build` flag when running `docker compose up`.

Since your set-up is now run by Docker Compose, it has assigned it a project name, so you get a new volume for your CockroachDB instance. This means that your application will fail to connect to the database, because the database doesn't exist in this new volume. The terminal displays an authentication error for the database:

```text
# ... omitted output ...
rest-server             | 2021/05/10 00:54:25 failed to initialise the store: pq: password authentication failed for user totoro
roach                   | *
roach                   | * INFO: Replication was disabled for this cluster.
roach                   | * When/if adding nodes in the future, update zone configurations to increase the replication factor.
roach                   | *
roach                   | CockroachDB node starting at 2021-05-10 00:54:26.398177 +0000 UTC (took 3.0s)
roach                   | build:               CCL v20.1.15 @ 2021/04/26 16:11:58 (go1.13.9)
roach                   | webui:               http://db:8080
roach                   | sql:                 postgresql://root@db:26257?sslmode=disable
roach                   | RPC client flags:    /cockroach/cockroach <client cmd> --host=db:26257 --insecure
roach                   | logs:                /cockroach/cockroach-data/logs
roach                   | temp dir:            /cockroach/cockroach-data/cockroach-temp349434348
roach                   | external I/O path:   /cockroach/cockroach-data/extern
roach                   | store[0]:            path=/cockroach/cockroach-data
roach                   | storage engine:      rocksdb
roach                   | status:              initialized new cluster
roach                   | clusterID:           b7b1cb93-558f-4058-b77e-8a4ddb329a88
roach                   | nodeID:              1
rest-server exited with code 0
rest-server             | 2021/05/10 00:54:25 failed to initialise the store: pq: password authentication failed for user totoro
rest-server             | 2021/05/10 00:54:26 failed to initialise the store: pq: password authentication failed for user totoro
rest-server             | 2021/05/10 00:54:29 failed to initialise the store: pq: password authentication failed for user totoro
rest-server             | 2021/05/10 00:54:25 failed to initialise the store: pq: password authentication failed for user totoro
rest-server             | 2021/05/10 00:54:26 failed to initialise the store: pq: password authentication failed for user totoro
rest-server             | 2021/05/10 00:54:29 failed to initialise the store: pq: password authentication failed for user totoro
rest-server exited with code 1
# ... omitted output ...
```

Because of the way you set up your deployment using `restart_policy`, the failing container is being restarted every 20 seconds. So, in order to fix the problem, you need to log in to the database engine and create the user. You've done it before in [Configure the database engine](#configure-the-database-engine).

This isn't a big deal. All you have to do is to connect to CockroachDB instance and run the three SQL commands to create the database and the user, as described in [Configure the database engine](#configure-the-database-engine).

So, log in to the database engine from another terminal:

```console
$ docker exec -it roach ./cockroach sql --insecure
```

And run the same commands as before to create the database `mydb`, the user `totoro`, and to grant that user necessary permissions. Once you do that (and the example application container is automatically restarts), the `rest-service` stops failing and restarting and the console goes quiet.

It would have been possible to connect the volume that you had previously used, but for the purposes of this example it's more trouble than it's worth and it also provided an opportunity to show how to introduce resilience into your deployment via the `restart_policy` Compose file feature.

#### Testing the application

Now, test your API endpoint. In the new terminal, run the following command:

```console
$ curl http://localhost/
```

You should receive the following response:

```json
Hello, Docker! (0)
```

#### Shutting down

To stop the containers started by Docker Compose, press `ctrl+c` in the terminal where you ran `docker compose up`. To remove those containers after they've been stopped, run `docker compose down`.

#### Detached mode

You can run containers started by the `docker compose` command in detached mode, just as you would with the `docker` command, by using the `-d` flag.

To start the stack, defined by the Compose file in detached mode, run:

```console
$ docker compose up --build -d
```

Then, you can use `docker compose stop` to stop the containers and `docker compose down` to remove them.

### Further exploration

You can run `docker compose` to see what other commands are available.

### Wrap up

There are some tangential, yet interesting points that were purposefully not covered in this chapter. For the more adventurous reader, this section offers some pointers for further study.

#### Persistent storage

A managed volume isn't the only way to provide your container with persistent storage. It is highly recommended to get acquainted with available storage options and their use cases, covered in [Manage data in Docker](/manuals/engine/storage/_index.md).

#### CockroachDB clusters

You ran a single instance of CockroachDB, which was enough for this example. But, it's possible to run a CockroachDB cluster, which is made of multiple instances of CockroachDB, each instance running in its own container. Since CockroachDB engine is distributed by design, it would have taken surprisingly little change to your procedure to run a cluster with multiple nodes.

Such distributed set-up offers interesting possibilities, such as applying Chaos Engineering techniques to simulate parts of the cluster failing and evaluating your application's ability to cope with such failures.

If you are interested in experimenting with CockroachDB clusters, check out:

- [Start a CockroachDB Cluster in Docker](https://www.cockroachlabs.com/docs/v20.2/start-a-local-cluster-in-docker-mac.html) article; and
- Documentation for Docker Compose keywords [`deploy`](/reference/compose-file/legacy-versions.md) and [`replicas`](/reference/compose-file/legacy-versions.md).

#### Other databases

Since you didn't run a cluster of CockroachDB instances, you might be wondering whether you could have used a non-distributed database engine. The answer is 'yes', and if you were to pick a more traditional SQL database, such as [PostgreSQL](https://www.postgresql.org/), the process described in this chapter would have been very similar.

### Next steps

In this module, you set up a containerized development environment with your application and the database engine running in different containers. You also wrote a Docker Compose file which links the two containers together and provides for easy starting up and tearing down of the development environment.

In the next module, you'll take a look at one possible approach to running functional tests in Docker.

## Run your tests using Go test

### Prerequisites

Complete the [Build your Go image](./) section of this guide.

### Overview

Testing is an essential part of modern software development. Testing can mean a
lot of things to different development teams. There are unit tests, integration
tests and end-to-end testing. In this guide you take a look at running your unit
tests in Docker when building.

For this section, use the `docker-gs-ping` project that you cloned in [Build
your Go image](./).

### Run tests when building

To run your tests when building, you need to add a test stage to the
`Dockerfile.multistage`. The `Dockerfile.multistage` in the sample application's
repository already has the following content:

```dockerfile {hl_lines="15-17"}
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

Run the following command to build an image using the `run-test-stage` stage as the target and view the test results. Include `--progress plain` to view the build output, `--no-cache` to ensure the tests always run, and `--target run-test-stage` to target the test stage.

```console
$ docker build -f Dockerfile.multistage -t docker-gs-ping-test --progress plain --no-cache --target run-test-stage .
```

You should see output containing the following.

```text
#13 [run-test-stage 1/1] RUN go test -v ./...
#13 4.915 === RUN   TestIntMinBasic
#13 4.915 --- PASS: TestIntMinBasic (0.00s)
#13 4.915 === RUN   TestIntMinTableDriven
#13 4.915 === RUN   TestIntMinTableDriven/0,1
#13 4.915 === RUN   TestIntMinTableDriven/1,0
#13 4.915 === RUN   TestIntMinTableDriven/2,-2
#13 4.915 === RUN   TestIntMinTableDriven/0,-1
#13 4.915 === RUN   TestIntMinTableDriven/-1,0
#13 4.915 --- PASS: TestIntMinTableDriven (0.00s)
#13 4.915     --- PASS: TestIntMinTableDriven/0,1 (0.00s)
#13 4.915     --- PASS: TestIntMinTableDriven/1,0 (0.00s)
#13 4.915     --- PASS: TestIntMinTableDriven/2,-2 (0.00s)
#13 4.915     --- PASS: TestIntMinTableDriven/0,-1 (0.00s)
#13 4.915     --- PASS: TestIntMinTableDriven/-1,0 (0.00s)
#13 4.915 PASS
```
