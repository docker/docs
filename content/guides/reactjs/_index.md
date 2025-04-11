---
title: React.js language-specific guide
linkTitle: React.js
description: Containerize and develop React.js apps using Docker
keywords: getting started, React.js, react.js, docker, language, Dockerfile
summary: |
  This guide explains how to containerize React.js applications using Docker.
toc_min: 1
toc_max: 2
aliases:
  - /language/reactjs/
  - /guides/language/reactjs/
languages: [js]
params:
  time: 20 minutes
---

The React.js language-specific guide shows you how to containerize a React.js application using Docker, following best practices for creating efficient, production-ready containers.

[React.js](https://react.dev/) is a widely used library for building interactive user interfaces. However, managing dependencies, environments, and deployments efficiently can be complex. **Docker** simplifies this process by providing a consistent and containerized environment.

> 
> **Acknowledgment**
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide. As a Docker Captain and experienced Front-end engineer, his expertise in Docker, DevOps, and modern web development has made this resource invaluable for the community, helping developers navigate and optimize their Docker workflows.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run a React.js application using Docker.
- Set up a local development environment for React.js inside a container. 
- Run tests for your React.js application within a Docker container.
- Configure a CI/CD pipeline using GitHub Actions for your containerized app.
- Deploy the containerized React.js application to a local Kubernetes cluster for testing and debugging.

To begin, you’ll start by containerizing an existing React.js application.

---

## Prerequisites

Before you begin, make sure you're familiar with the following:

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) or [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm) for managing dependencies and running scripts.
- Familiarity with [React.js](https://react.dev/) fundamentals.
- Understanding of Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Once you've completed the React.js getting started modules, you’ll be ready to containerize your own React.js application using the examples and instructions provided in this guide.
