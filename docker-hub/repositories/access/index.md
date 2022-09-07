---
description: Accessing repositories on Docker Hub
keywords: Docker, docker, trusted, registry, accounts, plans, Dockerfile, Docker Hub, webhooks, docs, documentation, collaborators, viewing, searching, starring
title: Accessing repositories
redirect_from:
---

## Collaborators and their role

A collaborator is someone you want to give access to a private repository. Once
designated, they can `push` and `pull` to your repositories. They are not
allowed to perform any administrative tasks such as deleting the repository or
changing its status from private to public.

> **Note**
>
> A collaborator cannot add other collaborators. Only the owner of
> the repository has administrative access.

You can also assign more granular collaborator rights ("Read", "Write", or
"Admin") on Docker Hub by using organizations and teams. For more information
see the [organizations documentation](orgs.md).


## Viewing repository tags

Docker Hub's individual repositories view shows you the available tags and the
size of the associated image. Go to the **Repositories** view and click on a
repository to see its tags.

![Repository View](/docker-hub/images/repos-create.png)

![View Repo Tags](/docker-hub/images/repo-overview.png)

Image sizes are the cumulative space taken up by the image and all its parent
images. This is also the disk space used by the contents of the `.tar` file
created when you `docker save` an image.

To view individual tags, click on the **Tags** tab.

![Manage Repo Tags](images/repo-tags-list.png)

An image is considered stale if there has been no push/pull activity for more
than 1 month, i.e.:

* It has not been pulled for more than 1 month
* And it has not been pushed for more than 1 month

A multi-architecture image is considered stale if all single-architecture images
part of its manifest are stale.

To delete a tag, select the corresponding checkbox and select **Delete** from the
**Action** drop-down list.

> **Note**
>
> Only a user with administrative access (owner or team member with Admin
> permission) over the repository can delete tags.

Select a tag's digest to view details.

![View Tag](images/repo-image-layers.png)

## Searching for Repositories

You can search the [Docker Hub](https://hub.docker.com) registry through its
search interface or by using the command line interface. Searching can find
images by image name, username, or description:

```console
$ docker search centos

NAME                                 DESCRIPTION                                     STARS     OFFICIAL   AUTOMATED
centos                               The official build of CentOS.                   1034      [OK]
ansible/centos7-ansible              Ansible on Centos7                              43                   [OK]
tutum/centos                         Centos image with SSH access. For the root...   13                   [OK]
...
```

There you can see two example results: `centos` and `ansible/centos7-ansible`.
The second result shows that it comes from the public repository of a user,
named `ansible/`, while the first result, `centos`, doesn't explicitly list a
repository which means that it comes from the top-level namespace for
[Docker Official Images](official_images.md). The `/` character separates
a user's repository from the image name.

Once you've found the image you want, you can download it with `docker pull <imagename>`:

```console
$ docker pull centos

latest: Pulling from centos
6941bfcbbfca: Pull complete
41459f052977: Pull complete
fd44297e2ddb: Already exists
centos:latest: The image you are pulling has been verified. Important: image verification is a tech preview feature and should not be relied on to provide security.
Digest: sha256:d601d3b928eb2954653c59e65862aabb31edefa868bd5148a41fa45004c12288
Status: Downloaded newer image for centos:latest
```

You now have an image from which you can run containers.

## Starring Repositories

Your repositories can be starred and you can star repositories in return. Stars
are a way to show that you like a repository. They are also an easy way of
bookmarking your favorites.
