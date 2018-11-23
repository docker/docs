---
description: Using repositories on Docker Hub
keywords: Docker, docker, trusted, registry, accounts, plans, Dockerfile, Docker Hub, webhooks, docs, documentation
title: Repositories
---

Docker Hub repositories let you share container images with your team,
customers, or the Docker community at large.

- Repositories hold Docker container images:
- One Docker Hub repository can hold many Docker images
- Docker images are pushed to Docker Hub via the [`docker push`](https://docs.docker.com/engine/reference/commandline/push/) command.
- Each image pushed to Docker Hub must have a **tag**
- Tags are named when images are pushed to Docker Hub (e.g. `latest`, `v1.0.0`, `1.0.0`)

## Creating Repositories


## Pushing a Docker container image to Docker Hub

To push a repository to the Docker Hub, you need to
name your local image using your Docker Hub username, and the
repository name that you created in the previous step.
You can add multiple images to a repository, by adding a specific `:<tag>` to
it (for example `docs/base:testing`). If it's not specified, the tag defaults to
`latest`.
You can name your local images either when you build it, using
`docker build -t <hub-user>/<repo-name>[:<tag>]`,
by re-tagging an existing local image `docker tag <existing-image> <hub-user>/<repo-name>[:<tag>]`,
or by using `docker commit <exiting-container> <hub-user>/<repo-name>[:<tag>]` to commit
changes.

Now you can push this repository to the registry designated by its name or tag.

    $ docker push <hub-user>/<repo-name>:<tag>

The image is then uploaded and available for use by your teammates and/or
the community.

## Private Repositories

Private repositories allow you to have repositories that contain images that you
want to keep private, either to your own account or within an organization or
team.

To work with a private repository on [Docker Hub](https://hub.docker.com), you
need to add one using the [Add Repository](https://hub.docker.com/add/repository/) button. You get one private
repository for free with your Docker Hub user account (not usable for
organizations you're a member of). If you need more private repositories for your user account, upgrade
your Docker Hub plan from your [Billing Information](https://hub.docker.com/account/billing-plans/) page.

Once the private repository is created, you can `push` and `pull` images to and
from it using Docker.

> **Note**: You need to be signed in and have access to work with a
> private repository.

Private repositories are just like public ones. However, it isn't possible to
browse them or search their content on the public registry. They do not get
cached the same way as a public repository either.

You can designate collaborators and manage their access to a private
repository from that repository's *Settings* page. You can also toggle the
repository's status between public and private, if you have an available
repository slot open. Otherwise, you can upgrade your
[Docker Hub](https://hub.docker.com/account/billing-plans/) plan.

## Public Repositories

## Editing Repository information

## Collaborators and their role

A collaborator is someone you want to give access to a private repository. Once
designated, they can `push` and `pull` to your repositories. They are not
allowed to perform any administrative tasks such as deleting the repository or
changing its status from private to public.

> **Note**:
> A collaborator cannot add other collaborators. Only the owner of
> the repository has administrative access.

You can also assign more granular collaborator rights ("Read", "Write", or
"Admin") on Docker Hub by using organizations and teams. For more information
see the [organizations documentation](/docker-hub/orgs.md).


## Viewing repository tags

Docker Hub's repository "Tags" view shows you the available tags and the size
of the associated image.

Image sizes are the cumulative space taken up by the image and all its parent
images. This is also the disk space used by the contents of the Tar file created
when you `docker save` an image.

![images/busybox-image-tags.png](/docker-hub/images/busybox-image-tags.png)

## Creating a new repository on Docker Hub

When you first create a Docker Hub user, you see a "Get started with
Docker Hub." screen, from which you can click directly into "Create Repository".
You can also use the "Create &#x25BC;" menu to "Create Repository".

When creating a new repository, you can choose to put it in your Docker ID
namespace, or that of any [organization](/docker-hub/orgs.md) that you are in the "Owners"
team. The Repository Name needs to be unique in that namespace, can be two
to 255 characters, and can only contain lowercase letters, numbers or `-` and
`_`.

The "Short Description" of 100 characters is used in the search results,
while the "Full Description" can be used as the Readme for the repository, and
can use Markdown to add simple formatting.

After you hit the "Create" button, you then need to `docker push` images to that
Hub based repository.

<!-- TODO: show a created example, and then use it in subsequent sections -->

## Searching for Repositories

You can search the [Docker Hub](https://hub.docker.com) registry via its search
interface or by using the command line interface. Searching can find images by
image name, user name, or description:

    $ docker search centos
    NAME                                 DESCRIPTION                                     STARS     OFFICIAL   AUTOMATED
    centos                               The official build of CentOS.                   1034      [OK]
    ansible/centos7-ansible              Ansible on Centos7                              43                   [OK]
    tutum/centos                         Centos image with SSH access. For the root...   13                   [OK]
    ...

There you can see two example results: `centos` and `ansible/centos7-ansible`.
The second result shows that it comes from the public repository of a user,
named `ansible/`, while the first result, `centos`, doesn't explicitly list a
repository which means that it comes from the top-level namespace for [Official
Images](/docker-hub/official_images.md). The `/` character separates a user's
repository from the image name.

Once you've found the image you want, you can download it with `docker pull <imagename>`:

    $ docker pull centos
    latest: Pulling from centos
    6941bfcbbfca: Pull complete
    41459f052977: Pull complete
    fd44297e2ddb: Already exists
    centos:latest: The image you are pulling has been verified. Important: image verification is a tech preview feature and should not be relied on to provide security.
    Digest: sha256:d601d3b928eb2954653c59e65862aabb31edefa868bd5148a41fa45004c12288
    Status: Downloaded newer image for centos:latest

You now have an image from which you can run containers.


## Starring Repositories

Your repositories can be starred and you can star repositories in return. Stars
are a way to show that you like a repository. They are also an easy way of
bookmarking your favorites.
