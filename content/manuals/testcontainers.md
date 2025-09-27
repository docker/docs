---
title: Testcontainers
weight: 40
description: Learn how to use Testcontainers to run containers programmatically in your preferred programming language.
keywords: docker APIs, docker, testcontainers documentation, testcontainers, testcontainers oss, testcontainers oss documentation,
  docker compose, docker-compose, java, golang, go
params:
  sidebar:
    group: Open source
intro:
- title: What is Testcontainers?
  description: Learn about what Testcontainers does and its key benefits
  icon: feature_search
  link: https://testcontainers.com/getting-started/#what-is-testcontainers
- title: The Testcontainers workflow
  description: Understand the Testcontainers workflow
  icon: explore
  link: https://testcontainers.com/getting-started/#testcontainers-workflow
quickstart:
- title: Testcontainers for Go
  description: A Go package that makes it simple to create and clean up container-based dependencies for automated integration/smoke tests.
  icon: /icons/go.svg
  link: https://golang.testcontainers.org/quickstart/
- title: Testcontainers for Java
  description: A Java library that supports JUnit tests, providing lightweight, throwaway instances of anything that can run in a Docker container.
  icon: /icons/java.svg
  link: https://java.testcontainers.org/
---

Testcontainers is a set of open source libraries that provides easy and lightweight APIs for bootstrapping local development and test dependencies with real services wrapped in Docker containers.
Using Testcontainers, you can write tests that depend on the same services you use in production without mocks or in-memory services.

{{< grid items=intro >}}

## Quickstart

### Supported languages

Testcontainers provide support for the most popular languages, and Docker sponsors the development of the following Testcontainers implementations:

- [Go](https://golang.testcontainers.org/quickstart/)
- [Java](https://java.testcontainers.org/quickstart/junit_5_quickstart/)

The rest are community-driven and maintained by independent contributors.

### Prerequisites

Testcontainers requires a Docker-API compatible container runtime.
During development, Testcontainers is actively tested against recent versions of Docker on Linux, as well as against Docker Desktop on Mac and Windows.
These Docker environments are automatically detected and used by Testcontainers without any additional configuration being necessary.

It is possible to configure Testcontainers to work for other Docker setups, such as a remote Docker host or Docker alternatives.
However, these are not actively tested in the main development workflow, so not all Testcontainers features might be available
and additional manual configuration might be necessary.

If you have further questions about configuration details for your setup or whether it supports running Testcontainers-based tests,
 contact the Testcontainers team and other users from the Testcontainers community on [Slack](https://slack.testcontainers.org/).

 {{< grid items=quickstart >}}
