---
title: Deploy Interlock offline
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
---

To install Interlock on a Docker cluster without internet access the Docker images will
need to be loaded.  This guide will show how to export the images from a local Docker
engine to then be loaded to the Docker Swarm cluster.

First, using an existing Docker engine save the images:

```bash
$> docker save docker/interlock:latest > interlock.tar
$> docker save docker/interlock-extension-nginx:latest > interlock-extension-nginx.tar
$> docker save nginx:alpine > nginx.tar
```

Note: replace `docker/interlock-extension-nginx:latest` and `nginx:alpine` with the corresponding
extension and proxy image if you are not using Nginx.

You should have two files:

- `interlock.tar`: This is the core Interlock application
- `interlock-extension-nginx.tar`: This is the Interlock extension for Nginx
- `nginx:alpine`: This is the official Nginx image based on Alpine

Copy these files to each node in the Docker Swarm cluster and run the following to load each image:

```bash
$> docker load < interlock.tar
$> docker load < interlock-extension-nginx.tar
$> docker load < nginx:alpine.tar
```

After running on each node, you can continue to the [Deployment](index.md#deployment) section to
continue the installation.
