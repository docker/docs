---
title: History and development of Docker Compose
linkTitle: History and development
description: Explore the evolution of Docker Compose from v1 to v2, including CLI changes, YAML versioning, and the Compose Specification.
keywords: compose, compose yaml, swarm, migration, compatibility, docker compose vs docker-compose
weight: 30
aliases:
- /compose/history/
---

This page provides:
 - A brief history of the development of the Docker Compose CLI
 - A clear explanation of the major versions and file formats that make up Compose v1 and Compose v2
 - The main differences between Compose v1 and Compose v2 

## Introduction

![Image showing the main differences between Compose v1 and Compose v2](../images/v1-versus-v2.png)

The previous image shows that the currently supported version of the Docker Compose CLI is Compose v2 which is defined by the [Compose Specification](/reference/compose-file/_index.md).

It also provides a quick snapshot of the differences in file formats, command-line syntax, and top-level elements. This is covered in more detail in the following sections.

### Docker Compose CLI versioning

Version one of the Docker Compose command-line binary was first released in 2014. It was written in Python, and is invoked with `docker-compose`.
Typically, Compose v1 projects include a top-level `version` element in the `compose.yaml` file, with values ranging from `2.0` to `3.8`, which refer to the specific [file formats](#compose-file-format-versioning).

Version two of the Docker Compose command-line binary was announced in 2020, is written in Go, and is invoked with `docker compose`.
Compose v2 ignores the `version` top-level element in the `compose.yaml` file.

### Compose file format versioning

The Docker Compose CLIs are defined by specific file formats. 

Three major versions of the Compose file format for Compose v1 were released:
- Compose file format 1 with Compose 1.0.0 in 2014
- Compose file format 2.x with Compose 1.6.0 in 2016
- Compose file format 3.x with Compose 1.10.0 in 2017

Compose file format 1 is substantially different to all the following formats as it lacks a top-level `services` key.
Its usage is historical and files written in this format don't run with Compose v2.

Compose file format 2.x and 3.x are very similar to each other, but the latter introduced many new options targeted at Swarm deployments.

To address confusion around Compose CLI versioning, Compose file format versioning, and feature parity depending on whether Swarm mode was in use, file format 2.x and 3.x were merged into the [Compose Specification](/reference/compose-file/_index.md). 

Compose v2 uses the Compose Specification for project definition. Unlike the prior file formats, the Compose Specification is rolling and makes the `version` top-level element optional. Compose v2 also makes use of optional specifications - [Deploy](/reference/compose-file/deploy.md), [Develop](/reference/compose-file/develop.md), and [Build](/reference/compose-file/build.md).

To make migration easier, Compose v2 has backwards compatibility for certain elements that have been deprecated or changed between Compose file format 2.x/3.x and the Compose Specification.

## What's next?

- [How Compose works](compose-application-model.md)
- [Compose Specification reference](/reference/compose-file/_index.md)
