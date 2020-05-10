---
description: Summary of samples related to Compose
keywords: documentation, docs, docker, compose, samples
title: Sample apps with Compose
---

The following samples show the various aspects of how to work with Docker
Compose. As a prerequisite, be sure to [install Docker Compose](install.md)
if you have not already done so.

## Key concepts these samples cover

The samples should help you to:

- define services based on Docker images using
  [Compose files](compose-file/index.md) `docker-compose.yml` and
  `docker-stack.yml` files
- understand the relationship between `docker-compose.yml` and
  [Dockerfiles](/engine/reference/builder/)
- learn how to make calls to your application services from Compose files
- learn how to deploy applications and services to a [swarm](../engine/swarm/index.md)

## Samples tailored to demo Compose

These samples focus specifically on Docker Compose:

- [Quickstart: Compose and Django](django.md) - Shows how to use Docker Compose to set up and run a simple Django/PostgreSQL app.

- [Quickstart: Compose and Rails](rails.md) - Shows how to use
Docker Compose to set up and run a Rails/PostgreSQL app.

- [Quickstart: Compose and WordPress](wordpress.md) - Shows how to
use Docker Compose to set up and run WordPress in an isolated environment
with Docker containers.
