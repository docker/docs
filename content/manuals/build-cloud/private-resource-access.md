---
title: Private resources access
description: Use Docker Build Cloud with private package repositories
keywords: docker build build, private packages, registry, package repository, vpn, forward-hosts
---

{{% experimental %}}
This feature is experimental and subject to change without notice.
{{% /experimental %}}

If your builds require access to images or other packages hosted on private or
on-premises artifact repositories, the builder must be able to access these
resources, or the build will fail. For example, if your organization hosts a
private [PyPI](https://pypi.org/) repository on a VPN, Docker Build Cloud would
not be able to access it by default, since it isn't on the same network as your
VPN.

To enable the cloud builder to access your internal resources, you can
configure your Buildx client to act as a forward proxy for the builder to
access your internal resources through a secure tunnel.

## How it works

When you create a cloud builder with the `forward-hosts` option, the builder
uses the client as a proxy to forward requests to the specified hostnames.

To use this feature, you need:

- Docker Desktop version 4.32.0 or later (Buildx version `0.15.1-desktop.1`)

## Configuration

Connection-forwarding is configured on a per-builder and per-client basis. To
enable it, set the `forward-hosts` driver option when creating the builder. For
example, to create a builder that forwards connections to an internal registry
on `registry.example.com:5000`, run the following command:

```console
$ docker buildx create \
  --driver cloud \
  --driver-opt "forward-hosts=registry.example.com:5000" \
  <org>/<builder>
```

Replace `registry.example.com` with the hostname of your internal registry, and
`<org>/<builder>` with the namespace of your Docker organization and the name
of your builder.

After creating the builder with the `forward-hosts` option, you can use it to
build images with private packages from the hostnames you specified.

### Multiple hostnames

If you need to forward connections to multiple hostnames, specify each hostname
separated by a semicolon. The following example enables forwarding for
`pypi.internal` and `npm.internal`:

```console
$ docker buildx create \
  --driver cloud \
  --driver-opt "forward-hosts=pypi.internal;npm.internal" \
  <org>/<builder>
```

You can also use wildcards to forward connections for all hostnames that match
a specific pattern. The following example forwards all requests to hostnames
ending with `*.internal`:

```console
$ docker buildx create \
  --driver cloud \
  --driver-opt "forward-hosts=*.internal" \
  <org>/<builder>
```

To forward all hosts, use `*`:

```console
$ docker buildx create \
  --driver cloud \
  --driver-opt "forward-hosts=*" \
  <org>/<builder>
```

## Authentication

If your internal artifacts require authentication, make sure that you
authenticate with the repository either before or during the build. For
internal packages like npm or PyPI, use [build secrets](/manuals/build/building/secrets.md)
to authenticate during the build. For internal OCI registries, use `docker
login` to authenticate before building.

Note that if you use a private registry that requires authentication, you will
need to authenticate with `docker login` twice before building. This is because
the cloud builder needs to authenticate with Docker to use the cloud builder,
and then again to authenticate with the private registry.

```console
$ echo $DOCKER_PAT | docker login docker.io -u <username> --password-stdin
$ echo $REGISTRY_PASSWORD | docker login registry.example.com -u <username> --password-stdin
$ docker build --builder <cloud-builder> --tag registry.example.com/<image> --push .
```
