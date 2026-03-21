---
title: Replace H2 with the Testcontainers JDBC URL
linkTitle: JDBC URL approach
description: Use the Testcontainers special JDBC URL to swap H2 for a real PostgreSQL database.
weight: 20
---

Replacing H2 with a real PostgreSQL database requires two test properties:

```java
@DataJpaTest
@TestPropertySource(properties = {
  "spring.test.database.replace=none",
  "spring.datasource.url=jdbc:tc:postgresql:16-alpine:///db"
})
class ProductRepositoryWithJdbcUrlTest {

  @Autowired
  ProductRepository productRepository;

  @Test
  @Sql("classpath:/sql/seed-data.sql")
  void shouldGetAllProducts() {
    List<Product> products = productRepository.findAll();
    assertEquals(2, products.size());
  }
}
```

That's it — two properties and your tests run against a real PostgreSQL
database.

## How the special JDBC URL works

A standard PostgreSQL JDBC URL looks like:

```text
jdbc:postgresql://localhost:5432/postgres
```

The Testcontainers special JDBC URL inserts `tc:` after `jdbc:`:

```text
jdbc:tc:postgresql:///db
```

The hostname, port, and database name are ignored — Testcontainers manages them
automatically. You can specify the Docker image tag after the database name:

```text
jdbc:tc:postgresql:16-alpine:///db
```

This creates a container from the `postgres:16-alpine` image.

## Initialize the database with a script

Pass `TC_INITSCRIPT` to run an SQL script when the container starts:

```text
jdbc:tc:postgresql:16-alpine:///db?TC_INITSCRIPT=sql/init-db.sql
```

Testcontainers runs the script automatically. For production applications,
use a database migration tool like Flyway or Liquibase instead.

The special JDBC URL also works for MySQL, MariaDB, PostGIS, YugabyteDB,
CockroachDB, and other databases with Testcontainers JDBC support.

## Testing JdbcTemplate-based repositories

The same approach works for `JdbcTemplate`-based repositories. Use `@JdbcTest`
instead of `@DataJpaTest`:

```java
@JdbcTest
@TestPropertySource(properties = {
  "spring.test.database.replace=none",
  "spring.datasource.url=jdbc:tc:postgresql:16-alpine:///db?TC_INITSCRIPT=sql/init-db.sql"
})
class JdbcProductRepositoryTest {

  @Autowired
  private JdbcTemplate jdbcTemplate;

  private JdbcProductRepository productRepository;

  @BeforeEach
  void setUp() {
    productRepository = new JdbcProductRepository(jdbcTemplate);
  }

  @Test
  @Sql("/sql/seed-data.sql")
  void shouldGetAllProducts() {
    List<Product> products = productRepository.getAllProducts();
    assertEquals(2, products.size());
  }
}
```
