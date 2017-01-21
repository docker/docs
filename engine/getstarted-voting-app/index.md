---
description: overview of voting app example
keywords: docker-stack.yml, stack deploy, compose, multi-container, services, swarm mode, cluster, voting app,
title: Tour the voting app
---

This example is built around a web-based voting application that collects,
tallies, and returns the results of votes (for cats and dogs, or other choices
you specify). The voting app includes several services, each one running in its
own container. We'll deploy the app as a _stack_ to introduce some new concepts
surfaced in [Compose Version 3](/compose/compose-file.md#version-3), and also
use [swarm mode](/engine/swarm/index.md), which is cluster management and
orchestration capability built into Docker Engine.

## Got Docker?

If you haven't yet downloaded Docker or installed it, go to [Get
Docker](/engine/getstarted/step_one.md#step-1-get-docker) and grab Docker for
your platform.  You can follow along and run this example using Docker for Mac,
Docker for Windows or Docker for Linux.

Once you have Docker installed, you can run `docker hello-world` or other
commands described in the Get Started with Docker tutorial to [verify your
installation](/engine/getstarted/step_one.md#step-3-verify-your-installation).
If you are totally new to Docker, you might continue through the full [Get
Started with Docker tutorial](/engine/getstarted/index.md) first, then come
back.

## What you'll learn and do

In this tutorial, you'll learn how to:

* Use `docker machine` to create multiple virtual local hosts or
dockerized cloud servers
* Use `docker` commands to set up and run a swarm with manager and worker nodes
* Deploy the `vote` app by feeding our example `docker-stack.yml` file to
the `docker stack deploy` command
* Test the app by voting for cats and dogs, and view the results
* Use the `visualizer` to explore and understand the runtime app and services
* Update the `docker-stack.yml` and redeploy the app using a different
`vote` image to implement a poll on different choices
* Use features new in Compose Version 3, highlighted in the sample app

## Anatomy of the voting app

The voting app you are about to deploy is composed of several services:


| Service        | Description | Base Image  |
| ------------- |--------------| -----|
| `vote`      | Presents the voting interface via port `5000`. Viewable at `<manager-IP>:5000` | Based on a Python image, `dockersamples/examplevotingapp_vote` |
| `result`      | Displays the voting results via port 5001.  Viewable at `<manager-IP>:5001`     |  Based on a Node.js image, `dockersamples/examplevotingapp_result` |
| `visualizer` | A web app that shows a map of the deployment of the various services across the available nodes via port `8080`. Viewable at `<manager-IP>:8080`  |  Based on a Node.js image, `dockersamples/visualizer` |
| `redis` | Collects raw voting data and stores it in a key/value queue     |  Based on a `redis` image, `redis:alpine` |
| `db` | A PostgreSQL service which provides permanent storage on a host volume     |  Based on a `postgres` image, `postgres:9.4` |
| `worker` | A background service that transfers votes from the queue to permanent storage     |  Based on a .NET image, `dockersamples/examplevotingapp_worker` |

Each service will run in its own container. Using swarm mode,
we can also scale the application to deploy replicas
of containerized services distributed across multiple nodes.

Here is an example of one of the services fully defined:

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

The `image` key defines which image the service will use. The `vote` service
uses `dockersamples/examplevotingapp_vote:before`.

The `depends_on` key allows you to specify that a service is only
deployed after another service. In our example, `vote` only deploys
after `redis`.

The `deploy` key specifies aspects of a swarm deployment, as described below in
[Compose Version 3 features and
compatibility](#compose-v3-features-and-compatibility).

## docker-stack.yml deployment configuration

We'll deploy the app using `docker-stack.yml`, which is a type of [Compose
file](/compose/compose-file.md) new in Compose Version 3.  

To follow along with the example, you need only have Docker running and the copy
of `docker-stack.yml` we provide here. This file defines all the services shown
in the [table above](#anatomy-of-the-voting-app), their base images,
configuration details such as ports and networks, application dependencies, and
the swarm configuration.

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

## Compose Version 3 features and compatibility

To deploy the voting application, we will run the `docker-stack deploy` command
with this `docker-stack.yml` file to pull the referenced images and launch the
services in a swarm as configured in the `.yml`.

Note that at the top of the `docker-stack.yml` file, the version is indicated as
`version: "3" `. The voting app example relies on Compose version 3, which is
designed to be cross-compatible with Compose and Docker Engine swarm mode.

Before we get started, let's take a look at some aspects of Compose files and
deployment options that are new in Compose Version 3, and that we want to highlight in
this walkthrough.

- [docker-stack.yml](#docker-stackyml)
  - [deploy key](#deploy-key)
- [docker stack deploy command](#docker-stack-deploy-command)
- [Docker stacks and services](#docker-stacks-and-services)

### docker-stack.yml

`docker-stack.yml` is a new type of [Compose file](/compose/compose-file.md)
only compatible with Compose Version 3.

#### deploy key

The `deploy` key allows you to specify various properties of a swarm deployment.

For example, the voting app configuration uses this to create replicas of the
`vote` and `result` services (2 containers of each will be deployed to the
swarm).

The voting app also uses the `deploy` key to constrain some services to run only
on the manager node.

For more about the `deploy key`, see [deploy](/compose/compose-file.md#deploy).

### docker stack deploy command

`docker stack deploy` is the command we will use to deploy with
`docker-stack.yml`.

* This command supports only `version: "3" ` Compose files.

* It does not support the `build` key supported in Compose files, which
builds based on a Dockerfile. You need to use pre-built images
with `docker stack deploy`.

* It can take the place of running `docker compose up` to deploy
Version 3 compatible applications.

See information about the [deploy](/compose/compose-file.md#deploy) key in the
Compose file reference and [`docker stack
deploy`](/engine/reference/commandline/stack_deploy.md) in the Docker Engine
command line reference.

### Docker stacks and services

Taken together, these new Compose features and deployment options can help when
mapping out distributed applications and clustering strategies. Rather than
thinking about running individual containers, we can start to model deployments
as Docker stacks and services.

### Compose file reference

For more on what's new in Compose Version 3:

* [Introducing Docker 1.13.0](https://blog.docker.com/2017/01/whats-new-in-docker-1-13/) blog post from [Docker Core Engineering](https://blog.docker.com/author/core_eng/)

* [Versioning](/compose/compose-file.md#versioning), [Version 3](/compose/compose-file.md#version-3), and  [Upgrading](/compose/compose-file.md#upgrading) in [Compose file reference](/compose/compose-file.md)

## What's next?

In the next step, we'll [set up two Dockerized hosts](node-setup.md).
