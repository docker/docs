---
title: Reference documentation
notoc: true
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, and file formats.

## File formats

| File format                                            | Description                                                     |
|:-------------------------------------------------------|:----------------------------------------------------------------|
| [Dockerfile](/engine/reference/builder/)               | Defines the contents and startup behavior of a single container |
| [Compose file](/compose/compose-file/)                 | Defines a multi-container application                           |
| [Docker Cloud Stack file](/docker-cloud/apps/stack-yaml-reference/) | Defines a multi-container application for Docker Cloud          |


## Command-line interfaces (CLIs)

| CLI                                                   | Description                                                                                                      |
|:------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------|
| [Engine CLI](/engine/reference/commandline/cli/)          | The main CLI for Docker, includes all `docker` and [`dockerd`](/engine/reference/commandline/dockerd/) commands|
| [Compose CLI](/compose/reference/overview/)           | The CLI for Docker Compose, which allows you to build and run multi-container applications                       |
| [Machine CLI](/machine/reference/)                    | Manages virtual machines that are pre-configured to run Docker                                                   |
| [UCP CLI](/datacenter/ucp/2.1/reference/cli/index.md) | Manages a Universal Control Plane instance                                                                       |
| [DTR CLI](/datacenter/dtr/2.2/reference/cli/index.md) | Manages a trusted registry                                                                                       |

## Application programming interfaces (APIs)

| API                                                        | Description                                                                            |
|:-----------------------------------------------------------|:---------------------------------------------------------------------------------------|
| [Cloud API](/apidocs/docker-cloud/)                        | Enables programmatic management of your Docker application running on a cloud provider |
| [Engine API](/engine/api/)                                 | The main API for Docker, provides programmatic access to a [daemon](/glossary/#daemon) |
| [Registry API](/registry/spec/api/)                        | Facilitates distribution of images to the engine                                       |
| [Trusted Registry API](/datacenter/dtr/2.2/reference/api/) | Provides programmatic access to a trusted registry                                     |


## Drivers and specifications

| Driver                                                 | Description                                                                        |
|:-------------------------------------------------------|:-----------------------------------------------------------------------------------|
| [Image specification](/registry/spec/manifest-v2-2/)   | Describes the various components of a Docker image                                 |
| [Machine drivers](/machine/drivers/os-base/)           | Enables support for given cloud providers when provisioning resources with Machine |
| [Registry token authentication](/registry/spec/auth/)  | Outlines the Docker registry authentication scheme                                 |
| [Registry storage drivers](/registry/storage-drivers/) | Enables support for given cloud providers when storing images with Registry        |
