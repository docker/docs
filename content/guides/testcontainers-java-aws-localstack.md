---
title: Testing AWS service integrations using LocalStack
linkTitle: AWS LocalStack
description: Learn how to test Spring Cloud AWS applications using LocalStack and Testcontainers.
keywords: testcontainers, java, spring boot, testing, aws, localstack, s3, sqs, spring cloud aws
summary: |
  Learn how to create a Spring Boot application with Spring Cloud AWS,
  then test S3 and SQS integrations using Testcontainers and LocalStack.
aliases:
  - /guides/testcontainers-java-aws-localstack/create-project/
  - /guides/testcontainers-java-aws-localstack/run-tests/
  - /guides/testcontainers-java-aws-localstack/write-tests/
params:
  tags: [testing]
  time: 25 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testing-aws-service-integrations-using-localstack -->

In this guide, you will learn how to:

- Create a Spring Boot application with Spring Cloud AWS integration
- Use AWS S3 and SQS services
- Test the application using Testcontainers and LocalStack

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

### Create the configuration properties

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

### Implement StorageService for S3

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

### Create the SQS message model

Create a `Message` record that represents the payload you send to the SQS
queue:

```java
package com.testcontainers.demo;

import java.util.UUID;

public record Message(UUID uuid, String content) {}
```

### Implement the message sender

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

### Implement the message listener

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

## Write tests with Testcontainers

To test the application, you need a running LocalStack instance that emulates
the AWS S3 and SQS services. Testcontainers spins up LocalStack in a Docker
container and `@DynamicPropertySource` connects it to Spring Cloud AWS.

### Configure the test container

You can start a LocalStack container and configure the Spring Cloud AWS
properties to talk to it instead of actual AWS services. The properties you
need to set are:

```properties
spring.cloud.aws.s3.endpoint=http://localhost:4566
spring.cloud.aws.sqs.endpoint=http://localhost:4566
spring.cloud.aws.credentials.access-key=noop
spring.cloud.aws.credentials.secret-key=noop
spring.cloud.aws.region.static=us-east-1
```

For testing, use an ephemeral container that starts on a random available port
so that you can run multiple builds in CI in parallel without port conflicts.

### Write the test

Create `MessageListenerTest.java`:

```java
package com.testcontainers.demo;

import static org.assertj.core.api.Assertions.assertThat;
import static org.awaitility.Awaitility.await;
import static org.testcontainers.containers.localstack.LocalStackContainer.Service.S3;
import static org.testcontainers.containers.localstack.LocalStackContainer.Service.SQS;

import java.io.IOException;
import java.time.Duration;
import java.util.UUID;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.localstack.LocalStackContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.DockerImageName;

@SpringBootTest
@Testcontainers
class MessageListenerTest {

  @Container
  static LocalStackContainer localStack = new LocalStackContainer(
    DockerImageName.parse("localstack/localstack:3.0")
  );

  static final String BUCKET_NAME = UUID.randomUUID().toString();
  static final String QUEUE_NAME = UUID.randomUUID().toString();

  @DynamicPropertySource
  static void overrideProperties(DynamicPropertyRegistry registry) {
    registry.add("app.bucket", () -> BUCKET_NAME);
    registry.add("app.queue", () -> QUEUE_NAME);
    registry.add(
      "spring.cloud.aws.region.static",
      () -> localStack.getRegion()
    );
    registry.add(
      "spring.cloud.aws.credentials.access-key",
      () -> localStack.getAccessKey()
    );
    registry.add(
      "spring.cloud.aws.credentials.secret-key",
      () -> localStack.getSecretKey()
    );
    registry.add(
      "spring.cloud.aws.s3.endpoint",
      () -> localStack.getEndpointOverride(S3).toString()
    );
    registry.add(
      "spring.cloud.aws.sqs.endpoint",
      () -> localStack.getEndpointOverride(SQS).toString()
    );
  }

  @BeforeAll
  static void beforeAll() throws IOException, InterruptedException {
    localStack.execInContainer("awslocal", "s3", "mb", "s3://" + BUCKET_NAME);
    localStack.execInContainer(
      "awslocal",
      "sqs",
      "create-queue",
      "--queue-name",
      QUEUE_NAME
    );
  }

  @Autowired
  StorageService storageService;

  @Autowired
  MessageSender publisher;

  @Autowired
  ApplicationProperties properties;

  @Test
  void shouldHandleMessageSuccessfully() {
    Message message = new Message(UUID.randomUUID(), "Hello World");
    publisher.publish(properties.queue(), message);

    await()
      .pollInterval(Duration.ofSeconds(2))
      .atMost(Duration.ofSeconds(10))
      .ignoreExceptions()
      .untilAsserted(() -> {
        String msg = storageService.downloadAsString(
          properties.bucket(),
          message.uuid().toString()
        );
        assertThat(msg).isEqualTo("Hello World");
      });
  }
}
```

Here's what the test does:

- `@SpringBootTest` starts the full Spring application context.
- The Testcontainers JUnit 5 annotations `@Testcontainers` and `@Container`
  manage the lifecycle of a `LocalStackContainer` instance.
- `@DynamicPropertySource` obtains the dynamic S3 and SQS endpoint URLs,
  region, access key, and secret key from the container, and registers them as
  Spring Cloud AWS configuration properties.
- `@BeforeAll` creates the required SQS queue and S3 bucket using the
  `awslocal` CLI tool that comes pre-installed in the LocalStack Docker image.
  The `localStack.execInContainer()` API runs commands inside the container.
- `shouldHandleMessageSuccessfully()` publishes a `Message` to the SQS queue.
  The listener receives the message and stores its content in the S3 bucket
  with the UUID as the key. Awaitility waits up to 10 seconds for the expected
  content to appear in the bucket.

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the LocalStack Docker container start and the test pass. After
the tests finish, the container stops and is removed automatically.

### Summary

LocalStack lets you develop and test AWS-based applications locally.
The Testcontainers LocalStack module makes it straightforward to write
integration tests by using ephemeral LocalStack containers that start on random
ports with no external setup required.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testcontainers LocalStack module](https://java.testcontainers.org/modules/localstack/)
- [Getting started with Testcontainers for Java](https://java.testcontainers.org/quickstart/junit_5_quickstart/)
- [Spring Cloud AWS documentation](https://docs.awspring.io/spring-cloud-aws/docs/3.0.3/reference/html/index.html)
