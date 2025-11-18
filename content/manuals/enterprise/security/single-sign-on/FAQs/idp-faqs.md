---
title: SSO identity provider FAQs
linkTitle: Identity providers
description: Frequently asked questions about Docker SSO and identity provider configuration
keywords: identity providers, SSO IdP, SAML, Azure AD, Entra ID, certificate management
tags: [FAQ]
aliases:
- /single-sign-on/idp-faqs/
- /faq/security/single-sign-on/idp-faqs/
- /security/faqs/single-sign-on/idp-faqs/
---

## Can I use multiple identity providers with Docker SSO?

Yes, Docker supports multiple IdP configurations. A domain can be associated with multiple IdPs. Docker supports Entra ID (formerly Azure AD) and identity providers that support SAML 2.0.

## Can I change my identity provider after configuring SSO?

Yes. Delete your existing IdP configuration in your Docker SSO connection, then [configure SSO using your new IdP](/manuals/enterprise/security/single-sign-on/connect.md). If you had already turned on enforcement, turn off enforcement before updating the provider connection.

## What information do I need from my identity provider to configure SSO?

To turn on SSO in Docker, you need the following from your IdP:

- SAML: Entity ID, ACS URL, Single Logout URL, and the public X.509 certificate
- Entra ID (formerly Azure AD): Client ID, Client Secret, AD Domain

## What happens if my existing certificate expires?

If your certificate expires, contact your identity provider to retrieve a new X.509 certificate. Then update the certificate in the [SSO configuration settings](/manuals/enterprise/security/single-sign-on/manage.md#manage-sso-connections) in the Docker Admin Console.

## What happens if my IdP goes down when SSO is turned on?

If SSO is enforced, users can't access Docker Hub when your IdP is down. Users can still access Docker Hub images from the CLI using personal access tokens.

If SSO is turned on but not enforced, users can fall back to username/password authentication.

## Do bot accounts need seats to access organizations using SSO?

Yes, bot accounts need seats like regular users, requiring a non-aliased domain email in the IdP and using a seat in Docker Hub. You can add bot accounts to your IdP and create access tokens to replace other credentials.

## Does SAML SSO use Just-in-Time provisioning?

The SSO implementation uses Just-in-Time (JIT) provisioning by default. You can optionally turn off JIT in the Admin Console if you turn on auto-provisioning using SCIM. See [Just-in-Time provisioning](/security/for-admins/provisioning/just-in-time/).

## My Entra ID SSO connection isn't working and shows an error. How can I troubleshoot this?

Confirm that you've configured the necessary API permissions in Entra ID for your SSO connection. You need to grant administrator consent within your Entra ID tenant. See [Entra ID (formerly Azure AD) documentation](https://learn.microsoft.com/en-us/azure/active-directory/manage-apps/grant-admin-consent?pivots=portal#grant-admin-consent-in-app-registrations).
