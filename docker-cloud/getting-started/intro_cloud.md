---
description: Introducing Docker Cloud concepts and terminology
keywords: node, create, understand
redirect_from:
- /docker-cloud/getting-started/beginner/intro_cloud/
title: Introducing Docker Cloud
---

This page introduces core Docker Cloud concepts and features so you can easily follow along with the tutorial.

The tutorial goes through the following steps:

1. Set up your hosts by linking to a cloud service provider or your own Linux hosts.
2. Deploy your first node cluster.
3. Deploy your first service.

Know all this stuff already? Skip to [Link to your infrastructure](connect-infra.md).

## What is a node?
A node is an individual Linux host used to deploy and run your applications. Docker Cloud does not provide hosting services, so all of your applications, services, and containers run on your own hosts. Your hosts can come from several different sources, including physical servers, virtual machines or cloud providers.

## What is a node cluster?
When launching a node from a cloud provider you actually create a node cluster. Node Clusters are groups of nodes of the same type and from the same cloud provider. Node clusters allow you to scale the infrastructure by provisioning more nodes with a drag of a slider.

### Use cloud service providers
Docker Cloud makes it easy to provision nodes from existing cloud providers. If you already have an account with an infrastructure as a service provider, you can provision new nodes directly from within Docker Cloud. Today we have native support for Amazon Web Services, DigitalOcean, Microsoft Azure, Packet.net, and IBM SoftLayer.

### Use your own hosts ("Bring your own nodes")
You can also provide your own node or nodes. This means you can use any Linux host connected to the Internet as a Docker Cloud node as long as you can install a Cloud agent. The agent registers itself with your Docker account, and allows you to use Docker Cloud to deploy containerized applications.

## What is a service?
Services are logical groups of containers from the same image. Services make it simple to scale your application across different nodes. In Docker Cloud you drag a slider to increase or decrease the availability, performance, and redundancy of the application. Services can also be linked one to another even if they are deployed on different nodes, regions, or even cloud providers.

## Let's get started!
Log in to <a href="https://cloud.docker.com" target="_blank">Docker Cloud</a> using your Docker ID. (These are the same credentials you used for Docker Hub if you had an account there.)

Start here [by linking your infrastructure to Docker Cloud](connect-infra.md).
