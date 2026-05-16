---
title: Working with jOOQ and Flyway using Testcontainers
linkTitle: jOOQ and Flyway
description: Learn how to generate jOOQ code from a database using Testcontainers and Flyway, then test your persistence layer.
keywords: testcontainers, java, testing, jooq, flyway, postgresql, spring boot, code generation
summary: |
  Generate typesafe jOOQ code from a real PostgreSQL database managed by
  Flyway migrations, then test repositories using Testcontainers.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 25 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-working-with-jooq-flyway-using-testcontainers -->

In this guide, you will learn how to:

- Create a Spring Boot application with jOOQ support
- Generate jOOQ code using Testcontainers, Flyway, and a Maven plugin
- Implement basic database operations using jOOQ
- Load complex object graphs using jOOQ's MULTISET feature
- Test the jOOQ persistence layer using Testcontainers

## Prerequisites

- Java 17+
- Maven
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
