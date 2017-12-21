---
title: Interlock overview
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
---

Interlock is an application routing and load balancing system for Docker Swarm.  It uses
the Docker Remote API to automatically configure extensions such as Nginx or HAProxy for
application traffic.

## About

- [Introduction](intro/index.md)
  - [What is Interlock](intro/index.md)
  - [Architecture](intro/architecture.md)

## Deployment

- [Get started](install/index.md)
- [Deploy Interlock manually](install/manual-deployment.md)
- [Deploy Interlock offline](install/offline.md)
- [Deploy Interlock for production](install/production.md)

## Configuration

- [Interlock configuration](configuration/index.md)
- [Service labels](configuration/service-labels.md)

## Extensions

- [NGINX](extensions/nginx.md)
- [HAProxy](extensions/haproxy.md)

## Usage

- [Basic deployment](usage/index.md)
- [Applications with SSL](usage/ssl.md)
- [Application redirects](usage/redirects.md)
- [Persistent (sticky) sessions](usage/sessions.md)
- [Websockets](usage/websockets.md)
- [Canary application instances](usage/canary.md)
- [Service clusters](usage/service-clusters.md)
- [Context/path based routing](usage/context.md)
- [Host mode networking](usage/host-mode-networking.md)

## Operations

- [Updates](ops/index.md)
- [Tuning](ops/tuning.md)
