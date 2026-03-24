---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Test the Quarkus REST API using Dev Services with Testcontainers, and test with services not supported by Dev Services.
weight: 20
---

## Quarkus Dev Services

Quarkus Dev Services automatically provisions unconfigured services in
development and test mode. When you include an extension and don't configure it,
Quarkus starts the relevant service using
[Testcontainers](https://www.testcontainers.org/) behind the scenes and wires
the application to use that service.

> [!NOTE]
> Dev Services requires a
> [supported Docker environment](https://www.testcontainers.org/supported_docker_environment/).

Quarkus Dev Services supports most commonly used services like SQL databases,
Kafka, RabbitMQ, Redis, and MongoDB. For more information, see the
[Quarkus Dev Services guide](https://quarkus.io/guides/dev-services).

## Write tests for the API endpoints

Test the `GET /api/customers` and `POST /api/customers` endpoints using REST
Assured. The `io.rest-assured:rest-assured` library was already added as a test
dependency when you generated the project.

Create `CustomerResourceTest.java` and annotate it with `@QuarkusTest`. This
bootstraps the application along with the required services using Dev Services.
Because you haven't configured datasource properties, Dev Services automatically
starts a PostgreSQL database using Testcontainers.

```java
package com.testcontainers.demo;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.junit.jupiter.api.Assertions.assertFalse;

import io.quarkus.test.junit.QuarkusTest;
import io.restassured.common.mapper.TypeRef;
import io.restassured.http.ContentType;
import java.util.List;
import org.junit.jupiter.api.Test;

@QuarkusTest
class CustomerResourceTest {

    @Test
    void shouldGetAllCustomers() {
        List<Customer> customers = given().when()
                .get("/api/customers")
                .then()
                .statusCode(200)
                .extract()
                .as(new TypeRef<>() {});
        assertFalse(customers.isEmpty());
    }

    @Test
    void shouldCreateCustomerSuccessfully() {
        Customer customer = new Customer(null, "John", "john@gmail.com");
        given().contentType(ContentType.JSON)
                .body(customer)
                .when()
                .post("/api/customers")
                .then()
                .statusCode(201)
                .body("name", is("John"))
                .body("email", is("john@gmail.com"));
    }
}
```

Here's what the test does:

- `@QuarkusTest` starts the full Quarkus application with Dev Services enabled.
- Dev Services starts a PostgreSQL container using Testcontainers and configures
  the datasource automatically.
- `shouldGetAllCustomers()` calls `GET /api/customers` and verifies that seeded
  data from the Flyway migration is returned.
- `shouldCreateCustomerSuccessfully()` sends a `POST /api/customers` request and
  verifies the response contains the created customer data.

## Customize test configuration

By default, the Quarkus test instance starts on port 8081 and uses a
`postgres:14` Docker image. Customize both by adding these properties to
`src/main/resources/application.properties`:

```properties
quarkus.http.test-port=0
quarkus.datasource.devservices.image-name=postgres:15.2-alpine
```

Setting `quarkus.http.test-port=0` starts the application on a random available
port, avoiding port conflicts. The `devservices.image-name` property lets you
pin the PostgreSQL image to a specific version that matches production.

## Test with services not supported by Dev Services

Your application might use a service that Dev Services doesn't support out of
the box. In that case, use `QuarkusTestResourceLifecycleManager` to start the
service before the Quarkus application starts for testing.

For example, suppose the application uses CockroachDB. First, add the
CockroachDB Testcontainers module dependency:

```xml
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>cockroachdb</artifactId>
    <scope>test</scope>
</dependency>
```

Create a `CockroachDBTestResource` that implements
`QuarkusTestResourceLifecycleManager`:

```java
package com.testcontainers.demo;

import io.quarkus.test.common.QuarkusTestResourceLifecycleManager;
import java.util.HashMap;
import java.util.Map;
import org.testcontainers.containers.CockroachContainer;

public class CockroachDBTestResource implements QuarkusTestResourceLifecycleManager {

    CockroachContainer cockroachdb;

    @Override
    public Map<String, String> start() {
        cockroachdb = new CockroachContainer("cockroachdb/cockroach:v22.2.0");
        cockroachdb.start();
        Map<String, String> conf = new HashMap<>();
        conf.put("quarkus.datasource.jdbc.url", cockroachdb.getJdbcUrl());
        conf.put("quarkus.datasource.username", cockroachdb.getUsername());
        conf.put("quarkus.datasource.password", cockroachdb.getPassword());
        return conf;
    }

    @Override
    public void stop() {
        cockroachdb.stop();
    }
}
```

Use the `CockroachDBTestResource` with `@QuarkusTestResource` in a test class:

```java
package com.testcontainers.demo;

import static io.restassured.RestAssured.given;
import static org.junit.jupiter.api.Assertions.assertFalse;

import io.quarkus.test.common.QuarkusTestResource;
import io.quarkus.test.junit.QuarkusTest;
import io.restassured.common.mapper.TypeRef;
import java.util.List;
import org.junit.jupiter.api.Test;

@QuarkusTest
@QuarkusTestResource(value = CockroachDBTestResource.class, restrictToAnnotatedClass = true)
class CockroachDBTest {

    @Test
    void shouldGetAllCustomers() {
        List<Customer> customers = given().when()
                .get("/api/customers")
                .then()
                .statusCode(200)
                .extract()
                .as(new TypeRef<>() {});
        assertFalse(customers.isEmpty());
    }
}
```

The `restrictToAnnotatedClass = true` attribute ensures the CockroachDB
container only starts when running this specific test class, rather than being
activated for all tests.
