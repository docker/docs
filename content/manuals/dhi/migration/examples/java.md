---
title: Java
description: Migrate a Java application to Docker Hardened Images
weight: 50
keywords: java, jvm, jdk, jre, maven, migration, dhi
---

This example shows how to migrate a Java application to Docker Hardened Images.

The following examples show Dockerfiles before and after migration to Docker
Hardened Images. Each example includes five variations:

- Before (Ubuntu): A sample Dockerfile using Ubuntu-based images, before migrating to DHI
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
{{< tab name="Before (Ubuntu)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM ubuntu/jre:21-24.04 AS builder

WORKDIR /app
COPY . ./

# Install any additional packages if needed using apt
# RUN apt-get update && apt-get install -y maven && rm -rf /var/lib/apt/lists/*

RUN mvn -B package -DskipTests

FROM ubuntu/jre:21-24.04

WORKDIR /app
COPY --from=builder /app/target/app.jar /app/app.jar

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

{{< /tab >}}
{{< tab name="Before (Wolfi)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM cgr.dev/chainguard/maven:latest-dev AS builder

WORKDIR /app
COPY . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN mvn -B package -DskipTests

FROM cgr.dev/chainguard/jre:latest

WORKDIR /app
COPY --from=builder /app/target/app.jar /app/app.jar

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

{{< /tab >}}
{{< tab name="Before (DOI)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM maven:3.9-eclipse-temurin-21 AS builder

WORKDIR /app
COPY . ./

# Install any additional packages if needed using apt
# RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

RUN mvn -B package -DskipTests

FROM eclipse-temurin:21-jre

WORKDIR /app
COPY --from=builder /app/target/app.jar /app/app.jar

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

{{< /tab >}}
{{< tab name="After (multi-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

# === Build stage: Compile and package the Java application with Maven ===
FROM dhi.io/maven:3-alpine3.21-dev AS builder

WORKDIR /app
COPY . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN mvn -B package -DskipTests

# === Final stage: Create minimal runtime image ===
FROM dhi.io/eclipse-temurin:21-alpine3.21

WORKDIR /app
COPY --from=builder /app/target/app.jar /app/app.jar

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

{{< /tab >}}
{{< tab name="After (single-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM dhi.io/maven:3-alpine3.21-dev

WORKDIR /app
COPY . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN mvn -B package -DskipTests

ENTRYPOINT ["java", "-jar", "/app/target/app.jar"]
```

{{< /tab >}}
{{< /tabs >}}
