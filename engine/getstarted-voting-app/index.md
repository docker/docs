---
description: Getting started with multi-container Docker apps
keywords: docker, container, multi-container, services, swarm mode, cluster, stack deploy, compose, voting app
title: Get started with multi-container apps and services in swarm mode
---

This tutorial is built around a web-based voting application that collects, tallies, and returns the results of votes (for cats and dogs, or other choices you specify). The voting app includes several services. We'll deploy the app as a _stack_ to introduce some new concepts surfaced in [Compose v.3](/compose/compose-file.md), and also use [swarm mode](/engine/swarm/index.md), which is built into Docker Engine.

## Got Docker?

If you haven't yet downloaded Docker or installed it, go to [Get Docker](https://www.docker.com/) and grab Docker for your platform. Once you have Docker installed, you can run `docker hello-world` or other commands described in the newcomer tutorial to [verify your installation](/engine/getstarted/step_one.md#step-3-verify-your-installation). If you are totally new to Docker, you might try the quick [newcomer tutorial](/engine/getstarted/index.md) first, then come back.

## What you'll learn and do

You'll learn how to:

* Use `docker machine` to create multiple virtual local hosts or dockerized cloud servers
* Use `docker` commands to set up and run a swarm with manager and worker nodes
* Deploy the `vote` app by feeding our example `docker-stack.yml` file to `docker stack deploy`
* Test the app by voting for cats and dogs, and viewing the results
* Use the `visualizer` to explore and understand the runtime app and services
* Update the `docker-stack.yml` and re-deploy the app using a different `vote` image to implement a vote on different choices


## Anatomy of the voting app

The voting app you are about to deploy is composed of several services:


| Service        | Description | Base Image  |
| ------------- |--------------| -----|
| `vote`      | Presents the voting interface via port `5000`. Viewable at `<manager-IP>:5000` | Based on a Python image, `dockersamples/examplevotingapp_vote` |
| `result`      | Displays the voting results via port 5001.  Viewable at `<manager-IP>:5001`     |  Based on a Node.js image, `dockersamples/examplevotingapp_result` |
| `visualizer` | A web app that shows a map of the deployment of the various services across the available nodes via port `8080`. Viewable at `<manager-IP>:8080`  |  Based on a .NET image, `dockersamples/examplevotingapp_worker` |
| `redis` | Collects raw voting data and stores it in a key/value queue     |  Based on a `redis` image, `redis:alpine` |
| `db` | A PostgreSQL service which provides permanent storage on a host volume     |  Based on a `postgres` image, `postgres:9.4` |
| `worker` | A background service that transfers votes from the queue to permanent storage     |  Based on a .NET image, `dockersamples/examplevotingapp_worker` |

Each service will run in its own container. Using swarm mode, we can also scale the application to deploy replicas of containerized services distributed across multiple nodes.

For this tutorial, you need only have Docker running and the copy of `docker-stack.yml` we provide here.

## Next topic

TBD

## Next topic

TBD

## Where to go next

The voting app is also available as a [lab on
GitHub](https://github.com/docker/labs/blob/master/beginner/chapters/votingapp.md)
along with the complete [source
code]((https://github.com/docker/example-voting-app).

The lab is a deeper dive, and includes a few more tasks, like cloning a GitHub
repository, manually changing source code, and rebuilding an image instead of
using the ready-baked images referenced here.


&nbsp;
