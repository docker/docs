---
title: Testing REST API integrations in Micronaut apps using WireMock
linkTitle: Micronaut WireMock
description: Learn how to test REST API integrations in a Micronaut application using the Testcontainers WireMock module.
keywords: testcontainers, java, micronaut, testing, wiremock, rest api
summary: |
  Learn how to create a Micronaut application that integrates with
  external REST APIs, then test those integrations using WireMock
  and the Testcontainers WireMock module.
aliases:
  - /guides/testcontainers-java-micronaut-wiremock/create-project/
  - /guides/testcontainers-java-micronaut-wiremock/run-tests/
  - /guides/testcontainers-java-micronaut-wiremock/write-tests/
params:
  tags: [testing]
  time: 20 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-in-micronaut-apps-using-wiremock -->

In this guide, you'll learn how to:

- Create a Micronaut application that talks to external REST APIs
- Test external API integrations using WireMock
- Use the Testcontainers WireMock module to run WireMock as a Docker container

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Micronaut project

### Set up the project

Create a Micronaut project from [Micronaut Launch](https://micronaut.io/launch)
by selecting the **http-client**, **micronaut-test-rest-assured**, and
**testcontainers** features.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-rest-api-integrations-in-micronaut-apps-using-wiremock).

After generating the project, add the **WireMock** and **Testcontainers
WireMock** libraries as test dependencies. The key dependencies in `pom.xml`
are:

```xml
<parent>
    <groupId>io.micronaut.platform</groupId>
    <artifactId>micronaut-parent</artifactId>
    <version>4.1.2</version>
</parent>

<properties>
    <jdk.version>17</jdk.version>
    <micronaut.version>4.1.2</micronaut.version>
    <micronaut.runtime>netty</micronaut.runtime>
</properties>

<repositories>
    <repository>
        <id>jitpack.io</id>
        <url>https://jitpack.io</url>
    </repository>
</repositories>

<dependencies>
    <dependency>
        <groupId>io.micronaut</groupId>
        <artifactId>micronaut-http-client</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut</groupId>
        <artifactId>micronaut-http-server-netty</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut.serde</groupId>
        <artifactId>micronaut-serde-jackson</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut.test</groupId>
        <artifactId>micronaut-test-junit5</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut.test</groupId>
        <artifactId>micronaut-test-rest-assured</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-junit-jupiter</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.wiremock</groupId>
        <artifactId>wiremock-standalone</artifactId>
        <version>3.2.0</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.wiremock.integrations.testcontainers</groupId>
        <artifactId>wiremock-testcontainers-module</artifactId>
        <version>1.0-alpha-13</version>
        <scope>test</scope>
    </dependency>
</dependencies>
```

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

Create `Album.java` using Java records. Annotate both records with `@Serdeable`
to allow serialization and deserialization:

```java
package com.testcontainers.demo;

import io.micronaut.serde.annotation.Serdeable;
import java.util.List;

@Serdeable
public record Album(Long albumId, List<Photo> photos) {}

@Serdeable
record Photo(Long id, String title, String url, String thumbnailUrl) {}
```

### Create the PhotoServiceClient

Micronaut provides
[declarative HTTP client](https://docs.micronaut.io/latest/guide/#httpClient)
support. Create an interface with a method that fetches photos for a given album
ID:

```java
package com.testcontainers.demo;

import io.micronaut.http.annotation.Get;
import io.micronaut.http.annotation.PathVariable;
import io.micronaut.http.client.annotation.Client;
import java.util.List;

@Client(id = "photosapi")
interface PhotoServiceClient {

    @Get("/albums/{albumId}/photos")
    List<Photo> getPhotos(@PathVariable Long albumId);
}
```

The `@Client(id = "photosapi")` annotation ties this client to a named
configuration. Add the following property to
`src/main/resources/application.properties` to set the base URL:

```properties
micronaut.http.services.photosapi.url=https://jsonplaceholder.typicode.com
```

### Create the REST API endpoint

Create `AlbumController.java`:

```java
package com.testcontainers.demo;

import static io.micronaut.scheduling.TaskExecutors.BLOCKING;

import io.micronaut.http.annotation.Controller;
import io.micronaut.http.annotation.Get;
import io.micronaut.http.annotation.PathVariable;
import io.micronaut.scheduling.annotation.ExecuteOn;

@Controller("/api")
class AlbumController {

    private final PhotoServiceClient photoServiceClient;

    AlbumController(PhotoServiceClient photoServiceClient) {
        this.photoServiceClient = photoServiceClient;
    }

    @ExecuteOn(BLOCKING)
    @Get("/albums/{albumId}")
    public Album getAlbumById(@PathVariable Long albumId) {
        return new Album(albumId, photoServiceClient.getPhotos(albumId));
    }
}
```

Here's what this controller does:

- `@Controller("/api")` maps the controller to the `/api` path.
- Constructor injection provides a `PhotoServiceClient` bean.
- `@ExecuteOn(BLOCKING)` offloads blocking I/O to a separate thread pool so it
  doesn't block the event loop.
- `@Get("/albums/{albumId}")` maps the `getAlbumById()` method to an HTTP GET
  request.

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

### Test with WireMock's JUnit 5 extension

The first approach uses WireMock's `WireMockExtension` to start an in-process
WireMock server on a dynamic port.

Create `AlbumControllerTest.java`:

```java
package com.testcontainers.demo;

import static com.github.tomakehurst.wiremock.client.WireMock.aResponse;
import static com.github.tomakehurst.wiremock.client.WireMock.urlMatching;
import static com.github.tomakehurst.wiremock.core.WireMockConfiguration.wireMockConfig;
import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.hasSize;

import com.github.tomakehurst.wiremock.client.WireMock;
import com.github.tomakehurst.wiremock.junit5.WireMockExtension;
import io.micronaut.context.ApplicationContext;
import io.micronaut.http.MediaType;
import io.micronaut.runtime.server.EmbeddedServer;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import java.util.Collections;
import java.util.Map;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.RegisterExtension;

class AlbumControllerTest {

    @RegisterExtension
    static WireMockExtension wireMock = WireMockExtension.newInstance()
            .options(wireMockConfig().dynamicPort())
            .build();

    private Map<String, Object> getProperties() {
        return Collections.singletonMap("micronaut.http.services.photosapi.url", wireMock.baseUrl());
    }

    @Test
    void shouldGetAlbumById() {
        try (EmbeddedServer server = ApplicationContext.run(EmbeddedServer.class, getProperties())) {
            RestAssured.port = server.getPort();
            Long albumId = 1L;
            String responseJson =
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
            """;
            wireMock.stubFor(WireMock.get(urlMatching("/albums/" + albumId + "/photos"))
                    .willReturn(aResponse()
                            .withHeader("Content-Type", MediaType.APPLICATION_JSON)
                            .withBody(responseJson)));

            given().contentType(ContentType.JSON)
                    .when()
                    .get("/api/albums/{albumId}", albumId)
                    .then()
                    .statusCode(200)
                    .body("albumId", is(albumId.intValue()))
                    .body("photos", hasSize(2));
        }
    }

    @Test
    void shouldReturnServerErrorWhenPhotoServiceCallFailed() {
        try (EmbeddedServer server = ApplicationContext.run(EmbeddedServer.class, getProperties())) {
            RestAssured.port = server.getPort();
            Long albumId = 2L;
            wireMock.stubFor(WireMock.get(urlMatching("/albums/" + albumId + "/photos"))
                    .willReturn(aResponse().withStatus(500)));

            given().contentType(ContentType.JSON)
                    .when()
                    .get("/api/albums/{albumId}", albumId)
                    .then()
                    .statusCode(500);
        }
    }
}
```

Here's what this test does:

- `WireMockExtension` starts a WireMock server on a dynamic port.
- The `getProperties()` method overrides `micronaut.http.services.photosapi.url`
  to point at the WireMock endpoint, so the application talks to WireMock
  instead of the real photo service.
- `shouldGetAlbumById()` configures a mock response for
  `/albums/{albumId}/photos`, sends a request to the application's
  `/api/albums/{albumId}` endpoint, and verifies the response body.
- `shouldReturnServerErrorWhenPhotoServiceCallFailed()` configures WireMock to
  return a 500 status and verifies the application propagates that error.

### Stub using JSON mapping files

Instead of stubbing with the WireMock Java API, you can use JSON mapping-based
configuration.

Create `src/test/resources/wiremock/mappings/get-album-photos.json`:

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

Create `src/test/resources/wiremock/__files/album-photos-resp-200.json`:

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

Then initialize WireMock to load stub mappings from these files:

```java
@RegisterExtension
static WireMockExtension wireMock = WireMockExtension.newInstance()
     .options(
         wireMockConfig()
            .dynamicPort()
            .usingFilesUnderClasspath("wiremock")
    )
    .build();
```

With mapping files-based stubbing in place, write tests without needing
programmatic stubs:

```java
@Test
void shouldGetAlbumById() {
    Long albumId = 1L;
    try (EmbeddedServer server = ApplicationContext.run(EmbeddedServer.class, getProperties())) {
        RestAssured.port = server.getPort();

        given().contentType(ContentType.JSON)
                .when()
                .get("/api/albums/{albumId}", albumId)
                .then()
                .statusCode(200)
                .body("albumId", is(albumId.intValue()))
                .body("photos", hasSize(2));
    }
}
```

### Use the Testcontainers WireMock module

The [Testcontainers WireMock module](https://testcontainers.com/modules/wiremock/)
provisions a WireMock server as a standalone container within your tests, based
on [WireMock Docker](https://github.com/wiremock/wiremock-docker).

Create `src/test/resources/mocks-config.json` with the stub mappings:

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

Create `src/test/resources/album-photos-response.json`:

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
import static org.hamcrest.Matchers.nullValue;

import io.micronaut.context.ApplicationContext;
import io.micronaut.core.annotation.NonNull;
import io.micronaut.runtime.server.EmbeddedServer;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import java.util.Collections;
import java.util.Map;
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.wiremock.integrations.testcontainers.WireMockContainer;

@Testcontainers(disabledWithoutDocker = true)
class AlbumControllerTestcontainersTests {

    @Container
    static WireMockContainer wiremockServer = new WireMockContainer("wiremock/wiremock:2.35.0")
            .withMappingFromResource("mocks-config.json")
            .withFileFromResource("album-photos-response.json");

    @NonNull public Map<String, Object> getProperties() {
        return Collections.singletonMap("micronaut.http.services.photosapi.url", wiremockServer.getBaseUrl());
    }

    @Test
    void shouldGetAlbumById() {
        Long albumId = 1L;
        try (EmbeddedServer server = ApplicationContext.run(EmbeddedServer.class, getProperties())) {
            RestAssured.port = server.getPort();

            given().contentType(ContentType.JSON)
                    .when()
                    .get("/api/albums/{albumId}", albumId)
                    .then()
                    .statusCode(200)
                    .body("albumId", is(albumId.intValue()))
                    .body("photos", hasSize(2));
        }
    }

    @Test
    void shouldReturnServerErrorWhenPhotoServiceCallFailed() {
        Long albumId = 2L;
        try (EmbeddedServer server = ApplicationContext.run(EmbeddedServer.class, getProperties())) {
            RestAssured.port = server.getPort();
            given().contentType(ContentType.JSON)
                    .when()
                    .get("/api/albums/{albumId}", albumId)
                    .then()
                    .statusCode(500);
        }
    }

    @Test
    void shouldReturnEmptyPhotos() {
        Long albumId = 3L;
        try (EmbeddedServer server = ApplicationContext.run(EmbeddedServer.class, getProperties())) {
            RestAssured.port = server.getPort();
            given().contentType(ContentType.JSON)
                    .when()
                    .get("/api/albums/{albumId}", albumId)
                    .then()
                    .statusCode(200)
                    .body("albumId", is(albumId.intValue()))
                    .body("photos", nullValue());
        }
    }
}
```

Here's what this test does:

- `@Testcontainers` and `@Container` annotations start a `WireMockContainer`
  using the `wiremock/wiremock:2.35.0` Docker image.
- `withMappingFromResource("mocks-config.json")` loads stub mappings from the
  classpath resource.
- `withFileFromResource("album-photos-response.json")` makes the response body
  file available to WireMock.
- `getProperties()` overrides the photo service URL to point at the WireMock
  container's base URL.
- `shouldGetAlbumById()` verifies that the application returns the expected
  album with two photos.
- `shouldReturnServerErrorWhenPhotoServiceCallFailed()` verifies that a 500
  from the photo service propagates to the caller.
- `shouldReturnEmptyPhotos()` verifies the application handles an empty photo
  list.

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

You built a Micronaut application that integrates with an external REST API
using declarative HTTP clients, then tested that integration using WireMock and
the Testcontainers WireMock module. Testing at the HTTP protocol level instead
of mocking Java methods lets you catch serialization issues and simulate
realistic failure scenarios.

> [!TIP]
> Testcontainers WireMock modules are available for Go and Python as well.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testcontainers WireMock module](https://testcontainers.com/modules/wiremock/)
- [WireMock documentation](https://wiremock.org/docs/)
- [Testcontainers JUnit 5 quickstart](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Testing REST API integrations in Spring Boot using WireMock](/guides/testcontainers-java-wiremock/)
