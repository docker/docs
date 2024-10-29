---
description: Learn how to manage Single Sign-On for your organization or company.
keywords: manage, single sign-on, SSO, sign-on, docker hub, admin console, admin, security
title: Manage single sign-on
linkTitle: Manage
aliases:
- /admin/company/settings/sso-management/
- /single-sign-on/manage/
---

## Manage organizations

> [!NOTE]
>
> You must have a [company](/admin/company/) to manage more than one organization.

{{< include "admin-early-access.md" >}}

{{% admin-sso-management-orgs product="admin" %}}

## Manage domains

{{< tabs >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-sso-management product="admin" %}}

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-management product="hub" %}}

{{< /tab >}}
{{< /tabs >}}

## Manage SSO connections

{{< tabs >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-sso-management-connections product="admin" %}}

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-management-connections product="hub" %}}

{{< /tab >}}
{{< /tabs >}}

## Manage users

{{< include "admin-early-access.md" >}}

{{% admin-sso-management-users product="admin" %}}

## Manage provisioning

Users are provisioned with Just-in-Time (JIT) provisioning by default. If you enable SCIM, you can disable JIT. For more information, see the [Provisioning overview](/manuals/security/for-admins/provisioning/_index.md) guide.

## What's next?

- [Set up SCIM](../provisioning/scim.md)
- [Enable Group mapping](../provisioning/group-mapping.md)

