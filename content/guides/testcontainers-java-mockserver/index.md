---
title: Testing REST API integrations using MockServer
linkTitle: MockServer
description: Learn how to test REST API integrations in a Spring Boot application using the Testcontainers MockServer module.
keywords: testcontainers, java, spring boot, testing, mockserver, rest api, rest assured
summary: |
  Learn how to create a Spring Boot application that integrates with
  external REST APIs, then test those integrations using Testcontainers
  and MockServer.
aliases:
  - /guides/testcontainers-java-mockserver/create-project/
  - /guides/testcontainers-java-mockserver/run-tests/
  - /guides/testcontainers-java-mockserver/write-tests/
params:
  tags: [testing]
  time: 20 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-using-mockserver -->

In this guide, you will learn how to:

- Create a Spring Boot application that talks to external REST APIs
- Test external API integrations using the Testcontainers MockServer module

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
by selecting the **Spring Web**, **Spring Reactive Web**, and **Testcontainers**
starters.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-using-mockserver).

After generating the project, add the **REST Assured** and **MockServer**
libraries as test dependencies. The key dependencies in `pom.xml` are:

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
        <artifactId>spring-boot-starter-webflux</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-test</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-junit-jupiter</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-mockserver</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.mock-server</groupId>
        <artifactId>mockserver-netty</artifactId>
        <version>5.15.0</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>io.rest-assured</groupId>
        <artifactId>rest-assured</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>
```

Using the Testcontainers BOM (Bill of Materials) is recommended so that you
don't have to repeat the version for every Testcontainers module dependency.

This guide builds an application that manages video albums. A third-party REST
API handles photo assets. For demonstration purposes, the application uses the
publicly available [JSONPlaceholder](https://jsonplaceholder.typicode.com/) API
as a photo service.

The application exposes a `GET /api/albums/{albumId}` endpoint that calls the
photo service to fetch photos for a given album.
[MockServer](https://www.mock-server.com/) is a library for mocking HTTP-based
services. Testcontainers provides a
[MockServer module](https://java.testcontainers.org/modules/mockserver/) that
runs MockServer as a Docker container.

### Create the Album and Photo models

Create `Album.java` using Java records:

```java
package com.testcontainers.demo;

import java.util.List;

public record Album(Long albumId, List<Photo> photos) {}

record Photo(Long id, String title, String url, String thumbnailUrl) {}
```

### Create the PhotoServiceClient interface

Spring Framework 6 introduced
[declarative HTTP client support](https://docs.spring.io/spring-framework/reference/integration/rest-clients.html#rest-http-interface).
Create an interface with a method that fetches photos for a given album ID:

```java
package com.testcontainers.demo;

import java.util.List;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.service.annotation.GetExchange;

interface PhotoServiceClient {
  @GetExchange("/albums/{albumId}/photos")
  List<Photo> getPhotos(@PathVariable Long albumId);
}
```

### Register PhotoServiceClient as a bean

To generate a runtime implementation of `PhotoServiceClient`, register it as a
Spring bean using `HttpServiceProxyFactory`. The factory requires an
`HttpClientAdapter` implementation. Spring Boot provides `WebClientAdapter` as
part of the `spring-webflux` library:

```java
package com.testcontainers.demo;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.support.WebClientAdapter;
import org.springframework.web.service.invoker.HttpServiceProxyFactory;

@Configuration
public class AppConfig {

  @Bean
  public PhotoServiceClient photoServiceClient(
    @Value("${photos.api.base-url}") String photosApiBaseUrl
  ) {
    WebClient client = WebClient.builder().baseUrl(photosApiBaseUrl).build();
    HttpServiceProxyFactory factory = HttpServiceProxyFactory
      .builder(WebClientAdapter.forClient(client))
      .build();
    return factory.createClient(PhotoServiceClient.class);
  }
}
```

The photo service base URL is externalized as a configuration property. Add the
following entry to `src/main/resources/application.properties`:

```properties
photos.api.base-url=https://jsonplaceholder.typicode.com
```

### Create the REST API endpoint

Create `AlbumController.java`:

```java
package com.testcontainers.demo;

import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.reactive.function.client.WebClientResponseException;

@RestController
@RequestMapping("/api")
class AlbumController {

  private static final Logger logger = LoggerFactory.getLogger(
    AlbumController.class
  );

  private final PhotoServiceClient photoServiceClient;

  AlbumController(PhotoServiceClient photoServiceClient) {
    this.photoServiceClient = photoServiceClient;
  }

  @GetMapping("/albums/{albumId}")
  public ResponseEntity<Album> getAlbumById(@PathVariable Long albumId) {
    try {
      List<Photo> photos = photoServiceClient.getPhotos(albumId);
      return ResponseEntity.ok(new Album(albumId, photos));
    } catch (WebClientResponseException e) {
      logger.error("Failed to get photos", e);
      return new ResponseEntity<>(e.getStatusCode());
    }
  }
}
```

This endpoint calls the photo service for a given album ID and returns a
response like:

```json
{
  "albumId": 1,
  "photos": [
    {
      "id": 51,
      "title": "non sunt voluptatem placeat consequuntur rem incidunt",
      "url": "https://via.placeholder.com/600/8e973b",
      "thumbnailUrl": "https://via.placeholder.com/150/8e973b"
    },
    {
      "id": 52,
      "title": "eveniet pariatur quia nobis reiciendis laboriosam ea",
      "url": "https://via.placeholder.com/600/121fa4",
      "thumbnailUrl": "https://via.placeholder.com/150/121fa4"
    }
  ]
}
```

## Write tests with Testcontainers MockServer

Mocking external API interactions at the HTTP protocol level, rather than
mocking Java methods, lets you verify marshalling and unmarshalling behavior and
simulate network issues.

Testcontainers provides a MockServer module that starts a
[MockServer](https://www.mock-server.com/) instance inside a Docker container.
You can then use `MockServerClient` to configure mock expectations.

### Write the test

Create `AlbumControllerTest.java`:

```java
package com.testcontainers.demo;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.hasSize;
import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;
import static org.mockserver.model.JsonBody.json;

import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockserver.client.MockServerClient;
import org.mockserver.model.Header;
import org.mockserver.verify.VerificationTimes;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.mockserver.MockServerContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@Testcontainers
class AlbumControllerTest {

  @LocalServerPort
  private Integer port;

  @Container
  static MockServerContainer mockServerContainer =
    new MockServerContainer("mockserver/mockserver:5.15.0");

  static MockServerClient mockServerClient;

  @DynamicPropertySource
  static void overrideProperties(DynamicPropertyRegistry registry) {
    mockServerClient =
    new MockServerClient(
      mockServerContainer.getHost(),
      mockServerContainer.getServerPort()
    );
    registry.add("photos.api.base-url", mockServerContainer::getEndpoint);
  }

  @BeforeEach
  void setUp() {
    RestAssured.port = port;
    mockServerClient.reset();
  }

  @Test
  void shouldGetAlbumById() {
    Long albumId = 1L;

    mockServerClient
      .when(
        request().withMethod("GET").withPath("/albums/" + albumId + "/photos")
      )
      .respond(
        response()
          .withStatusCode(200)
          .withHeaders(
            new Header("Content-Type", "application/json; charset=utf-8")
          )
          .withBody(
            json(
              """
              [
                   {
                       "id": 1,
                       "title": "accusamus beatae ad facilis cum similique qui sunt",
                       "url": "https://via.placeholder.com/600/92c952",
                       "thumbnailUrl": "https://via.placeholder.com/150/92c952"
                   },
                   {
                       "id": 2,
                       "title": "reprehenderit est deserunt velit ipsam",
                       "url": "https://via.placeholder.com/600/771796",
                       "thumbnailUrl": "https://via.placeholder.com/150/771796"
                   }
               ]
              """
            )
          )
      );

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(200)
      .body("albumId", is(albumId.intValue()))
      .body("photos", hasSize(2));

    verifyMockServerRequest("GET", "/albums/" + albumId + "/photos", 1);
  }

  @Test
  void shouldReturn404StatusWhenAlbumNotFound() {
    Long albumId = 1L;
    mockServerClient
      .when(
        request().withMethod("GET").withPath("/albums/" + albumId + "/photos")
      )
      .respond(response().withStatusCode(404));

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(404);

    verifyMockServerRequest("GET", "/albums/" + albumId + "/photos", 1);
  }

  private void verifyMockServerRequest(String method, String path, int times) {
    mockServerClient.verify(
      request().withMethod(method).withPath(path),
      VerificationTimes.exactly(times)
    );
  }
}
```

Here's what the test does:

- `@SpringBootTest` starts the full application on a random port.
- The `@Testcontainers` and `@Container` annotations start a
  `MockServerContainer` and create a `MockServerClient` connected to it.
- `@DynamicPropertySource` overrides `photos.api.base-url` to point at the
  MockServer endpoint, so the application talks to MockServer instead of the
  real photo service.
- `@BeforeEach` resets the `MockServerClient` before every test so that
  expectations from one test don't affect another.
- `shouldGetAlbumById()` configures a mock response for
  `/albums/{albumId}/photos`, sends a request to the application's
  `/api/albums/{albumId}` endpoint, and verifies the response body. It also
  uses `mockServerClient.verify()` to confirm that the expected API call
  reached MockServer.
- `shouldReturn404StatusWhenAlbumNotFound()` configures MockServer to return a
  404 status and verifies the application propagates that status to the caller.

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the MockServer Docker container start in the console output. It
acts as the photo service, serving mock responses based on the configured
expectations. All tests should pass.

### Summary

You built a Spring Boot application that integrates with an external REST API
using declarative HTTP clients, then tested that integration using the
Testcontainers MockServer module. Testing at the HTTP protocol level instead of
mocking Java methods lets you catch serialization issues and simulate realistic
failure scenarios.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testcontainers MockServer module](https://java.testcontainers.org/modules/mockserver/)
- [MockServer documentation](https://www.mock-server.com/)
- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
