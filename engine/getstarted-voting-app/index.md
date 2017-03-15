---
description: overview of voting app example
keywords: docker-stack.yml, stack deploy, compose, multi-container, services, swarm mode, cluster, voting app,
title: Sample app overview
---

This example is built around a web-based voting application
that collects, tallies, and returns the results of votes
(for cats and dogs, or other choices you specify). The voting
app includes several services, each one running in its
own container. We'll deploy the app as a _stack_ to introduce
some new concepts surfaced in
[Compose Version 3](/compose/compose-file.md#version-3), and
also use [swarm mode](/engine/swarm/index.md), which is
cluster management and orchestration capability built into
Docker Engine.

## Got Docker?

If you haven't yet downloaded Docker or installed it, go to [Get
Docker](/engine/getstarted/step_one.md#step-1-get-docker) and grab Docker for
your platform.  You can follow along and run this example using [Docker for Mac](/docker-for-mac/index.md),
[Docker for Windows](/docker-for-windows/index.md), or [Docker for Linux](/installation/linux/index.md).

If you are totally new to Docker, you might want to work through the [Get
Started with Docker tutorial](/engine/getstarted/index.md) first, then come
back here.

## What you'll learn and do

In this tutorial, you'll learn how to:

* Use `docker machine` to create multiple virtual local hosts or
dockerized cloud servers
* Use `docker` commands to set up and run a swarm with manager and worker nodes
* Deploy the `vote` app services across the two nodes by feeding our example `docker-stack.yml` file to
the `docker stack deploy` command
* Test the app by voting for cats and dogs, and view the results
* Use the `visualizer` to explore and understand the runtime app and services
* Update the `docker-stack.yml` and redeploy the app using a different
`vote` image to implement a poll on different choices
* Use features new in Compose Version 3, highlighted in the sample app

## Preview of voting app stack and how it works

These next few topics provide a quick tour of the services, deployment configuration, files, and commands we will use.

This diagram represents the application stack at runtime. It shows
dependencies among the services, and a potential division of services between
the manager and worker nodes in a swarm. As you'll discover in the tutorial,
some services are constrained to always run on a manager node, while others can
run on either a manager or workers, at the discretion of swarm load balancing.

![voting app diagram](images/vote-app-diagram.png)

### Services and images overview

A [service](/engine/reference/glossary.md#service) is a bit of executable code
designed to accomplish a specific task. A service can run in one or more
containers. Defining a service configuration for your app (above and beyond
`docker run` commands in a Dockerfile) enables you to deploy the app to a swarm
and manage it as a distributed, multi-container application.

The voting app you are about to deploy is made up of several services, each
based on an [image](/engine/reference/glossary.md#image) that the app will pull
from [Docker Hub](/engine/reference/glossary.md#docker-hub) at runtime:

| Service        | Description | Base Image on Docker Hub |
| ------------- |--------------| -----|
| `vote`      | Displays the web page where you cast your vote at `<manager-IP>:5000` | Based on a Python image, [dockersamples/examplevotingapp_vote](https://hub.docker.com/r/dockersamples/examplevotingapp_vote/) |
| `result`      | Shows the voting results in a web browser at `<manager-IP>:5001`     |  Based on a Node.js image, [dockersamples/examplevotingapp_result](https://hub.docker.com/r/dockersamples/examplevotingapp_result/) |
| `visualizer` | Shows a realtime map of services deployed across the available nodes, viewable at `<manager-IP>:8080`  |  Based on a Node.js image, [dockersamples/visualizer](https://hub.docker.com/r/dockersamples/visualizer/) |
| `redis` | Collects raw voting data and stores it in a key/value queue     |  Based on the Alpine version of the official `redis` image, [redis:alpine](https://hub.docker.com/_/redis/) |
| `db` | A PostgreSQL service which provides permanent storage on a host volume    |  Based on the official `postgres` image, [postgres:9.4](https://hub.docker.com/_/postgres/) |
| `worker` | A background service that transfers votes from the queue to permanent storage     |  Based on a .NET image, [dockersamples/examplevotingapp_worker](https://hub.docker.com/r/dockersamples/examplevotingapp_worker/) |

Each service will run in its own [container](/engine/reference/glossary.md#container). Using swarm mode,
we can also scale the application to deploy replicas
of containerized services distributed across multiple nodes.

### docker-stack.yml deployment configuration file

In the Getting Started with Docker tutorial, you wrote a
[Dockerfile for the whalesay app](/engine/getstarted/step_four.md) then used
it to build a single image and run it as a single container.

For this tutorial, the images are pre-built, and we use a _stack file_ instead
of a Dockerfile to specify the images. When we deploy, each image will run as a
service in a container (or in multiple containers, for those that have replicas
defined to scale the app).

To follow along, you need only have Docker running and a copy of the
`docker-stack.yml`
file that we provide.

This file defines all the services we want to use along with details about how
and where those services will run; their base images, configuration
details such as ports, networks, volumes, application dependencies, and the
swarm configuration.

This **example snip-it** taken from our `docker-stack.yml` shows one of the
services fully defined. (The full file is
[**here**](https://github.com/docker/example-voting-app/blob/master/docker-stack.yml).)

```
vote:
  image: dockersamples/examplevotingapp_vote:before
  ports:
    - 5000:80
  networks:
    - frontend
  depends_on:
    - redis
  deploy:
    replicas: 2
    update_config:
      parallelism: 2
    restart_policy:
      condition: on-failure
```

* The **image** key defines which image the service will use. The `vote` service
uses `dockersamples/examplevotingapp_vote:before`. This specifies the path to
the image on Docker Hub (as shown in the table above), and an [image
tag](/engine/reference/commandline/tag.md), `before` to indicate the version of
the image we want to start with. In the second part of the tutorial, we will
edit this file to call a different verson of this image with an `after` tag.

* The **depends_on** key allows you to specify that a service is only
deployed after another service. In our example, `vote` only deploys
after `redis`.

* The **deploy** key specifies aspects of a swarm deployment. For example,
in this configuration we create _replicas_ of the `vote` service (2 containers
for `vote` will be deployed to the swarm). The `result` service, not shown in
the file snip-it above, will also have 2 replicas. Additionally, we will use the
`deploy` key to constrain some other services (`db` and `visualizer`) to run
only on a manager node.


### docker stack deploy command

To deploy the voting app, we will run the [`docker stack
deploy`](/engine/reference/commandline/stack_deploy.md) command with appropriate
options, using the configuration in our `docker-stack.yml` file to pull the
referenced images and launch the services in a swarm.

### Where to learn more

If you are interested in reading more about Compose version 3.x, stack files,
Docker Engine 1.13.x, swarm mode integration, Docker CE, or Docker EE, jump
to the [list of resources](customize-app.md#resources) at the end of this
tutorial.

## What's next?

Ready to get started? In the next step, we'll [set up two Dockerized
hosts](node-setup.md).
