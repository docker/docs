---
title: Copy files into containers
linkTitle: Copy files
description: Initialize containers by copying files into specific locations.
weight: 10
---

Sometimes you need to initialize a container by placing files in a specific
location. For example, PostgreSQL runs SQL scripts from
`/docker-entrypoint-initdb.d/` when the container starts.

## Create the initialization script

Create `src/test/resources/init-db.sql`:

```sql
create table customers (
     id bigint not null,
     name varchar not null,
     primary key (id)
);
```

## Copy the file into the container

Use `withCopyFileToContainer()` to copy the SQL script into the container's
init directory:

```java
package com.testcontainers.demo;

import static org.junit.jupiter.api.Assertions.assertFalse;

import java.util.List;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.testcontainers.postgresql.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.MountableFile;

@Testcontainers
class CustomerServiceTest {

  @Container
  static PostgreSQLContainer postgres = new PostgreSQLContainer(
    "postgres:16-alpine"
  )
    .withCopyFileToContainer(
      MountableFile.forClasspathResource("init-db.sql"),
      "/docker-entrypoint-initdb.d/"
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
  }

  @Test
  void shouldGetCustomers() {
    customerService.createCustomer(new Customer(1L, "George"));
    customerService.createCustomer(new Customer(2L, "John"));

    List<Customer> customers = customerService.getAllCustomers();
    assertFalse(customers.isEmpty());
  }
}
```

The `withCopyFileToContainer(MountableFile, String)` method copies `init-db.sql`
from the classpath into `/docker-entrypoint-initdb.d/` inside the container.
PostgreSQL executes scripts in that directory automatically at startup.

You can also copy files from any path on the host:

```java
.withCopyFileToContainer(
    MountableFile.forHostPath("/host/path/to/init-db.sql"),
    "/docker-entrypoint-initdb.d/"
);
```
