---
description: Single sign-on user management FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: FAQs on SSO and managing users
tags: [FAQ]
aliases:
- /single-sign-on/users-faqs/
- /faq/security/single-sign-on/users-faqs/
---

### How do I manage users when using SSO?

You can manage users through organizations in Docker Hub or Admin Console. When you configure SSO in Docker, you need to make sure an account exists for each user in your IdP account. When a user signs in to Docker for the first time using their domain email address, they will be automatically added to the organization after a successful authentication.

### Do I need to manually add users to my organization?

No, you don’t need to manually add users to your organization in Docker or Admin Console. You just need to make sure an account for your users exists in your IdP. When users sign in to Docker, they're automatically assigned to the organization using their domain email address.

When a user signs in to Docker for the first time using their domain email address, they will be automatically added to the organization after a successful authentication.

### Can users in my organization use different email addresses to authenticate through SSO?

During the SSO setup, you’ll have to specify the company email domains that are allowed to authenticate. All users in your organization must authenticate using the email domain specified during SSO setup. Some of your users may want to maintain a different account for their personal projects.

If SSO isn't enforced, users with an email address that doesn't match the verified email domain can sign in with username and password to join the organization as guests.

### Can Docker organization and company owners approve users to join an organization and use a seat, rather than having them automatically added when SSO is enabled?

Organization owners and company owners can approve users by configuring their permissions through their IdP. If the user account is configured in the IdP, the user will be automatically added to the organization in Docker Hub as long as there’s an available seat.

### How will users be made aware that they're being made a part of a Docker organization?

When SSO is enabled, users will be prompted to authenticate through SSO the next time they try to sign in to Docker Hub or Docker Desktop. The system will see the end-user has a domain email associated with the Docker ID they're trying to authenticate with, and prompts them to sign in with SSO email and credentials instead.

If users attempt to sign in through the CLI, they must authenticate using a personal access token (PAT).

### Is it possible to force users of Docker Desktop to authenticate, and/or authenticate using their company’s domain?

Yes. Admins can [force users to authenticate with Docker Desktop](../../for-admins/enforce-sign-in/_index.md) using a registry key, `.plist` file, or `registry.json` file. 

Once SSO enforcement is set up on their Docker Business organization or company on Hub, when the user is forced to authenticate with Docker Desktop, the SSO enforcement will also force users to authenticate through SSO with their IdP (instead of authenticating using their username and password).

Users may still be able to authenticate as a guest account using an email address that doesn't match the verified domain. However, they can only authenticate as guests if that non-domain email was invited.

### Is it possible to convert existing users from non-SSO to SSO accounts?

Yes, you can convert existing users to an SSO account. To convert users from a non-SSO account:

- Ensure your users have a company domain email address and they have an account in your IdP.
- Verify that all users have Docker Desktop version 4.4.2 or later installed on their machines.
- Each user has created a PAT to replace their passwords to allow them to sign in through Docker CLI.
- Confirm that all CI/CD pipelines automation systems have replaced their passwords with PATs.

For detailed prerequisites and instructions on how to enable SSO, see [Configure Single Sign-on](../../../security/for-admins/single-sign-on/configure/_index.md).

### What impact can users expect once we start onboarding them to SSO accounts?

When SSO is enabled and enforced, your users just have to sign in using the verified domain email address.

### Is Docker SSO fully synced with the IdP?

Docker SSO provides Just-in-Time (JIT) provisioning by default, with an option to disable JIT. Users are provisioned when a user authenticates with SSO. If a user leaves the organization, administrators must sign in to Docker and manually [remove the user](../../../admin/organization/members.md#remove-a-member-or-invitee) from the organization.

[SCIM](../../../security/for-admins/provisioning/scim/) is available to provide full synchronization with users and groups. When you auto-provision users with SCIM, the recommended configuration is to disable JIT so that all auto-provisioning is handled by SCIM.

Additionally, you can use the [Docker Hub API](/reference/api/hub/latest/) to complete this process.

### How does disabling Just-in-Time provisioning impact user sign-in?

The option to disable JIT is available when you use the Admin Console and enable SCIM. If a user attempts to sign in to Docker using an email address that is a verified domain for your SSO connection, they need to be a member of the organization to access it, or have a pending invitation to the organization. Users who don't meet these criteria will encounter an `Access denied` error, and will need an administrator to invite them to the organization.

See [SSO authentication with JIT provisioning disabled](/security/for-admins/provisioning/just-in-time/#sso-authentication-with-jit-provisioning-disabled).

To auto-provision users without JIT provisioning, you can use [SCIM](/security/for-admins/provisioning/scim/).

### What's the best way to provision the Docker subscription without SSO?

Company or organization owners can invite users through Docker Hub or Admin Console, by email address (for any user) or by Docker ID (assuming the user has an existing Docker account).

### Can someone join an organization without an invitation? Is it possible to add specific users to an organization with existing email accounts?

Not without SSO. Joining requires an invite from an organization owner. When SSO is enforced, then the domains verified through SSO will let users automatically join the organization the next time they sign in as a user that has a domain email assigned.

### When we send an invitation to the user, will the existing account be consolidated and retained?

Yes, the existing user account will join the organization with all assets retained.

### How can I view, update, and remove multiple email addresses for my users?

We only support one email per user on the Docker platform.

### How can I remove invitees to the organization who haven't signed in?

You can go to the **Members** page for your organization in Docker Hub or Admin Console, view pending invites, and remove invitees as needed.

### Is the flow for service account authentication different from a UI user account?

No, we don't differentiate the two in product.

### Is user information visible in Docker Hub?

All Docker accounts have a public profile associated with their namespace. If you don't want user information (for example, full name) to be visible, you can remove those attributes from your SSO and SCIM mappings. Alternatively, you can use a different identifier to replace a user's full name.
