---
title: Namespace access control
linkTitle: Namespace access
description: Control whether organization members can push content to their personal namespaces on Docker Hub
keywords: namespace access, docker hub, personal namespace, organization security, docker business
tags: [admin]
weight: 50
---

{{< summary-bar feature_name="Namespace access" >}}

Namespace access control lets organization administrators control whether all
members of an organization can push content to their personal namespaces on
Docker Hub. This prevents organizations from accidentally publishing images
outside of approved, governed locations.

When namespace access control is enabled, organization members can still view and pull images
from their personal namespaces and continue accessing all existing repositories
and content. However, they will no longer be able to create new repositories or
push new images to their personal namespace.

> [!IMPORTANT]
>
> For users in multiple organizations, if namespace access control is enabled in
> any organization, that user cannot push to their personal namespace and cannot
> create new repositories in their personal namespace.

### Configure namespace access control

To configure namespace access control:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization from the top-left account drop-down.
2. Select **Admin Console**, then **Namespace access**.
3. Use the toggle to enable or disable namespace access control.
4. Select **Save changes**.

Once namespace access control is enabled, organization members can still view their
personal namespace and existing repositories but they will not be able to create
any new repositories or push any new images to existing repositories.

### Verify access restrictions

After configuring namespace access control, test that restrictions work correctly.

After any attempt to push to an existing repository in your personal namespace,
you'll see an error message like the following:

```console
$ docker push <personal-namespace>/<image>:<tag>
Unavailable
authentication required - namespace access restriction from an organization you belong to prevents pushing new content in your personal namespace. Restriction applied by: <organizations>. Please contact your organization administrator
```