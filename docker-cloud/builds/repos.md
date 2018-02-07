---
description: Create and edit Docker Cloud repositories
keywords: Docker Cloud repositories, automated, build, images
title: Docker Cloud repositories
---

Repositories in Docker Cloud store your Docker images. You can create
repositories and manually [push images](push-images.md) using `docker push`, or
you can link to a source code provider and use [automated builds](automated-build.md) to build the images for you. These repositories
can be either public or private.

![Docker Cloud repository General view](images/repo-general.png)

Additionally, you can access your Docker Hub repositories and automated builds
from within Docker Cloud.

## Create a new repository in Docker Cloud

To store your images in Docker Cloud, you create a repository. All individual users can create one private repository for free, and can create unlimited public repositories.

1. Click **Repositories** in the left navigation.
2. Click **Create**.
3. Enter a **name** and an optional **description**.
4. Choose a visibility setting for the repository.
5. Optionally, click a linked source code provider to set up automated builds.
    1. Select a namespace from that source code provider.
    2. From that namespace, select a repository to build.
    3. Optionally, expand the build settings section to set up build rules and enable or disable Autobuilds.

    > **Note**: You do not need to set up automated builds right away, and you can change the build settings at any time after the repository is created. If you choose not to enable automated builds, you can still push images to the repository using the `docker` or `docker-cloud` CLI.
6. Click **Create**.

  ![Create repository page](images/create-repository.png)

### Repositories for Organizations

Only members of an organization's `Owners` team can create new repositories for
the organization. Members of `Owners` can also change the organization's billing
information, and link the organization to a source code provider to set up
automated builds.

A member of the `Owners` team must also set up the repository's access
permissions so that other teams within the organization can use it. To learn
more, see the [organizations and teams documentation](../orgs.md#set-team-permissions).

## Edit an existing repository in Docker Cloud

You can edit repositories in Docker Cloud to change the description and build configuration.

From the **General** page, you can edit the repository's short description, or click to edit the version of the ReadMe displayed on the repository page.

> **Note**: Edits to the Docker Cloud **ReadMe** are not reflected in the source code linked to a repository.

To run a build, or to set up or change automated build settings, click the **Builds** tab, and click **Configure Automated Builds**. See the documentation on [configuring automated build settings](automated-build.md#configure-automated-build-settings) for more
information.

## Change repository privacy settings

Repositories in Docker Cloud can be either public or private. Public
repositories are visible from the Docker Store's Community Content section, and
can also be searched for from Docker Cloud's **Create Service** wizard. Private
repositories are only visible to the user account that created it (unless it
belongs to an Organization, see below).

> **Note**: These _privacy_ settings are separate from the [repository _access_ permissions](../orgs.md#change-team-permissions-for-an-individual-repository) available for repositories shared among members of an [organization](../orgs.md).

If a private repository belongs to an [Organization](../orgs.md), members of the
`Owners` team configure access. Only members of the `Owners` team can change an
organization's repository privacy settings.

Each Docker Cloud account comes with one free private repository. Additional
private repositories are available for subscribers on paid plans.

To change a repository's privacy settings:

1. Navigate to the repository in Docker Cloud.
2. Click the **Settings** tab.
3. Click the **Make public** or **Make private** button.
4. In the dialog that appears, enter the name of the repository to confirm the change.
5. Click the button to save the change.

## Delete a repository

When you delete a repository in Docker Cloud, all of the images in that
repository are also deleted.

If automated builds are configured for the repository, the build rules and
settings are deleted along with any Docker Security Scan results. However, this
does not affect the code in the linked source code repository, and does not
remove the source code provider link.

If you are running a service from deleted repository , the service continues
to run, but cannot be scaled up or redeployed. If any builds use the Docker
`FROM` directive and reference a deleted repository, those builds fail.

To delete a repository:

1. Navigate to the repository, and click the **Settings** tab.
2. Click **Delete**.
3. Enter the name of the repository to confirm deletion, and click **Delete**.

External (third-party) repositories cannot be deleted from within Docker Cloud,
however you can remove a link to them using the same process for a repository in
Docker Cloud. The link is removed, but images in the external repository are not
deleted.

> **Note**: If the repository to be deleted or removed belongs to an [Organization](../orgs.md), only members of the `Owners` team can delete it.


## Link to a repository from a third party registry

You can link to repositories hosted on a third party registry. This allows you
to deploy images from the third party registry to nodes in Docker Cloud, and
also allows you to enable automated builds which push built images back to the
registry.

> **Note**: To link to a repository that you want to share with an organization, contact a member of the organization's `Owners` team. Only the Owners team can import new external registry repositories for an organization.

1. Click **Repositories** in the side menu.

2. Click the down arrow menu next to the **Create** button.

3. Select **Import**.

4. Enter the name of the repository that you want to add.

    For example, `registry.com/namespace/reponame` where `registry.com` is the
    hostname of the registry.
    ![Import repository popup](images/third-party-images-modal.png)

5. Enter credentials for the registry.

    > **Note**: These credentials must have **push** permission to push
    built images back to the repository. If you provide **read-only**
    credentials, you can run automated tests and deploy from the
    repository to your nodes, but you cannot push built images to
    it.

6. Click **Import**.

7. Confirm that the repository on the third-party registry now appears in your **Repositories** dropdown list.

## What's next?

Once you create or link to a repository in Docker Cloud, you can set up [automated testing](automated-testing.md) and [automated builds](automated-build.md).
