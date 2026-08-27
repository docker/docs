---
title: Containerize a TanStack Start Application
linkTitle: Containerize
weight: 10
keywords: tanstack start, node, image, initialize, build
description: Learn how to containerize a TanStack Start application with Docker by creating an optimized, production-ready image using best practices for performance, security, and scalability.
---

## Prerequisites

Before you begin, make sure the following tools are installed and available on
your system:

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have a [git client](https://git-scm.com/downloads). The examples in this
  section use a command-line based git client, but you can use any client.

> [!NOTE]
> New to Docker? Start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide to get familiar with key concepts like images, containers, and Dockerfiles.

---

## Overview

This guide walks you through containerizing a TanStack Start application with
Docker. You'll learn how to create a production-ready Docker image using best
practices that improve performance, security, scalability, and deployment
efficiency.

By the end of this guide, you will:

- Containerize a TanStack Start application using Docker.
- Create and optimize a Dockerfile for production builds.
- Use multi-stage builds to minimize image size.
- Run the TanStack Start server from the `.output` build directory.
- Follow best practices for building secure and maintainable Docker images.

---

## Get the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following
commands:

```console
$ git clone https://github.com/kristiyan-velkov/docker-tanstack-start-sample
$ cd docker-tanstack-start-sample
```

The sample is a TanStack Start app that uses Vite and Nitro. The production
build writes a Node.js server bundle to `.output`.

---

## Build the Docker image

TanStack Start produces a server-side bundle at build time. The production
container runs that bundle with Node.js on port 3000.

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for
> your project. Ask Gordon to create a Dockerfile, Compose file, and
> `.dockerignore` tailored to your application.

### Step 1: Create the Dockerfile

Before creating a Dockerfile, choose a base image: the [Node.js Official Image](https://hub.docker.com/_/node) or a [Docker Hardened Image (DHI)](https://hub.docker.com/hardened-images/catalog) from the Hardened Image catalog.

> [!IMPORTANT]
> This guide uses a stable Node.js image tag that is considered secure when the
> guide is written. Because new releases and security patches are published
> regularly, always review the [official Node.js Docker images](https://hub.docker.com/_/node) and select a secure, up-to-date version before building or deploying.

Create a file named `Dockerfile` with the following contents (matching the
[sample project](https://github.com/kristiyan-velkov/docker-tanstack-start-sample)):

```dockerfile
# =========================================
# Stage 1: Build the TanStack Start Application
# =========================================
ARG NODE_VERSION=24.14.0-alpine

FROM node:${NODE_VERSION} AS builder

WORKDIR /app

COPY package.json package-lock.json* ./

RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .

RUN npm run build

# =========================================
# Stage 2: Run the TanStack Start Server
# =========================================
FROM node:${NODE_VERSION} AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

COPY --from=builder /app/.output ./

USER node

EXPOSE 3000

ENTRYPOINT ["node", "server/index.mjs"]
```

This Dockerfile:

- Uses a multi-stage build to keep the final image small.
- Installs dependencies with `npm ci` for reproducible builds.
- Runs `npm run build`, which outputs the server bundle to `.output`.
- Runs the app as the non-root `node` user.
- Starts the Nitro server with `node server/index.mjs`.

### Step 2: Create the compose.yml file

Create a file named `compose.yml` with the following contents:

```yaml
services:
  tanstack-start-app:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        NODE_VERSION: 24.14.0-alpine
    image: tanstack-start-image
    container_name: tanstack-start-container
    environment:
      NODE_ENV: production
      PORT: 3000
      HOST: 0.0.0.0
    ports:
      - "3000:3000"
    restart: unless-stopped
```

### Step 3: Create the .dockerignore file

The `.dockerignore` file tells Docker which files and folders to exclude when
building the image.

> [!NOTE]
> This helps reduce image size, speed up builds, and prevent sensitive or
> unnecessary files (like `.env`, `.git`, or `node_modules`) from being added
> to the build context. To learn more, see the
> [.dockerignore reference](/reference/dockerfile/#dockerignore-file).

Create a file named `.dockerignore` with the following contents:

```dockerignore
node_modules/
.output/
.nitro/
.vinxi/
dist/
build/
out/
.vite/
coverage/
*.test.ts
*.test.tsx
*.spec.ts
*.spec.tsx
.env
.env.*
!.env.example
.git/
.vscode/
.idea/
*.log
Dockerfile*
compose*.yml
.dockerignore
README.md
```

The [sample project](https://github.com/kristiyan-velkov/docker-tanstack-start-sample/blob/main/.dockerignore) includes a more exhaustive `.dockerignore` you can copy for production use.

### Step 4: Build the TanStack Start application image

Run the following command from the root of your project:

```console
$ docker build --tag tanstack-start .
```

What this command does:

- Uses the Dockerfile in the current directory (`.`)
- Packages the application and its dependencies into a Docker image
- Tags the image as `tanstack-start` so you can reference it later

### Step 5: View local images

After building your Docker image, list locally available images:

```console
$ docker images
```

Example output:

```shell
REPOSITORY        TAG       IMAGE ID       CREATED          SIZE
tanstack-start    latest    8c5fc80f098e   14 seconds ago   130MB
```

If the build was successful, you should see the `tanstack-start` image listed.

---

## Run the containerized application

Run the image in a container and verify that your application works:

```console
$ docker run -p 3000:3000 tanstack-start
```

Open a browser and view the application at
[http://localhost:3000](http://localhost:3000). You should see your TanStack
Start web application.

Press `ctrl+c` in the terminal to stop your application.

### Run the application in the background

Run the application detached from the terminal:

```console
$ docker run -d -p 3000:3000 --name tanstack-start-app tanstack-start
```

To confirm that the container is running:

```console
$ docker ps
```

Stop the container when you're done:

```console
$ docker stop tanstack-start-app
```

### Run with Docker Compose

From the project root, build and start the service defined in `compose.yml`:

```console
$ docker compose up --build
```

Open [http://localhost:3000](http://localhost:3000) in your browser. Press
`ctrl+c` to stop the services, or run `docker compose down` in another
terminal.

---

## Summary

In this section, you containerized a TanStack Start application using Docker.

What you accomplished:

- Created a multi-stage Dockerfile that builds and runs the TanStack Start
  server from `.output`
- Added `compose.yml` and `.dockerignore` for local orchestration and lean
  build contexts
- Built and ran the containerized application on port 3000

---

## Related resources

- [Multi-stage builds](/manuals/build/building/multi-stage/) – Create
  production-ready Docker images
- [Dockerfile best practices](/build/building/best-practices/) – Write clean,
  secure, and optimized Dockerfiles
- [TanStack Start documentation](https://tanstack.com/start/latest/docs/framework/react/overview) –
  Learn about routing, SSR, and deployment options

## Next steps

In the next section, you'll set up a development workflow with Docker Compose
and Compose Watch so you can iterate on your TanStack Start app inside
containers.
