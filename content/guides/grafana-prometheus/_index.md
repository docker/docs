---
description: Containerize a Golang application and monitor it with Prometheus and Grafana.
keywords: golang, prometheus, grafana, monitoring, containerize
title: Monitor a Golang application with Prometheus and Grafana
summary: |
  Learn how to containerize a Golang application and monitor it with Prometheus and Grafana.
linkTitle: Monitor with Prometheus and Grafana
languages: [go]
params:
  time: 60 minutes
---

The guide teaches you how to containerize a Golang application and monitor it with Prometheus and Grafana. In this guide, you'll learn how to:

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) for his contribution to this guide.

## What will you learn?

* Containerize and run a Golang application using Docker
* Set up a local environment to develop a Golang application using containers
* How to use Docker Compose to run multiple services and connect them together to monitor a Golang application with Prometheus and Grafana.

## Prerequisites

- A good understanding of Golang is assumed.
- You must have familiarity with Docker concepts like containers, images, and Dockerfiles. If you are new to Docker, you can start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

Start by containerizing an existing Bun application.
