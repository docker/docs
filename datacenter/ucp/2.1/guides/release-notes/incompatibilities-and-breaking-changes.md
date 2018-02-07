---
title: Incompatibilities and breaking changes
description: Learn about the incompatibilities and breaking changes introduced by Universal Control Plane version {{ page.ucp_version }}
keywords: docker, ucp, upgrade, incompatibilities
redirect_from:
- /datacenter/ucp/2.1/guides/admin/upgrade/incompatibilities-and-breaking-changes/
---

This is the list of incompatibilities introduced by Universal Control Plane
{{ page.ucp_version }}.

## HTTP routing mesh

When using the HTTP routing mesh you need to apply specific labels to your
services to make them accessible using a hostname.
The syntax used for these labels has changed on 2.1.

If you were using this feature on UCP 2.0, after you upgrade you need to
update your services to use a new label syntax.
You can do this from the UCP web UI or from the CLI using the
`docker service update` command.

There are two changes to consider.

First, a route was previously in the format of `internal_port=external_route`
or just `external_route`. Now the format is a comma separated list of
`key=value` pairs.

Second, if you have multiple routes to the same service, these were previously
written as a comma separated list of the above. These are now separate labels,
one per route, prefixed with `com.docker.ucp.mesh.http`.

Below you can find examples on how to upgrade from the old syntax to the new
one.

### A single route with a single internal port

```none
Old syntax
http://example.com

New syntax
external_route=http://example.com

How to upgrade
docker service update \
  --label-add com.docker.ucp.mesh.http=http://example.com \
  <service-name>
```

### A single route with an explicit internal port

```none
Old syntax
8080=http://example.com

New syntax
external_route=http://example.com,internal_port=8080

How to upgrade
docker service update \
  --label-add com.docker.ucp.mesh.http=external_route=http://example.com,internal_port=8080 \
  <service-name>
```

### Two routes, each to a different internal port

```none
Old syntax
8080=http://foo.example.com,1234=http://bar.example.com

New syntax
external_route=http://foo.example.com,internal_port=8080
external_route=http://bar.example.com,internal_port=1234

How to upgrade
docker service update \
  --label-remove com.docker.ucp.mesh.http \
  --label-add com.docker.ucp.mesh.http.foo=external_route=http://foo.example.com,internal_port=8080 \
  --label-add com.docker.ucp.mesh.http.bar=external_route=http://bar.example.com,internal_port=1234 \
  <service-name>
```

 <h2>Where to go next</h2>

 * [Release notes](index.md)
