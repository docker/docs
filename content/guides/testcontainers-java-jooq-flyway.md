---
title: Working with jOOQ and Flyway using Testcontainers
linkTitle: jOOQ and Flyway
description: Learn how to generate jOOQ code from a database using Testcontainers and Flyway, then test your persistence layer.
keywords: testcontainers, java, testing, jooq, flyway, postgresql, spring boot, code generation
summary: |
  Generate typesafe jOOQ code from a real PostgreSQL database managed by
  Flyway migrations, then test repositories using Testcontainers.
aliases:
  - /guides/testcontainers-java-jooq-flyway/create-project/
  - /guides/testcontainers-java-jooq-flyway/run-tests/
  - /guides/testcontainers-java-jooq-flyway/write-tests/
params:
  tags: [testing]
  time: 25 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-working-with-jooq-flyway-using-testcontainers -->

In this guide, you will learn how to:

- Create a Spring Boot application with jOOQ support
- Generate jOOQ code using Testcontainers, Flyway, and a Maven plugin
- Implement basic database operations using jOOQ
- Load complex object graphs using jOOQ's MULTISET feature
- Test the jOOQ persistence layer using Testcontainers

## Prerequisites

- Java 17+
- Maven
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Spring Boot project

### Set up the project

Create a Spring Boot project from [Spring Initializr](https://start.spring.io)
by selecting Maven as the build tool and adding the **JOOQ Access Layer**,
**Flyway Migration**, **Spring Boot DevTools**, **PostgreSQL Driver**, and
**Testcontainers** starters.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-working-with-jooq-flyway-using-testcontainers).

jOOQ (jOOQ Object Oriented Querying) provides a fluent API for building
typesafe SQL queries. To get the full benefit of its typesafe DSL, you need
to generate Java code from your database tables, views, and other objects.

> [!TIP]
> To learn more about how the jOOQ code generator helps, read
> [Why You Should Use jOOQ With Code Generation](https://blog.jooq.org/why-you-should-use-jooq-with-code-generation/).

The typical process for building and testing the application with jOOQ code
generation is:

1. Create a database instance using Testcontainers.
2. Apply Flyway database migrations.
3. Run the jOOQ code generator to produce Java code from the database objects.
4. Run integration tests.

The
[testcontainers-jooq-codegen-maven-plugin](https://github.com/testcontainers/testcontainers-jooq-codegen-maven-plugin)
automates this as part of the Maven build.

### Create Flyway migration scripts

The sample application has `users`, `posts`, and `comments` tables. Create
the first migration script following the Flyway naming convention.

Create `src/main/resources/db/migration/V1__create_tables.sql`:

```sql
create table users
(
    id         bigserial not null,
    name       varchar   not null,
    email      varchar   not null,
    created_at timestamp,
    updated_at timestamp,
    primary key (id),
    constraint user_email_unique unique (email)
);

create table posts
(
    id         bigserial                    not null,
    title      varchar                      not null,
    content    varchar                      not null,
    created_by bigint references users (id) not null,
    created_at timestamp,
    updated_at timestamp,
    primary key (id)
);

create table comments
(
    id         bigserial                    not null,
    name       varchar                      not null,
    content    varchar                      not null,
    post_id    bigint references posts (id) not null,
    created_at timestamp,
    updated_at timestamp,
    primary key (id)
);

ALTER SEQUENCE users_id_seq RESTART WITH 101;
ALTER SEQUENCE posts_id_seq RESTART WITH 101;
ALTER SEQUENCE comments_id_seq RESTART WITH 101;
```

The sequence values restart at 101 so that you can insert sample data with
explicit primary key values for testing.

### Configure jOOQ code generation

Add the `testcontainers-jooq-codegen-maven-plugin` to `pom.xml`:

```xml
<properties>
    <testcontainers.version>2.0.4</testcontainers.version>
    <testcontainers-jooq-codegen-maven-plugin.version>0.0.4</testcontainers-jooq-codegen-maven-plugin.version>
</properties>

<build>
    <plugins>
        <plugin>
            <groupId>org.testcontainers</groupId>
            <artifactId>testcontainers-jooq-codegen-maven-plugin</artifactId>
            <version>${testcontainers-jooq-codegen-maven-plugin.version}</version>
            <dependencies>
                <dependency>
                    <groupId>org.testcontainers</groupId>
                    <artifactId>testcontainers-postgresql</artifactId>
                    <version>${testcontainers.version}</version>
                </dependency>
                <dependency>
                    <groupId>org.postgresql</groupId>
                    <artifactId>postgresql</artifactId>
                    <version>${postgresql.version}</version>
                </dependency>
            </dependencies>
            <executions>
                <execution>
                    <id>generate-jooq-sources</id>
                    <goals>
                        <goal>generate</goal>
                    </goals>
                    <phase>generate-sources</phase>
                    <configuration>
                        <database>
                            <type>POSTGRES</type>
                            <containerImage>postgres:16-alpine</containerImage>
                        </database>
                        <flyway>
                            <locations>
                                filesystem:src/main/resources/db/migration
                            </locations>
                        </flyway>
                        <jooq>
                            <generator>
                                <database>
                                    <includes>.*</includes>
                                    <excludes>flyway_schema_history</excludes>
                                    <inputSchema>public</inputSchema>
                                </database>
                                <target>
                                    <packageName>com.testcontainers.demo.jooq</packageName>
                                    <directory>target/generated-sources/jooq</directory>
                                </target>
                            </generator>
                        </jooq>
                    </configuration>
                </execution>
            </executions>
        </plugin>
    </plugins>
</build>
```

Here's what the plugin configuration does:

- The `<configuration>/<database>` section sets the database type to
  `POSTGRES` and the Docker image to `postgres:16-alpine`.
- The `<configuration>/<flyway>` section points to the Flyway migration
  scripts.
- The `<configuration>/<jooq>` section configures the package name and
  output directory for the generated code. You can use any configuration
  option that the official `jooq-code-generator` plugin supports.

When you run `./mvnw clean package`, the plugin uses Testcontainers to
spin up a PostgreSQL container, applies the Flyway migrations, and generates
Java code under `target/generated-sources/jooq`.

### Create model classes

Create model classes to represent the data structures for various use cases.
These records hold a subset of column values from the tables.

`User.java`:

```java
package com.testcontainers.demo.domain;

public record User(Long id, String name, String email) {}
```

`Post.java`:

```java
package com.testcontainers.demo.domain;

import java.time.LocalDateTime;
import java.util.List;

public record Post(
  Long id,
  String title,
  String content,
  User createdBy,
  List<Comment> comments,
  LocalDateTime createdAt,
  LocalDateTime updatedAt
) {}
```

`Comment.java`:

```java
package com.testcontainers.demo.domain;

import java.time.LocalDateTime;

public record Comment(
  Long id,
  String name,
  String content,
  LocalDateTime createdAt,
  LocalDateTime updatedAt
) {}
```

### Implement repositories using jOOQ

Create `UserRepository.java` with methods to create a user and look up a user
by email:

```java
package com.testcontainers.demo.domain;

import static com.testcontainers.demo.jooq.tables.Users.USERS;
import static org.jooq.Records.mapping;

import java.time.LocalDateTime;
import java.util.Optional;
import org.jooq.DSLContext;
import org.springframework.stereotype.Repository;

@Repository
class UserRepository {

  private final DSLContext dsl;

  UserRepository(DSLContext dsl) {
    this.dsl = dsl;
  }

  public User createUser(User user) {
    return this.dsl.insertInto(USERS)
      .set(USERS.NAME, user.name())
      .set(USERS.EMAIL, user.email())
      .set(USERS.CREATED_AT, LocalDateTime.now())
      .returningResult(USERS.ID, USERS.NAME, USERS.EMAIL)
      .fetchOne(mapping(User::new));
  }

  public Optional<User> getUserByEmail(String email) {
    return this.dsl.select(USERS.ID, USERS.NAME, USERS.EMAIL)
      .from(USERS)
      .where(USERS.EMAIL.equalIgnoreCase(email))
      .fetchOptional(mapping(User::new));
  }
}
```

The jOOQ DSL looks similar to SQL but written in Java. Because the code is
generated from the database schema, it stays in sync with the database
structure and provides type safety. For example,
`where(USERS.EMAIL.equalIgnoreCase(email))` expects a `String` value. If you
pass a non-string value like `123`, you get a compiler error.

### Fetch complex object graphs

jOOQ shines when it comes to complex queries. The database has a many-to-one
relationship from `Post` to `User` and a one-to-many relationship from `Post`
to `Comment`.

Create `PostRepository.java` to load a `Post` with its creator and comments
using a single query with jOOQ's MULTISET feature:

```java
package com.testcontainers.demo.domain;

import static com.testcontainers.demo.jooq.Tables.COMMENTS;
import static com.testcontainers.demo.jooq.tables.Posts.POSTS;
import static org.jooq.Records.mapping;
import static org.jooq.impl.DSL.multiset;
import static org.jooq.impl.DSL.row;
import static org.jooq.impl.DSL.select;

import java.util.Optional;
import org.jooq.DSLContext;
import org.springframework.stereotype.Repository;

@Repository
class PostRepository {

  private final DSLContext dsl;

  PostRepository(DSLContext dsl) {
    this.dsl = dsl;
  }

  public Optional<Post> getPostById(Long id) {
    return this.dsl.select(
        POSTS.ID,
        POSTS.TITLE,
        POSTS.CONTENT,
        row(POSTS.users().ID, POSTS.users().NAME, POSTS.users().EMAIL)
          .mapping(User::new)
          .as("createdBy"),
        multiset(
          select(
            COMMENTS.ID,
            COMMENTS.NAME,
            COMMENTS.CONTENT,
            COMMENTS.CREATED_AT,
            COMMENTS.UPDATED_AT
          )
            .from(COMMENTS)
            .where(POSTS.ID.eq(COMMENTS.POST_ID))
        )
          .as("comments")
          .convertFrom(r -> r.map(mapping(Comment::new))),
        POSTS.CREATED_AT,
        POSTS.UPDATED_AT
      )
      .from(POSTS)
      .where(POSTS.ID.eq(id))
      .fetchOptional(mapping(Post::new));
  }
}
```

This uses jOOQ's
[nested records](https://www.jooq.org/doc/latest/manual/sql-building/column-expressions/nested-records/)
for the many-to-one `Post`-to-`User` association and
[MULTISET](https://www.jooq.org/doc/latest/manual/sql-building/column-expressions/multiset-value-constructor/)
for the one-to-many `Post`-to-`Comment` association.

## Write tests with Testcontainers

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

### Test with the @JooqTest slice

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

### Integration test with @SpringBootTest

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

### Test PostRepository

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

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

You should see the PostgreSQL Docker container start, jOOQ code generation
complete, and all tests pass. After the tests finish, the container stops and
is removed automatically.

### Summary

The Testcontainers library helps you generate Java code from the database
using the jOOQ code generator and test your persistence layer against the
same type of database (PostgreSQL) that you use in production, instead of
mocks or in-memory databases.

Because the code is always generated from the database's current state, you
can be confident that your code stays in sync with database changes. You're
free to refactor and still verify that the application works as expected.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [jOOQ documentation](https://www.jooq.org/)
- [jOOQ code generation](https://www.jooq.org/doc/latest/manual/code-generation/)
- [Spring Boot Testcontainers support](https://docs.spring.io/spring-boot/reference/testing/testcontainers.html)
- [Replace H2 with a real database for testing](/guides/testcontainers-java-replace-h2/)
