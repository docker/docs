---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Replace SQLite with a real Microsoft SQL Server using Testcontainers for .NET.
weight: 20
---

The existing tests use an in-memory SQLite database. While convenient, this
doesn't match production behavior. You can replace it with a real Microsoft SQL
Server instance managed by Testcontainers.

## Add dependencies

Change to the test project directory and add the SQL Server Entity Framework
provider and the Testcontainers MSSQL module:

```console
$ cd tests/RazorPagesProject.Tests
$ dotnet add package Microsoft.EntityFrameworkCore.SqlServer --version 7.0.0
$ dotnet add package Testcontainers.MsSql --version 3.0.0
```

> [!NOTE]
> Testcontainers for .NET offers a range of
> [modules](https://www.nuget.org/profiles/Testcontainers) that follow best
> practice configurations.

## Create the test class

Create a `MsSqlTests.cs` file in the `IntegrationTests` directory. This class
manages the SQL Server container lifecycle and contains a nested test class.

```csharp
using System.Data.Common;
using System.Net;
using AngleSharp.Html.Dom;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.EntityFrameworkCore;
using RazorPagesProject.Data;
using RazorPagesProject.Tests.Helpers;
using Testcontainers.MsSql;
using Xunit;

namespace RazorPagesProject.Tests.IntegrationTests;

public sealed class MsSqlTests : IAsyncLifetime
{
    private readonly MsSqlContainer _msSqlContainer = new MsSqlBuilder().Build();

    public Task InitializeAsync()
    {
        return _msSqlContainer.StartAsync();
    }

    public Task DisposeAsync()
    {
        return _msSqlContainer.DisposeAsync().AsTask();
    }

    public sealed class IndexPageTests : IClassFixture<MsSqlTests>, IDisposable
    {
        private readonly WebApplicationFactory<Program> _webApplicationFactory;

        private readonly HttpClient _httpClient;

        public IndexPageTests(MsSqlTests fixture)
        {
            var clientOptions = new WebApplicationFactoryClientOptions();
            clientOptions.AllowAutoRedirect = false;

            _webApplicationFactory = new CustomWebApplicationFactory(fixture);
            _httpClient = _webApplicationFactory.CreateClient(clientOptions);
        }

        public void Dispose()
        {
            _webApplicationFactory.Dispose();
        }

        [Fact]
        public async Task Post_DeleteAllMessagesHandler_ReturnsRedirectToRoot()
        {
            // Arrange
            var defaultPage = await _httpClient.GetAsync("/")
                .ConfigureAwait(false);

            var document = await HtmlHelpers.GetDocumentAsync(defaultPage)
                .ConfigureAwait(false);

            // Act
            var form = (IHtmlFormElement)document.QuerySelector("form[id='messages']");
            var submitButton = (IHtmlButtonElement)document.QuerySelector("button[id='deleteAllBtn']");

            var response = await _httpClient.SendAsync(form, submitButton)
                .ConfigureAwait(false);

            // Assert
            Assert.Equal(HttpStatusCode.OK, defaultPage.StatusCode);
            Assert.Equal(HttpStatusCode.Redirect, response.StatusCode);
            Assert.Equal("/", response.Headers.Location.OriginalString);
        }

        private sealed class CustomWebApplicationFactory : WebApplicationFactory<Program>
        {
            private readonly string _connectionString;

            public CustomWebApplicationFactory(MsSqlTests fixture)
            {
                _connectionString = fixture._msSqlContainer.GetConnectionString();
            }

            protected override void ConfigureWebHost(IWebHostBuilder builder)
            {
                builder.ConfigureServices(services =>
                {
                    services.Remove(services.SingleOrDefault(service => typeof(DbContextOptions<ApplicationDbContext>) == service.ServiceType));
                    services.Remove(services.SingleOrDefault(service => typeof(DbConnection) == service.ServiceType));
                    services.AddDbContext<ApplicationDbContext>((_, option) => option.UseSqlServer(_connectionString));
                });
            }
        }
    }
}
```

## Understand the test structure

### Container lifecycle with IAsyncLifetime

The outer `MsSqlTests` class implements `IAsyncLifetime`. xUnit calls
`InitializeAsync()` right after creating the class instance, which starts the
SQL Server container. After all tests complete, `DisposeAsync()` stops and
removes the container.

```csharp
private readonly MsSqlContainer _msSqlContainer = new MsSqlBuilder().Build();
```

`MsSqlBuilder().Build()` creates a pre-configured Microsoft SQL Server
container. Testcontainers modules follow best practices, so you don't need
to configure ports, passwords, or startup wait strategies yourself.

### Nested test class with IClassFixture

The `IndexPageTests` class is nested inside `MsSqlTests` and implements
`IClassFixture<MsSqlTests>`. This gives the test class access to the
container's private field and creates a clean hierarchy in the test explorer.

### Custom WebApplicationFactory

Instead of using the SQLite-based factory, the nested
`CustomWebApplicationFactory` retrieves the connection string from the running
SQL Server container and passes it to `UseSqlServer()`:

```csharp
private sealed class CustomWebApplicationFactory : WebApplicationFactory<Program>
{
    private readonly string _connectionString;

    public CustomWebApplicationFactory(MsSqlTests fixture)
    {
        _connectionString = fixture._msSqlContainer.GetConnectionString();
    }

    protected override void ConfigureWebHost(IWebHostBuilder builder)
    {
        builder.ConfigureServices(services =>
        {
            services.Remove(services.SingleOrDefault(service => typeof(DbContextOptions<ApplicationDbContext>) == service.ServiceType));
            services.Remove(services.SingleOrDefault(service => typeof(DbConnection) == service.ServiceType));
            services.AddDbContext<ApplicationDbContext>((_, option) => option.UseSqlServer(_connectionString));
        });
    }
}
```

This factory:

1. Removes the existing `DbContextOptions<ApplicationDbContext>` registration
2. Removes the existing `DbConnection` registration
3. Adds a new `ApplicationDbContext` configured with the SQL Server connection
   string from the Testcontainers-managed container

> [!NOTE]
> The Microsoft SQL Server Docker image isn't compatible with ARM devices, such
> as Macs with Apple Silicon. You can use the
> [SqlEdge](https://www.nuget.org/packages/Testcontainers.SqlEdge) module or
> [Testcontainers Cloud](https://www.testcontainers.cloud/) as alternatives.
