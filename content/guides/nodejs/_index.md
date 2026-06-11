---
title: Node.js language-specific guide
linkTitle: Node.js
description: Containerize and develop Node.js applications using Docker
keywords: getting started, node, node.js
summary: |
  This guide explains how to containerize Node.js applications using Docker.
toc_min: 1
toc_max: 2
aliases:
  - /language/nodejs/
  - /guides/language/nodejs/
languages: [js]
tags: [dhi]
params:
  time: 20 minutes
---

[Node.js](https://nodejs.org/en) is a JavaScript runtime for building server-side applications. This guide shows you how to containerize a TypeScript Node.js application using Docker, starting from a simple Express API and progressively adding features like a database and CI/CD.

This guide focuses on a backend Node.js API. If you're building a standalone frontend application, Docker has dedicated guides for [React.js](/guides/reactjs/), [Vue.js](/guides/vuejs/), [Angular](/guides/angular/), and [Next.js](/guides/nextjs/).

> **Acknowledgment**
>
> Docker thanks [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for his contribution to this guide.

## What will you learn?

In this guide, you'll learn how to:

- Containerize and run a Node.js application using Docker.
- Set up a local development environment using containers.
- Run tests inside a Docker container.
- Configure GitHub Actions for CI/CD with Docker.
- Inspect and generate supply chain attestations for your image.
- Deploy your containerized Node.js application to Kubernetes.

Start by containerizing a Node.js application.

## Prerequisites

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) and [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en) and [npm](https://docs.npmjs.com/about-npm).
- Familiarity with Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.
