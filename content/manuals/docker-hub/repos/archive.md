---
description: Learn how to archive or activate a repository on Docker Hub
keywords: Docker Hub, Hub, repositories, archive, activate
title: Archive or activate a repository
linkTitle: Archive
toc_max: 3
weight: 35
---

You can archive a repository on Docker Hub to mark it as read-only and indicate
that it's no longer actively maintained. This helps prevent the use of outdated
or unsupported images in workflows. Archived repositories can also be unarchived
if needed.

Docker Hub highlights repositories that haven't been updated in over a year by
displaying an icon ({{< inline-image src="./images/outdated-icon.webp"
alt="outdated icon" >}}) next to them on the [**Repositories**
page](https://hub.docker.com/repositories/). Consider reviewing these
highlighted repositories and archiving them if necessary.

When a repository is archived, the following occurs:

- The repository information can't be modified.
- New images can't be pushed to the repository.
- An **Archived** label is displayed on the public repository page.
- Users can still pull the images.

You can activate an archived repository to remove the archived state. When
activated, the following occurs:

- The repository information can be modified.
- New images can be pushed to the repository.
- The **Archived** label is removed on the public repository page.

## Archive a repository

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Settings** tab.
5. Select **Archive repository**.
6. Enter the name of your repository to confirm.
7. Select **Archive**.

## Activate an archived repository

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Settings** tab.
5. Select **Activate image**.