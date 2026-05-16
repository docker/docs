---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Write integration tests using Testcontainers for .NET and xUnit with a real PostgreSQL database.
weight: 20
---

## Add Testcontainers dependencies

Add the Testcontainers PostgreSQL module to the test project:

```console
$ dotnet add ./CustomerService.Tests/CustomerService.Tests.csproj package Testcontainers.PostgreSql
```

## Write the test

Create `CustomerServiceTest.cs` in the test project:

```csharp
using Testcontainers.PostgreSql;

namespace Customers.Tests;

public sealed class CustomerServiceTest : IAsyncLifetime
{
    private readonly PostgreSqlContainer _postgres = new PostgreSqlBuilder()
        .WithImage("postgres:16-alpine")
        .Build();

    public Task InitializeAsync()
    {
        return _postgres.StartAsync();
    }

    public Task DisposeAsync()
    {
        return _postgres.DisposeAsync().AsTask();
    }

    [Fact]
    public void ShouldReturnTwoCustomers()
    {
        // Given
        var customerService = new CustomerService(new DbConnectionProvider(_postgres.GetConnectionString()));

        // When
        customerService.Create(new Customer(1, "George"));
        customerService.Create(new Customer(2, "John"));
        var customers = customerService.GetCustomers();

        // Then
        Assert.Equal(2, customers.Count());
    }
}
```

Here's what the test does:

- Declares a `PostgreSqlContainer` using the `PostgreSqlBuilder` with the
  `postgres:16-alpine` Docker image.
- Implements `IAsyncLifetime` for container lifecycle management:
  - `InitializeAsync()` starts the container before the test runs.
  - `DisposeAsync()` stops and removes the container after the test finishes.
- `ShouldReturnTwoCustomers()` creates a `CustomerService` with connection
  details from the container, inserts two customers, fetches all customers, and
  asserts the count.
