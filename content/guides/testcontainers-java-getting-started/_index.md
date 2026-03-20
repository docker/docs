---
title: Getting started with Testcontainers for Java
linkTitle: Testcontainers for Java
description: Learn how to use Testcontainers for Java to test a customer service with a real PostgreSQL database.
keywords: testcontainers, java, testing, postgresql, integration testing, junit, maven
summary: |
  Learn how to create a Java application and test database interactions
  using Testcontainers for Java with a real PostgreSQL instance.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 20 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-getting-started-with-testcontainers-for-java -->

In this guide, you will learn how to:

- Create a Java project with Maven
- Implement a `CustomerService` that manages customer records in PostgreSQL
- Write integration tests using Testcontainers with a real Postgres database
- Run the tests and verify everything works

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
