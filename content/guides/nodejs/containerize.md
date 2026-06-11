---
title: Containerize a Node.js application
linkTitle: Containerize your app
weight: 10
keywords: node.js, node, containerize, initialize
description: Learn how to containerize a Node.js application with Docker.
aliases:
  - /get-started/nodejs/build-images/
  - /language/nodejs/build-images/
  - /language/nodejs/run-containers/
  - /language/nodejs/containerize/
  - /guides/language/nodejs/containerize/
---

## Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You're familiar with basic Docker concepts. If you're new to Docker, start with [Get started](/get-started/introduction/).

## Overview

Containerizing your application means packaging it together with its
dependencies, configuration, and runtime into a single portable unit called a
container image. Running that image creates a container, an isolated process
that behaves the same on any machine, whether it's your laptop, a CI runner, or
a production server.

In this section, you'll containerize a simple [Express.js](https://expressjs.com/) API written in TypeScript. You'll write a `Dockerfile` that describes how to build the image, add a `compose.yaml` file that defines how Docker runs your container, and then build and start the application with one command.

You'll use [Docker Hardened Images](/dhi/) as the base. These are minimal, secure Node.js images maintained by Docker.

This guide focuses on a backend Node.js API. If you're building a standalone frontend application, Docker has dedicated guides for [React.js](/guides/reactjs/), [Vue.js](/guides/vuejs/), [Angular](/guides/angular/), and [Next.js](/guides/nextjs/).

## Create the application

The sample application is a minimal Express API with a single endpoint that returns a JSON greeting. Create the following files in a new `nodejs-docker-example` directory. To create all the files at once, switch to the **Scaffold script** tab in the file browser and copy the shell command.

{{< files name="nodejs-docker-example" >}}

{{< file path="src/index.ts" status="new" >}}
```typescript
// A minimal Express application.
// The root endpoint (GET /) returns a JSON greeting.
// See https://expressjs.com/ for the framework reference.

import express, { type Request, type Response } from 'express';

const app = express();
const port = parseInt(process.env.PORT ?? '3000', 10);

app.get('/', (_req: Request, res: Response) => {
  res.json({ message: 'Hello World' });
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
```
{{< /file >}}

{{< file path="package.json" status="new" >}}
```json
{
  "name": "nodejs-docker-example",
  "version": "1.0.0",
  "description": "A minimal Node.js TypeScript application.",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch src/index.ts"
  },
  "dependencies": {
    "express": "^4.21.2"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "tsx": "^4.19.3",
    "typescript": "^5.8.3"
  }
}
```
{{< /file >}}

{{< file path="tsconfig.json" status="new" >}}
```json
{
  // TypeScript compiler configuration for the Node.js application.
  // Compiles src/ to dist/ as CommonJS modules targeting ES2022.
  // See https://www.typescriptlang.org/tsconfig/ for all options.
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```
{{< /file >}}

{{< file path=".gitignore" status="new" >}}
```text
# Files and directories that Git should ignore. Covers Node.js dependencies,
# TypeScript build output, environment files, and common editor artifacts.
# See https://git-scm.com/docs/gitignore for syntax reference.

node_modules/
dist/
.env
*.log
.DS_Store
coverage/
db/password.txt
```
{{< /file >}}

{{< /files >}}

If you have Node.js installed and want to verify the app works before containerizing it, you can run it locally.

To run in development mode with hot-reload:

```console
$ npm install
$ npm run dev
```

To run the compiled production build (matching what the Dockerfile does):

```console
$ npm install
$ npm run build
$ npm start
```

Then open [http://localhost:3000](http://localhost:3000) in your browser. You should see `{"message":"Hello World"}`.

If you don't have Node.js installed, skip ahead. The remaining steps run the application in a container, with no local Node.js required.

## Create the Docker assets

Sign in to the DHI registry so Docker can pull the Node.js base images during the build. The available Node.js images are listed in the [catalog](https://hub.docker.com/hardened-images/catalog/dhi/node).

```console
$ docker login dhi.io
```

Add the following three files to your `nodejs-docker-example` directory. The `Dockerfile` describes how to build the image, `compose.yaml` defines how Docker runs the container, and `.dockerignore` keeps unwanted files out of the build context.

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for
> your project. Ask Gordon to create a Dockerfile, Compose file, and
> `.dockerignore` tailored to your application.

{{< files name="nodejs-docker-example" >}}

{{< file path="src/index.ts" >}}
```typescript
// A minimal Express application.
// The root endpoint (GET /) returns a JSON greeting.
// See https://expressjs.com/ for the framework reference.

import express, { type Request, type Response } from 'express';

const app = express();
const port = parseInt(process.env.PORT ?? '3000', 10);

app.get('/', (_req: Request, res: Response) => {
  res.json({ message: 'Hello World' });
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
```
{{< /file >}}

{{< file path="package.json" >}}
```json
{
  "name": "nodejs-docker-example",
  "version": "1.0.0",
  "description": "A minimal Node.js TypeScript application.",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch src/index.ts"
  },
  "dependencies": {
    "express": "^4.21.2"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "tsx": "^4.19.3",
    "typescript": "^5.8.3"
  }
}
```
{{< /file >}}

{{< file path="tsconfig.json" >}}
```json
{
  // TypeScript compiler configuration for the Node.js application.
  // Compiles src/ to dist/ as CommonJS modules targeting ES2022.
  // See https://www.typescriptlang.org/tsconfig/ for all options.
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```
{{< /file >}}

{{< file path="Dockerfile" status="new" >}}
```dockerfile
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# This Dockerfile uses Docker Hardened Images (DHI) for enhanced security.
# For more information, see https://docs.docker.com/dhi/

# Builder stage: install all dependencies and compile TypeScript.
FROM dhi.io/node:24-alpine3.23-dev AS builder

WORKDIR /app

# Install dependencies as a separate step to take advantage of Docker's
# caching. Leverage a cache mount to /root/.npm to speed up subsequent
# builds. Leverage a bind mount to package.json to avoid having to copy
# it into this layer.
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm install
# Once you create a package-lock.json by running npm install locally, switch to npm ci and bind both files:
# RUN --mount=type=cache,target=/root/.npm \
#     --mount=type=bind,source=package.json,target=package.json \
#     --mount=type=bind,source=package-lock.json,target=package-lock.json \
#     npm ci

# Copy the source code into the container and compile TypeScript.
COPY . .
RUN npm run build


# Deps stage: install production dependencies only.
FROM dhi.io/node:24-alpine3.23-dev AS deps

WORKDIR /app

RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm install --omit=dev
# Once you create a package-lock.json by running npm install locally, switch to npm ci and bind both files:
# RUN --mount=type=cache,target=/root/.npm \
#     --mount=type=bind,source=package.json,target=package.json \
#     --mount=type=bind,source=package-lock.json,target=package-lock.json \
#     npm ci --omit=dev


# Runner stage: minimal runtime image with compiled app and production deps.
FROM dhi.io/node:24-alpine3.23 AS runner

ENV PATH=/app/node_modules/.bin:$PATH

WORKDIR /app

COPY --from=deps --chown=node:node /app/node_modules ./node_modules
COPY --from=builder --chown=node:node /app/dist ./dist

# Expose the port that the application listens on.
EXPOSE 3000

# Run the application.
CMD ["node", "dist/index.js"]
```
{{< /file >}}

{{< file path="compose.yaml" status="new" >}}
```yaml
# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Docker Compose reference guide at
# https://docs.docker.com/go/compose-spec-reference/

# Here the instructions define your application as a service called "server".
# This service is built from the Dockerfile in the current directory.
# You can add other services your application may depend on here, such as a
# database or a cache. For examples, see the Awesome Compose repository:
# https://github.com/docker/awesome-compose
services:
  server:
    build:
      context: .
    ports:
      - 3000:3000
```
{{< /file >}}

{{< file path=".dockerignore" status="new" >}}
```text
# Include any files or directories that you don't want to be copied to your
# container here (e.g., local build artifacts, temporary files, etc.).
#
# For more help, visit the .dockerignore file reference guide at
# https://docs.docker.com/go/build-context-dockerignore/

node_modules/
dist/
.env
.git
.gitignore
.DS_Store
npm-debug.log*
coverage/
db/
```
{{< /file >}}

{{< file path=".gitignore" >}}
```text
# Files and directories that Git should ignore. Covers Node.js dependencies,
# TypeScript build output, environment files, and common editor artifacts.
# See https://git-scm.com/docs/gitignore for syntax reference.

node_modules/
dist/
.env
*.log
.DS_Store
coverage/
db/password.txt
```
{{< /file >}}

{{< /files >}}

The `Dockerfile` uses three stages. The `builder` stage installs all dependencies and compiles TypeScript. The `deps` stage does a fresh install of production-only dependencies. The `runner` stage copies the compiled output and production node_modules into a minimal runtime image that contains only Node.js.

To learn more about each file, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yaml](/reference/compose-file/_index.md)

## Run the application

Inside the `nodejs-docker-example` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You should see `{"message":"Hello World"}`.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `nodejs-docker-example` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000).

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

## Summary

In this section, you learned how to containerize and run a Node.js application using Docker.

Related information:

- [Docker Hardened Images](/dhi/)
- [Dockerfile reference](/reference/dockerfile.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
- [Docker Compose overview](/manuals/compose/_index.md)

## Next steps

In the next section, you'll take a look at how to set up a local development environment using Docker containers.
