---
description: Summary of samples related to Compose
keywords: documentation, docs, docker, compose, samples
title: Sample apps with Compose
---

The following samples show the various aspects of how to work with Docker
Compose. As a prerequisite, be sure to [install Docker Compose](install/index.md)
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
https://github.com/docker/awesome-compose/tree/master/elasticsearch-logstash-kibana/logstash

- [Quickstart: Compose and Django](https://github.com/docker/awesome-compose/tree/master/official-documentation-samples/django/README.md) - Shows how to use Docker Compose to set up and run a simple Django/PostgreSQL app.

- [Quickstart: Compose and Rails](https://github.com/docker/awesome-compose/tree/master/official-documentation-samples/rails/README.md) - Shows how to use
Docker Compose to set up and run a Rails/PostgreSQL app.

- [Quickstart: Compose and WordPress](https://github.com/docker/awesome-compose/tree/master/official-documentation-samples/wordpress/README.md) - Shows how to
use Docker Compose to set up and run WordPress in an isolated environment
with Docker containers.

## Awesome Compose samples

The Awesome Compose samples provide a starting point on how to integrate different frameworks and technologies using Docker Compose. All samples are available in the [Awesome-compose GitHub repo](https://github.com/docker/awesome-compose){:target="_blank" rel="noopener" class="_"} and are ready to run with `docker compose up`.
