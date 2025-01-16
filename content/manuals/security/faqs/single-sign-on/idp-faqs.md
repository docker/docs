---
description: Single sign-on IdP FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, IdP
title: FAQs on SSO and identity providers
tags: [FAQ]
aliases:
- /single-sign-on/idp-faqs/
- /faq/security/single-sign-on/idp-faqs/
---

### Is it possible to use more than one IdP with Docker SSO?

No. You can only configure Docker SSO to work with a single IdP. A domain can only be associated with a single IdP. Docker supports Entra ID (formerly Azure AD) and identity providers that support SAML 2.0.

### Is it possible to change my identity provider after configuring SSO?

Yes. You must delete your existing IdP configuration in your Docker SSO connection and then [configure SSO using your new IdP](/manuals/security/for-admins/single-sign-on/connect.md). If you had already turned on enforcement, you should turn off enforcement before updating the provider SSO connection.

### What information do I need from my identity provider to configure SSO?

To enable SSO in Docker, you need the following from your IdP:

* **SAML**: Entity ID, ACS URL, Single Logout URL and the public X.509 certificate

* **Entra ID (formerly Azure AD)**: Client ID, Client Secret, AD Domain.

### What happens if my existing certificate expires?

If your existing certificate has expired, you may need to contact your identity provider to retrieve a new X.509 certificate. Then, you need to update the certificate in the [SSO configuration settings](/security/for-admins/single-sign-on/manage/#manage-sso-connections) in Docker Hub or Docker Admin Console.

### What happens if my IdP goes down when SSO is enabled?

If SSO is enforced, then it is not possible to access Docker Hub when your IdP is down. You can still access Docker Hub images from the CLI using your Personal Access Token.

If SSO is enabled but not enforced, then users could fallback to authenticate with username/password and trigger a reset password flow (if necessary).

### How do I handle accounts using Docker Hub as a secondary registry? Do I need a bot account?

You can add a bot account to your IdP and create an access token for it to replace the other credentials.

### Does a bot account need a seat to access an organization using SSO?

Yes, bot accounts need a seat, similar to a regular end user, having a non-aliased domain email enabled in the IdP and using a seat in Hub.

### Does SAML SSO use Just-in-Time provisioning?

The SSO implementation uses Just-in-Time (JIT) provisioning by default. You can optionally disable JIT in the Admin Console if you enable auto-provisioning using SCIM. See [Just-in-Time provisioning](/security/for-admins/provisioning/just-in-time/).

### Is IdP-initiated sign-in available?

Docker SSO doesn't support IdP-initiated sign-in, only Service Provider Initiated (SP-initiated) sign-in.

### Is it possible to connect Docker Hub directly with a Microsoft Entra (formerly Azure AD) group?

Yes, Entra ID (formerly Azure AD) is supported with SSO for Docker Business, both through a direct integration and through SAML.

### My SSO connection with Entra ID isn't working and I receive an error that the application is misconfigured. How can I troubleshoot this?

Confirm that you've configured the necessary API permissions in Entra ID (formerly Azure AD) for your SSO connection. You need to grant admin consent within your Entra ID (formerly Azure AD) tenant. See [Entra ID (formerly Azure AD) documentation](https://learn.microsoft.com/en-us/azure/active-directory/manage-apps/grant-admin-consent?pivots=portal#grant-admin-consent-in-app-registrations).
