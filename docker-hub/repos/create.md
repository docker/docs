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

  <div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample1" style="cursor: pointer">
    Description best practices
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample1">
      <p>A good image description is essential to help potential users understand why and how to use the image. The bellow covers the best practices to follow when adding a description to your image.</p>
      <h3>Describe the image</h3>
      <p>Include a description of what the image is, the features it offers, and why people might want to use it in their project or image.</p>
      <p>Optional information to include are examples of usage, history, the team behind the project, etc.</p>
      <h3>How to start and use the image</h3>
      <p>Provide instructions for a minimal “getting started” example with running a container using the image. Also provide a minimal example of how to use the image in a Dockerfile.</p>
      <p>Optional information to include, if relevant, is:</p>
      <ul>
      <li>Exposing ports</li>
      <li>Default environment variables</li>
      <li>Non-root user</li>
      <li>Running in debug mode</li>
      </ul>
      <h3>Image tags and variants</h3>
      <p>List the key image variants and the tags for using them along with what that variant offers and why someone might want to use that variant.</p>
      <h3>Where to find more information</h3>
      <p>Add links here for docs and support sites, communities, or mailing lists were users can find more and ask questions.</p>
      <h3>Who the image is maintained by</h3>
      <p>Who is the image maintained by and how can someone contact them with concerns.</p>
      <h3>License</h3>
      <p>What is the license for the image and where can people find more details if needed.</p>
    </div>

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