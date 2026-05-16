---
title: Use the JUnit 5 extension for more control
linkTitle: JUnit 5 extension
description: Use the Testcontainers JUnit 5 extension for more control over the PostgreSQL container.
weight: 30
---

If the special JDBC URL doesn't meet your needs, or you need more control over
container creation (for example, to copy initialization scripts), use the
Testcontainers JUnit 5 extension:

```java
@DataJpaTest
@TestPropertySource(properties = {
    "spring.test.database.replace=none"
})
@Testcontainers
class ProductRepositoryTest {

  @Container
  static PostgreSQLContainer postgres =
    new PostgreSQLContainer("postgres:16-alpine")
      .withCopyFileToContainer(
        MountableFile.forClasspathResource("sql/init-db.sql"),
        "/docker-entrypoint-initdb.d/init-db.sql");

  @DynamicPropertySource
  static void configureProperties(DynamicPropertyRegistry registry) {
    registry.add("spring.datasource.url", postgres::getJdbcUrl);
    registry.add("spring.datasource.username", postgres::getUsername);
    registry.add("spring.datasource.password", postgres::getPassword);
  }

  @Autowired
  ProductRepository productRepository;

  @Test
  @Sql("/sql/seed-data.sql")
  void shouldGetAllProducts() {
    List<Product> products = productRepository.findAll();
    assertEquals(2, products.size());
  }

  @Test
  @Sql("/sql/seed-data.sql")
  void shouldNotCreateAProductWithDuplicateCode() {
    Product product = new Product(3L, "p101", "Test Product");
    productRepository.createProductIfNotExists(product);
    Optional<Product> optionalProduct = productRepository.findById(
      product.getId()
    );
    assertThat(optionalProduct).isEmpty();
  }
}
```

This approach:

- Uses `@Testcontainers` and `@Container` to manage the container lifecycle.
- Copies `init-db.sql` into the container's init directory so PostgreSQL
  runs it at startup.
- Uses `@DynamicPropertySource` to register the container's connection details
  with Spring Boot.
- Tests PostgreSQL-specific features like `ON CONFLICT DO NOTHING` that
  wouldn't work with H2.

## Summary

- Use the **special JDBC URL** (`jdbc:tc:postgresql:...`) for the quickest way
  to switch from H2 to a real database — it's a one-property change.
- Use the **JUnit 5 extension** when you need more control over the container
  (custom init scripts, environment variables, etc.).
- Both approaches work with Spring Data JPA (`@DataJpaTest`) and JdbcTemplate
  (`@JdbcTest`) tests.

## Further reading

- [Testcontainers Postgres module](https://java.testcontainers.org/modules/databases/postgres/)
- [Testcontainers JDBC support](https://java.testcontainers.org/modules/databases/jdbc/)
- [Testing a Spring Boot REST API with Testcontainers](/guides/testcontainers-java-spring-boot-rest-api/)
