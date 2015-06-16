<!--[metadata]>
+++
aliases = [ "/docker-hub-enterprise/" ]
title = "Docker Trusted Registry: Overview"
description = "Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub,  registry"]
[menu.main]
parent="smn_dhe"
+++
<![end-metadata]-->


# Welcome to Docker Trusted Registry

Docker Trusted Registry (DTR) lets you run and manage your own Docker image
storage service, securely on your own infrastructure behind your company
firewall. This allows you to securely store, push, and pull the images used by
your enterprise to build, ship, and run applications. DTR also provides
monitoring and usage information to help you understand the workloads being
placed on it.

Specifically, DTR provides:

* An image registry to store, manage, and collaborate on Docker images
* Pluggable storage drivers
* Configuration options to let you run DTR in your particular enterprise
environment.
* Easy, transparent upgrades
* Logging, usage and system health metrics

DTR is perfect for:

* Providing a secure, on-premise development environment
* Creating a streamlined build pipeline
* Building a consistent, high-performance test/QA environment
* Managing image deployment

DTR is built on [version 2 of the Docker registry](https://github.com/docker/distribution).

> **Note:** This initial release of DTR has limited access. To get access,
> you will need an account on [Docker Hub](https://hub.docker.com/). Once you're
> logged in to the Hub with your account, visit the
> [early access registration page](https://registry.hub.docker.com/earlyaccess/)
> and follow the steps there to get signed up.

## Available Documentation

The following documentation for DTR is available:

* **Overview** This page.
* [**Quick Start: Basic User Workflow**](./quick-start.md) Go here to learn the
fundamentals of how DTR works and how you can set up a simple, but useful
workflow.
* [**User Guide**](./userguide.md) Go here to learn about using DTR from day to
day.
* [**Administrator Guide**](./adminguide.md) Go here if you are an administrator
responsible for running and maintaining DTR.
* [**Installation**](install.md) Go here for the steps you'll need to install
DTR and get it working.
* [**Configuration**](./configuration.md) Go here to find out details about
setting up and configuring DTR for your particular environment.
* [**Support**](./support.md) Go here for information on getting support for
DTR.

