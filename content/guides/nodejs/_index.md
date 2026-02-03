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
tags: []
params:
  time: 20 minutes
---

[Node.js](https://nodejs.org/en) is a JavaScript runtime for building web applications. This guide shows you how to containerize a TypeScript Node.js application with a React frontend and PostgreSQL database.

The sample application is a modern full-stack Todo application featuring:

- **Backend**: Express.js with TypeScript, PostgreSQL database, and RESTful API
- **Frontend**: React.js with Vite and Tailwind CSS 4


> **Acknowledgment**
>
> Docker extends its sincere gratitude to [Kristiyan Velkov](https://www.linkedin.com/in/kristiyan-velkov-763130b3/) for authoring this guide. As a Docker Captain and experienced Full-stack engineer, his expertise in Docker, DevOps, and modern web development has made this resource invaluable for the community, helping developers navigate and optimize their Docker workflows.

---

## What will you learn?

In this guide, you will learn how to:

- Containerize and run a Node.js application using Docker.
- Run tests inside a Docker container.
- Set up a development container environment.
- Configure GitHub Actions for CI/CD with Docker.
- Deploy your Dockerized Node.js app to Kubernetes.

To begin, you’ll start by containerizing an existing Node.js application.

---

## Prerequisites

Before you begin, make sure you're familiar with the following:

- Basic understanding of [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript) and [TypeScript](https://www.typescriptlang.org/).
- Basic knowledge of [Node.js](https://nodejs.org/en), [npm](https://docs.npmjs.com/about-npm), and [React](https://react.dev/) for modern web development.
- Understanding of Docker concepts such as images, containers, and Dockerfiles. If you're new to Docker, start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.
- Familiarity with [Express.js](https://expressjs.com/) for backend API development.

Once you've completed the Node.js getting started modules, you’ll be ready to containerize your own Node.js application using the examples and instructions provided in this guide.
