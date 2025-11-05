---
title: SSO domain FAQs
linkTitle: Domains
description: Frequently asked questions about domain verification and management for Docker single sign-on
keywords: SSO domains, domain verification, DNS, TXT records, single sign-on
tags: [FAQ]
aliases:
- /single-sign-on/domain-faqs/
- /faq/security/single-sign-on/domain-faqs/
- /security/faqs/single-sign-on/domain-faqs/
---

## Can I add sub-domains?

Yes, you can add sub-domains to your SSO connection. All email addresses must use domains you've added to the connection. Verify that your DNS provider supports multiple TXT records for the same domain.

## Do I need to keep the DNS TXT record permanently?

You can remove the TXT record after one-time verification to add the domain. However, if your organization changes identity providers and needs to set up SSO again, you'll need to verify the domain again.

## Can I verify the same domain for multiple organizations?

You can't verify the same domain for multiple organizations at the organization level. To verify one domain for multiple organizations, you must have a Docker Business subscription and create a company. Companies allow centralized management of organizations and domain verification at the company level.
