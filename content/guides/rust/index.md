---
title: Rust language-specific guide
linkTitle: Rust
description: Containerize Rust apps using Docker
keywords: Docker, getting started, Rust, language
summary: |
  This guide covers how to containerize Rust applications using Docker.
aliases:
  - /language/rust/
  - /guides/language/rust/
  - /language/rust/build-images/
  - /language/rust/run-containers/
  - /language/rust/develop/
  - /language/rust/configure-ci-cd/
  - /language/rust/deploy/
  - /guides/rust/build-images/
  - /guides/rust/configure-ci-cd/
  - /guides/rust/deploy/
  - /guides/rust/develop/
  - /guides/rust/run-containers/
params:
  tags: [languages]
  time: 20 minutes
---


The Rust language-specific guide teaches you how to create a containerized Rust application using Docker. In this guide, you'll learn how to:

- Containerize a Rust application
- Build an image and run the newly built image as a container
- Set up volumes and networking
- Orchestrate containers using Compose
- Use containers for development
- Configure a CI/CD pipeline for your application using GitHub Actions
- Deploy your containerized Rust application locally to Kubernetes to test and debug your deployment

After completing the Rust modules, you should be able to containerize your own Rust application based on the examples and instructions provided in this guide.

Start with building your first Rust image.

## Build your Rust image

### Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

### Overview

This guide walks you through building your first Rust image. An image
includes everything needed to run an application - the code or binary, runtime,
dependencies, and any other file system objects required.

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/docker/docker-rust-hello && cd docker-rust-hello
```

### Choose a base image

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for your project. Ask Gordon to create a Dockerfile, Compose file, and `.dockerignore` tailored to your application.

Before editing your Dockerfile, you need to choose a base image. You can use the [Rust Docker Official Image](https://hub.docker.com/_/rust),  
or a [Docker Hardened Image (DHI)](https://hub.docker.com/hardened-images/catalog/dhi/rust).

Docker Hardened Images (DHIs) are minimal, secure, and production-ready base images maintained by Docker.  
They help reduce vulnerabilities and simplify compliance. For more details, see [Docker Hardened Images](/dhi/).

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}

Docker Hardened Images (DHIs) are publicly available and can be used directly as base images.
To pull Docker Hardened Images, authenticate once with Docker:

```bash
docker login dhi.io
```

Use DHIs from the dhi.io registry, for example:

```bash
FROM dhi.io/rust:${RUST_VERSION}-alpine3.22-dev AS build
```

The following Dockerfile uses a Rust DHI as the build base image:

```dockerfile {title=Dockerfile}
# Make sure RUST_VERSION matches the Rust version
ARG RUST_VERSION=1.92
ARG APP_NAME=docker-rust-hello

################################################################################
# Create a stage for building the application.
################################################################################

FROM dhi.io/rust:${RUST_VERSION}-alpine3.22-dev AS build
ARG APP_NAME
WORKDIR /app

# Install host build dependencies.
RUN apk add --no-cache clang lld musl-dev git

# Build the application.
RUN --mount=type=bind,source=src,target=src \
    --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
    --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
    --mount=type=cache,target=/app/target/ \
    --mount=type=cache,target=/usr/local/cargo/git/db \
    --mount=type=cache,target=/usr/local/cargo/registry/ \
    cargo build --locked --release && \
    cp ./target/release/$APP_NAME /bin/server

################################################################################
# Create a new stage for running the application that contains the minimal
# We use dhi.io/static for the final stage because it’s a minimal Docker Hardened Image runtime (basically “just # enough OS to run the binary”), which helps keep the image small and with a lower attack surface compared to a # # full Alpine/Debian runtime.
################################################################################

FROM dhi.io/static:20250419 AS final

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

# Configure rocket to listen on all interfaces.
ENV ROCKET_ADDRESS=0.0.0.0

# Expose the port that the application listens on.
EXPOSE 8000

# What the container should run when it is started.
CMD ["/bin/server"]

```

{{< /tab >}}
{{< tab name="Using the Docker Official Images" >}}

```dockerfile {title=Dockerfile}
# Pin the Rust toolchain version used in the build stage.
ARG RUST_VERSION=1.92

# Name of the compiled binary produced by Cargo (must match Cargo.toml package name).
ARG APP_NAME=docker-rust-hello

################################################################################
# Build stage (DOI Rust image)
# This stage compiles the application.
################################################################################

FROM docker.io/library/rust:${RUST_VERSION}-alpine AS build

# Re-declare args inside the stage if you want to use them here.
ARG APP_NAME

# All build steps happen inside /app.
WORKDIR /app

# Install build dependencies needed to compile Rust crates on Alpine
RUN apk add --no-cache clang lld musl-dev git

# Build the application 
RUN --mount=type=bind,source=src,target=src \
    --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
    --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
    --mount=type=cache,target=/app/target/ \
    --mount=type=cache,target=/usr/local/cargo/git/db \
    --mount=type=cache,target=/usr/local/cargo/registry/ \
    cargo build --locked --release && \
    cp ./target/release/$APP_NAME /bin/server

################################################################################
# Runtime stage (DOI Alpine image)
# This stage runs the already-compiled binary with minimal dependencies.
################################################################################

FROM docker.io/library/alpine:3.18 AS final

# Create a non-privileged user (recommended best practice)
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# Drop privileges for runtime.
USER appuser

# Copy only the compiled binary from the build stage.
COPY --from=build /bin/server /bin/

# Rocket: listen on all interfaces inside the container.
ENV ROCKET_ADDRESS=0.0.0.0

# Document the port your app listens on.
EXPOSE 8000

# Start the application.
CMD ["/bin/server"]
```
{{< /tab >}}
{{< /tabs >}}



For building an image, only the Dockerfile is necessary. Open the Dockerfile
in your favorite IDE or text editor and see what it contains. To learn more
about Dockerfiles, see the [Dockerfile reference](/reference/dockerfile.md).

### .dockerignore file

The [`.dockerignore`](/reference/dockerfile.md#dockerignore-file) file specifies patterns and paths that you don't want copied into the image in order to keep the image as small as possible. Open up the `.dockerignore` file in your favorite IDE or text editor to review its contents.

### Build an image

Now that you’ve created the Dockerfile, you can build the image. To do this, use
the `docker build` command. The `docker build` command builds Docker images from
a Dockerfile and a context. A build's context is the set of files located in
the specified PATH or URL. The Docker build process can access any of the files
located in this context.

The build command optionally takes a `--tag` flag. The tag sets the name of the
image and an optional tag in the format `name:tag`. If you don't pass a tag,
Docker uses "latest" as its default tag.

Build the Docker image.

```console
$ docker build --tag docker-rust-image-dhi .
```

You should see output like the following.

```console
[+] Building 1.4s (13/13) FINISHED                                                                                                                                 docker:desktop-linux
 => [internal] load build definition from Dockerfile                                                                                                                               0.0s
 => => transferring dockerfile: 1.67kB                                                                                                                                             0.0s
 => [internal] load metadata for dhi.io/static:20250419                                                                                                                            1.1s
 => [internal] load metadata for dhi.io/rust:1.92-alpine3.22-dev                                                                                                                   1.2s
 => [auth] static:pull token for dhi.io                                                                                                                                            0.0s
 => [auth] rust:pull token for dhi.io                                                                                                                                              0.0s
 => [internal] load .dockerignore                                                                                                                                                  0.0s
 => => transferring context: 646B                                                                                                                                                  0.0s
 => [build 1/3] FROM dhi.io/rust:1.92-alpine3.22-dev@sha256:49eb72825a9e15fe48f2c4875a63c7e7f52a5b430bb52b8254b91d132aa5bf38                                                       0.0s
 => => resolve dhi.io/rust:1.92-alpine3.22-dev@sha256:49eb72825a9e15fe48f2c4875a63c7e7f52a5b430bb52b8254b91d132aa5bf38                                                             0.0s
 => [final 1/2] FROM dhi.io/static:20250419@sha256:74fc43fa240887b8159970e434244039aab0c6efaaa9cf044004cdc22aa2a34d                                                                0.0s
 => => resolve dhi.io/static:20250419@sha256:74fc43fa240887b8159970e434244039aab0c6efaaa9cf044004cdc22aa2a34d                                                                      0.0s
 => [internal] load build context                                                                                                                                                  0.0s
 => => transferring context: 117B                                                                                                                                                  0.0s
 => CACHED [build 2/3] WORKDIR /build                                                                                                                                              0.0s
 => CACHED [build 3/3] RUN --mount=type=bind,source=src,target=src     --mount=type=bind,source=Cargo.toml,target=Cargo.toml     --mount=type=bind,source=Cargo.lock,target=Cargo  0.0s
 => CACHED [final 2/2] COPY --from=build /build/target/release/docker-rust-hello /server                                                                                           0.0s
 => exporting to image                                                                                                                                                             0.1s
 => => exporting layers                                                                                                                                                            0.0s
 => => exporting manifest sha256:cc937bbdd712ef6e5445501f77e02ef8455ef64c567598786d46b7b21a4d4fa8                                                                                  0.0s
 => => exporting config sha256:077507b483af4b5e1a928e527e4bb3a4aaf0557e1eea81cd39465f564c187669                                                                                    0.0s
 => => exporting attestation manifest sha256:11b60e7608170493da1fdd88c120e2d2957f2a72a22edbc9cfbdd0dd37d21f89                                                                      0.0s
 => => exporting manifest list sha256:99a1b925a8d6ebf80e376b8a1e50cd806ec42d194479a3375e1cd9d2911b4db9                                                                             0.0s
 => => naming to docker.io/library/docker-rust-image-dhi:latest                                                                                                                    0.0s
 => => unpacking to docker.io/library/docker-rust-image-dhi:latest                                                                                                                 0.0s

View build details: docker-desktop://dashboard/build/desktop-linux/desktop-linux/yczk0ijw8kc5g20e8nbc8r6lj
```

### View local images

To see a list of images you have on your local machine, you have two options. One is to use the Docker CLI and the other is to use [Docker Desktop](/manuals/desktop/use-desktop/images.md). As you are working in the terminal already, take a look at listing images using the CLI.

To list images, run the `docker images` command.

```console
$ docker images
IMAGE                          ID             DISK USAGE   CONTENT SIZE   EXTRA
docker-rust-image-dhi:latest   99a1b925a8d6       11.6MB         2.45MB    U   
```

You should see at least one image listed, including the image you just built `docker-rust-image-dhi:latest`.

### Tag images

As mentioned earlier, an image name is made up of slash-separated name components. Name components may contain lowercase letters, digits, and separators. A separator can include a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.

An image is made up of a manifest and a list of layers. Don't worry too much about manifests and layers at this point other than a "tag" points to a combination of these artifacts. You can have multiple tags for an image. Create a second tag for the image you built and take a look at its layers.

To create a new tag for the image you built, run the following command.

```console
$ docker tag docker-rust-image-dhi:latest docker-rust-image-dhi:v1.0.0
```

The `docker tag` command creates a new tag for an image. It doesn't create a new image. The tag points to the same image and is just another way to reference the image.

Now, run the `docker images` command to see a list of the local images.

```console
$ docker images
IMAGE                          ID             DISK USAGE   CONTENT SIZE   EXTRA
docker-rust-image-dhi:latest   99a1b925a8d6       11.6MB         2.45MB    U   
docker-rust-image-dhi:v1.0.0   99a1b925a8d6       11.6MB         2.45MB    U  
```

You can see that two images start with `docker-rust-image-dhi`. You know they're the same image because if you take a look at the `IMAGE ID` column, you can see that the values are the same for the two images.

Remove the tag you just created. To do this, use the `rmi` command. The `rmi` command stands for remove image.

```console
$ docker rmi docker-rust-image-dhi:v1.0.0
Untagged: docker-rust-image-dhi:v1.0.0
```

Note that the response from Docker tells you that Docker didn't remove the image, but only "untagged" it. You can check this by running the `docker images` command.

```console
$ docker images
IMAGE                          ID             DISK USAGE   CONTENT SIZE   EXTRA
docker-rust-image-dhi:latest   99a1b925a8d6       11.6MB         2.45MB    U   
```

Docker removed the image tagged with `:v1.0.0`, but the `docker-rust-image-dhi:latest` tag is available on your machine.

### Summary

This section showed how to create a Dockerfile and `.dockerignore` file for a Rust application, build an image, and tag and list images.

Related information:

- [Dockerfile reference](/reference/dockerfile.md)
- [.dockerignore file](/reference/dockerfile.md#dockerignore-file)
- [docker build CLI reference](/reference/cli/docker/buildx/build/)
- [Docker Hardened Images](/dhi/)

### Next steps

In the next section learn how to run your image as a container.

## Run your Rust image as a container

### Prerequisite

You have completed [Build your Rust image](./) and you have built an image.

### Overview

A container is a normal operating system process except that Docker isolates this process so that it has its own file system, its own networking, and its own isolated process tree separate from the host.

To run an image inside of a container, you use the `docker run` command. The `docker run` command requires one parameter which is the name of the image.

### Run an image

Use `docker run` to run the image you built in [Build your Rust image](./).

```console
$ docker run docker-rust-image-dhi
```

After running this command, you’ll notice that you weren't returned to the command prompt. This is because your application is a server that runs in a loop waiting for incoming requests without returning control back to the OS until you stop the container.

Open a new terminal then make a request to the server using the `curl` command.

```console
$ curl http://localhost:8000
```

You should see output like the following.

```console
curl: (7) Failed to connect to localhost port 8000 after 2236 ms: Couldn't connect to server
```

As you can see, your `curl` command failed. This means you weren't able to connect to the localhost on port 8000. This is normal because your container is running in isolation which includes networking. Stop the container and restart with port 8000 published on your local network.

To stop the container, press ctrl-c. This will return you to the terminal prompt.

To publish a port for your container, you’ll use the `--publish` flag (`-p` for short) on the `docker run` command. The format of the `--publish` command is `[host port]:[container port]`. So, if you wanted to expose port 8000 inside the container to port 3001 outside the container, you would pass `3001:8000` to the `--publish` flag.

You didn't specify a port when running the application in the container and the default is 8000. If you want your previous request going to port 8000 to work, you can map the host's port 3001 to the container's port 8000:

```console
$ docker run --publish 3001:8000 docker-rust-image-dhi
```

Now, rerun the curl command. Remember to open a new terminal.

```console
$ curl http://localhost:3001
```

You should see output like the following.

```console
Hello, Docker!
```

Success! You were able to connect to the application running inside of your container on port 8000. Switch back to the terminal where your container is running and stop it.

Press ctrl-c to stop the container.

### Run in detached mode

This is great so far, but your sample application is a web server and you don't have to be connected to the container. Docker can run your container in detached mode or in the background. To do this, you can use the `--detach` or `-d` for short. Docker starts your container the same as before but this time will "detach" from the container and return you to the terminal prompt.

```console
$ docker run -d -p 3001:8000 docker-rust-image-dhi
3e4830e7f01304811d97dd3469d47a0c7a916a8b6c28ce0ef19c6f689a521144
```

Docker started your container in the background and printed the Container ID on the terminal.

Again, make sure that your container is running properly. Run the curl command again.

```console
$ curl http://localhost:3001
```

You should see output like the following.

```console
Hello, Docker!
```

### List containers

Since you ran your container in the background, how do you know if your container is running or what other containers are running on your machine? Well, to see a list of containers running on your machine, run `docker ps`. This is similar to how you use the ps command in Linux to see a list of processes.

You should see output like the following.

```console
CONTAINER ID   IMAGE                   COMMAND                  CREATED          STATUS          PORTS                                         NAMES
3e4830e7f013   docker-rust-image-dhi   "/server"                23 seconds ago   Up 22 seconds   0.0.0.0:3001->8000/tcp, [::]:3001->8000/tcp   youthful_lamport
```

The `docker ps` command provides a bunch of information about your running containers. You can see the container ID, the image running inside the container, the command that was used to start the container, when it was created, the status, ports that were exposed, and the name of the container.

You are probably wondering where the name of your container is coming from. Since you didn’t provide a name for the container when you started it, Docker generated a random name. You’ll fix this in a minute, but first you need to stop the container. To stop the container, run the `docker stop` command which does just that, stops the container. You need to pass the name of the container or you can use the container ID.

```console
$ docker stop youthful_lamport
youthful_lamport
```

Now, rerun the `docker ps` command to see a list of running containers.

```console
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

### Stop, start, and name containers

You can start, stop, and restart Docker containers. When you stop a container, it's not removed, but the status is changed to stopped and the process inside the container is stopped. When you ran the `docker ps` command in the previous module, the default output only shows running containers. When you pass the `--all` or `-a` for short, you see all containers on your machine, irrespective of their start or stop status.

```console
$ docker ps -a
CONTAINER ID   IMAGE                   COMMAND                  CREATED              STATUS                          PORTS                                         NAMES
3e4830e7f013   docker-rust-image-dhi   "/server"                About a minute ago   Exited (0) 28 seconds ago                                                     youthful_lamport
60009b7eaf40   docker-rust-image-dhi   "/server"                2 minutes ago        Exited (0) About a minute ago                                                 sharp_noyce
152e1d7d9eea   docker-rust-image-dhi   "/server ."              4 minutes ago        Exited (0) 2 minutes ago                                                      magical_bhabha
```

You should now see several containers listed. These are containers that you started and stopped but you haven't removed.

Restart the container that you just stopped. Locate the name of the container you just stopped and replace the name of the container in following restart command.

```console
$ docker restart youthful_lamport
```

Now list all the containers again using the `docker ps --all` command.

```console
$ docker ps --all
CONTAINER ID   IMAGE                   COMMAND                  CREATED             STATUS                         PORTS                                         NAMES
3e4830e7f013   docker-rust-image-dhi   "/server"                3 minutes ago       Up 7 seconds                   0.0.0.0:3001->8000/tcp, [::]:3001->8000/tcp   youthful_lamport
60009b7eaf40   docker-rust-image-dhi   "/server"                4 minutes ago       Exited (0) 3 minutes ago                                                     sharp_noyce
152e1d7d9eea   docker-rust-image-dhi   "/server ."              5 minutes ago       Exited (0) 4 minutes ago                                                     magical_bhabha
```

Notice that the container you just restarted has been started in detached mode. Also, observe the status of the container is "Up X seconds". When you restart a container, it starts with the same flags or commands that it was originally started with.

Now, stop and remove all of your containers and take a look at fixing the random naming issue. Stop the container you just started. Find the name of your running container and replace the name in the following command with the name of the container on your system.

```console
$ docker stop youthful_lamport
youthful_lamport
```

Now that you have stopped all of your containers, remove them. When you remove a container, it's no longer running, nor is it in the stopped status, but the process inside the container has been stopped and the metadata for the container has been removed.

To remove a container, run the `docker rm` command with the container name. You can pass multiple container names to the command using a single command. Again, replace the container names in the following command with the container names from your system.

```console
$ docker rm youthful_lamport friendly_montalcini tender_bose
youthful_lamport
sharp_noyce
magical_bhabha
```

Run the `docker ps --all` command again to see that Docker removed all containers.

Now, it's time to address the random naming issue. Standard practice is to name your containers for the simple reason that it's easier to identify what's running in the container and what application or service it's associated with.

To name a container, pass the `--name` flag to the `docker run` command.

```console
$ docker run -d -p 3001:8000 --name docker-rust-container docker-rust-image-dhi
1aa5d46418a68705c81782a58456a4ccdb56a309cb5e6bd399478d01eaa5cdda
$ docker ps
CONTAINER ID   IMAGE                   COMMAND                  CREATED         STATUS         PORTS                                         NAMES
219b2e3c7c38   docker-rust-image-dhi   "/server"                6 seconds ago   Up 5 seconds   0.0.0.0:3001->8000/tcp, [::]:3001->8000/tcp   docker-rust-container
```

Now you can identify your container based on the name.

### Summary

In this section, you took a look at running containers. You also took a look at managing containers by starting, stopping, and restarting them. And finally, you looked at naming your containers so they are more identifiable.

Related information:

- [docker run CLI reference](/reference/cli/docker/container/run/)

### Next steps

In the next section, you’ll learn how to run a database in a container and connect it to a Rust application.

## Develop your Rust application

### Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have completed the walkthroughs in the Docker Desktop [Learning Center](/manuals/desktop/use-desktop/_index.md) to learn about Docker concepts.
- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

### Overview

In this section, you’ll learn how to use volumes and networking in Docker. You’ll also use Docker to build your images and Docker Compose to make everything a whole lot easier.

First, you’ll take a look at running a database in a container and how you can use volumes and networking to persist your data and let your application talk with the database. Then you’ll pull everything together into a Compose file which lets you set up and run a local development environment with one command.

### Run a database in a container

Instead of downloading PostgreSQL, installing, configuring, and then running the PostgreSQL database as a service, you can use the Docker Official Image for PostgreSQL and run it in a container.

Before you run PostgreSQL in a container, create a volume that Docker can manage to store your persistent data and configuration. Use the named volumes feature that Docker provides instead of using bind mounts.

Run the following command to create your volume.

```console
$ docker volume create db-data
```

Now create a network that your application and database will use to talk to each other. The network is called a user-defined bridge network and gives you a nice DNS lookup service which you can use when creating your connection string.

```console
$ docker network create postgresnet
```

Now you can run PostgreSQL in a container and attach to the volume and network that you created previously. Docker pulls the image from Hub and runs it for you locally.
In the following command, option `--mount` is for starting the container with a volume. For more information, see [Docker volumes](/manuals/engine/storage/volumes.md).

```console
$ docker run --rm -d --mount \
  "type=volume,src=db-data,target=/var/lib/postgresql" \
  -p 5432:5432 \
  --network postgresnet \
  --name db \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=example \
  postgres:18
```

Now, make sure that your PostgreSQL database is running and that you can connect to it. Connect to the running PostgreSQL database inside the container.

```console
$ docker exec -it db psql -U postgres
```

You should see output like the following.

```console
psql (15.3 (Debian 15.3-1.pgdg110+1))
Type "help" for help.

postgres=#
```

In the previous command, you logged in to the PostgreSQL database by passing the `psql` command to the `db` container. Press ctrl-d to exit the PostgreSQL interactive terminal.

### Get and run the sample application

For the sample application, you'll use a variation of the backend from the react-rust-postgres application from [Awesome Compose](https://github.com/docker/awesome-compose/tree/master/react-rust-postgres).

1. Clone the sample application repository using the following command.

   ```console
   $ git clone https://github.com/docker/docker-rust-postgres
   ```

2. In the cloned repository's directory, create a `Dockerfile`. This application includes a `migrations` directory (in addition to `src`) to initialize the database, so the Dockerfile includes a bind mount for that directory in the build stage.

   ```dockerfile {hl_lines="28"}
   # syntax=docker/dockerfile:1

   # Comments are provided throughout this file to help you get started.
   # If you need more help, visit the Dockerfile reference guide at
   # https://docs.docker.com/reference/dockerfile/

   ################################################################################
   # Create a stage for building the application.

   ARG RUST_VERSION=1.70.0
   ARG APP_NAME=react-rust-postgres
   FROM rust:${RUST_VERSION}-slim-bullseye AS build
   ARG APP_NAME
   WORKDIR /app

   # Build the application.
   # Leverage a cache mount to /usr/local/cargo/registry/
   # for downloaded dependencies and a cache mount to /app/target/ for
   # compiled dependencies which will speed up subsequent builds.
   # Leverage a bind mount to the src directory to avoid having to copy the
   # source code into the container. Once built, copy the executable to an
   # output directory before the cache mounted /app/target is unmounted.
   RUN --mount=type=bind,source=src,target=src \
       --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
       --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
       --mount=type=cache,target=/app/target/ \
       --mount=type=cache,target=/usr/local/cargo/registry/ \
       --mount=type=bind,source=migrations,target=migrations \
       <<EOF
   set -e
   cargo build --locked --release
   cp ./target/release/$APP_NAME /bin/server
   EOF

   ################################################################################
   # Create a new stage for running the application that contains the minimal
   # runtime dependencies for the application. This often uses a different base
   # image from the build stage where the necessary files are copied from the build
   # stage.
   #
   # The example below uses the debian bullseye image as the foundation for    running the app.
   # By specifying the "bullseye-slim" tag, it will also use whatever happens to    be the
   # most recent version of that tag when you build your Dockerfile. If
   # reproducibility is important, consider using a digest
   # (e.g.,    debian@sha256:ac707220fbd7b67fc19b112cee8170b41a9e97f703f588b2cdbbcdcecdd8af57).
   FROM debian:bullseye-slim AS final

   # Create a non-privileged user that the app will run under.
   # See https://docs.docker.com/develop/develop-images/dockerfile_best-practices/   #user
   ARG UID=10001
   RUN adduser \
       --disabled-password \
       --gecos "" \
       --home "/nonexistent" \
       --shell "/sbin/nologin" \
       --no-create-home \
       --uid "${UID}" \
       appuser
   USER appuser

   # Copy the executable from the "build" stage.
   COPY --from=build /bin/server /bin/

   # Expose the port that the application listens on.
   EXPOSE 8000

   # What the container should run when it is started.
   CMD ["/bin/server"]
   ```

3. In the cloned repository's directory, run `docker build` to build the image.

   ```console
   $ docker build -t rust-backend-image .
   ```

4. Run `docker run` with the following options to run the image as a container on the same network as the database.

   ```console
   $ docker run \
     --rm -d \
     --network postgresnet \
     --name docker-develop-rust-container \
     -p 3001:8000 \
     -e PG_DBNAME=example \
     -e PG_HOST=db \
     -e PG_USER=postgres \
     -e PG_PASSWORD=mysecretpassword \
     -e ADDRESS=0.0.0.0:8000 \
     -e RUST_LOG=debug \
     rust-backend-image
   ```

5. Curl the application to verify that it connects to the database.

   ```console
   $ curl http://localhost:3001/users
   ```

   You should get a response like the following.

   ```json
   [{ "id": 1, "login": "root" }]
   ```

### Use Compose to develop locally

In the cloned repository's directory, create a `compose.yaml` file. Using Compose, you don't have to type all the parameters to pass to the `docker run` command — you can declare them in the file instead.

You need to update the following items in the `compose.yaml` file:

- Uncomment all of the database instructions.
- Add the environment variables under the server service.

The following is the updated `compose.yaml` file.

```yaml {hl_lines=["17-23","30-55"]}
# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Docker compose reference guide at
# https://docs.docker.com/reference/compose-file/

# Here the instructions define your application as a service called "server".
# This service is built from the Dockerfile in the current directory.
# You can add other services your application may depend on here, such as a
# database or a cache. For examples, see the Awesome Compose repository:
# https://github.com/docker/awesome-compose
services:
  server:
    build:
      context: .
      target: final
    ports:
      - 8000:8000
    environment:
      - PG_DBNAME=example
      - PG_HOST=db
      - PG_USER=postgres
      - PG_PASSWORD=mysecretpassword
      - ADDRESS=0.0.0.0:8000
      - RUST_LOG=debug
    # The commented out section below is an example of how to define a PostgreSQL
    # database that your application can use. `depends_on` tells Docker Compose to
    # start the database before your application. The `db-data` volume persists the
    # database data between container restarts. The `db-password` secret is used
    # to set the database password. You must create `db/password.txt` and add
    # a password of your choosing to it before running `docker compose up`.
    depends_on:
      db:
        condition: service_healthy
  db:
    image: postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

Note that the file doesn't specify a network for those 2 services. When you use Compose, it automatically creates a network and connects the services to it. For more information see [Networking in Compose](/manuals/compose/how-tos/networking.md).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

If you have any other containers running from the previous sections, [stop](#stop-start-and-name-containers) them now.

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

The command passes the `--build` flag so Docker will compile your image and then start the containers.

Now test your API endpoint. Open a new terminal then make a request to the server using the curl commands:

```console
$ curl http://localhost:8000/users
```

You should receive the following response:

```json
[{ "id": 1, "login": "root" }]
```

### Summary

In this section, you took a look at setting up your Compose file to run your Rust application and database with a single command.

Related information:

- [Docker volumes](/manuals/engine/storage/volumes.md)
- [Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

## Configure CI/CD for your Rust application

### Prerequisites

Complete the previous sections of this guide, starting with [Develop your Rust application](develop.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

### Overview

In this section, you'll learn how to set up and use GitHub Actions to build and push your Docker image to Docker Hub. You will complete the following steps:

1. Create a new repository on GitHub.
2. Define the GitHub Actions workflow.
3. Run the workflow.

### Step one: Create the repository

Create a GitHub repository, configure the Docker Hub credentials, and push your source code.

1. [Create a new repository](https://github.com/new) on GitHub.

2. Open the repository **Settings**, and go to **Secrets and variables** >
   **Actions**.

3. Create a new **Repository variable** named `DOCKER_USERNAME` and your Docker ID as a value.

4. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token) for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

5. Add the PAT as a **Repository secret** in your GitHub repository, with the name
   `DOCKERHUB_TOKEN`.

6. In your local repository on your machine, run the following command to change
   the origin to the repository you just created. Make sure you change
   `your-username` to your GitHub username and `your-repository` to the name of
   the repository you created.

   ```console
   $ git remote set-url origin https://github.com/your-username/your-repository.git
   ```

7. Run the following commands to stage, commit, and push your local repository to GitHub.

   ```console
   $ git add -A
   $ git commit -m "my commit"
   $ git push -u origin main
   ```

### Step two: Set up the workflow

Set up your GitHub Actions workflow for building, testing, and pushing the image
to Docker Hub.

1. Go to your repository on GitHub and then select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub actions workflow file in
   your repository, under `.github/workflows/main.yml` by default.

3. In the editor window, copy and paste the following YAML configuration.

   ```yaml
   name: ci

   on:
     push:
       branches:
         - main

   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - name: Login to Docker Hub
           uses: docker/login-action@{{% param "login_action_version" %}}
           with:
             username: ${{ vars.DOCKER_USERNAME }}
             password: ${{ secrets.DOCKERHUB_TOKEN }}

         - name: Set up Docker Buildx
           uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

         - name: Build and push
           uses: docker/build-push-action@{{% param "build_push_action_version" %}}
           with:
             push: true
             tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
   ```

   For more information about the YAML syntax for `docker/build-push-action`,
   refer to the [GitHub Action README](https://github.com/docker/build-push-action/blob/master/README.md).

### Step three: Run the workflow

Save the workflow file and run the job.

1. Select **Commit changes...** and push the changes to the `main` branch.

   After pushing the commit, the workflow starts automatically.

2. Go to the **Actions** tab. It displays the workflow.

   Selecting the workflow shows you the breakdown of all the steps.

3. When the workflow is complete, go to your
   [repositories on Docker Hub](https://hub.docker.com/repositories).

   If you see the new repository in that list, it means the GitHub Actions
   successfully pushed the image to Docker Hub.

### Summary

In this section, you learned how to set up a GitHub Actions workflow for your Rust application.

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### Next steps

Next, learn how you can locally test and debug your workloads on Kubernetes before deploying.

## Test your Rust deployment

### Prerequisites

- Complete the previous sections of this guide, starting with [Develop your Rust application](develop.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This lets you to test and debug your workloads on Kubernetes locally before deploying.

### Create a Kubernetes YAML file

In your `docker-rust-postgres` directory, create a file named
`docker-rust-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your Rust application](./).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    service: server
  name: server
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: server
  strategy: {}
  template:
    metadata:
      labels:
        service: server
    spec:
      initContainers:
        - name: wait-for-db
          image: busybox:1.28
          command:
            [
              "sh",
              "-c",
              'until nc -zv db 5432; do echo "waiting for db"; sleep 2; done;',
            ]
      containers:
        - image: DOCKER_USERNAME/REPO_NAME
          name: server
          imagePullPolicy: Always
          ports:
            - containerPort: 8000
              hostPort: 5000
              protocol: TCP
          env:
            - name: ADDRESS
              value: 0.0.0.0:8000
            - name: PG_DBNAME
              value: example
            - name: PG_HOST
              value: db
            - name: PG_PASSWORD
              value: mysecretpassword
            - name: PG_USER
              value: postgres
            - name: RUST_LOG
              value: debug
          resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    service: db
  name: db
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: db
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        service: db
    spec:
      containers:
        - env:
            - name: POSTGRES_DB
              value: example
            - name: POSTGRES_PASSWORD
              value: mysecretpassword
            - name: POSTGRES_USER
              value: postgres
          image: postgres:18
          name: db
          ports:
            - containerPort: 5432
              protocol: TCP
          resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    service: server
  name: server
  namespace: default
spec:
  type: NodePort
  ports:
    - name: "5000"
      port: 5000
      targetPort: 8000
      nodePort: 30001
  selector:
    service: server
status:
  loadBalancer: {}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    service: db
  name: db
  namespace: default
spec:
  ports:
    - name: "5432"
      port: 5432
      targetPort: 5432
  selector:
    service: db
status:
  loadBalancer: {}
```

In this Kubernetes YAML file, there are four objects, separated by the `---`. In addition to a Service and Deployment for the database, the other two objects are:

- A Deployment, describing a scalable group of identical pods. In this case,
  you'll get just one replica, or copy of your pod. That pod, which is
  described under `template`, has just one container in it. The container is
  created from the image built by GitHub Actions in [Configure CI/CD for your
  Rust application](./).
- A NodePort service, which will route traffic from port 30001 on your host to
  port 5000 inside the pods it routes to, allowing you to reach your app
  from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

### Deploy and check your application

1. In a terminal, navigate to `docker-rust-postgres` and deploy your application
   to Kubernetes.

   ```console
   $ kubectl apply -f docker-rust-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```shell
   deployment.apps/server created
   deployment.apps/db created
   service/server created
   service/db created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
   db       1/1     1            1           2m21s
   server   1/1     1            1           2m21s
   ```

   This indicates all of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   db           ClusterIP   10.105.167.81    <none>        5432/TCP         109s
   kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP          9d
   server       NodePort    10.101.235.213   <none>        5000:30001/TCP   109s
   ```

   In addition to the default `kubernetes` service, you can see your `service-entrypoint` service, accepting traffic on port 30001/TCP.

3. In a terminal, curl the service.

   ```console
   $ curl http://localhost:30001/users
   [{"id":1,"login":"root"}]
   ```

4. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-rust-kubernetes.yaml
   ```

### Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
- [Swarm mode overview](/manuals/engine/swarm/_index.md)
