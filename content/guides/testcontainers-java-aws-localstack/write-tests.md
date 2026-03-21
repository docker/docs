---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Test Spring Cloud AWS S3 and SQS integration using Testcontainers and LocalStack.
weight: 20
---

To test the application, you need a running LocalStack instance that emulates
the AWS S3 and SQS services. Testcontainers spins up LocalStack in a Docker
container and `@DynamicPropertySource` connects it to Spring Cloud AWS.

## Configure the test container

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

## Write the test

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
