---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based Spring Cloud AWS integration tests and explore next steps.
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

You should see the LocalStack Docker container start and the test pass. After
the tests finish, the container stops and is removed automatically.

## Summary

LocalStack lets you develop and test AWS-based applications locally.
The Testcontainers LocalStack module makes it straightforward to write
integration tests by using ephemeral LocalStack containers that start on random
ports with no external setup required.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers LocalStack module](https://java.testcontainers.org/modules/localstack/)
- [Getting started with Testcontainers for Java](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Spring Cloud AWS documentation](https://docs.awspring.io/spring-cloud-aws/docs/3.0.3/reference/html/index.html)
