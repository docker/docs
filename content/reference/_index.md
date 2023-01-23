---
title: Reference documentation
description: This section includes the reference documentation for the Docker platform’s various APIs, CLIs, and file formats.
notoc: true
fileGrid:
  - title: "Dockerfile"
    description: "Defines the contents and startup behavior of a single container."
    icon: "edit_document"
    link: "/reference/builder"
  - title: "Compose file"
    description: "Defines a multi-container application."
    icon: "flowsheet"
    link: "/compose/compose-file"
cliGrid:
  - title: "Docker CLI"
    description: "The main CLI for Docker, includes all `docker` commands."
    icon: "terminal"
    link: "/engine/reference/commandline/cli"
  - title: "Compose CLI"
    description: "The CLI for Docker Compose, which allows you to build and run multi-container applications."
    icon: "clear_all"
    link: "/compose/reference"
  - title: "Daemon CLI (dockerd)"
    description: "Persistent process that manages containers."
    icon: "toc"
    link: "/engine/reference/commandline/dockerd"
apiGrid:
  - title: "Engine API"
    description: "The main API for Docker, provides programmatic access to a daemon."
    icon: "domain_verification"
    link: "/engine/api"
  - title: "Registry API"
    description: "Facilitates distribution of images to the engine."
    icon: "storage"
    link: "/registry/spec/api"
  - title: "Docker Hub API"
    description: "API to interact with Docker Hub."
    icon: "sync"
    link: "/docker-hub/api/latest"
  - title: "DVP Data API"
    description: "API for Docker Verified Publishers to fetch analytics data."
    icon: "insights"
    link: "/docker-hub/api/dvp"
specGrid:
  - title: "Image specification"
    description: "Describes the various components of a Docker image."
    icon: "schema"
    link: "/registry/spec/manifest-v2-2"
  - title: "Registry token authentication"
    description: "Outlines the Docker Registry authentication schemes."
    icon: "lock_open"
    link: "/registry/spec/auth"
  - title: "Registry storage drivers"
    description: "Enables support for given cloud providers when storing images with Registry."
    icon: "database"
    link: "/registry/storage-drivers"
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, drivers and specifications, and file formats.

## File formats

{{< grid fileGrid >}}

## Command-line interfaces (CLIs)

{{< grid cliGrid >}}

## Application programming interfaces (APIs)

{{< grid apiGrid >}}

## Drivers and specifications

{{< grid specGrid >}}
