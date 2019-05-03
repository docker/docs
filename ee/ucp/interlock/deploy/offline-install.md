---
title: Offline installation considerations
description: Learn how to to install Interlock on a Docker cluster without internet access.
keywords: routing, proxy, interlock
---

To install Interlock on a Docker cluster without internet access, the Docker images must be loaded.  This topic describes how to export the images from a local Docker
engine and then loading them to the Docker Swarm cluster.

First, using an existing Docker engine, save the images:

```bash
$> docker save {{ page.ucp_org }}/ucp-interlock:{{ page.ucp_version }} > interlock.tar
$> docker save {{ page.ucp_org }}/ucp-interlock-extension:{{ page.ucp_version }} > interlock-extension-nginx.tar
$> docker save {{ page.ucp_org }}/ucp-interlock-proxy:{{ page.ucp_version }} > nginx.tar
```

Note: replace `{{ page.ucp_org }}/ucp-interlock-extension:{{ page.ucp_version
}}` and `{{ page.ucp_org }}/ucp-interlock-proxy:{{ page.ucp_version }}` with the
corresponding extension and proxy image if you are not using Nginx.

You should have the following two files:

- `interlock.tar`: This is the core Interlock application.
- `interlock-extension-nginx.tar`: This is the Interlock extension for Nginx.
- `nginx:alpine`: This is the official Nginx image based on Alpine.

Copy these files to each node in the Docker Swarm cluster and run the following commands to load each image:

```bash
$> docker load < interlock.tar
$> docker load < interlock-extension-nginx.tar
$> docker load < nginx:alpine.tar
```

## Next steps
After running on each node, refer to the [Deploy](./index.md) section to
continue the installation.
