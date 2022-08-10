---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local, Dockerfile, specify, base image
title: Specify a Dockerfile or base image
---

## Specify a Dockerfile 

Use a JSON file to specify a Dockerfile which in turn defines your Dev Environment. You must include this as part of the `.docker` folder and then add it as a `config.json` file. For example:

```jsx
{
    "dockerfile": "Dockerfile.devenv"
}
```

Next, define the dependencies you want to include in your `Dockerfile.devenv`.

While some images or Dockerfiles include a non-root user, many base images and Dockerfiles do not. Fortunately, you can add a non-root user named `vscode`. If you include Docker tooling, for example the Docker CLI or `docker compose`, in the `Dockerfile.devenv`, you need the `vscode` user to be included in the `docker` group.

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
> This configuration is to unblock users for the Beta release only. We may move this configuration for single and multi-container applications to a Compose-based implementation in future releases.
