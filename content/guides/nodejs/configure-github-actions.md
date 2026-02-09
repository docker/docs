---
title: Automate your builds with GitHub Actions
linkTitle: Automate your builds with GitHub Actions
weight: 50
keywords: CI/CD, GitHub Actions, Node.js, Docker
description: Learn how to configure CI/CD using GitHub Actions for your Node.js application.
aliases:
  - /language/nodejs/configure-ci-cd/
  - /guides/language/nodejs/configure-ci-cd/
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md).

You must also have:

- A [GitHub](https://github.com/signup) account.
- A verified [Docker Hub](https://hub.docker.com/signup) account.

---

## Overview

In this section, you'll set up a **CI/CD pipeline** using [GitHub Actions](https://docs.github.com/en/actions) to automatically:

- Build your Node.js application inside a Docker container.
- Run unit and integration tests, and make sure your application meets solid code quality standards.
- Perform security scanning and vulnerability assessment.
- Push production-ready images to [Docker Hub](https://hub.docker.com).

---

## Connect your GitHub repository to Docker Hub

To enable GitHub Actions to build and push Docker images, you'll securely store your Docker Hub credentials in your new GitHub repository.

### Step 1: Connect your GitHub repository to Docker Hub

1. Create a Personal Access Token (PAT) from [Docker Hub](https://hub.docker.com).
   1. From your Docker Hub account, go to **Account Settings → Security**.
   2. Generate a new Access Token with **Read/Write** permissions.
   3. Name it something like `docker-nodejs-sample`.
   4. Copy and save the token — you'll need it in Step 4.

2. Create a repository in [Docker Hub](https://hub.docker.com/repositories/).
   1. From your Docker Hub account, select **Create a repository**.
   2. For the Repository Name, use something descriptive — for example: `nodejs-sample`.
   3. Once created, copy and save the repository name — you'll need it in Step 4.

3. Create a new [GitHub repository](https://github.com/new) for your Node.js project.

4. Add Docker Hub credentials as GitHub repository secrets.

   In your newly created GitHub repository:
   1. From **Settings**, go to **Secrets and variables → Actions → New repository secret**.

   2. Add the following secrets:

   | Name                     | Value                                            |
   | ------------------------ | ------------------------------------------------ |
   | `DOCKER_USERNAME`        | Your Docker Hub username                         |
   | `DOCKERHUB_TOKEN`        | Your Docker Hub access token (created in Step 1) |
   | `DOCKERHUB_PROJECT_NAME` | Your Docker Project Name (created in Step 2)     |

   These secrets let GitHub Actions to authenticate securely with Docker Hub during automated workflows.

5. Connect your local project to GitHub.

   Link your local project `docker-nodejs-sample` to the GitHub repository you just created by running the following command from your project root:

   ```console
      $ git remote set-url origin https://github.com/{your-username}/{your-repository-name}.git
   ```

   > [!IMPORTANT]
   > Replace `{your-username}` and `{your-repository}` with your actual GitHub username and repository name.

   To confirm that your local project is correctly connected to the remote GitHub repository, run:

   ```console
   $ git remote -v
   ```

   You should see output similar to:

   ```console
   origin  https://github.com/{your-username}/{your-repository-name}.git (fetch)
   origin  https://github.com/{your-username}/{your-repository-name}.git (push)
   ```

   This confirms that your local repository is properly linked and ready to push your source code to GitHub.

6. Push your source code to GitHub.

   Follow these steps to commit and push your local project to your GitHub repository:
   1. Stage all files for commit.

      ```console
      $ git add -A
      ```

      This command stages all changes — including new, modified, and deleted files — preparing them for commit.

   2. Commit your changes.

      ```console
      $ git commit -m "Initial commit with CI/CD pipeline"
      ```

      This command creates a commit that snapshots the staged changes with a descriptive message.

   3. Push the code to the `main` branch.

      ```console
      $ git push -u origin main
      ```

      This command pushes your local commits to the `main` branch of the remote GitHub repository and sets the upstream branch.

Once completed, your code will be available on GitHub, and any GitHub Actions workflow you've configured will run automatically.

> [!NOTE]  
> Learn more about the Git commands used in this step:
>
> - [Git add](https://git-scm.com/docs/git-add) – Stage changes (new, modified, deleted) for commit
> - [Git commit](https://git-scm.com/docs/git-commit) – Save a snapshot of your staged changes
> - [Git push](https://git-scm.com/docs/git-push) – Upload local commits to your GitHub repository
> - [Git remote](https://git-scm.com/docs/git-remote) – View and manage remote repository URLs

---

### Step 2: Set up the workflow

Now you'll create a GitHub Actions workflow that builds your Docker image, runs tests, and pushes the image to Docker Hub.

1. From your repository on GitHub, select the **Actions** tab in the top menu.

2. When prompted, select **Set up a workflow yourself**.

   This opens an inline editor to create a new workflow file. By default, it will be saved to:
   `.github/workflows/main.yml`

3. Add the following workflow configuration to the new file:

```yaml
name: CI/CD – Node.js Application with Docker

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]

jobs:
  test:
    name: Run Node.js Tests
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:18-alpine
        env:
          POSTGRES_DB: todoapp_test
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Cache npm dependencies
        uses: actions/cache@v4
        with:
          path: ~/.npm
          key: ${{ runner.os }}-npm-${{ hashFiles('**/package-lock.json') }}
          restore-keys: ${{ runner.os }}-npm-

      - name: Build test image
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          context: .
          target: test
          tags: nodejs-app-test:latest
          platforms: linux/amd64
          load: true
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max

      - name: Run tests inside container
        run: |
          docker run --rm \
            --network host \
            -e NODE_ENV=test \
            -e POSTGRES_HOST=localhost \
            -e POSTGRES_PORT=5432 \
            -e POSTGRES_DB=todoapp_test \
            -e POSTGRES_USER=postgres \
            -e POSTGRES_PASSWORD=postgres \
            nodejs-app-test:latest
        env:
          CI: true
        timeout-minutes: 10

  build-and-push:
    name: Build and Push Docker Image
    runs-on: ubuntu-latest
    needs: test

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Cache Docker layers
        uses: actions/cache@v4
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: ${{ runner.os }}-buildx-

      - name: Extract metadata
        id: meta
        run: |
          echo "REPO_NAME=${GITHUB_REPOSITORY##*/}" >> "$GITHUB_OUTPUT"
          echo "SHORT_SHA=${GITHUB_SHA::7}" >> "$GITHUB_OUTPUT"

      - name: Log in to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build and push multi-arch production image
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          context: .
          target: production
          push: true
          platforms: linux/amd64,linux/arm64
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:latest
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:${{ steps.meta.outputs.SHORT_SHA }}
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max
```

This workflow performs the following tasks for your Node.js application:

- Triggers on every `push` or `pull request` to the `main` branch.
- Builds a test Docker image using the `test` stage.
- Runs tests in a containerized environment.
- Stops the workflow if any test fails.
- Caches Docker build layers and npm dependencies for faster runs.
- Authenticates with Docker Hub using GitHub secrets.
- Builds an image using the `production` stage.
- Tags and pushes the image to Docker Hub with `latest` and short SHA tags.

> [!NOTE]
> For more information about `docker/build-push-action`, refer to the [GitHub Action README](https://github.com/docker/build-push-action/blob/master/README.md).

---

### Step 3: Run the workflow

After adding your workflow file, trigger the CI/CD process.

1. Commit and push your workflow file

   From the GitHub editor, select **Commit changes…**.
   - This push automatically triggers the GitHub Actions pipeline.

2. Monitor the workflow execution
   1. From your GitHub repository, go to the **Actions** tab.
   2. Select the workflow run to follow each step: **test**, **build**, **security**, and (if successful) **push** and **deploy**.

3. Verify the Docker image on Docker Hub
   - After a successful workflow run, visit your [Docker Hub repositories](https://hub.docker.com/repositories).
   - You should see a new image under your repository with:
     - Repository name: `${your-repository-name}`
     - Tags include:
       - `latest` – represents the most recent successful build; ideal for quick testing or deployment.
       - `<short-sha>` – a unique identifier based on the commit hash, useful for version tracking, rollbacks, and traceability.

> [!TIP] Protect your main branch
> To maintain code quality and prevent accidental direct pushes, enable branch protection rules:
>
> - From your GitHub repository, go to **Settings → Branches**.
> - Under Branch protection rules, select **Add rule**.
> - Specify `main` as the branch name.
> - Enable options like:
>   - _Require a pull request before merging_.
>   - _Require status checks to pass before merging_.
>
> This ensures that only tested and reviewed code is merged into `main` branch.

---

## Summary

In this section, you set up a comprehensive CI/CD pipeline for your containerized Node.js application using GitHub Actions.

What you accomplished:

- Created a new GitHub repository specifically for your project.
- Generated a Docker Hub access token and added it as a GitHub secret.
- Created a GitHub Actions workflow that:
  - Builds your application in a Docker container.
  - Run tests in a containerized environment.
  - Pushes an image to Docker Hub if tests pass.
- Verified the workflow runs successfully.

Your Node.js application now has automated testing and deployment.

---

## Related resources

Deepen your understanding of automation and best practices for containerized apps:

- [Introduction to GitHub Actions](/guides/gha.md) – Learn how GitHub Actions automate your workflows
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md) – Set up container builds with GitHub Actions
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions) – Full reference for writing GitHub workflows
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry) – Learn about GHCR features and usage
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Optimize your image for performance and security

---

## Next steps

Next, learn how you can deploy your containerized Node.js application to Kubernetes with production-ready configuration. This helps you ensure your application behaves as expected in a production-like environment, reducing surprises during deployment.
