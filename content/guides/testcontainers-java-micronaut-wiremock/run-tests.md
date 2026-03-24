---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers WireMock integration tests and explore next steps.
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

You should see the WireMock Docker container start in the console output. It
acts as the photo service, serving mock responses based on the configured
expectations. All tests should pass.

## Summary

You built a Micronaut application that integrates with an external REST API
using declarative HTTP clients, then tested that integration using WireMock and
the Testcontainers WireMock module. Testing at the HTTP protocol level instead
of mocking Java methods lets you catch serialization issues and simulate
realistic failure scenarios.

> [!TIP]
> Testcontainers WireMock modules are available for Go and Python as well.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers WireMock module](https://testcontainers.com/modules/wiremock/)
- [WireMock documentation](https://wiremock.org/docs/)
- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Testing REST API integrations in Spring Boot using WireMock](/guides/testcontainers-java-wiremock/)
