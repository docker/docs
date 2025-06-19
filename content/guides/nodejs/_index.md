---
title: Node.js language-specific guide
linkTitle: Node.js
description: Containerize and develop Node.js apps using Docker
keywords: getting started, node, node.js
summary: |
  This guide explains how to containerize Node.js applications using Docker.
toc_min: 1
toc_max: 2
aliases:
  - /language/nodejs/
  - /guides/language/nodejs/
languages: [js]
params:
  time: 20 minutes
---

[Node.js](https://nodejs.org/en) is a powerful JavaScript runtime for building scalable, high-performance applications. While it's flexible and fast, managing different development and deployment environments can become complex. Docker offers a solution: a consistent, containerized environment that works everywhere.

> 
> **Acknowledgment**
>
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide. As a Docker Captain and experienced Front-end engineer, his expertise in Docker, DevOps, and modern web development has made this resource invaluable for the community, helping developers navigate and optimize their Docker workflows.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run a Node.js application using Docker.
- Run unit tests inside a Docker container.
- Set up a development container environment.
- Configure GitHub Actions for CI/CD with Docker.
- Deploy your Dockerized Node.js app to Kubernetes.

To begin, you’ll start by containerizing an existing Node.js application.

---

## Prerequisites

Before you begin, make sure you're familiar with the following:

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) or [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm) for managing dependencies and running scripts.
- Understanding of Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Once you've completed the Node.js getting started modules, you’ll be ready to containerize your own Node.js application using the examples and instructions provided in this guide.

