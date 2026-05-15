---
title: Replace H2 with a real database for testing
linkTitle: Replace H2 database
description: Learn how to replace an H2 in-memory database with a real PostgreSQL database for testing using Testcontainers.
keywords: testcontainers, java, testing, h2, postgresql, spring boot, spring data jpa, jdbc
summary: |
  Replace your H2 in-memory test database with a real PostgreSQL instance
  using the Testcontainers special JDBC URL — a one-line change.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 15 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-replace-h2-with-real-database-for-testing -->

In this guide, you will learn how to:

- Understand the drawbacks of using H2 in-memory databases for testing
- Replace H2 with a real PostgreSQL database using the Testcontainers special JDBC URL
- Use the Testcontainers JUnit 5 extension for more control over the container
- Test both Spring Data JPA and JdbcTemplate-based repositories

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
