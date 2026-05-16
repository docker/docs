---
title: Getting started with Testcontainers for Node.js
linkTitle: Testcontainers for Node.js
description: Learn how to use Testcontainers for Node.js to test database interactions with a real PostgreSQL instance.
keywords: testcontainers, nodejs, javascript, testing, postgresql, integration testing, jest
summary: |
  Learn how to create a Node.js application and test database interactions
  using Testcontainers for Node.js with a real PostgreSQL instance.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [js]
params:
  time: 15 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-getting-started-with-testcontainers-for-nodejs -->

In this guide, you will learn how to:

- Create a Node.js application that stores and retrieves customers from PostgreSQL
- Write integration tests using Testcontainers and Jest
- Run tests against a real PostgreSQL database in a Docker container

## Prerequisites

- Node.js 18+
- npm
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
