---
title: .NET
description: Migrate a .NET application to Docker Hardened Images
weight: 40
keywords: dotnet, .net, csharp, aspnet, migration, dhi
---

> **Acknowledgment**
>
> Docker would like to thank [Naga Santhosh Reddy Vootukuri](https://github.com/sunnynagavo) for his contribution to this guide.

This example shows how to migrate a .NET application to Docker Hardened Images.

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

FROM cgr.dev/chainguard/dotnet-sdk:latest-dev AS builder

WORKDIR /src
COPY . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN dotnet restore
RUN dotnet publish -c Release -o /src/out --no-restore

FROM cgr.dev/chainguard/aspnet-runtime:latest

WORKDIR /app
COPY --from=builder /src/out ./

ENTRYPOINT ["dotnet", "app.dll"]
```

{{< /tab >}}
{{< tab name="Before (DOI)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder

WORKDIR /src
COPY . ./

# Install any additional packages if needed using apt
# RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

RUN dotnet restore
RUN dotnet publish -c Release -o /app --no-restore

FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app
COPY --from=builder /app ./

ENTRYPOINT ["dotnet", "app.dll"]
```

{{< /tab >}}
{{< tab name="After (multi-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

# === Build stage: Restore, build, and publish the .NET application ===
FROM dhi.io/dotnet:8-sdk-alpine3.22 AS builder

WORKDIR /src
COPY . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN dotnet restore
RUN dotnet publish -c Release -o /app --no-restore

# === Final stage: Create minimal runtime image ===
FROM dhi.io/aspnetcore:8-alpine3.22

WORKDIR /app
COPY --from=builder /app ./

ENTRYPOINT ["dotnet", "app.dll"]
```

{{< /tab >}}
{{< tab name="After (single-stage)" >}}

```dockerfile
#syntax=docker/dockerfile:1

FROM dhi.io/dotnet:8-sdk-alpine3.22

WORKDIR /src
COPY . ./

# Install any additional packages if needed using apk
# RUN apk add --no-cache git

RUN dotnet restore
RUN dotnet publish -c Release -o /app --no-restore

WORKDIR /app

ENTRYPOINT ["dotnet", "/app/app.dll"]
```

{{< /tab >}}
{{< /tabs >}}
