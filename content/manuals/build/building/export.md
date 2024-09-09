---
title: Export binaries
weight: 50
description: Using Docker builds to create and export executable binaries
keywords: build, buildkit, buildx, guide, tutorial, build arguments, arg
aliases:
  - /build/guide/export/
---

Did you know that you can use Docker to build your application to standalone
binaries? Sometimes, you don’t want to package and distribute your application
as a Docker image. Use Docker to build your application, and use exporters to
save the output to disk.

The default output format for `docker build` is a container image. That image is
automatically loaded to your local image store, where you can run a container
from that image, or push it to a registry. Under the hood, this uses the default
exporter, called the `docker` exporter.

To export your build results as files instead, you can use the `--output` flag,
or `-o` for short. the `--output` flag lets you change the output format of
your build.

## Export binaries from a build

If you specify a filepath to the `docker build --output` flag, Docker exports
the contents of the build container at the end of the build to the specified
location on your host's filesystem. This uses the `local`
[exporter](/manuals/build/exporters/local-tar.md).

The neat thing about this is that you can use Docker's powerful isolation and
build features to create standalone binaries. This
works well for Go, Rust, and other languages that can compile to a single
binary.

The following example creates a simple Rust program that prints "Hello,
World!", and exports the binary to the host filesystem.

1. Create a new directory for this example, and navigate to it:

   ```console
   $ mkdir hello-world-bin
   $ cd hello-world-bin
   ```

2. Create a Dockerfile with the following contents:

   ```Dockerfile
   # syntax=docker/dockerfile:1
   FROM rust:alpine AS build
   WORKDIR /src
   COPY <<EOT hello.rs
   fn main() {
       println!("Hello World!");
   }
   EOT
   RUN rustc -o /bin/hello hello.rs
   
   FROM scratch
   COPY --from=build /bin/hello /
   ENTRYPOINT ["/hello"]
   ```

   > [!TIP]
   > The `COPY <<EOT` syntax is a [here-document](/reference/dockerfile.md#here-documents).
   > It lets you write multi-line strings in a Dockerfile. Here it's used to
   > create a simple Rust program inline in the Dockerfile.

   This Dockerfile uses a multi-stage build to compile the program in the first
   stage, and then copies the binary to a scratch image in the second. The
   final image is a minimal image that only contains the binary. This use case
   for the `scratch` image is common for creating minimal build artifacts for
   programs that don't require a full operating system to run.

3. Build the Dockerfile and export the binary to the current working directory:

   ```console
   $ docker build --output=. .
   ```

   This command builds the Dockerfile and exports the binary to the current
   working directory. The binary is named `hello`, and it's created in the
   current working directory.

## Exporting multi-platform builds

You use the `local` exporter to export binaries in combination with
[multi-platform builds](/manuals/build/building/multi-platform.md). This lets you
compile multiple binaries at once, that can be run on any machine of any
architecture, provided that the target platform is supported by the compiler
you use.

Continuing on the example Dockerfile in the
[Export binaries from a build](#export-binaries-from-a-build) section:

```dockerfile
# syntax=docker/dockerfile:1
FROM rust:alpine AS build
WORKDIR /src
COPY <<EOT hello.rs
fn main() {
    println!("Hello World!");
}
EOT
RUN rustc -o /bin/hello hello.rs

FROM scratch
COPY --from=build /bin/hello /
ENTRYPOINT ["/hello"]
```

You can build this Rust program for multiple platforms using the `--platform`
flag with the `docker build` command. In combination with the `--output` flag,
the build exports the binaries for each target to the specified directory.

For example, to build the program for both `linux/amd64` and `linux/arm64`:

```console
$ docker build --platform=linux/amd64,linux/arm64 --output=out .
$ tree out/
out/
├── linux_amd64
│   └── hello
└── linux_arm64
    └── hello

3 directories, 2 files
```

## Additional information

In addition to the `local` exporter, there are other exporters available. To
learn more about the available exporters and how to use them, see the
[exporters](/manuals/build/exporters/_index.md) documentation.
