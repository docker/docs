---
title: C++ language-specific guide
linkTitle: C++
description: Containerize and develop C++ applications using Docker.
keywords: getting started, c++
summary: |
  This guide explains how to containerize C++ applications using Docker.
aliases:
  - /language/cpp/
  - /guides/language/cpp/
  - /language/cpp/containerize/
  - /language/cpp/develop/
  - /language/cpp/configure-ci-cd/
  - /language/cpp/deploy/
  - /guides/cpp/configure-ci-cd/
  - /guides/cpp/containerize/
  - /guides/cpp/deploy/
  - /guides/cpp/develop/
  - /guides/cpp/multistage/
  - /guides/cpp/security/
params:
  tags: [languages]
  time: 20 minutes
---

The C++ getting started guide teaches you how to create a containerized C++ application using Docker. In this guide, you'll learn how to:

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) and [Mohammad-Ali A'râbi](https://twitter.com/MohammadAliEN) for their contribution to this guide.

- Containerize and run a C++ application using a multi-stage Docker build
- Build and run a C++ application using Docker Compose
- Set up a local environment to develop a C++ application using containers

After completing the C++ getting started modules, you should be able to containerize your own C++ application based on the examples and instructions provided in this guide.

Start by containerizing an existing C++ application.

## Create a multi-stage build for your C++ application

### Prerequisites

- You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

### Overview

This section walks you through creating a multi-stage Docker build for a C++ application.
A multi-stage build is a Docker feature that allows you to use different base images for different stages of the build process,
so you can optimize the size of your final image and separate build dependencies from runtime dependencies.

The standard practice for compiled languages like C++ is to have a build stage that compiles the code and a runtime stage that runs the compiled binary,
because the build dependencies are not needed at runtime.

### Get the sample application

Let's use a simple C++ application that prints `Hello, World!` to the terminal. To do so, clone the sample repository to use with this guide:

```bash
$ git clone https://github.com/dockersamples/c-plus-plus-docker.git
```

The example for this section is under the `hello` directory in the repository. Get inside it and take a look at the files:

```bash
$ cd c-plus-plus-docker/hello
$ ls
```

You should see the following files:

```text
Dockerfile  hello.cpp
```

### Check the Dockerfile

Open the `Dockerfile` in an IDE or text editor. The `Dockerfile` contains the instructions for building the Docker image.

```Dockerfile
# Stage 1: Build stage
FROM ubuntu:latest AS build

# Install build-essential for compiling C++ code
RUN apt-get update && apt-get install -y build-essential

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY hello.cpp .

# Compile the C++ code statically to ensure it doesn't depend on runtime libraries
RUN g++ -o hello hello.cpp -static

# Stage 2: Runtime stage
FROM scratch

# Copy the static binary from the build stage
COPY --from=build /app/hello /hello

# Command to run the binary
CMD ["/hello"]
```

The `Dockerfile` has two stages:

1. **Build stage**: This stage uses the `ubuntu:latest` image to compile the C++ code and create a static binary.
2. **Runtime stage**: This stage uses the `scratch` image, which is an empty image, to copy the static binary from the build stage and run it.

### Build the Docker image

To build the Docker image, run the following command in the `hello` directory:

```bash
$ docker build -t hello .
```

The `-t` flag tags the image with the name `hello`.

### Run the Docker container

To run the Docker container, use the following command:

```bash
$ docker run hello
```

You should see the output `Hello, World!` in the terminal.

### Summary

In this section, you learned how to create a multi-stage build for a C++ application. Multi-stage builds help you optimize the size of your final image and separate build dependencies from runtime dependencies.
In this example, the final image only contains the static binary and doesn't include any build dependencies.

As the image has an empty base, the usual OS tools are also absent. So, for example, you can't run a simple `ls` command in the container:

```bash
$ docker run hello ls
```

This makes the image very lightweight and secure.

## Containerize a C++ application

### Prerequisites

- You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

### Overview

This section walks you through containerizing and running a C++ application, using Docker Compose.

### Get the sample application

We're using the same sample repository that you used in the previous sections of this guide. If you haven't already cloned the repository, clone it now:

```console
$ git clone https://github.com/dockersamples/c-plus-plus-docker.git
```

You should now have the following contents in your `c-plus-plus-docker` (root)
directory.

```text
├── c-plus-plus-docker/
│ ├── compose.yml
│ ├── Dockerfile
│ ├── LICENSE
│ ├── ok_api.cpp
│ └── README.md

```

To learn more about the files in the repository, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yml](/reference/compose-file/_index.md)

### Run the application

Inside the `c-plus-plus-docker` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You will see a message `{"Status" : "OK"}` in the browser.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `c-plus-plus-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080).

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

### Summary

In this section, you learned how you can containerize and run your C++
application using Docker.

Related information:

- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll learn how you can develop your application using
containers.

## Use containers for C++ development

### Prerequisites

Complete [Containerize a C++ application](./).

### Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Configuring Compose to automatically update your running Compose services as you edit and save your code

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/dockersamples/c-plus-plus-docker.git && cd c-plus-plus-docker
```

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Open your `compose.yml` file in an IDE or text editor and then add the Compose Watch instructions. The following example shows how to add Compose Watch to your `compose.yml` file.

```yaml {hl_lines="11-14",linenos=true}
services:
  ok-api:
    image: ok-api
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    develop:
      watch:
        - action: rebuild
          path: .
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Now, if you modify your `ok_api.cpp` you will see the changes in real time without re-building the image.

To test it out, open the `ok_api.cpp` file in your favorite text editor and change the message from `{"Status" : "OK"}` to `{"Status" : "Updated"}`. Save the file and refresh your browser at [http://localhost:8080](http://localhost:8080). You should see the updated message.

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose file watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
