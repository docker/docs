---
title: Containerize a Deno application
linkTitle: Containerize your app
weight: 10
keywords: deno, containerize, initialize
description: Learn how to containerize a Deno application.
aliases:
  - /language/deno/containerize/
---

## Prerequisites

* You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

## Overview

For a long time, Node.js has been the go-to runtime for server-side JavaScript applications. However, recent years have introduced new alternative runtimes, including [Deno](https://deno.land/). Like Node.js, Deno is a JavaScript and TypeScript runtime, but it takes a fresh approach with modern security features, a built-in standard library, and native support for TypeScript.

Why develop Deno applications with Docker? Having a choice of runtimes is exciting, but managing multiple runtimes and their dependencies consistently across environments can be tricky. This is where Docker proves invaluable. Using containers to create and destroy environments on demand simplifies runtime management and ensures consistency. Additionally, as Deno continues to grow and evolve, Docker helps establish a reliable and reproducible development environment, minimizing setup challenges and streamlining the workflow.

## Get the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following
command to clone the repository:

```console
$ git clone https://github.com/dockersamples/docker-deno.git && cd docker-deno
```

You should now have the following contents in your `deno-docker` directory.

```text
├── deno-docker/
│ ├── compose.yml
│ ├── Dockerfile
│ ├── LICENSE
│ ├── server.ts
│ └── README.md
```

## Understand the sample application

The sample application is a simple Deno application that uses the Oak framework to create a simple API that returns a JSON response. The application listens on port 8000 and returns a message `{"Status" : "OK"}` when you access the application in a browser.

```typescript
// server.ts
import { Application, Router } from "https://deno.land/x/oak@v12.0.0/mod.ts";

const app = new Application();
const router = new Router();

// Define a route that returns JSON
router.get("/", (context) => {
  context.response.body = { Status: "OK" };
  context.response.type = "application/json";
});

app.use(router.routes());
app.use(router.allowedMethods());

console.log("Server running on http://localhost:8000");
await app.listen({ port: 8000 });
```

## Create a Dockerfile

Before creating a Dockerfile, you need to choose a base image. You can either use the [Deno Docker Official Image](https://hub.docker.com/r/denoland/deno) or a Docker Hardened Image (DHI) from the [Hardened Image catalog](https://hub.docker.com/hardened-images/catalog).

Choosing DHI offers the advantage of a production-ready image that is lightweight and secure. For more information, see [Docker Hardened Images](https://docs.docker.com/dhi/).

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}

Docker Hardened Images (DHIs) are available for Deno in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/deno). You can pull DHIs directly from the `dhi.io` registry.

1. Sign in to the DHI registry:

   ```console
   $ docker login dhi.io
   ```

2. Pull the Deno DHI as `dhi.io/deno:2`. The tag (`2`) in this example refers to the version to the latest 2.x version of Deno.

   ```console
   $ docker pull dhi.io/deno:2
   ```

For other available versions, refer to the [catalog](https://hub.docker.com/hardened-images/catalog/dhi/deno).

```dockerfile
# Use the DHI Deno image as the base image
FROM dhi.io/deno:2

# Set the working directory
WORKDIR /app

# Copy server code into the container
COPY server.ts .

# Set permissions (optional but recommended for security)
USER deno

# Expose port 8000
EXPOSE 8000

# Run the Deno server
CMD ["run", "--allow-net", "server.ts"]
```

{{< /tab >}}
{{< tab name="Using the official image" >}}

Using the Docker Official Image is straightforward. In the following Dockerfile, you'll notice that the `FROM` instruction uses `denoland/deno:latest` as the base image.

This is the official image for Deno. This image is [available on the Docker Hub](https://hub.docker.com/r/denoland/deno).

```dockerfile
# Use the official Deno image
FROM denoland/deno:latest

# Set the working directory
WORKDIR /app

# Copy server code into the container
COPY server.ts .

# Set permissions (optional but recommended for security)
USER deno

# Expose port 8000
EXPOSE 8000

# Run the Deno server
CMD ["run", "--allow-net", "server.ts"]
```

{{< /tab >}}
{{< /tabs >}}

In addition to specifying the base image, the Dockerfile also:

- Sets the working directory in the container to `/app`.
- Copies `server.ts` into the container.
- Sets the user to `deno` to run the application as a non-root user.
- Exposes port 8000 to allow traffic to the application.
- Runs the Deno server using the `CMD` instruction.
- Uses the `--allow-net` flag to allow network access to the application. The `server.ts` file uses the Oak framework to create a simple API that listens on port 8000.

## Run the application

Make sure you are in the `deno-docker` directory. Run the following command in a terminal to build and run the application.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000). You will see a message `{"Status" : "OK"}` in the browser.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `deno-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000).


In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

## Summary

In this section, you learned how you can containerize and run your Deno
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
