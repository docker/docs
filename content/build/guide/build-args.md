---
title: Build arguments
description: Introduction to configurable builds, using build args
keywords: build, buildkit, buildx, guide, tutorial, build arguments, arg
---

Build arguments is a great way to add flexibility to your builds. You can pass
build arguments at build-time, and you can set a default value that the builder
uses as a fallback.

## Change runtime versions

A practical use case for build arguments is to specify runtime versions for
build stages. Your image uses the `golang:{{% param "example_go_version" %}}-alpine`
image as a base image.
But what if someone wanted to use a different version of Go for building the
application? They could update the version number inside the Dockerfile, but
that’s inconvenient, it makes switching between versions more tedious than it
has to be. Build arguments make life easier:

```diff
  # syntax=docker/dockerfile:1
- FROM golang:{{% param "example_go_version" %}}-alpine AS base
+ ARG GO_VERSION={{% param "example_go_version" %}}
+ FROM golang:${GO_VERSION}-alpine AS base
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
  RUN --mount=type=cache,target=/go/pkg/mod/ \
      --mount=type=bind,target=. \
      go build -o /bin/server ./cmd/server

  FROM scratch AS client
  COPY --from=build /bin/client /bin/
  ENTRYPOINT [ "/bin/client" ]

  FROM scratch AS server
  COPY --from=build /bin/server /bin/
  ENTRYPOINT [ "/bin/server" ]
```

The `ARG` keyword is interpolated in the image name in the `FROM` instruction.
The default value of the `GO_VERSION` build argument is set to `{{% param "example_go_version" %}}`.
If the build doesn't receive a `GO_VERSION` build argument, the `FROM` instruction
resolves to `golang:{{% param "example_go_version" %}}-alpine`.

Try setting a different version of Go to use for building, using the
`--build-arg` flag for the build command:

```console
$ docker build --build-arg="GO_VERSION=1.19" .
```

Running this command results in a build using the `golang:1.19-alpine` image.

## Inject values

You can also make use of build arguments to modify values in the source code of
your program, at build time. This is useful for dynamically injecting
information, avoiding hard-coded values. With Go, consuming external values at
build time is done using linker flags, or `-ldflags`.

The server part of the application contains a conditional statement to print the
app version, if a version is specified:

```go
// cmd/server/main.go
var version string

func main() {
	if version != "" {
		log.Printf("Version: %s", version)
	}
```

You could declare the version string value directly in the code. But, updating
the version to line up with the release version of the application would require
updating the code ahead of every release. That would be both tedious and
error-prone. A better solution is to pass the version string as a build
argument, and inject the build argument into the code.

The following example adds an `APP_VERSION` build argument to the `build-server`
stage. The Go compiler uses the value of the build argument to set the value of
a variable in the code.

```diff
  # syntax=docker/dockerfile:1
  ARG GO_VERSION={{% param "example_go_version" %}}
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
+ ARG APP_VERSION="v0.0.0+unknown"
  RUN --mount=type=cache,target=/go/pkg/mod/ \
      --mount=type=bind,target=. \
-     go build -o /bin/server ./cmd/server
+     go build -ldflags "-X main.version=$APP_VERSION" -o /bin/server ./cmd/server

  FROM scratch AS client
  COPY --from=build-client /bin/client /bin/
  ENTRYPOINT [ "/bin/client" ]

  FROM scratch AS server
  COPY --from=build-server /bin/server /bin/
  ENTRYPOINT [ "/bin/server" ]
```

Now the version of the server is injected when building the binary, without having to update
the source code. To verify this, you can build the `server` target and start a
container with `docker run`. The server outputs `v0.0.1` as the version on
startup.

```console
$ docker build --target=server --build-arg="APP_VERSION=v0.0.1" --tag=buildme-server .
$ docker run buildme-server
2023/04/06 08:54:27 Version: v0.0.1
2023/04/06 08:54:27 Starting server...
2023/04/06 08:54:27 Listening on HTTP port 3000
```

## Summary

This section showed how you can use build arguments to make builds more
configurable, and inject values at build-time.

Related information:

- [`ARG` Dockerfile reference](../../engine/reference/builder.md#arg)

## Next steps

The next section of this guide shows how you can use Docker builds to create not
only container images, but executable binaries as well.

{{< button text="Export binaries" url="export.md" >}}
