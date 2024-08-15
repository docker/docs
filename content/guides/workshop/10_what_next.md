---
title: What next after the Docker workshop
keywords: get started, setup, orientation, quickstart, intro, concepts, containers,
  docker desktop
description: Making sure you have more ideas of what you could do next with your application
aliases:
- /get-started/11_what_next/
---

Although you're done with the workshop, there's still a lot more to learn about containers.

Here are a few other areas to look at next.

## Container orchestration

Running containers in production is tough. You don't want to log into a machine and simply run a
`docker run` or `docker compose up`. Why not? Well, what happens if the containers die? How do you
scale across several machines? Container orchestration solves this problem. Tools like Kubernetes,
Swarm, Nomad, and ECS all help solve this problem, all in slightly different ways.

The general idea is that you have managers who receive the expected state. This state might be
"I want to run two instances of my web app and expose port 80." The managers then look at all of the
machines in the cluster and delegate work to worker nodes. The managers watch for changes (such as
a container quitting) and then work to make the actual state reflect the expected state.

## Cloud Native Computing Foundation projects

The CNCF is a vendor-neutral home for various open-source projects, including Kubernetes, Prometheus, 
Envoy, Linkerd, NATS, and more. You can view the [graduated and incubated projects here](https://www.cncf.io/projects/)
and the entire [CNCF Landscape here](https://landscape.cncf.io/). There are a lot of projects to help
solve problems around monitoring, logging, security, image registries, messaging, and more.

## Getting started video workshop

Docker recommends watching the video workshop from DockerCon 2022. Watch the entire video or use the following links to open the video at a particular section.

* [Docker overview and installation](https://youtu.be/gAGEar5HQoU)
* [Pull, run, and explore containers](https://youtu.be/gAGEar5HQoU?t=1400)
* [Build a container image](https://youtu.be/gAGEar5HQoU?t=3185)
* [Containerize an app](https://youtu.be/gAGEar5HQoU?t=4683)
* [Connect a DB and set up a bind mount](https://youtu.be/gAGEar5HQoU?t=6305)
* [Deploy a container to the cloud](https://youtu.be/gAGEar5HQoU?t=8280)

<iframe src="https://www.youtube-nocookie.com/embed/gAGEar5HQoU" style="max-width: 100%; aspect-ratio: 16 / 9;" width="560" height="auto" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Creating a container from scratch

If you'd like to see how containers are built from scratch, Liz Rice from Aqua Security has a fantastic talk in which she creates a container from scratch in Go. While the talk does not go into networking, using images for the filesystem, and other advanced topics, it gives a deep dive into how things are working.

<iframe src="https://www.youtube-nocookie.com/embed/8fi7uSYlOdc" style="max-width: 100%; aspect-ratio: 16 / 9;" width="560" height="auto" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Language-specific guides

If you are looking for information on how to containerize an application using your favorite language, see the [Language-specific guides](../../language/_index.md).
