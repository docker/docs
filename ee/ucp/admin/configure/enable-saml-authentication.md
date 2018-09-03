---
title: Enable SAML authentication
description: Learn how configure user authentication with SAML 2.0
keywords: SAML, ucp, authentication
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

UCP supports Security Assertion Markup Language (SAML) for authentication. SAML is an open standard for exchanging authentication data between an identity provider and a service provider, such as UCP. SAML is commonly supported by enterprise authentication systems. SAML-based single sign-on (SSO) gives you access to UCP through a SAML 2.0-compliant identity provider. For more information about SAML, see the [SAML XML website] (http://saml.xml.org/)

A list of the identity providers that we support is found in the [Docker Compatibility matrix] (http://success.docker.com/article/compatibility-matrix).

## Identity provider integration steps

There are data your identity provider needs for successful integration with UCP. Typivally these are:

- URL for single signon (SSO).
- Service provider audience URI.

## Configure the SAML integration

To enable SAML authentication, go to the UCP web UI, then navigate to the **Admin Settings**. Select **Authentication & Authorization** to enable SAML.

[UI SCREENSHOT]

In the **SAML Enabled** section, select **Yes** to display the required settings.

## Security considerations

You can download a client bundle to access UCP. To ensure that access from the client bundle is synced with the identity provider, we recommend the following steps. Otherwise, a previously-authorized user could get access to UCP through their existing client bundle.

- Remove the user account from UCP granting client bundle access if access is removed from the identity provider.
- If group membership in the identity provider changes, replicate this change in UCP.
- Continue to use LDAP to sync group membership.
