---
description: Image Access Management
keywords: image, access, management, trusted content, permissions, Docker Business feature
title: Image Access Management
aliases:
- /docker-hub/image-access-management/
- /desktop/hardened-desktop/image-access-management/
- /admin/organization/image-access/
---

> **Note**
>
> Image Access Management is available to [Docker Business](../../subscription/core-subscription/details.md#docker-business) customers only.

Image Access Management gives administrators control over which types of images, such as Docker Official Images, Docker Verified Publisher Images, or community images, their developers can pull from Docker Hub.

For example, a developer, who is part of an organization, building a new containerized application could accidentally use an untrusted, community image as a component of their application. This image could be malicious and pose a security risk to the company. Using Image Access Management, the organization owner can ensure that the developer can only access trusted content like Docker Official Images, Docker Verified Publisher Images, or the organization’s own images, preventing such a risk.

## Prerequisites

You need to [configure a registry.json to enforce sign-in](enforce-sign-in/_index.md). For Image Access Management to take effect, Docker Desktop users must authenticate to your organization.

## Configure Image Access Management permissions

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-image-access product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-image-access product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## More resources

- [Video: Hardened Desktop Image Access Management](https://www.youtube.com/watch?v=r3QRKHA1A5U)
