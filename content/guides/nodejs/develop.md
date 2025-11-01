---
title: Use containers for Node.js development
linkTitle: Develop your app
weight: 30
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

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Adding a local database and persisting data
- Configuring your container to run a development environment
- Debugging your containerized application

## Add a local database and persist data

The application uses PostgreSQL for data persistence, providing a production-ready database solution. You'll need to add a database service to your Docker Compose configuration.

### Add database service to Docker Compose

If you haven't already created a `compose.yml` file in the previous section, or if you need to add the database service, update your `compose.yml` file to include the PostgreSQL database service:

```yaml
services:
  # ... existing app services ...

  # ========================================
  # PostgreSQL Database Service
  # ========================================
  db:
    image: postgres:16-alpine
    container_name: todoapp-db
    environment:
      POSTGRES_DB: '${POSTGRES_DB:-todoapp}'
      POSTGRES_USER: '${POSTGRES_USER:-todoapp}'
      POSTGRES_PASSWORD: '${POSTGRES_PASSWORD:-todoapp_password}'
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - '${DB_PORT:-5432}:5432'
    restart: unless-stopped
    healthcheck:
      test: ['CMD-SHELL', 'pg_isready -U ${POSTGRES_USER:-todoapp} -d ${POSTGRES_DB:-todoapp}']
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 5s
    networks:
      - todoapp-network

# ========================================
# Volume Configuration
# ========================================
volumes:
  postgres_data:
    name: todoapp-postgres-data
    driver: local

# ========================================
# Network Configuration
# ========================================
networks:
  todoapp-network:
    name: todoapp-network
    driver: bridge
```

### Update your application service

Make sure your application service in `compose.yml` is configured to connect to the database:

```yaml {hl_lines="18-20,42-44",collapse=true,title=compose.yml}
services:
  app-dev:
    build:
      context: .
      dockerfile: Dockerfile
      target: development
    container_name: todoapp-dev
    ports:
      - '${APP_PORT:-3000}:3000' # API server
      - '${VITE_PORT:-5173}:5173' # Vite dev server
      - '${DEBUG_PORT:-9229}:9229' # Node.js debugger
    environment:
      NODE_ENV: development
      DOCKER_ENV: 'true'
      POSTGRES_HOST: db
      POSTGRES_PORT: 5432
      POSTGRES_DB: todoapp
      POSTGRES_USER: todoapp
      POSTGRES_PASSWORD: '${POSTGRES_PASSWORD:-todoapp_password}'
      ALLOWED_ORIGINS: '${ALLOWED_ORIGINS:-http://localhost:3000,http://localhost:5173}'
    volumes:
      - ./src:/app/src:ro
      - ./package.json:/app/package.json
      - ./vite.config.ts:/app/vite.config.ts:ro
      - ./tailwind.config.js:/app/tailwind.config.js:ro
      - ./postcss.config.js:/app/postcss.config.js:ro
    depends_on:
      db:
        condition: service_healthy
    develop:
      watch:
        - action: sync
          path: ./src
          target: /app/src
          ignore:
            - '**/*.test.*'
            - '**/__tests__/**'
        - action: rebuild
          path: ./package.json
        - action: sync
          path: ./vite.config.ts
          target: /app/vite.config.ts
        - action: sync
          path: ./tailwind.config.js
          target: /app/tailwind.config.js
        - action: sync
          path: ./postcss.config.js
          target: /app/postcss.config.js
    restart: unless-stopped
    networks:
      - todoapp-network

  db:
    image: postgres:16-alpine
    container_name: todoapp-db
    environment:
      POSTGRES_DB: '${POSTGRES_DB:-todoapp}'
      POSTGRES_USER: '${POSTGRES_USER:-todoapp}'
      POSTGRES_PASSWORD: '${POSTGRES_PASSWORD:-todoapp_password}'
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - '${DB_PORT:-5432}:5432'
    restart: unless-stopped
    healthcheck:
      test: ['CMD-SHELL', 'pg_isready -U ${POSTGRES_USER:-todoapp} -d ${POSTGRES_DB:-todoapp}']
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 5s
    networks:
      - todoapp-network

volumes:
  postgres_data:
    name: todoapp-postgres-data
    driver: local

networks:
  todoapp-network:
    name: todoapp-network
    driver: bridge
```

1. The PostgreSQL database configuration is handled automatically by the application. The database is created and initialized when the application starts, with data persisted using the `postgres_data` volume.

2. Configure your environment by copying the example file:

   ```console
   $ cp .env.example .env
   ```

   Update the `.env` file with your preferred settings:

   ```env
   # Application Configuration
   NODE_ENV=development
   APP_PORT=3000
   VITE_PORT=5173
   DEBUG_PORT=9230

   # Database Configuration
   POSTGRES_HOST=db
   POSTGRES_PORT=5432
   POSTGRES_DB=todoapp
   POSTGRES_USER=todoapp
   POSTGRES_PASSWORD=todoapp_password

   # Security Configuration
   ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
   ```

3. Run the following command to start your application in development mode:

   ```console
   $ docker compose up app-dev --build
   ```

4. Open a browser and verify that the application is running at [http://localhost:5173](http://localhost:5173) for the frontend or [http://localhost:3000](http://localhost:3000) for the API. The React frontend is served by Vite dev server on port 5173, with API calls proxied to the Express server on port 3000.

5. Add some items to the todo list to test data persistence.

6. After adding some items to the todo list, press `CTRL + C` in the terminal to stop your application.

7. Run the application again:

   ```console
   $ docker compose up app-dev
   ```

8. Refresh [http://localhost:5173](http://localhost:5173) in your browser and verify that the todo items persisted, even after the containers were removed and ran again.

## Configure and run a development container

You can use a bind mount to mount your source code into the container. The container can then see the changes you make to the code immediately, as soon as you save a file. This means that you can run processes, like nodemon, in the container that watch for filesystem changes and respond to them. To learn more about bind mounts, see [Storage overview](/manuals/engine/storage/_index.md).

In addition to adding a bind mount, you can configure your Dockerfile and `compose.yaml` file to install development dependencies and run development tools.

### Update your Dockerfile for development

Your Dockerfile should be configured as a multi-stage build with separate stages for development, production, and testing. If you followed the previous section, your Dockerfile already includes a development stage that has all development dependencies and runs the application with hot reload enabled.

Here's the development stage from your multi-stage Dockerfile:

```dockerfile {hl_lines="5-26",collapse=true,title=Dockerfile}
# ========================================
# Development Stage
# ========================================
FROM build-deps AS development

# Set environment
ENV NODE_ENV=development \
    NPM_CONFIG_LOGLEVEL=warn

# Copy source files
COPY . .

# Ensure all directories have proper permissions
RUN mkdir -p /app/node_modules/.vite && \
    chown -R nodejs:nodejs /app && \
    chmod -R 755 /app

# Switch to non-root user
USER nodejs

# Expose ports
EXPOSE 3000 5173 9229

# Start development server
CMD ["npm", "run", "dev:docker"]
```

The development stage:

- Installs all dependencies including dev dependencies
- Exposes ports for the API server (3000), Vite dev server (5173), and Node.js debugger (9229)
- Runs `npm run dev` which starts both the Express server and Vite dev server concurrently
- Includes health checks for monitoring container status

Next, you'll need to update your Compose file to use the new stage.

### Update your Compose file for development

Now you need to configure your `compose.yml` file to run the development stage with comprehensive bind mounts for hot reloading. Update your development service configuration:

```yaml {hl_lines=[5,8-10,20-27],collapse=true,title=compose.yml}
services:
  app-dev:
    build:
      context: .
      dockerfile: Dockerfile
      target: development
    container_name: todoapp-dev
    ports:
      - '${APP_PORT:-3000}:3000' # API server
      - '${VITE_PORT:-5173}:5173' # Vite dev server
      - '${DEBUG_PORT:-9229}:9229' # Node.js debugger
    environment:
      NODE_ENV: development
      DOCKER_ENV: 'true'
      POSTGRES_HOST: db
      POSTGRES_PORT: 5432
      POSTGRES_DB: todoapp
      POSTGRES_USER: todoapp
      POSTGRES_PASSWORD: '${POSTGRES_PASSWORD:-todoapp_password}'
      ALLOWED_ORIGINS: '${ALLOWED_ORIGINS:-http://localhost:3000,http://localhost:5173}'
    volumes:
      - ./src:/app/src:ro
      - ./package.json:/app/package.json
      - ./vite.config.ts:/app/vite.config.ts:ro
      - ./tailwind.config.js:/app/tailwind.config.js:ro
      - ./postcss.config.js:/app/postcss.config.js:ro
    depends_on:
      db:
        condition: service_healthy
    develop:
      watch:
        - action: sync
          path: ./src
          target: /app/src
          ignore:
            - '**/*.test.*'
            - '**/__tests__/**'
        - action: rebuild
          path: ./package.json
        - action: sync
          path: ./vite.config.ts
          target: /app/vite.config.ts
        - action: sync
          path: ./tailwind.config.js
          target: /app/tailwind.config.js
        - action: sync
          path: ./postcss.config.js
          target: /app/postcss.config.js
    restart: unless-stopped
    networks:
      - todoapp-network
```

Key features of the development configuration:

- **Multi-port exposure**: API server (3000), Vite dev server (5173), and debugger (9229)
- **Comprehensive bind mounts**: Source code, configuration files, and package files for hot reloading
- **Environment variables**: Configurable through `.env` file or defaults
- **PostgreSQL database**: Production-ready database with persistent storage
- **Docker Compose watch**: Automatic file synchronization and container rebuilds
- **Health checks**: Database health monitoring with automatic dependency management

### Run your development container and debug your application

Run the following command to run your application with the development configuration:

```console
$ docker compose up app-dev --build
```

Or with file watching for automatic updates:

```console
$ docker compose up app-dev --watch
```

For local development without Docker:

```console
$ npm run dev:with-db
```

Or start services separately:

```console
$ npm run db:start    # Start PostgreSQL container
$ npm run dev         # Start both server and client
```

### Using Task Runner (Alternative)

The project includes a comprehensive Taskfile.yml for advanced workflows:

```console
# Development
$ task dev              # Start development environment
$ task dev:build        # Build development image
$ task dev:run          # Run development container

# Production
$ task build            # Build production image
$ task run              # Run production container
$ task build-run        # Build and run in one step

# Testing
$ task test             # Run all tests
$ task test:unit        # Run unit tests with coverage
$ task test:lint        # Run linting

# Kubernetes
$ task k8s:deploy       # Deploy to Kubernetes
$ task k8s:status       # Check deployment status
$ task k8s:logs         # View pod logs

# Utilities
$ task clean            # Clean up containers and images
$ task health           # Check application health
$ task logs             # View container logs
```

The application will start with both the Express API server and Vite development server:

- **API Server**: [http://localhost:3000](http://localhost:3000) - Express.js backend with REST API
- **Frontend**: [http://localhost:5173](http://localhost:5173) - Vite dev server with hot module replacement
- **Health Check**: [http://localhost:3000/health](http://localhost:3000/health) - Application health status

Any changes to the application's source files on your local machine will now be immediately reflected in the running container thanks to the bind mounts.

Try making a change to test hot reloading:

1. Open `src/client/components/TodoApp.tsx` in an IDE or text editor
2. Update the main heading text:

```diff
- <h1 className="text-3xl font-bold text-gray-900 mb-8">
-   Modern Todo App
- </h1>
+ <h1 className="text-3xl font-bold text-gray-900 mb-8">
+   My Todo App
+ </h1>
```

1. Save the file and the Vite dev server will automatically reload the page with your changes

**Debugging Support:**

You can connect a debugger to your application on port 9229. The Node.js inspector is enabled with `--inspect=0.0.0.0:9230` in the development script (`dev:server`).

### VS Code Debugger Setup

1. **Create a launch configuration** in `.vscode/launch.json`:

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

1. **Start your development container**:

```console
docker compose up app-dev --build
```

1. **Attach the debugger**:
   - Open VS Code
   - Go to the Debug panel (Ctrl/Cmd + Shift + D)
   - Select "Attach to Docker Container" from the drop-down
   - Select the green play button or press F5

### Chrome DevTools (Alternative)

You can also use Chrome DevTools for debugging:

1. **Start your container** (if not already running):

```console
docker compose up app-dev --build
```

1. **Open Chrome** and navigate to:

```text
chrome://inspect
```

1. **Select "Configure"** and add:

```text
localhost:9229
```

1. **Select "inspect"** under your Node.js target when it appears

### Debugging Configuration Details

The debugger configuration:

- **Container port**: 9230 (internal debugger port)
- **Host port**: 9229 (mapped external port)
- **Script**: `tsx watch --inspect=0.0.0.0:9230 src/server/index.ts`

The debugger listens on all interfaces (`0.0.0.0`) inside the container on port 9230 and is accessible on port 9229 from your host machine.

### Troubleshooting Debugger Connection

If the debugger doesn't connect:

1. **Check if the container is running**:

```console
docker ps
```

1. **Check if the port is exposed**:

```console
docker port todoapp-dev
```

1. **Check container logs**:

```console
docker compose logs app-dev
```

You should see a message like:

```text
Debugger listening on ws://0.0.0.0:9230/...
```

Now you can set breakpoints in your TypeScript source files and debug your containerized Node.js application!

For more details about Node.js debugging, see the [Node.js documentation](https://nodejs.org/en/docs/guides/debugging-getting-started).

## Summary

In this section, you took a look at setting up your Compose file to add a mock
database and persist data. You also learned how to create a multi-stage
Dockerfile and set up a bind mount for development.

Related information:

- [Volumes top-level element](/reference/compose-file/volumes/)
- [Services top-level element](/reference/compose-file/services/)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)

## Next steps

In the next section, you'll learn how to run unit tests using Docker.
