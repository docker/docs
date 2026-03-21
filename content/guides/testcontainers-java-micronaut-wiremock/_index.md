---
title: Testing REST API integrations in Micronaut apps using WireMock
linkTitle: Micronaut WireMock
description: Learn how to test REST API integrations in a Micronaut application using the Testcontainers WireMock module.
keywords: testcontainers, java, micronaut, testing, wiremock, rest api
summary: |
  Learn how to create a Micronaut application that integrates with
  external REST APIs, then test those integrations using WireMock
  and the Testcontainers WireMock module.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 20 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-in-micronaut-apps-using-wiremock -->

In this guide, you'll learn how to:

- Create a Micronaut application that talks to external REST APIs
- Test external API integrations using WireMock
- Use the Testcontainers WireMock module to run WireMock as a Docker container

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
