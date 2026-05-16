---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Test the Micronaut Kafka listener using Testcontainers Kafka and MySQL modules with Awaitility.
weight: 20
---

To test the Kafka listener, you need a running Kafka broker and a MySQL
database, plus a started Micronaut application context. Testcontainers spins up
both services in Docker containers and the `TestPropertyProvider` interface
connects them to Micronaut.

## Create a Kafka client for testing

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

## Write the test

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
