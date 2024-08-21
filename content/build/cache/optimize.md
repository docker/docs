---
title: Optimize cache usage in builds
description: An overview on how to optimize cache utilization in Docker builds.
keywords: build, buildkit, buildx, guide, tutorial, mounts, cache mounts, bind mounts
aliases:
  - /build/guide/mounts/
---

When building with Docker, a layer is reused from the build cache if the
instruction and the files it depends on hasn't changed since it was previously
built. Reusing layers from the cache speeds up the build process because Docker
doesn't have to rebuild the layer again.

Here are a few techniques you can use to optimize build caching and speed up
the build process:

- [Order your layers](#order-your-layers): Putting the commands in your
  Dockerfile into a logical order can help you avoid unnecessary cache
  invalidation.
- [Keep the context small](#keep-the-context-small): The context is the set of
  files and directories that are sent to the builder to process a build
  instruction. Keeping the context as small reduces the amount of data that
  needs to be sent to the builder, and reduces the likelihood of cache
  invalidation.
- [Use bind mounts](#use-bind-mounts): Bind mounts let you mount a file or
  directory from the host machine into the build container. Using bind mounts
  can help you avoid unnecessary layers in the image, which can slow down the
  build process.
- [Use cache mounts](#use-cache-mounts): Cache mounts let you specify a
  persistent package cache to be used during builds. The persistent cache helps
  speed up build steps, especially steps that involve installing packages using
  a package manager. Having a persistent cache for packages means that even if
  you rebuild a layer, you only download new or changed packages.
- [Use an external cache](#use-an-external-cache): An external cache lets you
  store build cache at a remote location. The external cache image can be
  shared between multiple builds, and across different environments.

## Order your layers

Putting the commands in your Dockerfile into a logical order is a great place
to start. Because a change causes a rebuild for steps that follow, try to make
expensive steps appear near the beginning of the Dockerfile. Steps that change
often should appear near the end of the Dockerfile, to avoid triggering
rebuilds of layers that haven't changed.

Consider the following example. A Dockerfile snippet that runs a JavaScript
build from the source files in the current directory:

```dockerfile
# syntax=docker/dockerfile:1
FROM node
WORKDIR /app
COPY . .          # Copy over all files in the current directory
RUN npm install   # Install dependencies
RUN npm build     # Run build
```

This Dockerfile is rather inefficient. Updating any file causes a reinstall of
all dependencies every time you build the Docker image even if the dependencies
didn't change since last time.

Instead, the `COPY` command can be split in two. First, copy over the package
management files (in this case, `package.json` and `yarn.lock`). Then, install
the dependencies. Finally, copy over the project source code, which is subject
to frequent change.

```dockerfile
# syntax=docker/dockerfile:1
FROM node
WORKDIR /app
COPY package.json yarn.lock .    # Copy package management files
RUN npm install                  # Install dependencies
COPY . .                         # Copy over project files
RUN npm build                    # Run build
```

By installing dependencies in earlier layers of the Dockerfile, there is
no need to rebuild those layers when a project file has changed.

## Keep the context small

The easiest way to make sure your context doesn't include unnecessary files is
to create a `.dockerignore` file in the root of your build context. The
`.dockerignore` file works similarly to `.gitignore` files, and lets you
exclude files and directories from the build context.

Here's an example `.dockerignore` file that excludes the `node_modules`
directory, all files and directories that start with `tmp`:

```plaintext {title=".dockerignore"}
node_modules
tmp*
```

Ignore-rules specified in the `.dockerignore` file apply to the entire build
context, including subdirectories. This means it's a rather coarse-grained
mechanism, but it's a good way to exclude files and directories that you know
you don't need in the build context, such as temporary files, log files, and
build artifacts.

## Use bind mounts

You might be familiar with bind mounts for when you run containers with `docker
run` or Docker Compose. Bind mounts let you mount a file or directory from the
host machine into a container.

```bash
# bind mount using the -v flag
docker run -v $(pwd):/path/in/container image-name
# bind mount using the --mount flag
docker run --mount=type=bind,src=.,dst=/path/in/container image-name
```

To use bind mounts in a build, you can use the `--mount` flag with the `RUN`
instruction in your Dockerfile:

```dockerfile
FROM golang:latest
WORKDIR /app
RUN --mount=type=bind,target=. go build -o /app/hello
```

In this example, the current directory is mounted into the build container
before the `go build` command gets executed. The source code is available in
the build container for the duration of that `RUN` instruction. When the
instruction is done executing, the mounted files are not persisted in the final
image, or in the build cache. Only the output of the `go build` command
remains.

The `COPY` and `ADD` instructions in a Dockerfile lets you copy files from the
build context into the build container. Using bind mounts is beneficial for
build cache optimization because you're not adding unnecessary layers to the
cache. If you have build context that's on the larger side, and it's only used
to generate an artifact, you're better off using bind mounts to temporarily
mount the source code required to generate the artifact into the build. If you
use `COPY` to add the files to the build container, BuildKit will include all
of those files in the cache, even if the files aren't used in the final image.

There are a few things to be aware of when using bind mounts in a build:

- Bind mounts are read-only by default. If you need to write to the mounted
  directory, you need to specify the `rw` option. However, even with the `rw`
  option, the changes are not persisted in the final image or the build cache.
  The file writes are sustained for the duration of the `RUN` instruction, and
  are discarded after the instruction is done.
- Mounted files are not persisted in the final image. Only the output of the
  `RUN` instruction is persisted in the final image. If you need to include
  files from the build context in the final image, you need to use the `COPY`
  or `ADD` instructions.
- If the target directory is not empty, the contents of the target directory
  are hidden by the mounted files. The original contents are restored after the
  `RUN` instruction is done.

  {{< accordion title="Example" >}}

  For example, given a build context with only a `Dockerfile` in it:

  ```plaintext
  .
  └── Dockerfile
  ```

  And a Dockerfile that mounts the current directory into the build container:

  ```dockerfile
  FROM alpine:latest
  WORKDIR /work
  RUN touch foo.txt
  RUN --mount=type=bind,target=. ls
  RUN ls
  ```

  The first `ls` command with the bind mount shows the contents of the mounted
  directory. The second `ls` lists the contents of the original build context.

  ```plaintext {title="Build log"}
  #8 [stage-0 3/5] RUN touch foo.txt
  #8 DONE 0.1s
  
  #9 [stage-0 4/5] RUN --mount=target=. ls -1
  #9 0.040 Dockerfile
  #9 DONE 0.0s
  
  #10 [stage-0 5/5] RUN ls -1
  #10 0.046 foo.txt
  #10 DONE 0.1s
  ```

  {{< /accordion >}}

## Use cache mounts

Regular cache layers in Docker correspond to an exact match of the instruction
and the files it depends on. If the instruction and the files it depends on
have changed since the layer was built, the layer is invalidated, and the build
process has to rebuild the layer.

Cache mounts are a way to specify a persistent cache location to be used during
builds. The cache is cumulative across builds, so you can read and write to the
cache multiple times. This persistent caching means that even if you need to
rebuild a layer, you only download new or changed packages. Any unchanged
packages are reused from the cache mount.

To use cache mounts in a build, you can use the `--mount` flag with the `RUN`
instruction in your Dockerfile:

```dockerfile
FROM node:latest
WORKDIR /app
RUN --mount=type=cache,target=/root/.npm npm install
```

In this example, the `npm install` command uses a cache mount for the
`/root/.npm` directory, the default location for the npm cache. The cache mount
is persisted across builds, so even if you end up rebuilding the layer, you
only download new or changed packages. Any changes to the cache are persisted
across builds, and the cache is shared between multiple builds.

How you specify cache mounts depends on the build tool you're using. If you're
unsure how to specify cache mounts, refer to the documentation for the build
tool you're using. Here are a few examples:

{{< tabs >}}
{{< tab name="Go" >}}

```dockerfile
RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o /app/hello
```

{{< /tab >}}
{{< tab name="Apt" >}}

```dockerfile
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
  --mount=type=cache,target=/var/lib/apt,sharing=locked \
  apt update && apt-get --no-install-recommends install -y gcc
```

{{< /tab >}}
{{< tab name="Python" >}}

```dockerfile
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install -r requirements.txt
```

{{< /tab >}}
{{< tab name="Ruby" >}}

```dockerfile
RUN --mount=type=cache,target=/root/.gem \
    bundle install
```

{{< /tab >}}
{{< tab name="Rust" >}}

```dockerfile
RUN --mount=type=cache,target=/app/target/ \
    --mount=type=cache,target=/usr/local/cargo/git/db \
    --mount=type=cache,target=/usr/local/cargo/registry/ \
    cargo build
```

{{< /tab >}}
{{< tab name=".NET" >}}

```dockerfile
RUN --mount=type=cache,target=/root/.nuget/packages \
    dotnet restore
```

{{< /tab >}}
{{< tab name="PHP" >}}  

```dockerfile
RUN --mount=type=cache,target=/tmp/cache \
    composer install
```

{{< /tab >}}
{{< /tabs >}}

It's important that you read the documentation for the build tool you're using
to make sure you're using the correct cache mount options. Package managers
have different requirements for how they use the cache, and using the wrong
options can lead to unexpected behavior. For example, Apt needs exclusive
access to its data, so the caches use the option `sharing=locked` to ensure
parallel builds using the same cache mount wait for each other and not access
the same cache files at the same time.

## Use an external cache

The default cache storage for builds is internal to the builder (BuildKit
instance) you're using. Each builder uses its own cache storage. When you
switch between different builders, the cache is not shared between them. Using
an external cache lets you define a remote location for pushing and pulling
cache data.

External caches are especially useful for CI/CD pipelines, where the builders
are often ephemeral, and build minutes are precious. Reusing the cache between
builds can drastically speed up the build process and reduce cost. You can even
make use of the same cache in your local development environment.

To use an external cache, you specify the `--cache-to` and `--cache-from`
options with the `docker buildx build` command.

- `--cache-to` exports the build cache to the specified location.
- `--cache-from` specifies remote caches for the build to use.

The following example shows how to set up a GitHub Actions workflow using
`docker/build-push-action`, and push the build cache layers to an OCI registry
image:

```yaml {title=".github/workflows/ci.yml"}
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      
      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          push: true
          tags: user/app:latest
          cache-from: type=registry,ref=user/app:buildcache
          cache-to: type=registry,ref=user/app:buildcache,mode=max
```

This setup tells BuildKit to look for cache in the `user/app:buildcache` image.
And when the build is done, the new build cache is pushed to the same image,
overwriting the old cache.

This cache can be used locally as well. To pull the cache in a local build,
you can use the `--cache-from` option with the `docker buildx build` command:

```console
$ docker buildx build --cache-from type=registry,ref=user/app:buildcache .
```

## Summary

Optimizing cache usage in builds can significantly speed up the build process.
Keeping the build context small, using bind mounts, cache mounts, and external
caches are all techniques you can use to make the most of the build cache and
speed up the build process.

For more information about the concepts discussed in this guide, see:

- [.dockerignore files](/build/concepts/context.md#dockerignore-files)
- [Cache invalidation](/build/cache/invalidation.md)
- [Cache mounts](/reference/dockerfile.md#run---mounttypecache)
- [Cache backend types](/build/cache/backends/_index.md)
- [Building best practices](/build/building/best-practices.md)
