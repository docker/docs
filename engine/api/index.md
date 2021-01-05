---
title: Develop with Docker Engine API
description: Using Docker APIs to automate Docker tasks in your language of choice
keywords: developing, api
redirect_from:
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
Engine API), as well as SDKs for Go and Python. The SDKs allow you to build and
scale Docker apps and solutions quickly and easily. If Go or Python don't work
for you, you can use the Docker Engine API directly.

For information about Docker Engine SDKs, see [Develop with Docker Engine SDKs](/engine/api/sdk/).

The Docker Engine API is a RESTful API accessed by an HTTP client such as `wget` or
`curl`, or the HTTP library which is part of most modern programming languages.

## View the API reference

You can
[view the reference for the latest version of the API](/engine/api/latest/)
or [choose a specific version](/engine/api/version-history/).

## Versioned API and SDK

The version of the Docker Engine API you should use depends upon the version of
your Docker daemon and Docker client.

A given version of the Docker Engine SDK supports a specific version of the
Docker Engine API, as well as all earlier versions. If breaking changes occur,
they are documented prominently.

> Daemon and client API mismatches
>
> The Docker daemon and client do not necessarily need to be the same version
> at all times. However, keep the following in mind.
>
> - If the daemon is newer than the client, the client does not know about new
>   features or deprecated API endpoints in the daemon.
>
> - If the client is newer than the daemon, the client can request API
>   endpoints that the daemon does not know about.

A new version of the API is released when new features are added. The Docker API
is backward-compatible, so you do not need to update code that uses the API
unless you need to take advantage of new features.

To see the highest version of the API your Docker daemon and client support, use
`docker version`:

```bash
$ docker version

Client: Docker Engine - Community
 Version:           20.10.0
 API version:       1.41
 Go version:        go1.13.15
 Git commit:        7287ab3
 Built:             Tue Dec  8 19:00:39 2020
 OS/Arch:           linux/amd64
 Context:           default
 Experimental:      true
Server: Docker Engine - Community
 Engine:
  Version:          20.10.0
  API version:      1.41 (minimum version 1.12)
  Go version:       go1.13.15
  Git commit:       eeddea2
  Built:            Tue Dec  8 18:58:12 2020
  OS/Arch:          linux/amd64
  ...
```

You can specify the API version to use, in one of the following ways:

- When using the SDK, use the latest version you can, but at least the version
  that incorporates the API version with the features you need.

- When using `curl` directly, specify the version as the first part of the URL.
  For instance, if the endpoint is `/containers/`, you can use
  `/v1.41/containers/`.

- To force the Docker CLI or the Docker Engine SDKs to use an old version
  version of the API than the version reported by `docker version`, set the
  environment variable `DOCKER_API_VERSION` to the correct version. This works
  on Linux, Windows, or macOS clients.

  ```bash
  DOCKER_API_VERSION='1.41'
  ```

  While the environment variable is set, that version of the API is used, even
  if the Docker daemon supports a newer version. This environment variable
  disables API version negotiation, and as such should only be used if you must
  use a specific version of the API, or for debugging purposes.

- The Docker Go SDK allows you to enable API version negotiation, automatically
  selects an API version that is supported by both the client, and the Docker Engine
  that is used.

- For the SDKs, you can also specify the API version programmatically, as a
  parameter to the `client` object. See the
  [Go constructor](https://github.com/moby/moby/blob/v19.03.6/client/client.go#L119){: target="_blank" rel="noopener" class="_"}
  or the
  [Python SDK documentation for `client`](https://docker-py.readthedocs.io/en/stable/client.html){: target="_blank" rel="noopener" class="_"}.

### API version matrix

{% include api-version-matrix.md %}

### Archived API versions

Documentation for older versions of the API has been archived, but can be found
in the [docker code repository on GitHub](https://github.com/moby/moby/tree/v1.9.1/docs/reference/api){: target="_blank" rel="noopener" class="_"}
