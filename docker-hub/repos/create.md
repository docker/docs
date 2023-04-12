---
description: Creating repositories on Docker Hub
keywords: Docker, docker, trusted, registry, accounts, plans, Dockerfile, Docker Hub, webhooks, docs, documentation, manage, repos
title: Create repositories
redirect_from:
- /docker-hub/repos/configure/
---

Repositories let you share container images with your team,
customers, or the Docker community at large.

A single Docker Hub repository can hold many Docker images which are stored as **tags**. Docker images are pushed to Docker Hub through the [`docker push`](/engine/reference/commandline/push/)
command.

## Create a repository

1. Sign in to Docker Hub.
2. Select **Repositories**.
3. Near the top-right corner, select **Create Repository**.

When creating a new repository:

- You can choose to locate it under your own user account, or under any
  [organization](../../docker-hub/orgs.md) where you are an [owner](../manage-a-team.md#the-owners-team).
- The repository name needs to:
    - Be unique 
    - Have between 2 and 255 characters
    - Only contain lowercase letters, numbers, hyphens (`-`), and underscores (`_`)

  > **Note**
  >
  > You can't rename a Docker Hub repository once it's created.

- The description can be up to 100 characters. It is used in the search results.
- If you are a Docker Verified Publisher (DVP) or Docker-Sponsored Open Source (DSOS) organization, you can also add a logo to a repository. The maximum size is 1000x1000.
- You can link a GitHub or Bitbucket account now, or choose to do it later in
  the repository settings.
- You can set the repository's default visibility to public or private.

  > **Note**
  >
  > For organizations creating a new repository, it's recommended you select **Private**.

## Push a Docker container image to Docker Hub

Once you have created a repository, you can start using `docker push` to push
images.

To push an image to Docker Hub, you must first name your local image using your
Docker Hub username and the repository name that you created.

If you want to add multiple images to a repository, add a specific `:<tag>` to them, for example `docs/base:testing`. If it's not specified, the tag defaults to `latest`.

Name your local images using one of these methods:

- When you build them, using `docker build -t <hub-user>/<repo-name>[:<tag>]`
- By re-tagging an existing local image `docker tag <existing-image> <hub-user>/<repo-name>[:<tag>]`
- By using `docker commit <existing-container> <hub-user>/<repo-name>[:<tag>]` to commit changes

Now you can push this image to the repository designated by its name or tag:

```console
$ docker push <hub-user>/<repo-name>:<tag>
```

The image is then uploaded and available for use by your teammates and/or the community.