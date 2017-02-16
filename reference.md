---
title: Reference documentation
skip-right-nav: true
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, and file formats.

## File formats

| File format | Description |
| ----------- | ----------- |
| [Dockerfile](/engine/reference/builder/) | Defines the contents and startup behavior of a single container |
| [Compose file](/compose/compose-file/) | Defines a multi-container application |
| [Stack file](/docker-cloud/apps/stack-yaml-reference/) | Defines a multi-container application for Docker Cloud |


## Command-line interfaces (CLIs)

| CLI | Description |
| --- | ----------- |
| [Engine CLI](/engine/reference/commandline/) | The main CLI for Docker, includes all `docker` and [`dockerd`](/engine/reference/commandline/dockerd/) commands. |
| [Compose CLI](/compose/reference/overview/) | The CLI for Docker Compose, which allows you to build and run multi-container applications |
| [Machine CLI](/machine/reference/) | Manages virtual machines that are pre-configured to run Docker |
| [UCP tool](/datacenter/ucp/2.0/reference/cli/) | Manages a Universal Control Plane instance |
| [Trusted Registry CLI](/docker-trusted-registry/reference/) | Manages a trusted registry |

## Application programming interfaces (APIs)

| API | Description |
| --- | ----------- |
| [Cloud API](/apidocs/docker-cloud/) | Enables programmatic management of your Docker application running on a cloud provider |
| [Docker ID accounts API](/docker-id/api-reference/) | An API for accessing and updating Docker ID accounts |
| [Engine API](/engine/api/) | The main API for Docker, provides programmatic access to a [daemon](/glossary/#daemon) |
| [Registry API](/registry/spec/api/) | Facilitates distribution of images to the engine |
| [Trusted Registry API](/apidocs/overview/) | Provides programmatic access to a trusted registry |


## Drivers and specifications

| Driver | Description |
| ------ | ----------- |
| [Image specification](/registry/spec/manifest-v2-2/) | Describes the various components of a Docker image |
| [Machine drivers](/machine/drivers/os-base/) | Enables support for given cloud providers when provisioning resources with Machine |
| [Registry token authentication](/registry/spec/auth/) | Outlines the Docker registry authentication scheme |
| [Registry storage drivers](/registry/storage-drivers/) | Enables support for given cloud providers when storing images with Registry |
