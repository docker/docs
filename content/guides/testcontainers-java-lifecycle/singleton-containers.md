---
title: Singleton containers pattern
linkTitle: Singleton containers
description: Share containers across multiple test classes using the singleton containers pattern.
weight: 40
---

As the number of test classes grows, starting containers for each class adds
up. The singleton containers pattern starts all required containers once in a
common base class and reuses them across all integration tests.

## Define the base class

Create an abstract base class that starts the containers in a static
initializer:

```java
package com.testcontainers.demo;

import org.testcontainers.postgresql.PostgreSQLContainer;
import org.testcontainers.kafka.ConfluentKafkaContainer;

public abstract class AbstractIntegrationTest {

   static PostgreSQLContainer postgres = new PostgreSQLContainer(
           "postgres:16-alpine");
   static ConfluentKafkaContainer kafka = new ConfluentKafkaContainer(
           "confluentinc/cp-kafka:7.8.0");

   static {
       postgres.start();
       kafka.start();
   }
}
```

The containers start once when the class loads and Testcontainers uses the
[Ryuk container](https://github.com/testcontainers/moby-ryuk) to remove them
after the JVM exits.

> [!TIP]
> Instead of starting containers sequentially, start them in parallel using
> `Startables.deepStart(postgres, kafka).join();`

## Extend the base class

Each test class inherits from the base class and reuses the same containers:

```java
class ProductControllerTest extends AbstractIntegrationTest {

   ProductRepository productRepository;

   @BeforeEach
   void setUp() {
       productRepository = new ProductRepository(...);
       productRepository.deleteAll();
   }

   @Test
   void shouldGetAllProducts() {
       // test logic using the shared postgres container
   }
}
```

## Avoid a common misconfiguration

A common mistake is combining singleton containers with the `@Testcontainers`
and `@Container` annotations:

```java
// DON'T DO THIS — containers will stop after each test class
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@Testcontainers
public abstract class AbstractIntegrationTest {

   @Container
   static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>(
           DockerImageName.parse("postgres:16-alpine"));

   @DynamicPropertySource
   static void configureProperties(DynamicPropertyRegistry registry) {
       registry.add("spring.datasource.url", postgres::getJdbcUrl);
       registry.add("spring.datasource.username", postgres::getUsername);
       registry.add("spring.datasource.password", postgres::getPassword);
   }
}
```

The `@Testcontainers` extension stops containers at the end of **each test
class**. Subsequent test classes reuse the cached Spring context, but the
containers are already stopped — causing connection failures.

Instead, use a static initializer or `@BeforeAll` to start the containers,
without the `@Testcontainers` and `@Container` annotations.

## Summary

- Use **JUnit 5 lifecycle callbacks** (`@BeforeAll`/`@AfterAll`) for
  explicit control over container startup and shutdown.
- Use **extension annotations** (`@Testcontainers`/`@Container`) for less
  boilerplate in single test classes.
- Use the **singleton containers pattern** (static initializer in a base class)
  to share containers across multiple test classes.
- Don't mix singleton containers with `@Testcontainers`/`@Container`
  annotations.

## Further reading

- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Testcontainers singleton containers pattern](https://java.testcontainers.org/test_framework_integration/manual_lifecycle_control/#singleton-containers)
- [Testing a Spring Boot REST API with Testcontainers](/guides/testcontainers-java-spring-boot-rest-api/)
