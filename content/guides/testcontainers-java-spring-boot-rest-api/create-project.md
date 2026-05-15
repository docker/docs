---
title: Create the Spring Boot project
linkTitle: Create the project
description: Set up a Spring Boot project with Spring Data JPA, PostgreSQL, and a REST API.
weight: 10
---

## Set up the project

Create a Spring Boot project from [Spring Initializr](https://start.spring.io)
by selecting the **Spring Web**, **Spring Data JPA**, **PostgreSQL Driver**, and
**Testcontainers** starters.

Alternatively, clone the
[guide repository](https://github.com/testcontainers/tc-guide-testing-spring-boot-rest-api).

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
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
    <dependency>
        <groupId>org.postgresql</groupId>
        <artifactId>postgresql</artifactId>
        <scope>runtime</scope>
    </dependency>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-test</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-junit-jupiter</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.testcontainers</groupId>
        <artifactId>testcontainers-postgresql</artifactId>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>io.rest-assured</groupId>
        <artifactId>rest-assured</artifactId>
        <scope>test</scope>
    </dependency>
</dependencies>
```

Using the Testcontainers BOM (Bill of Materials) is recommended so that you
don't have to repeat the version for every Testcontainers module dependency.

## Create the JPA entity

Create `Customer.java`:

```java
package com.testcontainers.demo;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

@Entity
@Table(name = "customers")
class Customer {

  @Id
  @GeneratedValue(strategy = GenerationType.IDENTITY)
  private Long id;

  @Column(nullable = false)
  private String name;

  @Column(nullable = false, unique = true)
  private String email;

  public Customer() {}

  public Customer(Long id, String name, String email) {
    this.id = id;
    this.name = name;
    this.email = email;
  }

  public Long getId() {
    return id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getEmail() {
    return email;
  }

  public void setEmail(String email) {
    this.email = email;
  }
}
```

## Create the Spring Data JPA repository

```java
package com.testcontainers.demo;

import org.springframework.data.jpa.repository.JpaRepository;

interface CustomerRepository extends JpaRepository<Customer, Long> {}
```

## Add the schema creation script

Create `src/main/resources/schema.sql`:

```sql
create table if not exists customers (
    id bigserial not null,
    name varchar not null,
    email varchar not null,
    primary key (id),
    UNIQUE (email)
);
```

Enable schema initialization in `src/main/resources/application.properties`:

```properties
spring.sql.init.mode=always
```

## Create the REST API endpoint

Create `CustomerController.java`:

```java
package com.testcontainers.demo;

import java.util.List;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
class CustomerController {

  private final CustomerRepository repo;

  CustomerController(CustomerRepository repo) {
    this.repo = repo;
  }

  @GetMapping("/api/customers")
  List<Customer> getAll() {
    return repo.findAll();
  }
}
```
