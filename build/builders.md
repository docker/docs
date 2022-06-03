---
title: Using Docker Build
description: About building with existing builders
keywords: build, buildx, buildkit
---


## Docker Build Components

Currently Engine holds two build components:
* Legacy Builder - Old builder implementation. It's executed through the `docker build` command.
* Buildkit - Active by default in new and latest installations of Docker. It's executed through the `docker buildx build` command.

### Docker Build Default Configuration

The backend you have active depends on your Docker's installation.
If you have installed:

* __Docker Desktop (all platforms)__: BuildKit is enabled by default; no additional configuration is required.
* In Linux, if you have:
  * Engine 18.09 and higher you have BuildKit enabled by default.
  * Engine releases older than 18.09 only include the Legacy Builder component, so we highly recommend you upgrade your installation or alternatively install Docker Desktop.


### Docker Build Configuration

Backend configuration is performed using the `DOCKER_BUILDKIT` environment variable. 

If not yet set as default, to set the BuildKit as the default backend behind the `docker build` command, you can:

__Method 1__

From your console, run:

```console
$ DOCKER_BUILDKIT=1 docker build .
```
` 
Depending on the value for this variable:
* __DOCKER_BUILDKIT=1__ - BuildKit is executed both when you issue the `docker build` and `docker buildx build` command. With this configuration `docker build` becomes an alias for `docker buildx build`.
* __DOCKER_BUILDKIT=0__ - When you issue `docker build` the build is executed using the Legacy Builder and BuildKit when you issue the `docker buildx build` command.

> See this [deprecation notice](https://github.com/docker/cli/blob/master/docs/deprecated.md#legacy-builder-for-linux-images).

__Method2__

To enable BuildKit backend by default, set [daemon configuration](/engine/reference/commandline/dockerd/#daemon-configuration-file)
in `/etc/docker/daemon.json` feature to `true` and restart the daemon. If the
`daemon.json` file doesn't exist, create new file called `daemon.json` and then
add the following to the file:

```json
{
  "features": {
    "buildkit": true
  }
}
```

### `docker buildx build` command and BuildKit

Regardless of your build backend configuration, when you issue the [`docker buildx build` command](../engine/reference/commandline/buildx_build.md) you are, by design of the command, always executing a build with Buildkit.
See also, [Build with Buildx](../buildx/working-with-buildx.md#build-with-buildx) guide for more details.

## `docker build` command and Legacy Builder
To perform a build with the Legacy Builder because you need, for example, to build a Windows image (Link to Backends support) you need to have `DOCKER_BUILDKIT` set as `0` and use the `docker build` command. 

