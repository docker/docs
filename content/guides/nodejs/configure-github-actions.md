---
title: Automate your builds with GitHub Actions
linkTitle: GitHub Actions CI
weight: 40
keywords: CI/CD, GitHub Actions, Node.js, Docker
description: Learn how to configure CI/CD using GitHub Actions for your Node.js application.
aliases:
  - /language/nodejs/configure-ci-cd/
  - /guides/language/nodejs/configure-ci-cd/
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

If you haven't created a [GitHub repository](https://github.com/new) for your project yet, do that now. After creating the repository, [add a remote](https://docs.github.com/en/get-started/getting-started-with-git/managing-remote-repositories) and make sure you can commit and [push your code](https://docs.github.com/en/get-started/using-git/pushing-commits-to-a-remote-repository#about-git-push) to GitHub.

1. In your project's GitHub repository, open **Settings**, and go to **Secrets and variables** > **Actions**.

2. Under the **Variables** tab, create a new **Repository variable** named `DOCKER_USERNAME` with your Docker ID as the value.

3. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token) for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

4. Add the PAT as a **Repository secret** in your GitHub repository, with the name `DOCKERHUB_TOKEN`.

## Overview

GitHub Actions is a CI/CD automation tool built into GitHub. A workflow is a YAML file that tells GitHub which jobs to run when something happens in your repository, like a push to a branch or a pull request opening. Workflows live in the `.github/workflows/` directory of your repository.

In this section, you'll add a workflow that runs your tests on every push to the main branch, then builds your Docker image and pushes it to Docker Hub.

## Define the GitHub Actions workflow

You can create a GitHub Actions workflow by creating a YAML file in the `.github/workflows/` directory of your repository. Use your favorite text editor or the GitHub web interface.

If you prefer to use the GitHub web interface:

1. Go to your repository on GitHub and select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub Actions workflow file in your repository. By default, the file is created under `.github/workflows/main.yml`. Change the filename to `build.yml`.

If you prefer to use your text editor, create a new file named `build.yml` in the `.github/workflows/` directory of your repository.

Add the following content to the file:

{{< files name="nodejs-docker-example" >}}

{{< file path=".github/workflows/build.yml" status="new" >}}
```yaml
# GitHub Actions workflow that runs on every push to main.
# - test: runs Vitest unit tests inside a container.
# - build_and_push: signs in to Docker Hub and the DHI registry, then
#   builds and pushes the image.
name: Build and push Docker image

on:
  push:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@{{% param "checkout_action_version" %}}

      - name: Run tests
        run: docker build --target test -t nodejs-test . && docker run --rm nodejs-test

  build_and_push:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@{{% param "checkout_action_version" %}}

      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Login to Docker Hardened Images
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          registry: dhi.io
          username: ${{ vars.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          push: true
          tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
```
{{< /file >}}

{{< /files >}}

The workflow has two jobs:

1. **test**: Builds the `test` stage of the Dockerfile and runs it. If tests fail, the workflow stops and `build_and_push` doesn't run.
2. **build_and_push**: Signs in to Docker Hub and the DHI registry, then builds and pushes the image.

## Run the workflow

Commit the changes and push them to the `main` branch. This workflow runs every time you push changes to `main`. You can find more information about workflow triggers in the [GitHub documentation](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows).

Go to the **Actions** tab of your GitHub repository. It displays the workflow. Selecting the workflow shows you the breakdown of all the steps.

When the workflow is complete, go to your [repositories on Docker Hub](https://hub.docker.com/repositories). If you see the new repository in that list, the GitHub Actions workflow successfully pushed the image to Docker Hub.

## Summary

In this section, you learned how to set up a GitHub Actions workflow for your Node.js application that includes:

- Running Vitest unit tests inside a container
- Building and pushing Docker images

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [docker/login-action](https://github.com/docker/login-action)
- [docker/build-push-action](https://github.com/docker/build-push-action)
- [Create a Docker Hub access token](/manuals/security/access-tokens.md#create-an-access-token)

## Next steps

In the next section, you'll learn how to inspect and generate supply chain
attestations for your image. See [Secure your supply chain](secure-supply-chain.md).
