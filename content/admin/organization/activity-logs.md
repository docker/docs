---
description: Learn about activity logs.
keywords: team, organization, activity, log, audit, activities
title: Activity logs
aliases:
- /docker-hub/audit-log/
---

Audit logs display a chronological list of activities that occur at organization and repository levels. It provides a report to owners on all their member activities.

With audit logs, owners can view and track:
 - What changes were made
 - The date when a change was made
 - Who initiated the change

 For example, audit logs display activities such as the date when a repository was created or deleted, the member who created the repository, the name of the repository, and when there was a change to the privacy settings.

Owners can also see the audit logs for their repository if the repository is part of the organization subscribed to a Docker Business or Team plan.

> **Note**
>
> Audit logs requires a [Docker Team or Business subscription](../../subscription/_index.md).

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-org-audit-log product="hub" %}}

{{< /tab >}}
{{< tab name="Docker Admin" >}}

{{< include "admin-early-access.md" >}}

{{% admin-org-audit-log product="admin" %}}

{{< /tab >}}
{{< /tabs >}}
