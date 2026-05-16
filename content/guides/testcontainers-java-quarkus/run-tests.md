---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based Quarkus integration tests and explore next steps.
weight: 30
---

## Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the PostgreSQL Docker container start and all tests pass. After
the tests finish, the container stops and is removed automatically.

## Run the application locally

Quarkus Dev Services automatically provisions unconfigured services in
development mode. Start the Quarkus application in dev mode:

```console
$ ./mvnw compile quarkus:dev
```

Or with Gradle:

```console
$ ./gradlew quarkusDev
```

Dev Services starts a PostgreSQL container automatically. If you're running a
PostgreSQL database on your system and want to use that instead, configure the
datasource properties in `src/main/resources/application.properties`:

```properties
quarkus.datasource.jdbc.url=jdbc:postgresql://localhost:5432/postgres
quarkus.datasource.username=postgres
quarkus.datasource.password=postgres
```

When these properties are set explicitly, Dev Services doesn't provision the
database container and instead connects to the configured database.

## Summary

Quarkus Dev Services improves the developer experience by automatically
provisioning the required services using Testcontainers during development and
testing. This guide covered:

- Building a REST API using JAX-RS with Hibernate ORM with Panache
- Testing API endpoints using REST Assured with Dev Services handling database
  provisioning
- Using `QuarkusTestResourceLifecycleManager` for services not supported by Dev
  Services
- Running the application locally with Dev Services

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Quarkus Dev Services overview](https://quarkus.io/guides/dev-services)
- [Quarkus testing guide](https://quarkus.io/guides/getting-started-testing)
- [Testcontainers Postgres module](https://java.testcontainers.org/modules/databases/postgres/)
