---
description: Frequently asked questions about Docker security, authentication, and organization management
keywords: Docker security, FAQs, authentication, SSO, vulnerability reporting, session management
title: General security FAQs
linkTitle: General
weight: 10
tags: [FAQ]
aliases:
- /faq/security/general/
---

## How do I report a vulnerability?

If you've discovered a security vulnerability in Docker, report it responsibly to security@docker.com so Docker can quickly address it.

## Does Docker lockout users after failed sign-ins?

Docker Hub locks out users after 10 failed sign-in attempts within 5 minutes. The lockout duration is 5 minutes. This policy applies to Docker Hub, Docker Desktop, and Docker Scout authentication.

## Do you support physical multi-factor authentication (MFA) with YubiKeys?

You can configure physical multi-factor authentication (MFA) through SSO using your identity provider (IdP). Check with your IdP if they support physical MFA devices like YubiKeys.

## How are sessions managed and do they expire?

Docker uses tokens to manage user sessions with different expiration periods:

- Docker Desktop: Signs you out after 90 days, or 30 days of inactivity
- Docker Hub and Docker Home: Sign you out after 24 hours

Docker also supports your IdP's default session timeout through SAML attributes. For more information, see [SSO attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes).

## How does Docker distinguish between employee users and contractor users?

Organizations use verified domains to distinguish user types. Team members with email domains other than verified domains appear as "Guest" users in the organization.

## How long are activity logs available?

Docker activity logs are available for 90 days. You're responsible for exporting logs or setting up drivers to send logs to your internal systems for longer retention.

## Can I export a list of users with their roles and privileges?

Yes, use the [Export Members](../../admin/organization/members.md#export-members) feature to export a CSV file containing your organization's users with role and team information.

## How does Docker Desktop handle authentication information?

Docker Desktop uses the host operating system's secure key management to store authentication tokens:

- macOS: [Keychain](https://support.apple.com/guide/security/keychain-data-protection-secb0694df1a/web)
- Windows: [Security and Identity API via Wincred](https://learn.microsoft.com/en-us/windows/win32/api/wincred/)
- Linux: [Pass](https://www.passwordstore.org/).

## How do I remove users who aren't part of my IdP when using SSO without SCIM?

If SCIM isn't turned on, you must manually remove users from the organization. SCIM can automate user removal, but only for users added after SCIM is turned on. Users added before SCIM was turned on must be removed manually.

For more information, see [Manage organization members](/manuals/admin/organization/members.md).

## What metadata does Scout collect from container images?

For information about metadata stored by Docker Scout, see [Data handling](/manuals/scout/deep-dive/data-handling.md).

## How are Marketplace extensions vetted for security?

Security vetting for extensions is on the roadmap but isn't currently implemented. Extensions aren't covered as part of Docker's Third-Party Risk Management Program.

## Can I prevent users from pushing images to Docker Hub private repositories?

No direct setting exists to disable private repositories. However, [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) lets administrators control which registries developers can access through Docker Desktop via the Admin Console.
