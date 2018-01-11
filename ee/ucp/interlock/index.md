---
title: Layer 7 routing overview
description: Learn about Layer 7 routing, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, layer 7, routing, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
next_steps:
- path: ops/
  title: Updates
- path: ops/tuning/
  title: Tuning
redirect_from:
  - /datacenter/ucp/3.0/guides/interlock/
---
{% if include.version=="ucp-3.0" %}

Layer 7 routing is an application routing and load balancing system for Docker Swarm. It uses
the Docker Remote API to automatically configure extensions such as Nginx or HAProxy for
application traffic.

## About

- [Introduction](intro/index.md)
  - [What is Layer 7 routing](intro/index.md)
  - [Architecture](intro/architecture.md)

## Deployment

- [Get started](install/index.md)
- [Deploy Layer 7 routing manually](install/manual-deployment.md)
- [Deploy Layer 7 routing offline](install/offline.md)
- [Deploy Layer 7 routing for production](install/production.md)

## Configuration

- [Layer 7 routing configuration](configuration/index.md)
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

{% endif %}
