---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local
title: Development Environments (Beta)
---

## Specify a Dockerfile 

In this preview, Dev Environments support a JSON file which allows you to specify a Dockerfile to define your Dev Environment. You must include this as part of the `.docker` folder and then add it as a `config.json` file. For example:

```jsx
{
    "dockerfile": "Dockerfile.devenv"
}
```

Next, you have to define the dependencies you want to include in your `Dockerfile.devenv`, alongside the following requisites:

While some images or Dockerfiles will include a non-root user, many base images and Dockerfiles do not. Fortunately, you can add a non-root user named `vscode`. If you were to include the Docker tooling (e.g. `docker` cli, `docker compose`, etc.) in the `Dockerfile.devenv`, you would need the `vscode` user to be included in the `docker` group.

```dockerfile
# syntax=docker/dockerfile:1

FROM <your base image>

RUN useradd -s /bin/bash -m vscode \
 && groupadd docker \
 && usermod -aG docker vscode

USER vscode
```

## Specify a base image

If you already have an image built, you can specify it as a base image to define your Dev Environment. You must include this as part of the `.docker` folder and then add it as a `config.json` file. For example, to use the Jekyll base image, add:

```jsx
{
  "image": "jekyll/jekyll"
}
```

> **Note**
>
> This configuration is to unblock users for the Preview release only. We may move this configuration for single and multi-container applications to a Compose-based implementation in future releases.
>
> To get involved with the discussion on how we are going to implement this as part of Compose, join the **#docker-dev-environments** channel in the [Docker Community Slack](https://dockercommunity.slack.com/messages){:target="_blank" rel="noopener" class="_"}, or let us know your feedback by creating an issue in the [Dev Environments](https://github.com/docker/dev-environments/issues){:target="_blank" rel="noopener" class="_"} GitHub repository.

