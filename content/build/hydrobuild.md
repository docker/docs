---
title: Hydrobuild
description: Get started with Docker Hydrobuild
sitemap: false
---

> **Early Access**
>
> Docker Hydrobuild is an early-access service that provides cloud-based
> builders for your Docker organization.
>
> If you want to get involved in testing Hydrobuild, you can
> [sign up for the early access program](https://www.docker.com/build-early-access-program/?utm_source=docs).
{ .restricted }

Hydrobuild is a service that lets you build your container images faster, both
locally and in CI. Builds run on cloud infrastructure optimally dimensioned for
your workloads, no configuration required. The service uses a remote build
cache, ensuring fast builds anywhere and for all team members.

## How Hydrobuild works

Using Hydrobuild is no different from running a regular build. You invoke a
build the same way you normally would, using `docker buildx build`. The
difference is in where and how that build gets executed.

By default when you invoke a build command, the build runs on a local instance
of BuildKit, bundled with the Docker daemon. With Hydrobuild, you send the
build request to a BuildKit instance running remotely, in the cloud.

The remote builder executes the build steps, and sends the resulting build
output to the destination that you specify. For example, back to your local
Docker Engine image store, or to an image registry.

Hydrobuild provides several benefits over local builds:

- Improved build speed
- Shared build cache
- Native multi-platform builds

And the best part: you don't need to worry about managing builders or
infrastructure. Just connect to your builders, and start building.

## Setup

To get started with Hydrobuild, you need to:

- Download and install a version of Buildx that supports Hydrobuild.

- Have a Docker ID that's part of a Docker organization participating in the
  [Hydrobuild early access program](https://www.docker.com/build-early-access-program/?utm_source=docs).

Docker Desktop 4.22.0 and later versions ship with a Hydrobuild-compatible
Buildx binary. Alternatively, you can download and install the binary manually
from [this repository](https://github.com/docker/buildx-desktop).

## Connecting to Hydrobuild

To start building with Hydrobuild, you must create a new builder using the
`docker buildx create` command. The builder is connected to Hydrobuild through
an endpoint that you specify. The endpoint represents a single, isolated
builder. Builder endpoints have the following format:

```text
cloud://<org>/<group>_<platform>
```

- `<org>` is the Docker organization that the builder is provisioned for
- `<group>` is the builder node group
- `<platform>` is the native OS and architecture of the builder

An organization can contain multiple builder groups. Each builder group is an
isolated builder, and by creating multiple groups you can assign unique
builders to different teams in the organization.

The platform suffix is optional, and if omitted creates a `linux/amd64` builder
by default. The supported values for `<platform>` are:

- `linux-amd64`
- `linux-arm64`

You can use the platform suffix to create a multi-node builder group with
multiple builders of different native architectures. This gives you a
high-performance build cluster for building multi-platform images natively. See
[Create a multi-platform builder](#create-a-multi-platform-builder).

You can omit the `cloud://` protocol prefix from the endpoint when you create a
builder using the `cloud` driver. The endpoint format then becomes
`<org>/<group>_<platform>`. The `docker buildx ls` command shows the full
endpoint URI, including the prefix.

### Create a single-platform builder

To create a `linux/amd64` builder:

1. Sign in to your Docker ID using the Docker Desktop UI or the `docker login`
   command.

2. Create a builder that uses the `cloud` driver.

   ```console
   $ docker buildx create --driver cloud --name hydrobuild \
     --platform linux/amd64 \
     <org>/<group>_linux-amd64
   ```

   Replace `<org>` with the Docker organization, and `<group>` with the name
   that you want to use for this builder group.

### Create a multi-platform builder

To create a builder with support for native `linux/amd64` and `linux/arm64`
builds:

1. Sign in to your Docker ID using the Docker Desktop UI or the `docker login`
   command.

2. Create a `linux/amd64` builder that uses the `cloud` driver.

   ```console
   $ docker buildx create --driver cloud --name hydrobuild \
     --platform linux/amd64 \
     <org>/<group>_linux-amd64
   ```

   Replace `<org>` with the Docker organization, and `<group>` with the name
   that you want to use for this builder group.

3. Create a `linux/arm64` builder and append it to the `hydrobuild` builder you
   just created.

   ```console
   $ docker buildx create --append --name hydrobuild \
     --platform linux/arm64 \
     <org>/<group>_linux-arm64
   ```

   `<org>` and `<group>` should be the same as for first builder, but this time
   use `linux-arm64` for the platform suffix.

## Use Hydrobuild

To build your applications with Hydrobuild, you can:

- [Use the Docker CLI](#cli) to build from your local development machine
- [Use GitHub Actions](#github-actions) to build with Hydrobuild in CI

### CLI

To run a build using Hydrobuild, invoke a build command and specify the
name of the builder using the `--builder` flag.

```console
$ docker buildx build --builder hydrobuild --tag myorg/some-tag .
```

> **Note**
>
> Specifying `--tag` ensures that the build result gets exported to your local
> image store when the build finishes. If you want to download the results from
> Hydrobuild without specifying a tag, you must pass the `--load` flag.
>
> If you use the containerd image store, you must always pass `--load` to
> download the results, even if you build with a tag.

If you created a [multi-platform builder](#create-a-multi-platform-builder),
you can build multi-platform images using the `--platform` flag:

```console
$ docker buildx build --builder hydrobuild \
  --platform linux/amd64,linux/arm64 \
  --tag myorg/some-tag --push .
```

> **Note**
>
> If you build multi-platform images, you won't be able to load the images back
> to your local image store unless you turn on the containerd image store
> feature, and use the `--load` flag.
>
> Using the containerd image store with Hydrobuild currently results in slower
> transfers of build output to the client, compared to when you use the default
> image store.
>
> When building multi-platform images, consider pushing the resulting image to
> a registry directly, using the `docker buildx build --push` flag.

### Use by default

If you want to use Hydrobuild by default, you can run the following command to
make it the selected builder:

```console
$ docker buildx use hydrobuild --global
```

> **Note**
>
> Changing your default builder with `docker buildx use` only changes the
> default builder for the `docker buildx build` command. The `docker build`
> command still uses the `default` builder, unless you specify the `--builder`
> flag explicitly.
>
> If you use build scripts, such as `make`, we recommend that you update your
> build commands from `docker build` to `docker buildx build`, to avoid any
> confusion with regards to builder selection. Alternatively, you can run
> `docker buildx install` to make the default `docker build` command behave
> like `docker buildx build`, without discrepancies.

### GitHub Actions

You can use GitHub Actions in combination with Hydrobuild to achieve faster
build times, while still leveraging the convenience of GitHub Action workflows.

With this approach, your CI workflows run on a GitHub Actions runner, and the
runner calls out to the builder to build the image.

To use Hydrobuild with GitHub Actions, you must first sign in with your Docker
ID, and then use the `lab` channel of `setup-buildx-action`:

```yaml
- name: Log in to Docker Hub
  uses: docker/login-action@v2
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}
- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v2
  with:
    version: "lab:latest"
    driver: cloud
    endpoint: "<org>/<group>"
```

The following example shows a basic workflow for GitHub Actions with Hydrobuild.

```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      - name: Log in to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          version: "lab:latest"
          driver: cloud
          endpoint: "<org>/<group>"
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: user/app:latest
```

This invokes the build from a GitHub Actions workflow, runs the build on
Hydrobuild, and pushes the image to a Docker Hub registry.

> **Note**
>
> The previous example uses a `push: true` configuration for the _Build and
> push_ GitHub Action. This ensures that the build result is pushed to a
> registry directly, rather than being loaded back to the image store of the
> GitHub Actions runner. When using Hydrobuild in CI, this is the recommended
> workflow, because it speeds up your builds and avoids unnecessary file
> transfers.
>
> If you're not using `push: true`, and if you build an image with a `tag`,
> Hydrobuild automatically loads the build results back to the client. If you
> only want to build the artifact without loading the results (as a validation
> step in pull requests, for example), you can add `outputs: type=cacheonly` to
> the action configuration:
>
> ```yaml
> - name: Build and push
>   uses: docker/build-push-action@v4
>   with:
>     context: .
>     tags: user/app:latest
>     # if this runs in a pull request, export results to build cache
>     outputs: ${{ github.event_name == 'pull_request' && 'type=cacheonly' || '' }}
>     # if this doesn't run in a pull request, push to a registry
>     push: ${{ github.event_name != 'pull_request' }}
> ```

## Hydrobuild in Docker Desktop

The Docker Desktop [Builds view](../desktop/use-desktop/builds.md) works with
Hydrobuild out of the box. With Hydrobuild, the Builds view becomes a
collaboration tool, showing information about not only your own builds, but
also builds initiated by your team members using the same builder.

Teams using a shared builder get access to information such as:

- Ongoing and completed builds
- Build configuration, statistics, dependencies, and results
- Build source (Dockerfile)
- Build logs and errors

This lets you and your team can work collaboratively on troubleshooting and
improving build speeds, without having to send build logs and benchmarks back
and forth between each other.

## Optimize for building in the cloud

Hydrobuild runs your builds remotely, and not on the machine where you invoke
the build. This means that file transfers between the client and builder happens
over the network.

Transferring files over the network has a higher latency and lower bandwidth
than local transfers. Hydrobuild has several features to mitigate this:

- It uses attached storage volumes for build cache, which makes reading and
  writing cache very fast.
- Loading build results back to the client only pulls the layers that were
  changed compared to previous builds.

Despite these optimizations, building remotely can still yield slow context
transfers and image loads, for large projects or if the network connection is
slow. Here are some ways that you can optimize your builds to make the transfer
more efficient:

- [Dockerignore files](#dockerignore-files)
- [Slim base images](#slim-base-images)
- [Multi-stage builds](#multi-stage-builds)
- [Fetch remote files in build](#fetch-remote-files-in-build)
- [Multi-threaded tools](#multi-threaded-tools)

### Dockerignore files

Using a [`.dockerignore` file](./building/context.md#dockerignore), you can be
explicit about which local files that you don’t want to include in the build
context. Files caught by the [glob patterns](../engine/reference/builder.md#dockerignore-file)
you specify in your ignore-file are not transferred to the remote builder.

Some examples of things you might want to add to your `.dockerignore` file are:

- `.git` — skip sending the version control history in the build context. Note
  that this means you won’t be able to run Git commands in your build steps,
  such as `git rev-parse` etc.
- Directories containing build artifacts, such as binaries. Build artifacts
  created locally during development.
- Vendor directories for package managers, such as `node_modules`.

In general, the contents of your `.dockerignore` file should be similar to what
you have in your `.gitignore`.

### Slim base images

Selecting smaller images for your `FROM` instructions in your Dockerfile can
help reduce the size of the final image. The [Alpine image](https://hub.docker.com/_/alpine)
is a good example of a minimal Docker image that provides all of the OS
utilities you would expect from a Linux container.

There’s also the [special `scratch` image](https://hub.docker.com/_/scratch),
which contains nothing at all. Useful for creating images of statically linked
binaries, for example.

### Multi-stage builds

[Multi-stage builds](./guide/multi-stage.md) can make your build run faster,
because stages can run in parallel. It can also make your end-result smaller.
Write your Dockerfile in such a way that the final runtime stage uses the
smallest possible base image, with only the resources that your program requires
to run.

It’s also possible to
[copy resources from other images or stages](./building/multi-stage.md#name-your-build-stages),
using the Dockerfile `COPY --from` instruction. This technique can reduce the
number of layers, and the size of those layers, in the final stage.

### Fetch remote files in build

When possible, you should fetch files from a remote location in the build,
rather than bundling the files into the build context. Downloading files on the
Hydrobuild server directly is better, because it will likely be faster than
transferring the files with the build context.

You can fetch remote files during the build using the
[Dockerfile `ADD` instruction](../engine/reference/builder.md#add),
or in your `RUN` instructions with tools like `wget` and `rsync`.

### Multi-threaded tools

Some tools that you use in your build instructions may not utilize multiple
cores by default. One such example is `make` which uses a single thread by
default, unless you specify the `make --jobs=<n>` option. For build steps
involving such tools, try checking if you can optimize the execution with
parallelization.

## Frequently asked questions

### How do I remove Hydrobuild from my system?

If you want to stop using Hydrobuild, and remove it from your system, remove
the builder using the `docker buildx rm` command.

```console
$ docker buildx rm hydrobuild
```

This doesn't deprovision the builder backend, it only removes the builder from
your local Docker client.

### Are builders shared between organizations?

No. Each Hydrobuild builder provisioned to an organization is completely
isolated to a single Amazon EC2 instance, with a dedicated EBS volume for build
cache, and end-to-end encryption. That means there are no shared processes or
data between Hydrobuild jobs.

### Do I need to add my secrets the builder to access private resources?

No. Your interface to Hydrobuild is Buildx, and you can use the existing
`--secret` and `--ssh` CLI flags for managing build secrets.

For more information, refer to:

- [docker buildx build --secret](../engine/reference/commandline/buildx_build.md#secret)
- [docker buildx build --ssh](../engine/reference/commandline/buildx_build.md#ssh)
