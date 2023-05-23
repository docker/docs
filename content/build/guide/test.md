---
title: Test
description: Running tests with Docker Build
keywords: build, buildkit, buildx, guide, tutorial, testing
---

{% include_relative nav.html selected="7" %}

This section focuses on testing. The example in this section focuses on linting,
but the same principles apply for other kinds of tests as well, such as unit
tests. Code linting is a static analysis of code that helps you detect errors,
style violations, and anti-patterns.

The exact steps for how to test your code can vary a lot depending on the
programming language or framework that you use. The example application used in
this guide is written in Go. You will add a build step that uses
`golangci-lint`, a popular linters runner for Go.

## Run tests

The `golangci-lint` tool is available as an image on Docker Hub. Before you add
the lint step to the Dockerfile, you can try it out using a `docker run`
command.

```console
$ docker run -v $PWD:/test -w /test \
  golangci/golangci-lint golangci-lint run
```

You will notice that `golangci-lint` works: it finds an issue in the code where
there's a missing error check.

```text
cmd/server/main.go:23:10: Error return value of `w.Write` is not checked (errcheck)
		w.Write([]byte(translated))
		      ^
```

Now you can add this as a step to the Dockerfile.

```diff
  # syntax=docker/dockerfile:1
  ARG GO_VERSION={{site.example_go_version}}
+ ARG GOLANGCI_LINT_VERSION={{site.example_golangci_lint_version}}
  FROM golang:${GO_VERSION}-alpine AS base
  WORKDIR /src
  RUN --mount=type=cache,target=/go/pkg/mod/ \
      --mount=type=bind,source=go.sum,target=go.sum \
      --mount=type=bind,source=go.mod,target=go.mod \
      go mod download -x

  FROM base AS build-client
  RUN --mount=type=cache,target=/go/pkg/mod/ \
      --mount=type=bind,target=. \
      go build -o /bin/client ./cmd/client

  FROM base AS build-server
  ARG APP_VERSION="0.0.0+unknown"
  RUN --mount=type=cache,target=/go/pkg/mod/ \
      --mount=type=bind,target=. \
      go build -ldflags "-X main.version=$APP_VERSION" -o /bin/server ./cmd/server

  FROM scratch AS client
  COPY --from=build-client /bin/client /bin/
  ENTRYPOINT [ "/bin/client" ]

  FROM scratch AS server
  COPY --from=build-server /bin/server /bin/
  ENTRYPOINT [ "/bin/server" ]

  FROM scratch AS binaries
  COPY --from=build-client /bin/client /
  COPY --from=build-server /bin/server /
+
+ FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} as lint
+ WORKDIR /test
+ RUN --mount=type=bind,target=. \
+     golangci-lint run
```

The added `lint` stage uses the `golangci/golangci-lint` image from Docker Hub
to invoke the `golangci-lint run` command with a bind-mount for the build
context.

The lint stage is independent of any of the other stages in the Dockerfile.
Therefore, running a regular build won’t cause the lint step to run. To lint the
code, you must specify the `lint` stage:

```console
$ docker build --target=lint .
```

## Export test results

In addition to running tests, it's sometimes useful to be able to export the
results of a test to a test report.

Exporting test results is no different to exporting binaries, as shown in the
previous section of this guide:

1. Save the test results to a file.
2. Create a new stage in your Dockerfile using the `scratch` base image.
3. Export that stage using the `local` exporter.

The exact steps on how to do this is left as a reader's exercise :-)

## Summary

This section has shown an example on how you can use Docker builds to run tests
(or as shown in this section, linters).

## Next steps

The next topic in this guide is multi-platform builds, using emulation and
cross-compilation.

[Multi-platform](multi-platform.md){: .button .primary-btn }
