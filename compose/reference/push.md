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

_Note: It is assumed you are building the image and should have the `build` key._

For example

```
version: '3'
services:
  service1:
    build: .
    image: localhost:5000/yourimage  # goes to local registry

  service2:
    build: .
    image: youruser/yourimage  # goes to youruser DockerHub registry
```
