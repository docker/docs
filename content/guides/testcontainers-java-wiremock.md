---
title: Testing REST API integrations using WireMock
linkTitle: WireMock
description: Learn how to test REST API integrations in a Spring Boot application using the Testcontainers WireMock module.
keywords: testcontainers, java, spring boot, testing, wiremock, rest api, rest assured
summary: |
  Learn how to create a Spring Boot application that integrates with
  external REST APIs, then test those integrations using Testcontainers
  and WireMock.
aliases:
  - /guides/testcontainers-java-wiremock/create-project/
  - /guides/testcontainers-java-wiremock/run-tests/
  - /guides/testcontainers-java-wiremock/write-tests/
params:
  tags: [testing]
  time: 20 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-using-wiremock -->

In this guide, you'll learn how to:

- Create a Spring Boot application that talks to external REST APIs
- Test external API integrations using WireMock with both the JUnit 5 extension
  and the Testcontainers WireMock module

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
by selecting the **Spring Web** and **Testcontainers** starters.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-using-wiremock).

After generating the project, add the **REST Assured**, **WireMock**, and
**WireMock Testcontainers module** libraries as test dependencies. The key
dependencies in `pom.xml` are:

```xml
<properties>
    <java.version>17</java.version>
    <testcontainers.version>2.0.4</testcontainers.version>
    <wiremock-testcontainers.version>1.0-alpha-13</wiremock-testcontainers.version>
</properties>
<dependencies>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
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
        <groupId>org.wiremock</groupId>
        <artifactId>wiremock-standalone</artifactId>
        <version>3.6.0</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.wiremock.integrations.testcontainers</groupId>
        <artifactId>wiremock-testcontainers-module</artifactId>
        <version>${wiremock-testcontainers.version}</version>
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
[WireMock](https://wiremock.org/) is a tool for building mock APIs.
Testcontainers provides a
[WireMock module](https://testcontainers.com/modules/wiremock/) that runs
WireMock as a Docker container.

### Create the Album and Photo models

Create `Album.java` using Java records:

```java
package com.testcontainers.demo;

import java.util.List;

public record Album(Long albumId, List<Photo> photos) {}

record Photo(Long id, String title, String url, String thumbnailUrl) {}
```

### Create the PhotoServiceClient

Create `PhotoServiceClient.java`, which uses `RestTemplate` to fetch photos for
a given album ID:

```java
package com.testcontainers.demo;

import java.util.List;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
class PhotoServiceClient {

  private final String baseUrl;
  private final RestTemplate restTemplate;

  PhotoServiceClient(
    @Value("${photos.api.base-url}") String baseUrl,
    RestTemplateBuilder builder
  ) {
    this.baseUrl = baseUrl;
    this.restTemplate = builder.build();
  }

  List<Photo> getPhotos(Long albumId) {
    String url = baseUrl + "/albums/{albumId}/photos";
    ResponseEntity<List<Photo>> response = restTemplate.exchange(
      url,
      HttpMethod.GET,
      null,
      new ParameterizedTypeReference<>() {},
      albumId
    );
    return response.getBody();
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
import org.springframework.web.client.RestClientResponseException;

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
    } catch (RestClientResponseException e) {
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

## Write tests with WireMock and Testcontainers

Mocking external API interactions at the HTTP protocol level, rather than
mocking Java methods, lets you verify marshalling and unmarshalling behavior and
simulate network issues.

### Test using WireMock JUnit 5 extension

WireMock provides a JUnit 5 extension that starts an in-process WireMock server.
You can configure stub responses using the WireMock Java API.

Create `AlbumControllerTest.java`:

```java
package com.testcontainers.demo;

import static com.github.tomakehurst.wiremock.client.WireMock.aResponse;
import static com.github.tomakehurst.wiremock.client.WireMock.urlMatching;
import static com.github.tomakehurst.wiremock.core.WireMockConfiguration.wireMockConfig;
import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.hasSize;
import static org.springframework.boot.test.context.SpringBootTest.WebEnvironment.RANDOM_PORT;

import com.github.tomakehurst.wiremock.client.WireMock;
import com.github.tomakehurst.wiremock.junit5.WireMockExtension;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.RegisterExtension;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.http.MediaType;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;

@SpringBootTest(webEnvironment = RANDOM_PORT)
class AlbumControllerTest {

  @LocalServerPort
  private Integer port;

  @RegisterExtension
  static WireMockExtension wireMock = WireMockExtension
    .newInstance()
    .options(wireMockConfig().dynamicPort())
    .build();

  @DynamicPropertySource
  static void configureProperties(DynamicPropertyRegistry registry) {
    registry.add("photos.api.base-url", wireMock::baseUrl);
  }

  @BeforeEach
  void setUp() {
    RestAssured.port = port;
  }

  @Test
  void shouldGetAlbumById() {
    Long albumId = 1L;

    wireMock.stubFor(
      WireMock
        .get(urlMatching("/albums/" + albumId + "/photos"))
        .willReturn(
          aResponse()
            .withHeader("Content-Type", MediaType.APPLICATION_JSON_VALUE)
            .withBody(
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
  }

  @Test
  void shouldReturnServerErrorWhenPhotoServiceCallFailed() {
    Long albumId = 2L;
    wireMock.stubFor(
      WireMock
        .get(urlMatching("/albums/" + albumId + "/photos"))
        .willReturn(aResponse().withStatus(500))
    );

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(500);
  }
}
```

Here's what the test does:

- `@SpringBootTest` starts the full application on a random port.
- `@RegisterExtension` creates a `WireMockExtension` that starts WireMock on a
  dynamic port.
- `@DynamicPropertySource` overrides `photos.api.base-url` to point at the
  WireMock endpoint, so the application talks to WireMock instead of the real
  photo service.
- `shouldGetAlbumById()` configures a stub response for
  `/albums/{albumId}/photos`, sends a request to the application's
  `/api/albums/{albumId}` endpoint, and verifies the response body.
- `shouldReturnServerErrorWhenPhotoServiceCallFailed()` configures WireMock to
  return a 500 status and verifies that the application propagates that status to
  the caller.

### Stub using JSON mapping files

Instead of using the WireMock Java API, you can configure stubs with JSON
mapping files. Create
`src/test/resources/wiremock/mappings/get-album-photos.json`:

```json
{
  "mappings": [
    {
      "request": {
        "method": "GET",
        "urlPattern": "/albums/([0-9]+)/photos"
      },
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "bodyFileName": "album-photos-resp-200.json"
      }
    },
    {
      "request": {
        "method": "GET",
        "urlPattern": "/albums/2/photos"
      },
      "response": {
        "status": 500,
        "headers": {
          "Content-Type": "application/json"
        }
      }
    },
    {
      "request": {
        "method": "GET",
        "urlPattern": "/albums/3/photos"
      },
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "jsonBody": []
      }
    }
  ]
}
```

Create the response body file at
`src/test/resources/wiremock/__files/album-photos-resp-200.json`:

```json
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
```

Initialize WireMock to load stubs from the mapping files:

```java
@RegisterExtension
static WireMockExtension wireMockServer = WireMockExtension
  .newInstance()
  .options(
    wireMockConfig().dynamicPort().usingFilesUnderClasspath("wiremock")
  )
  .build();
```

With mapping-based stubs in place, create
`AlbumControllerWireMockMappingTests.java`:

```java
package com.testcontainers.demo;

import static com.github.tomakehurst.wiremock.core.WireMockConfiguration.wireMockConfig;
import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.hasSize;
import static org.springframework.boot.test.context.SpringBootTest.WebEnvironment.RANDOM_PORT;

import com.github.tomakehurst.wiremock.junit5.WireMockExtension;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.RegisterExtension;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;

@SpringBootTest(webEnvironment = RANDOM_PORT)
class AlbumControllerWireMockMappingTests {

  @LocalServerPort
  private Integer port;

  @RegisterExtension
  static WireMockExtension wireMockServer = WireMockExtension
    .newInstance()
    .options(
      wireMockConfig().dynamicPort().usingFilesUnderClasspath("wiremock")
    )
    .build();

  @DynamicPropertySource
  static void configureProperties(DynamicPropertyRegistry registry) {
    registry.add("photos.api.base-url", wireMockServer::baseUrl);
  }

  @BeforeEach
  void setUp() {
    RestAssured.port = port;
  }

  @Test
  void shouldGetAlbumById() {
    Long albumId = 1L;

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(200)
      .body("albumId", is(albumId.intValue()))
      .body("photos", hasSize(2));
  }

  @Test
  void shouldReturnServerErrorWhenPhotoServiceCallFailed() {
    Long albumId = 2L;

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(500);
  }

  @Test
  void shouldReturnEmptyPhotos() {
    Long albumId = 3L;

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(200)
      .body("albumId", is(albumId.intValue()))
      .body("photos", hasSize(0));
  }
}
```

The tests don't need inline stub definitions because WireMock loads the mappings
automatically from the classpath.

### Test using the Testcontainers WireMock module

The [Testcontainers WireMock module](https://testcontainers.com/modules/wiremock/)
provisions WireMock as a standalone Docker container, based on
[WireMock Docker](https://github.com/wiremock/wiremock-docker). This approach is
useful when you want complete isolation between the test JVM and the mock server.

Create a mock configuration file at
`src/test/resources/com/testcontainers/demo/AlbumControllerTestcontainersTests/mocks-config.json`:

```json
{
  "mappings": [
    {
      "request": {
        "method": "GET",
        "urlPattern": "/albums/([0-9]+)/photos"
      },
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "bodyFileName": "album-photos-response.json"
      }
    },
    {
      "request": {
        "method": "GET",
        "urlPattern": "/albums/2/photos"
      },
      "response": {
        "status": 500,
        "headers": {
          "Content-Type": "application/json"
        }
      }
    },
    {
      "request": {
        "method": "GET",
        "urlPattern": "/albums/3/photos"
      },
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "jsonBody": []
      }
    }
  ]
}
```

Create the response body file at
`src/test/resources/com/testcontainers/demo/AlbumControllerTestcontainersTests/album-photos-response.json`:

```json
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
```

Create `AlbumControllerTestcontainersTests.java`:

```java
package com.testcontainers.demo;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.hasSize;
import static org.springframework.boot.test.context.SpringBootTest.WebEnvironment.RANDOM_PORT;

import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.wiremock.integrations.testcontainers.WireMockContainer;

@SpringBootTest(webEnvironment = RANDOM_PORT)
@Testcontainers
class AlbumControllerTestcontainersTests {

  @LocalServerPort
  private Integer port;

  @Container
  static WireMockContainer wiremockServer = new WireMockContainer(
    "wiremock/wiremock:3.6.0"
  )
    .withMapping(
      "photos-by-album",
      AlbumControllerTestcontainersTests.class,
      "mocks-config.json"
    )
    .withFileFromResource(
      "album-photos-response.json",
      AlbumControllerTestcontainersTests.class,
      "album-photos-response.json"
    );

  @DynamicPropertySource
  static void configureProperties(DynamicPropertyRegistry registry) {
    registry.add("photos.api.base-url", wiremockServer::getBaseUrl);
  }

  @BeforeEach
  void setUp() {
    RestAssured.port = port;
  }

  @Test
  void shouldGetAlbumById() {
    Long albumId = 1L;

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(200)
      .body("albumId", is(albumId.intValue()))
      .body("photos", hasSize(2));
  }

  @Test
  void shouldReturnServerErrorWhenPhotoServiceCallFailed() {
    Long albumId = 2L;

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(500);
  }

  @Test
  void shouldReturnEmptyPhotos() {
    Long albumId = 3L;

    given()
      .contentType(ContentType.JSON)
      .when()
      .get("/api/albums/{albumId}", albumId)
      .then()
      .statusCode(200)
      .body("albumId", is(albumId.intValue()))
      .body("photos", hasSize(0));
  }
}
```

Here's what the test does:

- The `@Testcontainers` and `@Container` annotations start a
  `WireMockContainer` using the `wiremock/wiremock:3.6.0` Docker image.
- `withMapping()` loads stub mappings from `mocks-config.json`, and
  `withFileFromResource()` loads the response body file.
- `@DynamicPropertySource` overrides `photos.api.base-url` to point at the
  WireMock container's base URL.
- The tests don't contain inline stub definitions because WireMock loads them
  from the JSON configuration files.

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the WireMock Docker container start in the console output. It
acts as the photo service, serving mock responses based on the configured
expectations. All tests should pass.

### Summary

You built a Spring Boot application that integrates with an external REST API,
then tested that integration using three different approaches:

- WireMock JUnit 5 extension with inline stubs
- WireMock JUnit 5 extension with JSON mapping files
- Testcontainers WireMock module running WireMock in a Docker container

Testing at the HTTP protocol level instead of mocking Java methods lets you
catch serialization issues and simulate realistic failure scenarios.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testcontainers WireMock module](https://testcontainers.com/modules/wiremock/)
- [WireMock documentation](https://wiremock.org/docs/)
- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
