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

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).

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

## Create the application

Create a new directory for your project and navigate into it:

```console
$ mkdir bun-docker && cd bun-docker
```

Create a file named `server.js` with the following contents:

```js {title="server.js"}
import { serve } from "bun";

serve({
  fetch(request) {
    const url = new URL(request.url);
    if (url.pathname === "/") {
      return new Response(JSON.stringify({ Status: "OK" }), {
        headers: { "content-type": "application/json" },
      });
    }
    return new Response("Not found", { status: 404 });
  },
  port: 3000,
});

console.log("Server is running.");
```

## Create Docker assets

This guide uses [Docker Hardened Images (DHI)](/dhi/) as the base image. Docker Hardened Images are minimal, secure, and production-ready base images maintained by Docker.

Docker Hardened Images (DHIs) are available for Bun in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/bun). Sign in to the DHI registry before pulling:

```console
$ docker login dhi.io
```

Create a file named `Dockerfile` with the following contents:

```dockerfile {title="Dockerfile"}
# Use the DHI Bun image as the base image
FROM dhi.io/bun:1

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Expose the port on which the API will listen
EXPOSE 3000

# Run the server when the container launches
CMD ["bun", "server.js"]
```

Create a file named `compose.yml` with the following contents:

```yaml {title="compose.yml"}
services:
  server:
    image: bun-server
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
```

You should now have the following files in your `bun-docker` directory:

```text
├── bun-docker/
│   ├── compose.yml
│   ├── Dockerfile
│   └── server.js
```

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
