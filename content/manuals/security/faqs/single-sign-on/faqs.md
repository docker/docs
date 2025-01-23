---
description: Single sign-on FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, administration, security
title: General FAQs on SSO
linkTitle: General
weight: 10
tags: [FAQ]
aliases:
- /single-sign-on/faqs/
- /faq/security/single-sign-on/faqs/
- /single-sign-on/saml-faqs/
- /faq/security/single-sign-on/saml-faqs/
- /security/faqs/single-sign-on/saml-faqs/
---

### Is Docker SSO available for all paid subscriptions?

Docker single sign-on (SSO) is only available with the Docker Business subscription. [Upgrade your existing subscription](../../../subscription/change.md) to start using Docker SSO.

### How does Docker SSO work?

Docker SSO lets users authenticate using their identity providers (IdPs) to access Docker. Docker supports Entra ID (formerly Azure AD) and any SAML 2.0 identity providers. When you enable SSO, this redirects users to your provider’s authentication page to authenticate using their email and password.

### What SSO flows does Docker support?

Docker supports Service Provider Initiated (SP-initiated) SSO flow. This means users must sign in to Docker Hub or Docker Desktop to initiate the SSO authentication process.

### Where can I find detailed instructions on how to configure Docker SSO?

You first need to establish an SSO connection with your identity provider, and the company email domain needs to be verified prior to establishing an SSO connection for your users. For detailed step-by-step instructions on how to configure Docker SSO, see [Single Sign-on](../../../security/for-admins/single-sign-on/configure/_index.md).

### Does Docker SSO support multi-factor authentication (MFA)?

When an organization uses SSO, MFA is determined on the IdP level, not on the Docker platform.

### Do I need a specific version of Docker Desktop for SSO?

Yes, all users in your organization must upgrade to Docker Desktop version 4.4.2 or later. Users on older versions of Docker Desktop won't be able to sign in after SSO is enforced if the company domain email is used to sign in or as the primary email associated with an existing Docker account. Your users with existing accounts can't sign in with their username and password.

### Can I retain my Docker ID when using SSO?

For a personal Docker ID, a user is the account owner. A Docker ID is associated with access to the user's repositories, images, assets. A user can choose to have a company domain email on the Docker account. When enforcing SSO, the account is connected to the organization account. When enforcing SSO for an organization(s) or company, any user logging in without an existing account using verified company domain email will automatically have an account provisioned, and a new Docker ID created.

### Does SAML authentication require additional attributes?

You must provide an email address as an attribute to authenticate through SAML. The ‘Name’ attribute is optional.

### Does the application recognize the NameID/Unique Identifier in the `SAMLResponse` subject?

The preferred format is your email address, which should also be your Name ID.

### Can I use group mapping with SSO and the Azure AD (OIDC) authentication method?

No. Group mapping with SSO isn't supported with the Azure AD (OIDC)
authentication method because it requires granting the OIDC app the
Directory.Read.All permission, which provides access to all users, groups, and
other sensitive data in the directory. Due to potential security risks, Docker
doesn't support this configuration. Instead, Docker recommends [configuring SCIM
to enable group sync
securely](/security/for-admins/provisioning/group-mapping/#use-group-mapping-with-scim).

### Are there any firewall rules required for SSO configuration?

No. There are no specific firewall rules required for configuring SSO, as long as the domain `login.docker.com` is accessible. This domain is commonly accessible by default. However, in rare cases, some organizations may have firewall restrictions in place that block this domain. If you encounter issues during SSO setup, ensure that `login.docker.com` is allowed in your network's firewall settings.