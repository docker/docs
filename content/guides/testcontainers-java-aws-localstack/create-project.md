---
title: Create the Spring Boot project
linkTitle: Create the project
description: Set up a Spring Boot project with Spring Cloud AWS, S3, and SQS.
weight: 10
---

## Set up the project

Create a Spring Boot project from [Spring Initializr](https://start.spring.io)
by selecting the **Testcontainers** starter. Spring Cloud AWS starters are not
available on Spring Initializr, so you need to add them manually.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-aws-service-integrations-using-localstack).

Add the Spring Cloud AWS BOM to your dependency management and add the S3, SQS
starters as dependencies. Testcontainers provides a
[LocalStack module](https://testcontainers.com/modules/localstack/) for testing
AWS service integrations. You also need
[Awaitility](http://www.awaitility.org/) for testing asynchronous SQS
processing.

The key dependencies in `pom.xml` are:

```xml
<properties>
    <java.version>17</java.version>
    <testcontainers.version>2.0.4</testcontainers.version>
    <awspring.version>3.0.3</awspring.version>
</properties>

<dependencies>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
    <dependency>
        <groupId>io.awspring.cloud</groupId>
        <artifactId>spring-cloud-aws-starter-s3</artifactId>
    </dependency>
    <dependency>
        <groupId>io.awspring.cloud</groupId>
        <artifactId>spring-cloud-aws-starter-sqs</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-test</artifactId>
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
        <artifactId>testcontainers-localstack</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.awaitility</groupId>
        <artifactId>awaitility</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>

<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>io.awspring.cloud</groupId>
            <artifactId>spring-cloud-aws-dependencies</artifactId>
            <version>${awspring.version}</version>
            <type>pom</type>
            <scope>import</scope>
        </dependency>
    </dependencies>
</dependencyManagement>
```

## Create the configuration properties

To make the SQS queue and S3 bucket names configurable, create an
`ApplicationProperties` record:

```java
package com.testcontainers.demo;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app")
public record ApplicationProperties(String queue, String bucket) {}
```

Then add `@ConfigurationPropertiesScan` to the main application class so that
Spring automatically scans for `@ConfigurationProperties`-annotated classes and
registers them as beans:

```java
package com.testcontainers.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;

@SpringBootApplication
@ConfigurationPropertiesScan
public class Application {

  public static void main(String[] args) {
    SpringApplication.run(Application.class, args);
  }
}
```

## Implement StorageService for S3

Spring Cloud AWS provides higher-level abstractions like `S3Template` with
convenience methods for uploading and downloading files. Create a
`StorageService` class:

```java
package com.testcontainers.demo;

import io.awspring.cloud.s3.S3Template;
import java.io.IOException;
import java.io.InputStream;
import org.springframework.stereotype.Service;

@Service
public class StorageService {

  private final S3Template s3Template;

  public StorageService(S3Template s3Template) {
    this.s3Template = s3Template;
  }

  public void upload(String bucketName, String key, InputStream stream) {
    this.s3Template.upload(bucketName, key, stream);
  }

  public InputStream download(String bucketName, String key)
    throws IOException {
    return this.s3Template.download(bucketName, key).getInputStream();
  }

  public String downloadAsString(String bucketName, String key)
    throws IOException {
    try (InputStream is = this.download(bucketName, key)) {
      return new String(is.readAllBytes());
    }
  }
}
```

## Create the SQS message model

Create a `Message` record that represents the payload you send to the SQS
queue:

```java
package com.testcontainers.demo;

import java.util.UUID;

public record Message(UUID uuid, String content) {}
```

## Implement the message sender

Create `MessageSender`, which uses `SqsTemplate` to publish messages:

```java
package com.testcontainers.demo;

import io.awspring.cloud.sqs.operations.SqsTemplate;
import org.springframework.stereotype.Service;

@Service
public class MessageSender {

  private final SqsTemplate sqsTemplate;

  public MessageSender(SqsTemplate sqsTemplate) {
    this.sqsTemplate = sqsTemplate;
  }

  public void publish(String queueName, Message message) {
    sqsTemplate.send(to -> to.queue(queueName).payload(message));
  }
}
```

## Implement the message listener

Create `MessageListener` with a handler method annotated with `@SqsListener`.
When a message arrives, the listener uploads the content to an S3 bucket using
the message UUID as the key:

```java
package com.testcontainers.demo;

import io.awspring.cloud.sqs.annotation.SqsListener;
import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;
import org.springframework.stereotype.Service;

@Service
public class MessageListener {

  private final StorageService storageService;
  private final ApplicationProperties properties;

  public MessageListener(
    StorageService storageService,
    ApplicationProperties properties
  ) {
    this.storageService = storageService;
    this.properties = properties;
  }

  @SqsListener(queueNames = { "${app.queue}" })
  public void handle(Message message) {
    String bucketName = this.properties.bucket();
    String key = message.uuid().toString();
    ByteArrayInputStream is = new ByteArrayInputStream(
      message.content().getBytes(StandardCharsets.UTF_8)
    );
    this.storageService.upload(bucketName, key, is);
  }
}
```

The `${app.queue}` expression reads the queue name from application
configuration instead of hard-coding it.
