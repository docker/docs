---
title: Containerize a Bun application
linkTitle: Containerize your app
weight: 10
keywords: bun, containerize, initialize
description: Learn how to containerize a Bun application.
aliases:
  - /language/bun/containerize/
---

## Prerequisites

* You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

## Overview

For a long time, Node.js has been the de-facto runtime for server-side
JavaScript applications. Recent years have seen a rise in new alternative
runtimes in the ecosystem, including [Bun website](https://bun.sh/). Like
Node.js, Bun is a JavaScript runtime. Bun is a comparatively lightweight
runtime that is designed to be fast and efficient.

Why develop Bun applications with Docker? Having multiple runtimes to choose
from is great. But as the number of runtimes increases, it becomes challenging
to manage the different runtimes and their dependencies consistently across
environments. This is where Docker comes in. Creating and destroying containers
on demand is a great way to manage the different runtimes and their
dependencies. Also, as it's fairly a new runtime, getting a consistent
development environment for Bun can be challenging. Docker can help you set up
a consistent development environment for Bun.

## Get the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following
command to clone the repository:

```console
$ git clone https://github.com/dockersamples/bun-docker.git && cd bun-docker
```

You should now have the following contents in your `bun-docker` directory.

```text
├── bun-docker/
│ ├── compose.yml
│ ├── Dockerfile
│ ├── LICENSE
│ ├── server.js
│ └── README.md
```

## Create a Dockerfile

Before creating a Dockerfile, you need to choose a base image. You can either use the [Bun Docker Official Image](https://hub.docker.com/r/oven/bun) or a Docker Hardened Image (DHI) from the [Hardened Image catalog](https://hub.docker.com/hardened-images/catalog).

Choosing DHI offers the advantage of a production-ready image that is lightweight and secure. For more information, see [Docker Hardened Images](https://docs.docker.com/dhi/).

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}
Docker Hardened Images (DHIs) are available for Bun on [Docker Hub](https://hub.docker.com/hardened-images/catalog/dhi/bun). Unlike using the Docker Official Image, you must first mirror the Bun image into your organization and then use it as your base image. Follow the instructions in the [DHI quickstart](/dhi/get-started/) to create a mirrored repository for Bun.

Mirrored repositories must start with `dhi-`, for example: `FROM <your-namespace>/dhi-bun:<tag>`. In the following Dockerfile, the `FROM` instruction uses `<your-namespace>/dhi-bun:1` as the base image.

```dockerfile
# Use the DHI Bun image as the base image
FROM <your-namespace>/dhi-bun:1

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Expose the port on which the API will listen
EXPOSE 3000

# Run the server when the container launches
CMD ["bun", "server.js"]
```

{{< /tab >}}
{{< tab name="Using the official image" >}}

Using the official image is straightforward. In the following Dockerfile, you'll notice that the `FROM` instruction uses `oven/bun` as the base image.

You can find the image on [Docker Hub](https://hub.docker.com/r/oven/bun). This is the Docker Official Image for Bun created by Oven, the company behind Bun, and it's available on Docker Hub.

```dockerfile
# Use the official Bun image
FROM oven/bun:latest

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Expose the port on which the API will listen
EXPOSE 3000

# Run the server when the container launches
CMD ["bun", "server.js"]
```

{{< /tab >}}
{{< /tabs >}}

In addition to specifying the base image, the Dockerfile also:

- Sets the working directory in the container to `/app`.
- Copies the content of the current directory to the `/app` directory in the container.
- Exposes port 3000, where the API is listening for requests.
- And finally, starts the server when the container launches with the command `bun server.js`.

## Run the application

Inside the `bun-docker` directory, run the following command in a terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You will see a message `{"Status" : "OK"}` in the browser.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `bun-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000).


In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

## Summary

In this section, you learned how you can containerize and run your Bun
application using Docker.

Related information:

 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file](/reference/dockerfile.md#dockerignore-file)
 - [Docker Compose overview](/manuals/compose/_index.md)
 - [Compose file reference](/reference/compose-file/_index.md)
 - [Docker Hardened Images](/dhi/)

## Next steps

In the next section, you'll learn how you can develop your application using
containers.
