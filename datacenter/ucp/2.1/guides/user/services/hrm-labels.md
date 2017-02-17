---
title: Use the Docker CLI to configure hostnames to access your service
description: Learn how to configure your services to make them accessible using a hostname instead of IP addresses and ports, using the Docker CLI.
keywords: ucp, services, http, dns
---

UCP's HTTP routing mesh provides configuration through the web UI, but it is
also possible to configure a service by using the Docker CLI directly. This
information is stored in the labels of the services and other configuration of
the service. This guide will walk you through configuring a sample service for
use with the HTTP routing mesh and also provide a full reference to the label
syntax.

This configuration can be performed manually using the Docker CLI. It is also
appropriate to use this in automation, either using the Docker CLI or the
Docker Remote API directly.

The HTTP routing mesh can route to a Docker service that runs a webserver (HTTP
or HTTPS). This service must meet three criteria:

* The service must be connected a network with a `com.docker.ucp.mesh.http` label
* The service must publish the ports that you wish to route
* The service must have one or more labels prefixed with
  `com.docker.ucp.mesh.http` to specify the ports to route (see the syntax
  below)

## Route to a service using the CLI

### Networks

Services must be connected to a network that has a `com.docker.ucp.mesh.http`
label. The value is not relevant. A `ucp-hrm` network is created for you
automatically when enabling the HTTP routing mesh, or your administrators may
create one for you. Refer to the administrator's guide for more information.

### Service Labels

The key of the label must begin with `com.docker.ucp.mesh.http`. For multiple
labels, some examples could be `com.docker.ucp.mesh.http.80` and
`com.docker.ucp.mesh.http.443`. Here `80` and `443` are used to differentiate
the HRM labels via port numbers. You can use whatever values you want, just
make sure they are different from each other and you can keep track of them.

Labels with the prefix `com.docker.ucp.mesh.http` allow you to configure a
single hostname and port to route to a service. If you wish to route multiple
ports or hostnames to the same service, then multiple labels with the prefix
`com.docker.ucp.mesh.http` may be created.

### Example using the CLI

A service based on the image `myimage/mywebserver:latest` with a webserver
running on port 8080 can be routed to `http://foo.example.com` can be created
using the following:

```sh
$ docker service create \
  -p 8080 \
  --network ucp-hrm \
  --label com.docker.ucp.mesh.http.8080=external_route=http://foo.example.com,internal_port=8080 \
  --name myservice \
  myimage/mywebserver:latest
```

## Service Label Syntax

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

* `com.docker.ucp.mesh.http.1=external_route=http://old.example.com.com,redirect=http://new.example.com`
* `com.docker.ucp.mesh.http.2=external_route=http://new.example.com`
