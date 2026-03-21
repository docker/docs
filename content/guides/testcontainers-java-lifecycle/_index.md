---
title: Testcontainers container lifecycle management using JUnit 5
linkTitle: Container lifecycle (Java)
description: Learn how to manage Testcontainers container lifecycle using JUnit 5 callbacks, extension annotations, and the singleton containers pattern.
keywords: testcontainers, java, testing, junit, lifecycle, singleton containers, postgresql
summary: |
  Learn different approaches to manage container lifecycle with Testcontainers
  using JUnit 5 lifecycle callbacks, extension annotations, and the singleton
  containers pattern.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 20 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-testcontainers-lifecycle -->

In this guide, you will learn how to:

- Start and stop containers using JUnit 5 lifecycle callbacks
- Manage containers using JUnit 5 extension annotations (`@Testcontainers` and `@Container`)
- Share containers across multiple test classes using the singleton containers pattern
- Avoid a common misconfiguration when combining extension annotations with singleton containers

## Prerequisites

- Java 17+
- Your preferred IDE
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
