---
title: "What next"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop
description: Making sure you have more ideas of what you could do next with your application
---

Although we're done with our workshop, there's still a LOT more to learn about containers!
We're not going to go deep-dive here, but here are a few other areas to look at next!

## Container orchestration

Running containers in production is tough. You don't want to log into a machine and simply run a
`docker run` or `docker-compose up`. Why not? Well, what happens if the containers die? How do you
scale across several machines? Container orchestration solves this problem. Tools like Kubernetes,
Swarm, Nomad, and ECS all help solve this problem, all in slightly different ways.

The general idea is that you have "managers" who receive **expected state**. This state might be
"I want to run two instances of my web app and expose port 80." The managers then look at all of the
machines in the cluster and delegate work to "worker" nodes. The managers watch for changes (such as
a container quitting) and then work to make **actual state** reflect the expected state.

## Cloud Native Computing Foundation projects

The CNCF is a vendor-neutral home for various open-source projects, including Kubernetes, Prometheus, 
Envoy, Linkerd, NATS, and more! You can view the [graduated and incubated projects here](https://www.cncf.io/projects/){:target="_blank" rel="noopener" class="_"}
and the entire [CNCF Landscape here](https://landscape.cncf.io/){:target="_blank" rel="noopener" class="_"}. There are a LOT of projects to help
solve problems around monitoring, logging, security, image registries, messaging, and more!

So, if you're new to the container landscape and cloud-native application development, welcome! Please
connect with the community, ask questions, and keep learning! We're excited to have you!
