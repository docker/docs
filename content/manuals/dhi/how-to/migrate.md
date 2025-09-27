---
title: Migrate an existing application to use Docker Hardened Images
linktitle: Migrate an app
description: Follow a step-by-step guide to update your Dockerfiles and adopt Docker Hardened Images for secure, minimal, and production-ready builds.
weight: 50
keywords: migrate dockerfile, hardened base image, multi-stage build, non-root containers, secure container build
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

This guide helps you migrate your existing Dockerfiles to use Docker Hardened
Images (DHIs) [manually](#step-1-update-the-base-image-in-your-dockerfile),
or with [Gordon](#use-gordon).
DHIs are minimal and security-focused, which may require
adjustments to your base images, build process, and runtime configuration.

This guide focuses on migrating framework images, such as images for building
applications from source using languages like Go, Python, or Node.js. If you're
migrating application images, such as databases, proxies, or other prebuilt
services, many of the same principles still apply.

## Migration considerations

DHIs omit common tools such as shells and package managers to
reduce the attack surface. They also default to running as a nonroot user. As a
result, migrating to DHI typically requires the following changes to your
Dockerfile:


| Item               | Migration note                                                                                                                                                                                                                                                                                                                 |
|:-------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Base image         | Replace your base images in your Dockerfile with a Docker Hardened Image.                                                                                                                                                                                                                                                      |
| Package management | Images intended for runtime, don't contain package managers. Use package managers only in images with a `dev` tag. Utilize multi-stage builds and copy necessary artifacts from the build stage to the runtime stage.                                                                                                                                                                        |
| Non-root user      | By default, images intended for runtime, run as the nonroot user. Ensure that necessary files and directories are accessible to the nonroot user.                                                                                                                                                                              |
| Multi-stage build  | Utilize images with a `dev` or `sdk` tags for build stages and non-dev images for runtime.                                                                                                                                                                                                                                     |
| TLS certificates   | DHIs contain standard TLS certificates by default. There is no need to install TLS certificates.                                                                                                                                                                                                                               |
| Ports              | DHIs intended for runtime run as a nonroot user by default. As a result, applications in these images can't bind to privileged ports (below 1024) when running in Kubernetes or in Docker Engine versions older than 20.10. To avoid issues, configure your application to listen on port 1025 or higher inside the container. |
| Entry point        | DHIs may have different entry points than images such as Docker Official Images. Inspect entry points for DHIs and update your Dockerfile if necessary.                                                                                                                                                                        |
| No shell           | DHIs intended for runtime don't contain a shell. Use dev images in build stages to run shell commands and then copy artifacts to the runtime stage.                                                                                                                                                                            |

For more details and troubleshooting tips, see the [Troubleshoot](/manuals/dhi/troubleshoot.md).

## Migrate an existing application

The following steps outline the migration process.

### Step 1: Update the base image in your Dockerfile

Update the base image in your application’s Dockerfile to a hardened image. This
is typically going to be an image tagged as `dev` or `sdk` because it has the tools
needed to install packages and dependencies.

The following example diff snippet from a Dockerfile shows the old base image
replaced by the new hardened image.

```diff
- ## Original base image
- FROM golang:1.22

+ ## Updated to use hardened base image
+ FROM <your-namespace>/dhi-golang:1.22-dev
```

### Step 2: Update the runtime image in your Dockerfile

> [!NOTE]
>
> Multi-stage builds are recommended to keep your final image minimal and
> secure. Single-stage builds are supported, but they include the full `dev` image
> and therefore result in a larger image with a broader attack surface.

To ensure that your final image is as minimal as possible, you should use a
[multi-stage build](/manuals/build/building/multi-stage.md). All stages in your
Dockerfile should use a hardened image. While intermediary stages will typically
use images tagged as `dev` or `sdk`, your final runtime stage should use a runtime image.

Utilize the build stage to compile your application and copy the resulting
artifacts to the final runtime stage. This ensures that your final image is
minimal and secure.

See the [Example Dockerfile migrations](#example-dockerfile-migrations) section for
examples of how to update your Dockerfile.

## Example Dockerfile migrations

The following examples show a Dockerfile before and after migration. Each
example includes both a multi-stage build (recommended for minimal, secure
images) and a single-stage build (supported, but results in a larger image with
a broader attack surface).

> [!NOTE]
>
> Multi-stage builds are recommended for most use cases. Single-stage builds are
> supported for simplicity, but come with tradeoffs in size and security.

### Go example

{{< tabs >}}
{{< tab name="Before" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< tab name="After (multi-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

# === Build stage: Compile Go application ===
FROM <your-namespace>/dhi-golang:1-alpine3.21-dev AS builder

WORKDIR /app
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

# === Final stage: Create minimal runtime image ===
FROM <your-namespace>/dhi-golang:1-alpine3.21

WORKDIR /app
COPY --from=builder /app/main  /app/main

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< tab name="After (single-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM <your-namespace>/dhi-golang:1-alpine3.21-dev

WORKDIR /app
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< /tabs >}}

### Node.js example

{{< tabs >}}
{{< tab name="Before" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM node:latest
WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install

COPY image.jpg ./image.jpg
COPY . .

CMD ["node", "index.js"]
```

{{< /tab >}}
{{< tab name="After (multi-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

#=== Build stage: Install dependencies and build application ===#
FROM <your-namespace>/dhi-node:23-alpine3.21-dev AS builder
WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install

COPY image.jpg ./image.jpg
COPY . .

#=== Final stage: Create minimal runtime image ===#
FROM <your-namespace>/dhi-node:23-alpine3.21
ENV PATH=/app/node_modules/.bin:$PATH

COPY --from=builder --chown=node:node /usr/src/app /app

WORKDIR /app

CMD ["index.js"]
```

{{< /tab >}}
{{< tab name="After (single-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM <your-namespace>/dhi-node:23-alpine3.21-dev
WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install

COPY image.jpg ./image.jpg
COPY . .

CMD ["index.js"]
```

{{< /tab >}}
{{< /tabs >}}

### Python example

{{< tabs >}}
{{< tab name="Before" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM python:latest AS builder

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PATH="/app/venv/bin:$PATH"

WORKDIR /app

RUN python -m venv /app/venv
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

FROM python:latest

WORKDIR /app

ENV PYTHONUNBUFFERED=1
ENV PATH="/app/venv/bin:$PATH"

COPY image.py image.png ./
COPY --from=builder /app/venv /app/venv

ENTRYPOINT [ "python", "/app/image.py" ]
```

{{< /tab >}}
{{< tab name="After (multi-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

#=== Build stage: Install dependencies and create virtual environment ===#
FROM <your-namespace>/dhi-python:3.13-alpine3.21-dev AS builder

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PATH="/app/venv/bin:$PATH"

WORKDIR /app

RUN python -m venv /app/venv
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

#=== Final stage: Create minimal runtime image ===#
FROM <your-namespace>/dhi-python:3.13-alpine3.21

WORKDIR /app

ENV PYTHONUNBUFFERED=1
ENV PATH="/app/venv/bin:$PATH"

COPY image.py image.png ./
COPY --from=builder /app/venv /app/venv

ENTRYPOINT [ "python", "/app/image.py" ]
```

{{< /tab >}}
{{< tab name="After (single-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM <your-namespace>/dhi-python:3.13-alpine3.21-dev

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PATH="/app/venv/bin:$PATH"

WORKDIR /app

RUN python -m venv /app/venv
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY image.py image.png ./

ENTRYPOINT [ "python", "/app/image.py" ]
```

{{< /tab >}}
{{< /tabs >}}

### Use Gordon

Alternatively, you can request assistance to
[Gordon](/manuals/ai/gordon/_index.md), Docker's AI-powered assistant, to
migrate your Dockerfile:

{{% include "gordondhi.md" %}}
