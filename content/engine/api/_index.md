---
title: Develop with Docker Engine API
description: Learn how you can use Docker Engine API and SDKs in the language of your choice.
keywords: developing, api, Docker Engine API, API versions, SDK versions
aliases:
- /engine/reference/api/
- /engine/reference/api/docker_remote_api/
- /reference/api/
- /reference/api/docker_remote_api/
- /reference/api/docker_remote_api_v1.0/
- /reference/api/docker_remote_api_v1.1/
- /reference/api/docker_remote_api_v1.2/
- /reference/api/docker_remote_api_v1.3/
- /reference/api/docker_remote_api_v1.4/
- /reference/api/docker_remote_api_v1.5/
- /reference/api/docker_remote_api_v1.6/
- /reference/api/docker_remote_api_v1.7/
- /reference/api/docker_remote_api_v1.8/
- /reference/api/docker_remote_api_v1.9/
- /reference/api/docker_remote_api_v1.10/
- /reference/api/docker_remote_api_v1.11/
- /reference/api/docker_remote_api_v1.12/
- /reference/api/docker_remote_api_v1.13/
- /reference/api/docker_remote_api_v1.14/
- /reference/api/docker_remote_api_v1.15/
- /reference/api/docker_remote_api_v1.16/
- /reference/api/docker_remote_api_v1.17/
---

Docker provides an API for interacting with the Docker daemon (called the Docker
Engine API), as well as SDKs for Go and Python. The SDKs allow you to efficiently build and
scale Docker apps and solutions. If Go or Python don't work
for you, you can use the Docker Engine API directly.

For information about Docker Engine SDKs, see [Develop with Docker Engine SDKs](sdk/index.md).

The Docker Engine API is a RESTful API accessed by an HTTP client such as `wget` or
`curl`, or the HTTP library which is part of most modern programming languages.

## View the API reference

You can
[view the reference for the latest version of the API](latest/index.html)
or [choose a specific version](/engine/api/version-history/).

## Versioned API and SDK

The version of the Docker Engine API you should use depends upon the version of
your Docker daemon and Docker client.

A given version of the Docker Engine SDK supports a specific version of the
Docker Engine API, as well as all earlier versions. If breaking changes occur,
they are documented prominently.

> **Note**
>
> The Docker daemon and client don't necessarily need to be the same version
> at all times. However, keep the following in mind.
>
> - If the daemon is newer than the client, the client doesn't know about new
>   features or deprecated API endpoints in the daemon.
>
> - If the client is newer than the daemon, the client can request API
>   endpoints that the daemon doesn't know about.

A new version of the API is released when new features are added. The Docker API
is backward-compatible, so you don't need to update code that uses the API
unless you need to take advantage of new features.

To see the highest version of the API your Docker daemon and client support, use
`docker version`:

```console
$ docker version
Client: Docker Engine - Community
 Version:           24.0.0
 API version:       1.43
 Go version:        go1.20.4
 Git commit:        98fdcd7
 Built:             Mon May 15 18:48:45 2023
 OS/Arch:           linux/arm64
 Context:           default

Server: Docker Engine - Community
 Engine:
  Version:          24.0.0
  API version:      1.43 (minimum version 1.12)
  Go version:       go1.20.4
  Git commit:       1331b8c
  Built:            Mon May 15 18:48:45 2023
  OS/Arch:          linux/arm64
  Experimental:     false
  ...
```

You can specify the API version to use in any of the following ways:

- When using the SDK, use the latest version. At a minimum, use the version
  that incorporates the API version with the features you need.
- When using `curl` directly, specify the version as the first part of the URL.
  For instance, if the endpoint is `/containers/` you can use
  `/v1.42/containers/`.
- To force the Docker CLI or the Docker Engine SDKs to use an older version
  of the API than the version reported by `docker version`, set the
  environment variable `DOCKER_API_VERSION` to the correct version. This works
  on Linux, Windows, or macOS clients.

  ```console
  $ DOCKER_API_VERSION='1.41'
  ```

  While the environment variable is set, that version of the API is used, even
  if the Docker daemon supports a newer version. This environment variable
  disables API version negotiation, so you should only use it if you must
  use a specific version of the API, or for debugging purposes.

- The Docker Go SDK allows you to enable API version negotiation, automatically
  selects an API version that's supported by both the client and the Docker Engine
  that's in use.
- For the SDKs, you can also specify the API version programmatically as a
  parameter to the `client` object. See the
  [Go constructor](https://pkg.go.dev/github.com/docker/docker/client#NewClientWithOpts)
  or the
  [Python SDK documentation for `client`](https://docker-py.readthedocs.io/en/stable/client.html).

### API version matrix

{{< include "api-version-matrix.md" >}}

### Archived API versions

You can find archived documentation for older versions of the API
in the [Docker code repository on GitHub](https://github.com/moby/moby/tree/v1.9.1/docs/reference/api)