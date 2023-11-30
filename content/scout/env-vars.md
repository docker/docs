---
title: Environment variables used by Docker Scout
linkTitle: Docker Scout environment variables
description: Configure how the behavior of Docker Scout CLI commands using these environment variables
keywords: scout, supply chain, cli, env, variables, configure
---

The following environment variables are available to configure the Docker Scout
CLI commands, and the corresponding `docker/scout-cli` container image:

| Name                                    | Format  | Description                                                                                 |
| :-------------------------------------- | ------- | :------------------------------------------------------------------------------------------ |
| DOCKER_SCOUT_CACHE_FORMAT               | String  | Format of the local image cache; can be `oci` or `tar` (default: `oci`)                     |
| DOCKER_SCOUT_CACHE_DIR                  | String  | Directory where the local SBOM cache is stored (default: `$HOME/.docker/scout`)             |
| DOCKER_SCOUT_NO_CACHE                   | Boolean | When set to `true`, disables the use of local SBOM cache                                    |
| DOCKER_SCOUT_OFFLINE                    | Boolean | Use [offline mode](#offline-mode) when indexing SBOM                                        |
| DOCKER_SCOUT_REGISTRY_TOKEN             | String  | Token for authenticating to a registry when pulling images                                  |
| DOCKER_SCOUT_REGISTRY_USER              | String  | Username for authenticating to a registry when pulling images                               |
| DOCKER_SCOUT_REGISTRY_PASSWORD          | String  | Password or personal access token for authenticating to a registry when pulling images      |
| DOCKER_SCOUT_HUB_USER                   | String  | Docker Hub username for authenticating to the Docker Scout backend                          |
| DOCKER_SCOUT_HUB_PASSWORD               | String  | Docker Hub password or personal access token for authenticating to the Docker Scout backend |
| DOCKER_SCOUT_NEW_VERSION_WARN           | Boolean | Warn about new versions of the Docker Scout CLI                                             |
| DOCKER_SCOUT_EXPERIMENTAL_WARN          | Boolean | Warn about experimental features                                                            |
| DOCKER_SCOUT_EXPERIMENTAL_POLICY_OUTPUT | Boolean | Disable experimental output for policy evaluation                                           |

## Offline mode

Under normal operation, Docker Scout cross-references external systems, such as
npm, NuGet, or proxy.golang.org, to retrieve additional information about
packages found in your image.

When `DOCKER_SCOUT_OFFLINE` is set to `true`, Docker Scout image analysis runs
in offline mode. Offline mode means Docker Scout doesn't make outbound requests
to external systems.

To use offline mode:

```console
$ export DOCKER_SCOUT_OFFLINE=true
```
