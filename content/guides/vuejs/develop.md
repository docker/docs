---
title: Use containers for Vue.js development
linkTitle: Develop your app
weight: 30
keywords: vuejs, development, node
description: Learn how to develop your Vue.js application locally using containers.

---

## Prerequisites

Complete [Containerize Vue.js application](containerize.md).

---

## Overview

In this section, you'll set up both production and development environments for your Vue.js application using Docker Compose. This approach streamlines your workflow—delivering a lightweight, static site via Nginx in production, and providing a fast, live-reloading dev server with Compose Watch for efficient local development.

You’ll learn how to:
- Configure isolated environments: Set up separate containers optimized for production and development use cases.
- Live-reload in development: Use Compose Watch to automatically sync file changes, enabling real-time updates without manual intervention.
- Preview and debug with ease: Develop inside containers with a seamless preview and debug experience—no rebuilds required after every change.

---

## Automatically update services (Development Mode)

Leverage Compose Watch to enable real-time file synchronization between your local machine and the containerized Vue.js development environment. This powerful feature eliminates the need to manually rebuild or restart containers, providing a fast, seamless, and efficient development workflow.

With Compose Watch, your code updates are instantly reflected inside the container—perfect for rapid testing, debugging, and live previewing changes.

## Step 1: Create a development Dockerfile

Create a file named `Dockerfile.dev` in your project root with the following content:

```dockerfile
# =========================================
# Stage 1: Develop the Vue.js Application
# =========================================
ARG NODE_VERSION=24.12.0-alpine

# Use a lightweight Node.js image for development
FROM node:${NODE_VERSION} AS dev

# Set environment variable to indicate development mode
ENV NODE_ENV=development

# Set the working directory inside the container
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json package-lock.json ./

# Install project dependencies
RUN --mount=type=cache,target=/root/.npm npm install

# Copy the rest of the application source code into the container
COPY . .

# Change ownership of the application directory to the node user
RUN chown -R node:node /app

# Switch to the node user
USER node

# Expose the port used by the Vite development server
EXPOSE 5173

# Use a default command, can be overridden in Docker compose.yml file
CMD [ "npm", "run", "dev", "--", "--host" ]

```

This file sets up a lightweight development environment for your Vue.js application using the dev server.

### Step 2: Update your `compose.yaml` file

Open your `compose.yaml` file and define two services: one for production (`vuejs-prod`) and one for development (`vuejs-dev`).

Here’s an example configuration for an Vue.js application:

```yaml
services:
  vuejs-prod:
    build:
      context: .
      dockerfile: Dockerfile
    image: docker-vuejs-sample
    ports:
      - "8080:8080"

  vuejs-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    develop:
      watch:
        - path: ./src
          target: /app/src
          action: sync
        - path: ./package.json
          target: /app/package.json
          action: restart
        - path: ./vite.config.js
          target: /app/vite.config.js
          action: restart
```
- The `vuejs-prod` service builds and serves your static production app using Nginx.
- The `vuejs-dev` service runs your Vue.js development server with live reload and hot module replacement.
- `watch` triggers file sync with Compose Watch.

> [!NOTE]
> For more details, see the official guide: [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

After completing the previous steps, your project directory should now contain the following files:

```text
├── docker-vuejs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── nginx.conf
│ └── README.Docker.md
```

### Step 4: Start Compose Watch

Run the following command from the project root to start the container in watch mode

```console
$ docker compose watch vuejs-dev
```

### Step 5: Test Compose Watch with Vue.js

To confirm that Compose Watch is functioning correctly:

1. Open the `src/App.vue` file in your text editor.

2. Locate the following line:

    ```html
    <HelloWorld msg="You did it!" />
    ```

3. Change it to:

    ```html
    <HelloWorld msg="Hello from Docker Compose Watch" />
    ```

4. Save the file.

5. Open your browser at [http://localhost:5173](http://localhost:5173).

You should see the updated text appear instantly, without needing to rebuild the container manually. This confirms that file watching and automatic synchronization are working as expected.

---

## Summary

In this section, you set up a complete development and production workflow for your Vue.js application using Docker and Docker Compose.

Here’s what you accomplished:
- Created a `Dockerfile.dev` to streamline local development with hot reloading  
- Defined separate `vuejs-dev` and `vuejs-prod` services in your `compose.yaml` file  
- Enabled real-time file syncing using Compose Watch for a smoother development experience  
- Verified that live updates work seamlessly by modifying and previewing a component

With this setup, you're now equipped to build, run, and iterate on your Vue.js app entirely within containers—efficiently and consistently across environments.

---

## Related resources

Deepen your knowledge and improve your containerized development workflow with these guides:

- [Using Compose Watch](/manuals/compose/how-tos/file-watch.md) – Automatically sync source changes during development  
- [Multi-stage builds](/manuals/build/building/multi-stage.md) – Create efficient, production-ready Docker images  
- [Dockerfile best practices](/build/building/best-practices/) – Write clean, secure, and optimized Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.
- [Docker volumes](/storage/volumes/) – Persist and manage data between container runs  

## Next steps

In the next section, you'll learn how to run unit tests for your Vue.js application inside Docker containers. This ensures consistent testing across all environments and removes dependencies on local machine setup.
