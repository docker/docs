---
title: Create the Quarkus project
linkTitle: Create the project
description: Set up a Quarkus project with Hibernate ORM with Panache, PostgreSQL, Flyway, and REST API endpoints.
weight: 10
---

## Set up the project

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

## Create the JPA entity

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

## Create the CustomerService CDI bean

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

## Add the Flyway database migration script

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

## Create the REST API endpoints

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
