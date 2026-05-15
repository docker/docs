---
title: Testing an ASP.NET Core web app with Testcontainers
linkTitle: ASP.NET Core testing
description: Learn how to use Testcontainers for .NET to replace SQLite with a real Microsoft SQL Server in ASP.NET Core integration tests.
keywords: testcontainers, dotnet, csharp, testing, mssql, asp.net core, integration testing, entity framework
summary: |
  Learn how to test an ASP.NET Core web app using Testcontainers for .NET
  with a real Microsoft SQL Server instance instead of SQLite.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [c-sharp]
params:
  time: 25 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-testing-aspnet-core -->

In this guide, you'll learn how to:

- Use Testcontainers for .NET to spin up a Microsoft SQL Server container for integration tests
- Replace SQLite with a production-like database provider in ASP.NET Core tests
- Customize `WebApplicationFactory` to configure test dependencies with Testcontainers
- Manage container lifecycle with xUnit's `IAsyncLifetime`

## Prerequisites

- .NET 8.0+ SDK
- A code editor or IDE (Visual Studio, VS Code, Rider)
- A Docker environment supported by Testcontainers. For details, see the
  [Testcontainers .NET system requirements](https://dotnet.testcontainers.org/supported_docker_environment/).

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
