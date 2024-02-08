---
title: Docker Build Cloud FAQ
description: Frequently asked questions about Docker Build Cloud
keywords: build, cloud build, faq, troubleshooting 
aliases:
  - /hydrobuild/faq/
---

<!--toc:start-->
- [How do I remove Docker Build Cloud from my system?](#how-do-i-remove-docker-build-cloud-from-my-system)
- [Are builders shared between organizations?](#are-builders-shared-between-organizations)
- [Do I need to add my secrets to the builder to access private resources?](#do-i-need-to-add-my-secrets-to-the-builder-to-access-private-resources)
- [How do I unset Docker Build Cloud as the default builder?](#how-do-i-unset-docker-build-cloud-as-the-default-builder)
- [How do I manage the build cache with Docker Build Cloud?](#how-do-i-manage-the-build-cache-with-docker-build-cloud)
- [Can I use Docker Build Cloud with a registry behind a VPN?](#can-i-use-docker-build-cloud-with-a-registry-behind-a-vpn)
<!--toc:end-->

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

### Can I use Docker Build Cloud with a registry behind a VPN?

No, you can't use Docker Build Cloud with a private registry or registry mirror
behind a VPN. All endpoints invoked with Docker Build Cloud, including OCI
registries, must be accessible over the internet.
