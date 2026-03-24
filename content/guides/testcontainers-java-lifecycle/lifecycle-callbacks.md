---
title: JUnit 5 lifecycle callbacks
linkTitle: Lifecycle callbacks
description: Manage Testcontainers container lifecycle using JUnit 5 @BeforeAll and @AfterAll callbacks.
weight: 20
---

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
