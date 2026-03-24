---
title: Create the Spring Boot project
linkTitle: Create the project
description: Set up a Spring Boot project with Kafka, Spring Data JPA, and MySQL.
weight: 10
---

## Set up the project

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

## Create the JPA entity

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

## Create the Spring Data JPA repository

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

## Add a schema creation script

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

## Create the event payload

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

## Implement the Kafka listener

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

## Configure Kafka serialization

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
