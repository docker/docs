---
description: Docker Universal Control Plane
keywords: networking, kv, engine-discovery, ucp
title: Enable container networking with UCP
---

UCP provides an HTTP routing mesh, that extends the networking capabilities
of Docker Engine. Docker Engine provides load balancing and service discovery
at the transport layer for TCP and UDP connections. UCP's HTTP routing mesh
allows you to extend service discovery to have name-based virtual hosting for
HTTP services.

See the
[Docker Engine documentation on overlay networks](/engine/swarm/networking.md)
for more information on what Docker Engine provides.

This feature is currently experimental.

## Enable the HTTP routing mesh

To enable the HTTP routing mesh, go to the **UCP web UI**, navigate to the
**Settings** page, and click the **Routing Mesh** tab.

<!-- todo: add screenshot -->

The default port for HTTP services is **80**. You may choose an alternate port
on this screen.

Check the checkbox to enable the HTTP routing mesh. This will create a service
called `ucp-hrm` and a network called `ucp-hrm`.

## Route to a service

The HTTP routing mesh can route to a Docker service that runs a webserver.
This service must meet three criteria:

* The service must be connected to the `ucp-hrm` network
* The service must publish one or more ports
* The service must have a `com.docker.ucp.mesh.http` label to specify the ports
to route

The syntax for the `com.docker.ucp.mesh.http` label is a list of one or more
values separated by commas. Each of these values is in the form of
`internal_port=protocol://host`, where:

* `internal_port` is the port the service is listening on (and may be omitted
if there is only one port published)
* `protocol` is `http`
* `host` is the hostname that should be routed to this service

Examples:

A service based on the image `myimage/mywebserver:latest` with a webserver running on port
8080 can be routed to `http://foo.example.com` can be created using the
following:

```sh
$ docker service create \
  -p 8080 \
  --network ucp-hrm \
  --label com.docker.ucp.mesh.http=8080=http://foo.example.com \
  --name myservice \
  myimage/mywebserver:latest
```

The HTTP Routing Mesh checks for new services every 60 seconds, so it may take
up to one minute for configuration to complete.

Next, you will need to route the referenced domains to the HTTP routing mesh.

## Route domains to the HTTP routing mesh

The HTTP routing mesh uses the `Host` HTTP header to determine which service
should receive a particular HTTP request. This is typically done using DNS and
pointing one or more domains to one or more nodes in the UCP cluster.

## Disable the HTTP routing mesh

To disable the HTTP routing mesh, first ensure that all services that are using
the HTTP routing mesh are disconnected from the **ucp-hrm** network.

Next, go to the **UCP web UI**, navigate to the **Settings** page, and click
the **Routing Mesh** tab. Uncheck the checkbox to disable the HTTP routing mesh.

## Access Control

To route a domain to the HTTP Routing Mesh, the service must be on the
`ucp-hrm` network which has the `ucp-hrm` access label. Adding a service to
this network either requires administrator-level access, or the user must be in
a group that gives them `ucp-hrm` access.

## Troubleshoot

Check the logs of the `ucp-controller` containers on your UCP controller nodes
for logging from the HTTP routing mesh.
