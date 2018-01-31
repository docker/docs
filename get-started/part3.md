---
title: "Get Started, Part 3: Docker Services"
keywords: services, replicas, scale, ports, compose, compose file, stack, networking
description: Learn how to define load-balanced and scalable service that runs containers.
---
{% include_relative nav.html selected="3" %}

## Prerequisites

Complete [Part 1](index.md) and [Part 2](part2.md).

## Introduction to part 3

In part 3, we scale our application to run mutliple containers (not just one)
and enable load-balancing. To do this, we must go one level up in the hierarchy
of a distributed application--the Docker service.

- Stack
- **Docker Service**   <span class="badge badge-danger">You are here</span>
- Container

> Docker Service vs Kubernetes Service
>
> We distinguish a [Docker serivce](https://docs.docker.com/engine/swarm/how-swarm-mode-works/services/){: target="_blank" class="_"}
> (with Swarm) from a [Kubernetes service](https://kubernetes.io/docs/concepts/services-networking/service/){: target="_blank" class="_"}.
> A Kubernetes service is a networking construct that load-balances pods (or
> sets of containers) behind a proxy. For now, you can think of a Docker service
> as one component of an application that is generated from one image.

## Run app as Docker service

Most applications are made up of multiple Docker services--a _stack_ of
services. For example, a video sharing site might include a service for the
frontend website, a service for a backend database, and a service for
transcoding uploaded user video. Each Docker service is an application
component--it is generated from one image, it is configurable, and it spawns one
or more containers.

We configure Docker services with Docker Compose YAML files to define how our
application containers should behave in production. For example, we can
define how many replicas of each container should run (called "scaling"), what
ports they should use, what computing resources they should have, and so on.

> Docker CLI vs Docker Compose
>
> In this tutorial, we use the Docker CLI (specifically, `docker stack`) to
> deploy a Docker Compose file. We are not using the tool, [Docker Compose](https://docs.docker.com/compose/overview/){: target="_blank" class="_"},
> with which you can also deploy applications.

### Create Docker compose file

In [part 2](part2.md), we deployed the `hellomoby` image with the Docker CLI and
used `docker run` commandline options (`-p` and `-d`) to apply simple
configuration. In this part, we deploy our application as if for production, as
a stack of services.

For now, our hellomoby application includes one service, a frontend web page. We
name the service, `web`, tell Docker to run three instances of `hellomoby`
(resulting in three containers), set a restart policy, and set resource limits.
Later in the tutorial, we add a service called `db` (to build a real stack) for
a backend redis database to store the number of web page visits.

1. Navigate to any directory, but preferably `~/hellomoby` just to stay organized.

2.  Create a file named `docker-compose-web.yml` with the following content:

    ```yaml
    version: "3"

    services:
      web:
        # replace "USERNAME" with your Docker Hub login name
        image: USERNAME/hellomoby:v1
        deploy:
          replicas: 3
          resources:
            limits:
              cpus: "0.1"
              memory: 50M
          restart_policy:
            condition: on-failure
        ports:
          - "80:80"
        networks:
          - webnet

    networks:
      webnet:
    ```

Our file, `docker-compose-web.yml`, tells Docker to do the following:

- Define a service called `web`.
- Pull the `hellomoby` image we uploaded to Docker Hub in [part 2](part2.md).
- Run three instances (containers) of the image.
- Limit each `web` container to use, at most, 10% of the CPU and 50MB of RAM.
- Restart any `web` container that fails.
- Map host port 80 to `web` port 80.
- Instruct `web` to use `webnet` to load-balance its containers. (Internally,
the containers themselves publish to port 80 at an ephemeral port.)
- Define a network we call `webnet` as a load-balanced overlay network (the default).

### Initiate swarm mode

Before we can deploy, we need to initiate [swarm mode](https://docs.docker.com/engine/swarm/){: target="_blank" class="_"}. Docker Swarm is the native Docker container orchestrator.

> Disable Kubernetes
>
> In [Part 4](part4-kube.md){: target="_blank" class="_"}, we explain how to
> deploy our app with Kubernetes. If you have Kubernetes enabled in Docker for
> Mac or Docker for Windows, disable it to continue. Another option is to
> prepend each command with `DOCKER_ORCHESTRATOR=swarm`.

1.  Initiate swarm mode to enable the use of Docker Compose to deploy applications
    as a stack of services:

    ```shell
    $ docker swarm init
    ```

2.  View the node and notice that the swarm MANAGER STATUS is "Leader". Later we add
    workers and build a real cluster. For now, we run a single node "cluster".

    ```shell
    $ docker node ls
    ID                             HOSTNAME       STATUS    AVAILABILITY    MANAGER STATUS
    yk08o9hvqww8k6bj5mflg2i6w *    <your host>    Ready     Active          Leader
    ```

### Deploy single-service stack

1.  Deploy the application as a stack and name it `hellomobylab`.

    ```shell
    $ docker stack deploy -c docker-compose-web.yml hellomobylab
    Creating network hellomobylab_webnet
    Creating service hellomobylab_web
    ```

    Our single service stack is running three container instances of our deployed
    image on one host. Let's investigate.

2.  List the stack:

    ```shell
    $ docker stack ls
    NAME                SERVICES
    hellomobylab        1
    ```

3.  List the one service in our application:

    ```shell
    $ docker service ls
    ID              NAME                MODE          REPLICAS    IMAGE                 PORTS
    5li63dwt7isq    hellomobylab_web    replicated    3/3         gordon/hellomoby:v1   *:80->80/tcp
    ```

4.  List the task IDs, with and without the quiet option:

    ```shell
    $ docker service ps hellomobylab_web
    $ docker service ps hellomobylab_web -q
    ```

5.  List the container IDs, with and without the quiet and no-truncate options.
    Notice the task IDs appended to the container NAMES.

    ```shell
    $ docker container ls
    $ docker container ls --no-trunc -q
    ```

6.  Pick any TASK_ID and use the `inspect` command with the [format option](https://docs.docker.com/engine/admin/formatting/){: target="_blank" class="_"} to find the associated CONTAINER ID:

    ```shell
    {% raw %}$ docker inspect --format='{{.Status.ContainerStatus.ContainerID}}' TASK_ID{% endraw %}
    ```

7.  Point a browser to [http://localhost](http://127.0.0.1/){: target="_blank" class="_"}
    and refresh three times. The CONTAINER ID changes, demonstrating Docker
    [service discovery and internal connection-based load-balancing](https://docs.docker.com/engine/swarm/networking/#configure-service-discovery).

    You can also run `curl` in a small for loop to see all three containers:

    ```shell
    $ for i in {1..3}; do curl -4 http://localhost; \n; done
    ```

> Slow response times?
>
> Depending on your networking configuration, it may take up to 30 seconds for
> containers to respond to HTTP requests. This probably stems from an unmet
> Redis dependency addressed later in the tutorial.

## Scale application

Scale the application by increasing the number of replicas to five, then rerun
`docker stack deploy`. Docker does an in-place update, with no need to tear down
the stack or kill any containers.

1.  Scale the number of container replicas by editing `docker-compose-web.yml`:

    ```yaml
    ...
        deploy:
          replicas: 4
    ...
    ```

2.  Redeploy:

    ```shell
    $ docker stack deploy -c docker-compose-web.yml hellomobylab
    Updating service hellomobylab_web (id: 5li63dwt7isqno5y7onxibapx)
    ```

3.  List all 4 running containers:

    ```shell
    $ docker container ls
    ```

4. Scale the number of container replicas at the commandline:

   ```shell
   $ docker service scale hellomobylab_web=5
   hellomobylab_web scaled to 5
   overall progress: 5 out of 5 tasks
   1/5: running   [==================================================>]
   2/5: running   [==================================================>]
   3/5: running   [==================================================>]
   4/5: running   [==================================================>]
   5/5: running   [==================================================>]
   verify: Service converged
   ```

5.  List all 5 running containers:

       ```shell
       $ docker container ls
       ```

## Tear down app and leave swarm

1.  Tear down the stack:

    ```shell
    $ docker stack rm hellomobylab
    ```

2.  Leave swarm mode:

    ```shell
    $ docker swarm leave --force
    ```

## Recap and cheat sheet

```shell
## Initiate swarm mode
docker swarm init

## Display swarm manager node
docker node ls

## Deploy application as stack named hellomobylab
docker stack deploy -c docker-compose-web.yml hellomobylab

## List stacks
docker stack ls

## List services
docker service ls

## List task IDs for service hellomobylab_web
docker service ps hellomobylab_web

## List container IDs
docker container ls

## Inspect task (and assoicated container)
{% raw %}docker inspect --format='{{.Status.ContainerStatus.ContainerID}}' TASK ID{% endraw %}

## Scale number of container replicas
docker service scale hellomobylab_web=5

## Tear down application stack
docker stack rm hellomobylab

## Leave swarm mode
docker swarm leave --force
```

## Conclusion of part 3

You have taken a huge step towards learning how to run containers in production.
Using `docker run` is simple, but the true implementation of a container in
production is running it as a Docker service. Changes to services can be applied
in place, as they run, using the same command that launched the service: `docker
stack deploy`.

Up next, we explain how to run our app on a real swarm cluster of Docker
machines.

[On to "Part 4" >>](part4.md){: class="button outline-btn"}
