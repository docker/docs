---
description: Frequently asked questions about Docker single sign-on, identity providers, and user management
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, identity providers, IdP, SAML, Entra ID, user management, SCIM, JIT, administration, security
title: SSO FAQs
linkTitle: SSO
weight: 40
tags: [FAQ]
aliases:
- /single-sign-on/faqs/
- /faq/security/single-sign-on/faqs/
- /single-sign-on/saml-faqs/
- /faq/security/single-sign-on/saml-faqs/
- /security/faqs/single-sign-on/saml-faqs/
- /security/faqs/single-sign-on/faqs/
- /platform/security/authentication/single-sign-on/FAQs/general/
- /single-sign-on/idp-faqs/
- /faq/security/single-sign-on/idp-faqs/
- /security/faqs/single-sign-on/idp-faqs/
- /platform/security/authentication/single-sign-on/FAQs/idp-faqs/
- /platform/security/faqs/idp-faqs/
- /single-sign-on/users-faqs/
- /faq/security/single-sign-on/users-faqs/
- /security/faqs/single-sign-on/users-faqs/
- /platform/security/authentication/single-sign-on/FAQs/users-faqs/
- /platform/security/faqs/users-faqs/
- /platform/security/faqs/sso-faqs/
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

Yes, Docker supports your IdP's session timeout using a custom `dockerSessionMinutes` SAML attribute instead of the standard `SessionNotOnOrAfter` element. See [SSO attributes](/manuals/platform/security/provisioning/_index.md#sso-attributes) for more information.

## Can I use multiple identity providers with Docker SSO?

Yes, Docker supports multiple IdP configurations. A domain can be associated with multiple IdPs. Docker supports Entra ID (formerly Azure AD) and identity providers that support SAML 2.0.

## Can I change my identity provider after configuring SSO?

Yes. Delete your existing IdP configuration in your Docker SSO connection, then [configure SSO using your new IdP](/manuals/platform/security/authentication/single-sign-on/connect.md). If you had already turned on enforcement, turn off enforcement before updating the provider connection.

## What information do I need from my identity provider to configure SSO?

To turn on SSO in Docker, you need the following from your IdP:

- SAML: Entity ID, ACS URL, Single Logout URL, and the public X.509 certificate
- Entra ID (formerly Azure AD): Client ID, Client Secret, AD Domain

## What happens if my existing certificate expires?

Contact your identity provider to retrieve a new X.509 certificate. Update with the new certificate in [SSO configuration settings](/manuals/platform/security/authentication/single-sign-on/manage.md#manage-sso-connections) from Docker Home.

- If your organization enforces SSO, username and password credentials won't work.
- If your organization doesn't enforce SSO, users can sign in with their username and password credentials.

If you need additional help, contact [Docker support](https://app.docker.com/support/contact).

## What happens if my IdP goes down when SSO is turned on?

If SSO is enforced, users can't access Docker Hub when your IdP is down. Users can still access Docker Hub images from the CLI using personal access tokens.

If SSO is turned on but not enforced, users can fall back to username/password authentication.

## Do bot accounts need seats to access organizations using SSO?

Yes, bot accounts need seats like regular users, requiring a non-aliased domain email in the IdP and using a seat in Docker Hub. You can add bot accounts to your IdP and create access tokens to replace other credentials.

## Does SAML SSO use Just-in-Time provisioning?

The SSO implementation uses Just-in-Time (JIT) provisioning by default. You can optionally turn off JIT in Docker Home if you turn on auto-provisioning using SCIM. See [Just-in-Time provisioning](/manuals/platform/security/provisioning/just-in-time.md).

## My Entra ID SSO connection isn't working and shows an error. How can I troubleshoot this?

Confirm that you've configured the necessary API permissions in Entra ID for your SSO connection. You need to grant administrator consent within your Entra ID tenant. See [Entra ID (formerly Azure AD) documentation](https://learn.microsoft.com/en-us/azure/active-directory/manage-apps/grant-admin-consent?pivots=portal#grant-admin-consent-in-app-registrations).

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

For detailed instructions, see [Configure single sign-on](/manuals/platform/security/authentication/single-sign-on/connect.md).

## Is Docker SSO fully synced with the IdP?

Docker SSO provides Just-in-Time (JIT) provisioning by default. Users are provisioned when they authenticate with SSO. If users leave the organization, administrators must manually [remove the user](/manuals/admin/organization/manage/members.md#remove-members-from-teams) from the organization.

[SCIM](/manuals/platform/security/provisioning/scim/_index.md) provides full synchronization with users and groups. When using SCIM, the recommended configuration is to turn off JIT so all auto-provisioning is handled by SCIM.

Additionally, you can use the [Docker Hub API](/reference/api/hub/latest.md) to complete this process.

## How does turning off Just-in-Time provisioning affect user sign-in?

When JIT is turned off (available with SCIM in Docker Home), users must be organization members or have pending invitations to access Docker. Users who don't meet these criteria get an "Access denied" error and need administrator invitations.

See [SSO authentication with JIT provisioning disabled](/manuals/platform/security/provisioning/just-in-time.md#sso-authentication-with-jit-provisioning-disabled).

## Can someone join an organization without an invitation?

Not without SSO. Joining requires an invite from an organization owner. When SSO is enforced, users with verified domain emails can automatically join the organization when they sign in.

## What happens to existing licensed users when SCIM is turned on?

Turning on SCIM doesn't immediately remove or modify existing licensed users. They retain current access and roles, but you'll manage them through your IdP after SCIM is active. If SCIM is later turned off, previously SCIM-managed users remain in Docker but are no longer automatically updated based on your IdP.

## Is user information visible in Docker Hub?

All Docker accounts have public profiles associated with their namespace. If you don't want user information (like full names) to be visible, remove those attributes from your SSO and SCIM mappings, or use different identifiers to replace users' full names.
