---
title: Develop with Docker
description: Overview of developer resources
keywords: developer, developing, apps, api, sdk
---

This page contains a list of resources for application developers who would like to build new applications using Docker.

## Prerequisites

Work through the learning modules in [Get started](../get-started/index.md) to understand how to build an image and run it as a containerized application.

## Develop new apps on Docker

If you're just getting started developing a brand new app on Docker, check out
these resources to understand some of the most common patterns for getting the
most benefits from Docker.

- Learn how to [build an image](../engine/reference/builder/){: target="_blank" rel="noopener" class="_"} using a Dockerfile
- Use [multi-stage builds](develop-images/multistage-build.md){: target="_blank" rel="noopener" class="_"} to keep your images lean
- Manage application data using [volumes](../storage/volumes.md) and [bind mounts](../storage/bind-mounts.md){: target="_blank" rel="noopener" class="_"}
- [Scale your app with Kubernetes](../get-started/kube-deploy.md){: target="_blank" rel="noopener" class="_"} 
- [Scale your app as a Swarm service](../get-started/swarm-deploy.md){: target="_blank" rel="noopener" class="_"} 
- [General application development best practices](dev-best-practices.md){: target="_blank" rel="noopener" class="_"}

## Learn about language-specific app development with Docker

- [Docker for Java developers lab](https://github.com/docker/labs/tree/master/developer-tools/java/){: target="_blank" rel="noopener" class="_"} 
- [Port a node.js app to Docker lab](https://github.com/docker/labs/tree/master/developer-tools/nodejs/porting){: target="_blank" rel="noopener" class="_"}
- [Ruby on Rails app on Docker lab](https://github.com/docker/labs/tree/master/developer-tools/ruby){: target="_blank" rel="noopener" class="_"}
- [Dockerize a .Net Core application](../samples/dotnetcore.md){: target="_blank" rel="noopener" class="_"}
- [Dockerize an ASP.NET Core application with SQL Server on Linux](../samples/aspnet-mssql-compose.md){: target="_blank" rel="noopener" class="_"} using Docker Compose

## Advanced development with the SDK or API

After you can write Dockerfiles or Compose files and use Docker CLI, take it to
the next level by using Docker Engine SDK for Go/Python or use the HTTP API
directly. Visit the [Develop with Docker Engine API](../engine/api/index.md)
section to learn more about developing with the Engine API, where to find SDKs
for your programming language of choice, and to see some examples.
