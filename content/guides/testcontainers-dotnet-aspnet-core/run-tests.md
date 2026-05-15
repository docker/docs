---
title: Run tests and next steps
linkTitle: Run tests
description: Run the Testcontainers-based integration tests and explore next steps.
weight: 30
---

## Run the tests

Run the tests from the solution root:

```console
$ dotnet test ./RazorPagesProject.sln
```

The first run may take longer because Docker needs to pull the Microsoft SQL
Server image. On subsequent runs, the image is cached locally.

You should see xUnit discover and run the tests, including the
`MsSqlTests.IndexPageTests` class. Testcontainers starts a SQL Server
container, the tests execute against it, and the container is stopped and
removed automatically after the tests finish.

## Summary

By replacing SQLite with a Testcontainers-managed Microsoft SQL Server
instance, the integration tests run against the same type of database used in
production. This approach catches database-specific issues early, such as
differences in SQL dialect, transaction behavior, or data type handling between
SQLite and SQL Server.

The `MsSqlTests` class uses `IAsyncLifetime` to manage the container lifecycle,
and a nested `CustomWebApplicationFactory` wires the container's connection
string into the application's service configuration. You can apply this same
pattern to any database or service that Testcontainers supports.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers for .NET documentation](https://dotnet.testcontainers.org/)
- [Testcontainers for .NET modules](https://dotnet.testcontainers.org/modules/)
- [Microsoft SQL Server module](https://www.nuget.org/packages/Testcontainers.MsSql)
- [Integration tests in ASP.NET Core](https://learn.microsoft.com/en-us/aspnet/core/test/integration-tests)
