---
title: Configure the HTTP Routing Mesh
description: Learn how to configure UCP's HTTP Routing Mesh
keywords: ucp, services, http, dns
---

UCP provides an HTTP routing mesh, that extends the networking capabilities
of Docker Engine. Docker Engine provides load balancing and service discovery
at the transport layer for TCP and UDP connections. UCP's HTTP routing mesh
allows you to extend service discovery to have name-based virtual hosting for
HTTP and HTTPS services.

See the
[Docker Engine documentation on overlay networks](/engine/swarm/networking.md)
for more information on what Docker Engine provides.

## Enable the HTTP routing mesh

To enable the HTTP routing mesh, go to the **UCP web UI**, navigate to the
**Settings** page, and click the **Routing Mesh** tab.

<!-- todo: add screenshot -->

The default port for HTTP services is **80**, and the default port for HTTPS
services is **8443**. You may choose an alternate port on this screen.

Check the checkbox to enable the HTTP routing mesh. This will create a service
called `ucp-hrm` and a network called `ucp-hrm`.

If the HTTP routing mesh receives a HTTP request for a domain that it does not
handle, it returns a 503 error (Bad Gateway). For HTTPS requests, all unknown
domains are routed to the UCP web interface.

## HTTPS support

The HTTP routing mesh has support for routing using HTTPS. Using a feature of
HTTPS called Server Name Indication, the HTTP routing mesh is able to route
connections to service backends without terminating the HTTPS connection.

To use HTTPS support, no certificates for the service are provided to the HTTP
routing mesh. Instead, the backend service **must** handle HTTPS connections
directly. Services that meet this criteria can use the `SNI` protocol to
indicate handling of HTTPS in this manner.

## Route to a service

The HTTP routing mesh can route to a Docker service that runs a webserver.
This service must meet three criteria:

* The service must be connected a network with a `com.docker.ucp.mesh.http` label
* The service must publish one or more ports
* The service must have one or more labels prefixed with
  `com.docker.ucp.mesh.http` to specify the ports to route (see the syntax
  below)

These options can be configured using the UCP UI, or can be entered manually
using the `docker service` command.

## Route domains to the HTTP routing mesh

The HTTP routing mesh uses the `Host` HTTP header (or the Server Name
Indication field for HTTPS requests) to determine which service should receive
a particular HTTP request. This is typically done using DNS and pointing one or
more domains to one or more nodes in the UCP cluster.

## Networks, Access Control, and the HTTP routing mesh

The HTTP routing mesh uses one or more overlay networks to communicate with the
backend services. By default, a single network is created called `ucp-hrm`,
with the access control label `ucp-hrm`. Adding a service to this network
either requires administrator-level access, or the user must be in a group that
gives them `ucp-hrm` access.

This default configuration does not provide any isolation between services
using the HTTP routing mesh.

Isolation between services may be implemented by creating one or more overlay
networks with the label `com.docker.ucp.mesh.http` prior to enabling the HTTP
routing mesh. Once the HTTP routing mesh is enabled, it will be able to route
to all services attached to any of these networks, but services on different
networks cannot communicate directly.

## Using the HTTP routing mesh

Once DNS and networks are configured, you can begin setting up services for
these domains. See the guides for the [UCP web
UI](../../user/services/use-hostnames-to-access-your-service.md) and [Docker
CLI](../../user/services/hrm-labels.md).

## Disable the HTTP routing mesh

To disable the HTTP routing mesh, first ensure that all services that are using
the HTTP routing mesh are disconnected from the **ucp-hrm** network.

Next, go to the **UCP web UI**, navigate to the **Settings** page, and click
the **Routing Mesh** tab. Uncheck the checkbox to disable the HTTP routing mesh.

## Troubleshoot

If a service is not configured properly for use of the HTTP routing mesh, this
information is available in the UI when inspecting the service.

More logging from the HTTP routing mesh is available in the logs of the
`ucp-controller` containers on your UCP manager nodes.
