---
title: TanStack Start language-specific guide
linkTitle: TanStack Start
description: Containerize, develop, test, and deploy TanStack Start apps with Docker and Kubernetes
keywords: getting started, TanStack Start, tanstack, react, docker, language, Dockerfile, CI/CD, Kubernetes
summary: |
  This guide explains how to containerize TanStack Start applications, set up
  development and testing in containers, automate builds with GitHub Actions,
  and deploy to Kubernetes.
toc_min: 1
toc_max: 2
languages: [js]
tags: [frameworks]
params:
  time: 20 minutes
---

This guide shows you how to containerize a [TanStack Start](https://tanstack.com/start)
application using Docker, following best practices for production-ready
containers.

[TanStack Start](https://tanstack.com/start) is a full-stack React framework
built on [TanStack Router](https://tanstack.com/router) and Vite. It supports
server-side rendering, streaming, and type-safe routing. Docker provides a
consistent containerized environment from development to production.

> **Acknowledgment**
>
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide and maintaining the [docker-tanstack-start-sample](https://github.com/kristiyan-velkov/docker-tanstack-start-sample) repository used throughout this guide.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run a TanStack Start application using Docker.
- Set up a local development environment for TanStack Start inside a container.
- Run tests for your TanStack Start application within a Docker container.
- Configure a CI/CD pipeline using GitHub Actions for your containerized app.
- Deploy the containerized TanStack Start application to a local Kubernetes
  cluster for testing and debugging.

To begin, you'll start by containerizing an existing TanStack Start
application.

---

## Prerequisites

Before you begin, make sure you're familiar with the following:

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) or [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm) for managing dependencies and running scripts.
- Familiarity with [React](https://react.dev/) and [TanStack Start](https://tanstack.com/start) fundamentals.
- Understanding of Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Once you've completed the TanStack Start getting started modules, you'll be
ready to containerize your own application using the examples and instructions
in this guide.
