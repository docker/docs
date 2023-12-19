---
title: Building with Docker Build Cloud
description: Invoke your cloud builds with the Buildx CLI client
keywords: build, cloud build, usage, cli, buildx, client
aliases:
  - /hydrobuild/usage/
---

To build using Docker Build Cloud, invoke a build command and specify the name of the
builder using the `--builder` flag.

```console
$ docker buildx build --builder cloud-<ORG>-<BUILDER_NAME> --tag <IMAGE> .
```

## Use by default

If you want to use Docker Build Cloud without having to specify the `--builder` flag
each time, you can set it as the default builder.

{{< tabs group="ui" >}}
{{< tab name="CLI" >}}

Run the following command:

```console
$ docker buildx use cloud-<ORG>-<BUILDER_NAME> --global
```

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

1. Open the Docker Desktop settings and navigate to the **Builders** tab.
2. Find the cloud builder under **Available builders**.
3. Open the drop-down menu and select **Use**.

   ![Selecting the cloud builder as default using the Docker Desktop GUI](/build/images/set-default-builder-gui.webp)

{{< /tab >}}
{{< /tabs >}}

Changing your default builder with `docker buildx use` only changes the default
builder for the `docker buildx build` command. The `docker build` command still
uses the `default` builder, unless you specify the `--builder` flag explicitly.

If you use build scripts, such as `make`, we recommend that you update your
build commands from `docker build` to `docker buildx build`, to avoid any
confusion with regards to builder selection. Alternatively, you can run `docker
buildx install` to make the default `docker build` command behave like `docker
buildx build`, without discrepancies.

## Use with Docker Compose

To build with Docker Build Cloud using `docker compose build`, first set the
cloud builder as your selected builder, then run your build.

> **Note**
>
> Make sure you're using a supported version of Docker Compose, see
> [Prerequisites](setup.md#prerequisites).

```console
$ docker buildx use cloud-<ORG>-<BUILDER_NAME>
$ docker compose build
```

In addition to `docker buildx use`, you can also use the `docker compose build
--builder` flag or the [`BUILDX_BUILDER` environment
variable](../building/env-vars.md#buildx_builder) to select the cloud builder.

## Loading build results

Building with `--tag` loads the build result to the local image store
automatically when the build finishes. To build without a tag and load the
result, you must pass the `--load` flag.

Loading the build result for multi-platform images is not supported. Use the
`docker buildx build --push` flag when building multi-platform images to push
the output to a registry.

```console
$ docker buildx build --builder cloud-<ORG>-<BUILDER_NAME> \
  --platform linux/amd64,linux/arm64 \
  --tag <IMAGE> \
  --push .
```

If you want to build with a tag, but you don't want to load the results to your
local image store, you can export the build results to the build cache only:

```console
$ docker buildx build --builder cloud-<ORG>-<BUILDER_NAME> \
  --platform linux/amd64,linux/arm64 \
  --tag <IMAGE> \
  --output type=cacheonly .
```

## Multi-platform builds

To run multi-platform builds, you must specify all of the platforms that you
want to build for using the `--platform` flag.

```console
$ docker buildx build --builder cloud-<ORG>-<BUILDER_NAME> \
  --platform linux/amd64,linux/arm64 \
  --tag <IMAGE> \
  --push .
```

If you don't specify the platform, the cloud builder automatically builds for the
architecture matching your local environment.

To learn more about building for multiple platforms, refer to [Multi-platform
builds](/build/building/multi-platform/).

## Cloud builds in Docker Desktop

The Docker Desktop [Builds view](/desktop/use-desktop/builds/) works with
Docker Build Cloud out of the box. This view can show information about not only your
own builds, but also builds initiated by your team members using the same
builder.

Teams using a shared builder get access to information such as:

- Ongoing and completed builds
- Build configuration, statistics, dependencies, and results
- Build source (Dockerfile)
- Build logs and errors

This lets you and your team work collaboratively on troubleshooting and
improving build speeds, without having to send build logs and benchmarks back
and forth between each other.
