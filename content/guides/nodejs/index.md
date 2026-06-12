---
title: Node.js language-specific guide
linkTitle: Node.js
description: Containerize and develop Node.js applications using Docker
keywords: getting started, node, node.js
summary: |
  This guide explains how to containerize Node.js applications using Docker.
aliases:
  - /guides/nodejs/containerize/
  - /guides/nodejs/develop/
  - /guides/nodejs/run-tests/
  - /guides/nodejs/configure-github-actions/
  - /guides/nodejs/deploy/
  - /guides/nodejs/secure-supply-chain/
  - /language/nodejs/
  - /guides/language/nodejs/
  - /get-started/nodejs/build-images/
  - /language/nodejs/build-images/
  - /language/nodejs/run-containers/
  - /language/nodejs/containerize/
  - /get-started/nodejs/develop/
  - /language/nodejs/develop/
  - /language/nodejs/run-tests/
  - /language/nodejs/configure-ci-cd/
  - /language/nodejs/deploy/
params:
  tags: [languages]
  time: 20 minutes
---


[Node.js](https://nodejs.org/en) is a JavaScript runtime for building server-side applications. This guide shows you how to containerize a TypeScript Node.js application using Docker, starting from a simple Express API and progressively adding features like a database and CI/CD.

This guide focuses on a backend Node.js API. If you're building a standalone frontend application, Docker has dedicated guides for [React.js](/guides/reactjs/), [Vue.js](/guides/vuejs/), [Angular](/guides/angular/), and [Next.js](/guides/nextjs/).

> **Acknowledgment**
>
> Docker thanks [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for his contribution to this guide.

## What will you learn?

In this guide, you'll learn how to:

- Containerize and run a Node.js application using Docker.
- Set up a local development environment using containers.
- Run tests inside a Docker container.
- Configure GitHub Actions for CI/CD with Docker.
- Inspect and generate supply chain attestations for your image.
- Deploy your containerized Node.js application to Kubernetes.

Start by containerizing a Node.js application.

## Prerequisites

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) and [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm).
- Familiarity with Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

## Containerize a Node.js application

### Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You're familiar with basic Docker concepts. If you're new to Docker, start with [Get started](/get-started/introduction/).

### Overview

Containerizing your application means packaging it together with its
dependencies, configuration, and runtime into a single portable unit called a
container image. Running that image creates a container, an isolated process
that behaves the same on any machine, whether it's your laptop, a CI runner, or
a production server.

In this section, you'll containerize a simple [Express.js](https://expressjs.com/) API written in TypeScript. You'll write a `Dockerfile` that describes how to build the image, add a `compose.yaml` file that defines how Docker runs your container, and then build and start the application with one command.

You'll use [Docker Hardened Images](/dhi/) as the base. These are minimal, secure Node.js images maintained by Docker.

This guide focuses on a backend Node.js API. If you're building a standalone frontend application, Docker has dedicated guides for [React.js](/guides/reactjs/), [Vue.js](/guides/vuejs/), [Angular](/guides/angular/), and [Next.js](/guides/nextjs/).

### Create the application

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

### Create the Docker assets

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

### Run the application

Inside the `nodejs-docker-example` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You should see `{"message":"Hello World"}`.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

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

### Summary

In this section, you learned how to containerize and run a Node.js application using Docker.

Related information:

- [Docker Hardened Images](/dhi/)
- [Dockerfile reference](/reference/dockerfile.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll take a look at how to set up a local development environment using Docker containers.

## Use containers for Node.js development

### Prerequisites

Complete [Containerize a Node.js application](./).

### Overview

Once your application runs in a container, the next step is making the container part of your everyday development workflow. Code changes should show up quickly, and services your app depends on, like databases, should run right alongside it.

In this section, you'll adapt the Dockerfile for local development by renaming the `builder` stage to `dev` and pointing Compose at it. You'll also update the application to connect to a PostgreSQL database, add a database service to `compose.yaml`, persist data in a named volume, enable Compose Watch so changes in your editor are picked up without a manual rebuild, and set up Node.js debugging so you can attach VS Code or Chrome DevTools to the running container.

### Update the application

You'll update your application to connect to a PostgreSQL database. Continue working in your `nodejs-docker-example` directory.

Replace `src/index.ts` and `package.json` with the following contents. The file browser shows only the files that change in this step.

> [!NOTE]
>
> The application won't run yet after this step. It tries to connect to a
> PostgreSQL database that doesn't exist. The next two sections add the
> database service and the Docker configuration needed to run everything
> together.

{{< files name="nodejs-docker-example" >}}

{{< file path="src/index.ts" status="modified" >}}
```typescript
// Express application backed by a PostgreSQL database.
// Creates a heroes table at startup.
// Endpoints: GET / (greeting), GET /health (health check), POST /heroes/ (create), GET /heroes/ (list).
// See https://expressjs.com/ and https://node-postgres.com/

import express, { type Request, type Response } from 'express';
import { Pool } from 'pg';
import { readFileSync } from 'fs';

const app = express();
const port = parseInt(process.env.PORT ?? '3000', 10);

app.use(express.json());

function getPassword(): string {
  const passwordFile = process.env.POSTGRES_PASSWORD_FILE;
  if (passwordFile) {
    return readFileSync(passwordFile, 'utf8').trim();
  }
  return process.env.POSTGRES_PASSWORD ?? '';
}

const pool = new Pool({
  host: process.env.POSTGRES_SERVER,
  port: 5432,
  database: process.env.POSTGRES_DB,
  user: process.env.POSTGRES_USER,
  password: getPassword(),
});

pool
  .query(
    `CREATE TABLE IF NOT EXISTS heroes (
      id SERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      secret_name TEXT NOT NULL,
      age INTEGER
    )`,
  )
  .catch(console.error);

app.get('/', (_req: Request, res: Response) => {
  res.json({ message: 'Hello World' });
});

app.get('/health', (_req: Request, res: Response) => {
  res.json({ status: 'ok' });
});

app.post('/heroes/', async (req: Request, res: Response) => {
  const { name, secret_name, age } = req.body as {
    name: string;
    secret_name: string;
    age?: number;
  };
  const result = await pool.query(
    'INSERT INTO heroes (name, secret_name, age) VALUES ($1, $2, $3) RETURNING *',
    [name, secret_name, age],
  );
  res.json(result.rows[0]);
});

app.get('/heroes/', async (_req: Request, res: Response) => {
  const result = await pool.query('SELECT * FROM heroes');
  res.json(result.rows);
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
```
{{< /file >}}

{{< file path="package.json" status="modified" hl_lines="13,18" >}}
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
    "express": "^4.21.2",
    "pg": "^8.16.0"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "@types/pg": "^8.11.0",
    "tsx": "^4.19.3",
    "typescript": "^5.8.3"
  }
}
```
{{< /file >}}

{{< /files >}}

### Update Docker assets

Replace `Dockerfile` and `compose.yaml` with the following.

{{< files name="nodejs-docker-example" >}}

{{< file path="Dockerfile" status="modified" hl_lines="12,34,37,63" >}}
```dockerfile
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# This Dockerfile uses Docker Hardened Images (DHI) for enhanced security.
# For more information, see https://docs.docker.com/dhi/

# Development stage: install all dependencies, compile TypeScript, and
# serve with hot-reload. Used directly in development via compose.yaml.
FROM dhi.io/node:24-alpine3.23-dev AS dev

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

# Expose the port that the application listens on.
EXPOSE 3000

# Run the application in development mode.
CMD ["npm", "run", "dev"]


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
COPY --from=dev --chown=node:node /app/dist ./dist

# Expose the port that the application listens on.
EXPOSE 3000

# Run the application.
CMD ["node", "dist/index.js"]
```
{{< /file >}}

{{< file path="compose.yaml" status="modified" hl_lines="8" >}}
```yaml
services:
  # Application service. The `target: dev` line builds the development
  # image (includes tsx and dev tooling); the runner stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: dev
    ports:
      - 3000:3000
```
{{< /file >}}

{{< /files >}}

#### About these changes

The `builder` stage from containerize is renamed to `dev` and gains `EXPOSE 3000` and `CMD ["npm", "run", "dev"]`, which runs `tsx watch` for hot-reload. The `deps` and `runner` stages are unchanged.

In `compose.yaml`, the new `target: dev` line tells Compose to build and run the `dev` stage during development. Unlike the production image, the development image includes `tsx` and other dev tooling. If you need a shell in a running production container, use [Docker Debug](/reference/cli/docker/debug/) instead.

The build step runs `tsc`, which compiles each TypeScript file into a corresponding JavaScript file. [esbuild](https://esbuild.github.io/) is a popular alternative that bundles everything into a single output file and builds significantly faster. To switch, replace the `tsc` call in `package.json` with an esbuild command and update the `COPY --from=dev` path in the `runner` stage to match esbuild's output.

### Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data, and add a `db/password.txt` file that holds the database password.

{{< files name="nodejs-docker-example" >}}

{{< file path="compose.yaml" status="modified" hl_lines="11-46" >}}
```yaml
services:
  # Application service. The `target: dev` line builds the development
  # image (includes tsx and dev tooling); the runner stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: dev
    ports:
      - 3000:3000
    environment:
      - POSTGRES_SERVER=db
      - POSTGRES_USER=postgres
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
  # Database service. Reads the password from a Docker secret mounted at
  # /run/secrets/db-password. Compose waits for the healthcheck to pass
  # before starting the server, via the server's depends_on.
  db:
    image: dhi.io/postgres:18
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
{{< /file >}}

{{< file path="db/password.txt" status="new" >}}
```text
mysecretpassword
```
{{< /file >}}

{{< /files >}}

> [!NOTE]
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/reference/compose-file/).

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

Now test your API endpoint. Open a new terminal and make a request to the server using the curl commands.

Create an object with a POST request:

```console
$ curl -X 'POST' \
  'http://localhost:3000/heroes/' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "my hero",
  "secret_name": "austing",
  "age": 12
}'
```

You should receive the following response:

```json
{
  "id": 1,
  "name": "my hero",
  "secret_name": "austing",
  "age": 12
}
```

Now make a GET request:

```console
$ curl http://localhost:3000/heroes/
```

You should receive the same response because it's the only object in the database.

Press `ctrl`+`c` in the terminal to stop your application.

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you edit and save your code. For more details about Compose Watch, see [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

Open your `compose.yaml` file in an IDE or text editor and add the highlighted Compose Watch instructions.

{{< files name="nodejs-docker-example" >}}

{{< file path="compose.yaml" status="modified" hl_lines="21-27" >}}
```yaml
services:
  # Application service. The `target: dev` line builds the development
  # image (includes tsx and dev tooling); the runner stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: dev
    ports:
      - 3000:3000
    environment:
      - POSTGRES_SERVER=db
      - POSTGRES_USER=postgres
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    develop:
      watch:
        - action: sync
          path: ./src
          target: /app/src
        - action: rebuild
          path: package.json
  db:
    image: dhi.io/postgres:18
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
{{< /file >}}

{{< /files >}}

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

In a terminal, curl the application to get a response.

```console
$ curl http://localhost:3000
{"message":"Hello World"}
```

Any changes to the application's source files on your local machine will now be immediately reflected in the running container.

Open `nodejs-docker-example/src/index.ts` in an IDE or text editor and update the `Hello World` string by adding a few more exclamation marks.

```diff
-  res.json({ message: 'Hello World' });
+  res.json({ message: 'Hello World!!!' });
```

Save the changes to `src/index.ts` and then wait a few seconds for the application to reload. Curl the application again and verify that the updated text appears.

```console
$ curl http://localhost:3000
{"message":"Hello World!!!"}
```

Press `ctrl`+`c` in the terminal to stop your application.

### Debug your application

`tsx watch` supports the Node.js inspector protocol, so you can attach a debugger from VS Code or Chrome DevTools and set breakpoints directly in your TypeScript source files.

Update the `dev` script in `package.json` to start the inspector. The `--inspect=0.0.0.0:9229` flag tells Node.js to listen for debugger connections on all network interfaces at port 9229. Using `0.0.0.0` rather than `localhost` is necessary so the debugger is reachable from outside the container. Also expose the debug port in `compose.yaml`, and add a `.vscode/launch.json` file that tells VS Code how to attach to the running inspector.

{{< files name="nodejs-docker-example" >}}

{{< file path="package.json" status="modified" hl_lines="9" >}}
```json
{
  "name": "nodejs-docker-example",
  "version": "1.0.0",
  "description": "A minimal Node.js TypeScript application.",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch --inspect=0.0.0.0:9229 src/index.ts"
  },
  "dependencies": {
    "express": "^4.21.2",
    "pg": "^8.16.0"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "@types/pg": "^8.11.0",
    "tsx": "^4.19.3",
    "typescript": "^5.8.3"
  }
}
```
{{< /file >}}

{{< file path=".vscode/launch.json" status="new" >}}
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Attach to Docker Container",
      "type": "node",
      "request": "attach",
      "port": 9229,
      "address": "localhost",
      "localRoot": "${workspaceFolder}",
      "remoteRoot": "/app",
      "protocol": "inspector",
      "restart": true,
      "sourceMaps": true,
      "skipFiles": ["<node_internals>/**"]
    }
  ]
}
```
{{< /file >}}

{{< file path="compose.yaml" status="modified" hl_lines="11" >}}
```yaml
services:
  # Application service. The `target: dev` line builds the development
  # image (includes tsx and dev tooling); the runner stage of the
  # Dockerfile is unused in development.
  server:
    build:
      context: .
      target: dev
    ports:
      - 3000:3000
      - 9229:9229
    environment:
      - POSTGRES_SERVER=db
      - POSTGRES_USER=postgres
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    develop:
      watch:
        - action: sync
          path: ./src
          target: /app/src
        - action: rebuild
          path: package.json
  db:
    image: dhi.io/postgres:18
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
{{< /file >}}

{{< /files >}}

Rebuild and restart with the updated configuration:

```console
$ docker compose up --build
```

When the inspector is ready, you'll see a line like the following in the logs:

```text
Debugger listening on ws://0.0.0.0:9229/...
```

#### VS Code

With `.vscode/launch.json` in place, attach the debugger using the Debug panel.

Open the Debug panel (`Ctrl+Shift+D` on Windows and Linux, `Cmd+Shift+D` on Mac), select **Attach to Docker Container**, and press `F5`. You can now set breakpoints in your TypeScript source files under `src/`.

#### Chrome DevTools

You can also use the built-in Node.js inspector in Chrome without any editor setup.

1. Open Chrome and go to `chrome://inspect`.

2. Select **Configure** and add `localhost:9229`.

3. When your Node.js target appears in the list, select **inspect**.

#### Troubleshoot the debugger

If the debugger doesn't connect, verify the container is running and the port is mapped correctly:

```console
$ docker compose ps
$ docker compose logs server
```

The logs should include a line like:

```text
Debugger listening on ws://0.0.0.0:9229/...
```

If that line is missing, confirm the `dev` script in `package.json` includes `--inspect=0.0.0.0:9229` and that `9229:9229` appears in the `ports` list for the `server` service in `compose.yaml`.

For more details about Node.js debugging, see the [Node.js debugging guide](https://nodejs.org/en/docs/guides/debugging-getting-started).

### Summary

In this section, you set up a Compose file with a local database and persistent storage, set up Compose Watch to automatically sync code changes, and configured a debugger that attaches from VS Code and Chrome DevTools.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose secrets](/reference/compose-file/secrets.md)
- [Compose Watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
- [Node.js debugging guide](https://nodejs.org/en/docs/guides/debugging-getting-started)

### Next steps

In the next section, you'll learn how to run tests using Docker.

## Run Node.js tests in a container

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](./).

### Overview

Testing is a core part of building reliable software. Docker makes it easy to
run your tests in the same environment used in CI and production, so failures
are caught before they reach your users.

In this section, you'll add [Vitest](https://vitest.dev/) to the project and
run tests both locally and inside a container.

### Update the application

You'll refactor `src/index.ts` to export the Express `app` instance so tests
can import it without starting a server. Add a test file and update
`package.json` to add Vitest and a test runner for HTTP requests. The file browser shows only the files that change in this step.

{{< files name="nodejs-docker-example" >}}

{{< file path="src/index.ts" status="modified" hl_lines="10,31,70-75" >}}
```typescript
// Express application backed by a PostgreSQL database.
// Creates a heroes table at startup.
// Endpoints: GET / (greeting), GET /health (health check), POST /heroes/ (create), GET /heroes/ (list).
// See https://expressjs.com/ and https://node-postgres.com/

import express, { type Request, type Response } from 'express';
import { Pool } from 'pg';
import { readFileSync } from 'fs';

export const app = express();
const port = parseInt(process.env.PORT ?? '3000', 10);

app.use(express.json());

function getPassword(): string {
  const passwordFile = process.env.POSTGRES_PASSWORD_FILE;
  if (passwordFile) {
    return readFileSync(passwordFile, 'utf8').trim();
  }
  return process.env.POSTGRES_PASSWORD ?? '';
}

const pool = new Pool({
  host: process.env.POSTGRES_SERVER,
  port: 5432,
  database: process.env.POSTGRES_DB,
  user: process.env.POSTGRES_USER,
  password: getPassword(),
});

if (process.env.POSTGRES_SERVER) {
  pool
    .query(
      `CREATE TABLE IF NOT EXISTS heroes (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        secret_name TEXT NOT NULL,
        age INTEGER
      )`,
    )
    .catch(console.error);
}

app.get('/', (_req: Request, res: Response) => {
  res.json({ message: 'Hello World' });
});

app.get('/health', (_req: Request, res: Response) => {
  res.json({ status: 'ok' });
});

app.post('/heroes/', async (req: Request, res: Response) => {
  const { name, secret_name, age } = req.body as {
    name: string;
    secret_name: string;
    age?: number;
  };
  const result = await pool.query(
    'INSERT INTO heroes (name, secret_name, age) VALUES ($1, $2, $3) RETURNING *',
    [name, secret_name, age],
  );
  res.json(result.rows[0]);
});

app.get('/heroes/', async (_req: Request, res: Response) => {
  const result = await pool.query('SELECT * FROM heroes');
  res.json(result.rows);
});

// Only start the server when this file is run directly.
if (require.main === module) {
  app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
  });
}
```
{{< /file >}}

{{< file path="src/index.test.ts" status="new" >}}
```typescript
// Unit tests for the Express application.
// Tests the root endpoint without starting a server.
// See https://vitest.dev/ for the test framework reference.

import { describe, it, expect } from 'vitest';
import request from 'supertest';
import { app } from './index';

describe('GET /', () => {
  it('returns a JSON greeting', async () => {
    const response = await request(app).get('/');
    expect(response.status).toBe(200);
    expect(response.body).toEqual({ message: 'Hello World' });
  });
});
```
{{< /file >}}

{{< file path="package.json" status="modified" hl_lines="10,20-22" >}}
```json
{
  "name": "nodejs-docker-example",
  "version": "1.0.0",
  "description": "A minimal Node.js TypeScript application.",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch src/index.ts",
    "test": "vitest run"
  },
  "dependencies": {
    "express": "^4.21.2",
    "pg": "^8.16.0"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "@types/pg": "^8.11.0",
    "supertest": "^7.0.0",
    "@types/supertest": "^6.0.0",
    "tsx": "^4.19.3",
    "typescript": "^5.8.3",
    "vitest": "^3.0.0"
  }
}
```
{{< /file >}}

{{< /files >}}

### Run tests locally

Run the following command to run the tests locally:

```console
$ npm install
$ npm test
```

You should see output like the following:

```console
 RUN  v3.0.0 /app

 ✓ src/index.test.ts (1)
   ✓ GET / (1)
     ✓ returns a JSON greeting

 Test Files  1 passed (1)
      Tests  1 passed (1)
   Start at  12:00:00
   Duration  500ms
```

### Run tests in a container

Run the tests using the dev stage of your Dockerfile:

```console
$ docker compose run --build --rm --no-deps server npm test
```

The `--no-deps` flag skips starting the database, since the unit tests don't require it. The `--rm` flag removes the container when the tests finish.

You should see the same test output as when running locally.

### Run tests when building

To run tests during the Docker build process, add a `test` stage to your Dockerfile that runs after the dev stage.

```dockerfile {hl_lines="32-36"}
FROM dhi.io/node:24-alpine3.23-dev AS dev

WORKDIR /app

RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm install

COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npm", "run", "dev"]


FROM dhi.io/node:24-alpine3.23-dev AS deps
WORKDIR /app
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm install --omit=dev

FROM dhi.io/node:24-alpine3.23 AS runner
ENV PATH=/app/node_modules/.bin:$PATH
WORKDIR /app
COPY --from=deps --chown=node:node /app/node_modules ./node_modules
COPY --from=dev --chown=node:node /app/dist ./dist

EXPOSE 3000
CMD ["node", "dist/index.js"]


FROM dev AS test

ENV CI=true

CMD ["npm", "test"]
```

Then build and run the test stage:

```console
$ docker build --target test -t nodejs-app-test .
$ docker run --rm nodejs-app-test
```

### Summary

In this section, you learned how to run tests when developing locally and inside a container.

Related information:

- [Dockerfile reference](/reference/dockerfile/)
- [Compose file reference](/compose/compose-file/)
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/)

### Next steps

In the next section, you'll learn how to set up a CI/CD pipeline using GitHub Actions.

## Automate your builds with GitHub Actions

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](./). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

If you haven't created a [GitHub repository](https://github.com/new) for your project yet, do that now. After creating the repository, [add a remote](https://docs.github.com/en/get-started/getting-started-with-git/managing-remote-repositories) and make sure you can commit and [push your code](https://docs.github.com/en/get-started/using-git/pushing-commits-to-a-remote-repository#about-git-push) to GitHub.

1. In your project's GitHub repository, open **Settings**, and go to **Secrets and variables** > **Actions**.

2. Under the **Variables** tab, create a new **Repository variable** named `DOCKER_USERNAME` with your Docker ID as the value.

3. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token) for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

4. Add the PAT as a **Repository secret** in your GitHub repository, with the name `DOCKERHUB_TOKEN`.

### Overview

GitHub Actions is a CI/CD automation tool built into GitHub. A workflow is a YAML file that tells GitHub which jobs to run when something happens in your repository, like a push to a branch or a pull request opening. Workflows live in the `.github/workflows/` directory of your repository.

In this section, you'll add a workflow that runs your tests on every push to the main branch, then builds your Docker image and pushes it to Docker Hub.

### Define the GitHub Actions workflow

You can create a GitHub Actions workflow by creating a YAML file in the `.github/workflows/` directory of your repository. Use your favorite text editor or the GitHub web interface.

If you prefer to use the GitHub web interface:

1. Go to your repository on GitHub and select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub Actions workflow file in your repository. By default, the file is created under `.github/workflows/main.yml`. Change the filename to `build.yml`.

If you prefer to use your text editor, create a new file named `build.yml` in the `.github/workflows/` directory of your repository.

Add the following content to the file:

{{< files name="nodejs-docker-example" >}}

{{< file path=".github/workflows/build.yml" status="new" >}}
```yaml
# GitHub Actions workflow that runs on every push to main.
# - test: runs Vitest unit tests inside a container.
# - build_and_push: signs in to Docker Hub and the DHI registry, then
#   builds and pushes the image.
name: Build and push Docker image

on:
  push:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@{{% param "checkout_action_version" %}}

      - name: Run tests
        run: docker build --target test -t nodejs-test . && docker run --rm nodejs-test

  build_and_push:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@{{% param "checkout_action_version" %}}

      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Login to Docker Hardened Images
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          registry: dhi.io
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
{{< /file >}}

{{< /files >}}

The workflow has two jobs:

1. **test**: Builds the `test` stage of the Dockerfile and runs it. If tests fail, the workflow stops and `build_and_push` doesn't run.
2. **build_and_push**: Signs in to Docker Hub and the DHI registry, then builds and pushes the image.

### Run the workflow

Commit the changes and push them to the `main` branch. This workflow runs every time you push changes to `main`. You can find more information about workflow triggers in the [GitHub documentation](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows).

Go to the **Actions** tab of your GitHub repository. It displays the workflow. Selecting the workflow shows you the breakdown of all the steps.

When the workflow is complete, go to your [repositories on Docker Hub](https://hub.docker.com/repositories). If you see the new repository in that list, the GitHub Actions workflow successfully pushed the image to Docker Hub.

### Summary

In this section, you learned how to set up a GitHub Actions workflow for your Node.js application that includes:

- Running Vitest unit tests inside a container
- Building and pushing Docker images

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [docker/login-action](https://github.com/docker/login-action)
- [docker/build-push-action](https://github.com/docker/build-push-action)
- [Create a Docker Hub access token](/manuals/security/access-tokens.md#create-an-access-token)

### Next steps

In the next section, you'll learn how to inspect and generate supply chain
attestations for your image. See [Secure your supply chain](./).

## Secure your Node.js image supply chain

### Prerequisites

Complete [Automate your builds with GitHub Actions](./).

### Overview

When you ship a container image, what's inside it and where it came from
matters. Supply chain attestations are signed records that answer questions
like which packages are in the image, what vulnerabilities affect them, how
the image was built, and what security checks it passed.

In this section, you'll inspect the attestations that ship with your Docker
Hardened Image base, generate your own SBOM and provenance attestations
during CI, and pin the base image by digest so your builds are reproducible.

The inspection commands in this topic are shown manually so you can see what
each one returns. In a real workflow you'd automate these checks with
[Docker Scout](/scout/), which runs the same scans on every push,
enforces policies in CI, and surfaces results in your registry and pull
requests.

### Inspect the base image attestations

Docker Hardened Images are built to SLSA Build Level 3 and ship with a set of
signed attestations covering bill-of-materials, vulnerabilities, build
provenance, and security scans. See
[DHI attestations](/manuals/dhi/core-concepts/attestations.md) for the full
list of types and how to verify their signatures with Cosign.

List all the attestations available on the Node.js DHI:

```console
$ docker scout attest list registry://dhi.io/node:24-alpine3.23-dev
```

View the SBOM:

```console
$ docker scout sbom registry://dhi.io/node:24-alpine3.23-dev
```

Check known vulnerabilities:

```console
$ docker scout cves registry://dhi.io/node:24-alpine3.23-dev
```

> [!NOTE]
>
> The `registry://` prefix forces `docker scout` to fetch the image and its
> attestations from the registry instead of reading a locally pulled copy. If
> you've already pulled or built against the base image, the local copy
> doesn't have the attached attestations, so the prefix is required to see
> them.

When you base your own image on a DHI image, these attestations stay attached to the base layer in the registry. Tools that inspect your image can follow the chain back to the DHI source.

### Generate attestations for your image

Update your GitHub Actions workflow to attach SBOM and provenance attestations to the image you push.

Edit `.github/workflows/build.yml` and update the build-and-push step:

```yaml {hl_lines="6-7"}
- name: Build and push Docker image
  uses: docker/build-push-action@{{% param "build_push_action_version" %}}
  with:
    context: .
    push: true
    sbom: true
    provenance: mode=max
    tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
```

- `sbom: true` tells BuildKit to scan the built image and attach an SBOM attestation.
- `provenance: mode=max` records detailed build provenance, including the source repository, commit, and build parameters.

The next time your workflow runs, the pushed image will carry these attestations alongside the image manifest in the registry.

### Inspect your pushed image's attestations

After your workflow pushes the image, inspect it the same way you inspected the base image:

```console
$ docker scout attest list registry://DOCKER_USERNAME/REPO_NAME:latest
$ docker scout sbom registry://DOCKER_USERNAME/REPO_NAME:latest
```

The SBOM includes packages from every layer, including those inherited from `dhi.io/node:24-alpine3.23-dev`. The provenance record references the DHI base image by digest, so consumers of your image can trace the build chain back to the DHI source.

### Pin the base image by digest

Image tags like `dhi.io/node:24-alpine3.23-dev` move over time as new patches land. For reproducible builds, pin to an immutable digest.

Look up the digest for each image:

```console
$ docker buildx imagetools inspect dhi.io/node:24-alpine3.23-dev --format "{{ .Manifest.Digest }}"
sha256:2bf01111c7dfe429362f64b3977f0cd6e63ff39023012f88487dec7e83aa26ca
$ docker buildx imagetools inspect dhi.io/node:24-alpine3.23 --format "{{ .Manifest.Digest }}"
sha256:868827fd45c6a01f7f3337ba7ff3f48ebb14da10d8cf3d347f98ded5481317a5
```

Each digest is a 64-character hex string. Update your `Dockerfile` to reference each digest on its corresponding `FROM` line:

```dockerfile
FROM dhi.io/node:24-alpine3.23-dev@sha256:2bf01111c7dfe429362f64b3977f0cd6e63ff39023012f88487dec7e83aa26ca AS dev
# ...
FROM dhi.io/node:24-alpine3.23@sha256:868827fd45c6a01f7f3337ba7ff3f48ebb14da10d8cf3d347f98ded5481317a5 AS runner
```

> [!TIP]
>
> Pinning by digest also pins you to that image's vulnerabilities. Use [Dependabot](https://docs.github.com/en/code-security/dependabot) or [Renovate](https://docs.renovatebot.com/) to automate digest updates so you get a PR when a new patched image is available, with a changelog to review before merging.

### Summary

In this section, you learned how to:

- Inspect the supply chain attestations that ship with the DHI base image, including SBOMs, CVE reports, and build provenance
- Generate SBOM and provenance attestations for your own image in CI
- Pin base images by digest for reproducible builds

Related information:

- [DHI attestations](/manuals/dhi/core-concepts/attestations.md)
- [Verify a Docker Hardened Image](/manuals/dhi/how-to/verify.md)
- [Docker Scout](/scout/)
- [Build attestations](/manuals/build/metadata/attestations/_index.md)

### Next steps

In the next section, you'll deploy your application to Kubernetes.

## Deploy your Node.js application

### Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a Node.js application](./).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

In this section, you'll deploy your containerized Node.js application to a local Kubernetes cluster using Docker Desktop. You'll create a Kubernetes manifest that describes how the application should run, including the application deployment, the PostgreSQL database, and the services that connect them.

### Create a Kubernetes manifest

Create a new file called `nodejs-docker-example-kubernetes.yaml` in your project root:

{{< files name="nodejs-docker-example" >}}

{{< file path="nodejs-docker-example-kubernetes.yaml" status="new" >}}
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: nodejs-docker-example

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: nodejs-docker-example
data:
  POSTGRES_SERVER: 'postgres'
  POSTGRES_DB: 'example'
  POSTGRES_USER: 'postgres'

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: nodejs-docker-example
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: nodejs-docker-example
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: dhi.io/postgres:18
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: POSTGRES_DB
            - name: POSTGRES_USER
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: POSTGRES_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: db-password
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql
          readinessProbe:
            exec:
              command: [pg_isready]
            initialDelaySeconds: 5
            periodSeconds: 5
      volumes:
        - name: postgres-storage
          persistentVolumeClaim:
            claimName: postgres-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: nodejs-docker-example
spec:
  type: ClusterIP
  ports:
    - port: 5432
      targetPort: 5432
  selector:
    app: postgres

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server
  namespace: nodejs-docker-example
spec:
  replicas: 2
  selector:
    matchLabels:
      app: server
  template:
    metadata:
      labels:
        app: server
    spec:
      containers:
        - name: server
          image: DOCKER_USERNAME/nodejs-docker-example:latest
          ports:
            - containerPort: 3000
          env:
            - name: POSTGRES_SERVER
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: POSTGRES_SERVER
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: POSTGRES_DB
            - name: POSTGRES_USER
              valueFrom:
                configMapKeyRef:
                  name: app-config
                  key: POSTGRES_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: db-password
          readinessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 10
            periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: server
  namespace: nodejs-docker-example
spec:
  type: ClusterIP
  ports:
    - port: 3000
      targetPort: 3000
  selector:
    app: server
```
{{< /file >}}

{{< /files >}}

Before applying the manifest, replace `DOCKER_USERNAME` in the `server` deployment's image field with your Docker Hub username.

### Deploy to Kubernetes

Apply the manifest to your local Kubernetes cluster:

```console
$ kubectl apply -f nodejs-docker-example-kubernetes.yaml
```

You should see output confirming that the resources were created:

```console
namespace/nodejs-docker-example created
configmap/app-config created
persistentvolumeclaim/postgres-pvc created
deployment.apps/postgres created
service/postgres created
deployment.apps/server created
service/server created
```

Then create the database secret from your password file:

```console
$ kubectl create secret generic app-secrets \
  --namespace nodejs-docker-example \
  --from-file=db-password=db/password.txt
```

### Verify the deployment

Check that your pods are running:

```console
$ kubectl get pods -n nodejs-docker-example
```

Wait until all pods show `Running` in the STATUS column. Then verify your services:

```console
$ kubectl get services -n nodejs-docker-example
```

### Access the application

Use port forwarding to access the application from your local machine:

```console
$ kubectl port-forward -n nodejs-docker-example service/server 3000:3000
```

Open a new terminal and make a request to the application:

```console
$ curl http://localhost:3000
{"message":"Hello World"}
```

You can also create a hero:

```console
$ curl -X POST http://localhost:3000/heroes/ \
  -H 'Content-Type: application/json' \
  -d '{"name": "my hero", "secret_name": "austing", "age": 12}'
```

### Clean up

When you're done testing, remove the deployment:

```console
$ kubectl delete -f nodejs-docker-example-kubernetes.yaml
```

### Summary

In this section, you deployed your containerized Node.js application to Kubernetes. You created a manifest that defines the application and database deployments, applied it to a local cluster, and verified the application is accessible.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
- [`kubectl` CLI reference](https://kubernetes.io/docs/reference/kubectl/)
- [Kubernetes Deployment resource](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [Kubernetes Service resource](https://kubernetes.io/docs/concepts/services-networking/service/)
