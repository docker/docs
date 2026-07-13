---
title: Local and tar exporters
keywords: build, buildx, buildkit, exporter, local, tar
description: >
  The local and tar exporters save the build result to the local filesystem
aliases:
  - /build/building/exporters/local-tar/
---

The `local` and `tar` exporters output the root filesystem of the build result
into a local directory. They're useful for producing artifacts that aren't
container images.

- `local` exports files and directories.
- `tar` exports the same, but bundles the export into a tarball.

## Synopsis

Build a container image using the `local` exporter:

```console
$ docker buildx build --output type=local[,parameters] .
$ docker buildx build --output type=tar[,parameters] .
```

The following table describes the available parameters:

| Parameter        | Type    | Default | Description                                                                       |
| ---------------- | ------- | ------- | --------------------------------------------------------------------------------- |
| `dest`           | String  |         | Path to copy files to                                                             |
| `platform-split` | Boolean | `true`  | `local` exporter only. Split multi-platform outputs into platform subdirectories. |

## Multi-platform builds with local exporter

The `platform-split` parameter controls how multi-platform build outputs are
organized.

Consider this Dockerfile that creates platform-specific files:

```dockerfile
FROM busybox AS build
ARG TARGETOS
ARG TARGETARCH
RUN mkdir /out && echo foo > /out/hello-$TARGETOS-$TARGETARCH

FROM scratch
COPY --from=build /out /
```

### Split by platform (default)

By default, the local exporter creates a separate subdirectory for each
platform:

```console
$ docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --output type=local,dest=./output \
  .
```

This produces the following directory structure:

```text
output/
├── linux_amd64/
│   └── hello-linux-amd64
└── linux_arm64/
    └── hello-linux-arm64
```

### Merge all platforms

To merge files from all platforms into the same directory, set
`platform-split=false`:

```console
$ docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --output type=local,dest=./output,platform-split=false \
  .
```

This produces a flat directory structure:

```text
output/
├── hello-linux-amd64
└── hello-linux-arm64
```

Files from all platforms merge into a single directory. If multiple platforms
produce files with identical names, the export fails with an error.

### Single-platform builds

Single-platform builds export directly to the destination directory without
creating a platform subdirectory:

```console
$ docker buildx build \
  --platform linux/amd64 \
  --output type=local,dest=./output \
  .
```

This produces:

```text
output/
└── hello-linux-amd64
```

To include the platform subdirectory even for single-platform builds, explicitly
set `platform-split=true`:

```console
$ docker buildx build \
  --platform linux/amd64 \
  --output type=local,dest=./output,platform-split=true \
  .
```

This produces:

```text
output/
└── linux_amd64/
    └── hello-linux-amd64
```

## Further reading

For more information on the `local` or `tar` exporters, see the
[BuildKit README](https://github.com/moby/buildkit/blob/master/README.md#local-directory).
