---
title: Reference documentation
linkTitle: Reference
layout: wide
description: Find reference documentation for the Docker platform’s various APIs, CLIs, and file formats
params:
  icon: terminal
  notoc: true
  grid_files:
  - title: Dockerfile
    description: Defines the contents and startup behavior of a single container.
    icon: edit_document
    link: /reference/dockerfile/
  - title: Compose file
    description: Defines a multi-container application.
    icon: polyline
    link: /reference/compose-file/
  grid_clis:
  - title: Docker CLI
    description: The main Docker CLI, includes all `docker` commands.
    icon: terminal
    link: /reference/cli/docker/
  - title: Compose CLI
    description: The CLI for Docker Compose, for building and running multi-container
      applications.
    icon: subtitles
    link: /reference/cli/docker/compose/
  - title: Daemon CLI (dockerd)
    description: Persistent process that manages containers.
    icon: developer_board
    link: /reference/cli/dockerd/
  grid_apis:
  - title: Engine API
    description: The main API for Docker, provides programmatic access to a daemon.
    icon: api
    link: /reference/api/engine/
  - title: Docker Hub API
    description: API to interact with Docker Hub.
    icon: communities
    link: /reference/api/hub/latest/
  - title: DVP Data API
    description: API for Docker Verified Publishers to fetch analytics data.
    icon: area_chart
    link: /reference/api/hub/dvp/
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, drivers and specifications, and file formats.

## File formats

{{< grid items="grid_files" >}}

## Command-line interfaces (CLIs)

{{< grid items="grid_clis" >}}

## Application programming interfaces (APIs)

{{< grid items="grid_apis" >}}
