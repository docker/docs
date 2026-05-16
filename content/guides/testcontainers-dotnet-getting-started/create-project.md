---
title: Create the .NET project
linkTitle: Create the project
description: Set up a .NET solution with a PostgreSQL-backed customer service.
weight: 10
---

## Set up the solution

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

## Implement the business logic

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
