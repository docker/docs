---
title: Docker Build Cloud FAQ
description: Frequently asked questions about Docker Build Cloud
keywords: build, cloud build, faq, troubleshooting 
aliases:
  - /hydrobuild/faq/
---

### How do I remove Docker Build Cloud from my system?

If you want to stop using Docker Build Cloud, remove the cloud builder using
the `docker buildx rm` command.

```console
$ docker buildx rm cloud-<ORG>-default
```

This doesn't deprovision the builder backend, it only removes the builder from
your local Docker client.

### Are builders shared between organizations?

No. Each cloud builder provisioned to an organization is completely
isolated to a single Amazon EC2 instance, with a dedicated EBS volume for build
cache, and end-to-end encryption. That means there are no shared processes or
data between cloud builders.

### Do I need to add my secrets to the builder to access private resources?

No. Your interface to Docker Build Cloud is Buildx, and you can use the existing
`--secret` and `--ssh` CLI flags for managing build secrets.

For more information, refer to:

- [`docker buildx build --secret`](/engine/reference/commandline/buildx_build/#secret)
- [`docker buildx build --ssh`](/engine/reference/commandline/buildx_build/#ssh)

### How do I unset Docker Build Cloud as the default builder?

If you've set a cloud builder as the default builder and want to revert to using the
default `docker` builder, run the following command:

```console
$ docker context use default
```

### How do I manage the build cache with Docker Build Cloud?

You don't need to manage the builder's cache manually. The system manages it
for you through [garbage collection](/build/cache/garbage-collection/).

Old cache is automatically removed if you hit your storage limit. You can check
your current cache state using the
[`docker buildx du` command](/engine/reference/commandline/buildx_du/).

To clear the builder's cache manually, you can use the
[`docker buildx prune` command](/engine/reference/commandline/buildx_prune/)
command. This works like pruning the cache for any other builder.

> **Note**
>
> Pruning a cloud builder's cache also removes the cache for other team members
> using the same builder.
