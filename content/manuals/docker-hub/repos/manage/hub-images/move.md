---
description: Discover how to move images between repositories.
keywords: Docker Hub, Hub, repository content, move
title: Move images between repositories
linkTitle: Move images
weight: 40
---

Consolidating and organizing your Docker images across repositories can
streamline your workflows, whether you're managing personal projects or
contributing to an organization. This topic explains how to move images between
Docker Hub repositories, ensuring that your content remains accessible and
organized under the correct accounts or namespaces.

## Personal to personal

When consolidating personal repositories, you can pull private images from the initial repository and push them into another repository owned by you. To avoid losing your private images, perform the following steps:

1. [Sign up](https://app.docker.com/signup) for a new Docker account with a personal subscription.
2. Sign in to [Docker](https://app.docker.com/login) using your original Docker account
3. Pull your images:

   ```console
   $ docker pull namespace1/docker101tutorial
   ```

4. Tag your private images with your newly created Docker username, for example:

   ```console
   $ docker tag namespace1/docker101tutorial new_namespace/docker101tutorial
   ```
5. Using `docker login` from the CLI, sign in with your newly created Docker account, and push your newly tagged private images to your new Docker account namespace:

   ```console
   $ docker push new_namespace/docker101tutorial
   ```

The private images that existed in your previous account are now available in your new account.

## Personal to an organization

To avoid losing your private images, you can pull your private images from your
personal account and push them to an organization that's owned by you.

1. Navigate to [Docker Hub](https://hub.docker.com) and select **Organizations**.
2. Select the applicable organization and verify that your user account is a member of the organization.
3. Sign in to [Docker Hub](https://hub.docker.com) using your original Docker account, and pull your images:

   ```console
   $ docker pull namespace1/docker101tutorial
   ```
4. Tag your images with your new organization namespace:

   ```console
   $ docker tag namespace1/docker101tutorial <new_org>/docker101tutorial
   ```
5. Push your newly tagged images to your new org namespace:

   ```console
   $ docker push new_org/docker101tutorial
   ```

The private images that existed in your user account are now available for your organization.