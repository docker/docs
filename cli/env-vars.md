---
title: Docker CLI environment variables
description: >
  Some aspects of the Docker CLI is configurable using environment variables
keywords: cli, environment variables, vars, configuration
---

In addition to the [configuration file](./config-file.md), the Docker CLI also
honors the following environment variables to configure or change aspects of
the Docker CLI client.

The following table describes the environment variables that you can use to
configure the Docker CLI:

| Variable                      | Description                                                                                                                                                                                                                                                  |
| :---------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `DOCKER_API_VERSION`          | Override the negotiated API version to use for debugging (e.g. `1.19`)                                                                                                                                                                                       |
| `DOCKER_CERT_PATH`            | Location of your authentication keys. This variable is used both by the `docker` CLI and the [`dockerd` daemon](../engine/reference/commandline/dockerd.md)                                                                                                  |
| `DOCKER_CONFIG`               | The location of your client configuration files.                                                                                                                                                                                                             |
| `DOCKER_CONTENT_TRUST_SERVER` | The URL of the Notary server to use. Defaults to the same URL as the registry.                                                                                                                                                                               |
| `DOCKER_CONTENT_TRUST`        | When set Docker uses notary to sign and verify images. Equates to `--disable-content-trust=false` for build, create, pull, push, run.                                                                                                                        |
| `DOCKER_CONTEXT`              | Name of the `docker context` to use (overrides `DOCKER_HOST` environment variable and default context set with `docker context use`)                                                                                                                         |
| `DOCKER_DEFAULT_PLATFORM`     | Default platform for commands that take the `--platform` flag.                                                                                                                                                                                               |
| `DOCKER_HIDE_LEGACY_COMMANDS` | When set, Docker hides "legacy" top-level commands (such as `docker rm`, and `docker pull`) in `docker help` output, and only `Management commands` per object-type (e.g., `docker container`) are printed. This may become the default in a future release. |
| `DOCKER_HOST`                 | Daemon socket to connect to.                                                                                                                                                                                                                                 |
| `DOCKER_TLS_VERIFY`           | When set Docker uses TLS and verifies the remote. This variable is used both by the `docker` CLI and the [`dockerd` daemon](../engine/reference/commandline/dockerd.md)                                                                                      |
| `BUILDKIT_PROGRESS`           | Set type of progress output (`auto`, `plain`, `tty`) when [building](../engine/reference/commandline/build.md). Use plain to show container output (default `auto`).                                                                                         |

Additionally, the Docker CLI honors the de-facto standard environment variables
for configuring HTTP proxies:

| Variable      | Description                                                                    |
| :------------ | :----------------------------------------------------------------------------- |
| `HTTP_PROXY`  | Proxy URL, unless overridden by NoProxy.                                       |
| `HTTPS_PROXY` | Proxy URL, unless overridden by NoProxy.                                       |
| `NO_PROXY`    | Comma-separated values specifying hosts that should be excluded from proxying. |

For more options on configuring the Docker CLI to use a proxy, see
[Configure Docker to use a proxy](../network/proxy.md)
