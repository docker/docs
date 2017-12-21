---
title: Interlock architecture
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

The following are definitions that are used:

- Cluster: A group of compute resources running Docker
- Swarm: A Docker cluster running in Swarm mode
- Upstream: An upstream container that serves an application
- Proxy Service: A service that provides load balancing and proxying (such as Nginx)
- Extension Service: A helper service that configures the proxy service
- Service Cluster: A service cluster is an Interlock extension+proxy service
- GRPC: A high-performance RPC framework

## Services
Interlock runs entirely as Docker Swarm services.  There are three core services
in an Interlock routing layer: core, extension and proxy.

## Core
The core service is responsible for interacting with the Docker Remote API and building
an upstream configuration for the extensions.  This is served on a GRPC API that the
extensions are configured to access.

## Extension
The extension service is a helper service that queries the Interlock GRPC API for the
upstream configuration.  The extension service uses this to configure
the proxy service.  For proxy services that use files such as Nginx or HAProxy the
extension service generates the file and sends it to Interlock using the GRPC API.  Interlock
then updates the corresponding Docker Config object for the proxy service.

## Proxy
The proxy service handles the actual requests for the upstream application services.  These
are configured using the data created by the corresponding extension service.

Interlock manages both the extension and proxy service updates for both configuration changes
and application service deployments.  There is no intervention from the operator required.

{% endif %}
