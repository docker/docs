---
title: Testing Micronaut Kafka Listener using Testcontainers
linkTitle: Micronaut Kafka
description: Learn how to test a Micronaut Kafka listener using Testcontainers with Kafka and MySQL modules.
keywords: testcontainers, java, micronaut, testing, kafka, mysql, jpa, awaitility
summary: |
  Learn how to create a Micronaut application with a Kafka listener that persists data in MySQL,
  then test it using Testcontainers Kafka and MySQL modules with Awaitility.
aliases:
  - /guides/testcontainers-java-micronaut-kafka/create-project/
  - /guides/testcontainers-java-micronaut-kafka/run-tests/
  - /guides/testcontainers-java-micronaut-kafka/write-tests/
params:
  tags: [testing]
  time: 25 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testing-micronaut-kafka-listener -->

In this guide, you'll learn how to:

- Create a Micronaut application with Kafka integration
- Implement a Kafka listener and persist data in a MySQL database
- Test the Kafka listener using Testcontainers and Awaitility

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

Create a Micronaut project from [Micronaut Launch](https://micronaut.io/launch) by
selecting the **kafka**, **data-jpa**, **mysql**, **awaitility**, **assertj**, and
**testcontainers** features.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-micronaut-kafka-listener).

You'll use the [Awaitility](http://www.awaitility.org/) library to assert the
expectations of an asynchronous process flow.

The key dependencies in `pom.xml` are:

```xml
<parent>
    <groupId>io.micronaut.platform</groupId>
    <artifactId>micronaut-parent</artifactId>
    <version>4.1.4</version>
</parent>
<dependencies>
    <dependency>
        <groupId>io.micronaut.data</groupId>
        <artifactId>micronaut-data-hibernate-jpa</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut.kafka</groupId>
        <artifactId>micronaut-kafka</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut.serde</groupId>
        <artifactId>micronaut-serde-jackson</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>io.micronaut.sql</groupId>
        <artifactId>micronaut-jdbc-hikari</artifactId>
        <scope>compile</scope>
    </dependency>
    <dependency>
        <groupId>mysql</groupId>
        <artifactId>mysql-connector-java</artifactId>
        <scope>runtime</scope>
    </dependency>
    <dependency>
        <groupId>org.awaitility</groupId>
        <artifactId>awaitility</artifactId>
        <version>4.2.0</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-junit-jupiter</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-kafka</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-mysql</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>
```

The Micronaut parent POM manages the Testcontainers BOM, so you don't need to
specify versions for Testcontainers modules individually.

### Create the JPA entity

The application listens to a topic called `product-price-changes`. When a
message arrives, it extracts the product code and price from the event payload
and updates the price for that product in the MySQL database.

Create `Product.java`:

```java
package com.testcontainers.demo;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import java.math.BigDecimal;

@Entity
@Table(name = "products")
public class Product {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, unique = true)
    private String code;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    private BigDecimal price;

    public Product() {}

    public Product(Long id, String code, String name, BigDecimal price) {
        this.id = id;
        this.code = code;
        this.name = name;
        this.price = price;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public BigDecimal getPrice() {
        return price;
    }

    public void setPrice(BigDecimal price) {
        this.price = price;
    }
}
```

### Create the Micronaut Data JPA repository

Create a repository interface for the `Product` entity with a method to find a
product by code and a method to update the price for a given product code:

```java
package com.testcontainers.demo;

import io.micronaut.data.annotation.Query;
import io.micronaut.data.annotation.Repository;
import io.micronaut.data.jpa.repository.JpaRepository;
import java.math.BigDecimal;
import java.util.Optional;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {

    Optional<Product> findByCode(String code);

    @Query("update Product p set p.price = :price where p.code = :productCode")
    void updateProductPrice(String productCode, BigDecimal price);
}
```

Unlike Spring Data JPA, Micronaut Data uses compile-time annotation processing
to implement repository methods, avoiding runtime reflection.

### Create the event payload

Create a record named `ProductPriceChangedEvent` that represents the structure
of the event payload received from the Kafka topic:

```java
package com.testcontainers.demo;

import io.micronaut.serde.annotation.Serdeable;
import java.math.BigDecimal;

@Serdeable
public record ProductPriceChangedEvent(String productCode, BigDecimal price) {}
```

The `@Serdeable` annotation tells Micronaut Serialization that this type can be
serialized and deserialized.

The sender and receiver agree on the following JSON format:

```json
{
  "productCode": "P100",
  "price": 25.0
}
```

### Implement the Kafka listener

Create `ProductPriceChangedEventHandler.java`, which handles messages from the
`product-price-changes` topic and updates the product price in the database:

```java
package com.testcontainers.demo;

import static io.micronaut.configuration.kafka.annotation.OffsetReset.EARLIEST;

import io.micronaut.configuration.kafka.annotation.KafkaListener;
import io.micronaut.configuration.kafka.annotation.Topic;
import jakarta.inject.Singleton;
import jakarta.transaction.Transactional;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Singleton
@Transactional
class ProductPriceChangedEventHandler {

    private static final Logger LOG = LoggerFactory.getLogger(ProductPriceChangedEventHandler.class);

    private final ProductRepository productRepository;

    ProductPriceChangedEventHandler(ProductRepository productRepository) {
        this.productRepository = productRepository;
    }

    @Topic("product-price-changes")
    @KafkaListener(offsetReset = EARLIEST, groupId = "demo")
    public void handle(ProductPriceChangedEvent event) {
        LOG.info("Received a ProductPriceChangedEvent with productCode:{}: ", event.productCode());
        productRepository.updateProductPrice(event.productCode(), event.price());
    }
}
```

Key details:

- The `@KafkaListener` annotation marks this class as a Kafka message listener.
  Setting `offsetReset` to `EARLIEST` makes the listener start consuming
  messages from the beginning of the partition, which is useful during testing.
- The `@Topic` annotation specifies which topic to subscribe to.
- Micronaut handles JSON deserialization of the `ProductPriceChangedEvent`
  automatically using Micronaut Serialization.

### Configure the datasource

Add the following properties to `src/main/resources/application.properties`:

```properties
micronaut.application.name=tc-guide-testing-micronaut-kafka-listener
datasources.default.db-type=mysql
datasources.default.dialect=MYSQL
jpa.default.properties.hibernate.hbm2ddl.auto=update
jpa.default.entity-scan.packages=com.testcontainers.demo
datasources.default.driver-class-name=com.mysql.cj.jdbc.Driver
```

Hibernate's `hbm2ddl.auto=update` creates and updates the database schema
automatically. For testing, you'll override this to `create-drop` in the test
properties file.

Create `src/test/resources/application-test.properties`:

```properties
jpa.default.properties.hibernate.hbm2ddl.auto=create-drop
```

## Write tests with Testcontainers

To test the Kafka listener, you need a running Kafka broker and a MySQL
database, plus a started Micronaut application context. Testcontainers spins up
both services in Docker containers and the `TestPropertyProvider` interface
connects them to Micronaut.

### Create a Kafka client for testing

First, create a `@KafkaClient` interface to publish events in the test:

```java
package com.testcontainers.demo;

import io.micronaut.configuration.kafka.annotation.KafkaClient;
import io.micronaut.configuration.kafka.annotation.KafkaKey;
import io.micronaut.configuration.kafka.annotation.Topic;

@KafkaClient
public interface ProductPriceChangesClient {

    @Topic("product-price-changes")
    void send(@KafkaKey String productCode, ProductPriceChangedEvent event);
}
```

Key details:

- The `@KafkaClient` annotation designates this interface as a Kafka producer.
- The `@Topic` annotation specifies the target topic.
- The `@KafkaKey` annotation marks the parameter used as the Kafka message key.
  If no such parameter exists, Micronaut sends the record with a null key.

### Write the test

Create `ProductPriceChangedEventHandlerTest.java`:

```java
package com.testcontainers.demo;

import static java.util.concurrent.TimeUnit.SECONDS;
import static org.assertj.core.api.Assertions.assertThat;
import static org.awaitility.Awaitility.await;

import io.micronaut.context.annotation.Property;
import io.micronaut.core.annotation.NonNull;
import io.micronaut.test.extensions.junit5.annotation.MicronautTest;
import io.micronaut.test.support.TestPropertyProvider;
import java.math.BigDecimal;
import java.time.Duration;
import java.util.Collections;
import java.util.Map;
import java.util.Optional;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.testcontainers.kafka.ConfluentKafkaContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

@MicronautTest(transactional = false)
@Property(name = "datasources.default.driver-class-name", value = "org.testcontainers.jdbc.ContainerDatabaseDriver")
@Property(name = "datasources.default.url", value = "jdbc:tc:mysql:8.0.32:///db")
@Testcontainers(disabledWithoutDocker = true)
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class ProductPriceChangedEventHandlerTest implements TestPropertyProvider {

    @Container
    static final ConfluentKafkaContainer kafka = new ConfluentKafkaContainer("confluentinc/cp-kafka:7.8.0");

    @Override
    public @NonNull Map<String, String> getProperties() {
        if (!kafka.isRunning()) {
            kafka.start();
        }
        return Collections.singletonMap("kafka.bootstrap.servers", kafka.getBootstrapServers());
    }

    @Test
    void shouldHandleProductPriceChangedEvent(
            ProductPriceChangesClient productPriceChangesClient, ProductRepository productRepository) {
        Product product = new Product(null, "P100", "Product One", BigDecimal.TEN);
        Long id = productRepository.save(product).getId();

        ProductPriceChangedEvent event = new ProductPriceChangedEvent("P100", new BigDecimal("14.50"));

        productPriceChangesClient.send(event.productCode(), event);

        await().pollInterval(Duration.ofSeconds(3)).atMost(10, SECONDS).untilAsserted(() -> {
            Optional<Product> optionalProduct = productRepository.findByCode("P100");
            assertThat(optionalProduct).isPresent();
            assertThat(optionalProduct.get().getCode()).isEqualTo("P100");
            assertThat(optionalProduct.get().getPrice()).isEqualTo(new BigDecimal("14.50"));
        });

        productRepository.deleteById(id);
    }
}
```

Here's what the test does:

- `@MicronautTest` initializes the Micronaut application context and the
  embedded server. Setting `transactional` to `false` prevents each test method
  from running inside a rolled-back transaction, which is necessary because the
  Kafka listener processes messages in a separate thread.
- The `@Property` annotations override the datasource driver and URL to use the
  Testcontainers special JDBC URL (`jdbc:tc:mysql:8.0.32:///db`). This spins up
  a MySQL container and configures it as the datasource automatically.
- `@Testcontainers` and `@Container` manage the Kafka container lifecycle.
  The `TestPropertyProvider` interface registers the Kafka bootstrap servers
  with Micronaut so that the producer and consumer connect to the test container.
- `@TestInstance(TestInstance.Lifecycle.PER_CLASS)` creates a single test
  instance for all test methods, which is required when implementing
  `TestPropertyProvider`.
- The test creates a `Product` record in the database, then sends a
  `ProductPriceChangedEvent` to the `product-price-changes` topic using the
  `ProductPriceChangesClient`.
- Because Kafka message processing is asynchronous, the test uses
  [Awaitility](http://www.awaitility.org/) to poll every 3 seconds (up to a
  maximum of 10 seconds) until the product price in the database matches the
  expected value.

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the Kafka and MySQL Docker containers start and all tests pass.
After the tests finish, the containers stop and are removed automatically.

### Summary

Testing with real Kafka and MySQL instances gives you more confidence in the
correctness of your code than mocks or embedded alternatives. The
Testcontainers library manages the container lifecycle so that your integration
tests run against the same services you use in production.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testing REST API integrations in Micronaut apps using WireMock](/guides/testcontainers-java-micronaut-wiremock/)
- [Testing Spring Boot Kafka Listener using Testcontainers](/guides/testcontainers-java-spring-boot-kafka/)
- [Getting started with Testcontainers in a Java Spring Boot project](https://testcontainers.com/guides/testing-spring-boot-rest-api-using-testcontainers/)
- [Awaitility](http://www.awaitility.org/)
- [Testcontainers Kafka module](https://java.testcontainers.org/modules/kafka/)
- [Testcontainers MySQL module](https://java.testcontainers.org/modules/databases/mysql/)
