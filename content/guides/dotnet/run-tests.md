---
title: Run .NET tests in a container
linkTitle: Run your tests
weight: 30
keywords: .NET, test
description: Learn how to run your .NET tests in a container.
aliases:
  - /language/dotnet/run-tests/
  - /guides/language/dotnet/run-tests/
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a .NET application](containerize.md).

## Overview

Testing is an essential part of modern software development. Testing can mean a
lot of things to different development teams. There are unit tests, integration
tests and end-to-end testing. In this guide you take a look at running your unit
tests in Docker when developing and when building.

## Run tests when developing locally

The sample application already has an xUnit test inside the `tests` directory. When developing locally, you can use Compose to run your tests.

Run the following command in the `docker-dotnet-sample` directory to run the tests inside a container.

```console
$ docker compose run --build --rm server dotnet test /source/tests
```

You should see output that contains the following.

```console
Starting test execution, please wait...
A total of 1 test files matched the specified pattern.

Passed!  - Failed:     0, Passed:     1, Skipped:     0, Total:     1, Duration: < 1 ms - /source/tests/bin/Debug/net8.0/tests.dll (net8.0)
```

To learn more about the command, see [docker compose run](/reference/cli/docker/compose/run/).

## Run tests when building

To run your tests when building, you need to update your Dockerfile. You can create a new test stage that runs the tests, or run the tests in the existing build stage. For this guide, update the Dockerfile to run the tests in the build stage.

The following is the updated Dockerfile.

```dockerfile {hl_lines="9"}
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS build
ARG TARGETARCH
COPY . /source
WORKDIR /source/src
RUN --mount=type=cache,id=nuget,target=/root/.nuget/packages \
    dotnet publish -a ${TARGETARCH/amd64/x64} --use-current-runtime --self-contained false -o /app
RUN dotnet test /source/tests

FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS development
COPY . /source
WORKDIR /source/src
CMD dotnet run --no-launch-profile

FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine AS final
WORKDIR /app
COPY --from=build /app .
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
ENTRYPOINT ["dotnet", "myWebApp.dll"]
```

Run the following command to build an image using the build stage as the target and view the test results. Include `--progress=plain` to view the build output, `--no-cache` to ensure the tests always run, and `--target build` to target the build stage.

```console
$ docker build -t dotnet-docker-image-test --progress=plain --no-cache --target build .
```

You should see output containing the following.

```console
#11 [build 5/5] RUN dotnet test /source/tests
#11 1.564   Determining projects to restore...
#11 3.421   Restored /source/src/myWebApp.csproj (in 1.02 sec).
#11 19.42   Restored /source/tests/tests.csproj (in 17.05 sec).
#11 27.91   myWebApp -> /source/src/bin/Debug/net8.0/myWebApp.dll
#11 28.47   tests -> /source/tests/bin/Debug/net8.0/tests.dll
#11 28.49 Test run for /source/tests/bin/Debug/net8.0/tests.dll (.NETCoreApp,Version=v8.0)
#11 28.67 Microsoft (R) Test Execution Command Line Tool Version 17.3.3 (x64)
#11 28.67 Copyright (c) Microsoft Corporation.  All rights reserved.
#11 28.68
#11 28.97 Starting test execution, please wait...
#11 29.03 A total of 1 test files matched the specified pattern.
#11 32.07
#11 32.08 Passed!  - Failed:     0, Passed:     1, Skipped:     0, Total:     1, Duration: < 1 ms - /source/tests/bin/Debug/net8.0/tests.dll (net8.0)
#11 DONE 32.2s
```

## Summary

In this section, you learned how to run tests when developing locally using Compose and how to run tests when building your image.

Related information:

- [docker compose run](/reference/cli/docker/compose/run/)

## Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions.
