---
title: Create the Micronaut project
linkTitle: Create the project
description: Set up a Micronaut project with Kafka, Micronaut Data JPA, and MySQL.
weight: 10
---

## Set up the project

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

## Create the Micronaut Data JPA repository

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

## Create the event payload

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

## Implement the Kafka listener

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

## Configure the datasource

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
