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

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).

## Overview

For a long time, Node.js has been the go-to runtime for server-side JavaScript applications. However, recent years have introduced new alternative runtimes, including [Deno](https://deno.land/). Like Node.js, Deno is a JavaScript and TypeScript runtime, but it takes a fresh approach with modern security features, a built-in standard library, and native support for TypeScript.

Why develop Deno applications with Docker? Having a choice of runtimes is exciting, but managing multiple runtimes and their dependencies consistently across environments can be tricky. This is where Docker proves invaluable. Using containers to create and destroy environments on demand simplifies runtime management and ensures consistency. Additionally, as Deno continues to grow and evolve, Docker helps establish a reliable and reproducible development environment, minimizing setup challenges and streamlining the workflow.

## Create the application

Create a new directory for your project and navigate into it:

```console
$ mkdir deno-docker && cd deno-docker
```

Create a file named `server.ts` with the following contents:

```typescript {title="server.ts"}
import { Application, Router } from "https://deno.land/x/oak@v12.0.0/mod.ts";

const app = new Application();
const router = new Router();

router.get("/", (context) => {
  context.response.body = { Status: "OK" };
  context.response.type = "application/json";
});

app.use(router.routes());
app.use(router.allowedMethods());

console.log("Server running on http://localhost:8000");
await app.listen({ port: 8000 });
```

## Create Docker assets

This guide uses [Docker Hardened Images (DHI)](/dhi/) as the base image. Docker Hardened Images are minimal, secure, and production-ready base images maintained by Docker.

Docker Hardened Images (DHIs) are available for Deno in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/deno). Sign in to the DHI registry before pulling:

```console
$ docker login dhi.io
```

Create a file named `Dockerfile` with the following contents:

```dockerfile {title="Dockerfile"}
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

Create a file named `compose.yml` with the following contents:

```yaml {title="compose.yml"}
services:
  server:
    image: deno-server
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8000:8000"
```

You should now have the following files in your `deno-docker` directory:

```text
├── deno-docker/
│   ├── compose.yml
│   ├── Dockerfile
│   └── server.ts
```

## Run the application

From the `deno-docker` directory, run the following command in a terminal.

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
