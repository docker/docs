---
title: Execute commands inside containers
linkTitle: Execute commands
description: Run commands inside running containers to initialize services for testing.
weight: 20
---

Some Docker containers provide CLI tools for performing actions. You can use
`container.execInContainer(String...)` to run any available command inside a
running container.

## Example: Create an S3 bucket in LocalStack

The [LocalStack](https://localstack.cloud/) module emulates AWS services. To
test S3 file uploads, create a bucket before running tests:

```java
package com.testcontainers.demo;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.testcontainers.containers.localstack.LocalStackContainer.Service.S3;

import java.io.IOException;
import java.net.URI;
import java.util.List;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.testcontainers.containers.localstack.LocalStackContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.DockerImageName;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.auth.credentials.StaticCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.Bucket;

@Testcontainers
class LocalStackTest {

  static final String bucketName = "mybucket";

  @Container
  static LocalStackContainer localStack = new LocalStackContainer(
    DockerImageName.parse("localstack/localstack:3.4.0")
  );

  @BeforeAll
  static void beforeAll() throws IOException, InterruptedException {
    localStack.execInContainer("awslocal", "s3", "mb", "s3://" + bucketName);

    org.testcontainers.containers.Container.ExecResult execResult =
      localStack.execInContainer("awslocal", "s3", "ls");
    String stdout = execResult.getStdout();
    int exitCode = execResult.getExitCode();
    assertTrue(stdout.contains(bucketName));
    assertEquals(0, exitCode);
  }

  @Test
  void shouldListBuckets() {
    URI s3Endpoint = localStack.getEndpointOverride(S3);
    StaticCredentialsProvider credentialsProvider =
      StaticCredentialsProvider.create(
        AwsBasicCredentials.create(
          localStack.getAccessKey(),
          localStack.getSecretKey()
        )
      );
    S3Client s3 = S3Client
      .builder()
      .endpointOverride(s3Endpoint)
      .credentialsProvider(credentialsProvider)
      .region(Region.of(localStack.getRegion()))
      .build();

    List<String> s3Buckets = s3
      .listBuckets()
      .buckets()
      .stream()
      .map(Bucket::name)
      .toList();

    assertTrue(s3Buckets.contains(bucketName));
  }
}
```

The `execInContainer("awslocal", "s3", "mb", "s3://mybucket")` call runs the
`awslocal` CLI tool (provided by the LocalStack image) to create an S3 bucket.

You can capture the output and exit code from any command:

```java
Container.ExecResult execResult =
    localStack.execInContainer("awslocal", "s3", "ls");
String stdout = execResult.getStdout();
int exitCode = execResult.getExitCode();
```

> [!NOTE]
> The `withCopyFileToContainer()` and `execInContainer()` methods are inherited
> from `GenericContainer`, so they're available for all Testcontainers modules.

## Summary

- Use `withCopyFileToContainer()` to place initialization files inside
  containers before they start.
- Use `execInContainer()` to run commands inside running containers for
  setup tasks like creating buckets, topics, or queues.

## Further reading

- [Getting started with Testcontainers for Java](/guides/testcontainers-java-getting-started/)
- [Testcontainers Postgres module](https://java.testcontainers.org/modules/databases/postgres/)
- [Testcontainers LocalStack module](https://java.testcontainers.org/modules/localstack/)
