---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers MockServer integration tests and explore next steps.
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

You should see the MockServer Docker container start in the console output. It
acts as the photo service, serving mock responses based on the configured
expectations. All tests should pass.

## Summary

You built a Spring Boot application that integrates with an external REST API
using declarative HTTP clients, then tested that integration using the
Testcontainers MockServer module. Testing at the HTTP protocol level instead of
mocking Java methods lets you catch serialization issues and simulate realistic
failure scenarios.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers MockServer module](https://java.testcontainers.org/modules/mockserver/)
- [MockServer documentation](https://www.mock-server.com/)
- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
