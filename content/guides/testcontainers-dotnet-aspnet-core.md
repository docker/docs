---
title: Testing an ASP.NET Core web app with Testcontainers
linkTitle: ASP.NET Core testing
description: Learn how to use Testcontainers for .NET to replace SQLite with a real Microsoft SQL Server in ASP.NET Core integration tests.
keywords: testcontainers, dotnet, csharp, testing, mssql, asp.net core, integration testing, entity framework
summary: |
  Learn how to test an ASP.NET Core web app using Testcontainers for .NET
  with a real Microsoft SQL Server instance instead of SQLite.
aliases:
  - /guides/testcontainers-dotnet-aspnet-core/create-project/
  - /guides/testcontainers-dotnet-aspnet-core/run-tests/
  - /guides/testcontainers-dotnet-aspnet-core/write-tests/
params:
  tags: [testing]
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

## Set up the project

### Background

This guide builds on top of Microsoft's
[Integration tests in ASP.NET Core](https://learn.microsoft.com/en-us/aspnet/core/test/integration-tests)
documentation. The original sample uses an in-memory SQLite database as the
backing store for integration tests. You'll replace SQLite with a real
Microsoft SQL Server instance running in a Docker container using
Testcontainers.

You can find the original code sample in the
[dotnet/AspNetCore.Docs.Samples](https://github.com/dotnet/AspNetCore.Docs.Samples/tree/main/test/integration-tests/IntegrationTestsSample)
repository.

### Clone the repository

Clone the Testcontainers guide repository and change into the project
directory:

```console
$ git clone https://github.com/testcontainers/tc-guide-testing-aspnet-core.git
$ cd tc-guide-testing-aspnet-core
```

### Project structure

The solution contains two projects:

```text
RazorPagesProject.sln
├── src/RazorPagesProject/              # ASP.NET Core Razor Pages app
└── tests/RazorPagesProject.Tests/      # xUnit integration tests
```

#### Application project

The application project (`src/RazorPagesProject/RazorPagesProject.csproj`)
is a Razor Pages web app that uses Entity Framework Core with SQLite as its
default database provider:

```xml
<Project Sdk="Microsoft.NET.Sdk.Web">

  <PropertyGroup>
    <TargetFramework>net9.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="Microsoft.EntityFrameworkCore.Sqlite" Version="7.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Diagnostics.EntityFrameworkCore" Version="7.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Identity.EntityFrameworkCore" Version="7.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Identity.UI" Version="7.0.0" />
    <PackageReference Include="Microsoft.EntityFrameworkCore.Tools" Version="7.0.0">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers; buildtransitive</IncludeAssets>
    </PackageReference>
  </ItemGroup>

</Project>
```

The `ApplicationDbContext` stores `Message` entities and provides methods to
query and manage them:

```csharp
public class ApplicationDbContext : IdentityDbContext
{
    public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
        : base(options)
    {
    }

    public virtual DbSet<Message> Messages { get; set; }

    public async virtual Task<List<Message>> GetMessagesAsync()
    {
        return await Messages
            .OrderBy(message => message.Text)
            .AsNoTracking()
            .ToListAsync();
    }

    public async virtual Task AddMessageAsync(Message message)
    {
        await Messages.AddAsync(message);
        await SaveChangesAsync();
    }

    public async virtual Task DeleteAllMessagesAsync()
    {
        foreach (Message message in Messages)
        {
            Messages.Remove(message);
        }

        await SaveChangesAsync();
    }

    public async virtual Task DeleteMessageAsync(int id)
    {
        var message = await Messages.FindAsync(id);

        if (message != null)
        {
            Messages.Remove(message);
            await SaveChangesAsync();
        }
    }

    public void Initialize()
    {
        Messages.AddRange(GetSeedingMessages());
        SaveChanges();
    }

    public static List<Message> GetSeedingMessages()
    {
        return new List<Message>()
        {
            new Message(){ Text = "You're standing on my scarf." },
            new Message(){ Text = "Would you like a jelly baby?" },
            new Message(){ Text = "To the rational mind, nothing is inexplicable; only unexplained." }
        };
    }
}
```

#### Test project

The test project (`tests/RazorPagesProject.Tests/RazorPagesProject.Tests.csproj`)
includes xUnit, the ASP.NET Core testing infrastructure, and the
Testcontainers MSSQL module:

```xml
<Project Sdk="Microsoft.NET.Sdk.Web">

  <PropertyGroup>
    <TargetFramework>net9.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="AngleSharp" Version="0.17.1" />
    <PackageReference Include="Microsoft.AspNetCore.Diagnostics.EntityFrameworkCore" Version="7.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Identity.EntityFrameworkCore" Version="7.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Identity.UI" Version="7.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Mvc.Testing" Version="7.0.0" />
    <PackageReference Include="Microsoft.EntityFrameworkCore" Version="7.0.0" />
    <PackageReference Include="Microsoft.EntityFrameworkCore.Sqlite" Version="7.0.0" />
    <PackageReference Include="Microsoft.EntityFrameworkCore.SqlServer" Version="7.0.0" />
    <PackageReference Include="Microsoft.EntityFrameworkCore.Tools" Version="7.0.0">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers; buildtransitive</IncludeAssets>
    </PackageReference>

    <PackageReference Include="Microsoft.NET.Test.Sdk" Version="17.4.0" />

    <PackageReference Include="Testcontainers.MsSql" Version="3.0.0" />
    <PackageReference Include="xunit" Version="2.4.2" />
    <PackageReference Include="xunit.runner.visualstudio" Version="2.4.5">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers; buildtransitive</IncludeAssets>
    </PackageReference>
  </ItemGroup>

  <ItemGroup>
    <ProjectReference Include="..\..\src\RazorPagesProject\RazorPagesProject.csproj" />
  </ItemGroup>

  <ItemGroup>
    <Content Update="xunit.runner.json">
      <CopyToOutputDirectory>Always</CopyToOutputDirectory>
    </Content>
  </ItemGroup>

</Project>
```

The key dependencies are:

- `Microsoft.AspNetCore.Mvc.Testing` - provides `WebApplicationFactory` for
  bootstrapping the app in tests
- `Microsoft.EntityFrameworkCore.SqlServer` - the SQL Server database provider
  for Entity Framework Core
- `Testcontainers.MsSql` - the Testcontainers module for Microsoft SQL Server

#### Existing SQLite-based test factory

The original project includes a `CustomWebApplicationFactory` that replaces
the application's database with an in-memory SQLite instance:

```csharp
public class CustomWebApplicationFactory<TProgram>
    : WebApplicationFactory<TProgram> where TProgram : class
{
    protected override void ConfigureWebHost(IWebHostBuilder builder)
    {
        builder.ConfigureServices(services =>
        {
            var dbContextDescriptor = services.SingleOrDefault(
                d => d.ServiceType ==
                    typeof(DbContextOptions<ApplicationDbContext>));

            services.Remove(dbContextDescriptor);

            var dbConnectionDescriptor = services.SingleOrDefault(
                d => d.ServiceType ==
                    typeof(DbConnection));

            services.Remove(dbConnectionDescriptor);

            // Create open SqliteConnection so EF won't automatically close it.
            services.AddSingleton<DbConnection>(container =>
            {
                var connection = new SqliteConnection("DataSource=:memory:");
                connection.Open();

                return connection;
            });

            services.AddDbContext<ApplicationDbContext>((container, options) =>
            {
                var connection = container.GetRequiredService<DbConnection>();
                options.UseSqlite(connection);
            });
        });

        builder.UseEnvironment("Development");
    }
}
```

While this approach works, SQLite has behavioral differences from the database
you'd use in production. In the next section, you'll replace it with a
Testcontainers-managed Microsoft SQL Server instance.

## Write tests with Testcontainers

The existing tests use an in-memory SQLite database. While convenient, this
doesn't match production behavior. You can replace it with a real Microsoft SQL
Server instance managed by Testcontainers.

### Add dependencies

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

### Create the test class

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

### Understand the test structure

#### Container lifecycle with IAsyncLifetime

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

#### Nested test class with IClassFixture

The `IndexPageTests` class is nested inside `MsSqlTests` and implements
`IClassFixture<MsSqlTests>`. This gives the test class access to the
container's private field and creates a clean hierarchy in the test explorer.

#### Custom WebApplicationFactory

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

## Run tests and next steps

### Run the tests

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

### Summary

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

### Further reading

- [Testcontainers for .NET documentation](https://dotnet.testcontainers.org/)
- [Testcontainers for .NET modules](https://dotnet.testcontainers.org/modules/)
- [Microsoft SQL Server module](https://www.nuget.org/packages/Testcontainers.MsSql)
- [Integration tests in ASP.NET Core](https://learn.microsoft.com/en-us/aspnet/core/test/integration-tests)
