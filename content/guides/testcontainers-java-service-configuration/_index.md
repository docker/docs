---
title: Configuration of services running in a container
linkTitle: Service configuration (Java)
description: Learn how to configure services running in Testcontainers by copying files and executing commands inside containers.
keywords: testcontainers, java, testing, postgresql, localstack, container configuration
summary: |
  Learn how to initialize and configure Docker containers for testing
  by copying files into containers and executing commands inside them.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 15 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-configuration-of-services-running-in-container -->

In this guide, you will learn how to:

- Initialize containers by copying files into them
- Run commands inside running containers using `execInContainer()`
- Set up a PostgreSQL database with SQL scripts
- Create AWS S3 buckets in LocalStack containers

## Prerequisites

- Java 17+
- Your preferred IDE
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
