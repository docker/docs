---
description: Containerize Go apps using Docker
keywords: docker, getting started, go, golang, language, dockerfile
title: What will you learn in this module?
toc_min: 1
toc_max: 2
---

In this guide, you will learn how to create a containerized Go application using Docker.

Why [Go](https://golang.org/){:target="_blank" rel="noopener" class="_"}? Go is an open-source programming language that lets you build simple, reliable, and efficient software. Go is undeniably a major player in the modern Cloud ecosystem; both Docker and Kubernetes are written in Go. 

[golang]: https://golang.org/

> **Acknowledgment**
>
> We'd like to thank [Oliver Frolovs](https://twitter.com/nocturnalgopher){:target="_blank" rel="noopener" class="_"} for his contribution to the Golang get started guide.

In this guide, you’ll learn how to:

* Create a new `Dockerfile` which contains instructions required to build a Docker image for a simple Go program
* Run the newly built image as a container
* Set up a local development environment to connect a database to the container
* Use Docker Compose to run your Go application and other services it requires
* Configure a CI/CD pipeline for your application using [GitHub Actions](https://docs.github.com/en/actions){:target="_blank" rel="noopener" class="_"}

You can containerize your own Go application using the examples and resources provided after you complete the Go getting started modules.

Let's get started!

[Build your Go image](build-images.md){: .button .outline-btn}
