---
title: Testing Spring Boot Kafka Listener using Testcontainers
linkTitle: Spring Boot Kafka
description: Learn how to test a Spring Boot Kafka listener using Testcontainers with Kafka and MySQL modules.
keywords: testcontainers, java, spring boot, testing, kafka, mysql, jpa, awaitility
summary: |
  Learn how to create a Spring Boot application with a Kafka listener that persists data in MySQL,
  then test it using Testcontainers Kafka and MySQL modules with Awaitility.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 25 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-testing-spring-boot-kafka-listener -->

In this guide, you will learn how to:

- Create a Spring Boot application with Kafka integration
- Implement a Kafka listener and persist data in a MySQL database
- Test the Kafka listener using Testcontainers and Awaitility

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
