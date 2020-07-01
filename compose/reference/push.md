---
description: Pushes service images.
keywords: fig, composition, compose, docker, orchestration, cli,  push
title: docker-compose push
notoc: true
---

```
Usage: push [options] [SERVICE...]

Options:
    --ignore-push-failures  Push what it can and ignores images with push failures.
```

Pushes images for services to their respective `registry/repository`.

The following assumptions are made:

- You are pushing an image you have built locally

- You have access to the build key

## Example

```yaml
version: '3'
services:
  service1:
    build: .
    image: localhost:5000/yourimage  # goes to local registry

  service2:
    build: .
    image: your-dockerid/yourimage  # goes to your repository on Docker Hub
```
