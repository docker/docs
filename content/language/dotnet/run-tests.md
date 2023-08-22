---
title: Run your tests
keywords: .NET, build, test
description: How to build and run your tests
---

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Use containers for development](develop.md).

## Introduction

Testing is an essential part of modern software development. In combination with 3rd party frameworks or services, Docker helps to test applications without mocks or complicated environment configurations fast and reliably.

## Add .NET test project

To test our sample application, we create a standalone test project from a template using the .NET CLI. On your local machine, open a terminal, change the directory to the `dotnet-docker` directory and run the following command:

```console
$ cd /path/to/dotnet-docker
$ dotnet new xunit -n myWebApp.Tests -o tests
```

Next, we'll update the test project and add the Testcontainers for .NET package that allows us to run tests against Docker resources (PostgreSQL container). Switch to the `tests` directory and run the following command:

```console
$ dotnet add package Testcontainers --version 2.3.0
```

## Add a test

Open the test project in your favorite IDE and replace the contents of `UnitTest1.cs` with the following code:

```c#
using System.Net;
using DotNet.Testcontainers.Builders;
using DotNet.Testcontainers.Containers;
using DotNet.Testcontainers.Networks;

public sealed class UnitTest1 : IAsyncLifetime, IDisposable
{
    private const ushort HttpPort = 80;

    private readonly CancellationTokenSource _cts = new(TimeSpan.FromMinutes(1));

    private readonly IDockerNetwork _network;

    private readonly IDockerContainer _dbContainer;

    private readonly IDockerContainer _appContainer;

    public UnitTest1()
    {
        _network = new TestcontainersNetworkBuilder()
            .WithName(Guid.NewGuid().ToString("D"))
            .Build();

        _dbContainer = new TestcontainersBuilder<TestcontainersContainer>()
            .WithImage("postgres")
            .WithNetwork(_network)
            .WithNetworkAliases("db")
            .WithVolumeMount("postgres-data", "/var/lib/postgresql/data")
            .Build();

        _appContainer = new TestcontainersBuilder<TestcontainersContainer>()
            .WithImage("dotnet-docker")
            .WithNetwork(_network)
            .WithPortBinding(HttpPort, true)
            .WithWaitStrategy(Wait.ForUnixContainer().UntilPortIsAvailable(HttpPort))
            .Build();
    }

    public async Task InitializeAsync()
    {
        await _network.CreateAsync(_cts.Token)
            .ConfigureAwait(false);

        await _dbContainer.StartAsync(_cts.Token)
            .ConfigureAwait(false);

        await _appContainer.StartAsync(_cts.Token)
            .ConfigureAwait(false);
    }

    public Task DisposeAsync()
    {
        return Task.CompletedTask;
    }

    public void Dispose()
    {
        _cts.Dispose();
    }

    [Fact]
    public async Task Test1()
    {
        using var httpClient = new HttpClient();
        httpClient.BaseAddress = new UriBuilder("http", _appContainer.Hostname, _appContainer.GetMappedPublicPort(HttpPort)).Uri;

        var httpResponseMessage = await httpClient.GetAsync(string.Empty)
            .ConfigureAwait(false);

        var body = await httpResponseMessage.Content.ReadAsStringAsync()
            .ConfigureAwait(false);

        Assert.Equal(HttpStatusCode.OK, httpResponseMessage.StatusCode);
        Assert.Contains("Welcome", body);
    }
}
```

The test class picks up the configurations and lessons we learned in the previous steps. It connects our application and database through a custom Docker network and runs an HTTP request against our application. As you can see, running containerized tests allows us to test applications without mocks or complicated environment configurations. The tests run on any Docker-API compatible environments including CI.

## Run the test

Before you run the test, [stop](run-containers.md#stop-start-and-name-containers) any running containers from the previous sections.

To run the test, change directory to the `dotnet-docker` directory and run the following `dotnet test` command:

```console
$ dotnet test tests
```

You should see output like the following:

```console
Starting test execution, please wait...
A total of 1 test files matched the specified pattern.

Passed!  - Failed:     0, Passed:     1, Skipped:     0, Total:     1, Duration: < 1 ms - myWebApp.Tests.dll (net6.0)
```

## Next steps

In the next module, we’ll take a look at how to set up a CI/CD pipeline using GitHub Actions. See:

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[dotnet%20docs%20feedback]) GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls) to suggest updates.
