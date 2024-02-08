---
title: Multi-platform
description: Building for multiple operating systems and architectures
keywords: build, buildkit, buildx, guide, tutorial, multi-platform, emulation, cross-compilation
---

Up until this point in the guide, you've built Linux binaries. This section
describes how you can support other operating systems, and architectures, using
multi-platform builds via emulation and cross-compilation.

The easiest way to get started with building for multiple platforms is using
emulation. With emulation, you can build your app to multiple architectures
without having to make any changes to your Dockerfile. All you need to do is to
pass the `--platform` flag to the build command, specifying the OS and
architecture you want to build for.

The following command builds the server image for the `linux/arm/v7` platform:

```console
$ docker build --target=server --platform=linux/arm/v7 .
```

You can also use emulation to produce outputs for multiple platforms at once.
However, the default image store in Docker Engine doesn't support building
and loading multi-platform images. You need to enable the containerd image store
which supports concurrent multi-platform builds.

## Enable the containerd image store

{{< tabs >}}
{{< tab name="Docker Desktop" >}}

To enable the containerd image store in Docker Desktop,
go to **Settings** and select **Use containerd for pulling and storing images**
in the **General** tab.

Note that changing the image store means you'll temporarily lose access to
images and containers in the classic image store.
Those resources still exist, but to view them, you'll need to
disable the containerd image store.

{{< /tab >}}
{{< tab name="Docker Engine" >}}

If you're not using Docker Desktop,
enable the containerd image store by adding the following feature configuration
to your `/etc/docker/daemon.json` configuration file.

```json {hl_lines=3}
{
  "features": {
    "containerd-snapshotters": true
  }
}
```

Restart the daemon after updating the configuration file.

```console
$ systemctl restart docker
```

{{< /tab >}}
{{< /tabs >}}

## Build using emulation

To run multi-platform builds, invoke the `docker build` command,
and pass it the same arguments as you did before.
Only this time, also add a `--platform` flag specifying multiple architectures.

```console {hl_lines=4}
$ docker build \
    --target=binaries \
    --output=bin \
    --platform=linux/amd64,linux/arm64,linux/arm/v7 .
```

This command uses emulation to run the same build four times, once for each
platform. The build results are exported to a `bin` directory.

```text
bin
├── linux_amd64
│   ├── client
│   └── server
├── linux_arm64
│   ├── client
│   └── server
└── linux_arm_v7
    ├── client
    └── server
```

When you build for multiple platforms concurrently,
BuildKit runs all of the build steps under emulation for each platform that you specify.
Effectively forking the build into two concurrent processes.

![Build pipelines using emulation](./images/emulation.png)

There are, however, a few downsides to running multi-platform builds using
emulation:

- If you tried running the command above, you may have noticed that it took a
  long time to finish. Emulation can be much slower than native execution for
  CPU-intensive tasks.
- Emulation only works when the architecture is supported by the base image
  you’re using. The example in this guide uses the Alpine Linux version of the
  `golang` image, which means you can only build Linux images this way, for a
  limited set of CPU architectures, without having to change the base image.

As an alternative to emulation, the next step explores cross-compilation.
Cross-compiling makes multi-platform builds much faster and versatile.

## Build using cross-compilation

Using cross-compilation means leveraging the capabilities of a compiler to build
for multiple platforms, without the need for emulation.

The first thing you'll need to do is pinning the builder to use the node’s
native architecture as the build platform. This is to prevent emulation. Then,
from the node's native architecture, the builder cross-compiles the application
to a number of other target platforms.

### Platform build arguments

This approach involves using a few pre-defined build arguments that you have
access to in your Docker builds: `BUILDPLATFORM` and `TARGETPLATFORM` (and
derivatives, like `TARGETOS`). These build arguments reflect the values you pass
to the `--platform` flag.

For example, if you invoke a build with `--platform=linux/amd64`, then the build
arguments resolve to:

- `TARGETPLATFORM=linux/amd64`
- `TARGETOS=linux`
- `TARGETARCH=amd64`

When you pass more than one value to the platform flag, build stages that use
the pre-defined platform arguments are forked automatically for each platform.
This is in contrast to builds running under emulation, where the entire build
pipeline runs per platform.

![Build pipelines using cross-compilation](./images/cross-compilation.png)

### Update the Dockerfile

To build the app using the cross-compilation technique, update the Dockerfile as
follows:

- Add `--platform=$BUILDPLATFORM` to the `FROM` instruction for the initial
  `base` stage, pinning the platform of the `golang` image to match the
  architecture of the host machine.
- Add `ARG` instructions for the Go compilation stages, making the `TARGETOS`
  and `TARGETARCH` build arguments available to the commands in this stage.
- Set the `GOOS` and `GOARCH` environment variables to the values of `TARGETOS`
  and `TARGETARCH`. The Go compiler uses these variables to do
  cross-compilation.

```diff
  # syntax=docker/dockerfile:1
  ARG GO_VERSION={{% param "example_go_version" %}}
  ARG GOLANGCI_LINT_VERSION={{% param "example_golangci_lint_version" %}}
- FROM golang:${GO_VERSION}-alpine AS base
+ FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
  WORKDIR /src
  RUN --mount=type=cache,target=/go/pkg/mod \
      --mount=type=bind,source=go.mod,target=go.mod \
      --mount=type=bind,source=go.sum,target=go.sum \
      go mod download -x

  FROM base AS build-client
+ ARG TARGETOS
+ ARG TARGETARCH
  RUN --mount=type=cache,target=/go/pkg/mod \
      --mount=type=bind,target=. \
-     go build -o /bin/client ./cmd/client
+     GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/client ./cmd/client

  FROM base AS build-server
+ ARG TARGETOS
+ ARG TARGETARCH
  RUN --mount=type=cache,target=/go/pkg/mod \
      --mount=type=bind,target=. \
-     go build -o /bin/server ./cmd/server
+     GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/server ./cmd/server

  FROM scratch AS client
  COPY --from=build-client /bin/client /bin/
  ENTRYPOINT [ "/bin/client" ]

  FROM scratch AS server
  COPY --from=build-server /bin/server /bin/
  ENTRYPOINT [ "/bin/server" ]

  FROM scratch AS binaries
  COPY --from=build-client /bin/client /
  COPY --from=build-server /bin/server /

  FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} as lint
  WORKDIR /test
  RUN --mount=type=bind,target=. \
      golangci-lint run
```

The only thing left to do now is to run the actual build. To run a
multi-platform build, set the `--platform` option, and specify a CSV string of
the OS and architectures that you want to build for. The following command
illustrates how to build, and export, binaries for Mac (ARM64), Windows, and
Linux:

```console
$ docker build \
  --target=binaries \
  --output=bin \
  --platform=darwin/arm64,windows/amd64,linux/amd64 .
```

When the build finishes, you’ll find client and server binaries for all of the
selected platforms in the `bin` directory:

```diff
bin
├── darwin_arm64
│   ├── client
│   └── server
├── linux_amd64
│   ├── client
│   └── server
└── windows_amd64
    ├── client
    └── server
```

## Summary

This section has demonstrated how you can get started with multi-platform builds
using emulation and cross-compilation.

Related information:

- [Multi-platfom images](../building/multi-platform.md)
- [containerd image store (Docker Desktop)](../../desktop/containerd.md)
- [containerd image store (Docker Engine)](../../storage/containerd.md)

You may also want to consider checking out
[xx - Dockerfile cross-compilation helpers](https://github.com/tonistiigi/xx).
`xx` is a Docker image containing utility scripts that make cross-compiling with Docker builds easier.

## Next steps

This section is the final part of the Build with Docker guide. The following
page contains some pointers for where to go next.

{{< button text="Next steps" url="next-steps.md" >}}
