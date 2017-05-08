---
description: Summary of samples related to Compose
keywords: documentation, docs, docker, compose, samples
title: Samples for Compose
notoc: true
---

The following samples show the various aspects of how to work with Docker Compose. As a prerequisite, be sure to [install Docker Compose](/compose/install/) if you have not already done so.

The samples should help you to:

- define services based on Docker images in [Compose files](/compose/compose-file.md) `docker-compose.yml` and
`docker-stack.yml` files
- understand the relationship between `docker-compose.yml` and [Dockerfiles](/engine/reference/builder.md)
- learn how to make calls to your application services from Compose files
- learn how to deploy applications and services to a [swarm](/engine/swarm.md)

These samples focus specifically on Docker Compose:

- [Quickstart: Compose and Django](/compose/django.md) - This quickstart
guide demonstrates how to use Docker Compose to set up and run a simple
Django/PostgreSQL app.

- [Quickstart: Compose and Rails](/compose/rails.md) - This quickstart
guide will show you how to use Docker Compose to set up and run a Rails/PostgreSQL
app.

- [Quickstart: Compose and WordPress](/compose/wordpress.md)

These samples include working with Docker Compose as part of broader learning goals:

- [Get Started with Docker](/get-started/index.md) - This multi-part tutorial covers writing your first app, data storage, networking, and swarms,
and ends with your app running on production servers in the cloud.

- [Deploying an app to a Swarm](https://github.com/docker/labs/blob/master/beginner/chapters/votingapp.md) - This tutorial from [Docker Labs](https://github.com/docker/labs/blob/master/README.md) shows you how to create and customize a sample voting app, deploy it to a [swarm](/engine/swarm.md), test it, reconfigure the app, and redeploy.
