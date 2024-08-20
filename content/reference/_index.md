---
title: Reference documentation
description: Find reference documentation for the Docker platform’s various APIs, CLIs, and file formats
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
  link: /engine/reference/commandline/cli/
- title: Compose CLI
  description: The CLI for Docker Compose, for building and running multi-container
    applications.
  icon: subtitles
  link: /compose/reference/
- title: Daemon CLI (dockerd)
  description: Persistent process that manages containers.
  icon: developer_board
  link: /reference/cli/dockerd/
grid_apis:
- title: Engine API
  description: The main API for Docker, provides programmatic access to a daemon.
  icon: api
  link: /engine/api/
- title: Registry API
  description: Facilitates distribution of images to the engine.
  icon: storage
  link: /registry/spec/api/
- title: Docker Hub API
  description: API to interact with Docker Hub.
  icon: communities
  link: /reference/api/docker-hub/latest/
- title: DVP Data API
  description: API for Docker Verified Publishers to fetch analytics data.
  icon: area_chart
  link: /reference/api/docker-hub/dvp/
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, drivers and specifications, and file formats.

## File formats

{{< grid items="grid_files" >}}

## Command-line interfaces (CLIs)

{{< grid items="grid_clis" >}}

## Application programming interfaces (APIs)

{{< grid items="grid_apis" >}}
