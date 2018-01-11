---
title: What is Layer 7 routing
description: Learn about Layer 7 routing, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

Layer 7 routing is an application routing proxy service for Docker.

## Design Goals

- Fully integrate with Docker (Swarm, Services, Secrets, Configs)
- Enhanced configuration (context roots, TLS, zero downtime deploy, rollback)
- Support external load balancers (nginx, haproxy, F5, etc) via extensions
- Least privilege for extensions (no Docker API access)

Layer 7 routing was designed to be a first class application routing layer for Docker.
The following are the high level features it provides:

## Automatic Configuration
Layer 7 routing uses the Docker API for configuration. The user does not have to manually
update or restart anything to make services available.

## Native Swarm Support
Layer 7 routing is fully Docker native.  It runs on Docker Swarm and routes traffic using
cluster networking and Docker services.

## High Availability
Layer 7 routing runs as Docker services which are highly available and handle failures gracefully.

## Scalability
Layer 7 routing uses a modular design where the proxy service is separate.  This allows an
operator to individually customize and scale the proxy layer to whatever demand.  This is
transparent to the user and causes no downtime.

## SSL
Layer 7 routing leverages Docker Secrets to securely store and use SSL certificates for services.  Both
SSL termination and TCP passthrough are supported.

## Context Based Routing
Layer 7 routing supports advanced application request routing by context or path.

## Host Mode Networking
Layer 7 routing supports running the proxy and application services in "host" mode networking allowing
the operator to bypass the routing mesh completely.  This is beneficial if you want
maximum performance for your applications.

## Blue-Green and Canary Service Deployment
Layer 7 routing supports blue-green service deployment allowing an operator to deploy a new application
while the current version is serving.  Once traffic is verified to the new application the operator
can scale the older version to zero.  If there is a problem the operation is quickly reversible.

## Service Cluster Support
Layer 7 routing supports multiple extension+proxy combinations allowing for operators to partition load
balancing resources for uses such as region or organization based load balancing.

## Least Privilege
Layer 7 routing supports (and recommends) being deployed where the load balancing
proxies do not need to be colocated with a Swarm manager.  This makes the
deployment more secure by not exposing the Docker API access to the extension or proxy services.

{% endif %}
