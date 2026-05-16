---
title: Use containers for Next.js development
linkTitle: Develop your app
weight: 30
keywords: next.js, development, node
description: Learn how to develop your Next.js application locally using containers.

---

## Prerequisites

Complete [Containerize Next.js application](containerize.md).

---

## Overview

In this section, you'll learn how to set up both production and development environments for your containerized Next.js application using Docker Compose. This setup allows you to run a production build using the standalone server and to develop efficiently inside containers using Next.js's built-in hot reloading with Compose Watch.

You'll learn how to:
- Configure separate containers for production and development
- Enable automatic file syncing using Compose Watch in development
- Debug and live-preview your changes in real-time without manual rebuilds

---

## Automatically update services (Development Mode)

Use Compose Watch to automatically sync source file changes into your
containerized development environment. This automatically syncs file changes
without needing to restart or rebuild containers manually.

## Step 1: Create a development Dockerfile

Create a file named `Dockerfile.dev` in your project root with the following content (matching the [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample)):

```dockerfile
# ============================================
# Development Dockerfile for Next.js
# ============================================
ARG NODE_VERSION=24.14.0-slim

FROM node:${NODE_VERSION} AS dev

WORKDIR /app

COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* .npmrc* ./

RUN --mount=type=cache,target=/root/.npm \
    --mount=type=cache,target=/usr/local/share/.cache/yarn \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
  if [ -f package-lock.json ]; then \
    npm ci --no-audit --no-fund; \
  elif [ -f yarn.lock ]; then \
    corepack enable yarn && yarn install --frozen-lockfile --production=false; \
  elif [ -f pnpm-lock.yaml ]; then \
    corepack enable pnpm && pnpm install --frozen-lockfile; \
  else \
    echo "No lockfile found." && exit 1; \
  fi

COPY . .

ENV WATCHPACK_POLLING=true
ENV HOSTNAME="0.0.0.0"

RUN chown -R node:node /app
USER node

EXPOSE 3000

CMD ["sh", "-c", "if [ -f package-lock.json ]; then npm run dev; elif [ -f yarn.lock ]; then yarn dev; elif [ -f pnpm-lock.yaml ]; then pnpm dev; else npm run dev; fi"]
```

This file sets up a development environment for your Next.js app with hot module replacement and supports npm, yarn, and pnpm.


### Step 2: Update your `compose.yaml` file

Open your `compose.yaml` file and define two services: one for production (`nextjs-prod-standalone`) and one for development (`nextjs-dev`). This matches the [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample) structure.

Here's an example configuration for a Next.js application:

```yaml
services:
  nextjs-prod-standalone:
    build:
      context: .
      dockerfile: Dockerfile
    image: nextjs-sample:prod
    container_name: nextjs-sample-prod
    ports:
      - "3000:3000"

  nextjs-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    image: nextjs-sample:dev
    container_name: nextjs-sample-dev
    ports:
      - "3000:3000"
    environment:
      - WATCHPACK_POLLING=true
    develop:
      watch:
        - action: sync
          path: .
          target: /app
          ignore:
            - node_modules/
            - .next/
        - action: rebuild
          path: package.json
```

- The `nextjs-prod-standalone` service builds and runs your production Next.js app using the standalone output.
- The `nextjs-dev` service runs your Next.js development server with hot module replacement.
- `watch` triggers file sync with Compose Watch.
- `WATCHPACK_POLLING=true` ensures file changes are detected properly inside Docker.
- The `rebuild` action for `package.json` ensures dependencies are reinstalled when the file changes.

> [!NOTE]
> For more details, see the official guide: [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

### Step 3: Configure Next.js for Docker development

Next.js works well inside Docker containers out of the box, but there are a few configurations that can improve the development experience.

The `next.config.ts` file you created during containerization already includes the `output: "standalone"` option for production. For development, Next.js automatically uses its built-in development server with hot reloading enabled.

> [!NOTE]
> The Next.js development server automatically:
> - Enables Hot Module Replacement (HMR) for instant updates
> - Watches for file changes and recompiles automatically
> - Provides detailed error messages in the browser
>
> The `WATCHPACK_POLLING=true` environment variable in the compose file ensures file watching works correctly inside Docker containers.


After completing the previous steps, your project directory should now contain the following files:

```text
├── docker-nextjs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── next.config.ts
│ └── README.Docker.md
```

### Step 4: Start Compose Watch

Run the following command from your project root to start your container in watch mode:

```console
$ docker compose watch nextjs-dev
```

### Step 5: Test Compose Watch with Next.js

To verify that Compose Watch is working correctly:

1. Open the `app/page.tsx` file in your text editor (or `src/app/page.tsx` if your project uses a `src` directory).

2. Locate the main content area and find a text element to modify.

3. Make a visible change, for example, update a heading:

    ```tsx
    <h1>Hello from Docker Compose Watch!</h1>
    ```

4. Save the file.

5. Open your browser at [http://localhost:3000](http://localhost:3000).

You should see the updated text appear instantly, without needing to rebuild the container manually. This confirms that file watching and automatic synchronization are working as expected.

---

## Summary

In this section, you set up a complete development and production workflow for your Next.js application using Docker and Docker Compose.

Here's what you achieved:
- Created a `Dockerfile.dev` to streamline local development with hot reloading  
- Defined separate `nextjs-dev` and `nextjs-prod-standalone` services in your `compose.yaml` file  
- Enabled real-time file syncing using Compose Watch for a smoother development experience  
- Verified that live updates work seamlessly by modifying and previewing a component

With this setup, you can build, run, and iterate on your Next.js app
entirely within containers across environments.

---

## Related resources

Deepen your knowledge and improve your containerized development workflow with these guides:

- [Using Compose Watch](/manuals/compose/how-tos/file-watch.md) – Automatically sync source changes during development  
- [Multi-stage builds](/manuals/build/building/multi-stage.md) – Create efficient, production-ready Docker images  
- [Dockerfile best practices](/build/building/best-practices/) – Write clean, secure, and optimized Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.
- [Docker volumes](/storage/volumes/) – Persist and manage data between container runs  

## Next steps

In the next section, you'll learn how to run unit tests for your Next.js application inside Docker containers. This ensures consistent testing across all environments and removes dependencies on local machine setup.
