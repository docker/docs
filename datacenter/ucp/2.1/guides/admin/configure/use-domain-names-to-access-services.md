---
description: Docker Universal Control Plane
keywords: networking, kv, engine-discovery, ucp
title: Enable container networking with UCP
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

### Route to a service using the UI

<!-- todo: expand, add screenshots -->

When using the wizard to create services in the UCP UI, you may enable the HTTP
routing mesh for a service by publishing a port, and filling in the options
relating to "Routing Mesh Host".

### Route to a service using the CLI

The key of the label must begin with `com.docker.ucp.mesh.http`. For multiple
labels, some examples could be `com.docker.ucp.mesh.http.80` and
`com.docker.ucp.mesh.http.443`. Here `80` and `443` are used to differentiate
the HRM labels via port numbers. You can use whatever values you want, just
make sure they are different from each other and you can keep track of them.

Labels with the prefix `com.docker.ucp.mesh.http` allow you to configure a
single hostname and port to route to a service. If you wish to route multiple
ports or hostnames to the same service, then multiple labels with the prefix
`com.docker.ucp.mesh.http` may be created.

The syntax of this label is as follows:

The key of the label must begin with `com.docker.ucp.mesh.http`, for example
`com.docker.ucp.mesh.http.80` and `com.docker.ucp.mesh.http.443`.

The value of the label is a comma separated list of key/value pairs separated
by equals signs. These pairs are optional unless noted below, and are as
follows:

* `external_route` **(required)** the external URL to route to this service.
  Examples: `http://myapp.example.com` and `sni://myapp.example.com`
* `internal_port`: the internal port to use for the service.  Examples: `80`,
  `8443`. This is **required** if more one port is published by the service.
* `sticky_sessions`: if present, use the named cookie to route the user to the
  same backend task for this service. See the "Sticky Sessions" section below.
* `redirect`: if present, perform redirection to the specified URL. See the
  "Redirection" section below.

Examples:

A service based on the image `myimage/mywebserver:latest` with a webserver running on port
8080 can be routed to `http://foo.example.com` can be created using the
following:

```sh
$ docker service create \
  -p 8080 \
  --network ucp-hrm \
  --label com.docker.ucp.mesh.http.8080=external_route=http://foo.example.com,internal_port=8080 \
  --name myservice \
  myimage/mywebserver:latest
```

Next, you will need to route the referenced domains to the HTTP routing mesh.

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

## Disable the HTTP routing mesh

To disable the HTTP routing mesh, first ensure that all services that are using
the HTTP routing mesh are disconnected from the **ucp-hrm** network.

Next, go to the **UCP web UI**, navigate to the **Settings** page, and click
the **Routing Mesh** tab. Uncheck the checkbox to disable the HTTP routing mesh.

## Additional Routing Features

The HTTP routing mesh provides some additional features for some specific use
cases.

### Sticky Sessions

Enable the sticky sessions option for a route if your application requires that
a user's session continues to use the same task of a backend service. This
option uses HTTP cookies to choose which task receives a given connection.

The cookie name for this feature is configured as the value of this option
within the label. The cookie must be created by the application, and its value
is used to pick a backend task.

Stickyness may be lost temporarily if the number of tasks for a service
changes, or if a service is reconfigured in a way that requires all of its
tasks to restart.

This option is incompatible with the `sni` protocol (routing HTTPS connections
without termination).

### Redirection

The `redirect` option indicates that all requests to this route should be
redirected to another domain name using a HTTP redirect.

One use of this feature is for a service which only listens using HTTPS, with
HTTP traffic to it being redirected to HTTPS. If the service is on
`example.com`, then this can be accomplished with two labels:

* `com.docker.ucp.mesh.http.1=external_route=http://example.com,redirect=https://example.com`
* `com.docker.ucp.mesh.http.2=external_route=sni://example.com`

Another use is a service expecting traffic only on a single domain, but other
domains should be redirected to it. For example, a website that has been
renamed might use this functionality. The following labels accomplish this for
`new.example.com` and `old.example.com`

* `com.docker.ucp.mesh.http.1=external_route=http://old.example.com,redirect=http://new.example.com`
* `com.docker.ucp.mesh.http.2=external_route=http://new.example.com`

## Troubleshoot

If a service is not configured properly for use of the HTTP routing mesh, this
information is available in the UI when inspecting the service.

More logging from the HTTP routing mesh is available in the logs of the
`ucp-controller` containers on your UCP manager nodes.
