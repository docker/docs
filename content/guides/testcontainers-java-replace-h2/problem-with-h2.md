---
title: The problem with H2 for testing
linkTitle: The H2 problem
description: Understand why using H2 in-memory databases for testing gives false confidence.
weight: 10
---

A common practice is to use lightweight databases like H2 or HSQL as
in-memory databases for testing while using PostgreSQL, MySQL, or Oracle in
production. This approach has significant drawbacks:

- The test database might not support all features of your production database.
- SQL syntax might not be compatible between H2 and your production database.
- Tests passing with H2 don't guarantee they'll work in production.

## Example: PostgreSQL-specific syntax

Consider implementing an "upsert" — insert a product only if it doesn't
already exist. In PostgreSQL, you can use:

```sql
INSERT INTO products(id, code, name) VALUES(?,?,?) ON CONFLICT DO NOTHING;
```

This query doesn't work with H2 by default:

```text
Caused by: org.h2.jdbc.JdbcSQLException: Syntax error in SQL statement
"INSERT INTO products (id, code, name) VALUES (?, ?, ?) ON[*] CONFLICT DO NOTHING";
```

You can run H2 in PostgreSQL compatibility mode, but not all features are
supported. The inverse is also true — H2 supports `ROWNUM()` which PostgreSQL
doesn't.

Testing with a different database than production means you can't trust your
test results and must verify after deployment, defeating the purpose of
automated tests.

## The Spring Boot test using H2

A typical H2-based test looks like this:

```java
@DataJpaTest
class ProductRepositoryTest {

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

Spring Boot uses H2 automatically when it's on the classpath. The test passes,
but it doesn't catch PostgreSQL-specific issues.
