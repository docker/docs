---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based integration tests and explore next steps.
weight: 30
---

## Run the tests

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

## Summary

The Testcontainers for .NET library helps you write integration tests using the
same type of database (Postgres) that you use in production, instead of mocks.
Because you aren't using mocks and instead talk to real services, you're free
to refactor code and still verify that the application works as expected.

In addition to Postgres, Testcontainers provides dedicated
[modules](https://www.nuget.org/profiles/Testcontainers) for many SQL
databases, NoSQL databases, messaging queues, and more.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testing an ASP.NET Core web app](https://testcontainers.com/guides/testing-an-aspnet-core-web-app/)
