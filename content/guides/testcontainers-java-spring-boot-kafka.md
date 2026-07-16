---
title: Testing Spring Boot Kafka Listener using Testcontainers
linkTitle: Spring Boot Kafka
description: Learn how to test a Spring Boot Kafka listener using Testcontainers with Kafka and MySQL modules.
keywords: testcontainers, java, spring boot, testing, kafka, mysql, jpa, awaitility
summary: |
  Learn how to create a Spring Boot application with a Kafka listener that persists data in MySQL,
  then test it using Testcontainers Kafka and MySQL modules with Awaitility.
aliases:
  - /guides/testcontainers-java-spring-boot-kafka/create-project/
  - /guides/testcontainers-java-spring-boot-kafka/run-tests/
  - /guides/testcontainers-java-spring-boot-kafka/write-tests/
params:
  tags: [testing]
  time: 25 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testing-spring-boot-kafka-listener -->

In this guide, you will learn how to:

- Create a Spring Boot application with Kafka integration
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

## Create the Spring Boot project

### Set up the project

Create a Spring Boot project from [Spring Initializr](https://start.spring.io)
by selecting the **Spring for Apache Kafka**, **Spring Data JPA**, **MySQL Driver**,
and **Testcontainers** starters.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-spring-boot-kafka-listener).

After generating the application, add the Awaitility library as a test
dependency. You'll use it later to assert the expectations of an asynchronous
process flow.

The key dependencies in `pom.xml` are:

```xml
<properties>
    <java.version>17</java.version>
    <testcontainers.version>2.0.4</testcontainers.version>
</properties>
<dependencies>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-data-jpa</artifactId>
    </dependency>
    <dependency>
        <groupId>org.springframework.kafka</groupId>
        <artifactId>spring-kafka</artifactId>
    </dependency>
    <dependency>
        <groupId>com.mysql</groupId>
        <artifactId>mysql-connector-j</artifactId>
        <scope>runtime</scope>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-test</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.springframework.kafka</groupId>
        <artifactId>spring-kafka-test</artifactId>
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
    <dependency>
        <groupId>org.awaitility</groupId>
        <artifactId>awaitility</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>
```

Using the Testcontainers BOM (Bill of Materials) is recommended so that you
don't have to repeat the version for every Testcontainers module dependency.

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
class Product {

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

### Create the Spring Data JPA repository

Create a repository interface for the `Product` entity with a method to find a
product by code and a method to update the price for a given product code:

```java
package com.testcontainers.demo;

import java.math.BigDecimal;
import java.util.Optional;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

interface ProductRepository extends JpaRepository<Product, Long> {
  Optional<Product> findByCode(String code);

  @Modifying
  @Query("update Product p set p.price = :price where p.code = :productCode")
  void updateProductPrice(
    @Param("productCode") String productCode,
    @Param("price") BigDecimal price
  );
}
```

### Add a schema creation script

Because the application doesn't use an in-memory database, you need to create
the MySQL tables. The recommended approach for production is a migration tool
like Flyway or Liquibase, but for this guide the built-in Spring Boot schema
initialization is sufficient.

Create `src/main/resources/schema.sql`:

```sql
create table products (
      id int NOT NULL AUTO_INCREMENT,
      code varchar(255) not null,
      name varchar(255) not null,
      price numeric(5,2) not null,
      PRIMARY KEY (id),
      UNIQUE (code)
);
```

Enable schema initialization in `src/main/resources/application.properties`:

```properties
spring.sql.init.mode=always
```

### Create the event payload

Create a record named `ProductPriceChangedEvent` that represents the structure
of the event payload received from the Kafka topic:

```java
package com.testcontainers.demo;

import java.math.BigDecimal;

record ProductPriceChangedEvent(String productCode, BigDecimal price) {}
```

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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

@Component
@Transactional
class ProductPriceChangedEventHandler {

  private static final Logger log = LoggerFactory.getLogger(
    ProductPriceChangedEventHandler.class
  );

  private final ProductRepository productRepository;

  ProductPriceChangedEventHandler(ProductRepository productRepository) {
    this.productRepository = productRepository;
  }

  @KafkaListener(topics = "product-price-changes", groupId = "demo")
  public void handle(ProductPriceChangedEvent event) {
    log.info(
      "Received a ProductPriceChangedEvent with productCode:{}: ",
      event.productCode()
    );
    productRepository.updateProductPrice(event.productCode(), event.price());
  }
}
```

The `@KafkaListener` annotation specifies the topic name to listen to. Spring
Kafka handles serialization and deserialization based on the properties
configured in `application.properties`.

### Configure Kafka serialization

Add the following Kafka properties to
`src/main/resources/application.properties`:

```properties
######## Kafka Configuration  #########
spring.kafka.bootstrap-servers=localhost:9092
spring.kafka.producer.key-serializer=org.apache.kafka.common.serialization.StringSerializer
spring.kafka.producer.value-serializer=org.springframework.kafka.support.serializer.JsonSerializer

spring.kafka.consumer.group-id=demo
spring.kafka.consumer.auto-offset-reset=latest
spring.kafka.consumer.key-deserializer=org.apache.kafka.common.serialization.StringDeserializer
spring.kafka.consumer.value-deserializer=org.springframework.kafka.support.serializer.JsonDeserializer
spring.kafka.consumer.properties.spring.json.trusted.packages=com.testcontainers.demo
```

The `productCode` key is (de)serialized using `StringSerializer`/`StringDeserializer`,
and the `ProductPriceChangedEvent` value is (de)serialized using
`JsonSerializer`/`JsonDeserializer`.

## Write tests with Testcontainers

To test the Kafka listener, you need a running Kafka broker and a MySQL
database, plus a started Spring context. Testcontainers spins up both services
in Docker containers and `@DynamicPropertySource` connects them to Spring.

### Write the test

Create `ProductPriceChangedEventHandlerTest.java`:

```java
package com.testcontainers.demo;

import static java.util.concurrent.TimeUnit.SECONDS;
import static org.assertj.core.api.Assertions.assertThat;
import static org.awaitility.Awaitility.await;

import java.math.BigDecimal;
import java.time.Duration;
import java.util.Optional;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.context.TestPropertySource;
import org.testcontainers.kafka.ConfluentKafkaContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

@SpringBootTest
@TestPropertySource(
  properties = {
    "spring.kafka.consumer.auto-offset-reset=earliest",
    "spring.datasource.url=jdbc:tc:mysql:8.0.32:///db",
  }
)
@Testcontainers
class ProductPriceChangedEventHandlerTest {

  @Container
  static final ConfluentKafkaContainer kafka =
    new ConfluentKafkaContainer("confluentinc/cp-kafka:7.8.0");

  @DynamicPropertySource
  static void overrideProperties(DynamicPropertyRegistry registry) {
    registry.add("spring.kafka.bootstrap-servers", kafka::getBootstrapServers);
  }

  @Autowired
  private KafkaTemplate<String, Object> kafkaTemplate;

  @Autowired
  private ProductRepository productRepository;

  @BeforeEach
  void setUp() {
    Product product = new Product(null, "P100", "Product One", BigDecimal.TEN);
    productRepository.save(product);
  }

  @Test
  void shouldHandleProductPriceChangedEvent() {
    ProductPriceChangedEvent event = new ProductPriceChangedEvent(
      "P100",
      new BigDecimal("14.50")
    );

    kafkaTemplate.send("product-price-changes", event.productCode(), event);

    await()
      .pollInterval(Duration.ofSeconds(3))
      .atMost(10, SECONDS)
      .untilAsserted(() -> {
        Optional<Product> optionalProduct = productRepository.findByCode(
          "P100"
        );
        assertThat(optionalProduct).isPresent();
        assertThat(optionalProduct.get().getCode()).isEqualTo("P100");
        assertThat(optionalProduct.get().getPrice())
          .isEqualTo(new BigDecimal("14.50"));
      });
  }
}
```

Here's what the test does:

- `@SpringBootTest` starts the full Spring application context.
- The Testcontainers special JDBC URL (`jdbc:tc:mysql:8.0.32:///db`) in
  `@TestPropertySource` spins up a MySQL container and configures it as the
  datasource automatically.
- `@Testcontainers` and `@Container` manage the lifecycle of the Kafka
  container. `@DynamicPropertySource` registers the Kafka bootstrap servers
  with Spring so that the producer and consumer connect to the test container.
- `@BeforeEach` creates a `Product` record in the database before each test.
- The test sends a `ProductPriceChangedEvent` to the `product-price-changes`
  topic using `KafkaTemplate`. Spring Boot converts the object to JSON using
  `JsonSerializer`.
- Because Kafka message processing is asynchronous, the test uses
  [Awaitility](http://www.awaitility.org/) to poll every 3 seconds (up to a
  maximum of 10 seconds) until the product price in the database matches the
  expected value.
- The property `spring.kafka.consumer.auto-offset-reset` is set to `earliest`
  so that the listener consumes messages even if they're sent to the topic
  before the listener is ready. This setting is helpful when running tests.

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

- [Getting started with Testcontainers in a Java Spring Boot project](https://testcontainers.com/guides/testing-spring-boot-rest-api-using-testcontainers/)
- [The simplest way to replace H2 with a real database for testing](https://testcontainers.com/guides/replace-h2-with-real-database-for-testing/)
- [Awaitility](http://www.awaitility.org/)
- [Testcontainers Kafka module](https://java.testcontainers.org/modules/kafka/)
- [Testcontainers MySQL module](https://java.testcontainers.org/modules/databases/mysql/)
