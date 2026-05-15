---
title: Write tests with WireMock and Testcontainers
linkTitle: Write tests
description: Test external REST API integrations using WireMock with both the JUnit 5 extension and the Testcontainers WireMock module.
weight: 20
---

Mocking external API interactions at the HTTP protocol level, rather than
mocking Java methods, lets you verify marshalling and unmarshalling behavior and
simulate network issues.

## Test using WireMock JUnit 5 extension

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

## Stub using JSON mapping files

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

## Test using the Testcontainers WireMock module

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
