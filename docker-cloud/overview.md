---
description: Docker Cloud Overview
keywords: Docker, cloud, three
title: Docker Cloud overview
---

Docker Cloud is a hosted service that provides a Registry with build and testing
facilities for Dockerized application images, tools to help you set up and
manage your host infrastructure, and deployment features to help you automate
deploying your images to your infrastructure.

#### Your Docker Cloud account and Docker ID

You log in to Docker Cloud using your free Docker ID. Your Docker ID is the same
set of credentials you used to log in to Docker Hub, and this allows you to
access your Docker Hub repositories from Docker Cloud.

#### Images, Builds, and Testing

Docker Cloud uses Docker Hub as an online registry service. This allows you to
publish Dockerized images on the internet either publicly or privately. Along
with the ability to store pre-built images, Docker Cloud can link to your source
code repositories and manage building and testing your images before pushing the
images.

![](images/cloud-build.png)

#### Infrastructure management

Before you can do anything with images, you need somewhere to run them. Docker
Cloud allows you to link to your infrastructure or cloud services provider which
lets you provision new nodes automatically, and deploy images directly from your
Docker Cloud repositories onto your infrastructure hosts.

![](images/cloud-clusters.png)

#### Services, Stacks, and Applications

Images are just one layer in containerized applications. Once you've built an
image, you can use it to produce containers, which make up a service, or use
Docker Cloud's stackfiles to combine it with other services and microservices,
to form a full application.

![](images/cloud-stack.png)