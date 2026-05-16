---
title: Testing a Spring Boot REST API with Testcontainers
linkTitle: Spring Boot REST API
description: Learn how to test a Spring Boot REST API using Testcontainers with PostgreSQL and REST Assured.
keywords: testcontainers, java, spring boot, testing, postgresql, rest api, rest assured, jpa
summary: |
  Learn how to create a Spring Boot REST API with Spring Data JPA and PostgreSQL,
  then test it using Testcontainers and REST Assured.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 25 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-testing-spring-boot-rest-api -->

In this guide, you will learn how to:

- Create a Spring Boot application with a REST API endpoint
- Use Spring Data JPA with PostgreSQL to store and retrieve data
- Test the REST API using Testcontainers and REST Assured

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
