---
title: Introduction to Azure Pipelines with Docker
linkTitle: Azure Pipelines and Docker
summary: |
  Learn how to automate Docker image build and push using Azure Pipelines.
params:
  tags: [devops]
  time: 10 minutes
---

## Prerequisites

Before you begin, ensure the following requirements are met:

- A [Docker Hub account](https://hub.docker.com) with a generated access token.
- An active [Azure DevOps project](https://dev.azure.com/) with a connected [Git repository](https://learn.microsoft.com/en-us/azure/devops/repos/git/?view=azure-devops).
- A project that includes a valid [`Dockerfile`](https://docs.docker.com/engine/reference/builder/) at its root or appropriate build context.

## Overview

This guide walks you through building and pushing Docker images using [Azure Pipelines](https://azure.microsoft.com/en-us/products/devops/pipelines), enabling a streamlined and secure CI workflow for containerized applications. You’ll learn how to:

- Configure Docker authentication securely.
- Set up an automated pipeline to build and push images.

## Step 1: Configure Docker credentials

To authenticate securely with Docker Hub:

1. In your Azure DevOps project, navigate to **Project Settings > Pipelines > Library**.
2. Create a new **Variable Group** and add the following variables:
   - `DOCKER_USERNAME` — your Docker Hub username.
   - `DOCKER_PASSWORD` — your Docker Hub access token (mark this as **secret**).

These credentials will be used by the pipeline to log in to Docker Hub.

## Step 2: Create your pipeline

Add the following `azure-pipelines.yml` file to the root of your repository:

```yaml
# Trigger pipeline on commits to the main branch
trigger:
  - main

# Trigger pipeline on pull requests targeting the main branch
pr:
  - main

# Define variables for reuse across the pipeline
variables:
  imageName: '$(dockerRegistry)/$(dockerUsername)/my-image'  # Docker image name with registry and username
  dockerRegistry: 'docker.io'  # Docker registry URL (e.g., Docker Hub or private registry)
  dockerUsername: '$(DOCKER_USERNAME)'  # Docker username from variable group or pipeline variable
  buildTag: '$(Build.BuildId)'  # Tag the image with the unique build ID
  latestTag: 'latest'  # Additional tag for the latest image
  - group: my-variable-group  # Link to Azure DevOps variable group containing DOCKER_USERNAME and DOCKER_PASSWORD

# Define stages for the pipeline
stages:
  - stage: BuildAndPush
    displayName: Build and Push Docker Image
    # Only run this stage for the main branch and if previous steps succeeded
    condition: and(succeeded(), eq(variables['Build.SourceBranch'], 'refs/heads/main'))
    jobs:
      - job: DockerJob
        displayName: Build and Push
        # Use the latest Ubuntu VM image for compatibility with Docker
        pool:
          vmImage: ubuntu-latest
        steps:
          # Check out the repository code
          - checkout: self
            displayName: Checkout Code

          # Log in to Docker registry using a service connection for secure credential management
          - task: Docker@2
            displayName: Docker Login
            inputs:
              command: login
              containerRegistry: 'my-docker-registry'  # Service connection defined in Azure DevOps

          # Build the Docker image with BuildKit for caching and efficiency
          - task: Docker@2
            displayName: Build Docker Image
            inputs:
              command: build
              repository: $(imageName)  # Image name with registry and username
              tags: |
                $(buildTag)  # Tag with build ID
                $(latestTag)  # Tag as latest
              dockerfile: './Dockerfile'  # Path to Dockerfile
              arguments: '--cache-from $(imageName):latest'  # Use cached layers from latest image
            env:
              DOCKER_BUILDKIT: 1  # Enable BuildKit for faster builds and better caching

          # Validate the built Docker image by running a simple command (e.g., version check)
          - script: |
              docker run --rm $(imageName):$(buildTag) --version
            displayName: Validate Docker Image
            continueOnError: true  # Continue even if validation fails (optional)

          # Push the Docker image to the registry
          - task: Docker@2
            displayName: Push Docker Image
            inputs:
              command: push
              repository: $(imageName)
              tags: |
                $(buildTag)  # Push build ID tag
                $(latestTag)  # Push latest tag
```

## What this pipeline does

- Triggers on commits and pull requests targeting the `main` branch.
- Authenticates with Docker Hub (or another specified registry) using a secure Azure DevOps service connection for credential management.
- Builds and tags the Docker image with the Azure build ID and a latest tag, utilizing Docker BuildKit for efficient caching.
- Validates the built Docker image by running a simple command (e.g., version check) to ensure it functions as expected.
- Pushes the tagged Docker images to the specified Docker registry (e.g., Docker Hub).

---

## Summary

With a streamlined configuration, this Azure Pipelines CI workflow:

- Automatically triggers on commits and pull requests to the main branch, building and pushing Docker images.
- Authenticates securely with Docker Hub (or another registry) using an Azure DevOps service connection for credential management.
- Builds and tags Docker images with the Azure build ID and a latest tag, leveraging Docker BuildKit for efficient caching.
- Validates the built image with a simple command (e.g., version check) to ensure functionality.
- Pushes the tagged images to the specified Docker registry (e.g., Docker Hub).

**You can extend this pipeline to support:**

- Multi-platform image builds for broader architecture compatibility.
- Deployment stages for Helm or Kubernetes-based environments.
- Integration with [Azure Container Registry (ACR)](https://learn.microsoft.com/en-us/azure/container-registry/) for private registry storage.

---

## Further Reading

- [Azure Pipelines Documentation](https://learn.microsoft.com/en-us/azure/devops/pipelines/?view=azure-devops) - Comprehensive guide to configuring and managing CI/CD pipelines in Azure DevOps.
- [Docker Task for Azure Pipelines](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/build/docker) - Detailed reference for using the Docker task in Azure Pipelines to build and push images.
- [Docker Buildx Bake](/manuals/build/bake/_index.md) - Explore Docker's advanced build tool for complex, multi-stage, and multi-platform build setups. See also the [Mastering Buildx Bake Guide](/guides/bake/index.md) for practical examples and best practices.
- [Docker Build Cloud](/guides/docker-build-cloud/_index.md) - Learn about Docker's managed build service for faster, scalable, and multi-platform image builds in the cloud.