---
title: "Build your Go image"
keywords: containers, images, go, golang, dockerfiles, coding, build, push, run
description: Learn how to build your first Docker image by writing a Dockerfile
redirect_from:
- /get-started/golang/build-images/
---

{% include_relative nav.html selected="1" %}

## Prerequisites

Work through the orientation and setup in Get started [Part 1](/get-started/) to understand Docker concepts.

## Overview

Now that we have a good overview of containers and the Docker platform, let’s take a look at building our first image. An image includes everything you need to run an application - the code or binary, runtime, dependencies, and any other file system objects required.

To complete this tutorial, you need the following:

- Go version 1.13 or later: [Download Go](https://golang.org/dl/){: target="_blank" rel="noopener" class="_"}.
- Docker running locally: Follow the instructions to [download and install Docker](https://docs.docker.com/desktop/).
- An IDE or a text editor to edit files. We recommend using [Visual Studio Code](https://code.visualstudio.com/){: target="_blank" rel="noopener" class="_"}.

## Sample application

Let’s create a simple Go application that we can use as our example &ndash; a (naive) key-value store with two operations: `add` and `list`, accessible over a REST API end-point.

If you are in a hurry, the source code is in the [go-docker](https://github.com/olliefr/go-docker) repo. 

Alternatively, to recreate the application from scratch, create a directory on your local machine named `go-docker` and follow the steps below to initialise the Go project.

```shell
$ cd [path to your go-docker directory]
$ go mod init go-docker
$ go get github.com/labstack/echo/v4
$ go get github.com/labstack/echo/v4/middleware
$ touch main.go
```

Now, let’s add some code to handle our REST requests. We’ll use the [echo](https://echo.labstack.com/){: target="_blank" rel="noopener" class="_"} framework to quickly build a minimal server so we can focus on Dockerizing the application.

Open this working directory in your IDE and add the following code into the `main.go` file.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Entry is a single key-value pair
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// store is an ordered list of key-value pairs
var store = make([]*Entry, 0)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/list", listEntries)
	e.POST("/add", addEntry)

	e.Logger.Fatal(e.Start(":8000"))
}

// listEntries returns a full copy of the store in JSON format to the client.
func listEntries(c echo.Context) error {
	return c.JSON(http.StatusOK, store)
}

// addEntry adds a new entry to the key-value store
// and returns the value in JSON format to the client.
func addEntry(c echo.Context) error {
	e := new(Entry)
	if err := c.Bind(e); err != nil {
		return err
	}
	store = append(store, e)
	return c.JSON(http.StatusOK, e)
}
```

The server will listen on port 8000. It responds to three kinds of query:

* GET requests to the root (`/`) return server status in JSON format.
* POST requests to `/add` store a key-value pair on the server.
* GET requests to `/list` return an array of JSON objects that you have previously POSTed.

The key-value pairs POSTed to the server must be in (valid) JSON format.

## Test application

Let’s start our application and make sure it’s running properly. Open your terminal and navigate to your working directory you created.

```shell
$ go run main.go
```

This should compile and start the server as a foreground application, outputting the banner and some logging information, as illustrated in the next figure.

```shell
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.1.17
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8000
```

To test that the application is working properly, let's perform a _smoke test_ first. In a new terminal run the following command.

```shell
$ curl http://localhost:8000/
```

The server should respond with `{"Status":"OK"}` in the same terminal. 

Having established that the server is running and is accessible, we’ll POST some JSON to the API and then make a GET request to see that the data has been saved.

Run the following `curl` command in the terminal.

```shell
$ curl --request POST \
  --url http://localhost:8000/add \
  --header 'content-type: application/json' \
  --data '{"key": "name", "value": "Docker"}'
```

With this command, we have sent the following JSON structure to the server.

```json
{
	"key": "name",
	"value": "Docker"
}
```

On success, the output from the server should read just that.

```json
{"key":"name","value":"Docker"}
```

To see what data has been saved on the server, run the following command.

```shell
$ curl http://localhost:8000/list
```

This should return a JSON list of length one, comprising of the element we've just added to it.

```json
[{"key":"name","value":"Docker"}]
```

Switch back to the terminal where our server is running. You should now see thee requests in the server logs. The exact format might be different from the following, and the timestamps will differ too.

```json
{"time":"2021-02-08T23:49:16.0027265+02:00","id":"","remote_ip":"127.0.0.1","host":"localhost:8000","method":"GET","uri":"/","user_agent":"curl/7.68.0","status":200,"error":"","latency":28300,"latency_human":"28.3µs","bytes_in":0,"bytes_out":16}
{"time":"2021-02-08T23:49:25.8262003+02:00","id":"","remote_ip":"127.0.0.1","host":"localhost:8000","method":"POST","uri":"/add","user_agent":"curl/7.68.0","status":200,"error":"","latency":85500,"latency_human":"85.5µs","bytes_in":34,"bytes_out":32}
{"time":"2021-02-08T23:49:29.695339+02:00","id":"","remote_ip":"127.0.0.1","host":"localhost:8000","method":"GET","uri":"/list","user_agent":"curl/7.68.0","status":200,"error":"","latency":27700,"latency_human":"27.7µs","bytes_in":0,"bytes_out":34}
```

You can add more key-value pairs to the store, if you'd like. Once we are satisfied, that the store "works", we can proceed to an exciting part of running this application in Docker.

## Create a Dockerfile for Go

A Dockerfile is a text document that contains all the commands a user could call on the command line to assemble an image. When we tell Docker to build our image by executing the `docker build` command, Docker reads these instructions and executes them one by one and creates a Docker image as a result.

Let’s walk through the process of creating a Dockerfile for our application. In the root of your working directory, create a file named `Dockerfile` and open this file in your text editor.

> **Note**
>
> The name of the Dockerfile is not important but the default filename for many commands is simply `Dockerfile`. So, we’ll use that as our filename throughout this series.

The first thing we need to do is to add a line in our Dockerfile that tells Docker what base image we would like to use for our application.

```dockerfile
FROM golang:1.15-alpine
```

Docker images can be inherited from other images. Therefore, instead of creating our own base image, we’ll use the official Go image that already has all the tools and packages to compile and run a Go application. You can think of this in the same way you would think about class inheritance in object oriented programming. For example, if we were able to create Docker images in JavaScript, we might write something like the following.

`class MyImage extends NodeBaseImage {}`

This would create a class called `MyImage` that inherited functionality from the base class `NodeBaseImage`.

In the same way, when we use the `FROM` command, we tell Docker to include in our image all the functionality from the `golang:1.15-alpine` image.

> **Note**
>
> If you want to learn more about creating your own base images, see [Creating base images](https://docs.docker.com/develop/develop-images/baseimages/).

<!-- TODO something about env variables?
The `NODE_ENV` environment variable specifies the environment in which an application is running (usually, development or production). One of the simplest things you can do to improve performance is to set `NODE_ENV` to `production`.

```dockerfile
ENV NODE_ENV=production
```
-->

To make things easier when running the rest of our commands, let’s create a working directory. This instructs Docker to use this path as the default location for all subsequent commands. This way we do not have to type out full file paths but can use relative paths based on the working directory.

```dockerfile
WORKDIR /app
```

Usually the very first thing you do once you’ve downloaded a project written in Go is to install the modules necessary to compile it.

Before we can run `go mod download`, we need to get our `go.mod` and `go.sum` files into our images. We use the `COPY` command to do this. The  `COPY` command takes two parameters. The first parameter tells Docker what file(s) you would like to copy into the image. The second parameter tells Docker where you want that file(s) to be copied to. We’ll copy the `go.mod` and `go.sum` file into our working directory `/app`.

```dockerfile
COPY ["go.mod", "go.sum", "./"]
```

Once we have our module files inside the image, we can use the `RUN` command to execute the command `go mod download`. This works exactly the same as if we were running `go` locally on our machine, but this time these Go modules will be installed into the a directory inside our image.

```dockerfile
RUN go mod download
```

At this point, we have an image that is based on Go environment version 1.13 (or later minor version) and we have installed our dependencies. The next thing we need to do is to add our source code into the image. We’ll use the COPY command just like we did with our module files above.

```dockerfile
COPY . .
```

The COPY command takes all the files located in the current directory and copies them into the image. Now, all we have to do is to tell Docker what command we want to run when our image is run inside of a container. We do this with the `CMD` command.

```dockerfile
CMD [ "go", "run", "main.go" ]
```

Here's the complete Dockerfile.

```dockerfile
FROM golang:1.15-alpine

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .

CMD [ "go", "run", "main.go" ]
```

You can also find this `Dockerfile` for this part in [go-docker](https://github.com/olliefr/go-docker) repo.

## Build image

Now that we’ve created our Dockerfile, let’s build our image. To do this, we use the `docker build` command. The `docker build` command builds Docker images from a Dockerfile and a “context”. A build’s context is the set of files located in the specified PATH or URL. The Docker build process can access any of the files located in the context.

The build command optionally takes a `--tag` flag. The tag is used to set the name of the image and an optional tag in the format `‘name:tag’`. We’ll leave off the optional “tag” for now to help simplify things. If you do not pass a tag, Docker will use “latest” as its default tag. You’ll see this in the last line of the build output.

Let’s build our first Docker image.

```shell
$ docker build --tag go-docker .
```

```shell
[+] Building 15.4s (10/10) FINISHED
 => [internal] load build definition from Dockerfile                               0.0s
 => => transferring dockerfile: 179B                                               0.0s
 => [internal] load .dockerignore                                                  0.0s
 => => transferring context: 2B                                                    0.0s
 => [internal] load metadata for docker.io/library/golang:1.15-alpine              0.0s
 => CACHED [1/5] FROM docker.io/library/golang:1.15-alpine                         0.0s
 => [internal] load build context                                                  0.0s
 => => transferring context: 5.75kB                                                0.0s
 => [2/5] WORKDIR /app                                                             0.0s
 => [3/5] COPY [go.mod, go.sum, ./]                                                0.0s
 => [4/5] RUN go mod download                                                     14.7s
 => [5/5] COPY . .                                                                 0.0s
 => exporting to image                                                             0.5s
 => => exporting layers                                                            0.5s
 => => writing image sha256:94a23bb01d3d00f1866efc1be139785e9c3e83f1363df7dc63d41457b8a06423  0.0s
 => => naming to docker.io/library/go-docker                                       0.0s
```

## Viewing local images

To see a list of images we have on our local machine, we have two options. One is to use the CLI and the other is to use [Docker Desktop](https://www.docker.com/products/docker-desktop){: target="_blank" rel="noopener" class="_"}. Since we are currently working in the terminal, let’s take a look at listing images with the CLI.

To list images, simply run the `images` command.

```shell
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
go-docker           latest              94a23bb01d3d        6 minutes ago       378MB
golang              1.15-alpine         1463476d8605        5 weeks ago         299MB
```

You should see _at least_ two images listed. One for the base image `1.15-alpine` and the other for our image we just build `go-docker:latest`.

## Tag images

An image name is made up of slash-separated name components. Name components may contain lowercase letters, digits and separators. A separator is defined as a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.

An image is made up of a manifest and a list of layers. In simple terms, a “tag” points to a combination of these artifacts. You can have multiple tags for an image. Let’s create a second tag for the image we built and take a look at its layers.

To create a new tag for the image we built above, run the following command.

```shell
$ docker tag go-docker:latest go-docker:v1.0.0
```

The Docker tag command creates a new tag for an image. It does not create a new image. The tag points to the same image and is just another way to reference the image.

Now run the `docker images` command to see a list of our local images.

```
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
go-docker           v1.0.0              94a23bb01d3d        42 minutes ago      378MB
go-docker           latest              94a23bb01d3d        6 minutes ago       378MB
golang              1.15-alpine         1463476d8605        5 weeks ago         299MB
```

You can see that we have two images that start with `go-docker`. We know they are the same image because if you look at the IMAGE ID column, you can see that the values are the same for the two images.

Let’s remove the tag that we just created. To do this, we’ll use the rmi command. The rmi command stands for “remove image”.

```shell
$ docker rmi go-docker:v1.0.0
Untagged: go-docker:v1.0.0
```

Notice that the response from Docker tells us that the image has not been removed but only “untagged”. Verify this by running the images command.

```shell
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
go-docker           latest              94a23bb01d3d        45 minutes ago      378MB
golang              1.15-alpine         1463476d8605        5 weeks ago         299MB
```

Our image that was tagged with `:v1.0.0` has been removed but we still have the `go-docker:latest` tag available on our machine.

## Next steps

In this module, we took a look at setting up our example Go application that we will use for the rest of the tutorial. We also created a Dockerfile that we used to build our Docker image. Then, we took a look at tagging our images and removing images. In the next module, we’ll take a look at how to:

[Run your image as a container](run-containers.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs ](https://github.com/docker/docker.github.io/issues/new?title=[Golang%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
