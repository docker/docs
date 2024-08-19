---
title: Building images
keywords: build images, Dockerfile, layers, tag, push, cache, multi-stage
description: |
  Learn how to build Docker images from a Dockerfile. You'll understand the
  structure of a Dockerfile, how to build an image, and how to customize the
  build process.
summary: |
  Building container images is both technical and an art. You want to keep the
  image small and focused to increase your security posture, but also need to
  balance potential tradeoffs, such as caching impacts. In this series, you’ll
  deep dive into the secrets of images, how they are built and best practices.
layout: series
params:
  skill: Beginner
  time: 25 minutes
  prereq: None
---

## About this series

Learn how to build production-ready images that are lean and efficient Docker
images, essential for minimizing overhead and enhancing deployment in
production environments.

## What you'll learn

- Understanding image layers
- Writing a Dockerfile
- Build, tag and publish an image
- Using the build cache
- Multi-stage builds
