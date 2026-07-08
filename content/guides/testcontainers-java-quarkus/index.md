---
title: Testing Quarkus applications with Testcontainers
linkTitle: Quarkus
description: Learn how to test a Quarkus REST API using Testcontainers with PostgreSQL, Hibernate ORM with Panache, and REST Assured.
keywords: testcontainers, java, quarkus, testing, postgresql, rest api, rest assured, panache, dev services
summary: |
  Learn how to create a Quarkus REST API with Hibernate ORM with Panache and PostgreSQL,
  then test it using Quarkus Dev Services, Testcontainers, and REST Assured.
aliases:
  - /guides/testcontainers-java-quarkus/create-project/
  - /guides/testcontainers-java-quarkus/run-tests/
  - /guides/testcontainers-java-quarkus/write-tests/
params:
  tags: [testing]
  time: 25 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-testcontainers-in-quarkus-applications -->

In this guide, you'll learn how to:

- Create a Quarkus application with REST API endpoints
- Use Hibernate ORM with Panache and PostgreSQL for persistence
- Test the REST API using Quarkus Dev Services, which uses Testcontainers behind
  the scenes
- Test with services not supported by Dev Services using
  `QuarkusTestResourceLifecycleManager`

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Quarkus project

### Set up the project

Create a Quarkus project from [code.quarkus.io](https://code.quarkus.io/) by
selecting the **RESTEasy Classic**, **RESTEasy Classic Jackson**,
**Hibernate Validator**, **Hibernate ORM with Panache**, **JDBC Driver -
PostgreSQL**, and **Flyway** extensions.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testcontainers-in-quarkus-applications).

The key dependencies in `pom.xml` are:

```xml
<properties>
    <quarkus.platform.version>3.22.3</quarkus.platform.version>
</properties>
<dependencies>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-hibernate-orm-panache</artifactId>
    </dependency>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-flyway</artifactId>
    </dependency>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-hibernate-validator</artifactId>
    </dependency>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-resteasy</artifactId>
    </dependency>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-resteasy-jackson</artifactId>
    </dependency>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-jdbc-postgresql</artifactId>
    </dependency>
    <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-junit5</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>io.rest-assured</groupId>
        <artifactId>rest-assured</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>
```

### Create the JPA entity

Hibernate ORM with Panache supports the Active Record pattern and the
Repository pattern to simplify JPA usage. This guide uses the Active Record
pattern.

Create `Customer.java` by extending `PanacheEntity`. This gives the entity
built-in persistence methods such as `persist()`, `listAll()`, and
`findById()`.

```java
package com.testcontainers.demo;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;

@Entity
@Table(name = "customers")
public class Customer extends PanacheEntity {

    @Column(nullable = false)
    public String name;

    @Column(nullable = false, unique = true)
    public String email;

    public Customer() {}

    public Customer(Long id, String name, String email) {
        this.id = id;
        this.name = name;
        this.email = email;
    }
}
```

### Create the CustomerService CDI bean

Create a `CustomerService` class annotated with `@ApplicationScoped` and
`@Transactional` to handle persistence operations:

```java
package com.testcontainers.demo;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.transaction.Transactional;
import java.util.List;

@ApplicationScoped
@Transactional
public class CustomerService {

    public List<Customer> getAll() {
        return Customer.listAll();
    }

    public Customer create(Customer customer) {
        customer.persist();
        return customer;
    }
}
```

### Add the Flyway database migration script

Create `src/main/resources/db/migration/V1__init_database.sql`:

```sql
create sequence customers_seq start with 1 increment by 50;

create table customers
(
    id    bigint DEFAULT nextval('customers_seq') not null,
    name  varchar                                 not null,
    email varchar                                 not null,
    primary key (id)
);

insert into customers(name, email)
values ('john', 'john@mail.com'),
       ('rambo', 'rambo@mail.com');
```

Enable Flyway migrations in `src/main/resources/application.properties`:

```properties
quarkus.flyway.migrate-at-start=true
```

### Create the REST API endpoints

Create `CustomerResource.java` with endpoints for fetching all customers and
creating a customer:

```java
package com.testcontainers.demo;

import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import java.util.List;

@Path("/api/customers")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class CustomerResource {
    private final CustomerService customerService;

    public CustomerResource(CustomerService customerService) {
        this.customerService = customerService;
    }

    @GET
    public List<Customer> getAllCustomers() {
        return customerService.getAll();
    }

    @POST
    public Response createCustomer(Customer customer) {
        var savedCustomer = customerService.create(customer);
        return Response.status(Response.Status.CREATED).entity(savedCustomer).build();
    }
}
```

## Write tests with Testcontainers

### Quarkus Dev Services

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

### Write tests for the API endpoints

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

### Customize test configuration

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

### Test with services not supported by Dev Services

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

## Run tests and next steps

### Run the tests

```console
$ ./mvnw test
```

Or with Gradle:

```console
$ ./gradlew test
```

You should see the PostgreSQL Docker container start and all tests pass. After
the tests finish, the container stops and is removed automatically.

### Run the application locally

Quarkus Dev Services automatically provisions unconfigured services in
development mode. Start the Quarkus application in dev mode:

```console
$ ./mvnw compile quarkus:dev
```

Or with Gradle:

```console
$ ./gradlew quarkusDev
```

Dev Services starts a PostgreSQL container automatically. If you're running a
PostgreSQL database on your system and want to use that instead, configure the
datasource properties in `src/main/resources/application.properties`:

```properties
quarkus.datasource.jdbc.url=jdbc:postgresql://localhost:5432/postgres
quarkus.datasource.username=postgres
quarkus.datasource.password=postgres
```

When these properties are set explicitly, Dev Services doesn't provision the
database container and instead connects to the configured database.

### Summary

Quarkus Dev Services improves the developer experience by automatically
provisioning the required services using Testcontainers during development and
testing. This guide covered:

- Building a REST API using JAX-RS with Hibernate ORM with Panache
- Testing API endpoints using REST Assured with Dev Services handling database
  provisioning
- Using `QuarkusTestResourceLifecycleManager` for services not supported by Dev
  Services
- Running the application locally with Dev Services

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Quarkus Dev Services overview](https://quarkus.io/guides/dev-services)
- [Quarkus testing guide](https://quarkus.io/guides/getting-started-testing)
- [Testcontainers Postgres module](https://java.testcontainers.org/modules/databases/postgres/)
