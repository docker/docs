---
title: Angular language-specific guide
linkTitle: Angular
description: Containerize and develop Angular apps using Docker
keywords: getting started, angular, docker, language, Dockerfile
summary: |
  This guide explains how to containerize Angular applications using Docker.
aliases:
  - /guides/angular/configure-github-actions/
  - /guides/angular/containerize/
  - /guides/angular/deploy/
  - /guides/angular/develop/
  - /guides/angular/run-tests/
params:
  tags: [cicd]
  time: 20 minutes
---


The Angular language-specific guide shows you how to containerize an Angular application using Docker, following best practices for creating efficient, production-ready containers.

[Angular](https://angular.dev/) is a robust and widely adopted framework for building dynamic, enterprise-grade web applications. However, managing dependencies, environments, and deployments can become complex as applications scale. Docker streamlines these challenges by offering a consistent, isolated environment for development and production.

> 
> **Acknowledgment**
>
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide. As a Docker Captain and experienced Front-end engineer, his expertise in Docker, DevOps, and modern web development has made this resource essential for the community, helping developers navigate and optimize their Docker workflows.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run an Angular application using Docker.
- Set up a local development environment for Angular inside a container.
- Run tests for your Angular application within a Docker container.
- Configure a CI/CD pipeline using GitHub Actions for your containerized app.
- Deploy the containerized Angular application to a local Kubernetes cluster for testing and debugging.

You'll start by containerizing an existing Angular application and work your way up to production-level deployments.

---

## Prerequisites

Before you begin, ensure you have a working knowledge of:

- Basic understanding of [TypeScript](https://www.typescriptlang.org/) and [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript).
- Familiarity with [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm) for managing dependencies and running scripts.
- Familiarity with [Angular](https://angular.io/) fundamentals.
- Understanding of core Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Once you've completed the Angular getting started modules, you’ll be fully prepared to containerize your own Angular application using the detailed examples and best practices outlined in this guide.

## Containerize an Angular Application

### Prerequisites

Before you begin, make sure the following tools are installed and available on your system:

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

> **New to Docker?**  
> Start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide to get familiar with key concepts like images, containers, and Dockerfiles.

---

### Overview

This guide walks you through the complete process of containerizing an Angular application with Docker. You’ll learn how to create a production-ready Docker image using best practices that improve performance, security, scalability, and deployment efficiency.

By the end of this guide, you will:

- Containerize an Angular application using Docker.
- Create and optimize a Dockerfile for production builds. 
- Use multi-stage builds to minimize image size.
- Serve the application efficiently with a custom Nginx configuration.
- Build secure and maintainable Docker images by following best practices.

---

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, navigate to the directory where you want to work, and run the following command
to clone the git repository:

```console
$ git clone https://github.com/kristiyan-velkov/docker-angular-sample
```
---

### Build the Docker image

Angular is a front-end framework that compiles into static assets, so the Dockerfile uses a multi-stage build: one stage compiles the app with Node.js, and a second minimal stage serves the static output with Nginx.

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

In the following Dockerfile, the `FROM` instruction uses `dhi.io/node:24-alpine3.22-dev` as the base image.

```dockerfile
# =========================================
# Stage 1: Build the Angular Application
# =========================================

# Use a lightweight DHI Node.js image for building
FROM dhi.io/node:24-alpine3.22-dev AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json package-lock.json* ./

# Install project dependencies using npm ci (ensures a clean, reproducible install)
RUN --mount=type=cache,target=/root/.npm npm ci

# Copy the rest of the application source code into the container
COPY . .

# Build the Angular application
RUN npm run build 

# =========================================
# Stage 2: Prepare Nginx to Serve Static Files
# =========================================

FROM dhi.io/nginx:1.28.0-alpine3.21-dev AS runner

# Copy custom Nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the static build output from the build stage to Nginx's default HTML serving directory
COPY --chown=nginx:nginx --from=builder /app/dist/*/browser /usr/share/nginx/html

# Use a non-root user for security best practices
USER nginx

# Expose port 8080 to allow HTTP traffic
# Note: The default Nginx container now listens on port 8080 instead of 80 
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
# Stage 1: Build the Angular Application
# =========================================
ARG NODE_VERSION=24.12.0-alpine
ARG NGINX_VERSION=alpine3.22

# Use a lightweight Node.js image for building (customizable via ARG)
FROM node:${NODE_VERSION} AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy package-related files first to leverage Docker's caching mechanism
COPY package.json *package-lock.json* ./

# Install project dependencies using npm ci (ensures a clean, reproducible install)
RUN --mount=type=cache,target=/root/.npm npm ci

# Copy the rest of the application source code into the container
COPY . .

# Build the Angular application
RUN npm run build 

# =========================================
# Stage 2: Prepare Nginx to Serve Static Files
# =========================================

FROM nginxinc/nginx-unprivileged:${NGINX_VERSION} AS runner

# Copy custom Nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the static build output from the build stage to Nginx's default HTML serving directory
COPY --chown=nginx:nginx --from=builder /app/dist/*/browser /usr/share/nginx/html

# Use a built-in non-root user for security best practices
USER nginx

# Expose port 8080 to allow HTTP traffic
# Note: The default Nginx container now listens on port 8080 instead of 80 
EXPOSE 8080

# Start Nginx directly with custom config
ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]
```

> [!NOTE]
> We are using nginx-unprivileged instead of the standard Nginx image to follow security best practices.
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
# ================================
# Node and build output
# ================================
node_modules
dist
out-tsc
.angular
.cache
.tmp

# ================================
# Testing & Coverage
# ================================
coverage
jest
cypress
cypress/screenshots
cypress/videos
reports
playwright-report
.vite
.vitepress

# ================================
# Environment & log files
# ================================
*.env*
!*.env.production
*.log
*.tsbuildinfo

# ================================
# IDE & OS-specific files
# ================================
.vscode
.idea
.DS_Store
Thumbs.db
*.swp

# ================================
# Version control & CI files
# ================================
.git
.gitignore

# ================================
# Docker & local orchestration
# ================================
Dockerfile
Dockerfile.*
.dockerignore
docker-compose.yml
docker-compose*.yml

# ================================
# Miscellaneous
# ================================
*.bak
*.old
*.tmp
```

#### Step 4: Create the `nginx.conf` file

To serve your Angular application efficiently inside the container, you’ll configure Nginx with a custom setup. This configuration is optimized for performance, browser caching, gzip compression, and support for client-side routing.

Create a file named `nginx.conf` in the root of your project directory, and add the following content:

> [!NOTE]
> To learn more about configuring Nginx, see the [official Nginx documentation](https://nginx.org/en/docs/).


```nginx
worker_processes auto;

pid /tmp/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    client_body_temp_path /tmp/client_temp;
    proxy_temp_path       /tmp/proxy_temp_path;
    fastcgi_temp_path     /tmp/fastcgi_temp;
    uwsgi_temp_path       /tmp/uwsgi_temp;
    scgi_temp_path        /tmp/scgi_temp;

    # Logging
    access_log off;
    error_log  /dev/stderr warn;

    # Performance
    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;
    keepalive_timeout  65;
    keepalive_requests 1000;

    # Compression
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_min_length 256;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/x-javascript
        application/json
        application/xml
        application/xml+rss
        font/ttf
        font/otf
        image/svg+xml;

    server {
        listen       8080;
        server_name  localhost;

        root /usr/share/nginx/html;
        index index.html;

        # Angular Routing
        location / {
            try_files $uri $uri/ /index.html;
        }

        # Static Assets Caching
        location ~* \.(?:ico|css|js|gif|jpe?g|png|woff2?|eot|ttf|svg|map)$ {
            expires 1y;
            access_log off;
            add_header Cache-Control "public, immutable";
        }

        # Optional: Explicit asset route
        location /assets/ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

#### Step 5: Build the Angular application image

With your custom configuration in place, you're now ready to build the Docker image for your Angular application.

The updated setup includes:

- The updated setup includes a clean, production-ready Nginx configuration tailored specifically for Angular.
- Efficient multi-stage Docker build, ensuring a small and secure final image.

After completing the previous steps, your project directory should now contain the following files:

```text
├── docker-angular-sample/
│ ├── Dockerfile
│ ├── .dockerignore
│ ├── compose.yaml
│ └── nginx.conf
```

Now that your Dockerfile is configured, you can build the Docker image for your Angular application.

> [!NOTE]
> The `docker build` command packages your application into an image using the instructions in the Dockerfile. It includes all necessary files from the current directory (called the [build context](/build/concepts/context/#what-is-a-build-context)).

Run the following command from the root of your project:

```console
$ docker build --tag docker-angular-sample .
```

What this command does:
- Uses the Dockerfile in the current directory (.)
- Packages the application and its dependencies into a Docker image
- Tags the image as docker-angular-sample so you can reference it later


#### Step 6: View local images

After building your Docker image, you can check which images are available on your local machine using either the Docker CLI or [Docker Desktop](/manuals/desktop/use-desktop/images.md). Since you're already working in the terminal, let's use the Docker CLI.

To list all locally available Docker images, run the following command:

```console
$ docker images
```

Example Output:

```shell
REPOSITORY                TAG               IMAGE ID       CREATED         SIZE
docker-angular-sample     latest            34e66bdb9d40   14 seconds ago   76.4MB
```

This output provides key details about your images:

- **Repository** – The name assigned to the image.
- **Tag** – A version label that helps identify different builds (e.g., latest).
- **Image ID** – A unique identifier for the image.
- **Created** – The timestamp indicating when the image was built.
- **Size** – The total disk space used by the image.

If the build was successful, you should see `docker-angular-sample` image listed. 

---

### Run the containerized application

In the previous step, you created a Dockerfile for your Angular application and built a Docker image using the docker build command. Now it’s time to run that image in a container and verify that your application works as expected.


Inside the `docker-angular-sample` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple Angular web application.

Press `ctrl+c` in the terminal to stop your application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `docker-angular-sample` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see your Angular application running in the browser.


To confirm that the container is running, use `docker ps` command:

```console
$ docker ps
```

This will list all active containers along with their ports, names, and status. Look for a container exposing port 8080.

Example Output:

```shell
CONTAINER ID   IMAGE                          COMMAND                  CREATED             STATUS             PORTS                    NAMES
eb13026806d1   docker-angular-sample-server   "nginx -c /etc/nginx…"   About a minute ago  Up About a minute  0.0.0.0:8080->8080/tcp   docker-angular-sample-server-1
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

In this guide, you learned how to containerize, build, and run an Angular application using Docker. By following best practices, you created a secure, optimized, and production-ready setup.

What you accomplished:
- Created a multi-stage `Dockerfile` that compiles the Angular application and serves the static files using Nginx.
- Created a `.dockerignore` file to exclude unnecessary files and keep the image clean and efficient.
- Built your Docker image using `docker build`.
- Ran the container using `docker compose up`, both in the foreground and in detached mode.
- Verified that the app was running by visiting [http://localhost:8080](http://localhost:8080).
- Learned how to stop the containerized application using `docker compose down`.

You now have a fully containerized Angular application, running in a Docker container, and ready for deployment across any environment with confidence and consistency.

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

With your Angular application now containerized, you're ready to move on to the next step.

In the next section, you'll learn how to develop your application using Docker containers, enabling a consistent, isolated, and reproducible development environment across any machine.

## Use containers for Angular development

### Prerequisites

Complete [Containerize Angular application](containerize.md).

---

### Overview

In this section, you'll learn how to set up both production and development environments for your containerized Angular application using Docker Compose. This setup allows you to serve a static production build via Nginx and to develop efficiently inside containers using a live-reloading dev server with Compose Watch.

You’ll learn how to:
- Configure separate containers for production and development
- Enable automatic file syncing using Compose Watch in development
- Debug and live-preview your changes in real-time without manual rebuilds

---

### Automatically update services (development mode)

Use Compose Watch to automatically sync source file changes into your containerized development environment. This provides a seamless, efficient development experience without restarting or rebuilding containers manually.

### Step 1: Create a development Dockerfile

Create a file named `Dockerfile.dev` in your project root with the following content:

```dockerfile
# =========================================
# Stage 1: Development - Angular Application
# =========================================

# Define the Node.js version to use (Alpine for a small footprint)
ARG NODE_VERSION=24.12.0-alpine

# Set the base image for development
FROM node:${NODE_VERSION} AS dev

# Set environment variable to indicate development mode
ENV NODE_ENV=development

# Set the working directory inside the container
WORKDIR /app

# Copy only the dependency files first to optimize Docker caching
COPY package.json package-lock.json* ./

# Install dependencies using npm with caching to speed up subsequent builds
RUN --mount=type=cache,target=/root/.npm npm install

# Copy all application source files into the container
COPY . .

# Expose the port Angular uses for the dev server (default is 4200)
EXPOSE 4200

# Start the Angular dev server and bind it to all network interfaces
CMD ["npm", "start", "--", "--host=0.0.0.0"]

```

This file sets up a lightweight development environment for your Angular application using the dev server.


#### Step 2: Update your `compose.yaml` file

Open your `compose.yaml` file and define two services: one for production (`angular-prod`) and one for development (`angular-dev`).

Here’s an example configuration for an Angular application:

```yaml
services:
  angular-prod:
    build:
      context: .
      dockerfile: Dockerfile
    image: docker-angular-sample
    ports:
      - "8080:8080"

  angular-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "4200:4200"
    develop:
      watch:
        - action: sync
          path: .
          target: /app
```
- The `angular-prod` service builds and serves your static production app using Nginx.
- The `angular-dev` service runs your Angular development server with live reload and hot module replacement.
- `watch` triggers file sync with Compose Watch.

> [!NOTE]
> For more details, see the official guide: [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

After completing the previous steps, your project directory should now contain the following files:

```text
├── docker-angular-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ └── nginx.conf
```

#### Step 4: Start Compose Watch

Run the following command from the project root to start the container in watch mode

```console
$ docker compose watch angular-dev
```

#### Step 5: Test Compose Watch with Angular

To verify that Compose Watch is working correctly:

1. Open the `src/app/app.component.html` file in your text editor.

2. Locate the following line:

    ```html
    <h1>Docker Angular Sample Application</h1>
    ```

3. Change it to:

    ```html
    <h1>Hello from Docker Compose Watch</h1>
    ```

4. Save the file.

5. Open your browser at [http://localhost:4200](http://localhost:4200).

You should see the updated text appear instantly, without needing to rebuild the container manually. This confirms that file watching and automatic synchronization are working as expected.

---

### Summary

In this section, you set up a complete development and production workflow for your Angular application using Docker and Docker Compose.

Here’s what you accomplished:
- Created a `Dockerfile.dev` to streamline local development with hot reloading  
- Defined separate `angular-dev` and `angular-prod` services in your `compose.yaml` file  
- Enabled real-time file syncing using Compose Watch for a smoother development experience  
- Verified that live updates work seamlessly by modifying and previewing a component

With this setup, you're now equipped to build, run, and iterate on your Angular app entirely within containers—efficiently and consistently across environments.

---

### Related resources

Deepen your knowledge and improve your containerized development workflow with these guides:

- [Using Compose Watch](/manuals/compose/how-tos/file-watch.md) – Automatically sync source changes during development  
- [Multi-stage builds](/manuals/build/building/multi-stage.md) – Create efficient, production-ready Docker images  
- [Dockerfile best practices](/build/building/best-practices/) – Write clean, secure, and optimized Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.
- [Docker volumes](/storage/volumes/) – Persist and manage data between container runs  

### Next steps

In the next section, you'll learn how to run unit tests for your Angular application inside Docker containers. This ensures consistent testing across all environments and removes dependencies on local machine setup.

## Run Angular tests in a container

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize Angular application](containerize.md).

### Overview

Testing is a critical part of the development process. In this section, you'll learn how to:

- Run Jasmine unit tests using the Angular CLI inside a Docker container.
- Use Docker Compose to isolate your test environment.
- Ensure consistency between local and container-based testing.


The `docker-angular-sample` project comes pre-configured with Jasmine, so you can get started quickly without extra setup.

---

### Run tests during development

The `docker-angular-sample` application includes a sample test file at the following location:

```console
$ src/app/app.component.spec.ts
```

This test uses Jasmine to validate the AppComponent logic.

#### Step 1: Update compose.yaml

Add a new service named `angular-test` to your `compose.yaml` file. This service allows you to run your test suite in an isolated, containerized environment.

```yaml {hl_lines="22-26",linenos=true}
services:
  angular-dev:
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

  angular-prod:
    build:
      context: .
      dockerfile: Dockerfile
    image: docker-angular-sample
    ports:
      - "8080:8080"

  angular-test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    command: ["npm", "run", "test"]

```

The angular-test service reuses the same `Dockerfile.dev` used for [development](develop.md) and overrides the default command to run tests with `npm run test`. This setup ensures a consistent test environment that matches your local development configuration.


After completing the previous steps, your project directory should contain the following files:

```text
├── docker-angular-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ └── nginx.conf
```

#### Step 2: Run the tests

To execute your test suite inside the container, run the following command from your project root:

```console
$ docker compose run --rm angular-test
```

This command will:
- Start the `angular-test` service defined in your `compose.yaml` file.
- Execute the `npm run test` script using the same environment as development.
- Automatically removes the container after tests complete, using the [`docker compose run --rm`](/reference/cli/docker/compose/run/) command.

You should see output similar to the following:

```shell
Test Suites: 1 passed, 1 total
Tests:       3 passed, 3 total
Snapshots:   0 total
Time:        1.529 s
```

> [!NOTE]
> For more information about Compose commands, see the [Compose CLI
> reference](/reference/cli/docker/compose/).

---

### Summary

In this section, you learned how to run unit tests for your Angular application inside a Docker container using Jasmine and Docker Compose.

What you accomplished:
- Created a `angular-test` service in `compose.yaml` to isolate test execution.
- Reused the development `Dockerfile.dev` to ensure consistency between dev and test environments.
- Ran tests inside the container using `docker compose run --rm angular-test`.
- Ensured reliable, repeatable testing across environments without depending on your local machine setup.

---

### Related resources

Explore official references and best practices to sharpen your Docker testing workflow:

- [Dockerfile reference](/reference/dockerfile/) – Understand all Dockerfile instructions and syntax.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.  
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) – Run one-off commands in a service container.
---

### Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions to automatically build and test your Angular application in a containerized environment. This ensures your code is validated on every push or pull request, maintaining consistency and reliability across your development workflow.

## Test your Angular deployment

### Prerequisites

Before you begin, make sure you’ve completed the following:
- Complete all the previous sections of this guide, starting with [Containerize Angular application](containerize.md).
- [Enable Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

> **New to Kubernetes?**  
> Visit the [Kubernetes basics tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/) to get familiar with how clusters, pods, deployments, and services work.

---

### Overview

This section guides you through deploying your containerized Angular application locally using [Docker Desktop’s built-in Kubernetes](/desktop/kubernetes/). Running your app in a local Kubernetes cluster closely simulates a real production environment, enabling you to test, validate, and debug your workloads with confidence before promoting them to staging or production.

---

### Create a Kubernetes YAML file

Follow these steps to define your deployment configuration:

1. In the root of your project, create a new file named: angular-sample-kubernetes.yaml

2. Open the file in your IDE or preferred text editor.

3. Add the following configuration, and be sure to replace `{DOCKER_USERNAME}` and `{DOCKERHUB_PROJECT_NAME}` with your actual Docker Hub username and repository name from the previous [Automate your builds with GitHub Actions](configure-github-actions.md).


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: angular-sample
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: angular-sample
  template:
    metadata:
      labels:
        app: angular-sample
    spec:
      containers:
        - name: angular-container
          image: {DOCKER_USERNAME}/{DOCKERHUB_PROJECT_NAME}:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
          resources:
            limits:
              cpu: "500m"
              memory: "256Mi"
            requests:
              cpu: "250m"
              memory: "128Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: angular-sample-service
  namespace: default
spec:
  type: NodePort
  selector:
    app: angular-sample
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 30001
```

This manifest defines two key Kubernetes resources, separated by `---`:

- Deployment
  Deploys a single replica of your Angular application inside a pod. The pod uses the Docker image built and pushed by your GitHub Actions CI/CD workflow  
  (refer to [Automate your builds with GitHub Actions](configure-github-actions.md)).  
  The container listens on port `8080`, which is typically used by [Nginx](https://nginx.org/en/docs/) to serve your production Angular app.

- Service (NodePort) 
  Exposes the deployed pod to your local machine.  
  It forwards traffic from port `30001` on your host to port `8080` inside the container.  
  This lets you access the application in your browser at [http://localhost:30001](http://localhost:30001).

> [!NOTE]
> To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

---

### Deploy and check your application

Follow these steps to deploy your containerized Angular app into a local Kubernetes cluster and verify that it’s running correctly.

#### Step 1. Apply the Kubernetes configuration

In your terminal, navigate to the directory where your `angular-sample-kubernetes.yaml` file is located, then deploy the resources using:

```console
  $ kubectl apply -f angular-sample-kubernetes.yaml
```

If everything is configured properly, you’ll see confirmation that both the Deployment and the Service were created:

```shell
  deployment.apps/angular-sample created
  service/angular-sample-service created
```
   
This confirms that both the Deployment and the Service were successfully created and are now running inside your local cluster.

#### Step 2. Check the deployment status

Run the following command to check the status of your deployment:
   
```console
  $ kubectl get deployments
```

You should see output similar to the following:

```shell
  NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
  angular-sample       1/1     1            1           14s
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
angular-sample-service   NodePort    10.100.185.105    <none>        8080:30001/TCP   1m
```

This output confirms that your app is available via NodePort on port 30001.

#### Step 4. Access your app in the browser

Open your browser and navigate to [http://localhost:30001](http://localhost:30001).

You should see your production-ready Angular Sample application running — served by your local Kubernetes cluster.

#### Step 5. Clean up Kubernetes resources

Once you're done testing, you can delete the deployment and service using:

```console
  $ kubectl delete -f angular-sample-kubernetes.yaml
```

Expected output:

```shell
  deployment.apps "angular-sample" deleted
  service "angular-sample-service" deleted
```

This ensures your cluster stays clean and ready for the next deployment.
   
---

### Summary

In this section, you learned how to deploy your Angular application to a local Kubernetes cluster using Docker Desktop. This setup allows you to test and debug your containerized app in a production-like environment before deploying it to the cloud.

What you accomplished:

- Created a Kubernetes Deployment and NodePort Service for your Angular app  
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

Complete all the previous sections of this guide, starting with [Containerize an Angular application](containerize.md).

You must also have:
- A [GitHub](https://github.com/signup) account.
- A verified [Docker Hub](https://hub.docker.com/signup) account.

---

### Overview

In this section, you'll set up a CI/CD pipeline using [GitHub Actions](https://docs.github.com/en/actions) to automatically:

- Build your Angular application inside a Docker container.
- Run tests in a consistent environment.
- Push the production-ready image to [Docker Hub](https://hub.docker.com).

---

### Connect your GitHub repository to Docker Hub

To enable GitHub Actions to build and push Docker images, you’ll securely store your Docker Hub credentials in your new GitHub repository.

#### Step 1: Generate Docker Hub credentials and set GitHub secrets

1. Create a Personal Access Token (PAT) from [Docker Hub](https://hub.docker.com)
   1. Go to your **Docker Hub account → Account Settings → Security**.
   2. Generate a new Access Token with **Read/Write** permissions.
   3. Name it something like `docker-angular-sample`.
   4. Copy and save the token — you’ll need it in Step 4.

2. Create a repository in [Docker Hub](https://hub.docker.com/repositories/)
   1. Go to your **Docker Hub account → Create a repository**.
   2. For the Repository Name, use something descriptive — for example: `angular-sample`.
   3. Once created, copy and save the repository name — you’ll need it in Step 4.

3. Create a new [GitHub repository](https://github.com/new) for your Angular project

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

   These secrets allow GitHub Actions to authenticate securely with Docker Hub during automated workflows.

5. Connect Your Local Project to GitHub

   Link your local project `docker-angular-sample` to the GitHub repository you just created by running the following command from your project root:

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

6. Push your source code to GitHub

   Follow these steps to commit and push your local project to your GitHub repository:

   1. Stage all files for commit.

      ```console
      $ git add -A
      ```
      This command stages all changes — including new, modified, and deleted files — preparing them for commit.


   2. Commit the staged changes with a descriptive message.

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
name: CI/CD – Angular Application with Docker

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]

jobs:
  build-test-push:
    name: Build, Test, and Push Docker Image
    runs-on: ubuntu-latest

    steps:
      # 1. Checkout source code
      - name: Checkout source code
        uses: actions/checkout@{{% param "checkout_action_version" %}}
        with:
          fetch-depth: 0

      # 2. Set up Docker Buildx
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      # 3. Cache Docker layers
      - name: Cache Docker layers
        uses: actions/cache@{{% param "cache_action_version" %}}
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: |
            ${{ runner.os }}-buildx-

      # 4. Cache npm dependencies
      - name: Cache npm dependencies
        uses: actions/cache@{{% param "cache_action_version" %}}
        with:
          path: ~/.npm
          key: ${{ runner.os }}-npm-${{ hashFiles('**/package-lock.json') }}
          restore-keys: |
            ${{ runner.os }}-npm-

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
          load: true
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max

      # 7. Run Angular tests with Jasmine
      - name: Run Angular Jasmine tests inside container
        run: |
          docker run --rm \
            --workdir /app \
            --entrypoint "" \
            ${{ steps.meta.outputs.REPO_NAME }}-dev:latest \
            sh -c "npm ci && npm run test -- --ci --runInBand"
        env:
          CI: true
          NODE_ENV: test
        timeout-minutes: 10

      # 8. Log in to Docker Hub
      - name: Log in to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      # 9. Build and push production image
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

This workflow performs the following tasks for your Angular application:
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

   - Select "Commit changes…" in the GitHub editor.

   - This push will automatically trigger the GitHub Actions pipeline.

2. Monitor the workflow execution

   - Go to the Actions tab in your GitHub repository.
   - Click into the workflow run to follow each step: **build**, **test**, and (if successful) **push**.

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

In this section, you set up a complete CI/CD pipeline for your containerized Angular application using GitHub Actions.

Here's what you accomplished:

- Created a new GitHub repository specifically for your project.
- Generated a secure Docker Hub access token and added it to GitHub as a secret.
- Defined a GitHub Actions workflow that:
   - Build your application inside a Docker container.
   - Run tests in a consistent, containerized environment.
   - Push a production-ready image to Docker Hub if tests pass.
- Triggered and verified the workflow execution through GitHub Actions.
- Confirmed that your image was successfully published to Docker Hub.

With this setup, your Angular application is now ready for automated testing and deployment across environments — increasing confidence, consistency, and team productivity.

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

Next, learn how you can locally test and debug your Angular workloads on Kubernetes before deploying. This helps you ensure your application behaves as expected in a production-like environment, reducing surprises during deployment.
