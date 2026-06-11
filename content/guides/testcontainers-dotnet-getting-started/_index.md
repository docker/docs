---
title: Getting started with Testcontainers for .NET
linkTitle: Testcontainers for .NET
description: Learn how to use Testcontainers for .NET to test database interactions with a real PostgreSQL instance.
keywords: testcontainers, dotnet, csharp, testing, postgresql, integration testing, xunit
summary: |
  Learn how to create a .NET application and test database interactions
  using Testcontainers for .NET with a real PostgreSQL instance.
aliases:
  - /guides/testcontainers-dotnet-getting-started/create-project/
  - /guides/testcontainers-dotnet-getting-started/run-tests/
  - /guides/testcontainers-dotnet-getting-started/write-tests/
params:
  tags: [testing]
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

## Create the .NET project

### Set up the solution

Create a .NET solution with source and test projects:

```console
$ dotnet new sln -o TestcontainersDemo
$ cd TestcontainersDemo
$ dotnet new classlib -o CustomerService
$ dotnet sln add ./CustomerService/CustomerService.csproj
$ dotnet new xunit -o CustomerService.Tests
$ dotnet sln add ./CustomerService.Tests/CustomerService.Tests.csproj
$ dotnet add ./CustomerService.Tests/CustomerService.Tests.csproj reference ./CustomerService/CustomerService.csproj
```

Add the Npgsql dependency to the source project:

```console
$ dotnet add ./CustomerService/CustomerService.csproj package Npgsql
```

### Implement the business logic

Create a `Customer` record type:

```csharp
namespace Customers;

public readonly record struct Customer(long Id, string Name);
```

Create a `DbConnectionProvider` class to manage database connections:

```csharp
using System.Data.Common;
using Npgsql;

namespace Customers;

public sealed class DbConnectionProvider
{
    private readonly string _connectionString;

    public DbConnectionProvider(string connectionString)
    {
        _connectionString = connectionString;
    }

    public DbConnection GetConnection()
    {
        return new NpgsqlConnection(_connectionString);
    }
}
```

Create the `CustomerService` class:

```csharp
namespace Customers;

public sealed class CustomerService
{
    private readonly DbConnectionProvider _dbConnectionProvider;

    public CustomerService(DbConnectionProvider dbConnectionProvider)
    {
        _dbConnectionProvider = dbConnectionProvider;
        CreateCustomersTable();
    }

    public IEnumerable<Customer> GetCustomers()
    {
        IList<Customer> customers = new List<Customer>();

        using var connection = _dbConnectionProvider.GetConnection();
        using var command = connection.CreateCommand();
        command.CommandText = "SELECT id, name FROM customers";
        command.Connection?.Open();

        using var dataReader = command.ExecuteReader();
        while (dataReader.Read())
        {
            var id = dataReader.GetInt64(0);
            var name = dataReader.GetString(1);
            customers.Add(new Customer(id, name));
        }

        return customers;
    }

    public void Create(Customer customer)
    {
        using var connection = _dbConnectionProvider.GetConnection();
        using var command = connection.CreateCommand();

        var id = command.CreateParameter();
        id.ParameterName = "@id";
        id.Value = customer.Id;

        var name = command.CreateParameter();
        name.ParameterName = "@name";
        name.Value = customer.Name;

        command.CommandText = "INSERT INTO customers (id, name) VALUES(@id, @name)";
        command.Parameters.Add(id);
        command.Parameters.Add(name);
        command.Connection?.Open();
        command.ExecuteNonQuery();
    }

    private void CreateCustomersTable()
    {
        using var connection = _dbConnectionProvider.GetConnection();
        using var command = connection.CreateCommand();
        command.CommandText = "CREATE TABLE IF NOT EXISTS customers (id BIGINT NOT NULL, name VARCHAR NOT NULL, PRIMARY KEY (id))";
        command.Connection?.Open();
        command.ExecuteNonQuery();
    }
}
```

Here's what `CustomerService` does:

- The constructor calls `CreateCustomersTable()` to ensure the table exists.
- `GetCustomers()` fetches all rows from the `customers` table and returns them as `Customer` objects.
- `Create()` inserts a customer record into the database.

## Write tests with Testcontainers

### Add Testcontainers dependencies

Add the Testcontainers PostgreSQL module to the test project:

```console
$ dotnet add ./CustomerService.Tests/CustomerService.Tests.csproj package Testcontainers.PostgreSql
```

### Write the test

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

## Run tests and next steps

### Run the tests

Run the tests:

```console
$ dotnet test
```

You can see in the output that Testcontainers pulls the Postgres Docker image
from Docker Hub (if not already available locally), starts the container, and
runs the test.

Writing an integration test using Testcontainers works like writing a unit test
that you can run from your IDE. Your teammates can clone the project and run
tests without installing Postgres on their machines.

### Summary

The Testcontainers for .NET library helps you write integration tests using the
same type of database (Postgres) that you use in production, instead of mocks.
Because you aren't using mocks and instead talk to real services, you're free
to refactor code and still verify that the application works as expected.

In addition to Postgres, Testcontainers provides dedicated
[modules](https://www.nuget.org/profiles/Testcontainers) for many SQL
databases, NoSQL databases, messaging queues, and more.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testing an ASP.NET Core web app](https://testcontainers.com/guides/testing-an-aspnet-core-web-app/)
