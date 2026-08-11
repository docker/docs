---
title: Testcontainers container lifecycle management using JUnit 5
linkTitle: Container lifecycle (Java)
description: Learn how to manage Testcontainers container lifecycle using JUnit 5 callbacks, extension annotations, and the singleton containers pattern.
keywords: testcontainers, java, testing, junit, lifecycle, singleton containers, postgresql
summary: |
  Learn different approaches to manage container lifecycle with Testcontainers
  using JUnit 5 lifecycle callbacks, extension annotations, and the singleton
  containers pattern.
aliases:
  - /guides/testcontainers-java-lifecycle/create-project/
  - /guides/testcontainers-java-lifecycle/extension-annotations/
  - /guides/testcontainers-java-lifecycle/lifecycle-callbacks/
  - /guides/testcontainers-java-lifecycle/singleton-containers/
params:
  tags: [testing]
  time: 20 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testcontainers-lifecycle -->

In this guide, you will learn how to:

- Start and stop containers using JUnit 5 lifecycle callbacks
- Manage containers using JUnit 5 extension annotations (`@Testcontainers` and `@Container`)
- Share containers across multiple test classes using the singleton containers pattern
- Avoid a common misconfiguration when combining extension annotations with singleton containers

## Prerequisites

- Java 17+
- Your preferred IDE
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the project and business logic

### Set up the project

Create a Java project with Maven and add the required dependencies:

```xml
<dependencies>
    <dependency>
        <groupId>org.postgresql</groupId>
        <artifactId>postgresql</artifactId>
        <version>42.7.3</version>
    </dependency>
    <dependency>
        <groupId>ch.qos.logback</groupId>
        <artifactId>logback-classic</artifactId>
        <version>1.5.6</version>
    </dependency>
    <dependency>
        <groupId>org.junit.jupiter</groupId>
        <artifactId>junit-jupiter</artifactId>
        <version>5.10.2</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-junit-jupiter</artifactId>
        <version>2.0.4</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-postgresql</artifactId>
        <version>2.0.4</version>
        <scope>test</scope>
    </dependency>
</dependencies>
```

### Create the business logic

Create a `Customer` record:

```java
package com.testcontainers.demo;

public record Customer(Long id, String name) {}
```

Create a `CustomerService` class with methods to create, retrieve, and delete
customers:

```java
package com.testcontainers.demo;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

public class CustomerService {

  private final String url;
  private final String username;
  private final String password;

  public CustomerService(String url, String username, String password) {
    this.url = url;
    this.username = username;
    this.password = password;
    createCustomersTableIfNotExists();
  }

  public void createCustomer(Customer customer) {
    try (Connection conn = this.getConnection()) {
      PreparedStatement pstmt = conn.prepareStatement(
        "insert into customers(id,name) values(?,?)"
      );
      pstmt.setLong(1, customer.id());
      pstmt.setString(2, customer.name());
      pstmt.execute();
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
  }

  public List<Customer> getAllCustomers() {
    List<Customer> customers = new ArrayList<>();
    try (Connection conn = this.getConnection()) {
      PreparedStatement pstmt = conn.prepareStatement(
        "select id,name from customers"
      );
      ResultSet rs = pstmt.executeQuery();
      while (rs.next()) {
        long id = rs.getLong("id");
        String name = rs.getString("name");
        customers.add(new Customer(id, name));
      }
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
    return customers;
  }

  public Optional<Customer> getCustomer(Long customerId) {
    try (Connection conn = this.getConnection()) {
      PreparedStatement pstmt = conn.prepareStatement(
        "select id,name from customers where id = ?"
      );
      pstmt.setLong(1, customerId);
      ResultSet rs = pstmt.executeQuery();
      if (rs.next()) {
        long id = rs.getLong("id");
        String name = rs.getString("name");
        return Optional.of(new Customer(id, name));
      }
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
    return Optional.empty();
  }

  public void deleteAllCustomers() {
    try (Connection conn = this.getConnection()) {
      PreparedStatement pstmt = conn.prepareStatement("delete from customers");
      pstmt.execute();
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
  }

  private void createCustomersTableIfNotExists() {
    try (Connection conn = this.getConnection()) {
      PreparedStatement pstmt = conn.prepareStatement(
        """
        create table if not exists customers (
            id bigint not null,
            name varchar not null,
            primary key (id)
        )
        """
      );
      pstmt.execute();
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
  }

  private Connection getConnection() {
    try {
      return DriverManager.getConnection(url, username, password);
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }
}
```

## JUnit 5 lifecycle callbacks

When testing with Testcontainers, you want to start the required containers
before executing any tests and remove them afterwards. You can use JUnit 5
`@BeforeAll` and `@AfterAll` lifecycle callback methods for this:

```java
package com.testcontainers.demo;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.util.List;
import java.util.Optional;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.testcontainers.postgresql.PostgreSQLContainer;

class CustomerServiceWithLifeCycleCallbacksTest {

  static PostgreSQLContainer postgres = new PostgreSQLContainer(
    "postgres:16-alpine"
  );

  CustomerService customerService;

  @BeforeAll
  static void startContainers() {
    postgres.start();
  }

  @AfterAll
  static void stopContainers() {
    postgres.stop();
  }

  @BeforeEach
  void setUp() {
    customerService =
    new CustomerService(
      postgres.getJdbcUrl(),
      postgres.getUsername(),
      postgres.getPassword()
    );
    customerService.deleteAllCustomers();
  }

  @Test
  void shouldCreateCustomer() {
    customerService.createCustomer(new Customer(1L, "George"));

    Optional<Customer> customer = customerService.getCustomer(1L);
    assertTrue(customer.isPresent());
    assertEquals(1L, customer.get().id());
    assertEquals("George", customer.get().name());
  }

  @Test
  void shouldGetCustomers() {
    customerService.createCustomer(new Customer(1L, "George"));
    customerService.createCustomer(new Customer(2L, "John"));

    List<Customer> customers = customerService.getAllCustomers();
    assertEquals(2, customers.size());
  }
}
```

Here's what the code does:

- `PostgreSQLContainer` is declared as a **static field**. The container starts
  before all tests and stops after all tests in this class.
- `@BeforeAll` starts the container, `@AfterAll` stops it.
- `@BeforeEach` initializes `CustomerService` with the container's JDBC
  parameters and deletes all rows to give each test a clean database.

Key observations:

- Because the container is a **static field**, it's shared across all test
  methods in the class. You can declare it as a non-static field and use
  `@BeforeEach`/`@AfterEach` to start a new container per test, but this
  isn't recommended as it's resource-intensive.
- Even without explicitly stopping the container in `@AfterAll`, Testcontainers
  uses the [Ryuk container](https://github.com/testcontainers/moby-ryuk) to
  clean up containers automatically when the JVM exits.

## JUnit 5 extension annotations

The Testcontainers library provides a JUnit 5 extension that simplifies
starting and stopping containers using annotations. To use it, add the
`org.testcontainers:testcontainers-junit-jupiter` test dependency.

```java
package com.testcontainers.demo;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.util.List;
import java.util.Optional;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.testcontainers.postgresql.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

@Testcontainers
class CustomerServiceWithJUnit5ExtensionTest {

  @Container
  static PostgreSQLContainer postgres = new PostgreSQLContainer(
    "postgres:16-alpine"
  );

  CustomerService customerService;

  @BeforeEach
  void setUp() {
    customerService =
    new CustomerService(
      postgres.getJdbcUrl(),
      postgres.getUsername(),
      postgres.getPassword()
    );
    customerService.deleteAllCustomers();
  }

  @Test
  void shouldCreateCustomer() {
    customerService.createCustomer(new Customer(1L, "George"));

    Optional<Customer> customer = customerService.getCustomer(1L);
    assertTrue(customer.isPresent());
    assertEquals(1L, customer.get().id());
    assertEquals("George", customer.get().name());
  }

  @Test
  void shouldGetCustomers() {
    customerService.createCustomer(new Customer(1L, "George"));
    customerService.createCustomer(new Customer(2L, "John"));

    List<Customer> customers = customerService.getAllCustomers();
    assertEquals(2, customers.size());
  }
}
```

Instead of manually starting and stopping the container in `@BeforeAll` and
`@AfterAll`, the `@Testcontainers` annotation on the class and the
`@Container` annotation on the field handle it automatically:

- The extension finds all `@Container`-annotated fields.
- **Static fields** start once before all tests and stop after all tests.
- **Instance fields** start before each test and stop after each test (not
  recommended — it's resource-intensive).

## Singleton containers pattern

As the number of test classes grows, starting containers for each class adds
up. The singleton containers pattern starts all required containers once in a
common base class and reuses them across all integration tests.

### Define the base class

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

### Extend the base class

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

### Avoid a common misconfiguration

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

### Summary

- Use **JUnit 5 lifecycle callbacks** (`@BeforeAll`/`@AfterAll`) for
  explicit control over container startup and shutdown.
- Use **extension annotations** (`@Testcontainers`/`@Container`) for less
  boilerplate in single test classes.
- Use the **singleton containers pattern** (static initializer in a base class)
  to share containers across multiple test classes.
- Don't mix singleton containers with `@Testcontainers`/`@Container`
  annotations.

### Further reading

- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Testcontainers singleton containers pattern](https://java.testcontainers.org/test_framework_integration/manual_lifecycle_control/#singleton-containers)
- [Testing a Spring Boot REST API with Testcontainers](/guides/testcontainers-java-spring-boot-rest-api/)
