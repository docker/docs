---
title: JUnit 5 extension annotations
linkTitle: Extension annotations
description: Manage Testcontainers container lifecycle using @Testcontainers and @Container annotations.
weight: 30
---

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
