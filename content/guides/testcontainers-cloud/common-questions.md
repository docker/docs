---
title: Common challenges and questions
description: Explore common challenges and questions related to Testcontainers Cloud by Docker.
weight: 40
---

<!-- vale Docker.HeadingLength = NO -->

### How is Testcontainers Cloud different from the open-source Testcontainers framework?

While the open-source Testcontainers framework runs containers locally on your development machine or CI server, Testcontainers Cloud moves this process to the cloud. This reduces the resource strain on local environments and provides more scalability, especially in CI/CD workflows.

### What types of containers can I run with Testcontainers Cloud?

Testcontainers Cloud supports any containers you would typically use with the Testcontainers framework, including databases (PostgreSQL, MySQL, MongoDB), message brokers (Kafka, RabbitMQ), and other services required for integration testing.

### Do I need to change my existing test code to use Testcontainers Cloud?

No, you don't need to change your existing test code. Testcontainers Cloud integrates seamlessly with the open-source Testcontainers framework. Once the cloud configuration is set up, it automatically manages the containers in the cloud without requiring code changes.

### How do I integrate Testcontainers Cloud into my project?

To integrate Testcontainers Cloud, you need to install the Testcontainers Cloud CLI, configure your project for cloud integration, and ensure your containers are running through the cloud service. No major code changes are required beyond setting up the CLI and enabling cloud integration.

### Can I use Testcontainers Cloud in a CI/CD pipeline?

Yes, Testcontainers Cloud is designed to work efficiently in CI/CD pipelines. It helps reduce build times and resource bottlenecks by offloading containerized tests to the cloud, making it a perfect fit for continuous testing environments.

### What are the benefits of using Testcontainers Cloud?

The key benefits include faster test execution, reduced resource usage on local machines and CI servers, scalability (run more containers without performance degradation), consistent testing environments, and simplified container management.

### Does Testcontainers Cloud support all programming languages?

Testcontainers Cloud supports any language that works with the open-source Testcontainers framework, including Java, Python, Node.js, and Go. As long as your project uses Testcontainers, it can be offloaded to Testcontainers Cloud.

### How is container cleanup handled in Testcontainers Cloud?

Testcontainers Cloud automatically handles container lifecycle management. This means that containers are spun up, monitored, and cleaned up after tests are completed, freeing developers from manually managing containers.

### Is there a free tier or pricing model for Testcontainers Cloud?

Yes, Testcontainers Cloud typically offers a free tier for small-scale projects and testing environments. For larger teams and heavier usage, there are various pricing models based on container usage, resources, and features. You can find detailed pricing information on the Testcontainers Cloud website.
