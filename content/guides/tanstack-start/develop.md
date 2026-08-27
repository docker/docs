---
title: Use containers for TanStack Start development
linkTitle: Develop your app
weight: 30
keywords: tanstack start, development, node
description: Learn how to develop your TanStack Start application locally using containers.
---

## Prerequisites

Complete [Containerize TanStack Start application](containerize.md).

---

## Overview

In this section, you'll set up production and development environments for your
containerized TanStack Start application using Docker Compose. This setup lets
you run a production build with the Nitro server and develop inside containers
using Vite's dev server with Compose Watch.

You'll learn how to:

- Configure separate containers for production and development
- Enable automatic file syncing using Compose Watch in development
- Debug and live-preview your changes without manual rebuilds

---

## Automatically update services (development mode)

Use Compose Watch to automatically sync source file changes into your
containerized development environment. File changes sync without needing to
restart or rebuild containers manually.

## Step 1: Create a development Dockerfile

Create a file named `Dockerfile.dev` in your project root with the following
content (matching the [sample project](https://github.com/kristiyan-velkov/docker-tanstack-start-sample)):

```dockerfile
# =========================================
# Development Dockerfile for TanStack Start
# =========================================
ARG NODE_VERSION=24.14.0-alpine

FROM node:${NODE_VERSION} AS dev

WORKDIR /app

COPY package.json package-lock.json* ./

RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .

ENV HOST=0.0.0.0

EXPOSE 3000

CMD ["npm", "run", "dev"]
```

This file sets up a development environment that runs `npm run dev`, which
starts the Vite dev server on port 3000.

### Step 2: Update your `compose.yml` file

Open your `compose.yml` file and define two services: one for production
(`tanstack-start-prod`) and one for development (`tanstack-start-dev`). This
matches the [sample project](https://github.com/kristiyan-velkov/docker-tanstack-start-sample) structure.

```yaml
services:
  tanstack-start-prod:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        NODE_VERSION: 24.14.0-alpine
    image: tanstack-start:prod
    container_name: tanstack-start-prod
    environment:
      NODE_ENV: production
      PORT: 3000
      HOST: 0.0.0.0
    ports:
      - "3000:3000"

  tanstack-start-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
      args:
        NODE_VERSION: 24.14.0-alpine
    image: tanstack-start:dev
    container_name: tanstack-start-dev
    ports:
      - "3000:3000"
    develop:
      watch:
        - action: sync
          path: .
          target: /app
          ignore:
            - node_modules/
            - .output/
        - action: rebuild
          path: package.json
```

- The `tanstack-start-prod` service builds and runs your production TanStack
  Start app from the `.output` bundle.
- The `tanstack-start-dev` service runs the Vite dev server with hot module
  replacement.
- `watch` triggers file sync with Compose Watch.
- The `rebuild` action for `package.json` reinstalls dependencies when the file
  changes.

> [!NOTE]
> For more details, see the official guide:
> [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

### Step 3: Update vite.config.ts for Docker development

To make Vite's development server reachable from outside the container, add
`server` options to your `vite.config.ts`:

```ts {hl_lines="15-19",linenos=true}
import { defineConfig } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import viteTsConfigPaths from "vite-tsconfig-paths";
import tailwindcss from "@tailwindcss/vite";
import { nitro } from "nitro/vite";

const config = defineConfig({
  plugins: [
    devtools(),
    nitro(),
    viteTsConfigPaths({ projects: ["./tsconfig.json"] }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
  server: {
    host: true,
    port: 3000,
    strictPort: true,
  },
});

export default config;
```

> [!NOTE]
> The `server` options are required for running Vite inside Docker:
>
> - `host: true` lets the dev server accept connections from outside the
>   container.
> - `port: 3000` matches the port exposed in Docker and the sample
>   `package.json` dev script.
> - `strictPort: true` fails clearly if the port is unavailable.
>
> For full details, see the
> [Vite server configuration docs](https://vitejs.dev/config/server-options.html).

After completing the previous steps, your project directory should contain:

```text
├── docker-tanstack-start-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yml
│ └── vite.config.ts
```

### Step 4: Start Compose Watch

Run the following command from your project root:

```console
$ docker compose watch tanstack-start-dev
```

### Step 5: Test Compose Watch with TanStack Start

To verify that Compose Watch is working:

1. Open a route file under `src/routes/` in your text editor.
2. Update visible text in the page component.
3. Save the file.
4. Open your browser at [http://localhost:3000](http://localhost:3000).

You should see the updated content without rebuilding the container manually.

---

## Summary

In this section, you set up development and production workflows for your
TanStack Start application using Docker and Docker Compose.

Here's what you achieved:

- Created a `Dockerfile.dev` for local development with hot reloading
- Defined separate `tanstack-start-dev` and `tanstack-start-prod` services in
  `compose.yml`
- Enabled file syncing using Compose Watch
- Verified live updates by modifying a route component

With this setup, you can build, run, and iterate on your TanStack Start app
entirely within containers across environments.

---

## Related resources

- [Using Compose Watch](/manuals/compose/how-tos/file-watch.md) – Automatically
  sync source changes during development
- [Multi-stage builds](/manuals/build/building/multi-stage/) – Create
  production-ready Docker images
- [Dockerfile best practices](/build/building/best-practices/) – Write clean,
  secure, and optimized Dockerfiles
- [Compose file reference](/compose/compose-file/) – Configure services in
  `compose.yml`

## Next steps

In the next section, you'll learn how to run unit tests for your TanStack Start
application inside Docker containers.
