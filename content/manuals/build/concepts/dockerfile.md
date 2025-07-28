---
title: Dockerfile overview
weight: 20
description: Learn how to use Dockerfiles to build and package your software into Docker images.
keywords: dockerfile, docker build, buildx, buildkit, container image, getting started, image layers, dockerfile, instructions
aliases:
- /build/hellobuild/
- /build/building/packaging/
---

## Dockerfile

Docker builds images by reading instructions from a Dockerfile. A
Dockerfile is a text file that contains instructions for building your source
code. The Dockerfile instruction syntax is defined in the
[Dockerfile reference](/reference/dockerfile.md).

Here are the most common types of instructions:

| Instruction                                               | Description                                                                                                                            |
|-----------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|
| [`FROM <image>`](/reference/dockerfile.md#from)           | Defines a base for your image.                                                                                                         |
| [`RUN <command>`](/reference/dockerfile.md#run)           | Executes commands in a new layer on top of the current image and commits the result. `RUN` also has a shell form for running commands. |
| [`WORKDIR <directory>`](/reference/dockerfile.md#workdir) | Sets the working directory for any `RUN`, `CMD`, `ENTRYPOINT`, `COPY`, and `ADD` instructions that follow it in the Dockerfile.        |
| [`COPY <src> <dest>`](/reference/dockerfile.md#copy)      | Copies new files or directories from `<src>` and adds them to the container at the path `<dest>`.                                      |
| [`CMD <command>`](/reference/dockerfile.md#cmd)           | Defines the default program that runs when you start the container. Only the last `CMD` in the Dockerfile is used if multiple exist.   |

Dockerfiles are crucial inputs for image builds and can facilitate automated,
multi-layer image builds based on your unique configurations. Dockerfiles can
start simple and grow with your needs to support more complex scenarios.

### Filename

The default filename for a Dockerfile is `Dockerfile`, without a file extension. Using
the default name lets you run the `docker build` command without extra flags.

Some projects may need distinct Dockerfiles for specific purposes. A common
convention is to name these `<something>.Dockerfile`. You can specify the
Dockerfile filename using the `--file` flag with the `docker build` command. See the
[`docker build` CLI reference](/reference/cli/docker/buildx/build.md#file) for details.

> [!NOTE]
> We recommend using the default (`Dockerfile`) for your project's main Dockerfile.

## Docker images

Docker images consist of layers. Each layer is the result of a build
instruction in the Dockerfile. Layers are stacked sequentially, and each one is
a delta representing the changes applied to the previous layer.

### Example

A typical workflow for building applications with Docker:

The following example shows a small "Hello World" application in Python using Flask.

```python
from flask import Flask
app = Flask(__name__)

@app.route("/")
def hello():
    return "Hello World!"
```

Without Docker Build, you need to:

- Install the required runtime dependencies on the server.
- Upload the Python code to the server's filesystem.
- Start your application on the server with the necessary parameters.

The following Dockerfile creates a container image with all dependencies installed and
automatically starts your application.

```dockerfile
# syntax=docker/dockerfile:1
FROM ubuntu:22.04

# install app dependencies
RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip install flask==3.0.*

# install app
COPY hello.py /

# final configuration
ENV FLASK_APP=hello
EXPOSE 8000
CMD ["flask", "run", "--host", "0.0.0.0", "--port", "8000"]
```

This Dockerfile does the following:

- [Dockerfile syntax](#dockerfile-syntax)
- [Base image](#base-image)
- [Environment setup](#environment-setup)
- [Comments](#comments)
- [Installing dependencies](#installing-dependencies)
- [Copying files](#copying-files)
- [Setting environment variables](#setting-environment-variables)
- [Exposed ports](#exposed-ports)
- [Starting the application](#starting-the-application)

### Dockerfile syntax

The first line is a [`# syntax` parser directive](/reference/dockerfile.md#syntax).
This optional directive tells Docker which syntax to use when parsing the Dockerfile.
It lets older Docker versions with [BuildKit enabled](../buildkit/_index.md#getting-started)
use a specific [Dockerfile frontend](../buildkit/frontend.md) before starting the
build. [Parser directives](/reference/dockerfile.md#parser-directives) must appear
before any other comment, whitespace, or instruction, and should be the first line.

```dockerfile
# syntax=docker/dockerfile:1
```

> [!TIP]
> Use `docker/dockerfile:1` to always get the latest version 1 syntax. BuildKit
> checks for updates before building, so you use the most current version.

### Base image

The next line defines the base image:

```dockerfile
FROM ubuntu:22.04
```

The [`FROM` instruction](/reference/dockerfile.md#from) sets your base
image to the 22.04 release of Ubuntu. All following instructions run in this Ubuntu environment. The
`ubuntu:22.04` notation follows the `name:tag` standard for Docker images. You can
use many public images in your projects by importing them with the `FROM` instruction.

[Docker Hub](https://hub.docker.com/search?image_filter=official&q=&type=image)
offers many official images you can use.

### Environment setup

This line runs a build command inside the base image.

```dockerfile
# install app dependencies
RUN apt-get update && apt-get install -y python3 python3-pip
```

This [`RUN` instruction](/reference/dockerfile.md#run) executes a
shell in Ubuntu that updates the APT package index and installs Python tools in
the container.

### Comments

Note the `# install app dependencies` line. This is a comment. Comments in
Dockerfiles begin with the `#` symbol. As your Dockerfile evolves, comments can
be instrumental to document how your Dockerfile works for any future readers
and editors of the file, including your future self!

### Installing dependencies

The second `RUN` instruction installs the `flask` dependency for the Python app.

```dockerfile
RUN pip install flask==3.0.*
```

A prerequisite for this instruction is that `pip` is installed into the build
container. The first `RUN` command installs `pip`, which ensures that we can
use the command to install the flask web framework.

### Copying files

The next instruction uses the
[`COPY` instruction](/reference/dockerfile.md#copy) to copy the
`hello.py` file from the local build context into the root directory of our image.

```dockerfile
COPY hello.py /
```

A [build context](./context.md) is the set of files that you can access
in Dockerfile instructions such as `COPY` and `ADD`.

After the `COPY` instruction, the `hello.py` file is added to the filesystem
of the build container.

### Setting environment variables

If your application uses environment variables, you can set environment variables
in your Docker build using the [`ENV` instruction](/reference/dockerfile.md#env).

```dockerfile
ENV FLASK_APP=hello
```

This sets a Linux environment variable needed by Flask to start the app. Without this,
Flask cannot find the app to run it.

### Exposed ports

The [`EXPOSE` instruction](/reference/dockerfile.md#expose) marks that
our final image has a service listening on port `8000`.

```dockerfile
EXPOSE 8000
```

This instruction isn't required, but it is a good practice and helps tools and
team members understand what this application is doing.

### Starting the application

Finally, [`CMD` instruction](/reference/dockerfile.md#cmd) sets the
command that is run when the user starts a container based on this image.

```dockerfile
CMD ["flask", "run", "--host", "0.0.0.0", "--port", "8000"]
```

This command starts the flask development server listening on all addresses
on port `8000`. The example here uses the "exec form" version of `CMD`.
It's also possible to use the "shell form":

```dockerfile
CMD flask run --host 0.0.0.0 --port 8000
```

There are subtle differences between these two versions,
for example in how they trap signals like `SIGTERM` and `SIGKILL`.
For more information about these differences, see
[Shell and exec form](/reference/dockerfile.md#shell-and-exec-form)

## Building

To build a container image using the Dockerfile example from the
[previous section](#example), you use the `docker build` command:

```console
$ docker build -t test:latest .
```

The `-t test:latest` option specifies the name and tag of the image.

The single dot (`.`) at the end of the command sets the
[build context](./context.md) to the current directory. This means that the
build expects to find the Dockerfile and the `hello.py` file in the directory
where the command is invoked. If those files aren't there, the build fails.

After the image has been built, you can run the application as a container with
`docker run`, specifying the image name:

```console
$ docker run -p 127.0.0.1:8000:8000 test:latest
```

This publishes the container's port 8000 to `http://localhost:8000` on the
Docker host.

> [!TIP]
>
> Want a better editing experience for Dockerfiles in VS Code?
> Check out the [Docker VS Code Extension (Beta)](https://marketplace.visualstudio.com/items?itemName=docker.docker) for linting, code navigation, and vulnerability scanning.
