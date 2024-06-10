---
description: Learn how to audit your domains for uncaptured users.
keywords: domain audit, security, identify users, manage users
title: Domain audit
aliases:
- /docker-hub/domain-audit/
- /admin/company/settings/domains/
- /admin/organization/security-settings/domains/
---

Domain audit identifies uncaptured users in an organization. Uncaptured users are Docker users who have authenticated to Docker using an email address associated with one of your verified domains, but they're not a member of your organization in Docker. You can audit domains on organizations that are part of the Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

Uncaptured users who access Docker Desktop in your environment may pose a security risk because your organization's security settings, like Image Access Management and Registry Access Management, aren't applied to a user's session. In addition, you won't have visibility into the activity of uncaptured users. You can add uncaptured users to your organization to gain visibility into their activity and apply your organization's security settings.

Domain audit can't identify the following Docker users in your environment:

- Users who access Docker Desktop without authenticating
- Users who authenticate using an account that doesn't have an email address associated with one of your verified domains

Although domain audit can't identify all Docker users in your environment, you can enforce sign-in to prevent unidentifiable users from accessing Docker Desktop in your environment. For more details about enforcing sign-in, see [Configure registry.json to enforce sign-in](configure-sign-in.md).

> **Tip**
>
> You can use endpoint management (MDM) software to search for the number of Docker Desktop instances and the Docker Desktop versions in your environment. This can provide accurate license reporting and ensure your machines use the latest Docker Desktop version. 
> - [Intune](https://learn.microsoft.com/en-us/mem/intune/apps/app-discovered-apps)
> - [Jamf](https://docs.jamf.com/10.25.0/jamf-pro/administrator-guide/Application_Usage.html)
> - [Kandji](https://support.kandji.io/support/solutions/articles/72000559793-view-a-device-application-list)
> - [Kolid](https://www.kolide.com/features/device-inventory/properties/mac-apps)
> - [Workspace One](https://blogs.vmware.com/euc/2022/11/how-to-use-workspace-one-intelligence-to-manage-app-licenses-and-reduce-costs.html)
{ .tip }

## Prerequisites

Before you audit your domains, review the following required prerequisites:

- Your organization must be part of a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](../../subscription/core-subscription/upgrade.md).
- You must [add and verify your domains](./single-sign-on/configure/_index.md#step-one-add-and-verify-your-domain).

> **Important**
>
> Domain audit is not supported for companies or organizations within a company.
{ .important }

## Audit your domains for uncaptured users

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-domain-audit product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-domain-audit product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

