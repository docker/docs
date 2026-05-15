---
title: Set up the project
linkTitle: Create the project
description: Set up an ASP.NET Core Razor Pages project with integration test dependencies.
weight: 10
---

## Background

This guide builds on top of Microsoft's
[Integration tests in ASP.NET Core](https://learn.microsoft.com/en-us/aspnet/core/test/integration-tests)
documentation. The original sample uses an in-memory SQLite database as the
backing store for integration tests. You'll replace SQLite with a real
Microsoft SQL Server instance running in a Docker container using
Testcontainers.

You can find the original code sample in the
[dotnet/AspNetCore.Docs.Samples](https://github.com/dotnet/AspNetCore.Docs.Samples/tree/main/test/integration-tests/IntegrationTestsSample)
repository.

## Clone the repository

Clone the Testcontainers guide repository and change into the project
directory:

```console
$ git clone https://github.com/testcontainers/tc-guide-testing-aspnet-core.git
$ cd tc-guide-testing-aspnet-core
```

## Project structure

The solution contains two projects:

```text
RazorPagesProject.sln
├── src/RazorPagesProject/              # ASP.NET Core Razor Pages app
└── tests/RazorPagesProject.Tests/      # xUnit integration tests
```

### Application project

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

### Test project

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

### Existing SQLite-based test factory

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
