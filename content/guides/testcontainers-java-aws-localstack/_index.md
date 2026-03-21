---
title: Testing AWS service integrations using LocalStack
linkTitle: AWS LocalStack
description: Learn how to test Spring Cloud AWS applications using LocalStack and Testcontainers.
keywords: testcontainers, java, spring boot, testing, aws, localstack, s3, sqs, spring cloud aws
summary: |
  Learn how to create a Spring Boot application with Spring Cloud AWS,
  then test S3 and SQS integrations using Testcontainers and LocalStack.
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [java]
params:
  time: 25 minutes
---

<!-- Source: https://github.com/testcontainers/tc-guide-testing-aws-service-integrations-using-localstack -->

In this guide, you will learn how to:

- Create a Spring Boot application with Spring Cloud AWS integration
- Use AWS S3 and SQS services
- Test the application using Testcontainers and LocalStack

## Prerequisites

- Java 17+
- Maven or Gradle
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.
