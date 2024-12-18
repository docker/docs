---
title: "Faster development and testing with container-supported development"
linkTitle: Container-supported development
summary: |
  Containers don't have to be just for your app. Learn how to run your app's dependent services and other debugging tools to enhance your development environment.
description: |
  Use containers in your local development loop to develop and test faster… even if your main app isn't running in containers.
tags: [app-dev]
params:
  image: images/learning-paths/container-supported-development.png
  time: 20 minutes
  resource_links: []
---

Containers offer a consistent way to build, share, and run applications across different environments. While containers are typically used to containerize your application, they also make it incredibly easy to run essential services needed for development. Instead of installing or connecting to a remote database, you can easily launch your own database. But the possibilities don't stop there.

With container-supported development, you use containers to enhance your development environment by emulating or running your own instances of the services your app needs. This provides faster feedback loops, less coupling with remote services, and a greater ability to test error states.

And best of all, you can have these benefits regardless of whether the main app under development is running in containers.

## What you'll learn

- The meaning of container-supported development
- How to connect non-containerized applications to containerized services
- Several examples of using containers to emulate or run local instances of services
- How to use containers to add additional troubleshooting and debugging tools to your development environment

## Who's this for?

- Teams that want to reduce the coupling they have on shared or deployed infrastructure or remote API endpoints
- Teams that want to reduce the complexity and costs associated with using cloud services directly during development
- Developers that want to make it easier to visualize what's going on in their databases, queues, etc.
- Teams that want to reduce the complexity of setting up their development environment without impacting the development of the app itself


## Tools integration

Works well with Docker Compose and Testcontainers.

## Modules

### What is container-supported development?

Container-supported development is the idea of using containers to enhance your development environment by running local instances or emulators of the services your application relies on. Once you're using containers, it's easy to add additional services to visualize or troubleshoot what's going on in your services.

{{< youtube-embed pNcrto_wGi0 >}}

### Demo: running databases locally

With container-supported development, it's easy to run databases locally. In this demo, you'll see how to do so, as well as how to connect a non-containerized application to the database.

{{< youtube-embed VieWeXOwKLU >}}

> [!TIP]
>
> Learn more about running databases in containers in the [Use containerized databases](/guides/databases.md) guide.

### Demo: mocking API endpoints

Many APIs require data from other data endpoints. In development, this adds complexities such as the sharing of credentials, uptime/availability, and rate limiting. Instead of relying on those services directly, your application can interact with a mock API server.

This demo will demonstrate how using WireMock can make it easy to develop and test an application, including the APIs various error states.

{{< youtube-embed VXSmX6f8vo0 >}}

> [!TIP]
>
> Learn more about using WireMock to mock API in the [Mocking API services with WireMock](/guides/wiremock.md) guide.

### Demo: developing the cloud locally

When developing apps, it's often easier to outsource aspects of the application to cloud services, such as Amazon S3. However, connecting to those services in local development introduces IAM policies, networking constraints, and provisioning complications. While these requirements are important in a production setting, they complicate development environments significantly. 

With container-supported development, you can run local instances of these services during development and testing, removing the need for complex setups. In this demo, you'll see how LocalStack makes it easy to develop and test applications entirely from the developer's workstation.

{{< youtube-embed JtwUMvR5xlY >}}

> [!TIP]
>
> Learn more about using LocalStack in the [Develop and test AWS Cloud applications using LocalStack](/guides/localstack.md) guide.

### Demo: adding additional debug and troubleshooting tools

Once you start using containers in your development environment, it becomes much easier to add additional containers to visualize the contents of the databases or message queues, seed document stores, or event publishers. In this demo, you'll see a few of these examples, as well as how you can connect multiple containers together to make testing even easier.

{{< youtube-embed TCZX15aKSu4 >}}

<div id="lp-survey-anchor"></div>
