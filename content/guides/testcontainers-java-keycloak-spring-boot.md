---
title: Securing Spring Boot microservice using Keycloak and Testcontainers
linkTitle: Keycloak with Spring Boot
description: Learn how to secure a Spring Boot microservice using Keycloak and test it with the Testcontainers Keycloak module.
keywords: testcontainers, java, spring boot, testing, keycloak, security, oauth2, jwt
summary: |
  Learn how to create an OAuth 2.0 Resource Server using Spring Boot, secure API
  endpoints with Keycloak, and test the application using the Testcontainers Keycloak module.
aliases:
  - /guides/testcontainers-java-keycloak-spring-boot/create-project/
  - /guides/testcontainers-java-keycloak-spring-boot/run-tests/
  - /guides/testcontainers-java-keycloak-spring-boot/write-tests/
params:
  tags: [testing]
  time: 30 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-securing-spring-boot-microservice-using-keycloak-and-testcontainers -->

In this guide, you'll learn how to:

- Create an OAuth 2.0 Resource Server using Spring Boot
- Secure API endpoints using Keycloak
- Test the APIs using the Testcontainers Keycloak module
- Run the application locally using the Testcontainers Keycloak module

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Spring Boot project

### Set up the project

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

### Create the domain model

Create a `Product` record that represents the domain object:

```java
package com.testcontainers.products.domain;

import jakarta.validation.constraints.NotEmpty;

public record Product(Long id, @NotEmpty String title, String description) {}
```

### Create the repository

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

### Add a schema creation script

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

### Implement the API endpoints

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

### Configure OAuth 2.0 security

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

### Export the Keycloak realm configuration

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
   - **Client Authentication**: `On`
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

## Write tests with Testcontainers

To test the secured API endpoints, you need a running Keycloak instance and a
PostgreSQL database, plus a started Spring context. Testcontainers spins up both
services in Docker containers and connects them to Spring through dynamic
property registration.

### Configure the test containers

Spring Boot's Testcontainers support lets you declare containers as beans. For
Keycloak, `@ServiceConnection` isn't available, but you can use
`DynamicPropertyRegistry` to set the JWT issuer URI dynamically.

Create `ContainersConfig.java` under `src/test/java`:

```java
package com.testcontainers.products;

import dasniko.testcontainers.keycloak.KeycloakContainer;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.context.annotation.Bean;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.testcontainers.postgresql.PostgreSQLContainer;

@TestConfiguration(proxyBeanMethods = false)
public class ContainersConfig {

  static String POSTGRES_IMAGE = "postgres:16-alpine";
  static String KEYCLOAK_IMAGE = "quay.io/keycloak/keycloak:25.0";
  static String realmImportFile = "/keycloaktcdemo-realm.json";
  static String realmName = "keycloaktcdemo";

  @Bean
  @ServiceConnection
  PostgreSQLContainer postgres() {
    return new PostgreSQLContainer(POSTGRES_IMAGE);
  }

  @Bean
  KeycloakContainer keycloak(DynamicPropertyRegistry registry) {
    var keycloak = new KeycloakContainer(KEYCLOAK_IMAGE)
      .withRealmImportFile(realmImportFile);
    registry.add(
      "spring.security.oauth2.resourceserver.jwt.issuer-uri",
      () -> keycloak.getAuthServerUrl() + "/realms/" + realmName
    );
    return keycloak;
  }
}
```

This configuration:

- Declares a `PostgreSQLContainer` bean with `@ServiceConnection`, which starts
  a PostgreSQL container and automatically registers the datasource properties.
- Declares a `KeycloakContainer` bean using the `quay.io/keycloak/keycloak:25.0`
  image, imports the realm configuration file, and dynamically registers the JWT
  issuer URI from the Keycloak container's auth server URL.

### Write the test

Create `ProductControllerTests.java`:

```java
package com.testcontainers.products.api;

import static io.restassured.RestAssured.given;
import static io.restassured.RestAssured.when;
import static java.util.Collections.singletonList;
import static org.springframework.boot.test.context.SpringBootTest.WebEnvironment.RANDOM_PORT;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.testcontainers.products.ContainersConfig;
import io.restassured.RestAssured;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.security.oauth2.resource.OAuth2ResourceServerProperties;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.context.annotation.Import;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.client.RestTemplate;

@SpringBootTest(webEnvironment = RANDOM_PORT)
@Import(ContainersConfig.class)
class ProductControllerTests {

  static final String GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials";
  static final String CLIENT_ID = "product-service";
  static final String CLIENT_SECRET = "jTJJqdzeCSt3DmypfHZa42vX8U9rQKZ9";

  @LocalServerPort
  private int port;

  @Autowired
  OAuth2ResourceServerProperties oAuth2ResourceServerProperties;

  @BeforeEach
  void setup() {
    RestAssured.port = port;
  }

  @Test
  void shouldGetProductsWithoutAuthToken() {
    when().get("/api/products").then().statusCode(200);
  }

  @Test
  void shouldGetUnauthorizedWhenCreateProductWithoutAuthToken() {
    given()
      .contentType("application/json")
      .body(
        """
            {
                "title": "New Product",
                "description": "Brand New Product"
            }
        """
      )
      .when()
      .post("/api/products")
      .then()
      .statusCode(401);
  }

  @Test
  void shouldCreateProductWithAuthToken() {
    String token = getToken();

    given()
      .header("Authorization", "Bearer " + token)
      .contentType("application/json")
      .body(
        """
            {
                "title": "New Product",
                "description": "Brand New Product"
            }
        """
      )
      .when()
      .post("/api/products")
      .then()
      .statusCode(201);
  }

  private String getToken() {
    RestTemplate restTemplate = new RestTemplate();
    HttpHeaders httpHeaders = new HttpHeaders();
    httpHeaders.setContentType(MediaType.APPLICATION_FORM_URLENCODED);

    MultiValueMap<String, String> map = new LinkedMultiValueMap<>();
    map.put("grant_type", singletonList(GRANT_TYPE_CLIENT_CREDENTIALS));
    map.put("client_id", singletonList(CLIENT_ID));
    map.put("client_secret", singletonList(CLIENT_SECRET));

    String authServerUrl =
      oAuth2ResourceServerProperties.getJwt().getIssuerUri() +
      "/protocol/openid-connect/token";

    var request = new HttpEntity<>(map, httpHeaders);
    KeyCloakToken token = restTemplate.postForObject(
      authServerUrl,
      request,
      KeyCloakToken.class
    );

    assert token != null;
    return token.accessToken();
  }

  record KeyCloakToken(@JsonProperty("access_token") String accessToken) {}
}
```

Here's what the tests cover:

- `shouldGetProductsWithoutAuthToken()` invokes `GET /api/products` without an
  `Authorization` header. Because this endpoint is configured to permit
  unauthenticated access, the response status code is 200.
- `shouldGetUnauthorizedWhenCreateProductWithoutAuthToken()` invokes the secured
  `POST /api/products` endpoint without an `Authorization` header and asserts
  the response status code is 401 (Unauthorized).
- `shouldCreateProductWithAuthToken()` first obtains an `access_token` using the
  Client Credentials flow. It then includes the token as a Bearer token in the
  `Authorization` header when invoking `POST /api/products` and asserts the
  response status code is 201 (Created).

The `getToken()` helper method requests an access token from the Keycloak token
endpoint using the client ID and client secret that were configured in the
exported realm.

### Use Testcontainers for local development

Spring Boot's Testcontainers support also works for local development. Create
`TestApplication.java` under `src/test/java`:

```java
package com.testcontainers.products;

import org.springframework.boot.SpringApplication;

public class TestApplication {

  public static void main(String[] args) {
    SpringApplication
      .from(Application::main)
      .with(ContainersConfig.class)
      .run(args);
  }
}
```

Run `TestApplication.java` from your IDE instead of the main `Application.java`.
It starts the containers defined in `ContainersConfig` and configures the
application to use the dynamically registered properties, so you don't have to
install or configure PostgreSQL and Keycloak manually.

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the Keycloak and PostgreSQL Docker containers start with the
realm settings imported and the tests pass. After the tests finish, the
containers stop and are removed automatically.

### Summary

The Testcontainers Keycloak module lets you develop and test applications using a
real Keycloak server instead of mocks. Testing against a real OAuth 2.0
provider that mirrors your production setup gives you more confidence in your
security configuration and token-based authentication flows.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Getting started with Testcontainers in a Java Spring Boot project](https://testcontainers.com/guides/testing-spring-boot-rest-api-using-testcontainers/)
- [Testcontainers Keycloak module](https://testcontainers.com/modules/keycloak/)
- [testcontainers-keycloak GitHub repository](https://github.com/dasniko/testcontainers-keycloak)
- [Spring Boot OAuth 2.0 Resource Server](https://docs.spring.io/spring-security/reference/servlet/oauth2/resource-server/index.html)
