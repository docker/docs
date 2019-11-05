---
title: Reference documentation
notoc: true
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, and file formats.

## File formats

| File format                                                         | Description                                                     |
|:--------------------------------------------------------------------|:----------------------------------------------------------------|
| [Dockerfile](/engine/reference/builder/)                            | Defines the contents and startup behavior of a single container |
| [Compose file](/compose/compose-file/)                              | Defines a multi-container application                           |
| [Daemon CLI (dockerd)]()| 


## Command-line interfaces (CLIs)

| CLI                                                           | Description                                                                                                     |
|:--------------------------------------------------------------|:----------------------------------------------------------------------------------------------------------------|
| [Docker CLI](/engine/reference/commandline/cli/)              | The main CLI for Docker, includes all `docker` commands |
| [Compose CLI](/compose/reference/overview/)                   | The CLI for Docker Compose, which allows you to build and run multi-container applications                      |
| [Daemon CLI (dockerd)](/engine/reference/commandline/dockerd/)                            | Persistent process that manages containers                                                 |
| [DTR CLI](/reference/dtr/{{ site.dtr_version }}/cli/index.md) | Deploy and manage Docker Trusted Registry                                                                       |
| [UCP CLI](/reference/ucp/{{ site.ucp_version }}/cli/index.md) | Deploy and manage Universal Control Plane                                                                       |

## Application programming interfaces (APIs)

| API                                                   | Description                                                                            |
|:------------------------------------------------------|:---------------------------------------------------------------------------------------|
| [Engine API](/engine/api/)                            | The main API for Docker, provides programmatic access to a daemon |
| [DTR API](/reference/dtr/{{ site.dtr_version }}/api/) | Provides programmatic access to a Docker Trusted Registry deployment                   |
| [Registry API](/registry/spec/api/)                   | Facilitates distribution of images to the engine                                       |
| [Template API](app-template/api-reference)| Allows users to create new Docker applications by using a library of templates.|
| [UCP API](/reference/ucp/{{ site.ucp_version }}/api/) | Provides programmatic access to a Universal Control Plane deployment                   |

## Drivers and specifications

| Driver                                                 | Description                                                                        |
|:-------------------------------------------------------|:-----------------------------------------------------------------------------------|
| [Image specification](/registry/spec/manifest-v2-2/)   | Describes the various components of a Docker image                                 |
| [Registry token authentication](/registry/spec/auth/)  | Outlines the Docker registry authentication scheme                                 |
| [Registry storage drivers](/registry/storage-drivers/) | Enables support for given cloud providers when storing images with Registry        |

## Compliance control reference

| Reference                                                      | Description                                                                                                       |
|:---------------------------------------------------------------|:------------------------------------------------------------------------------------------------------------------|
| [NIST 800-53 control reference](/compliance/reference/800-53/) | All of the NIST 800-53 Rev. 4 controls applicable to Docker Enterprise Edition can be referenced in this section. |
