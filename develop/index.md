---
title: Develop with Docker
description: Overview of developer resources
keywords: developer, developing, apps, api, sdk
---

This page contains a list of resources for application developers who would like to build new applications using Docker.

## Prerequisites

Work through the learning modules in [Get started](/get-started/index.md) to understand how to build an image and run it as a containerized application.

## Develop new apps on Docker

If you're just getting started developing a brand new app on Docker, check out
these resources to understand some of the most common patterns for getting the
most benefits from Docker.

- Use [multistage builds](/engine/userguide/eng-image/multistage-build.md){: target="_blank" class="_"} to keep your images lean
- Manage application data using [volumes](/engine/admin/volumes/volumes.md) and [bind mounts](/engine/admin/volumes/bind-mounts.md){: target="_blank" class="_"}
- [Scale your app](/get-started/kube-deploy.md){: target="_blank" class="_"} with kubernetes
- [Scale your app](/get-started/swarm-deploy.md){: target="_blank" class="_"} as a swarm service
- [General application development best practices](/develop/dev-best-practices.md){: target="_blank" class="_"}

## Learn about language-specific app development with Docker

- [Docker for Java developers](https://github.com/docker/labs/tree/master/developer-tools/java/){: target="_blank" class="_"} lab
- [Port a node.js app to Docker](https://github.com/docker/labs/tree/master/developer-tools/nodejs/porting){: target="_blank" class="_"}
- [Ruby on Rails app on Docker](https://github.com/docker/labs/tree/master/developer-tools/ruby){: target="_blank" class="_"} lab
- [Dockerize a .Net Core application](/engine/examples/dotnetcore/){: target="_blank" class="_"}
- [Dockerize an ASP.NET Core application with SQL Server on Linux](/compose/aspnet-mssql-compose/){: target="_blank" class="_"} using Docker Compose

## Advanced development with the SDK or API

After you can write Dockerfiles or Compose files and use Docker CLI, take it to the next level by using Docker Engine SDK for Go/Python or use the HTTP API directly.
