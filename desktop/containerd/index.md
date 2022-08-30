---
title: Overview
description: Viewing the containerd integration in Docker for image and file system management
keywords: Docker, containerd, engine, image store, lazy-pull
toc_min: 1
toc_max: 2  
---

This page presents information about the on-going integration of containerd for image and file system management in the Docker Engine.

> **Warning**
>
> The containerd image store management feature is currently experimental. This may change or be removed from future releases.
{: .warning }

## 2022-09-01
Beta release support of containerd for image management as an Experimental Feature in Docker Desktop 4.12.0.

### New
Initial implementation of the Docker commands: `run`, `commit`, `build`, `push`, `load`, `search` and `save`.

### Known issues
- The containerd integration requires Buildx version 0.9.0 or newer.  On Docker Desktop for Linux (DD4L), validate the locally installed version meets this requirement.  If an older version is installed, the Docker daemon will report an error: `multiple platforms feature is currently not supported for docker driver. Please switch to a different driver`. You can install a newer version of Buildx following the instructions at [Docker Buildx Manual download](https://docs.docker.com/build/buildx/install/#manual-download).
- The containerd integration feature and Kubernetes cluster support in Docker Desktop 4.12.0 are incompatible at the moment. Disable the containerd integration feature from the Experimental features tab if you are using the Kubernetes   from Docker Desktop.

## Feedback

Thanks for trying the new features available with containerd. We’d love to hear from you. You can provide feedback and report any bugs through the Issues tracker in the
[feedback form](https://dockr.ly/3PODIhD){: target="_blank" rel="noopener" class="_"}.
