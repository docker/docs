---
description: Single sign-on enforcement FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, enforce SSO, SSO enforcement
title: Enforcement
tags: [FAQ]
aliases:
- /single-sign-on/enforcement-faqs/
- /faq/security/single-sign-on/enforcement-faqs/
---

### We currently have a Docker Team subscription. How do we enable SSO?

SSO is available with a Docker Business subscription. To enable SSO, you must first upgrade your subscription to a Docker Business subscription. To learn how to upgrade your existing account, see [Upgrade your subscription](../../../subscription/core-subscription/upgrade.md).

### How do service accounts work with SSO?

Service accounts work like any other user when SSO is turned on. If the service account is using an email for a domain with SSO turned on, it needs a PAT for CLI and API usage.

### Is DNS verification required to enable SSO?

Yes. You must verify a domain before using it with an SSO connection.

### Does Docker SSO support authenticating through the command line?

When SSO is enforced, you can't use passwords to access the Docker CLI, but you can still access the Docker CLI using a personal access token (PAT) for authentication.

Each user must create a PAT to access the CLI. To learn how to create a PAT, see [Manage access tokens](/security/for-developers/access-tokens/). Users who already used a PAT to sign in before SSO enforcement will still be able to use that PAT to authenticate.

### How does SSO affect our automation systems and CI/CD pipelines?

Before enforcing SSO, you must create PATs for automation systems and CI/CD pipelines and use the tokens instead of a password.

### What can organization users who authenticated with personal emails prior to enforcement expect?

Ensure your users have their organization email on their account, so that the accounts will be migrated to SSO for authentication.

### Can I enable SSO and hold off on the enforcement option?

Yes, you can choose to not enforce, and users have the option to use either Docker ID (standard email and password) or domain-verified email address (SSO) at the sign-in screen.

### SSO is enforced, but one of our users is connected to several organizations (and several email addresses) and is able to bypass SSO and sign in through username and password. Why is this happening?

Users can bypass SSO if the email they're using to sign in doesn't match the organization email that's used for SSO enforcement.

### Is there a way to test this functionality in a test tenant with Okta before going to production?

Yes, you can create a test organization. Companies can set up a new 5 seat Business plan on a new organization to test with (making sure to only enable SSO, not enforce it or all domain email users will be forced to sign in to that test tenant).

### Once we enable SSO for Docker Desktop, what's the impact to the flow for Build systems that use service accounts?

If you enable SSO, there is no impact. Both username/password or personal access token (PAT) sign-in are supported.
However, if you enforce SSO:

- Service Account domain email addresses must not be aliased and must be enabled in their IdP
- Username/password authentication won’t work, so you should update the build system to use a PAT instead of a password
- Those who know the IdP credentials can sign in as that Service Account through SSO on Hub and create or change the personal access token for that service account.

### Is the sign in required tracking at runtime or install time?

At runtime for Docker Desktop if it’s configured to require authentication to the organization.

### What is enforcing SSO versus enforcing sign-in?

Enforcing SSO and enforcing sign-in to Docker Desktop are different features that you can use separately or together.

Enforcing SSO ensures that users sign in using their SSO credentials instead of their Docker ID. One of the benefits is that SSO enables you to better manage user credentials.

Enforcing sign-in to Docker Desktop ensures that users always sign in to an

account that's a member of your organization. The benefits are that your organization's security settings are always applied to the user's session and your users always receive the benefits of your subscription. For more details, see [Enforce sign-in for Desktop](../../../security/for-admins/enforce-sign-in/_index.md).

