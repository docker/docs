---
title: Exporters overview
linkTitle: Exporters
weight: 90
description: Build exporters define the output format of your build result
keywords: build, buildx, buildkit, exporter, image, registry, local, tar, oci, docker, cacheonly
aliases:
  - /build/building/exporters/
---

Exporters save your build results to a specified output type. You specify the
exporter to use with the
[`--output` CLI option](/reference/cli/docker/buildx/build.md#output).
Buildx supports the following exporters:

- `image`: exports the build result to a container image.
- `registry`: exports the build result into a container image, and pushes it to
  the specified registry.
- `local`: exports the build root filesystem into a local directory.
- `tar`: packs the build root filesystem into a local tarball.
- `oci`: exports the build result to the local filesystem in the
  [OCI image layout](https://github.com/opencontainers/image-spec/blob/v1.0.1/image-layout.md)
  format.
- `docker`: exports the build result to the local filesystem in the
  [Docker Image Specification v1.2.0](https://github.com/moby/moby/blob/v25.0.0/image/spec/v1.2.md)
  format.
- `cacheonly`: doesn't export a build output, but runs the build and creates a
  cache.

## Using exporters

To specify an exporter, use the following command syntax:

```console
$ docker buildx build --tag <registry>/<image> \
  --output type=<TYPE> .
```

Most common use cases don't require that you specify which exporter to use
explicitly. You only need to specify the exporter if you intend to customize
the output, or if you want to save it to disk. The `--load` and `--push`
options allow Buildx to infer the exporter settings to use.

For example, if you use the `--push` option in combination with `--tag`, Buildx
automatically uses the `image` exporter, and configures the exporter to push the
results to the specified registry.

To get the full flexibility out of the various exporters BuildKit has to offer,
you use the `--output` flag that lets you configure exporter options.

## Use cases

Each exporter type is designed for different use cases. The following sections
describe some common scenarios, and how you can use exporters to generate the
output that you need.

### Load to image store

Buildx is often used to build container images that can be loaded to an image
store. That's where the `docker` exporter comes in. The following example shows
how to build an image using the `docker` exporter, and have that image loaded to
the local image store, using the `--output` option:

```console
$ docker buildx build \
  --output type=docker,name=<registry>/<image> .
```

Buildx CLI will automatically use the `docker` exporter and load it to the image
store if you supply the `--tag` and `--load` options:

```console
$ docker buildx build --tag <registry>/<image> --load .
```

Building images using the `docker` driver are automatically loaded to the local
image store.

Images loaded to the image store are available to `docker run` immediately
after the build finishes, and you'll see them in the list of images when you run
the `docker images` command.

### Push to registry

To push a built image to a container registry, you can use the `registry` or
`image` exporters.

When you pass the `--push` option to the Buildx CLI, you instruct BuildKit to
push the built image to the specified registry:

```console
$ docker buildx build --tag <registry>/<image> --push .
```

Under the hood, this uses the `image` exporter, and sets the `push` parameter.
It's the same as using the following long-form command using the `--output`
option:

```console
$ docker buildx build \
  --output type=image,name=<registry>/<image>,push=true .
```

You can also use the `registry` exporter, which does the same thing:

```console
$ docker buildx build \
  --output type=registry,name=<registry>/<image> .
```

### Export image layout to file

You can use either the `oci` or `docker` exporters to save the build results to
image layout on your local filesystem. Both of these exporters generate a tar
archive file containing the corresponding image layout. The `dest` parameter
defines the target output path for the tarball.

```console
$ docker buildx build --output type=oci,dest=./image.tar .
[+] Building 0.8s (7/7) FINISHED
 ...
 => exporting to oci image format                                                                     0.0s
 => exporting layers                                                                                  0.0s
 => exporting manifest sha256:c1ef01a0a0ef94a7064d5cbce408075730410060e253ff8525d1e5f7e27bc900        0.0s
 => exporting config sha256:eadab326c1866dd247efb52cb715ba742bd0f05b6a205439f107cf91b3abc853          0.0s
 => sending tarball                                                                                   0.0s
$ mkdir -p out && tar -C out -xf ./image.tar
$ tree out
out
├── blobs
│   └── sha256
│       ├── 9b18e9b68314027565b90ff6189d65942c0f7986da80df008b8431276885218e
│       ├── c78795f3c329dbbbfb14d0d32288dea25c3cd12f31bd0213be694332a70c7f13
│       ├── d1cf38078fa218d15715e2afcf71588ee482352d697532cf316626164699a0e2
│       ├── e84fa1df52d2abdfac52165755d5d1c7621d74eda8e12881f6b0d38a36e01775
│       └── fe9e23793a27fe30374308988283d40047628c73f91f577432a0d05ab0160de7
├── index.json
├── manifest.json
└── oci-layout
```

### Export filesystem

If you don't want to build an image from your build results, but instead export
the filesystem that was built, you can use the `local` and `tar` exporters.

The `local` exporter unpacks the filesystem into a directory structure in the
specified location. The `tar` exporter creates a tarball archive file.

```console
$ docker buildx build --output type=local,dest=<path/to/output> .
```

The `local` exporter is useful in [multi-stage builds](../building/multi-stage.md)
since it allows you to export only a minimal number of build artifacts, such as
self-contained binaries.

### Cache-only export

The `cacheonly` exporter can be used if you just want to run a build, without
exporting any output. This can be useful if, for example, you want to run a test
build. Or, if you want to run the build first, and create exports using
subsequent commands. The `cacheonly` exporter creates a build cache, so any
successive builds are instant.

```console
$ docker buildx build --output type=cacheonly
```

If you don't specify an exporter, and you don't provide short-hand options like
`--load` that automatically selects the appropriate exporter, Buildx defaults to
using the `cacheonly` exporter. Except if you build using the `docker` driver,
in which case you use the `docker` exporter.

Buildx logs a warning message when using `cacheonly` as a default:

```console
$ docker buildx build .
WARNING: No output specified with docker-container driver.
         Build result will only remain in the build cache.
         To push result image into registry use --push or
         to load image into docker use --load
```

## Multiple exporters

{{< summary-bar feature_name="Build multiple exporters" >}}

You can use multiple exporters for any given build by specifying the `--output`
flag multiple times. This requires **both Buildx and BuildKit** version 0.13.0
or later.

The following example runs a single build, using three
different exporters:

- The `registry` exporter to push the image to a registry
- The `local` exporter to extract the build results to the local filesystem
- The `--load` flag (a shorthand for the `image` exporter) to load the results to the local image store.

```console
$ docker buildx build \
  --output type=registry,tag=<registry>/<image> \
  --output type=local,dest=<path/to/output> \
  --load .
```

## Configuration options

This section describes some configuration options available for exporters.

The options described here are common for at least two or more exporter types.
Additionally, the different exporters types support specific parameters as well.
See the detailed page about each exporter for more information about which
configuration parameters apply.

The common parameters described here are:

- [Compression](#compression)
- [OCI media type](#oci-media-types)

### Compression

When you export a compressed output, you can configure the exact compression
algorithm and level to use. While the default values provide a good
out-of-the-box experience, you can tweak the parameters to optimize for
storage versus compute costs. Changing the compression parameters can reduce storage
space required, and improve image download times, but will increase build times.

To select the compression algorithm, you can use the `compression` option. For
example, to build an `image` with `compression=zstd`:

```console
$ docker buildx build \
  --output type=image,name=<registry>/<image>,push=true,compression=zstd .
```

Use the `compression-level=<value>` option alongside the `compression` parameter
to choose a compression level for the algorithms which support it:

- 0-9 for `gzip` and `estargz`
- 0-22 for `zstd`

As a general rule, the higher the number, the smaller the resulting file will
be, and the longer the compression will take to run.

Use the `force-compression=true` option to force re-compressing layers imported
from a previous image, if the requested compression algorithm is different from
the previous compression algorithm.

> [!NOTE]
>
> The `gzip` and `estargz` compression methods use the [`compress/gzip` package](https://pkg.go.dev/compress/gzip),
> while `zstd` uses the [`github.com/klauspost/compress/zstd` package](https://github.com/klauspost/compress/tree/master/zstd).

### OCI media types

The `image`, `registry`, `oci` and `docker` exporters create container images.
These exporters support both Docker media types (default) and OCI media types

To export images with OCI media types set, use the `oci-mediatypes` property.

```console
$ docker buildx build \
  --output type=image,name=<registry>/<image>,push=true,oci-mediatypes=true .
```

## What's next

Read about each of the exporters to learn about how they work and how to use
them:

- [Image and registry exporters](image-registry.md)
- [OCI and Docker exporters](oci-docker.md).
- [Local and tar exporters](local-tar.md)
