---
title: "Overview"
keywords: docker basics, how to start a docker container, container settings, setup docker, how to setup docker, setting up docker, docker container guide, how to get started with docker
description: Get started with the Docker basics in this comprehensive overview, You'll learn about containers, images, and how to containerize your first application.
redirect_from:
- /engine/getstarted-voting-app/
- /engine/getstarted-voting-app/cleanup/
- /engine/getstarted-voting-app/create-swarm/
- /engine/getstarted-voting-app/customize-app/
- /engine/getstarted-voting-app/deploy-app/
- /engine/getstarted-voting-app/node-setup/
- /engine/getstarted-voting-app/test-drive/
- /engine/getstarted/
- /engine/getstarted/last_page/
- /engine/getstarted/step_five/
- /engine/getstarted/step_four/
- /engine/getstarted/step_one/
- /engine/getstarted/step_six/
- /engine/getstarted/step_three/
- /engine/getstarted/step_two/
- /engine/quickstart/
- /engine/tutorials/
- /engine/tutorials/dockerimages/
- /engine/tutorials/dockerizing/
- /engine/tutorials/usingdocker/
- /engine/userguide/containers/dockerimages/
- /engine/userguide/dockerimages/
- /engine/userguide/intro/
- /get-started/part1/
- /get-started/part5/
- /get-started/part6/
- /getstarted/
- /getting-started/
- /learn/
- /linux/last_page/
- /linux/started/
- /linux/step_four/
- /linux/step_one/
- /linux/step_six/
- /linux/step_three/
- /linux/step_two/
- /mac/last_page/
- /mac/started/
- /mac/step_four/
- /mac/step_one/
- /mac/step_six/
- /mac/step_three/
- /mac/step_two/
- /userguide/dockerimages/
- /userguide/dockerrepos/
- /windows/last_page/
- /windows/started/
- /windows/step_four/
- /windows/step_one/
- /windows/step_six/
- /windows/step_three/
- /windows/step_two/
---

Welcome! We're excited that you want to learn Docker.

This guide contains step-by-step instructions on how to get started with Docker. Some of the things you'll learn and do in this guide are:

- Build and run an image as a container.
- Share images using Docker Hub.
- Deploy Docker applications using multiple containers with a database.
- Run applications using Docker Compose.

Before you get to the hands on part of the guide, you should learn about containers and images.

## What is a container?

Simply put, a container is a sandboxed process on your machine that is isolated from all other processes on the host machine. That isolation leverages [kernel namespaces and cgroups](https://medium.com/@saschagrunert/demystifying-containers-part-i-kernel-space-2c53d6979504),
features that have been in Linux for a long time. Docker has worked to make these capabilities approachable and easy to use. To summarize, a container:

- Is a runnable instance of an image. You can create, start, stop, move, or delete a container using the DockerAPI or CLI.
- Can be run on local machines, virtual machines or deployed to the cloud.
- Is portable (can be run on any OS).
- Is isolated from other containers and runs its own software, binaries, and configurations.

## What is a container image?

When running a container, it uses an isolated filesystem. This custom filesystem is provided by a container image. Since the image contains the container's filesystem, it must contain everything needed to run an application - all dependencies, configurations, scripts, binaries, etc. The image also contains other configuration for the container, such as environment variables, a default command to run, and other metadata.

You'll dive deeper into images later on in this guide, covering topics such as layering, best practices, and more.

> **Note**
>
> If you're familiar with `chroot`, think of a container as an extended version of `chroot`. The filesystem is simply coming from the image. But, a container adds additional isolation not available when simply using chroot.

## Next steps

In this section, you learned about containers and images.

In the next section, you'll containerize your first application.

[Containerize an application](02_our_app.md){: .button  .primary-btn}

