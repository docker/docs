---
title: Create the Micronaut project
linkTitle: Create the project
description: Set up a Micronaut project with an external REST API integration using declarative HTTP clients.
weight: 10
---

## Set up the project

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

## Create the Album and Photo models

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

## Create the PhotoServiceClient

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

## Create the REST API endpoint

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
