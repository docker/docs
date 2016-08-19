---
description: Create and edit Docker Cloud repositories
keywords:
- Docker Cloud repositories, automated, build, images
menu:
  main:
    parent: builds
    weight: -70
title: Docker Cloud repositories
---

# Docker Cloud repositories

Repositories in Docker Cloud store your Docker images. You can create
repositories and manually [push images](push-images.md) using `docker push`, or
you can link to a source code provider and use [automated builds](automated-build.md) to automatically build images. These repositories
can be either public or private.

Additionally, you can access your Docker Hub repositories and automated builds
from within Docker Cloud.

## Create a new repository in Docker Cloud

To store your images in Docker Cloud, you create a repository.

1. Click **Repositories** in the side menu.
2. Click **Create**.
3. Enter a **name** and an optional **description**, and click **Create**.
  ![](images/create-repository.png)

Once you create a repository, you can either [configure automated builds](automated-build.md#configure-automated-build-settings) from the UI, or push to it manually using the `docker` or `docker-cloud` CLI.

## Edit an existing repository in Docker Cloud

You can edit your repositories in Docker Cloud to change the description and
build configuration.

1. Click **Repositories** in the left menu.
2. Click the existing repository that you want to edit.
3. From the repository details page, click **Edit repository**.

To enable or change the automated build settings, click **Configure Automated
Builds**. See the Automated build documentation on [configuring automated build settings](automated-build.md#configure-automated-build-settings) for more
information.

## Link to a repository from a third party registry

You can link to repositories hosted on a third party registry. This allows you to deploy images from the third party registry to your nodes in Docker Cloud, and also allows you to enable automated builds which push built images back to the registry.

1. Click **Repositories** in the side menu.

2. Click the down arrow menu next to the **Create** button.

3. Select **Import**.

4. Enter the name of the repository that you want to add.

    For example, `registry.com/namespace/reponame` where `registry.com` is the
    hostname of the registry.
    ![](images/third-party-images-modal.png)

5. Enter credentials for the registry.

    > **Note**: These credentials must have **push** permission in order to push
    built images back to the repository. If you provide **read-only**
    credentials, you will be able to run automated tests and deploy from the
    repository to your nodes, but you will not be able to push built images to
    it.

6. Click **Import**.

7. Confirm that the repository on the third-party registry now appears in your **Repositories** dropdown list.

## What's next?

Once you create or link to a repository in Docker Cloud, you can set up [automated testing](automated-testing.md) and [automated builds](automated-build.md).
