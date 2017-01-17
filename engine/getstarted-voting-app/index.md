---
description: overview of voting app example
keywords: docker-stack.yml, stack deploy, compose, container, multi-container, services, swarm mode, cluster, voting app,
title: Learn about the voting app example and setup
---

This example is built around a web-based voting application that collects,
tallies, and returns the results of votes (for cats and dogs, or other choices
you specify). The voting app includes several services, each one running in its
own container. We'll deploy the app as a _stack_ to introduce some new concepts
surfaced in [Compose v.3](/compose/compose-file.md), and also use [swarm
mode](/engine/swarm/index.md), which is cluster management and orchestration
capability built into Docker Engine.

## Got Docker?

If you haven't yet downloaded Docker or installed it, go to [Get
Docker](/engine/getstarted/step_one.md#step-1-get-docker) and grab Docker for
your platform.  You can follow along and run this example using Docker for Mac,
Docker for Windows or Docker for Linux.

Once you have Docker installed, you can run `docker hello-world`
or other commands described in the newcomer tutorial to [verify your
installation]
(/engine/getstarted/step_one.md#step-3-verify-your-installation).
If you are totally new to Docker, you might continue through the full [newcomer
tutorial](/engine/getstarted/index.md) first, then come back.

## What you'll learn and do

In this tutorial, you'll learn how to:

* Use `docker machine` to create multiple virtual local hosts or
dockerized cloud servers
* Use `docker` commands to set up and run a swarm with manager and worker nodes
* Deploy the `vote` app by feeding our example `docker-stack.yml` file to
the `docker stack deploy` command
* Test the app by voting for cats and dogs, and viewing the results
* Use the `visualizer` to explore and understand the runtime app and services
* Update the `docker-stack.yml` and re-deploy the app using a different
`vote` image to implement a vote on different choices
* Identify the differences between `docker-stack.yml` and
`docker-compose.yml` files, which serve similar purposes and
can be used somewhat interchangeably
* Identify the differences between `docker compose up` and
`docker stack deploy` commands, which serve similar purposes and
can be used somewhat interchangeably

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

Each service will run in its own container. Using swarm mode, we can also scale
the application to deploy replicas of containerized services distributed across
multiple nodes.

## docker-stack.yml

We'll deploy the app using `docker-stack.yml`.  To follow along with the
example, you need only have Docker running and the copy of `docker-stack.yml` we
provide here. This file defines all the services shown in the [table above](#anatomy-of-the-voting-app), their base images,
configuration details such as ports and dependencies, and the swarm
configuration.

```
version: "3"
services:

  redis:
    image: redis:alpine
    ports:
      - "6379"
    networks:
      - frontend
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure
  db:
    image: postgres:9.4
    volumes:
      - db-data:/var/lib/postgresql/data
    networks:
      - backend
    deploy:
      placement:
        constraints: [node.role == manager]
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
  result:
    image: dockersamples/examplevotingapp_result:before
    ports:
      - 5001:80
    networks:
      - backend
    depends_on:
      - db
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure

  worker:
    image: dockersamples/examplevotingapp_worker
    networks:
      - frontend
      - backend
    deploy:
      mode: replicated
      replicas: 1
      labels: [APP=VOTING]
      restart_policy:
        condition: on-failure
        delay: 10s
        max_attempts: 3
        window: 120s

  visualizer:
    image: dockersamples/visualizer:stable
    ports:
      - "8080:8080"
    stop_grace_period: 1m30s
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
    deploy:
      placement:
        constraints: [node.role == manager]

networks:
  frontend:
  backend:

volumes:
  db-data:
```

To deploy the application, we will use the `docker-stack deploy` command with this `docker-stack.yml` file to pull the referenced images and launch the services in a swarm as configured in the `.yml`.

But first, we need to set up the hosts and create a swarm.

## What's next?

In the next step, we'll [set up two Dockerized hosts](node-setup.md).
