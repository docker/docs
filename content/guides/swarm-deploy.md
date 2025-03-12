---
title: Deploy to Swarm
keywords: swarm, swarm services, stacks
description: Learn how to describe and deploy a simple application on Docker Swarm.
aliases:
  - /get-started/part4/
  - /get-started/swarm-deploy/
  - /guides/deployment-orchestration/swarm-deploy/
summary: |
  Discover how to deploy and manage Docker containers using Docker Swarm.
tags: [deploy]
params:
  time: 10 minutes
---

{{% include "swarm-mode.md" %}}

## Prerequisites

- Download and install Docker Desktop as described in [Get Docker](/get-started/get-docker.md).
- Work through containerizing an application in [Docker workshop part 2](/get-started/workshop/02_our_app.md)
- Make sure that Swarm is enabled on your Docker Desktop by typing `docker system info`, and looking for a message `Swarm: active` (you might have to scroll up a little).

  If Swarm isn't running, simply type `docker swarm init` in a shell prompt to set it up.

## Introduction

Now that you've demonstrated that the individual components of your application run as stand-alone containers and shown how to deploy it using Kubernetes, you can look at how to arrange for them to be managed by Docker Swarm. Swarm provides many tools for scaling, networking, securing and maintaining your containerized applications, above and beyond the abilities of containers themselves.

In order to validate that your containerized application works well on Swarm, you'll use Docker Desktop's built in Swarm environment right on your development machine to deploy your application, before handing it off to run on a full Swarm cluster in production. The Swarm environment created by Docker Desktop is fully featured, meaning it has all the Swarm features your app will enjoy on a real cluster, accessible from the convenience of your development machine.

## Describe apps using stack files

Swarm never creates individual containers like you did in the previous step of this tutorial. Instead, all Swarm workloads are scheduled as services, which are scalable groups of containers with added networking features maintained automatically by Swarm. Furthermore, all Swarm objects can and should be described in manifests called stack files. These YAML files describe all the components and configurations of your Swarm app, and can be used to create and destroy your app in any Swarm environment.

Now you can write a simple stack file to run and manage your Todo app, the container `getting-started` image created in [Part 2](02_our_app.md) of the tutorial. Place the following in a file called `bb-stack.yaml`:

{{% include "swarm-compose-compat.md" %}}

```yaml
version: "3.7"

services:
  bb-app:
    image: getting-started
    ports:
      - "8000:3000"
```

In this Swarm YAML file, there is one object, a `service`, describing a scalable group of identical containers. In this case, you'll get just one container (the default), and that container will be based on your `getting-started` image created in [Part 2](02_our_app.md) of the tutorial. In addition, you've asked Swarm to forward all traffic arriving at port 8000 on your development machine to port 3000 inside our getting-started container.

> **Kubernetes Services and Swarm Services are very different**
>
> Despite the similar name, the two orchestrators mean very different things by
> the term 'service'. In Swarm, a service provides both scheduling and
> networking facilities, creating containers and providing tools for routing
> traffic to them. In Kubernetes, scheduling and networking are handled
> separately, deployments (or other controllers) handle the scheduling of
> containers as pods, while services are responsible only for adding
> networking features to those pods.

## Deploy and check your application

1. Deploy your application to Swarm:

   ```console
   $ docker stack deploy -c bb-stack.yaml demo
   ```

   If all goes well, Swarm will report creating all your stack objects with no complaints:

   ```shell
   Creating network demo_default
   Creating service demo_bb-app
   ```

   Notice that in addition to your service, Swarm also creates a Docker network by default to isolate the containers deployed as part of your stack.

2. Make sure everything worked by listing your service:

   ```console
   $ docker service ls
   ```

   If all has gone well, your service will report with 1/1 of its replicas created:

   ```shell
   ID                  NAME                MODE                REPLICAS            IMAGE               PORTS
   il7elwunymbs        demo_bb-app         replicated          1/1                 getting-started:latest   *:8000->3000/tcp
   ```

   This indicates 1/1 containers you asked for as part of your services are up and running. Also, you see that port 8000 on your development machine is getting forwarded to port 3000 in your getting-started container.

3. Open a browser and visit your Todo app at `localhost:8000`; you should see your Todo application, the same as when you ran it as a stand-alone container in [Part 2](02_our_app.md) of the tutorial.

4. Once satisfied, tear down your application:

   ```console
   $ docker stack rm demo
   ```

## Conclusion

At this point, you've successfully used Docker Desktop to deploy your application to a fully-featured Swarm environment on your development machine. You can now add other components to your app and taking advantage of all the features and power of Swarm, right on your own machine.

In addition to deploying to Swarm, you've also described your application as a stack file. This simple text file contains everything you need to create your application in a running state; you can check it in to version control and share it with your colleagues, letting you to distribute your applications to other clusters (like the testing and production clusters that probably come after your development environments).

## Swarm and CLI references

Further documentation for all new Swarm objects and CLI commands used in this article are available here:

- [Swarm Mode](/manuals/engine/swarm/_index.md)
- [Swarm Mode Services](/manuals/engine/swarm/how-swarm-mode-works/services.md)
- [Swarm Stacks](/manuals/engine/swarm/stack-deploy.md)
- [`docker stack *`](/reference/cli/docker/stack/)
- [`docker service *`](/reference/cli/docker/service/)
