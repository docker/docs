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

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-sso-connect product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{% admin-sso-connect product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## What's next

Learn how you can [manage your SSO connection](../manage/_index.md), domain, and users for your organization or company.
