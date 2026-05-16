---
title: Getting started with Testcontainers for .NET
linkTitle: Testcontainers for .NET
description: Learn how to use Testcontainers for .NET to test database interactions with a real PostgreSQL instance.
keywords: testcontainers, dotnet, csharp, testing, postgresql, integration testing, xunit
summary: |
  Learn how to create a .NET application and test database interactions
  using Testcontainers for .NET with a real PostgreSQL instance.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [c-sharp]
params:
  time: 20 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-getting-started-with-testcontainers-for-dotnet -->

In this guide, you will learn how to:

- Create a .NET solution with a source and test project
- Implement a `CustomerService` that manages customer records in PostgreSQL
- Write integration tests using Testcontainers and xUnit
- Manage container lifecycle with `IAsyncLifetime`

## Prerequisites

- .NET 8.0+ SDK
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
