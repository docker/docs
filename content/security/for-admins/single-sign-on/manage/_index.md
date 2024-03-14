---
description: Learn how to manage Single Sign-On for your organization or company.
keywords: manage, single sign-on, SSO, sign-on, docker hub, admin console, admin, security
title: Manage Single Sign-On
aliases:
- /admin/company/settings/sso-management/
- /single-sign-on/manage/
---

## Manage organizations

> **Note**
>
> You must have a [company](/admin/company/) to manage more than one organization.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-management-orgs product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-sso-management-orgs product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## Manage domains

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-management product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-sso-management product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## What's next?

- [Set up SCIM](../../scim.md)
- [Enable Group mapping](../../group-mapping.md)

