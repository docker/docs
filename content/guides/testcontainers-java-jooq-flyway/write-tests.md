---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Test jOOQ repositories using Testcontainers with the @JooqTest slice and @SpringBootTest.
weight: 20
---

Before writing the tests, create an SQL script to seed test data at
`src/test/resources/test-data.sql`:

```sql
DELETE FROM comments;
DELETE FROM posts;
DELETE FROM users;

INSERT INTO users(id, name, email) VALUES
(1, 'Siva', 'siva@gmail.com'),
(2, 'Oleg', 'oleg@gmail.com');

INSERT INTO posts(id, title, content, created_by, created_at) VALUES
(1, 'Post 1 Title', 'Post 1 content', 1, CURRENT_TIMESTAMP),
(2, 'Post 2 Title', 'Post 2 content', 2, CURRENT_TIMESTAMP);

INSERT INTO comments(id, name, content, post_id, created_at) VALUES
(1, 'Ron', 'Comment 1', 1, CURRENT_TIMESTAMP),
(2, 'James', 'Comment 2', 1, CURRENT_TIMESTAMP),
(3, 'Robert', 'Comment 3', 2, CURRENT_TIMESTAMP);
```

## Test with the @JooqTest slice

The `@JooqTest` annotation loads only the persistence layer components and
auto-configures jOOQ's `DSLContext`. Use the Testcontainers special JDBC URL
to start a Postgres container.

Create `UserRepositoryJooqTest.java`:

```java
package com.testcontainers.demo.domain;

import static org.assertj.core.api.Assertions.assertThat;

import org.jooq.DSLContext;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jooq.JooqTest;
import org.springframework.test.context.jdbc.Sql;

@JooqTest(
  properties = {
    "spring.test.database.replace=none",
    "spring.datasource.url=jdbc:tc:postgresql:16-alpine:///db",
  }
)
@Sql("/test-data.sql")
class UserRepositoryJooqTest {

  @Autowired
  DSLContext dsl;

  UserRepository repository;

  @BeforeEach
  void setUp() {
    this.repository = new UserRepository(dsl);
  }

  @Test
  void shouldCreateUserSuccessfully() {
    User user = new User(null, "John", "john@gmail.com");

    User savedUser = repository.createUser(user);

    assertThat(savedUser.id()).isNotNull();
    assertThat(savedUser.name()).isEqualTo("John");
    assertThat(savedUser.email()).isEqualTo("john@gmail.com");
  }

  @Test
  void shouldGetUserByEmail() {
    User user = repository.getUserByEmail("siva@gmail.com").orElseThrow();

    assertThat(user.id()).isEqualTo(1L);
    assertThat(user.name()).isEqualTo("Siva");
    assertThat(user.email()).isEqualTo("siva@gmail.com");
  }
}
```

Here's what the test does:

- `@JooqTest` loads only the persistence layer and auto-configures
  `DSLContext`.
- The Testcontainers special JDBC URL
  (`jdbc:tc:postgresql:16-alpine:///db`) starts a PostgreSQL container
  automatically.
- Because `flyway-core` is on the classpath, Spring Boot runs the Flyway
  migrations from `src/main/resources/db/migration` on startup.
- `@Sql("/test-data.sql")` loads the test data before each test.
- The `UserRepository` is instantiated manually with the injected
  `DSLContext`.

## Integration test with @SpringBootTest

For a full integration test, use `@SpringBootTest` with the Testcontainers
`@ServiceConnection` support introduced in Spring Boot 3.1.

Create `UserRepositoryTest.java`:

```java
package com.testcontainers.demo.domain;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.test.context.jdbc.Sql;
import org.testcontainers.postgresql.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

@SpringBootTest
@Sql("/test-data.sql")
@Testcontainers
class UserRepositoryTest {

  @Container
  @ServiceConnection
  static PostgreSQLContainer postgres = new PostgreSQLContainer(
    "postgres:16-alpine"
  );

  @Autowired
  UserRepository repository;

  @Test
  void shouldCreateUserSuccessfully() {
    User user = new User(null, "John", "john@gmail.com");

    User savedUser = repository.createUser(user);

    assertThat(savedUser.id()).isNotNull();
    assertThat(savedUser.name()).isEqualTo("John");
    assertThat(savedUser.email()).isEqualTo("john@gmail.com");
  }

  @Test
  void shouldGetUserByEmail() {
    User user = repository.getUserByEmail("siva@gmail.com").orElseThrow();

    assertThat(user.id()).isEqualTo(1L);
    assertThat(user.name()).isEqualTo("Siva");
    assertThat(user.email()).isEqualTo("siva@gmail.com");
  }
}
```

Here's what the test does:

- `@SpringBootTest` loads the entire application context, so
  `UserRepository` is injected directly.
- `@Testcontainers` and `@Container` manage the PostgreSQL container
  lifecycle.
- `@ServiceConnection` auto-configures the datasource properties from the
  running container, replacing the need for `@DynamicPropertySource`.
- `@Sql("/test-data.sql")` initializes the test data.

## Test PostRepository

Test the `PostRepository` that fetches complex object graphs using the
Testcontainers special JDBC URL.

Create `PostRepositoryTest.java`:

```java
package com.testcontainers.demo.domain;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.jdbc.Sql;

@SpringBootTest(
  properties = {
    "spring.test.database.replace=none",
    "spring.datasource.url=jdbc:tc:postgresql:16-alpine:///db",
  }
)
@Sql("/test-data.sql")
class PostRepositoryTest {

  @Autowired
  PostRepository repository;

  @Test
  void shouldGetPostById() {
    Post post = repository.getPostById(1L).orElseThrow();

    assertThat(post.id()).isEqualTo(1L);
    assertThat(post.title()).isEqualTo("Post 1 Title");
    assertThat(post.content()).isEqualTo("Post 1 content");
    assertThat(post.createdBy().id()).isEqualTo(1L);
    assertThat(post.createdBy().name()).isEqualTo("Siva");
    assertThat(post.createdBy().email()).isEqualTo("siva@gmail.com");
    assertThat(post.comments()).hasSize(2);
  }
}
```

This test verifies that `getPostById` loads the post along with its creator
and comments in a single query using jOOQ's MULTISET feature.
