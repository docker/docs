---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Write your first integration test using Testcontainers for Java and PostgreSQL.
weight: 20
---

You have the `CustomerService` implementation ready, but for testing you need a
PostgreSQL database. You can use Testcontainers to spin up a Postgres database
in a Docker container and run your tests against it.

## Add Testcontainers dependencies

Add the Testcontainers PostgreSQL module as a test dependency in `pom.xml`:

```xml
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>testcontainers-postgresql</artifactId>
    <version>2.0.4</version>
    <scope>test</scope>
</dependency>
```

Since the application uses a Postgres database, the Testcontainers Postgres
module provides a `PostgreSQLContainer` class for managing the container.

## Write the test

Create `CustomerServiceTest.java` under `src/test/java`:

```java
package com.testcontainers.demo;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.List;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.testcontainers.postgresql.PostgreSQLContainer;

class CustomerServiceTest {

  static PostgreSQLContainer postgres = new PostgreSQLContainer(
    "postgres:16-alpine"
  );

  CustomerService customerService;

  @BeforeAll
  static void beforeAll() {
    postgres.start();
  }

  @AfterAll
  static void afterAll() {
    postgres.stop();
  }

  @BeforeEach
  void setUp() {
    DBConnectionProvider connectionProvider = new DBConnectionProvider(
      postgres.getJdbcUrl(),
      postgres.getUsername(),
      postgres.getPassword()
    );
    customerService = new CustomerService(connectionProvider);
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

Here's what the test does:

- Declares a `PostgreSQLContainer` with the `postgres:16-alpine` Docker image.
- The `@BeforeAll` callback starts the Postgres container before any test
  methods run.
- The `@BeforeEach` callback creates a `DBConnectionProvider` using the JDBC
  connection parameters from the container, then creates a `CustomerService`.
  The `CustomerService` constructor creates the `customers` table if it
  doesn't exist.
- `shouldGetCustomers()` inserts 2 customer records, fetches all customers,
  and asserts the count.
- The `@AfterAll` callback stops the container after all test methods finish.
