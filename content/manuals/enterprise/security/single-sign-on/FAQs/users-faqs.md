---
title: SSO user management FAQs
linkTitle: User management
description: Frequently asked questions about managing users with Docker single sign-ons
keywords: SSO user management, user provisioning, SCIM, just-in-time provisioning, organization members
tags: [FAQ]
aliases:
- /single-sign-on/users-faqs/
- /faq/security/single-sign-on/users-faqs/
- /security/faqs/single-sign-on/users-faqs/
---

## Do I need to manually add users to my organization?

No, you don't need to manually add users to your organization. Just ensure user accounts exist in your IdP. When users sign in to Docker with their domain email address, they're automatically added to the organization after successful authentication.

## Can users use different email addresses to authenticate through SSO?

All users must authenticate using the email domain specified during SSO setup. Users with email addresses that don't match the verified domain can sign in as guests with username and password if SSO isn't enforced, but only if they've been invited.

## How will users know they're being added to a Docker organization?

When SSO is turned on, users are prompted to authenticate through SSO the next time they sign in to Docker Hub or Docker Desktop. The system detects their domain email and prompts them to sign in with SSO credentials instead.

For CLI access, users must authenticate using personal access tokens.

## Can I convert existing users from non-SSO to SSO accounts?

Yes, you can convert existing users to SSO accounts. Ensure users have:

- Company domain email addresses and accounts in your IdP
- Docker Desktop version 4.4.2 or later
- Personal access tokens created to replace passwords for CLI access
- CI/CD pipelines updated to use PATs instead of passwords

For detailed instructions, see [Configure single sign-on](/manuals/enterprise/security/single-sign-on/configure.md).

## Is Docker SSO fully synced with the IdP?

Docker SSO provides Just-in-Time (JIT) provisioning by default. Users are provisioned when they authenticate with SSO. If users leave the organization, administrators must manually [remove the user](/manuals/admin/organization/members.md#remove-a-member-or-invitee) from the organization.

[SCIM](/manuals/enterprise/security/provisioning/scim.md) provides full synchronization with users and groups. When using SCIM, the recommended configuration is to turn off JIT so all auto-provisioning is handled by SCIM.

Additionally, you can use the [Docker Hub API](/reference/api/hub/latest/) to complete this process.

## How does turning off Just-in-Time provisioning affect user sign-in?

When JIT is turned off (available with SCIM in the Admin Console), users must be organization members or have pending invitations to access Docker. Users who don't meet these criteria get an "Access denied" error and need administrator invitations.

See [SSO authentication with JIT provisioning disabled](/manuals/enterprise/security/provisioning/just-in-time.md#sso-authentication-with-jit-provisioning-disabled).

## Can someone join an organization without an invitation?

Not without SSO. Joining requires an invite from an organization owner. When SSO is enforced, users with verified domain emails can automatically join the organization when they sign in.

## What happens to existing licensed users when SCIM is turned on?

Turning on SCIM doesn't immediately remove or modify existing licensed users. They retain current access and roles, but you'll manage them through your IdP after SCIM is active. If SCIM is later turned off, previously SCIM-managed users remain in Docker but are no longer automatically updated based on your IdP.

## Is user information visible in Docker Hub?

All Docker accounts have public profiles associated with their namespace. If you don't want user information (like full names) to be visible, remove those attributes from your SSO and SCIM mappings, or use different identifiers to replace users' full names.
