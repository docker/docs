---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based Spring Boot integration tests and explore next steps.
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

You should see the Postgres Docker container start and all tests pass. After
the tests finish, the container stops and is removed automatically.

## Summary

The Testcontainers library helps you write integration tests by using the same
type of database (Postgres) that you use in production, instead of mocks or
in-memory databases. Because you test against real services, you're free to
refactor code and still verify that the application works as expected.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Testcontainers Postgres module](https://java.testcontainers.org/modules/databases/postgres/)
- [Testcontainers JDBC support](https://java.testcontainers.org/modules/databases/jdbc/)
