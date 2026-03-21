---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based Micronaut Kafka integration tests and explore next steps.
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

You should see the Kafka and MySQL Docker containers start and all tests pass.
After the tests finish, the containers stop and are removed automatically.

## Summary

Testing with real Kafka and MySQL instances gives you more confidence in the
correctness of your code than mocks or embedded alternatives. The
Testcontainers library manages the container lifecycle so that your integration
tests run against the same services you use in production.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testing REST API integrations in Micronaut apps using WireMock](/guides/testcontainers-java-micronaut-wiremock/)
- [Testing Spring Boot Kafka Listener using Testcontainers](/guides/testcontainers-java-spring-boot-kafka/)
- [Getting started with Testcontainers in a Java Spring Boot project](https://testcontainers.com/guides/testing-spring-boot-rest-api-using-testcontainers/)
- [Awaitility](http://www.awaitility.org/)
- [Testcontainers Kafka module](https://java.testcontainers.org/modules/kafka/)
- [Testcontainers MySQL module](https://java.testcontainers.org/modules/databases/mysql/)
