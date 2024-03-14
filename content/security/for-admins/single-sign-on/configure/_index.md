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
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-domains product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## Step two: Create an SSO connection

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-config product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{% admin-sso-config product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## More resources

The following video provides an overview of configuring SSO with SAML in Entra ID (formerly Azure AD).

<iframe title="Configure SSO with SAML in Entra ID overview" class="border-0 w-full aspect-video mb-8" allow="fullscreen" src="https://www.loom.com/embed/0a30409381f340cfb01790adbd9aa9b3?sid=7e4e10a7-7f53-437d-b593-8a4886775632"></iframe>

## What's next?

- [Set up SCIM](../../scim.md)
- [Enable Group mapping](../../group-mapping.md)
- [Manage your SSO connections](../manage/_index.md)

