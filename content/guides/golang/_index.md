---
title: Go language-specific guide
linkTitle: Go
description: Containerize Go apps using Docker
keywords: docker, getting started, go, golang, language, dockerfile
summary: |
  This guide teaches you how to containerize Go applications using Docker.
toc_min: 1
toc_max: 2
aliases:
  - /language/golang/
  - /guides/language/golang/
languages: [go]
params:
  time: 30 minutes
---

This guide will show you how to create, test, and deploy containerized Go applications using Docker.

> **Acknowledgment**
>
> Docker would like to thank [Oliver Frolovs](https://www.linkedin.com/in/ofr/) for his contribution to this guide.

## What will you learn?

In this guide, you’ll learn how to:

- Create a `Dockerfile` which contains the instructions for building a container image for a program written in Go.
- Run the image as a container in your local Docker instance and manage the container's lifecycle.
- Use multi-stage builds for building small images efficiently while keeping your Dockerfiles easy to read and maintain.
- Use Docker Compose to orchestrate running of multiple related containers together in a development environment.
- Configure a CI/CD pipeline for your application using [GitHub Actions](https://docs.github.com/en/actions)
- Deploy your containerized Go application.

## Prerequisites

Some basic understanding of Go and its toolchain is assumed. This isn't a Go tutorial. If you are new to the : languages:,
the [Go website](https://golang.org/) is a great place to explore,
so _go_ (pun intended) check it out!

You also must know some basic [Docker concepts](/get-started/docker-concepts/the-basics/what-is-a-container.md) as well as to
be at least vaguely familiar with the [Dockerfile format](/manuals/build/concepts/dockerfile.md).

Your Docker set-up must have BuildKit enabled. BuildKit is enabled by default for all users on [Docker Desktop](/manuals/desktop/_index.md).
If you have installed Docker Desktop, you don’t have to manually enable BuildKit. If you are running Docker on Linux,
please check out BuildKit [getting started](/manuals/build/buildkit/_index.md#getting-started) page.

Some familiarity with the command line is also expected.

## What's next?

The aim of this guide is to provide enough examples and instructions for you to containerize your own Go application and deploy it into the Cloud.

Start by building your first Go image.
