---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Test the secured Spring Boot API endpoints using Testcontainers Keycloak and PostgreSQL modules.
weight: 20
---

To test the secured API endpoints, you need a running Keycloak instance and a
PostgreSQL database, plus a started Spring context. Testcontainers spins up both
services in Docker containers and connects them to Spring through dynamic
property registration.

## Configure the test containers

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

## Write the test

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

## Use Testcontainers for local development

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
