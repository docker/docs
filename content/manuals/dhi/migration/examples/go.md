---
title: Go
description: Migrate a Go application to Docker Hardened Images
weight: 10
keywords: go, golang, migration, dhi
---

This example shows how to migrate a Go application to Docker Hardened Images.

The following examples show Dockerfiles before and after migration to Docker
Hardened Images. Each example includes four variations:

- Before (Wolfi): A sample Dockerfile using Wolfi distribution images, before migrating to DHI
- Before (DOI): A sample Dockerfile using Docker Official Images, before migrating to DHI
- After (multi-stage): A sample Dockerfile after migrating to DHI with multi-stage builds (recommended for minimal, secure images)
- After (single-stage): A sample Dockerfile after migrating to DHI with single-stage builds (simpler but results in a larger image with a broader attack surface)

> [!NOTE]
>
> Multi-stage builds are recommended for most use cases. Single-stage builds are
> supported for simplicity, but come with tradeoffs in size and security.
>
> You must authenticate to `dhi.io` before you can pull Docker Hardened Images.
> Use your Docker ID credentials (the same username and password you use for
> Docker Hub). If you don't have a Docker account, [create
> one](../../../accounts/create-account.md) for free.
>
> Run `docker login dhi.io` to authenticate.

{{< tabs >}}
{{< tab name="Before (Wolfi)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM cgr.dev/chainguard/go:latest-dev

WORKDIR /app
ADD . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< tab name="Before (DOI)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app
ADD . ./

# Install any additional packages if needed using apt
# RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< tab name="After (multi-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

# === Build stage: Compile Go application ===
FROM dhi.io/golang:1-alpine3.21-dev AS builder

WORKDIR /app
ADD . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

# === Final stage: Create minimal runtime image ===
FROM dhi.io/golang:1-alpine3.21

WORKDIR /app
COPY --from=builder /app/main  /app/main

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< tab name="After (single-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM dhi.io/golang:1-alpine3.21-dev

WORKDIR /app
ADD . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" --installsuffix cgo -o main .

ENTRYPOINT ["/app/main"]
```

{{< /tab >}}
{{< /tabs >}}
