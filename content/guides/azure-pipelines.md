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

## Step 1: Configure a Docker Hub service connection

To securely authenticate with Docker Hub using Azure Pipelines:

1. Navigate to **Project Settings > Service Connections** in your Azure DevOps project.
2. Click **New service connection > Docker Registry**.
3. Choose **Docker Hub** and provide your Docker Hub credentials or access token.
4. Give the service connection a recognizable name, such as `my-docker-registry`.
5. Grant access only to the specific pipeline(s) that require it for improved security and least privilege.

> [!IMPORTANT]
> Avoid selecting the option to grant access to all pipelines unless absolutely necessary. Always apply the principle of least privilege.

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
  imageName: 'docker.io/$(dockerUsername)/my-image'
  dockerUsername: 'your-dockerhub-username'  # Replace with your Docker Hub username
  buildTag: '$(Build.BuildId)'
  latestTag: 'latest'

stages:
  - stage: BuildAndPush
    displayName: Build and Push Docker Image
    condition: and(succeeded(), eq(variables['Build.SourceBranch'], 'refs/heads/main'))
    jobs:
      - job: DockerJob
        displayName: Build and Push
        pool:
          vmImage: ubuntu-latest
        steps:
          - checkout: self
            displayName: Checkout Code

          - task: Docker@2
            displayName: Docker Login
            inputs:
              command: login
              containerRegistry: 'my-docker-registry'  # Service connection name

          - task: Docker@2
            displayName: Build Docker Image
            inputs:
              command: build
              repository: $(imageName)
              tags: |
                $(buildTag)
                $(latestTag)
              dockerfile: './Dockerfile'
              arguments: '--cache-from $(imageName):latest'
            env:
              DOCKER_BUILDKIT: 1

          - task: Docker@2
            displayName: Push Docker Image
            inputs:
              command: push
              repository: $(imageName)
              tags: |
                $(buildTag)
                $(latestTag)

          # Optional: logout for self-hosted agents
          - script: docker logout
            displayName: Docker Logout (Self-hosted only)
            condition: ne(variables['Agent.OS'], 'Windows_NT')
```

## What this pipeline does

This pipeline automates the Docker image build and deployment process for the main branch. It ensures a secure and efficient workflow with best practices like caching, tagging, and conditional cleanup. Here's what it does:

- Triggers on commits and pull requests targeting the `main` branch.
- Authenticates securely with Docker Hub using an Azure DevOps service connection.
- Builds and tags the Docker image using Docker BuildKit for caching.
- Pushes both buildId and latest tags to Docker Hub.
- Logs out from Docker if running on a self-hosted Linux agent.


## Detailed Step-by-Step Explanation 

### Step 1: Define Pipeline Triggers 

```yaml
trigger:
  - main

pr:
  - main
```

This pipeline is triggered automatically on:
- Commits pushed to the `main` branch
- Pull requests targeting `main` main branch

> [!NOTE]
> Learn more: [Define pipeline triggers in Azure Pipelines](https://learn.microsoft.com/en-us/azure/devops/pipelines/build/triggers?view=azure-devops)


### Step 2: Define Common Variables

```yaml
variables:
  imageName: 'docker.io/$(dockerUsername)/my-image'
  dockerUsername: 'your-dockerhub-username'  # Replace with your actual Docker Hub username
  buildTag: '$(Build.BuildId)'
  latestTag: 'latest'
```

These variables ensure consistent naming, versioning, and reuse throughout the pipeline steps:

- `imageName`: your image path on Docker Hub
- `buildTag`: a unique tag for each pipeline run
- `latestTag`: a stable alias for your most recent image

> [!NOTE]
> Learn more: [Define and use variables in Azure Pipelines](https://learn.microsoft.com/en-us/azure/devops/pipelines/process/variables?view=azure-devops&tabs=yaml%2Cbatch)


### Step 3: Define Pipeline Stages and Jobs

```yaml
stages:
  - stage: BuildAndPush
    displayName: Build and Push Docker Image
    condition: and(succeeded(), eq(variables['Build.SourceBranch'], 'refs/heads/main'))
```

This stage executes only if:

- The pipeline completes successfully.
- The source branch is main.

> [!NOTE]
> Learn more: [Stage conditions in Azure Pipelines](https://learn.microsoft.com/en-us/azure/devops/pipelines/process/stages?view=azure-devops&tabs=yaml)

### Step 4: Job Configuration

```yaml
jobs:
  - job: DockerJob
  displayName: Build and Push
  pool:
  vmImage: ubuntu-latest
```

This job uses the latest Ubuntu VM image provided by Microsoft-hosted agents. It can be swapped with a custom pool for self-hosted agents if needed.

> [!NOTE]
> Learn more: [Specify a pool in your pipeline](https://learn.microsoft.com/en-us/azure/devops/pipelines/agents/pools-queues?view=azure-devops&tabs=yaml%2Cbrowser)

#### Step 4.1 Checkout Code

```yaml
steps:
  - checkout: self
    displayName: Checkout Code

```

This step pulls your repository code into the build agent, so the pipeline can access the Dockerfile and application files.

> [!NOTE]
> Learn more: [checkout step documentation](https://learn.microsoft.com/en-us/azure/devops/pipelines/yaml-schema/steps-checkout?view=azure-pipelines)


#### Step 4.2 Authenticate to Docker Hub

```yaml
- task: Docker@2
  displayName: Docker Login
  inputs:
    command: login
    containerRegistry: 'my-docker-registry'  # Replace with your service connection name
```

Uses a preconfigured Azure DevOps Docker registry service connection to authenticate securely without exposing credentials directly.

> [!NOTE]
> Learn more: [Use service connections for Docker Hub](https://learn.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops#docker-hub-or-others)


#### Step 4.3 Build the Docker Image

```yaml
- task: Docker@2
  displayName: Build Docker Image
  inputs:
    command: build
    repository: $(imageName)
    tags: |
        $(buildTag)
        $(latestTag)
    dockerfile: './Dockerfile'
    arguments: '--cache-from $(imageName):latest'
  env:
    DOCKER_BUILDKIT: 1
```

This builds the image with:

- Two tags: one with the build ID and one as latest
- Docker BuildKit for faster builds and layer caching
- Cache pull from the last pushed latest tag

> [!NOTE]
> Learn more: [Docker task for Azure Pipelines](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/reference/docker-v2?view=azure-pipelines&tabs=yaml)


#### Step 4.4 Push the Docker Image
```yaml
- task: Docker@2
  displayName: Push Docker Image
  inputs:
      command: push
      repository: $(imageName)
      tags: |
        $(buildTag)
        $(latestTag)

```

This uploads both tags to Docker Hub:
- `$(buildTag)` ensures traceability per run.
- `latest` is used for most recent image references.


5. Logout from Docker (Self-Hosted Agents)

```yaml
- script: docker logout
  displayName: Docker Logout (Self-hosted only)
  condition: ne(variables['Agent.OS'], 'Windows_NT')

```

Executes docker logout at the end of the pipeline on Linux-based self-hosted agents to proactively clean up credentials and enhance security posture.

---

## Summary

With this Azure Pipelines CI setup, you get:
- Secure Docker authentication using a built-in service connection.
- Automated image building and tagging triggered by code changes.
- Efficient builds leveraging Docker BuildKit cache.
- Safe cleanup with logout on persistent agents.

---

## Further Reading

- [Azure Pipelines Documentation](https://learn.microsoft.com/en-us/azure/devops/pipelines/?view=azure-devops) - Comprehensive guide to configuring and managing CI/CD pipelines in Azure DevOps.
- [Docker Task for Azure Pipelines](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/build/docker) - Detailed reference for using the Docker task in Azure Pipelines to build and push images.
- [Docker Buildx Bake](/manuals/build/bake/_index.md) - Explore Docker's advanced build tool for complex, multi-stage, and multi-platform build setups. See also the [Mastering Buildx Bake Guide](/guides/bake/index.md) for practical examples and best practices.
- [Docker Build Cloud](/guides/docker-build-cloud/_index.md) - Learn about Docker's managed build service for faster, scalable, and multi-platform image builds in the cloud.