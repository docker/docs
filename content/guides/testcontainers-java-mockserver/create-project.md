---
title: Create the Spring Boot project
linkTitle: Create the project
description: Set up a Spring Boot project with an external REST API integration using declarative HTTP clients.
weight: 10
---

## Set up the project

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

## Create the Album and Photo models

Create `Album.java` using Java records:

```java
package com.testcontainers.demo;

import java.util.List;

public record Album(Long albumId, List<Photo> photos) {}

record Photo(Long id, String title, String url, String thumbnailUrl) {}
```

## Create the PhotoServiceClient interface

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

## Register PhotoServiceClient as a bean

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

## Create the REST API endpoint

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
