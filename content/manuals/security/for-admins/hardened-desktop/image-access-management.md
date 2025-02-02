---
description: Manage Docker Hub image access with Image Access Management, restricting developers to trusted images for enhanced security
keywords: image, access, management, trusted content, permissions, Docker Business feature, security, admin
title: Image Access Management
tags: [admin]
aliases:
 - /docker-hub/image-access-management/
 - /desktop/hardened-desktop/image-access-management/
 - /admin/organization/image-access/
 - /security/for-admins/image-access-management/
weight: 40
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Image Access Management gives you control over which types of images, such as Docker Official Images, Docker Verified Publisher Images, or community images, your developers can pull from Docker Hub.

For example, a developer, who is part of an organization, building a new containerized application could accidentally use an untrusted, community image as a component of their application. This image could be malicious and pose a security risk to the company. Using Image Access Management, the organization owner can ensure that the developer can only access trusted content like Docker Official Images, Docker Verified Publisher Images, or the organization’s own images, preventing such a risk.

## Prerequisites

You first need to [enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) to ensure that all Docker Desktop developers authenticate with your organization. Since Image Access Management requires a Docker Business subscription, enforced sign-in guarantees that only authenticated users have access and that the feature consistently takes effect across all users, even though it may still work without enforced sign-in.

## Configure

{{< summary-bar feature_name="Admin console early access" >}}

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-image-access product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{% admin-image-access product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## More resources

- [Video: Hardened Desktop Image Access Management](https://www.youtube.com/watch?v=r3QRKHA1A5U)
