---
title: Configure CI/CD for your Node.js application
keywords: ci/cd, github actions, node.js, node
description: Learn how to configure CI/CD using GitHub Actions for your Node.js application.
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a [Docker](https://hub.docker.com/signup) account to complete this section.

## Overview

In this section, you'll learn how to set up and use GitHub Actions to build and test your Docker image as well as push it to Docker Hub. You will complete the following steps:

1. Create a new repository on GitHub.
2. Define the GitHub Actions workflow.
3. Run the workflow.

## Step one: Create the repository

Create a GitHub repository, configure the Docker Hub secrets, and push your source code.

1. [Create a new repository](https://github.com/new) on GitHub.

2. Open the repository **Settings**, and go to **Secrets and variables** >
   **Actions**.

3. Create a new secret named `DOCKER_USERNAME` and your Docker ID as value.

4. Create a new [Personal Access Token
   (PAT)](/security/for-developers/access-tokens/#create-an-access-token) for Docker Hub. You
   can name this token `node-docker`.

5. Add the PAT as a second secret in your GitHub repository, with the name
   `DOCKERHUB_TOKEN`.

6. In your local repository on your machine, run the following command to change
   the origin to the repository you just created. Make sure you change
   `your-username` to your GitHub username and `your-repository` to the name of
   the repository you created.

   ```console
   $ git remote set-url origin https://github.com/your-username/your-repository.git
   ```

7. Run the following command to push your local repository to GitHub.

   ```console
   $ git push -u origin main
   ```

## Step two: Set up the workflow

Set up your GitHub Actions workflow for building, testing, and pushing the image
to Docker Hub.

1. Go to your repository on GitHub and then select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub actions workflow file in
   your repository, under `.github/workflows/main.yml` by default.

3. In the editor window, copy and paste the following YAML configuration.

   ```yaml
   name: ci
   
   on:
     push:
       branches:
         - main
   
   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         -
           name: Checkout
           uses: actions/checkout@v4
         -
           name: Login to Docker Hub
           uses: docker/login-action@v3
           with:
             username: ${{ secrets.DOCKER_USERNAME }}
             password: ${{ secrets.DOCKERHUB_TOKEN }}
         -
           name: Set up Docker Buildx
           uses: docker/setup-buildx-action@v3
         -
           name: Build and test
           uses: docker/build-push-action@v5
           with:
             context: .
             target: test
             load: true
         -
           name: Build and push
           uses: docker/build-push-action@v5
           with:
             context: .
             push: true
             target: prod
             tags: ${{ secrets.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
   ```

   For more information about the YAML syntax used here, see [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions).

## Step three: Run the workflow

Save the workflow file and run the job.

1. Select **Commit changes...** and push the changes to the `main` branch.

   After pushing the commit, the workflow starts automatically.

2. Go to the **Actions** tab. It displays the workflow.

   Selecting the workflow shows you the breakdown of all the steps.

3. When the workflow is complete, go to your
   [repositories on Docker Hub](https://hub.docker.com/repositories).

   If you see the new repository in that list, it means the GitHub Actions
   successfully pushed the image to Docker Hub.

## Summary

In this section, you learned how to set up a GitHub Actions workflow for your Node.js application.

Related information:
 - [Introduction to GitHub Actions](../../build/ci/github-actions/index.md)
 - [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

## Next steps

Next, learn how you can locally test and debug your workloads on Kubernetes before deploying.

{{< button text="Develop using Kubernetes" url="./deploy.md" >}}