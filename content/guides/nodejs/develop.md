---
title: Use containers for Node.js development
linkTitle: Develop your app
weight: 20
keywords: node, node.js, development
description: Learn how to develop your Node.js application locally using containers.
aliases:
  - /get-started/nodejs/develop/
  - /language/nodejs/develop/
  - /guides/language/nodejs/develop/
---

## Prerequisites

Complete [Containerize a Node.js application](containerize.md).

## Overview

Once your application runs in a container, the next step is making the container part of your everyday development workflow. Code changes should show up quickly, and services your app depends on, like databases, should run right alongside it.

In this section, you'll adapt the Dockerfile for local development by renaming the `builder` stage to `dev` and pointing Compose at it. You'll also update the application to connect to a PostgreSQL database, add a database service to `compose.yaml`, persist data in a named volume, enable Compose Watch so changes in your editor are picked up without a manual rebuild, and set up Node.js debugging so you can attach VS Code or Chrome DevTools to the running container.

## Update the application

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

## Update Docker assets

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

### About these changes

The `builder` stage from containerize is renamed to `dev` and gains `EXPOSE 3000` and `CMD ["npm", "run", "dev"]`, which runs `tsx watch` for hot-reload. The `deps` and `runner` stages are unchanged.

In `compose.yaml`, the new `target: dev` line tells Compose to build and run the `dev` stage during development. Unlike the production image, the development image includes `tsx` and other dev tooling. If you need a shell in a running production container, use [Docker Debug](/reference/cli/docker/debug/) instead.

The build step runs `tsc`, which compiles each TypeScript file into a corresponding JavaScript file. [esbuild](https://esbuild.github.io/) is a popular alternative that bundles everything into a single output file and builds significantly faster. To switch, replace the `tsc` call in `package.json` with an esbuild command and update the `COPY --from=dev` path in the `runner` stage to match esbuild's output.

## Add a local database and persist data

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

## Automatically update services

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

## Debug your application

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

### VS Code

With `.vscode/launch.json` in place, attach the debugger using the Debug panel.

Open the Debug panel (`Ctrl+Shift+D` on Windows and Linux, `Cmd+Shift+D` on Mac), select **Attach to Docker Container**, and press `F5`. You can now set breakpoints in your TypeScript source files under `src/`.

### Chrome DevTools

You can also use the built-in Node.js inspector in Chrome without any editor setup.

1. Open Chrome and go to `chrome://inspect`.

2. Select **Configure** and add `localhost:9229`.

3. When your Node.js target appears in the list, select **inspect**.

### Troubleshoot the debugger

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

## Summary

In this section, you set up a Compose file with a local database and persistent storage, set up Compose Watch to automatically sync code changes, and configured a debugger that attaches from VS Code and Chrome DevTools.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose secrets](/reference/compose-file/secrets.md)
- [Compose Watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
- [Node.js debugging guide](https://nodejs.org/en/docs/guides/debugging-getting-started)

## Next steps

In the next section, you'll learn how to run tests using Docker.
