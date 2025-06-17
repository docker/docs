---
title: Enforce secure authentication
description: Enforce secure authentication for all users in your organization.
weight: 10
---

In regulated and security-sensitive environments, enforcing single sign-on
(SSO) ensures all users authenticate through a centralized identity provider
(IdP). This strengthens security, simplifies user management, and allows you to
enforce organization-wide authentication policies.

This module walks you through how to configure SSO for your Docker organization,
enforce it for all users, and disable fallback sign-in methods.

## Prerequisites

Before you begin, ensure you have:

- A Docker Business subscription
- Admin access to your Docker organization or company
- Access to your DNS provider
- Access to your Identity Provider (IdP) admin console (e.g., Okta, Azure AD)

## Step one: Add and verify your domain

Verifying your organization’s domain is the first step in securing access. This
process confirms ownership and allows you to enforce SSO and auto-provisioning.

1. Sign in to the [Docker Admin Console](https://app.docker.com/admin) and
select your organization from the **Choose profile** page.
1. Add your domain:
    1. Under **Security and access**, select **Domain management**.
    1. Select **Add a domain**.
    1. Enter your domain (e.g., `example.com`) and select **Add domain**.
1. Verify your domain:
    1. A pop-up modal will display a **TXT Record Value.**
    1. Sign in to your DNS provider and add a TXT record using the provided value.
    1. It may take up to 72 hours for DNS changes to propagate.
    1. Once the TXT record is recognized, return to the Admin Console’s **Domain management** page and select **Verify**.

> [!NOTE]
>
> For detailed instructions on adding TXT records with specific DNS providers,
see [Domain management]().

## Step two: Set up SSO

Docker offers two types of SSO integration:

- OIDC: For IdPs like Entra ID, Auth0, or Google Workspace
- SAML 2.0: Widely supported by enterprise IdPs like Okta, Ping, and legacy
providers

Docker’s SSO configuration supports:

- Just-in-Time (JIT) user provisioning
- Multi-domain SSO
- Group mapping for team assignment (covered in [Module 2]())

To compare protocols and choose your setup path, start with the
[SSO overview]().

Then follow the instructions for your IdP:

- [Set up OIDC SSO]()
- [Set up SAML SSO]()

Each guide walks you through:

- Linking your verified domain to your IdP
- Entering credentials
- Mapping user claims
- Testing the connection with a non-admin account

## Step three: Enforce SSO

Once you’ve confirmed the SSO connection works, you can enforce it across your
organization to ensure all users authenticate through your IdP.

To enforce SSO:

1. In the [Admin Console](https://app.docker.com/admin), navigate
to **Security and access** > **Authentication**.
2. Under **SSO enforcement**, select **Enforce SSO for all users**.
3. Confirm your changes.

This step blocks users from signing in with Docker credentials and requires
authentication via your IdP for any domain-matched account.

## Step four: Enforce Docker Desktop sign-in

To prevent users from running Docker Desktop anonymously or without
organizational control, you can enforce sign-in at the Desktop client level.
When enabled, users must sign in with a Docker ID to use Docker Desktop.

This setting is enforced using centralized configuration methods like:

- `admin-settings.json` for local testing and smaller rollouts
- Mobile Device Management (MDM) tools for larger fleets

To enable it:

1. In your settings configuration, set:

    ```json
    {
    	"enforceSignIn": true
    }
    ```

2. Distribute the setting using one of the supported configuration
methods (e.g., MDM, file copy, registry edit).

For full details, see [Enforce sign-in]().

## Best practices

- Enable Just-in-Time (JIT) provisioning to streamline user onboarding.
- Set up Multi-Factor Authentication (MFA) in your IdP for stronger
authentication.
- Use Enforce Sign-In on Docker Desktop to prevent unauthenticated or offline
usage.
- Avoid fallback authentication paths by enforcing SSO per domain.
- Test with sample accounts before rolling out enforcement org-wide.
