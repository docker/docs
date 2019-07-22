---
title: Docker Assemble images
description: Building Docker Assemble images
keywords: Assemble, Docker Enterprise, plugin, Spring Boot, .NET, c#, F#
---

## Multi-platform images

By default, Docker Assemble builds images for the `linux/amd64` platform and exports them to the local Docker image store. This is also true when running Docker Assemble on Windows or macOS. For some application frameworks, Docker Assemble can build multi-platform images to support running on several host platforms. For example, `linux/amd64` and `windows/amd64`.

To support multi-platform images, images must be pushed to a registry instead of the local image store. This is because the local image store can only import uni-platform images which match its platform.

To enable the multi-platform mode, use the `--push` option. For example:

```bash
$ docker assemble build --push /path/to/my/project
```

To push to an insecure (unencrypted) registry, use `--push-insecure` instead of `--push`.

## Custom base images

Docker Assemble allows you to override the base images for building and running your project.  For example, the following `docker-assemble.yaml` file defines `maven:3-ibmjava-8-alpine` as the base build image and `openjdk:8-jre-alpine` as the base runtime image (for linux/amd64 platform).

```
version: "0.2.0"
springboot:
  enabled: true
  build-image: "maven:3-ibmjava-8-alpine"
  runtime-images:
    linux/amd64: "openjdk:8-jre-alpine"
```

Linux-based images must be Debian, Red Hat, or Alpine-based and have a standard environment with:

- `find`

- `xargs`

- `grep`

- `true`

- a standard POSIX shell (located at `/bin/sh`)

These tools are required for internal inspection that Docker Assemble performs on the images. Depending on the type of your project and your configuration, the base images must meet other requirements as described in the following sections.

### Spring Boot

Install Java JDK and maven on the base build image and ensure it is available in `$PATH`. Install a maven settings file as `/usr/share/maven/ref/settings-docker.xml` (irrespective of the install location of Maven).

Ensure the base runtime image has Java JRE installed and is available in `$PATH`. The build and runtime image must have the same version of Java installed.

Supported build platform:

- `linux/amd64`

Supported runtime platforms:

- `linux/amd64`

- `windows/amd64`

### ASP.NET Core

Install .NET Core SDK on the base build image and ensure it includes the [.NET Core command-line interface tools](https://docs.microsoft.com/en-us/dotnet/core/tools/?tabs=netcore2x).

Install [.NET Core command-line interface tools](https://docs.microsoft.com/en-us/dotnet/core/tools/?tabs=netcore2x) on the base runtime image.

Supported build platform:

- `linux/amd64`

Supported runtime platforms:

- `linux/amd64`

- `windows/amd64`

## Bill of lading

Docker Assemble generates a bill of lading when building an image. This contains information about the tools, base images, libraries, and packages used by Assemble to build the image and that are included in the runtime image. The bill of lading has two parts – one for build and one for runtime.

The build part includes:

- The base image used
- A map of packages installed and their versions
- A map of libraries used for the build and their versions
- A map of build tools and their corresponding versions

The runtime part includes:

- The base image used
- A map of packages installed and their versions
- A map of runtime tools and their versions

You can find the bill of lading by inspecting the resulting image. It is stored using the label `com.docker.assemble.bill-of-lading`:

{% raw %}
```bash
$ docker image inspect --format '{{ index .Config.Labels "com.docker.assemble.bill-of-lading" }}' <image>
```
{% endraw %}

> **Note:** The bill of lading is only supported on the `linux/amd64` platform and only for images which are based on Alpine (`apk`), Red Hat (`rpm`) or Debian (`dpkg-query`).

## Health checks

Docker Assemble only supports health checks on `linux/amd64` based runtime images and require certain additional commands to be present depending on the value of `image.healthcheck.kind`:

- `simple-tcpport-open:` requires the `nc` command
- `springboot:` requires the `curl` and `jq` commands

On Alpine (`apk`) and Debian (`dpkg`) based images, these dependencies are installed automatically. For other base images, you must ensure they are present in the images you specify.

If your base runtime image lacks the necessary commands, you may need to set `image.healthcheck.kind` to `none` in your `docker-assemble.yaml` file.
