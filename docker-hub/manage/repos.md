---
title: Docker Hub repositories
description: Create and edit Docker Hub repositories
keywords: Docker Hub, repositories, repos
redirect_from:
- /docker-hub/repos/
- /docker-cloud/builds/repos/
---

Repositories in Docker Hub store your Docker images. You can create repositories
and manually [push images](push-images) using `docker push`, or you can link
to a source code provider and use [automated builds](../build/) to
build the images for you. These repositories can be either public or private.

## Create new user repository

All individual users can create one private repository for free, and can create
unlimited public repositories.

1.  Click **Repositories** in the left navigation.

2.  Click **Create**.

3.  Enter a **name** and an optional **description**.

4.  Choose a visibility setting for the repository.

5.  Optionally, click a linked source code provider to set up [automated builds](../build/).

    a. Select a namespace from that source code provider.

    b. From that namespace, select a repository to build.

    c. Optionally, expand the build settings section to set up build rules and enable or disable Autobuilds.

    > You can set up autobuilds later
    >
    > Repos are configurable and you can change build settings at any time after
    > the repository is created. If you choose not to enable automated builds,
    > you can still push images to the repository.

6.  Click **Create**.

### Repositories for organizations

Every organization has and `Owners` team. Members of the `Owners` team can:

- Create new repositories for that organization.
- Configure repo access permissions for other teams in the organization.
- Change the organization billing information.
- Link the organization to a source code provider to set up automated builds.

To learn more, see the [organizations and teams documentation](orgs-teams/).

## Edit an existing repository

You can edit repositories in Docker Hub to change the description and build
configuration.

From the **General** page, edit the repository short description, or click to
edit the version of the ReadMe displayed on the repository page.

> Edits to the Docker Hub **ReadMe** are not reflected in the source code linked to a repository.

## Change repository privacy settings

Repositories in Docker Hub can be either public or private.

Public repositories are visible in  ..., and can be searched in ...

Private repositories are only visible to the user account that created it
(unless it belongs to an organization, see below).

> Privacy settings vs access permissions
>
> _Privacy_ settings for an individual repo differ from differ from
> [_access_ permissions](orgs-teams#edit-permissions-for-individual-repos)
> of a repo shared among members of an [organization](orgs-teams/).

If a private repository belongs to an organization, members of the `Owners` team
can configure access. Only members of the `Owners` team can change an organization's
repository privacy settings.

Each Docker Hub account comes with one free private repository. Additional
private repositories are available for subscribers on paid plans.

To change a repository's privacy settings:

1.  Navigate to the repository in Docker Hub.

2.  Click the **Settings** tab.

3.  Click the **Make public** or **Make private** button.

4.  In the dialog that appears, enter the name of the repository to confirm the change.

5.  Click the button to save the change.

## Link to repo from third party registry

You can link to repositories hosted on a third party registry. This allows you
to enable automated builds and push built images back to the registry.

To link to a repository that you want to share with an organization, contact a
member of the organization's `Owners` team. Only the Owners team can import new
external registry repositories for an organization.

1.  Click **Repositories** in the side menu.

2.  Click the down arrow menu next to the **Create** button.

3.  Select **Import**.

4.  Enter the name of the repository that you want to add.

    For example, `registry.com/namespace/reponame` where `registry.com` is the
    hostname of the registry.

5.  Enter credentials for the registry.

    > Push vs read-only permissions
    >
    > Credentials must have **push** permissions to push built images back to
    > the repository. If you provide **read-only** credentials, you can run
    > automated tests but you cannot push built images to it.

6.  Click **Import**.

7.  Confirm that the repository on the third-party registry now appears in your **Repositories** dropdown list.

## Delete a repository

When you delete a repository in Docker Hub, all of the images in that
repository are also deleted.

If automated builds are configured for the repository, the build rules and
settings are deleted along with any Docker Security Scan results. However, this
does not affect the code in the linked source code repository, and does not
remove the source code provider link.

If you are running a service from deleted repository, the service continues
to run, but cannot be scaled up or redeployed. If any builds use the Docker
`FROM` directive and reference a deleted repository, those builds fail.

To delete a repository:

1.  Navigate to the repository, and click the **Settings** tab.

2.  Click **Delete**.

3.  Enter the name of the repository to confirm deletion, and click **Delete**.

External (third-party) repositories cannot be deleted from within Docker Hub,but
you can remove a link to them. The link is removed, but images in the external
repository are not deleted.

> To delete an organization, you must be a member of the `Owners` team for that organization.

## What's next?

Once you create or link to a repository in Docker Hub, you can set up [automated testing](../build/autotest/) and [automated builds](../build/).
