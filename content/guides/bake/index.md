---
title: Mastering multi-platform builds, testing, and more with Docker Buildx Bake
linkTitle: Mastering Docker Buildx Bake
description: >
  Learn how to manage simple and complex build configurations with Buildx Bake.
summary: >
  Learn to automate Docker builds and testing with declarative configurations using Buildx Bake.
tags: [devops]
languages: [go]
params:
  time: 30 minutes
  featured: true
  image: /images/guides/bake.webp
---

This guide demonstrates how to simplify and automate the process of building
images, testing, and generating build artifacts using Docker Buildx Bake. By
defining build configurations in a declarative `docker-bake.hcl` file, you can
eliminate manual scripts and enable efficient workflows for complex builds,
testing, and artifact generation.

## Assumptions

This guide assumes that you're familiar with:

- Docker
- [Buildx](/manuals/build/concepts/overview.md#buildx)
- [BuildKit](/manuals/build/concepts/overview.md#buildkit)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
- [Multi-platform builds](/manuals/build/building/multi-platform.md)

## Prerequisites

- You have a recent version of Docker installed on your machine.
- You have Git installed for cloning repositories.
- You're using the [containerd](/manuals/desktop/features/containerd.md) image store.

## Introduction

This guide uses an example project to demonstrate how Docker Buildx Bake can
streamline your build and test workflows. The repository includes both a
Dockerfile and a `docker-bake.hcl` file, giving you a ready-to-use setup to try
out Bake commands.

Start by cloning the example repository:

```bash
git clone https://github.com/dvdksn/bakeme.git
cd bakeme
```

The Bake file, `docker-bake.hcl`, defines the build targets in a declarative
syntax, using targets and groups, allowing you to manage complex builds
efficiently.

Here's what the Bake file looks like out-of-the-box:

```hcl
target "default" {
  target = "image"
  tags = [
    "bakeme:latest",
  ]
  attest = [
    "type=provenance,mode=max",
    "type=sbom",
  ]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/riscv64",
  ]
}
```

The `target` keyword defines a build target for Bake. The `default` target
defines the target to build when no specific target is specified on the command
line. Here's a quick summary of the options for the `default` target:

- `target`: The target build stage in the Dockerfile.
- `tags`: Tags to assign to the image.
- `attest`: [Attestations](/manuals/build/metadata/attestations/_index.md) to attach to the image.

  > [!TIP]
  > The attestations provide metadata such as build provenance, which tracks
  > the source of the image's build, and an SBOM (Software Bill of Materials),
  > useful for security audits and compliance.

- `platforms`: Platform variants to build.

To execute this build, simply run the following command in the root of the
repository:

```console
$ docker buildx bake
```

With Bake, you avoid long, hard-to-remember command-line incantations,
simplifying build configuration management by replacing manual, error-prone
scripts with a structured configuration file.

For contrast, here's what this build command would look like without Bake:

```console
$ docker buildx build \
  --target=image \
  --tag=bakeme:latest \
  --provenance=true \
  --sbom=true \
  --platform=linux/amd64,linux/arm64,linux/riscv64 \
  .
```

## Testing and linting

Bake isn't just for defining build configurations and running builds. You can
also use Bake to run your tests, effectively using BuildKit as a task runner.
Running your tests in containers is great for ensuring reproducible results.
This section shows how to add two types of tests:

- Unit testing with `go test`.
- Linting for style violations with `golangci-lint`.

In Test-Driven Development (TDD) fashion, start by adding a new `test` target
to the Bake file:

```hcl
target "test" {
  target = "test"
  output = ["type=cacheonly"]
}
```

> [!TIP]
> Using `type=cacheonly` ensures that the build output is effectively
> discarded; the layers are saved to BuildKit's cache, but Buildx will not
> attempt to load the result to the Docker Engine's image store.
>
> For test runs, you don't need to export the build output — only the test
> execution matters.

To execute this Bake target, run `docker buildx bake test`. At this time,
you'll receive an error indicating that the `test` stage does not exist in the
Dockerfile.

```console
$ docker buildx bake test
[+] Building 1.2s (6/6) FINISHED
 => [internal] load local bake definitions
...
ERROR: failed to solve: target stage "test" could not be found
```

To satisfy this target, add the corresponding Dockerfile target. The `test`
stage here is based on the same base stage as the build stage.

```dockerfile
FROM base AS test
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    go test .
```

> [!TIP]
> The [`--mount=type=cache` directive](/manuals/build/cache/optimize.md#use-cache-mounts)
> caches Go modules between builds, improving build performance by avoiding the
> need to re-download dependencies. This shared cache ensures that the same
> dependency set is available across build, test, and other stages.

Now, running the `test` target with Bake will evaluate the unit tests for this
project. If you want to verify that it works, you can make an arbitrary change
to `main_test.go` to cause the test to fail.

Next, to enable linting, add another target to the Bake file, named `lint`:

```hcl
target "lint" {
  target = "lint"
  output = ["type=cacheonly"]
}
```

And in the Dockerfile, add the build stage. This stage will use the official
`golangci-lint` image on Docker Hub.

> [!TIP]
> Because this stage relies on executing an external dependency, it's generally
> a good idea to define the version you want to use as a build argument. This
> lets you more easily manage version upgrades in the future by collocating
> dependency versions to the beginning of the Dockerfile.

```dockerfile {hl_lines=[2,"6-8"]}
ARG GO_VERSION="1.23"
ARG GOLANGCI_LINT_VERSION="1.61"

#...

FROM golangci/golangci-lint:v${GOLANGCI_LINT_VERSION}-alpine AS lint
RUN --mount=target=.,rw \
    golangci-lint run
```

Lastly, to enable running both tests simultaneously, you can use the `groups`
construct in the Bake file. A group can specify multiple targets to run with a
single invocation.

```hcl
group "validate" {
  targets = ["test", "lint"]
}
```

Now, running both tests is as simple as:

```console
$ docker buildx bake validate
```

## Building variants

Sometimes you need to build more than one version of a program. The following
example uses Bake to build separate "release" and "debug" variants of the
program, using [matrices](/manuals/build/bake/matrices.md). Using matrices lets
you run parallel builds with different configurations, saving time and ensuring
consistency.

A matrix expands a single build into multiple builds, each representing a
unique combination of matrix parameters. This means you can orchestrate Bake
into building both the production and development build of your program in
parallel, with minimal configuration changes.

The example project for this guide is set up to use a build-time option to
conditionally enable debug logging and tracing capabilities.

- If you compile the program with `go build -tags="debug"`, the additional
  logging and tracing capabilities are enabled (development mode).
- If you build without the `debug` tag, the program is compiled with a default
  logger (production mode).

Update the Bake file by adding a matrix attribute which defines the variable
combinations to build:

```diff {title="docker-bake.hcl"}
 target "default" {
+  matrix = {
+    mode = ["release", "debug"]
+  }
+  name = "image-${mode}"
   target = "image"
```

The `matrix` attribute defines the variants to build ("release" and "debug").
The `name` attribute defines how the matrix gets expanded into multiple
distinct build targets. In this case, the matrix attribute expands the build
into two workflows: `image-release` and `image-debug`, each using different
configuration parameters.

Next, define a build argument named `BUILD_TAGS` which takes the value of the
matrix variable.

```diff {title="docker-bake.hcl"}
   target = "image"
+  args = {
+    BUILD_TAGS = mode
+  }
   tags = [
```

You'll also want to change how the image tags are assigned to these builds.
Currently, both matrix paths would generate the same image tag names, and
overwrite each other. Update the `tags` attribute use a conditional operator to
set the tag depending on the matrix variable value.

```diff {title="docker-bake.hcl"}
   tags = [
-    "bakeme:latest",
+    mode == "release" ? "bakeme:latest" : "bakeme:dev"
   ]
```

- If `mode` is `release`, the tag name is `bakeme:latest`
- If `mode` is `debug`, the tag name is `bakeme:dev`

Finally, update the Dockerfile to consume the `BUILD_TAGS` argument during the
compilation stage. When the `-tags="${BUILD_TAGS}"` option evaluates to
`-tags="debug"`, the compiler uses the `configureLogging` function in the
[`debug.go`](https://github.com/dvdksn/bakeme/blob/75c8a41e613829293c4bd3fc3b4f0c573f458f42/debug.go#L1)
file.

```diff {title=Dockerfile}
 # build compiles the program
 FROM base AS build
-ARG TARGETOS TARGETARCH
+ARG TARGETOS TARGETARCH BUILD_TAGS
 ENV GOOS=$TARGETOS
 ENV GOARCH=$TARGETARCH
 RUN --mount=target=. \
        --mount=type=cache,target=/go/pkg/mod \
-       go build -o "/usr/bin/bakeme" .
+       go build -tags="${BUILD_TAGS}" -o "/usr/bin/bakeme" .
```

That's all. With these changes, your `docker buildx bake` command now builds
two multi-platform image variants. You can introspect the canonical build
configuration that Bake generates using the `docker buildx bake --print`
command. Running this command shows that Bake will run a `default` group with
two targets with different build arguments and image tags.

```json {collapse=true}
{
  "group": {
    "default": {
      "targets": ["image-release", "image-debug"]
    }
  },
  "target": {
    "image-debug": {
      "attest": ["type=provenance,mode=max", "type=sbom"],
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "BUILD_TAGS": "debug"
      },
      "tags": ["bakeme:dev"],
      "target": "image",
      "platforms": ["linux/amd64", "linux/arm64", "linux/riscv64"]
    },
    "image-release": {
      "attest": ["type=provenance,mode=max", "type=sbom"],
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "BUILD_TAGS": "release"
      },
      "tags": ["bakeme:latest"],
      "target": "image",
      "platforms": ["linux/amd64", "linux/arm64", "linux/riscv64"]
    }
  }
}
```

Factoring in all of the platform variants as well, this means that the build
configuration generates 6 different images.

```console
$ docker buildx bake
$ docker image ls --tree

IMAGE                   ID             DISK USAGE   CONTENT SIZE   USED
bakeme:dev              f7cb5c08beac       49.3MB         28.9MB
├─ linux/riscv64        0eae8ba0367a       9.18MB         9.18MB
├─ linux/arm64          56561051c49a         30MB         9.89MB
└─ linux/amd64          e8ca65079c1f        9.8MB          9.8MB

bakeme:latest           20065d2c4d22       44.4MB         25.9MB
├─ linux/riscv64        7cc82872695f       8.21MB         8.21MB
├─ linux/arm64          e42220c2b7a3       27.1MB         8.93MB
└─ linux/amd64          af5b2dd64fde       8.78MB         8.78MB
```

## Exporting build artifacts

Exporting build artifacts like binaries can be useful for deploying to
environments without Docker or Kubernetes. For example, if your programs are
meant to be run on a user's local machine.

> [!TIP]
> The techniques discussed in this section can be applied not only to build
> output like binaries, but to any type of artifacts, such as test reports.

With programming languages like Go and Rust where the compiled binaries are
usually portable, creating alternate build targets for exporting only the
binary is trivial. All you need to do is add an empty stage in the Dockerfile
containing nothing but the binary that you want to export.

First, let's add a quick way to build a binary for your local platform and
export it to `./build/local` on the local filesystem.

In the `docker-bake.hcl` file, create a new `bin` target. In this stage, set
the `output` attribute to a local filesystem path. Buildx automatically detects
that the output looks like a filepath, and exports the results to the specified
path using the [local exporter](/manuals/build/exporters/local-tar.md).

```hcl
target "bin" {
  target = "bin"
  output = ["build/bin"]
  platforms = ["local"]
}
```

Notice that this stage specifies a `local` platform. By default, if `platforms`
is unspecified, builds target the OS and architecture of the BuildKit host. If
you're using Docker Desktop, this often means builds target `linux/amd64` or
`linux/arm64`, even if your local machine is macOS or Windows, because Docker
runs in a Linux VM. Using the `local` platform forces the target platform to
match your local environment.

Next, add the `bin` stage to the Dockerfile which copies the compiled binary
from the build stage.

```dockerfile
FROM scratch AS bin
COPY --from=build "/usr/bin/bakeme" /
```

Now you can export your local platform version of the binary with `docker
buildx bake bin`. For example, on macOS, this build target generates an
executable in the [Mach-O format](https://en.wikipedia.org/wiki/Mach-O) — the
standard executable format for macOS.

```console
$ docker buildx bake bin
$ file ./build/bin/bakeme
./build/bin/bakeme: Mach-O 64-bit executable arm64
```

Next, let's add a target to build all of the platform variants of the program.
To do this, you can [inherit](/manuals/build/bake/inheritance.md) the `bin`
target that you just created, and extend it by adding the desired platforms.

```hcl
target "bin-cross" {
  inherits = ["bin"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/riscv64",
  ]
}
```

Now, building the `bin-cross` target creates binaries for all platforms.
Subdirectories are automatically created for each variant.

```console
$ docker buildx bake bin-cross
$ tree build/
build/
└── bin
    ├── bakeme
    ├── linux_amd64
    │   └── bakeme
    ├── linux_arm64
    │   └── bakeme
    └── linux_riscv64
        └── bakeme

5 directories, 4 files
```

To also generate "release" and "debug" variants, you can use a matrix just like
you did with the default target. When using a matrix, you also need to
differentiate the output directory based on the matrix value, otherwise the
binary gets written to the same location for each matrix run.

```hcl
target "bin-all" {
  inherits = ["bin-cross"]
  matrix = {
    mode = ["release", "debug"]
  }
  name = "bin-${mode}"
  args = {
    BUILD_TAGS = mode
  }
  output = ["build/bin/${mode}"]
}
```

```console
$ rm -r ./build/
$ docker buildx bake bin-all
$ tree build/
build/
└── bin
    ├── debug
    │   ├── linux_amd64
    │   │   └── bakeme
    │   ├── linux_arm64
    │   │   └── bakeme
    │   └── linux_riscv64
    │       └── bakeme
    └── release
        ├── linux_amd64
        │   └── bakeme
        ├── linux_arm64
        │   └── bakeme
        └── linux_riscv64
            └── bakeme

10 directories, 6 files
```

## Conclusion

Docker Buildx Bake streamlines complex build workflows, enabling efficient
multi-platform builds, testing, and artifact export. By integrating Buildx Bake
into your projects, you can simplify your Docker builds, make your build
configuration portable, and wrangle complex configurations more easily.

Experiment with different configurations and extend your Bake files to suit
your project's needs. You might consider integrating Bake into your CI/CD
pipelines to automate builds, testing, and artifact deployment. The flexibility
and power of Buildx Bake can significantly improve your development and
deployment processes.

### Further reading

For more information about how to use Bake, check out these resources:

- [Bake documentation](/manuals/build/bake/_index.md)
- [Matrix targets](/manuals/build/bake/matrices.md)
- [Bake file reference](/manuals/build/bake/reference.md)
- [Bake GitHub Action](https://github.com/docker/bake-action)
