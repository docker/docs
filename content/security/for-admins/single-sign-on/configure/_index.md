---
description: Learn how to configure Single Sign-On for your organization or company.
keywords: configure, sso, docker hub, hub, docker admin, admin, security 
title: Configure Single Sign-On
aliases:
- /docker-hub/domains/
- /docker-hub/sso-connection/
- /docker-hub/enforcing-sso/
- /single-sign-on/configure/
- /admin/company/settings/sso-configuration/
- /admin/organization/security-settings/sso-configuration/
---

Follow the steps on this page to configure SSO for your organization or company.

## Step one: Add and verify your domain

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-domains product="hub" %}}

{{< /tab >}}
{{< tab name="Docker Admin" >}}

{{< include "admin-early-access.md" >}}

{{% admin-domains product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## Step two: Create an SSO connection

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-config product="hub" %}}

{{< /tab >}}
{{< tab name="Docker Admin" >}}

{{% admin-sso-config product="admins" %}}

{{< /tab >}}
{{< /tabs >}}



