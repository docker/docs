---
title: Configuration of services running in a container
linkTitle: Service configuration (Java)
description: Learn how to configure services running in Testcontainers by copying files and executing commands inside containers.
keywords: testcontainers, java, testing, postgresql, localstack, container configuration
summary: |
  Learn how to initialize and configure Docker containers for testing
  by copying files into containers and executing commands inside them.
aliases:
  - /guides/testcontainers-java-service-configuration/copy-files/
  - /guides/testcontainers-java-service-configuration/exec-in-container/
params:
  tags: [testing]
  time: 15 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-configuration-of-services-running-in-container -->

In this guide, you will learn how to:

- Initialize containers by copying files into them
- Run commands inside running containers using `execInContainer()`
- Set up a PostgreSQL database with SQL scripts
- Create AWS S3 buckets in LocalStack containers

## Prerequisites

- Java 17+
- Your preferred IDE
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Copy files into containers

Sometimes you need to initialize a container by placing files in a specific
location. For example, PostgreSQL runs SQL scripts from
`/docker-entrypoint-initdb.d/` when the container starts.

### Create the initialization script

Create `src/test/resources/init-db.sql`:

```sql
create table customers (
     id bigint not null,
     name varchar not null,
     primary key (id)
);
```

### Copy the file into the container

Use `withCopyFileToContainer()` to copy the SQL script into the container's
init directory:

```java
package com.testcontainers.demo;

import static org.junit.jupiter.api.Assertions.assertFalse;

import java.util.List;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.testcontainers.postgresql.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.MountableFile;

@Testcontainers
class CustomerServiceTest {

  @Container
  static PostgreSQLContainer postgres = new PostgreSQLContainer(
    "postgres:16-alpine"
  )
    .withCopyFileToContainer(
      MountableFile.forClasspathResource("init-db.sql"),
      "/docker-entrypoint-initdb.d/"
    );

  CustomerService customerService;

  @BeforeEach
  void setUp() {
    customerService =
    new CustomerService(
      postgres.getJdbcUrl(),
      postgres.getUsername(),
      postgres.getPassword()
    );
  }

  @Test
  void shouldGetCustomers() {
    customerService.createCustomer(new Customer(1L, "George"));
    customerService.createCustomer(new Customer(2L, "John"));

    List<Customer> customers = customerService.getAllCustomers();
    assertFalse(customers.isEmpty());
  }
}
```

The `withCopyFileToContainer(MountableFile, String)` method copies `init-db.sql`
from the classpath into `/docker-entrypoint-initdb.d/` inside the container.
PostgreSQL executes scripts in that directory automatically at startup.

You can also copy files from any path on the host:

```java
.withCopyFileToContainer(
    MountableFile.forHostPath("/host/path/to/init-db.sql"),
    "/docker-entrypoint-initdb.d/"
);
```

## Execute commands inside containers

Some Docker containers provide CLI tools for performing actions. You can use
`container.execInContainer(String...)` to run any available command inside a
running container.

### Example: Create an S3 bucket in LocalStack

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

### Summary

- Use `withCopyFileToContainer()` to place initialization files inside
  containers before they start.
- Use `execInContainer()` to run commands inside running containers for
  setup tasks like creating buckets, topics, or queues.

### Further reading

- [Getting started with Testcontainers for Java](/guides/testcontainers-java-getting-started/)
- [Testcontainers Postgres module](https://java.testcontainers.org/modules/databases/postgres/)
- [Testcontainers LocalStack module](https://java.testcontainers.org/modules/localstack/)
