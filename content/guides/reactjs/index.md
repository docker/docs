---
title: React.js language-specific guide
linkTitle: React.js
description: Containerize and develop React.js apps using Docker
keywords: getting started, React.js, react.js, docker, language, Dockerfile
summary: |
  This guide explains how to containerize React.js applications using Docker.
aliases:
  - /guides/reactjs/configure-github-actions/
  - /guides/reactjs/containerize/
  - /guides/reactjs/deploy/
  - /guides/reactjs/develop/
  - /guides/reactjs/run-tests/
params:
  tags: [languages]
  time: 20 minutes
---


The React.js language-specific guide shows you how to containerize a React.js application using Docker, following best practices for creating efficient, production-ready containers.

[React.js](https://react.dev/) is a widely used library for building interactive user interfaces. However, managing dependencies, environments, and deployments efficiently can be complex. Docker simplifies this process by providing a consistent and containerized environment.

> 
> **Acknowledgment**
>
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide. As a Docker Captain and experienced Front-end engineer, his expertise in Docker, DevOps, and modern web development has made this resource invaluable for the community, helping developers navigate and optimize their Docker workflows.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run a React.js application using Docker.
- Set up a local development environment for React.js inside a container. 
- Run tests for your React.js application within a Docker container.
- Configure a CI/CD pipeline using GitHub Actions for your containerized app.
- Deploy the containerized React.js application to a local Kubernetes cluster for testing and debugging.

To begin, you’ll start by containerizing an existing React.js application.

---

## Prerequisites

Before you begin, make sure you're familiar with the following:

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) or [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm) for managing dependencies and running scripts.
- Familiarity with [React.js](https://react.dev/) fundamentals.
- Understanding of Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Once you've completed the React.js getting started modules, you’ll be ready to containerize your own React.js application using the examples and instructions provided in this guide.

## Containerize a React.js Application

### Prerequisites

Before you begin, make sure the following tools are installed and available on your system:

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

> **New to Docker?**  
> Start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide to get familiar with key concepts like images, containers, and Dockerfiles.

---

### Overview

This guide walks you through the complete process of containerizing a React.js application with Docker. You’ll learn how to create a production-ready Docker image using best practices that improve performance, security, scalability, and deployment efficiency.

By the end of this guide, you will:

- Containerize a React.js application using Docker.
- Create and optimize a Dockerfile for production builds. 
- Use multi-stage builds to minimize image size.
- Serve the application efficiently with a custom NGINX configuration.
- Follow best practices for building secure and maintainable Docker images. 

---

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following command
to clone the git repository:

```console
$ git clone https://github.com/kristiyan-velkov/docker-reactjs-sample
```
---

### Build the Docker image

React.js is a front-end library that compiles into static assets, so the Dockerfile is tailored to optimize how React applications are built and served in a production environment.

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for your project. Ask Gordon to create a Dockerfile, Compose file, and `.dockerignore` tailored to your application.

#### Step 1: Create the Dockerfile

Before creating a Dockerfile, you need to choose a base image. You can either use the [Node.js Official Image](https://hub.docker.com/_/node) or a Docker Hardened Image (DHI) from the [Hardened Image catalog](https://hub.docker.com/hardened-images/catalog).

Choosing DHI offers the advantage of a production-ready image that is lightweight and secure. For more information, see [Docker Hardened Images](https://docs.docker.com/dhi/).

> [!IMPORTANT]
> This guide uses a stable Node.js LTS image tag that is considered secure when the guide is written. Because new releases and security patches are published regularly, the tag shown here may no longer be the safest option when you follow the guide. Always review the latest available image tags and select a secure, up-to-date version before building or deploying your application.
>
> Official Node.js Docker Images: https://hub.docker.com/_/node

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}
Docker Hardened Images (DHIs) are available for Node.js in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/node). Docker Hardened Images are freely available to everyone with no subscription required. You can pull and use them like any other Docker image after signing in to the DHI registry. For more information, see the [DHI quickstart](/dhi/get-started/) guide.

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

In the following Dockerfile, the `FROM` instructions use `dhi.io/node:24-alpine3.22-dev` and `dhi.io/nginx:1.28.0-alpine3.21-dev` as the base images.

```dockerfile
# =========================================
# Stage 1: Build the React.js Application
# =========================================

# Use a lightweight Node.js image for building (customizable via ARG)
FROM dhi.io/node:24-alpine3.22-dev AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json package-lock.json* ./

# Install project dependencies using npm ci (ensures a clean, reproducible install)
RUN --mount=type=cache,target=/root/.npm npm ci

# Copy the rest of the application source code into the container
COPY . .

# Build the React.js application (outputs to /app/dist)
RUN npm run build

# =========================================
# Stage 2: Prepare Nginx to Serve Static Files
# =========================================

FROM dhi.io/nginx:1.28.0-alpine3.21-dev AS runner

# Copy custom Nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the static build output from the build stage to Nginx's default HTML serving directory
COPY --chown=nginx:nginx --from=builder /app/dist /usr/share/nginx/html

# Use a non-root user for security best practices
USER nginx

# Expose port 8080 to allow HTTP traffic
# Note: The default NGINX container now listens on port 8080 instead of 80 
EXPOSE 8080

# Start Nginx directly with custom config
ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]
```

{{< /tab >}}
{{< tab name="Using the Docker Official Image" >}}

Create a file named `Dockerfile` with the following contents:

```dockerfile
# =========================================
# Stage 1: Build the React.js Application
# =========================================
ARG NODE_VERSION=24.12.0-alpine
ARG NGINX_VERSION=alpine3.22

# Use a lightweight Node.js image for building (customizable via ARG)
FROM node:${NODE_VERSION} AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json package-lock.json* ./

# Install project dependencies using npm ci (ensures a clean, reproducible install)
RUN --mount=type=cache,target=/root/.npm npm ci

# Copy the rest of the application source code into the container
COPY . .

# Build the React.js application (outputs to /app/dist)
RUN npm run build

# =========================================
# Stage 2: Prepare Nginx to Serve Static Files
# =========================================

FROM nginxinc/nginx-unprivileged:${NGINX_VERSION} AS runner

# Copy custom Nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the static build output from the build stage to Nginx's default HTML serving directory
COPY --chown=nginx:nginx --from=builder /app/dist /usr/share/nginx/html

# Use a built-in non-root user for security best practices
USER nginx

# Expose port 8080 to allow HTTP traffic
# Note: The default NGINX container now listens on port 8080 instead of 80 
EXPOSE 8080

# Start Nginx directly with custom config
ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]
```

> [!NOTE]
> We are using nginx-unprivileged instead of the standard NGINX image to follow security best practices.
> Running as a non-root user in the final image:
>- Reduces the attack surface
>- Aligns with Docker’s recommendations for container hardening
>- Helps comply with stricter security policies in production environments

{{< /tab >}}
{{< /tabs >}}

#### Step 2: Create the compose.yaml file

Create a file named `compose.yaml` with the following contents:

```yaml {collapse=true,title=compose.yaml}
services:
  server:
    build:
      context: .
    ports:
      - 8080:8080
```

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
# Ignore dependencies and build output
node_modules/
dist/
out/
.tmp/
.cache/

# Ignore Vite, Webpack, and React-specific build artifacts
.vite/
.vitepress/
.eslintcache
.npm/
coverage/
jest/
cypress/
cypress/screenshots/
cypress/videos/
reports/

# Ignore environment and config files (sensitive data)
*.env*
*.log

# Ignore TypeScript build artifacts (if using TypeScript)
*.tsbuildinfo

# Ignore lockfiles (optional if using Docker for package installation)
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

# Ignore local development files
.git/
.gitignore
.vscode/
.idea/
*.swp
.DS_Store
Thumbs.db

# Ignore Docker-related files (to avoid copying unnecessary configs)
Dockerfile
.dockerignore
docker-compose.yml
docker-compose.override.yml

# Ignore build-specific cache files
*.lock

```

#### Step 4: Create the `nginx.conf` file

To serve your React.js application efficiently inside the container, you’ll configure NGINX with a custom setup. This configuration is optimized for performance, browser caching, gzip compression, and support for client-side routing.

Create a file named `nginx.conf` in the root of your project directory, and add the following content:

> [!NOTE]
> To learn more about configuring NGINX, see the [official NGINX documentation](https://nginx.org/en/docs/).


```nginx
worker_processes auto;

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
    error_log  /dev/stderr warn;

    # Optimize static file serving
    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;
    keepalive_timeout  65;
    keepalive_requests 1000;

    # Gzip compression for optimized delivery
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;
    gzip_min_length 256;
    gzip_vary on;

    server {
        listen       8080;
        server_name  localhost;

        # Root directory where React.js build files are placed
        root /usr/share/nginx/html;
        index index.html;

        # Serve React.js static files with proper caching
        location / {
            try_files $uri /index.html;
        }

        # Serve static assets with long cache expiration
        location ~* \.(?:ico|css|js|gif|jpe?g|png|woff2?|eot|ttf|svg|map)$ {
            expires 1y;
            access_log off;
            add_header Cache-Control "public, immutable";
        }

        # Handle React.js client-side routing
        location /static/ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

#### Step 5: Build the React.js application image

With your custom configuration in place, you're now ready to build the Docker image for your React.js application.

The updated setup includes:

- Optimized browser caching and gzip compression  
- Secure, non-root logging to avoid permission issues  
- Support for React client-side routing by redirecting unmatched routes to `index.html`  

After completing the previous steps, your project directory should now contain the following files:

```text
├── docker-reactjs-sample/
│ ├── Dockerfile
│ ├── .dockerignore
│ ├── compose.yaml
│ └── nginx.conf
```

Now that your Dockerfile is configured, you can build the Docker image for your React.js application.

> [!NOTE]
> The `docker build` command packages your application into an image using the instructions in the Dockerfile. It includes all necessary files from the current directory (called the [build context](/build/concepts/context/#what-is-a-build-context)).

Run the following command from the root of your project:

```console
$ docker build --tag docker-reactjs-sample .
```

What this command does:
- Uses the Dockerfile in the current directory (.)
- Packages the application and its dependencies into a Docker image
- Tags the image as docker-reactjs-sample so you can reference it later


#### Step 6: View local images

After building your Docker image, you can check which images are available on your local machine using either the Docker CLI or [Docker Desktop](/manuals/desktop/use-desktop/images.md). Since you're already working in the terminal, let's use the Docker CLI.

To list all locally available Docker images, run the following command:

```console
$ docker images
```

Example Output:

```shell
REPOSITORY                TAG               IMAGE ID       CREATED         SIZE
docker-reactjs-sample     latest            f39b47a97156   14 seconds ago   75.8MB
```

This output provides key details about your images:

- **Repository** – The name assigned to the image.
- **Tag** – A version label that helps identify different builds (e.g., latest).
- **Image ID** – A unique identifier for the image.
- **Created** – The timestamp indicating when the image was built.
- **Size** – The total disk space used by the image.

If the build was successful, you should see `docker-reactjs-sample` image listed. 

---

### Run the containerized application

In the previous step, you created a Dockerfile for your React.js application and built a Docker image using the docker build command. Now it’s time to run that image in a container and verify that your application works as expected.


Inside the `docker-reactjs-sample` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple React.js web application.

Press `ctrl+c` in the terminal to stop your application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `docker-reactjs-sample` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple web application preview.


To confirm that the container is running, use `docker ps` command:

```console
$ docker ps
```

This will list all active containers along with their ports, names, and status. Look for a container exposing port 8080.

Example Output:

```shell
CONTAINER ID   IMAGE                          COMMAND                  CREATED             STATUS             PORTS                    NAMES
88bced6ade95   docker-reactjs-sample-server   "nginx -c /etc/nginx…"   About a minute ago  Up About a minute  0.0.0.0:8080->8080/tcp   docker-reactjs-sample-server-1
```


To stop the application, run:

```console
$ docker compose down
```


> [!NOTE]
> For more information about Compose commands, see the [Compose CLI
> reference](/reference/cli/docker/compose/).

---

### Summary

In this guide, you learned how to containerize, build, and run a React.js application using Docker. By following best practices, you created a secure, optimized, and production-ready setup.

What you accomplished:
- Created a multi-stage `Dockerfile` that compiles the React.js application and serves the static files using Nginx.
- Created a `.dockerignore` file to exclude unnecessary files and keep the image clean and efficient.
- Built your Docker image using `docker build`.
- Ran the container using `docker compose up`, both in the foreground and in detached mode.
- Verified that the app was running by visiting [http://localhost:8080](http://localhost:8080).
- Learned how to stop the containerized application using `docker compose down`.

You now have a fully containerized React.js application, running in a Docker container, and ready for deployment across any environment with confidence and consistency.

---

### Related resources

Explore official references and best practices to sharpen your Docker workflow:

- [Multi-stage builds](/build/building/multi-stage/) – Learn how to separate build and runtime stages.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.  
- [Build context in Docker](/build/concepts/context/) – Learn how context affects image builds.  
- [`docker build` CLI reference](/reference/cli/docker/image/build/) – Build Docker images from a Dockerfile.
- [`docker images` CLI reference](/reference/cli/docker/image/ls/) – Manage and inspect local Docker images.
- [`docker compose up` CLI reference](/reference/cli/docker/compose/up/) – Start and run multi-container applications.
- [`docker compose down` CLI reference](/reference/cli/docker/compose/down/) – Stop and remove containers, networks, and volumes.

---

### Next steps

With your React.js application now containerized, you're ready to move on to the next step.

In the next section, you'll learn how to develop your application using Docker containers, enabling a consistent, isolated, and reproducible development environment across any machine.

## Use containers for React.js development

### Prerequisites

Complete [Containerize React.js application](./).

---

### Overview

In this section, you'll learn how to set up both production and development environments for your containerized React.js application using Docker Compose. This setup allows you to serve a static production build via Nginx and to develop efficiently inside containers using a live-reloading dev server with Compose Watch.

You’ll learn how to:
- Configure separate containers for production and development
- Enable automatic file syncing using Compose Watch in development
- Debug and live-preview your changes in real-time without manual rebuilds

---

### Automatically update services (development mode)

Use Compose Watch to automatically sync source file changes into your containerized development environment. This provides a seamless, efficient development experience without needing to restart or rebuild containers manually.

### Step 1: Create a development Dockerfile

Create a file named `Dockerfile.dev` in your project root with the following content:

```dockerfile
# =========================================
# Stage 1: Develop the React.js Application
# =========================================
ARG NODE_VERSION=24.12.0-alpine

# Use a lightweight Node.js image for development
FROM node:${NODE_VERSION} AS dev

# Set the working directory inside the container
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json package-lock.json* ./

# Install project dependencies
RUN --mount=type=cache,target=/root/.npm npm install

# Copy the rest of the application source code into the container
COPY . .

# Expose the port used by the Vite development server
EXPOSE 5173

# Use a default command, can be overridden in Docker compose.yml file
CMD ["npm", "run", "dev"]
```

This file sets up a lightweight development environment for your React app using the dev server.


#### Step 2: Update your `compose.yaml` file

Open your `compose.yaml` file and define two services: one for production (`react-prod`) and one for development (`react-dev`).

Here’s an example configuration for a React.js application:

```yaml
services:
  react-prod:
    build:
      context: .
      dockerfile: Dockerfile
    image: docker-reactjs-sample
    ports:
      - "8080:8080"

  react-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    develop:
      watch:
        - action: sync
          path: .
          target: /app

```
- The `react-prod` service builds and serves your static production app using Nginx.
- The `react-dev` service runs your React development server with live reload and hot module replacement.
- `watch` triggers file sync with Compose Watch.

> [!NOTE]
> For more details, see the official guide: [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

#### Step 3: Update vite.config.ts to ensure it works properly inside Docker

To make Vite’s development server work reliably inside Docker, you need to update your vite.config.ts with the correct settings.

Open the `vite.config.ts` file in your project root and update it as follows:

```ts
/// <reference types="vitest" />

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: "/",
  plugins: [react()],
  server: {
    host: true,
    port: 5173,
    strictPort: true,
  },
});
```

> [!NOTE]
> The `server` options in `vite.config.ts` are essential for running Vite inside Docker:
> - `host: true` allows the dev server to be accessible from outside the container.
> - `port: 5173` sets a consistent development port (must match the one exposed in Docker).
> - `strictPort: true` ensures Vite fails clearly if the port is unavailable, rather than switching silently.
> 
> For full details, refer to the [Vite server configuration docs](https://vitejs.dev/config/server-options.html).


After completing the previous steps, your project directory should now contain the following files:

```text
├── docker-reactjs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ └── nginx.conf
```

#### Step 4: Start Compose Watch

Run the following command from your project root to start your container in watch mode:

```console
$ docker compose watch react-dev
```

#### Step 5: Test Compose Watch with React

To verify that Compose Watch is working correctly:

1. Open the `src/App.tsx` file in your text editor.

2. Locate the following line:

    ```html
    <h1>Vite + React</h1>
    ```

3. Change it to:

    ```html
    <h1>Hello from Docker Compose Watch</h1>
    ```

4. Save the file.

5. Open your browser at [http://localhost:5173](http://localhost:5173).

You should see the updated text appear instantly, without needing to rebuild the container manually. This confirms that file watching and automatic synchronization are working as expected.

---

### Summary

In this section, you set up a complete development and production workflow for your React.js application using Docker and Docker Compose.

Here's what you achieved:
- Created a `Dockerfile.dev` to streamline local development with hot reloading  
- Defined separate `react-dev` and `react-prod` services in your `compose.yaml` file  
- Enabled real-time file syncing using Compose Watch for a smoother development experience  
- Verified that live updates work seamlessly by modifying and previewing a component

With this setup, you're now equipped to build, run, and iterate on your React.js app entirely within containers—efficiently and consistently across environments.

---

### Related resources

Deepen your knowledge and improve your containerized development workflow with these guides:

- [Using Compose Watch](/manuals/compose/how-tos/file-watch.md) – Automatically sync source changes during development  
- [Multi-stage builds](/manuals/build/building/multi-stage.md) – Create efficient, production-ready Docker images  
- [Dockerfile best practices](/build/building/best-practices/) – Write clean, secure, and optimized Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.
- [Docker volumes](/storage/volumes/) – Persist and manage data between container runs  

### Next steps

In the next section, you'll learn how to run unit tests for your React.js application inside Docker containers. This ensures consistent testing across all environments and removes dependencies on local machine setup.

## Run React.js tests in a container

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize React.js application](./).

### Overview

Testing is a critical part of the development process. In this section, you'll learn how to:

- Run unit tests using Vitest inside a Docker container.
- Use Docker Compose to run tests in an isolated, reproducible environment.

You’ll use [Vitest](https://vitest.dev) — a blazing fast test runner designed for Vite — along with [Testing Library](https://testing-library.com/) for assertions.

---

### Run tests during development

`docker-reactjs-sample` application includes a sample test file at location:

```console
$ src/App.test.tsx
```

This file uses Vitest and React Testing Library to verify the behavior of `App` component.

#### Step 1: Install Vitest and React Testing Library

If you haven’t already added the necessary testing tools, install them by running:

```console
$ npm install --save-dev vitest @testing-library/react @testing-library/jest-dom jsdom
```

Then, update the scripts section of your `package.json` file to include the following:

```json
"scripts": {
  "test": "vitest run"
}
```

---

#### Step 2: Configure Vitest

Update `vitest.config.ts` file in your project root with the following configuration:

```ts {hl_lines="14-18",linenos=true}
/// <reference types="vitest" />

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: "/",
  plugins: [react()],
  server: {
    host: true,
    port: 5173,
    strictPort: true,
  },
  test: {
    environment: "jsdom",
    setupFiles: "./src/setupTests.ts",
    globals: true,
  },
});
```

> [!NOTE]
> The `test` options in `vitest.config.ts` are essential for reliable testing inside Docker:
> - `environment: "jsdom"` simulates a browser-like environment for rendering and DOM interactions.  
> - `setupFiles: "./src/setupTests.ts"` loads global configuration or mocks before each test file (optional but recommended).  
> - `globals: true` enables global test functions like `describe`, `it`, and `expect` without importing them.
>
> For more details, see the official [Vitest configuration docs](https://vitest.dev/config/).

#### Step 3: Update compose.yaml

Add a new service named `react-test` to your `compose.yaml` file. This service allows you to run your test suite in an isolated containerized environment.

```yaml {hl_lines="22-26",linenos=true}
services:
  react-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    develop:
      watch:
        - action: sync
          path: .
          target: /app

  react-prod:
    build:
      context: .
      dockerfile: Dockerfile
    image: docker-reactjs-sample
    ports:
      - "8080:8080"

  react-test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    command: ["npm", "run", "test"]

```

The react-test service reuses the same `Dockerfile.dev` used for [development](develop.md) and overrides the default command to run tests with `npm run test`. This setup ensures a consistent test environment that matches your local development configuration.


After completing the previous steps, your project directory should contain the following files:

```text
├── docker-reactjs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ └── nginx.conf
```

#### Step 4: Run the tests

To execute your test suite inside the container, run the following command from your project root:

```console
$ docker compose run --rm react-test
```

This command will:
- Start the `react-test` service defined in your `compose.yaml` file.
- Execute the `npm run test` script using the same environment as development.
- Automatically remove the container after the tests complete [`docker compose run --rm`](/reference/cli/docker/compose/run/) command.

> [!NOTE]
> For more information about Compose commands, see the [Compose CLI
> reference](/reference/cli/docker/compose/).

---

### Summary

In this section, you learned how to run unit tests for your React.js application inside a Docker container using Vitest and Docker Compose.

What you accomplished:
- Installed and configured Vitest and React Testing Library for testing React components.
- Created a `react-test` service in `compose.yaml` to isolate test execution.
- Reused the development `Dockerfile.dev` to ensure consistency between dev and test environments.
- Ran tests inside the container using `docker compose run --rm react-test`.
- Ensured reliable, repeatable testing across environments without relying on local machine setup.

---

### Related resources

Explore official references and best practices to sharpen your Docker testing workflow:

- [Dockerfile reference](/reference/dockerfile/) – Understand all Dockerfile instructions and syntax.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.  
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) – Run one-off commands in a service container.
---

### Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions to automatically build and test your React.js application in a containerized environment. This ensures your code is validated on every push or pull request, maintaining consistency and reliability across your development workflow.

## Test your React.js deployment

### Prerequisites

Before you begin, make sure you’ve completed the following:
- Complete all the previous sections of this guide, starting with [Containerize React.js application](./).
- [Enable Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

> **New to Kubernetes?**  
> Visit the [Kubernetes basics tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/) to get familiar with how clusters, pods, deployments, and services work.

---

### Overview

This section guides you through deploying your containerized React.js application locally using [Docker Desktop’s built-in Kubernetes](/desktop/kubernetes/). Running your app in a local Kubernetes cluster allows you to closely simulate a real production environment, enabling you to test, validate, and debug your workloads with confidence before promoting them to staging or production.

---

### Create a Kubernetes YAML file

Follow these steps to define your deployment configuration:

1. In the root of your project, create a new file named: reactjs-sample-kubernetes.yaml

2. Open the file in your IDE or preferred text editor.

3. Add the following configuration, and be sure to replace `{DOCKER_USERNAME}` and `{DOCKERHUB_PROJECT_NAME}` with your actual Docker Hub username and repository name from the previous [Automate your builds with GitHub Actions](./).


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: reactjs-sample
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: reactjs-sample
  template:
    metadata:
      labels:
        app: reactjs-sample
    spec:
      containers:
        - name: reactjs-container
          image: {DOCKER_USERNAME}/{DOCKERHUB_PROJECT_NAME}:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name:  reactjs-sample-service
  namespace: default
spec:
  type: NodePort
  selector:
    app:  reactjs-sample
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 30001
```

This manifest defines two key Kubernetes resources, separated by `---`:

- Deployment
  Deploys a single replica of your React.js application inside a pod. The pod uses the Docker image built and pushed by your GitHub Actions CI/CD workflow  
  (refer to [Automate your builds with GitHub Actions](./)).  
  The container listens on port `8080`, which is typically used by [Nginx](https://nginx.org/en/docs/) to serve your production React app.

- Service (NodePort) 
  Exposes the deployed pod to your local machine.  
  It forwards traffic from port `30001` on your host to port `8080` inside the container.  
  This lets you access the application in your browser at [http://localhost:30001](http://localhost:30001).

> [!NOTE]
> To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

---

### Deploy and check your application

Follow these steps to deploy your containerized React.js app into a local Kubernetes cluster and verify that it’s running correctly.

#### Step 1. Apply the Kubernetes configuration

In your terminal, navigate to the directory where your `reactjs-sample-kubernetes.yaml` file is located, then deploy the resources using:

```console
  $ kubectl apply -f reactjs-sample-kubernetes.yaml
```

If everything is configured properly, you’ll see confirmation that both the Deployment and the Service were created:

```shell
  deployment.apps/reactjs-sample created
  service/reactjs-sample-service created
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
  reactjs-sample       1/1     1            1           14s
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
reactjs-sample-service   NodePort    10.100.244.65    <none>        8080:30001/TCP   1m
```

This output confirms that your app is available via NodePort on port 30001.

#### Step 4. Access your app in the browser

Open your browser and navigate to [http://localhost:30001](http://localhost:30001).

You should see your production-ready React.js Sample application running — served by your local Kubernetes cluster.

#### Step 5. Clean up Kubernetes resources

Once you're done testing, you can delete the deployment and service using:

```console
  $ kubectl delete -f reactjs-sample-kubernetes.yaml
```

Expected output:

```shell
  deployment.apps "reactjs-sample" deleted
  service "reactjs-sample-service" deleted
```

This ensures your cluster stays clean and ready for the next deployment.
   
---

### Summary

In this section, you learned how to deploy your React.js application to a local Kubernetes cluster using Docker Desktop. This setup allows you to test and debug your containerized app in a production-like environment before deploying it to the cloud.

What you accomplished:

- Created a Kubernetes Deployment and NodePort Service for your React.js app  
- Used `kubectl apply` to deploy the application locally  
- Verified the app was running and accessible at `http://localhost:30001`  
- Cleaned up your Kubernetes resources after testing

---

### Related resources

Explore official references and best practices to sharpen your Kubernetes deployment workflow:

- [Kubernetes documentation](https://kubernetes.io/docs/home/) – Learn about core concepts, workloads, services, and more.  
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md) – Use Docker Desktop’s built-in Kubernetes support for local testing and development.
- [`kubectl` CLI reference](https://kubernetes.io/docs/reference/kubectl/) – Manage Kubernetes clusters from the command line.  
- [Kubernetes Deployment resource](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) – Understand how to manage and scale applications using Deployments.  
- [Kubernetes Service resource](https://kubernetes.io/docs/concepts/services-networking/service/) – Learn how to expose your application to internal and external traffic.

## Automate your builds with GitHub Actions

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize React.js application](./).

You must also have:
- A [GitHub](https://github.com/signup) account.
- A verified [Docker Hub](https://hub.docker.com/signup) account.

---

### Overview

In this section, you'll set up a CI/CD pipeline using [GitHub Actions](https://docs.github.com/en/actions) to automatically:

- Build your React.js application inside a Docker container.
- Run tests in a consistent environment.
- Push the production-ready image to [Docker Hub](https://hub.docker.com).

---

### Connect your GitHub repository to Docker Hub

To enable GitHub Actions to build and push Docker images, you’ll securely store your Docker Hub credentials in your new GitHub repository.

#### Step 1: Connect your GitHub repository to Docker Hub

1. Create a Personal Access Token (PAT) from [Docker Hub](https://hub.docker.com)
   1. Go to your **Docker Hub account → Account Settings → Security**.
   2. Generate a new Access Token with **Read/Write** permissions.
   3. Name it something like `docker-reactjs-sample`.
   4. Copy and save the token — you’ll need it in Step 4.

2. Create a repository in [Docker Hub](https://hub.docker.com/repositories/)
   1. Go to your **Docker Hub account → Create a repository**.
   2. For the Repository Name, use something descriptive — for example: `reactjs-sample`.
   3. Once created, copy and save the repository name — you’ll need it in Step 4.

3. Create a new [GitHub repository](https://github.com/new) for your React.js project

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

   These secrets let GitHub Actions to authenticate securely with Docker Hub during automated workflows.

5. Connect Your Local Project to GitHub

   Link your local project `docker-reactjs-sample` to the GitHub repository you just created by running the following command from your project root:

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

Once completed, your code will be available on GitHub, and any GitHub Actions workflow you’ve configured will run automatically.

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
name: CI/CD – React.js Application with Docker

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
        uses: actions/checkout@{{% param "checkout_action_version" %}}
        with:
          fetch-depth: 0 # Fetches full history for better caching/context

      # 2. Set up Docker Buildx
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      # 3. Cache Docker layers
      - name: Cache Docker layers
        uses: actions/cache@{{% param "cache_action_version" %}}
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: ${{ runner.os }}-buildx-

      # 4. Cache npm dependencies
      - name: Cache npm dependencies
        uses: actions/cache@{{% param "cache_action_version" %}}
        with:
          path: ~/.npm
          key: ${{ runner.os }}-npm-${{ hashFiles('**/package-lock.json') }}
          restore-keys: ${{ runner.os }}-npm-

      # 5. Extract metadata
      - name: Extract metadata
        id: meta
        run: |
          echo "REPO_NAME=${GITHUB_REPOSITORY##*/}" >> "$GITHUB_OUTPUT"
          echo "SHORT_SHA=${GITHUB_SHA::7}" >> "$GITHUB_OUTPUT"

      # 6. Build dev Docker image
      - name: Build Docker image for tests
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          context: .
          file: Dockerfile.dev
          tags: ${{ steps.meta.outputs.REPO_NAME }}-dev:latest
          load: true # Load to local Docker daemon for testing
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max

      # 7. Run Vitest tests
      - name: Run Vitest tests and generate report
        run: |
          docker run --rm \
            --workdir /app \
            --entrypoint "" \
            ${{ steps.meta.outputs.REPO_NAME }}-dev:latest \
            sh -c "npm ci && npx vitest run --reporter=verbose"
        env:
          CI: true
          NODE_ENV: test
        timeout-minutes: 10

      # 8. Login to Docker Hub
      - name: Log in to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      # 9. Build and push prod image
      - name: Build and push production image
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          context: .
          file: Dockerfile
          push: true
          platforms: linux/amd64,linux/arm64
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:latest
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:${{ steps.meta.outputs.SHORT_SHA }}
          cache-from: type=local,src=/tmp/.buildx-cache
```

This workflow performs the following tasks for your React.js application:
- Triggers on every `push` or `pull request` targeting the `main` branch.
- Builds a development Docker image using `Dockerfile.dev`, optimized for testing.
- Executes unit tests using Vitest inside a clean, containerized environment to ensure consistency.
- Halts the workflow immediately if any test fails — enforcing code quality.
- Caches both Docker build layers and npm dependencies for faster CI runs.
- Authenticates securely with Docker Hub using GitHub repository secrets.
- Builds a production-ready image using the `prod` stage in `Dockerfile`.
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
>  - Under Branch protection rules, click **Add rule**.
>  - Specify `main` as the branch name.
>  - Enable options like:
>     - *Require a pull request before merging*.
>     - *Require status checks to pass before merging*.
>
>  This ensures that only tested and reviewed code is merged into `main` branch.
---

### Summary

In this section, you set up a complete CI/CD pipeline for your containerized React.js application using GitHub Actions.

Here's what you accomplished:

- Created a new GitHub repository specifically for your project.
- Generated a secure Docker Hub access token and added it to GitHub as a secret.
- Defined a GitHub Actions workflow to:
   - Build your application inside a Docker container.
   - Run tests in a consistent, containerized environment.
   - Push a production-ready image to Docker Hub if tests pass.
- Triggered and verified the workflow execution through GitHub Actions.
- Confirmed that your image was successfully published to Docker Hub.

With this setup, your React.js application is now ready for automated testing and deployment across environments — increasing confidence, consistency, and team productivity.

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

Next, learn how you can locally test and debug your React.js workloads on Kubernetes before deploying. This helps you ensure your application behaves as expected in a production-like environment, reducing surprises during deployment.
