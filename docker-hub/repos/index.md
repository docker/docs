---
description: Using repositories on Docker Hub
keywords: Docker, docker, trusted, registry, accounts, plans, Dockerfile, Docker Hub, webhooks, docs, documentation, manage, repos
title: Manage repositories
redirect_from:
- /engine/tutorials/dockerrepos/
---

## Consolidating a repository

### Personal to personal

When consolidating personal repositories, you can pull private images from the initial repository and push them into another repository owned by you. To avoid losing your private images, perform the following steps:

1. Navigate to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} create a Docker ID and select the personal subscription.
2. Using `docker login` from the CLI, sign in using your original Docker ID and pull your private images.
3. Tag your private images with your newly created Docker ID using:
`docker tag namespace1/docker101tutorial new_namespace/docker101tutorial`
4. Using `docker login` from the CLI, sign in with your newly created Docker ID, and push your newly tagged private images to your new Docker ID namespace.
`docker push new_namespace/docker101tutorial`
5. The private images that existed in your previous namespace are now available in your new Docker ID namespace.

### Personal to an organization

To avoid losing your private images, you can pull your private images from your personal namespace and push them to an organization that's owned by you.

1. Navigate to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and select **Organizations**.
2. Select the applicable organization and verify that your user account is a member of the organization.
3. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} using your original Docker ID, and pull your images from the initial namespace.
`docker pull namespace1/docker101tutorial`
4. Tag your images with your new organization namespace.
`docker tag namespace1/docker101tutorial <new_org>/docker101tutorial`
5. Push your newly tagged images to your new org namespace.
`docker push new_org/docker101tutorial`

The private images that existed in the initial namespace are now available for your organization.

## Change a repository from public to private

> **Note**
>
> To update your public repository to private, navigate to your repository, select **Settings** and **Make private**.

## Deleting a repository

1. Navigate to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and select **Repositories**.

2. Select a repository from the list, select **Settings**, and then Delete Repository.

    > **Note:**
    >
    > Deleting a repository deletes all the images it contains and its build settings. This action can't be undone.

3. Enter the name of the repository to confirm the deletion and select **Delete**.