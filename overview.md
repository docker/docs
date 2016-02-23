
+++
aliases = [ "/docker-hub-enterprise/" ]
title = "Overview Trusted Registry"
description = "Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub,  registry"]
[menu.main]
parent="workw_dtr"
weight=-99
+++

# Overview of Docker Trusted Registry

Docker Trusted Registry  lets you run and manage your own Docker image
storage service, securely on your own infrastructure behind your company
firewall. This allows you to securely store, push, and pull the images used by
your enterprise to build, ship, and run applications. Docker Trusted Registry also provides
monitoring and usage information to help you understand the workloads being
placed on it.

Specifically, Docker Trusted Registry provides:

* An image registry to store, manage, and collaborate on Docker images
* Pluggable storage drivers
* Configuration options to let you run Docker Trusted Registry in your particular enterprise
environment.
* Easy, transparent upgrades
* Logging, usage and system health metrics

Docker Trusted Registry is perfect for:

* Providing a secure, on-premises development environment
* Creating a streamlined build pipeline
* Building a consistent, high-performance test/QA environment
* Managing image deployment

Docker Trusted Registry is built on [version 2 of the Docker registry](https://github.com/docker/distribution).

To get your copy of Docker Trusted Registry, including a free trial, visit [the Docker Subscription page](https://hub.docker.com/enterprise/). For more information on acquiring Docker Trusted Registry, see the [install page](install/index.md).

>   **Important**: Docker Trusted Registry must be used with the current version of the commercially
>   supported Docker Engine. You must install this version of Docker before
>   installing Docker Trusted Registry. For instructions on accessing and installing commercially
>   supported Docker Engine, visit the [install page](install/index.md#download-the-commercially-supported-docker-engine-installation-script).

## Available Documentation

The following documentation for Docker Trusted Registry is available:

* **Overview**, this page.
* [**Release Notes**](release-notes.md) See the latest additions, fixes, and known issues.
* [**Quick Start: Basic User Workflow**](quick-start.md) Go here to learn the
fundamentals of how Docker Trusted Registry works and how you can set up a simple, but useful
workflow.
* [**User Guide**](userguide.md) Go here to learn about using Docker Trusted Registry from day to
day.
* [**Administrator Guide**](adminguide.md) Go here if you are an administrator
responsible for running and maintaining Docker Trusted Registry.
* [**Installation**](install/index.md) Go here for the steps you'll need to install
Docker Trusted Registry and get it working.
* [**Configuration**](configure/configuration.md) Go here to find out details about
setting up and configuring Docker Trusted Registry for your particular environment.
* [**Support**](install/index.md) Go here for information on getting support for Docker Trusted Registry.
* [**The Docker Trusted Registry product page**](https://www.docker.com/docker-trusted-registry).
* [**Docker Trusted Registry Use Cases page**](https://www.docker.com/products/use-cases) showing an example CI/CD pipeline.
* [**Docker Trusted Registry and Docker tutorials and webinars**](https://www.docker.com/products/resources).

Note: Docker Trusted Registry requires that you use the commercially supported Docker Engine.
