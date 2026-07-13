---
title: |
  Mastering Testcontainers Cloud by Docker: streamlining integration testing with containers
linkTitle: Testcontainers Cloud by Docker
summary: |
  Automate, scale, and optimize testing workflows with Testcontainers Cloud
description: |
  Testcontainers Cloud by Docker streamlines integration testing by offloading container management to the cloud. It enables faster, consistent tests for containerized services like databases, improving performance and scalability in CI/CD pipelines without straining local or CI resources. Ideal for developers needing efficient, reliable testing environments.
keywords: testcontainers cloud, integration testing, ci/cd, containerized tests, cloud testing, scalable testing
aliases:
  - /guides/testcontainers-cloud/common-questions/
  - /guides/testcontainers-cloud/demo-ci/
  - /guides/testcontainers-cloud/demo-local/
  - /guides/testcontainers-cloud/why/
params:
  tags: [testing]
  image: images/learning-paths/testcontainers-cloud-learning-path.png
  time: 12 minutes
---


Testcontainers Cloud is a cloud-based solution designed to streamline and enhance the process of running integration tests using Testcontainers. Testcontainers is the open source framework, which allows developers to easily spin up containerized dependencies such as databases, message brokers, and other services required for testing. By shifting the management of Testcontainers-based services to the cloud, Testcontainers Cloud optimizes performance, reduces resource constraints on local machines or CI servers, and ensures consistent test environments. This solution is particularly beneficial for teams working on complex, distributed systems, as it allows for scalable, isolated, and reliable testing without the typical overhead of managing containers locally.

## What you'll learn

- Understand the fundamentals of Docker Testcontainers Cloud and its role in integration testing.
- Learn how to set up and configure Docker Testcontainers Cloud for automated testing in various environments.
- Explore how Testcontainers Cloud integrates with CI/CD pipelines to streamline testing workflows.

## Tools integration

Works well with Docker Desktop, GitHub Actions, Jenkins, Kubernetes, and other CI solutions

Docker Pro, Team, and Business subscriptions come with Testcontainers Cloud runtime minutes, and additional minutes are available via consumption pricing. Testcontainers Cloud runtime minutes do not rollover month to month.

## Who’s this for?

- Teams that build cloud-native applications and are already using Testcontainers open source.
- DevOps Teams that integrate automated container-based testing into CI/CD pipelines for continuous testing.
- QA Teams that seek scalable and consistent test environments for comprehensive integration and end-to-end testing.
- Developers who need reliable, containerized test environments for testing microservices and databases.

## Why Testcontainers Cloud?

{{< youtube-embed "6dRRlk5Vd0E" >}}

Testcontainers Cloud is a powerful cloud-based solution designed to optimize integration testing with Testcontainers by offloading container management to the cloud. It helps developers and teams overcome the limitations of traditional local and CI-based testing, ensuring consistent environments, faster test execution, and scalable workflows. Whether you're new to Testcontainers or looking to enhance your existing setup, Testcontainers Cloud offers a seamless way to manage containerized tests, improving efficiency and reliability in your development pipeline.

Testcontainers Cloud provides several benefits:

- **Offloading to the Cloud:** Frees up local resources by shifting container management to the cloud, keeping your laptop responsive.
- **Consistent Testing Environments:** Ensures that tests run in isolated, reliable environments, reducing inconsistencies across platforms from Dev to CI.
- **Scalability:** Allows running large numbers of containers simultaneously without being limited by local or CI resources.
- **Faster CI/CD Pipelines:** Reduces configuration bottlenecks and speeds up build times by offloading containers to multiple on-demand cloud workers with the Turbo-mode feature.

Testcontainers Cloud streamlines integration testing by offloading container management to the cloud, ensuring consistent environments and faster test execution resulting in reduced resource strain, making it an essential tool for improving the stability of your Testcontainers-based workflows.

## Setting up Testcontainers Cloud by Docker

{{< youtube-embed "7c3xLAG560U" >}}

This demo shows the process of setting up Testcontainers Cloud by Docker to
work in your local development environment using the Testcontainers Desktop
application. By the end of this walkthrough, you'll have Testcontainers Cloud
by Docker up and running, ready to offload container management from your local
machine to the cloud for more efficient testing.

- Install and configure Testcontainers Cloud and the CLI to seamlessly integrate with your local development environment.
- Set up and configure the Testcontainers Desktop application to monitor and manage cloud-based containers during local tests.
- Create and run integration tests using Testcontainers that leverage cloud-based container resources.
- Monitor and manage containers efficiently, understanding how Testcontainers Cloud automates cleanup and ensures consistent testing environments.
- Review options for monitoring and troubleshooting in the Testcontainers Cloud Dashboard.

## Configuring Testcontainers Cloud in the CI Pipeline

{{< youtube-embed "NlZY9aumKJU" >}}

This demo shows how Testcontainers Cloud can be seamlessly integrated into a
Continuous Integration (CI) pipeline using GitHub Workflows, providing a
powerful solution for running containerized integration tests without
overloading local or CI runner resources. By leveraging GitHub Actions,
developers can automate the process of spinning up and managing containers for
testing in the cloud, ensuring faster and more reliable test execution. With
just a few configuration steps, including setting up Testcontainers Cloud
authentication and adding it to your workflow, you can offload container
orchestration to the cloud. This approach improves the scalability of your
pipeline, ensures consistency across tests, and simplifies resource management,
making it an ideal solution for modern, containerized development workflows.

- Understand how to set up a GitHub Actions workflow to automate the build and testing of a project.
- Learn how to configure Testcontainers Cloud within GitHub Actions to offload containerized testing to the cloud, improving efficiency and resource management.
- Explore how Testcontainers Cloud integrates with GitHub workflows to run integration tests that require containerized services, such as databases and message brokers.

## Common challenges and questions

<!-- vale Docker.HeadingLength = NO -->

#### How is Testcontainers Cloud different from the open-source Testcontainers framework?

While the open-source Testcontainers is a library that provides a lightweight APIs for bootstrapping local development and test dependencies with real services wrapped in Docker containers, Testcontainers Cloud provides a cloud runtime for these containers. This reduces the resource strain on local environments and provides more scalability, especially in CI/CD workflows, that enables consistent Testcontainers experience across the organization.

#### What types of containers can I run with Testcontainers Cloud?

Testcontainers Cloud supports any containers you would typically use with the Testcontainers framework, including databases (PostgreSQL, MySQL, MongoDB), message brokers (Kafka, RabbitMQ), and other services required for integration testing.

#### Do I need to change my existing test code to use Testcontainers Cloud?

No, you don't need to change your existing test code. Testcontainers Cloud integrates seamlessly with the open-source Testcontainers framework. Once the cloud configuration is set up, it automatically manages the containers in the cloud without requiring code changes.

#### How do I integrate Testcontainers Cloud into my project?

To integrate Testcontainers Cloud, you need to install the Testcontainers Desktop app and select run with Testcontainers Cloud option in the menu. In CI you’ll need to add a workflow step that downloads Testcontainers Cloud agent. No code changes are required beyond enabling Cloud runtime via the Testcontainers Desktop app locally or installing Testcontainers Cloud agent in CI.

#### Can I use Testcontainers Cloud in a CI/CD pipeline?

Yes, Testcontainers Cloud is designed to work efficiently in CI/CD pipelines. It helps reduce build times and resource bottlenecks by offloading containers that you spin up with Testcontainers library to the cloud, making it a perfect fit for continuous testing environments.

#### What are the benefits of using Testcontainers Cloud?

The key benefits include reduced resource usage on local machines and CI servers, scalability (run more containers without performance degradation), consistent testing environments, centralized monitoring, ease of CI configuration with removed security concerns of running Docker-in-Docker or a privileged daemon.

#### Does Testcontainers Cloud support all programming languages?

Testcontainers Cloud supports any language that works with the open-source Testcontainers libraries, including Java, Python, Node.js, Go, and others. As long as your project uses Testcontainers, it can be offloaded to Testcontainers Cloud.

#### How is container cleanup handled in Testcontainers Cloud?

While Testcontainers library automatically handles container lifecycle management, Testcontainers Cloud manages the allocated cloud worker lifetime. This means that containers are spun up, monitored, and cleaned up after tests are completed by Testcontainers library, and the worker where these containers have being running will be removed automatically after the ~35 min idle period by Testcontainers Cloud. This approach frees developers from manually managing containers and associated cloud resources.

#### Is there a free tier or pricing model for Testcontainers Cloud?

Pricing details for Testcontainers Cloud can be found on the [pricing page](https://testcontainers.com/cloud/pricing/).
