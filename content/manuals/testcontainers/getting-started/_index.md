---
description: Learn how to install Testcontainers OSS and run your first container
keywords: testcontainers, testcontainers quickstart, testcontainers oss, testcontainers oss quickstart, testcontainers quickstart,
    java, go, golang
title: Testcontainers Quickstart
linkTitle: Testcontainers Quickstart
weight: 20
toc_max: 3
aliases:
- /testcontainers/getting-started/
---

This page contains summary information about Testcontainers OSS.

## Supported languages

Testcontainers provide support for the most popular languages, and Docker sponsors the development of the following Testcontainers implementations:

- [Go](/manuals/testcontainers/getting-started/go.md)
- [Java](/manuals/testcontainers/getting-started/java.md)

The rest are community-driven and maintained by independent contributors.

## Prerequisites

Testcontainers requires a Docker-API compatible container runtime. 
During development, Testcontainers is actively tested against recent versions of Docker on Linux, as well as against Docker Desktop on Mac and Windows. 
These Docker environments are automatically detected and used by Testcontainers without any additional configuration being necessary.

It is possible to configure Testcontainers to work for other Docker setups, such as a remote Docker host or Docker alternatives.
However, these are not actively tested in the main development workflow, so not all Testcontainers features might be available
and additional manual configuration might be necessary.

If you have further questions about configuration details for your setup or whether it supports running Testcontainers-based tests,
 contact the Testcontainers team and other users from the Testcontainers community on [Slack](https://slack.testcontainers.org/).
