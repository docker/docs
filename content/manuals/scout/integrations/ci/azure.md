---
description: How to integrate Docker Scout with Microsoft Azure DevOps Pipelines
keywords: supply chain, security, ci, continuous integration, azure, devops
title: Integrate Docker Scout with Microsoft Azure DevOps Pipelines
linkTitle: Azure DevOps Pipelines
---

The following examples runs in an Azure DevOps-connected repository containing
a Docker image's definition and contents. Triggered by a commit to the main
branch, the pipeline builds the image and uses Docker Scout to create a CVE
report.

First, set up the rest of the workflow and set up the variables available to all
pipeline steps. Add the following to an _azure-pipelines.yml_ file:

```yaml
trigger:
  - main

resources:
  - repo: self

variables:
  tag: "$(Build.BuildId)"
  image: "vonwig/nodejs-service"
```

This sets up the workflow to use a particular container image for the
application and tag each new image build with the build ID.

Add the following to the YAML file:

```yaml
stages:
  - stage: Build
    displayName: Build image
    jobs:
      - job: Build
        displayName: Build
        pool:
          vmImage: ubuntu-latest
        steps:
          - task: Docker@2
            displayName: Build an image
            inputs:
              command: build
              dockerfile: "$(Build.SourcesDirectory)/Dockerfile"
              repository: $(image)
              tags: |
                $(tag)
          - task: CmdLine@2
            displayName: Find CVEs on image
            inputs:
              script: |
                # Install the Docker Scout CLI
                curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s --
                # Login to Docker Hub required for Docker Scout CLI
                echo $(DOCKER_HUB_PAT) | docker login -u $(DOCKER_HUB_USER) --password-stdin
                # Get a CVE report for the built image and fail the pipeline when critical or high CVEs are detected
                docker scout cves $(image):$(tag) --exit-code --only-severity critical,high
```

This creates the flow mentioned previously. It builds and tags the image using
the checked-out Dockerfile, downloads the Docker Scout CLI, and then runs the
`cves` command against the new tag to generate a CVE report. It only shows
critical or high-severity vulnerabilities.
