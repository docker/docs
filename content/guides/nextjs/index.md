---
title: Containerize a Next.js application
linkTitle: Next.js
description: Containerize, develop, test, and deploy Next.js apps with Docker and Kubernetes
keywords: getting started, Next.js, next.js, docker, language, Dockerfile, CI/CD, Kubernetes
summary: |
  This guide explains how to containerize Next.js applications, set up
  development and testing in containers, automate builds with GitHub Actions,
  and deploy to Kubernetes.
aliases:
  - /guides/nextjs/configure-github-actions/
  - /guides/nextjs/containerize/
  - /guides/nextjs/deploy/
  - /guides/nextjs/develop/
  - /guides/nextjs/run-tests/
params:
  tags: [languages]
  time: 20 minutes
---


This guide shows you how to containerize a Next.js application using Docker, following best practices for creating efficient, production-ready containers.

[Next.js](https://nextjs.org/) is a React framework that enables server-side
rendering, static site generation, and full-stack capabilities. Docker
provides a consistent containerized environment from development to
production.

> **Acknowledgment**
>
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide and contributing the official [Next.js Docker examples](https://github.com/vercel/next.js/tree/canary/examples/with-docker) to the Vercel Next.js repository, including the standalone and export output examples. As a Docker Captain and experienced engineer, his expertise in Docker, DevOps, and modern web development has made this resource invaluable for the community, helping developers navigate and optimize their Docker workflows.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run a Next.js application using Docker.
- Set up a local development environment for Next.js inside a container. 
- Run tests for your Next.js application within a Docker container.
- Configure a CI/CD pipeline using GitHub Actions for your containerized app.
- Deploy the containerized Next.js application to a local Kubernetes cluster for testing and debugging.

To begin, you'll start by containerizing an existing Next.js application.

---

## Prerequisites

Before you begin, make sure you're familiar with the following:

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) or [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm) for managing dependencies and running scripts.
- Familiarity with [React](https://react.dev/) and [Next.js](https://nextjs.org/) fundamentals.
- Understanding of Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Once you've completed the Next.js getting started modules, you'll be ready to containerize your own Next.js application using the examples and instructions provided in this guide.

## Containerize a Next.js Application

### Prerequisites

Before you begin, make sure the following tools are installed and available on your system:

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

> [!NOTE]
> New to Docker? Start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide to get familiar with key concepts like images, containers, and Dockerfiles.

---

### Overview

This guide walks you through containerizing a Next.js application with Docker.
You'll learn how to create a production-ready Docker image using best
practices that improve performance, security, scalability, and deployment
efficiency.

By the end of this guide, you will:

- Containerize a Next.js application using Docker.
- Create and optimize a Dockerfile for production builds. 
- Use multi-stage builds to minimize image size.
- Leverage Next.js standalone or export output for efficient containerization.
- Follow best practices for building secure and maintainable Docker images. 

---

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following command
to clone the git repository:

```console
$ git clone https://github.com/kristiyan-velkov/docker-nextjs-sample
```
---

### Build the Docker image

Next.js has specific requirements for production deployments. This guide shows two approaches: `standalone` output (Node.js server) and `export` output (static files with Nginx).

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for your project. Ask Gordon to create a Dockerfile, Compose file, and `.dockerignore` tailored to your application.

#### Step 1: Configure Next.js and create the Dockerfile

Before creating a Dockerfile, choose a base image: the [Node.js Official Image](https://hub.docker.com/_/node) or a [Docker Hardened Image (DHI)](https://hub.docker.com/hardened-images/catalog) from the Hardened Image catalog. Choosing DHI gives you a production-ready, lightweight, and secure image. For more information, see [Docker Hardened Images](https://docs.docker.com/dhi/).

> [!IMPORTANT]
> This guide uses stable Node.js LTS image tags that are considered secure when the guide is written. Because new releases and security patches are published regularly, always review the [official Node.js Docker images](https://hub.docker.com/_/node) and select a secure, up-to-date version before building or deploying.

---

##### 1.1 Next.js with standalone output

Standalone output (`output: "standalone"`) makes Next.js build a self-contained output that includes only the files and dependencies needed to run the application. A single `node server.js` can serve the app, which is ideal for Docker and supports server-side rendering, API routes, and incremental static regeneration. For details, see the [Next.js output configuration documentation](https://nextjs.org/docs/app/api-reference/config/next-config-js/output) (including the "standalone" option).

The container runs the Next.js server with Node.js on port 3000.

Configure Next.js — Open or create `next.config.ts` in your project root:

```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
};

export default nextConfig;
```

Choose either a Docker Hardened Image or the Docker Official Image, then create a `Dockerfile` using the content from the selected tab below.

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}

Docker Hardened Images (DHIs) are available for Node.js in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/node). For more information, see the [DHI quickstart](/dhi/get-started/) guide.

1. Sign in to the DHI registry:
   ```console
   $ docker login dhi.io
   ```

2. Pull the Node.js DHI (check the catalog for available versions):
   ```console
   $ docker pull dhi.io/node:24-alpine3.22-dev
   ```

3. Create a file named `Dockerfile` with the following contents. The `FROM` instructions use `dhi.io/node:24-alpine3.22-dev`. Check the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog) for the latest versions and update the image tags as needed for security and compatibility.

    ```dockerfile
    # ============================================
    # Stage 1: Dependencies Installation Stage
    # ============================================

    # IMPORTANT: Docker Hardened Image (DHI) Version Maintenance
    # This Dockerfile uses dhi.io/node. Regularly validate and update to the latest DHI versions in the catalog for security and compatibility.

    FROM dhi.io/node:24-alpine3.22-dev AS dependencies

    # Set working directory
    WORKDIR /app

    # Copy package-related files first to leverage Docker's caching mechanism
    COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* .npmrc* ./

    # Install project dependencies with frozen lockfile for reproducible builds
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

    # ============================================
    # Stage 2: Build Next.js application in standalone mode
    # ============================================

    FROM dhi.io/node:24-alpine3.22-dev AS builder

    # Set working directory
    WORKDIR /app

    # Copy project dependencies from dependencies stage
    COPY --from=dependencies /app/node_modules ./node_modules

    # Copy application source code
    COPY . .

    ENV NODE_ENV=production

    # Next.js collects completely anonymous telemetry data about general usage.
    # Learn more here: https://nextjs.org/telemetry
    # Uncomment the following line in case you want to disable telemetry during the build.
    # ENV NEXT_TELEMETRY_DISABLED=1

    # Build Next.js application
    # If you want to speed up Docker rebuilds, you can cache the build artifacts
    # by adding: --mount=type=cache,target=/app/.next/cache
    # This caches the .next/cache directory across builds, but it also prevents
    # .next/cache/fetch-cache from being included in the final image, meaning
    # cached fetch responses from the build won't be available at runtime.
    RUN if [ -f package-lock.json ]; then \
        npm run build; \
      elif [ -f yarn.lock ]; then \
        corepack enable yarn && yarn build; \
      elif [ -f pnpm-lock.yaml ]; then \
        corepack enable pnpm && pnpm build; \
      else \
        echo "No lockfile found." && exit 1; \
      fi

    # ============================================
    # Stage 3: Run Next.js application
    # ============================================

    FROM dhi.io/node:24-alpine3.22-dev AS runner

    # Set working directory
    WORKDIR /app

    # Set production environment variables
    ENV NODE_ENV=production
    ENV PORT=3000
    ENV HOSTNAME="0.0.0.0"

    # Next.js collects completely anonymous telemetry data about general usage.
    # Learn more here: https://nextjs.org/telemetry
    # Uncomment the following line in case you want to disable telemetry during the run time.
    # ENV NEXT_TELEMETRY_DISABLED=1

    # Copy production assets
    COPY --from=builder --chown=node:node /app/public ./public

    # Set the correct permission for prerender cache
    RUN mkdir .next
    RUN chown node:node .next

    # Automatically leverage output traces to reduce image size
    # https://nextjs.org/docs/advanced-features/output-file-tracing
    COPY --from=builder --chown=node:node /app/.next/standalone ./
    COPY --from=builder --chown=node:node /app/.next/static ./.next/static

    # If you want to persist the fetch cache generated during the build so that
    # cached responses are available immediately on startup, uncomment this line:
    # COPY --from=builder --chown=node:node /app/.next/cache ./.next/cache

    # Switch to non-root user for security best practices
    USER node

    # Expose port 3000 to allow HTTP traffic
    EXPOSE 3000

    # Start Next.js standalone server
    CMD ["node", "server.js"]
    ```

{{< /tab >}}
{{< tab name="Using the Docker Official Image" >}}

Create a file named `Dockerfile` with the following contents (uses `node`):

```dockerfile
  # ============================================
  # Stage 1: Dependencies Installation Stage
  # ============================================

  ARG NODE_VERSION=24.14.0-slim

  FROM node:${NODE_VERSION} AS dependencies

  # Set working directory
  WORKDIR /app

  # Copy package-related files first to leverage Docker's caching mechanism
  COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* .npmrc* ./

  # Install project dependencies with frozen lockfile for reproducible builds
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

  # ============================================
  # Stage 2: Build Next.js application in standalone mode
  # ============================================

  FROM node:${NODE_VERSION} AS builder

  # Set working directory
  WORKDIR /app

  # Copy project dependencies from dependencies stage
  COPY --from=dependencies /app/node_modules ./node_modules

  # Copy application source code
  COPY . .

  ENV NODE_ENV=production

  # Next.js collects completely anonymous telemetry data about general usage.
  # Learn more here: https://nextjs.org/telemetry
  # Uncomment the following line in case you want to disable telemetry during the build.
  # ENV NEXT_TELEMETRY_DISABLED=1

  # Build Next.js application
  # If you want to speed up Docker rebuilds, you can cache the build artifacts
  # by adding: --mount=type=cache,target=/app/.next/cache
  # This caches the .next/cache directory across builds, but it also prevents
  # .next/cache/fetch-cache from being included in the final image, meaning
  # cached fetch responses from the build won't be available at runtime.
  RUN if [ -f package-lock.json ]; then \
      npm run build; \
    elif [ -f yarn.lock ]; then \
      corepack enable yarn && yarn build; \
    elif [ -f pnpm-lock.yaml ]; then \
      corepack enable pnpm && pnpm build; \
    else \
      echo "No lockfile found." && exit 1; \
    fi

  # ============================================
  # Stage 3: Run Next.js application
  # ============================================

  FROM node:${NODE_VERSION} AS runner

  # Set working directory
  WORKDIR /app

  # Set production environment variables
  ENV NODE_ENV=production
  ENV PORT=3000
  ENV HOSTNAME="0.0.0.0"

  # Next.js collects completely anonymous telemetry data about general usage.
  # Learn more here: https://nextjs.org/telemetry
  # Uncomment the following line in case you want to disable telemetry during the run time.
  # ENV NEXT_TELEMETRY_DISABLED=1

  # Copy production assets
  COPY --from=builder --chown=node:node /app/public ./public

  # Set the correct permission for prerender cache
  RUN mkdir .next
  RUN chown node:node .next

  # Automatically leverage output traces to reduce image size
  # https://nextjs.org/docs/advanced-features/output-file-tracing
  COPY --from=builder --chown=node:node /app/.next/standalone ./
  COPY --from=builder --chown=node:node /app/.next/static ./.next/static

  # If you want to persist the fetch cache generated during the build so that
  # cached responses are available immediately on startup, uncomment this line:
  # COPY --from=builder --chown=node:node /app/.next/cache ./.next/cache

  # Switch to non-root user for security best practices
  USER node

  # Expose port 3000 to allow HTTP traffic
  EXPOSE 3000

  # Start Next.js standalone server
  CMD ["node", "server.js"]
```

> [!NOTE]
> This Dockerfile uses three stages: `dependencies`, `builder`, and `runner`. The final image runs `node server.js` and listens on port 3000.

{{< /tab >}}
{{< /tabs >}}

---

##### 1.2 Next.js with export output

Output export (`output: "export"`) makes Next.js build a fully static site at build time. It generates HTML, CSS, and JavaScript into an `out` directory that can be served by any static host or CDN—no Node.js server at runtime. Use this when you don't need server-side rendering or API routes. For details, see the [Next.js output configuration documentation](https://nextjs.org/docs/app/api-reference/config/next-config-js/output).

Configure Next.js — Open  `next.config.ts` in your project root and add the following code:

```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
```

Choose either a Docker Hardened Image or the Docker Official Image, then create a `Dockerfile` using the content from the selected tab below.

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}

Docker Hardened Images (DHIs) are available for Node.js and Nginx in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog). For more information, see the [DHI quickstart](/dhi/get-started/) guide.

1. Sign in to the DHI registry:
   ```console
   $ docker login dhi.io
   ```

2. Pull the Node.js DHI (check the catalog for available versions):
   ```console
   $ docker pull dhi.io/node:24-alpine3.22-dev
   ```

3. Pull the Nginx DHI (check the catalog for available versions):
   ```console
   $ docker pull dhi.io/nginx:1.28.0-alpine3.21-dev
   ```

4. Create a file named `Dockerfile` with the following contents. The `FROM` instructions use Docker Hardened Images: `dhi.io/node:24-alpine3.22-dev` and `dhi.io/nginx:1.28.0-alpine3.21-dev`. Check the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog) for the latest versions and update the image tags as needed for security and compatibility.

    ```dockerfile
    # ============================================
    # Stage 1: Dependencies Installation Stage
    # ============================================

    # IMPORTANT: Docker Hardened Image (DHI) Version Maintenance
    # This Dockerfile uses dhi.io/node and dhi.io/nginx. Regularly validate and update to the latest DHI versions in the catalog for security and compatibility.

    FROM dhi.io/node:24-alpine3.22-dev AS dependencies

    # Set the working directory
    WORKDIR /app

    # Copy package-related files first to leverage Docker's caching mechanism
    COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* .npmrc* ./

    # Install project dependencies with frozen lockfile for reproducible builds
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

    # ============================================
    # Stage 2: Build Next.js Application
    # ============================================

    FROM dhi.io/node:24-alpine3.22-dev AS builder

    # Set the working directory
    WORKDIR /app

    # Copy project dependencies from dependencies stage
    COPY --from=dependencies /app/node_modules ./node_modules

    # Copy application source code
    COPY . .

    ENV NODE_ENV=production

    # Next.js collects completely anonymous telemetry data about general usage.
    # Learn more here: https://nextjs.org/telemetry
    # Uncomment the following line in case you want to disable telemetry during the build.
    # ENV NEXT_TELEMETRY_DISABLED=1

    # Build Next.js application
    RUN --mount=type=cache,target=/app/.next/cache \
      if [ -f package-lock.json ]; then \
        npm run build; \
      elif [ -f yarn.lock ]; then \
        corepack enable yarn && yarn build; \
      elif [ -f pnpm-lock.yaml ]; then \
        corepack enable pnpm && pnpm build; \
      else \
        echo "No lockfile found." && exit 1; \
      fi

    # =========================================
    # Stage 3: Serve Static Files with Nginx
    # =========================================

    FROM dhi.io/nginx:1.28.0-alpine3.21-dev AS runner

    # Set the working directory
    WORKDIR /app

    # Next.js collects completely anonymous telemetry data about general usage.
    # Learn more here: https://nextjs.org/telemetry
    # Uncomment the following line in case you want to disable telemetry during the run time.
    # ENV NEXT_TELEMETRY_DISABLED=1

    # Copy custom Nginx config
    COPY nginx.conf /etc/nginx/nginx.conf

    # Copy the static build output from the build stage to Nginx's default HTML serving directory
    COPY --chown=nginx:nginx --from=builder /app/out /usr/share/nginx/html

    # Non-root user for security best practices
    USER nginx

    # Expose port 8080 to allow HTTP traffic
    EXPOSE 8080

    # Start Nginx directly with custom config
    ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
    CMD ["-g", "daemon off;"]
    ```

{{< /tab >}}
{{< tab name="Using the Docker Official Image" >}}

Create a file named `Dockerfile` with the following contents (uses `node` and `nginxinc/nginx-unprivileged`):

```dockerfile
# ============================================
# Stage 1: Dependencies Installation Stage
# ============================================

ARG NODE_VERSION=24.14.0-slim
ARG NGINXINC_IMAGE_TAG=alpine3.22

FROM node:${NODE_VERSION} AS dependencies

# Set the working directory
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* .npmrc* ./

# Install project dependencies with frozen lockfile for reproducible builds
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

# ============================================
# Stage 2: Build Next.js Application
# ============================================

FROM node:${NODE_VERSION} AS builder

# Set the working directory
WORKDIR /app

# Copy project dependencies from dependencies stage
COPY --from=dependencies /app/node_modules ./node_modules

# Copy application source code
COPY . .

ENV NODE_ENV=production

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
# ENV NEXT_TELEMETRY_DISABLED=1

# Build Next.js application
RUN --mount=type=cache,target=/app/.next/cache \
  if [ -f package-lock.json ]; then \
    npm run build; \
  elif [ -f yarn.lock ]; then \
    corepack enable yarn && yarn build; \
  elif [ -f pnpm-lock.yaml ]; then \
    corepack enable pnpm && pnpm build; \
  else \
    echo "No lockfile found." && exit 1; \
  fi

# =========================================
# Stage 3: Serve Static Files with Nginx
# =========================================

FROM nginxinc/nginx-unprivileged:${NGINXINC_IMAGE_TAG} AS runner

# Set the working directory
WORKDIR /app

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the run time.
# ENV NEXT_TELEMETRY_DISABLED=1

# Copy custom Nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the static build output from the build stage to Nginx's default HTML serving directory
COPY --from=builder /app/out /usr/share/nginx/html

# Non-root user for security best practices
USER nginx

# Expose port 8080 to allow HTTP traffic
EXPOSE 8080

# Start Nginx directly with custom config
ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]
```

> [!NOTE]
> This guide uses [nginx-unprivileged](https://hub.docker.com/r/nginxinc/nginx-unprivileged) instead of the standard Nginx image to run as a non-root user, following security best practices.

{{< /tab >}}
{{< /tabs >}}

1. Create `nginx.conf` (required for export output only) — Create a file named `nginx.conf` in the root of your project:

    ```nginx
    # Minimal Nginx config for static Next.js app
    worker_processes 1;

    # Store PID in /tmp (always writable)
    pid /tmp/nginx.pid;

    events {
        worker_connections 1024;
    }

    http {
        include       /etc/nginx/mime.types;
        default_type  application/octet-stream;

        # Disable logging to avoid permission issues
        access_log off;
        error_log  /dev/stderr;

        # Optimize static file serving
        sendfile        on;
        tcp_nopush      on;
        tcp_nodelay     on;
        keepalive_timeout  65;

        # Gzip compression
        gzip on;
        gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
        gzip_min_length 256;

        server {
            listen       8080;
            server_name  localhost;

            # Serve static files
            root /usr/share/nginx/html;
            index index.html;

            # Handle Next.js static export routing
            # See: https://nextjs.org/docs/app/guides/static-exports#deploying
            location / {
                try_files $uri $uri.html $uri/ =404;
            }

            # This is necessary when `trailingSlash: false` (default).
            # You can omit this when `trailingSlash: true` in next.config.
            # Handles nested routes like /blog/post -> /blog/post.html
            location ~ ^/(.+)/$ {
                rewrite ^/(.+)/$ /$1.html break;
            }

            # Serve Next.js static assets
            location ~ ^/_next/ {
                try_files $uri =404;
                expires 1y;
                add_header Cache-Control "public, immutable";
            }

            # Optional 404 handling
            error_page 404 /404.html;
            location = /404.html {
                internal;
            }
        }
    }
    ```

      > [!NOTE]
      > Export uses port 8080. For more details, see the [Next.js output configuration](https://nextjs.org/docs/app/api-reference/config/next-config-js/output) and [Nginx documentation](https://nginx.org/en/docs/).

#### Step 2: Create the compose.yaml file

Create a file named `compose.yaml` with the following contents:

```yaml {collapse=true,title=compose.yaml}
services:
  server:
    build:
      context: .
    ports:
      - 3000:3000
```

> [!NOTE]
> If using export output (Nginx), change the port mapping to `8080:8080`.

#### Step 3: Create the .dockerignore file

The `.dockerignore` file tells Docker which files and folders to exclude when building the image.


> [!NOTE]
>This helps:
>- Reduce image size  
>- Speed up the build process  
>- Prevent sensitive or unnecessary files (like `.env`, `.git`, or `node_modules`) from being added to the final image.
>
> To learn more, visit the [.dockerignore reference](/reference/dockerfile.md#dockerignore-file).

Create a file named `.dockerignore` with the following contents:

```dockerignore
# Dependencies (installed inside the image, never copy from host)
node_modules/
.pnp/
.pnp.js
.pnpm-store/

# Next.js build output (generated during the image build)
.next/
out/
dist/
build/
.vercel/

# Testing (not needed in the production image)
coverage/
.nyc_output/
__tests__/
__mocks__/
jest/
cypress/
playwright-report/
test-results/
.vitest/

# Environment files (avoid leaking secrets into the build context)
.env
.env*
.env.local
.env.development.local
.env.test.local
.env.production.local

# Debug and log files
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*
lerna-debug.log*
*.log

# IDE and editor files
.vscode/
.idea/
.cursor/
.cursorrules
.copilot/
*.swp
*.swo
*~

# Git
.git/
.gitignore
.gitattributes

# Docker files (reduce build context; not needed inside the image)
Dockerfile*
.dockerignore
docker-compose*.yml
compose*.yaml

# Documentation (not needed in the image)
*.md
docs/

# CI/CD (not needed in the image)
.github/
.gitlab-ci.yml
.travis.yml
.circleci/
Jenkinsfile

# TypeScript and build metadata
*.tsbuildinfo

# Cache and temporary directories
.cache/
.parcel-cache/
.eslintcache
.stylelintcache
.turbo/
.tmp/
.temp/

# Sensitive or dev-only config (optional; omit if your build needs these)
.pem
.editorconfig
.prettierrc*
.eslintrc*
.stylelintrc*
.babelrc*
*.iml

# OS-specific files
.DS_Store
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
Desktop.ini
```

#### Step 4: Build the Next.js application image

With your custom configuration in place, you're now ready to build the Docker image. Use the Dockerfile you created in Step 1 (standalone or export).

The setup includes:

- Multi-stage builds for optimized image size  
- Standalone: Node.js server on port 3000; Export: Nginx serving static files on port 8080  
- Non-root user for enhanced security  
- Proper file permissions and ownership  

After completing the previous steps, your project directory should contain at least the following files (export also requires `nginx.conf`):

```text
├── docker-nextjs-sample/
│ ├── Dockerfile
│ ├── .dockerignore
│ ├── compose.yaml
│ └── next.config.ts
```

Now that your Dockerfile is configured, you can build the Docker image for your Next.js application.

> [!NOTE]
> The `docker build` command packages your application into an image using the instructions in the Dockerfile. It includes all necessary files from the current directory (called the [build context](/build/concepts/context/#what-is-a-build-context)).

Run the following command from the root of your project:

```console
$ docker build --tag nextjs-sample .
```

What this command does:
- Uses the Dockerfile in the current directory (.)
- Packages the application and its dependencies into a Docker image
- Tags the image as nextjs-sample so you can reference it later


#### Step 5: View local images

After building your Docker image, you can check which images are available on your local machine using either the Docker CLI or [Docker Desktop](/manuals/desktop/use-desktop/images.md). Since you're already working in the terminal, let's use the Docker CLI.

To list all locally available Docker images, run the following command:

```console
$ docker images
```

Example Output:

```shell
REPOSITORY                TAG               IMAGE ID       CREATED         SIZE
nextjs-sample             latest            8c5fc80f098e   14 seconds ago   130MB
```

This output provides key details about your images:

- Repository – The name assigned to the image.
- Tag – A version label that helps identify different builds (e.g., latest).
- Image ID – A unique identifier for the image.
- Created – The timestamp indicating when the image was built.
- Size – The total disk space used by the image.

If the build was successful, you should see `nextjs-sample` image listed. 

---

### Run the containerized application

In the previous step, you created a Dockerfile for your Next.js application and built a Docker image using the docker build command. Now it's time to run that image in a container and verify that your application works as expected.

Run the following command in a terminal. Use the port that matches your setup: standalone uses port 3000, export uses port 8080.

```console
$ docker run -p 3000:3000 nextjs-sample
```

For export output, use port 8080 instead:

```console
$ docker run -p 8080:8080 nextjs-sample
```

Open a browser and view the application: [http://localhost:3000](http://localhost:3000) for standalone or [http://localhost:8080](http://localhost:8080) for export. You should see your Next.js web application.

Press `ctrl+c` in the terminal to stop your application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d` option and `--name` to give the container a name so you can stop it later:

```console
$ docker run -d -p 3000:3000 --name nextjs-app nextjs-sample
```

For export output, use port 8080:

```console
$ docker run -d -p 8080:8080 --name nextjs-app nextjs-sample
```

Open a browser and view the application: [http://localhost:3000](http://localhost:3000) for standalone or [http://localhost:8080](http://localhost:8080) for export. You should see your web application.

To confirm that the container is running, use the `docker ps` command:

```console
$ docker ps
```

This will list all active containers along with their ports, names, and status. Look for a container exposing port 3000 (standalone) or 8080 (export).

Example Output:

```shell
CONTAINER ID   IMAGE           COMMAND                  CREATED             STATUS             PORTS                    NAMES
f49b74736a9d   nextjs-sample   "node server.js"         About a minute ago   Up About a minute   0.0.0.0:3000->3000/tcp nextjs-app
```

To stop the application, run:

```console
$ docker stop nextjs-app
```

> [!NOTE]
> For more information about running containers, see the [`docker run` CLI reference](/reference/cli/docker/container/run/) and the [`docker stop` CLI reference](/reference/cli/docker/container/stop/).

---

### Summary

In this guide, you learned how to containerize, build, and run a Next.js application using Docker. By following best practices, you created a secure, optimized, and production-ready setup.

What you accomplished:
- Configured Next.js for either standalone output (Node.js server) or export output (static files with Nginx).
- Added a multi-stage Dockerfile for your chosen approach: standalone (port 3000) or export (port 8080, with `nginx.conf`).
- Created a `.dockerignore` file to exclude unnecessary files and keep the image clean and efficient.
- Built your Docker image using `docker build`.
- Ran the container using `docker run` with the image name `nextjs-sample`, both in the foreground and in detached mode.
- Verified that the app was running by visiting [http://localhost:3000](http://localhost:3000) (standalone) or [http://localhost:8080](http://localhost:8080) (export).
- Learned how to stop the containerized application using `docker stop nextjs-app`.

You now have a fully containerized Next.js application, running in a Docker container, and ready for deployment across any environment with confidence and consistency.

---

### Related resources

Explore official references and best practices to sharpen your Docker workflow:

- [Multi-stage builds](/build/building/multi-stage/) – Learn how to separate build and runtime stages.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.  
- [Build context in Docker](/build/concepts/context/) – Learn how context affects image builds.  
- [Next.js output configuration](https://nextjs.org/docs/app/api-reference/config/next-config-js/output) – Learn about Next.js production optimization (standalone and export).
- [Next.js with Docker (standalone)](https://github.com/vercel/next.js/tree/canary/examples/with-docker) – Official Next.js example: standalone output with Node.js.
- [Next.js with Docker (export)](https://github.com/vercel/next.js/tree/canary/examples/with-docker-export-output) – Official Next.js example: static export with Nginx or serve.
- [`docker build` CLI reference](/reference/cli/docker/image/build/) – Build Docker images from a Dockerfile.
- [`docker images` CLI reference](/reference/cli/docker/image/ls/) – Manage and inspect local Docker images.
- [`docker run` CLI reference](/reference/cli/docker/container/run/) – Run a command in a new container.
- [`docker stop` CLI reference](/reference/cli/docker/container/stop/) – Stop one or more running containers.

---

### Next steps

With your Next.js application now containerized, you're ready to move on to the next step.

In the next section, you'll learn how to develop your application using Docker containers, enabling a consistent, isolated, and reproducible development environment across any machine.

## Use containers for Next.js development

### Prerequisites

Complete [Containerize Next.js application](./).

---

### Overview

In this section, you'll learn how to set up both production and development environments for your containerized Next.js application using Docker Compose. This setup allows you to run a production build using the standalone server and to develop efficiently inside containers using Next.js's built-in hot reloading with Compose Watch.

You'll learn how to:
- Configure separate containers for production and development
- Enable automatic file syncing using Compose Watch in development
- Debug and live-preview your changes in real-time without manual rebuilds

---

### Automatically update services (development mode)

Use Compose Watch to automatically sync source file changes into your
containerized development environment. This automatically syncs file changes
without needing to restart or rebuild containers manually.

### Step 1: Create a development Dockerfile

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


#### Step 2: Update your `compose.yaml` file

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

#### Step 3: Configure Next.js for Docker development

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
│ └── next.config.ts
```

#### Step 4: Start Compose Watch

Run the following command from your project root to start your container in watch mode:

```console
$ docker compose watch nextjs-dev
```

#### Step 5: Test Compose Watch with Next.js

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

### Summary

In this section, you set up a complete development and production workflow for your Next.js application using Docker and Docker Compose.

Here's what you achieved:
- Created a `Dockerfile.dev` to streamline local development with hot reloading  
- Defined separate `nextjs-dev` and `nextjs-prod-standalone` services in your `compose.yaml` file  
- Enabled real-time file syncing using Compose Watch for a smoother development experience  
- Verified that live updates work seamlessly by modifying and previewing a component

With this setup, you can build, run, and iterate on your Next.js app
entirely within containers across environments.

---

### Related resources

Deepen your knowledge and improve your containerized development workflow with these guides:

- [Using Compose Watch](/manuals/compose/how-tos/file-watch.md) – Automatically sync source changes during development  
- [Multi-stage builds](/manuals/build/building/multi-stage.md) – Create efficient, production-ready Docker images  
- [Dockerfile best practices](/build/building/best-practices/) – Write clean, secure, and optimized Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.
- [Docker volumes](/storage/volumes/) – Persist and manage data between container runs  

### Next steps

In the next section, you'll learn how to run unit tests for your Next.js application inside Docker containers. This ensures consistent testing across all environments and removes dependencies on local machine setup.

## Run Next.js tests in a container

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize Next.js application](./).

### Overview

Testing is a critical part of the development process. In this section, you'll learn how to:

- Run unit tests using Vitest (or Jest) inside a Docker container.
- Run lint (e.g. ESLint) inside a Docker container.
- Use Docker Compose to run tests and lint in an isolated, reproducible environment.

The [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample) uses [Vitest](https://vitest.dev/) with [Testing Library](https://testing-library.com/) for component testing. You can use the same setup or follow the alternative Jest configuration later.

---

### Run tests during development

The [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample) already includes lint (ESLint) and sample tests (Vitest, `app/page.test.tsx`) in place. If you're using the sample app, you can skip to **Step 3: Update compose.yaml** and run tests or lint with the commands below. If you're using your own project, follow the install and configuration steps to add the packages and scripts.

The sample includes a test file at:

```text
app/page.test.tsx
```

This file uses Vitest and React Testing Library to verify the behavior of page components.

#### Step 1: Install Vitest and React Testing Library (custom projects)

If you're using a custom project and haven't already added the necessary testing tools, install them by running:

```console
$ npm install --save-dev vitest @vitejs/plugin-react @testing-library/react @testing-library/dom jsdom
```

Then, update the scripts section of your `package.json` file to include:

```json
"scripts": {
  "test": "vitest",
  "test:run": "vitest run"
}
```

For lint, add a `lint` script (and optionally `lint:fix`). For example, with [ESLint](https://eslint.org/):

```json
"scripts": {
  "test": "vitest",
  "test:run": "vitest run",
  "lint": "eslint .",
  "lint:fix": "eslint . --fix"
}
```

The sample project uses `eslint` and `eslint-config-next` for Next.js. Install them in a custom project with:

```console
$ npm install --save-dev eslint eslint-config-next @eslint/eslintrc
```

Create an ESLint config file (e.g. `eslint.config.cjs`) in your project root with Next.js rules and global ignores:

```js
const { defineConfig, globalIgnores } = require("eslint/config");
const { FlatCompat } = require("@eslint/eslintrc");

const compat = new FlatCompat({ baseDirectory: __dirname });

module.exports = defineConfig([
  ...compat.extends(
    "eslint-config-next/core-web-vitals",
    "eslint-config-next/typescript"
  ),
  globalIgnores([
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    "node_modules/**",
    "eslint.config.cjs",
  ]),
]);
```

---

#### Step 2: Configure Vitest (custom projects)

If you're using a custom project, create a `vitest.config.ts` file in your project root (matching the [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample)):

```ts
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: "./vitest.setup.ts",
    globals: true,
  },
});
```

Create a `vitest.setup.ts` file in your project root:

```ts
import "@testing-library/jest-dom/vitest";
```

> [!NOTE]
> Vitest works well with Next.js and provides fast execution and ESM support. For more details, see the [Next.js testing documentation](https://nextjs.org/docs/app/building-your-application/testing) and [Vitest docs](https://vitest.dev/).

#### Step 3: Update compose.yaml

Add `nextjs-test` and `nextjs-lint` services to your `compose.yaml` file. In the sample project these services use the `tools` profile so they don't start with a normal `docker compose up`. Both reuse `Dockerfile.dev` and run the test or lint command:

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

  nextjs-test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    image: nextjs-sample:dev
    container_name: nextjs-sample-test
    command:
      [
        "sh",
        "-c",
        "if [ -f package-lock.json ]; then npm run test:run 2>/dev/null || npm run test -- --run; elif [ -f yarn.lock ]; then yarn test:run 2>/dev/null || yarn test --run; elif [ -f pnpm-lock.yaml ]; then pnpm run test:run; else npm run test -- --run; fi",
      ]
    profiles:
      - tools

  nextjs-lint:
    build:
      context: .
      dockerfile: Dockerfile.dev
    image: nextjs-sample:dev
    container_name: nextjs-sample-lint
    command:
      [
        "sh",
        "-c",
        "if [ -f package-lock.json ]; then npm run lint; elif [ -f yarn.lock ]; then yarn lint; elif [ -f pnpm-lock.yaml ]; then pnpm lint; else npm run lint; fi",
      ]
    profiles:
      - tools
```

The `nextjs-test` and `nextjs-lint` services reuse the same `Dockerfile.dev` used for [development](develop.md) and override the default command to run tests or lint. The `profiles: [tools]` means these services only run when you use the `--profile tools` option.

After completing the previous steps, your project directory should contain:

```text
├── docker-nextjs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── vitest.config.ts
│ ├── vitest.setup.ts
│ └── next.config.ts
```

#### Step 4: Run the tests

To execute your test suite inside the container, run from your project root:

```console
$ docker compose --profile tools run --rm nextjs-test
```

This command will:
- Start the `nextjs-test` service (because of `--profile tools`).
- Run your test script (`test:run` or `test -- --run`) in the same environment as development.
- Remove the container after the tests complete ([`docker compose run --rm`](/reference/cli/docker/compose/run/)).

> [!NOTE]
> For more information about Compose commands and profiles, see the [Compose CLI reference](/reference/cli/docker/compose/).

#### Step 5: Run lint in the container

To run your linter (e.g. ESLint) inside the container, use the `nextjs-lint` service with the same `tools` profile:

```console
$ docker compose --profile tools run --rm nextjs-lint
```

This command will:
- Start the `nextjs-lint` service (because of `--profile tools`).
- Run your lint script (`npm run lint`, `yarn lint`, or `pnpm lint` depending on your lockfile) in the same environment as development.
- Remove the container after lint completes.

Ensure your `package.json` includes a `lint` script. The sample project already has `"lint": "eslint ."` and `"lint:fix": "eslint . --fix"`; for a custom project, add the same and install `eslint` and `eslint-config-next` if needed.

---

### Summary

In this section, you learned how to run unit tests for your Next.js application inside a Docker container using Vitest and Docker Compose.

What you accomplished:
- Installed and configured Vitest and React Testing Library for testing Next.js components.
- Created `nextjs-test` and `nextjs-lint` services in `compose.yaml` (with `tools` profile) to isolate test and lint execution.
- Reused the development `Dockerfile.dev` to ensure consistency between dev, test, and lint environments.
- Ran tests inside the container using `docker compose --profile tools run --rm nextjs-test`.
- Ran lint inside the container using `docker compose --profile tools run --rm nextjs-lint`.
- Ensured reliable, repeatable testing and linting across environments without relying on local machine setup.

---

### Related resources

Explore official references and best practices to sharpen your Docker testing workflow:

- [Dockerfile reference](/reference/dockerfile/) – Understand all Dockerfile instructions and syntax.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.  
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) – Run one-off commands in a service container.
- [Next.js Testing Documentation](https://nextjs.org/docs/app/building-your-application/testing) – Official Next.js testing guide.
---

### Next steps

Next, you'll learn how to set up a CI/CD pipeline using GitHub Actions to automatically build and test your Next.js application in a containerized environment. This ensures your code is validated on every push or pull request, maintaining consistency and reliability across your development workflow.

## Automate your builds with GitHub Actions

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize Next.js application](./).

You must also have:
- A [GitHub](https://github.com/signup) account.
- A verified [Docker Hub](https://hub.docker.com/signup) account.

---

### Overview

In this section, you'll set up a CI/CD pipeline using [GitHub Actions](https://docs.github.com/en/actions) to automatically:

- Build your Next.js application inside a Docker container.
- Run tests in a consistent environment.
- Push the production-ready image to [Docker Hub](https://hub.docker.com).

---

### Integrate GitHub and Docker Hub

To enable GitHub Actions to build and push Docker images, you'll securely
store your Docker Hub credentials in your new GitHub repository.

#### Step 1: Connect your GitHub repository to Docker Hub

1. Create a Personal Access Token (PAT) from [Docker Hub](https://hub.docker.com)
   1. Go to your **Docker Hub account → Account Settings → Security**.
   2. Generate a new Access Token with **Read/Write** permissions.
   3. Name it something like `nextjs-sample`.
   4. Copy and save the token — you'll need it in Step 4.

2. Create a repository in [Docker Hub](https://hub.docker.com/repositories/)
   1. Go to your **Docker Hub account → Create a repository**.
   2. For the Repository Name, use something descriptive — for example: `nextjs-sample`.
   3. Once created, copy and save the repository name — you'll need it in Step 4.

3. Create a new [GitHub repository](https://github.com/new) for your Next.js project

4. Add Docker Hub credentials as GitHub repository secrets

   In your newly created GitHub repository:
   
   1. Navigate to:
   **Settings → Secrets and variables → Actions → New repository secret**.

   2. Add the following secrets:

   | Name              | Value                          |
   |-------------------|--------------------------------|
   | `DOCKER_USERNAME` | Your Docker Hub username       |
   | `DOCKERHUB_TOKEN` | Your Docker Hub access token (created in Step 1)   |
   | `DOCKERHUB_PROJECT_NAME` | Your Docker Project Name (created in Step 2)   |

   These secrets let GitHub Actions authenticate securely with Docker Hub
   during automated workflows.

5. Connect Your Local Project to GitHub

   Link your local project to the GitHub repository you just created by running the following command from your project root:

   ```console
      $ git remote set-url origin https://github.com/{your-username}/{your-repository-name}.git
   ```

   >[!IMPORTANT]
   >Replace `{your-username}` and `{your-repository}` with your actual GitHub username and repository name.

   To confirm that your local project is correctly connected to the remote GitHub repository, run:

   ```console
   $ git remote -v
   ```

   You should see output similar to:

   ```console
   origin  https://github.com/{your-username}/{your-repository-name}.git (fetch)
   origin  https://github.com/{your-username}/{your-repository-name}.git (push)
   ```

   This confirms that your local repository is properly linked and ready to push your source code to GitHub.

6. Push Your Source Code to GitHub

   Follow these steps to commit and push your local project to your GitHub repository:

   1. Stage all files for commit.

      ```console
      $ git add -A
      ```
      This command stages all changes — including new, modified, and deleted files — preparing them for commit.


   2. Commit your changes.

      ```console
      $ git commit -m "Initial commit"
      ```
      This command creates a commit that snapshots the staged changes with a descriptive message.  

   3. Push the code to the `main` branch.

      ```console
      $ git push -u origin main
      ```
      This command pushes your local commits to the `main` branch of the remote GitHub repository and sets the upstream branch.

Once completed, your code will be available on GitHub, and any GitHub Actions workflow you've configured will run automatically.

> [!NOTE]  
> Learn more about the Git commands used in this step:
> - [Git add](https://git-scm.com/docs/git-add) – Stage changes (new, modified, deleted) for commit  
> - [Git commit](https://git-scm.com/docs/git-commit) – Save a snapshot of your staged changes  
> - [Git push](https://git-scm.com/docs/git-push) – Upload local commits to your GitHub repository  
> - [Git remote](https://git-scm.com/docs/git-remote) – View and manage remote repository URLs

---

#### Step 2: Set up the workflow

Now you'll create a GitHub Actions workflow that builds your Docker image, runs tests, and pushes the image to Docker Hub.

1. Go to your repository on GitHub and select the **Actions** tab in the top menu.

2. Select **Set up a workflow yourself** when prompted.

    This opens an inline editor to create a new workflow file. By default, it will be saved to:
   `.github/workflows/main.yml`

   
3. Add the following workflow configuration to the new file:

```yaml
# CI/CD – Next.js Application with Docker
# Builds the app, runs tests in a container, and pushes the production image to Docker Hub.

name: CI/CD – Next.js Application with Docker

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]

jobs:
  build-test-push:
    name: Build, Test and Push Docker Image
    runs-on: ubuntu-latest

    steps:
      # 1. Checkout source code
      - name: Checkout source code
        uses: actions/checkout@v5
        with:
          fetch-depth: 0

      # 2. Set up Docker Buildx
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v4

      # 3. Cache Docker layers
      - name: Cache Docker layers
        uses: actions/cache@v5
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: ${{ runner.os }}-buildx-

      # 4. Cache pnpm dependencies
      - name: Cache pnpm dependencies
        uses: actions/cache@v4
        with:
          path: ~/.local/share/pnpm/store
          key: ${{ runner.os }}-pnpm-${{ hashFiles('**/pnpm-lock.yaml') }}
          restore-keys: ${{ runner.os }}-pnpm-

      # 5. Extract metadata
      - name: Extract metadata
        id: meta
        run: |
          echo "REPO_NAME=${GITHUB_REPOSITORY##*/}" >> "$GITHUB_OUTPUT"
          echo "SHORT_SHA=${GITHUB_SHA::7}" >> "$GITHUB_OUTPUT"

      # 6. Build dev Docker image (for running tests)
      - name: Build Docker image for tests
        uses: docker/build-push-action@v6
        with:
          context: .
          file: Dockerfile.dev
          tags: ${{ steps.meta.outputs.REPO_NAME }}-dev:latest
          load: true
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max

      # 7. Run Vitest tests inside the container
      # Use same package-manager detection as Dockerfile (no corepack at runtime; node user can't write to /usr/local/bin)
      - name: Run tests
        run: |
          docker run --rm \
            --workdir /app \
            --entrypoint "" \
            -e CI=true \
            ${{ steps.meta.outputs.REPO_NAME }}-dev:latest \
            sh -c "if [ -f package-lock.json ]; then npm run test:run; elif [ -f yarn.lock ]; then yarn test:run; elif [ -f pnpm-lock.yaml ]; then pnpm run test:run; else npm run test:run; fi"
        env:
          CI: true
          NODE_ENV: test
        timeout-minutes: 10

      # 8. Log in to Docker Hub (only needed for push)
      - name: Log in to Docker Hub
        if: github.event_name == 'push' && github.ref == 'refs/heads/main'
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      # 9. Build and push production image (only on push to main)
      - name: Build and push production image
        if: github.event_name == 'push' && github.ref == 'refs/heads/main'
        uses: docker/build-push-action@v6
        with:
          context: .
          file: Dockerfile
          push: true
          platforms: linux/amd64,linux/arm64
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:latest
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:${{ steps.meta.outputs.SHORT_SHA }}
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max

```

This workflow performs the following tasks for your Next.js application:
- Triggers on every `push` or `pull request` targeting the `main` branch.
- Builds a development Docker image using `Dockerfile.dev`, optimized for testing.
- Executes unit tests using Jest inside a clean, containerized environment to ensure consistency.
- Halts the workflow immediately if any test fails — enforcing code quality.
- Caches both Docker build layers and npm dependencies for faster CI runs.
- Authenticates securely with Docker Hub using GitHub repository secrets.
- Builds a production-ready image using the `Dockerfile`.
- Tags and pushes the final image to Docker Hub with both `latest` and short SHA tags for traceability.

> [!NOTE]
>  For more information about  `docker/build-push-action`, refer to the [GitHub Action README](https://github.com/docker/build-push-action/blob/master/README.md).

---

#### Step 3: Run the workflow

After you've added your workflow file, it's time to trigger and observe the CI/CD process in action.

1. Commit and push your workflow file

   Select "Commit changes…" in the GitHub editor.

   - This push will automatically trigger the GitHub Actions pipeline.

2. Monitor the workflow execution

   1. Go to the Actions tab in your GitHub repository.
   2. Select the workflow run to follow each step: `build`, `test`, and (if successful) `push`.

3. Verify the Docker image on Docker Hub

   - After a successful workflow run, visit your [Docker Hub repositories](https://hub.docker.com/repositories).
   - You should see a new image under your repository with:
      - Repository name: `${your-repository-name}`
      - Tags include:
         - `latest` – represents the most recent successful build; ideal for quick testing or deployment.
         - `<short-sha>` – a unique identifier based on the commit hash, useful for version tracking, rollbacks, and traceability.

> [!TIP] Protect your main branch
> To maintain code quality and prevent accidental direct pushes, enable branch protection rules:
>  - Navigate to your **GitHub repo → Settings → Branches**.
>  - Under Branch protection rules, select **Add rule**.
>  - Specify `main` as the branch name.
>  - Enable options like:
>     - *Require a pull request before merging*.
>     - *Require status checks to pass before merging*.
>
>  This ensures that only tested and reviewed code is merged into `main` branch.
---

### Summary

In this section, you set up a complete CI/CD pipeline for your containerized Next.js application using GitHub Actions.

Here's what you accomplished:

- Created a new GitHub repository specifically for your project.
- Generated a secure Docker Hub access token and added it to GitHub as a secret.
- Defined a GitHub Actions workflow to:
   - Build your application inside a Docker container.
   - Run tests in a consistent, containerized environment.
   - Push a production-ready image to Docker Hub if tests pass.
- Triggered and verified the workflow execution through GitHub Actions.
- Confirmed that your image was successfully published to Docker Hub.

With this setup, your Next.js application is now ready for automated testing and deployment across environments — increasing confidence, consistency, and team productivity.

---

### Related resources

Deepen your understanding of automation and best practices for containerized apps:

- [Introduction to GitHub Actions](/guides/gha.md) – Learn how GitHub Actions automate your workflows  
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md) – Set up container builds with GitHub Actions  
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions) – Full reference for writing GitHub workflows  
- [Compose file reference](/compose/compose-file/) – Full configuration reference for `compose.yaml`  
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Optimize your image for performance and security  

---

### Next steps

Next, learn how you can locally test and debug your Next.js workloads on Kubernetes before deploying. This helps you ensure your application behaves as expected in a production-like environment, reducing surprises during deployment.

## Test your Next.js deployment

### Prerequisites

Before you begin, make sure you've completed the following:
- Complete all the previous sections of this guide, starting with [Containerize Next.js application](./).
- [Enable Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

> **New to Kubernetes?**  
> Visit the [Kubernetes basics tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/) to get familiar with how clusters, pods, deployments, and services work.

---

### Overview

This section guides you through deploying your containerized Next.js application locally using [Docker Desktop's built-in Kubernetes](/desktop/kubernetes/). Running your app in a local Kubernetes cluster allows you to closely simulate a real production environment, enabling you to test, validate, and debug your workloads with confidence before promoting them to staging or production.

---

### Create a Kubernetes YAML file

Follow these steps to define your deployment configuration:

1. In the root of your project, create a new file named: nextjs-sample-kubernetes.yaml

2. Open the file in your IDE or preferred text editor.

3. Add the following configuration, and be sure to replace `{DOCKER_USERNAME}` and `{DOCKERHUB_PROJECT_NAME}` with your actual Docker Hub username and repository name from the previous [Automate your builds with GitHub Actions](./).


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nextjs-sample
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nextjs-sample
  template:
    metadata:
      labels:
        app: nextjs-sample
    spec:
      containers:
        - name: nextjs-container
          image: {DOCKER_USERNAME}/{DOCKERHUB_PROJECT_NAME}:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 3000
          env:
            - name: NODE_ENV
              value: "production"
            - name: HOSTNAME
              value: "0.0.0.0"
---
apiVersion: v1
kind: Service
metadata:
  name:  nextjs-sample-service
  namespace: default
spec:
  type: NodePort
  selector:
    app:  nextjs-sample
  ports:
    - port: 3000
      targetPort: 3000
      nodePort: 30001
```

This manifest defines two key Kubernetes resources, separated by `---`:

- Deployment
  Deploys a single replica of your Next.js application inside a pod. The pod uses the Docker image built and pushed by your GitHub Actions CI/CD workflow  
  (refer to [Automate your builds with GitHub Actions](./)).  
  The container listens on port `3000`, which is the default port for Next.js applications.

- Service (NodePort) 
  Exposes the deployed pod to your local machine.  
  It forwards traffic from port `30001` on your host to port `3000` inside the container.  
  This lets you access the application in your browser at [http://localhost:30001](http://localhost:30001).

> [!NOTE]
> To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

---

### Deploy and check your application

Follow these steps to deploy your containerized Next.js app into a local Kubernetes cluster and verify that it's running correctly.

#### Step 1. Apply the Kubernetes configuration

In your terminal, navigate to the directory where your `nextjs-sample-kubernetes.yaml` file is located, then deploy the resources using:

```console
  $ kubectl apply -f nextjs-sample-kubernetes.yaml
```

If everything is configured properly, you'll see confirmation that both the Deployment and the Service were created:

```shell
  deployment.apps/nextjs-sample created
  service/nextjs-sample-service created
```
   
This output means that both the Deployment and the Service were successfully created and are now running inside your local cluster.

#### Step 2. Check the deployment status

Run the following command to check the status of your deployment:
   
```console
  $ kubectl get deployments
```

You should see an output similar to:

```shell
  NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
  nextjs-sample        1/1     1            1           14s
```

This confirms that your pod is up and running with one replica available.

#### Step 3. Verify the service exposure

Check if the NodePort service is exposing your app to your local machine:

```console
$ kubectl get services
```

You should see something like:

```shell
NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
nextjs-sample-service    NodePort    10.100.244.65    <none>        3000:30001/TCP   1m
```

This output confirms that your app is available via NodePort on port 30001.

#### Step 4. Access your app in the browser

Open your browser and navigate to [http://localhost:30001](http://localhost:30001).

You should see your production-ready Next.js Sample application running — served by your local Kubernetes cluster.

#### Step 5. Clean up Kubernetes resources

Once you're done testing, you can delete the deployment and service using:

```console
  $ kubectl delete -f nextjs-sample-kubernetes.yaml
```

Expected output:

```shell
  deployment.apps "nextjs-sample" deleted
  service "nextjs-sample-service" deleted
```

This ensures your cluster stays clean and ready for the next deployment.
   
---

### Summary

In this section, you learned how to deploy your Next.js application to a local Kubernetes cluster using Docker Desktop. This setup allows you to test and debug your containerized app in a production-like environment before deploying it to the cloud.

What you accomplished:

- Created a Kubernetes Deployment and NodePort Service for your Next.js app  
- Used `kubectl apply` to deploy the application locally  
- Verified the app was running and accessible at `http://localhost:30001`  
- Cleaned up your Kubernetes resources after testing

---

### Related resources

Explore official references and best practices to sharpen your Kubernetes deployment workflow:

- [Kubernetes documentation](https://kubernetes.io/docs/home/) – Learn about core concepts, workloads, services, and more.  
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md) – Use Docker Desktop's built-in Kubernetes support for local testing and development.
- [`kubectl` CLI reference](https://kubernetes.io/docs/reference/kubectl/) – Manage Kubernetes clusters from the command line.  
- [Kubernetes Deployment resource](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) – Understand how to manage and scale applications using Deployments.  
- [Kubernetes Service resource](https://kubernetes.io/docs/concepts/services-networking/service/) – Learn how to expose your application to internal and external traffic.
