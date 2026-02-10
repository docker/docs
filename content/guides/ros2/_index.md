---
title: Introduction to ROS 2 Development with Docker
linkTitle: ROS 2
description: Learn how to containerize and develop ROS 2 applications using Docker.
keywords: ros2, robotics, devcontainers, python, cpp, Dockerfile, rviz
summary: |
  This guide details how to containerize ROS 2 applications using Docker.
toc_min: 1
toc_max: 2
languages: []
params:
  tags: [frameworks]
  time: 30 minutes
---

> **Acknowledgment**
>
> This guide is a community contribution. Docker would like to thank
> [Shakirth Anisha](https://www.linkedin.com/in/shakirth-anisha/) for her contribution
> to this guide.

[ROS 2](https://www.ros.org/) is a set of software libraries and tools for building robot applications. It uses Data Distribution Service (DDS) for real-time, secure communication between distributed nodes, making it ideal for robotics and autonomous systems.

---

## What will you learn?

In this guide, you'll learn how to:

- Use official ROS 2 base images from Docker Hub
- Run ROS 2 in an Ubuntu container
- Install ROS 2 packages and dependencies
- Set up a development container for local development
- Run a complete end-to-end example with Turtlesim

## Prerequisites

Before you begin, make sure you're familiar with the following:

- [Docker Desktop](https://docs.docker.com/desktop/): You must have Docker Desktop installed and running.
- [Homebrew](https://brew.sh/) (macOS users): If you are using macOS, you must have Homebrew installed to manage dependencies.
- [Docker concepts](/get-started/docker-concepts/the-basics/what-is-a-container.md): You must understand core Docker concepts, such as images and containers.
- [ROS 2 concepts](https://www.ros.org): Basic understanding of concepts like nodes, packages, topics, and services.

## What's next?

Start by setting up your ROS 2 development environment using Docker and dev containers.
