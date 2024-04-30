---
description: Learn how to complete your single-sign on connection and next steps for enabling SSO.
keywords: configure, sso, docker hub, hub, docker admin, admin, security
title: Complete your single sign-on connection
---

The steps to set up your SSO configuration are:

1. [Add and verify the domain or domains](/security/for-admins/single-sign-on/configure#step-one-add-and-verify-your-domain) that your members use to sign in to Docker.
2. [Create your SSO connection](/security/for-admins/single-sign-on/configure#step-two-create-an-sso-connection-in-docker) in Docker.
3. [Configure your IdP](/security/for-admins/single-sign-on/configure/configure-idp#step-three-configure-your-idp-to-work-with-docker) to work with Docker.
4. [Complete your SSO connection](#step-four-complete-your-sso-connection) in Docker.

This page walks you through the final steps of creating your SSO connection. You can then test your connection and optionally enforce SSO for your organization.

## Prerequisites

Make sure you have completed the following before you begin:

- Your domain is verified
- You have created your SSO connection in Docker
- You configured your IdP using the appropriate values from your Docker connection
- You have pasted the following from your IdP into the settings in the Docker console:
    - SAML: **SAML Sign-on URL**, **x509 Certificate**
    - Azure AD (OIDC): **Client ID**, **Client Secret**, **Azure AD Domain**

## Step four: Complete your SSO connection

> **Beta feature**
>
> Optional Just-in-Time (JIT) provisioning is available in Private Beta when you use the Admin Console. If you're participating in this program, you have the option to turn off this default provisioning and disable JIT. This configuration is recommended if you're using SCIM to auto-provision users. See [SSO authentication with JIT provisioning disabled](/security/for-admins/group-mapping/#sso-authentication-with-jit-provisioning-disabled).
{ .experimental }

{{< tabs >}}
{{< tab name="Admin Console" >}}

{{% admin-sso-connect product="admin" %}}

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-connect product="hub" %}}

{{< /tab >}}
{{< /tabs >}}

## More resources

The following videos demonstrate how to enforce SSO.

- [Video: Enforce SSO with Okta SAML](https://youtu.be/c56YECO4YP4?feature=shared&t=1072)
- [Video: Enforce SSO with Azure AD (OIDC)](https://youtu.be/bGquA8qR9jU?feature=shared&t=1087)


## What's next

Learn how you can [manage your SSO connection](../manage/_index.md), domain, and users for your organization or company.
