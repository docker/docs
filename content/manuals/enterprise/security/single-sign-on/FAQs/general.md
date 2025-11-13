---
description: Frequently asked questions about Docker single sign-on
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, administration, security
title: General SSO FAQs
linkTitle: General
weight: 10
tags: [FAQ]
aliases:
- /single-sign-on/faqs/
- /faq/security/single-sign-on/faqs/
- /single-sign-on/saml-faqs/
- /faq/security/single-sign-on/saml-faqs/
- /security/faqs/single-sign-on/saml-faqs/
- /security/faqs/single-sign-on/faqs/
---

## What SSO flows does Docker support?

Docker supports Service Provider Initiated (SP-initiated) SSO flow. Users must sign in to Docker Hub or Docker Desktop to initiate the SSO authentication process.

## Does Docker SSO support multi-factor authentication?

When an organization uses SSO, multi-factor authentication is controlled at the identity provider level, not on the Docker platform.

## Can I retain my Docker ID when using SSO?

Users with personal Docker IDs retain ownership of their repositories, images, and assets. When SSO is enforced, existing accounts with company domain emails are connected to the organization. Users signing in without existing accounts automatically have new accounts and Docker IDs created.

## Are there any firewall rules required for SSO configuration?

No specific firewall rules are required as long as `login.docker.com` is accessible. This domain is commonly accessible by default, but some organizations may need to allow it in their firewall settings if SSO setup encounters issues.

## Does Docker use my IdP's default session timeout?

Yes, Docker supports your IdP's session timeout using a custom `dockerSessionMinutes` SAML attribute instead of the standard `SessionNotOnOrAfter` element. See [SSO attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes) for more information.
