---
title: Docker Engine API
description: Learn how you can use Docker Engine API and SDKs in the language of your choice.
keywords: developing, api, Docker Engine API, API versions, SDK versions
aliases:
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
  - /engine/reference/api/
  - /engine/reference/api/docker_remote_api/
  - /engine/api/
---

Docker provides an API for interacting with the Docker daemon (called the Docker
Engine API), as well as SDKs for Go and Python. The SDKs allow you to efficiently build and
scale Docker apps and solutions. If Go or Python don't work
for you, you can use the Docker Engine API directly.

For information about Docker Engine SDKs, see [Develop with Docker Engine SDKs](sdk/_index.md).

The Docker Engine API is a RESTful API accessed by an HTTP client such as `wget` or
`curl`, or the HTTP library which is part of most modern programming languages.

## View the API reference

You can
[view the reference for the latest version of the API](/reference/api/engine/version/v{{% param latest_engine_api_version %}}.md)
or [choose a specific version](/reference/api/engine/version-history/).

## Versioned API and SDK

The version of the Docker Engine API you should use depends upon the version of
your Docker daemon and Docker client.

A given version of the Docker Engine SDK supports a specific version of the
Docker Engine API, as well as all earlier versions. If breaking changes occur,
they are documented prominently.

> [!NOTE]
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
 Version:           28.4.0
 API version:       1.51
 Go version:        go1.24.7
 Git commit:        d8eb465
 Built:             Wed Sep  3 20:57:32 2025
 OS/Arch:           linux/amd64
 Context:           default

Server: Docker Engine - Community
 Engine:
  Version:          28.4.0
  API version:      1.51 (minimum version 1.24)
  Go version:       go1.24.7
  Git commit:       249d679
  Built:            Wed Sep  3 20:57:32 2025
  OS/Arch:          linux/amd64
  ...
```

You can specify the API version to use in any of the following ways:

- When using the SDK, use the latest version. At a minimum, use the version
  that incorporates the API version with the features you need.
- When using `curl` directly, specify the version as the first part of the URL.
  For instance, if the endpoint is `/containers/` you can use
  `/v{{% param "latest_engine_api_version" %}}/containers/`.
- To force the Docker CLI or the Docker Engine SDKs to use an older version
  of the API than the version reported by `docker version`, set the
  environment variable `DOCKER_API_VERSION` to the correct version. This works
  on Linux, Windows, or macOS clients.

  {{% apiVersionPrevious.inline %}}
  {{- $version := site.Params.latest_engine_api_version }}
  {{- $parts := strings.Split $version "." }}
  {{- $major := cast.ToInt (index $parts 0) }}
  {{- $minor := cast.ToInt (index $parts 1) }}
  ```console
  $ DOCKER_API_VERSION={{ $major }}.{{ math.Sub $minor 1 }}
  ```
  {{% /apiVersionPrevious.inline %}}

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

| Docker version | Maximum API version                          | Change log                                                         |
|:---------------|:---------------------------------------------|:-------------------------------------------------------------------|
| 28.3           | [1.51](/reference/api/engine/version/v1.51/) | [changes](/reference/api/engine/version-history/#v151-api-changes) |
| 28.2           | [1.50](/reference/api/engine/version/v1.50/) | [changes](/reference/api/engine/version-history/#v150-api-changes) |
| 28.1           | [1.49](/reference/api/engine/version/v1.49/) | [changes](/reference/api/engine/version-history/#v149-api-changes) |
| 28.0           | [1.48](/reference/api/engine/version/v1.48/) | [changes](/reference/api/engine/version-history/#v148-api-changes) |
| 27.5           | [1.47](/reference/api/engine/version/v1.47/) | [changes](/reference/api/engine/version-history/#v147-api-changes) |
| 27.4           | [1.47](/reference/api/engine/version/v1.47/) | [changes](/reference/api/engine/version-history/#v147-api-changes) |
| 27.3           | [1.47](/reference/api/engine/version/v1.47/) | [changes](/reference/api/engine/version-history/#v147-api-changes) |
| 27.2           | [1.47](/reference/api/engine/version/v1.47/) | [changes](/reference/api/engine/version-history/#v147-api-changes) |
| 27.1           | [1.46](/reference/api/engine/version/v1.46/) | [changes](/reference/api/engine/version-history/#v146-api-changes) |
| 27.0           | [1.46](/reference/api/engine/version/v1.46/) | [changes](/reference/api/engine/version-history/#v146-api-changes) |
| 26.1           | [1.45](/reference/api/engine/version/v1.45/) | [changes](/reference/api/engine/version-history/#v145-api-changes) |
| 26.0           | [1.45](/reference/api/engine/version/v1.45/) | [changes](/reference/api/engine/version-history/#v145-api-changes) |
| 25.0           | [1.44](/reference/api/engine/version/v1.44/) | [changes](/reference/api/engine/version-history/#v144-api-changes) |
| 24.0           | [1.43](/reference/api/engine/version/v1.43/) | [changes](/reference/api/engine/version-history/#v143-api-changes) |
| 23.0           | [1.42](/reference/api/engine/version/v1.42/) | [changes](/reference/api/engine/version-history/#v142-api-changes) |
| 20.10          | [1.41](/reference/api/engine/version/v1.41/) | [changes](/reference/api/engine/version-history/#v141-api-changes) |
| 19.03          | [1.40](/reference/api/engine/version/v1.40/) | [changes](/reference/api/engine/version-history/#v140-api-changes) |
| 18.09          | [1.39](/reference/api/engine/version/v1.39/) | [changes](/reference/api/engine/version-history/#v139-api-changes) |
| 18.06          | [1.38](/reference/api/engine/version/v1.38/) | [changes](/reference/api/engine/version-history/#v138-api-changes) |
| 18.05          | [1.37](/reference/api/engine/version/v1.37/) | [changes](/reference/api/engine/version-history/#v137-api-changes) |
| 18.04          | [1.37](/reference/api/engine/version/v1.37/) | [changes](/reference/api/engine/version-history/#v137-api-changes) |
| 18.03          | [1.37](/reference/api/engine/version/v1.37/) | [changes](/reference/api/engine/version-history/#v137-api-changes) |
| 18.02          | [1.36](/reference/api/engine/version/v1.36/) | [changes](/reference/api/engine/version-history/#v136-api-changes) |
| 17.12          | [1.35](/reference/api/engine/version/v1.35/) | [changes](/reference/api/engine/version-history/#v135-api-changes) |
| 17.11          | [1.34](/reference/api/engine/version/v1.34/) | [changes](/reference/api/engine/version-history/#v134-api-changes) |
| 17.10          | [1.33](/reference/api/engine/version/v1.33/) | [changes](/reference/api/engine/version-history/#v133-api-changes) |
| 17.09          | [1.32](/reference/api/engine/version/v1.32/) | [changes](/reference/api/engine/version-history/#v132-api-changes) |
| 17.07          | [1.31](/reference/api/engine/version/v1.31/) | [changes](/reference/api/engine/version-history/#v131-api-changes) |
| 17.06          | [1.30](/reference/api/engine/version/v1.30/) | [changes](/reference/api/engine/version-history/#v130-api-changes) |
| 17.05          | [1.29](/reference/api/engine/version/v1.29/) | [changes](/reference/api/engine/version-history/#v129-api-changes) |
| 17.04          | [1.28](/reference/api/engine/version/v1.28/) | [changes](/reference/api/engine/version-history/#v128-api-changes) |
| 17.03.1        | [1.27](/reference/api/engine/version/v1.27/) | [changes](/reference/api/engine/version-history/#v127-api-changes) |
| 17.03          | [1.26](/reference/api/engine/version/v1.27/) | [changes](/reference/api/engine/version-history/#v126-api-changes) |
| 1.13.1         | [1.26](/reference/api/engine/version/v1.26/) | [changes](/reference/api/engine/version-history/#v126-api-changes) |
| 1.13           | [1.25](/reference/api/engine/version/v1.26/) | [changes](/reference/api/engine/version-history/#v125-api-changes) |
| 1.12           | [1.24](/reference/api/engine/version/v1.24/) | [changes](/reference/api/engine/version-history/#v124-api-changes) |

### Deprecated API versions

API versions before v1.24 are [deprecated](/engine/deprecated/#deprecate-legacy-api-versions).
You can find archived documentation for deprecated versions of the API in the
code repository on GitHub:

- [Documentation for API versions 1.23 and before](https://github.com/moby/moby/tree/v25.0.0/docs/api).
- [Documentation for API versions 1.17 and before](https://github.com/moby/moby/tree/v1.9.1/docs/reference/api).
