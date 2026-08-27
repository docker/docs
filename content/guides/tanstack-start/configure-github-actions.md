---
title: Automate your builds with GitHub Actions
linkTitle: GitHub Actions CI
weight: 60
keywords: CI/CD, GitHub Actions, TanStack Start
description: Learn how to configure CI/CD using GitHub Actions for your TanStack Start application.
---

## Prerequisites

Complete all the previous sections of this guide, starting with
[Containerize TanStack Start application](containerize.md).

You must also have:

- A [GitHub](https://github.com/signup) account.
- A verified [Docker Hub](https://hub.docker.com/signup) account.

---

## Overview

In this section, you'll set up a CI/CD pipeline using
[GitHub Actions](https://docs.github.com/en/actions) to automatically:

- Build your TanStack Start application inside a Docker container.
- Run tests in a consistent environment.
- Push the production-ready image to [Docker Hub](https://hub.docker.com).

---

## Integrate GitHub and Docker Hub

To enable GitHub Actions to build and push Docker images, you'll securely
store your Docker Hub credentials in your GitHub repository.

### Step 1: Connect your GitHub repository to Docker Hub

1. Create a Personal Access Token (PAT) from [Docker Hub](https://hub.docker.com)
   1. Go to your **Docker Hub account → Account Settings → Security**.
   2. Generate a new Access Token with **Read/Write** permissions.
   3. Name it something like `tanstack-start-sample`.
   4. Copy and save the token — you'll need it in Step 4.

2. Create a repository in [Docker Hub](https://hub.docker.com/repositories/)
   1. Go to your **Docker Hub account → Create a repository**.
   2. For the Repository Name, use something descriptive — for example:
      `tanstack-start-sample`.
   3. Once created, copy and save the repository name — you'll need it in
      Step 4.

3. Create a new [GitHub repository](https://github.com/new) for your TanStack
   Start project.

4. Add Docker Hub credentials as GitHub repository secrets

   In your GitHub repository:
   1. Navigate to:
      **Settings → Secrets and variables → Actions → New repository secret**.

   2. Add the following secrets:

   | Name                     | Value                                               |
   | ------------------------ | --------------------------------------------------- |
   | `DOCKER_USERNAME`        | Your Docker Hub username                            |
   | `DOCKERHUB_TOKEN`        | Your Docker Hub access token (created in Step 1)    |
   | `DOCKERHUB_PROJECT_NAME` | Your Docker Hub repository name (created in Step 2) |

   These secrets let GitHub Actions authenticate securely with Docker Hub
   during automated workflows.

5. Connect your local project to GitHub

   Link your local project to the GitHub repository you created:

   ```console
   $ git remote set-url origin https://github.com/{your-username}/{your-repository-name}.git
   ```

   > [!IMPORTANT]
   > Replace `{your-username}` and `{your-repository-name}` with your actual
   > GitHub username and repository name.

   Confirm the remote is configured:

   ```console
   $ git remote -v
   ```

6. Push your source code to GitHub

   Stage, commit, and push your project files to the `main` branch. Once
   completed, your code is on GitHub and configured workflows run
   automatically.

> [!NOTE]
> Learn more about the Git commands used in this step:
>
> - [Git add](https://git-scm.com/docs/git-add) – Stage changes for commit
> - [Git commit](https://git-scm.com/docs/git-commit) – Save staged changes
> - [Git push](https://git-scm.com/docs/git-push) – Upload commits to GitHub
> - [Git remote](https://git-scm.com/docs/git-remote) – Manage remote URLs

---

### Step 2: Set up the workflow

Create a GitHub Actions workflow that builds your Docker image, runs tests,
and pushes the image to Docker Hub.

1. Go to your repository on GitHub and select the **Actions** tab.

2. Select **Set up a workflow yourself**.

   By default, the file is saved to `.github/workflows/main.yml`.

3. Add the following workflow configuration:

```yaml
# CI/CD – TanStack Start Application with Docker
# Builds the app, runs tests in a container, and pushes the production image to Docker Hub.

name: CI/CD – TanStack Start Application with Docker

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]

jobs:
  build-test-push:
    name: Build, Test and Push Docker Image
    runs-on: ubuntu-latest

    steps:
      - name: Checkout source code
        uses: actions/checkout@v5
        with:
          fetch-depth: 0

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v4

      - name: Cache Docker layers
        uses: actions/cache@v5
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: ${{ runner.os }}-buildx-

      - name: Extract metadata
        id: meta
        run: |
          echo "REPO_NAME=${GITHUB_REPOSITORY##*/}" >> "$GITHUB_OUTPUT"
          echo "SHORT_SHA=${GITHUB_SHA::7}" >> "$GITHUB_OUTPUT"

      - name: Build Docker image for tests
        uses: docker/build-push-action@v6
        with:
          context: .
          file: Dockerfile.dev
          tags: ${{ steps.meta.outputs.REPO_NAME }}-dev:latest
          load: true
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max

      - name: Run tests
        run: |
          docker run --rm \
            --workdir /app \
            --entrypoint "" \
            -e CI=true \
            ${{ steps.meta.outputs.REPO_NAME }}-dev:latest \
            npm run test
        env:
          CI: true
          NODE_ENV: test
        timeout-minutes: 10

      - name: Log in to Docker Hub
        if: github.event_name == 'push' && github.ref == 'refs/heads/main'
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build and push production image
        if: github.event_name == 'push' && github.ref == 'refs/heads/main'
        uses: docker/build-push-action@v6
        with:
          context: .
          file: Dockerfile
          push: true
          platforms: linux/amd64,linux/arm64
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:latest
            ${{ secrets.DOCKER_USERNAME }}/${{ secrets.DOCKERHUB_PROJECT_NAME }}:${{ steps.meta.outputs.SHORT_SHA }}
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache,mode=max
```

This workflow:

- Triggers on every `push` or `pull request` targeting the `main` branch.
- Builds a development Docker image using `Dockerfile.dev` for testing.
- Runs Vitest inside a container with `npm run test`.
- Halts the workflow if any test fails.
- Caches Docker build layers for faster CI runs.
- Authenticates with Docker Hub using GitHub repository secrets.
- Builds and pushes a production image with `latest` and short SHA tags on
  pushes to `main`.

> [!NOTE]
> For more information about `docker/build-push-action`, see the
> [GitHub Action README](https://github.com/docker/build-push-action/blob/master/README.md).

---

### Step 3: Run the workflow

1. Commit and push your workflow file from the GitHub editor or your local
   repository.
2. Open the **Actions** tab and select the workflow run to follow each step.
3. After a successful run on `main`, verify the image on
   [Docker Hub](https://hub.docker.com/repositories).

> [!TIP] Protect your main branch
> To maintain code quality and prevent accidental direct pushes, enable branch
> protection rules:
>
> - Navigate to your **GitHub repo → Settings → Branches**.
> - Under Branch protection rules, select **Add rule**.
> - Specify `main` as the branch name.
> - Enable options like _Require a pull request before merging_ and _Require
>   status checks to pass before merging_.

---

## Summary

In this section, you set up a CI/CD pipeline for your containerized TanStack
Start application using GitHub Actions.

What you accomplished:

- Stored Docker Hub credentials as GitHub repository secrets
- Defined a workflow to build, test, and push your application image
- Triggered and verified the workflow through GitHub Actions

---

## Related resources

- [GitHub Actions documentation](https://docs.github.com/en/actions) – Learn
  about workflows, jobs, and steps
- [Docker Hub](https://hub.docker.com/) – Store and share container images
- [docker/build-push-action](https://github.com/docker/build-push-action) –
  Build and push Docker images in GitHub Actions

## Next steps

In the next section, you'll deploy your TanStack Start application to a local
Kubernetes cluster using Docker Desktop.
