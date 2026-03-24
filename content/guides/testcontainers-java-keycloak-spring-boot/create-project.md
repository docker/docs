---
title: Create the Spring Boot project
linkTitle: Create the project
description: Set up a Spring Boot OAuth 2.0 Resource Server with Keycloak, PostgreSQL, and Testcontainers.
weight: 10
---

## Set up the project

Create a Spring Boot project from [Spring Initializr](https://start.spring.io)
by selecting the **Spring Web**, **Validation**, **JDBC API**,
**PostgreSQL Driver**, **Spring Security**, **OAuth2 Resource Server**, and
**Testcontainers** starters.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-securing-spring-boot-microservice-using-keycloak-and-testcontainers).

After generating the application, add the
[testcontainers-keycloak](https://github.com/dasniko/testcontainers-keycloak)
community module and [REST Assured](https://rest-assured.io/) as test
dependencies.

The key dependencies in `pom.xml` are:

```xml
<properties>
    <java.version>17</java.version>
    <testcontainers.version>2.0.4</testcontainers.version>
</properties>
<dependencies>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-validation</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-jdbc</artifactId>
    </dependency>
    <dependency>
        <groupId>org.postgresql</groupId>
        <artifactId>postgresql</artifactId>
        <scope>runtime</scope>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-security</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-oauth2-resource-server</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-test</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.springframework.security</groupId>
        <artifactId>spring-security-test</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-testcontainers</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-junit-jupiter</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-postgresql</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>com.github.dasniko</groupId>
        <artifactId>testcontainers-keycloak</artifactId>
        <version>3.4.0</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>io.rest-assured</groupId>
        <artifactId>rest-assured</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>
```

## Create the domain model

Create a `Product` record that represents the domain object:

```java
package com.testcontainers.products.domain;

import jakarta.validation.constraints.NotEmpty;

public record Product(Long id, @NotEmpty String title, String description) {}
```

## Create the repository

Implement `ProductRepository` using Spring `JdbcClient` to interact with a
PostgreSQL database:

```java
package com.testcontainers.products.domain;

import java.util.List;
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.jdbc.support.GeneratedKeyHolder;
import org.springframework.jdbc.support.KeyHolder;
import org.springframework.stereotype.Repository;

@Repository
public class ProductRepository {

  private final JdbcClient jdbcClient;

  public ProductRepository(JdbcClient jdbcClient) {
    this.jdbcClient = jdbcClient;
  }

  public List<Product> getAll() {
    return jdbcClient.sql("SELECT * FROM products").query(Product.class).list();
  }

  public Product create(Product product) {
    String sql =
      "INSERT INTO products(title, description) VALUES (:title,:description) RETURNING id";
    KeyHolder keyHolder = new GeneratedKeyHolder();
    jdbcClient
      .sql(sql)
      .param("title", product.title())
      .param("description", product.description())
      .update(keyHolder);
    Long id = keyHolder.getKeyAs(Long.class);
    return new Product(id, product.title(), product.description());
  }
}
```

## Add a schema creation script

Create `src/main/resources/schema.sql` to initialize the `products` table:

```sql
CREATE TABLE products (
    id bigserial primary key,
    title varchar not null,
    description text
);
```

Enable schema initialization in `src/main/resources/application.properties`:

```properties
spring.sql.init.mode=always
```

For production applications, use a database migration tool like Flyway or
Liquibase instead.

## Implement the API endpoints

Create `ProductController` with endpoints to fetch all products and create a
product:

```java
package com.testcontainers.products.api;

import com.testcontainers.products.domain.Product;
import com.testcontainers.products.domain.ProductRepository;
import jakarta.validation.Valid;
import java.util.List;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/products")
class ProductController {

  private final ProductRepository productRepository;

  ProductController(ProductRepository productRepository) {
    this.productRepository = productRepository;
  }

  @GetMapping
  List<Product> getAll() {
    return productRepository.getAll();
  }

  @PostMapping
  @ResponseStatus(HttpStatus.CREATED)
  Product createProduct(@RequestBody @Valid Product product) {
    return productRepository.create(product);
  }
}
```

## Configure OAuth 2.0 security

Create a `SecurityConfig` class that protects the API endpoints using JWT
token-based authentication:

```java
package com.testcontainers.products.config;

import static org.springframework.security.config.Customizer.withDefaults;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.CorsConfigurer;
import org.springframework.security.config.annotation.web.configurers.CsrfConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
class SecurityConfig {

  @Bean
  SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
    http
      .authorizeHttpRequests(c ->
        c
          .requestMatchers(HttpMethod.GET, "/api/products")
          .permitAll()
          .requestMatchers(HttpMethod.POST, "/api/products")
          .authenticated()
          .anyRequest()
          .authenticated()
      )
      .sessionManagement(c ->
        c.sessionCreationPolicy(SessionCreationPolicy.STATELESS)
      )
      .cors(CorsConfigurer::disable)
      .csrf(CsrfConfigurer::disable)
      .oauth2ResourceServer(oauth2 -> oauth2.jwt(withDefaults()));
    return http.build();
  }
}
```

This configuration:

- Permits unauthenticated access to `GET /api/products`.
- Requires authentication for `POST /api/products` and all other endpoints.
- Configures the OAuth 2.0 Resource Server with JWT token-based authentication.
- Disables CORS and CSRF because this is a stateless API.

Add the JWT issuer URI to `application.properties`:

```properties
spring.security.oauth2.resourceserver.jwt.issuer-uri=http://localhost:9090/realms/keycloaktcdemo
```

## Export the Keycloak realm configuration

Before writing the tests, export a Keycloak realm configuration so that the test
environment can import it automatically. Start a temporary Keycloak instance:

```console
$ docker run -p 9090:8080 \
    -e KEYCLOAK_ADMIN=admin \
    -e KEYCLOAK_ADMIN_PASSWORD=admin \
    quay.io/keycloak/keycloak:25 start-dev
```

Open `http://localhost:9090` and sign in to the Admin Console with `admin/admin`.
Then set up the realm:

1. In the top-left corner, select the realm drop-down and create a realm named
   `keycloaktcdemo`.
2. Under the `keycloaktcdemo` realm, create a client with the following
   settings:
   - **Client ID**: `product-service`
   - **Client Authentication**: **On**
   - **Authentication flow**: select only **Service accounts roles**
3. On the **Client details** screen, go to the **Credentials** tab and copy the
   **Client secret** value.

Export the realm configuration:

```console
$ docker ps
# copy the keycloak container id

$ docker exec -it <container-id> /bin/bash

$ /opt/keycloak/bin/kc.sh export --dir /opt/keycloak/data/import --realm keycloaktcdemo

$ exit

$ docker cp <container-id>:/opt/keycloak/data/import/keycloaktcdemo-realm.json keycloaktcdemo-realm.json
```

Copy the exported `keycloaktcdemo-realm.json` file into `src/test/resources`.
