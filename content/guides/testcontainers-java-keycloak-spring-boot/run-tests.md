---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based Spring Boot Keycloak integration tests and explore next steps.
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

You should see the Keycloak and PostgreSQL Docker containers start with the
realm settings imported and the tests pass. After the tests finish, the
containers stop and are removed automatically.

## Summary

The Testcontainers Keycloak module lets you develop and test applications using a
real Keycloak server instead of mocks. Testing against a real OAuth 2.0
provider that mirrors your production setup gives you more confidence in your
security configuration and token-based authentication flows.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Getting started with Testcontainers in a Java Spring Boot project](https://testcontainers.com/guides/testing-spring-boot-rest-api-using-testcontainers/)
- [Testcontainers Keycloak module](https://testcontainers.com/modules/keycloak/)
- [testcontainers-keycloak GitHub repository](https://github.com/dasniko/testcontainers-keycloak)
- [Spring Boot OAuth 2.0 Resource Server](https://docs.spring.io/spring-security/reference/servlet/oauth2/resource-server/index.html)
