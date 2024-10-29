---
title: Configuring Testcontainers Cloud in the CI Pipeline
description: Use Testcontainers Cloud with GitHub Workflows to automate testing in a CI pipeline.
weight: 30
---

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
