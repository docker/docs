---
title: Manage an account
weight: 20
description: Learn how to manage settings for your Docker account.
keywords: accounts, docker ID, account settings, account management, docker home
---

You can centrally manage your Docker account settings using Docker Home. Here
you can also take administrative actions for your account and manage your
account security.

> [!TIP]
>
> If your account is associated with an organization that enforces single
> sign-on (SSO), you may not have permissions to update your account settings.
> You must contact your administrator to update your settings.

## Update general settings

1. Sign in to your [Docker account](https://app.docker.com/login).
2. Select your avatar in the top-right corner and select **Account settings**.

From the Account settings page, you can take any of the following actions.

### Update account information

Account information is visible on your account profile in Docker Hub. You can
update the following account information:

- Full name
- Company
- Location
- Website
- Gravatar email: To add an avatar to your Docker account, create a
[Gravatar account](https://gravatar.com/) and create your avatar. Next, add your
Gravatar email to your Docker account settings. It may take some time for your
avatar to update in Docker.

Make your changes here, then select **Save** to save your settings.

### Update email address

To update your email address, select **Email**:

1. Enter your new email address.
2. Enter your password to confirm the change.
3. Select **Send verification email** to send a verification email to your new
email address.

Your new email address will appear as unverified until you complete the
verification process. You can:

- Resend the verification email if needed.
- Removed the unverified email address at any time before verification.

To verify your email, open your email client and follow the instructions
in the Docker verification email.

> [!NOTE]
>
> Docker accounts only support one verified email address at a time, which
is used for account notifications and security-related communications. You
can't add multiple verified email addresses to your account.

### Change your password

You can change your password by initiating a password reset via email.

To change your password, select **Password** and then **Reset password**.
Follow the instructions in the password reset email.

## Manage security settings

To update your two-factor authentication (2FA) settings, select **2FA**.
For information on two-factor authentication (2FA) for your account, see
[Enable two-factor authentication](../security/2fa/_index.md)
to get started.

To manage personal access tokens, select **Personal access tokens**.
For information on personal access tokens, see
[Create and manage access tokens](../security/access-tokens.md).

## Manage connected accounts

You can unlink Google or GitHub accounts that are linked to your Docker account
using the Account settings page:

1. Select **Connected accounts**.
2. Select **Disconnect** on your connected account.
3. To fully unlink your Docker account, you must also unlink Docker from Google
or GitHub. See Google or GitHub's documentation for more information:
    - [Manage connections between your Google Account and third-parties](https://support.google.com/accounts/answer/13533235?hl=en)
    - [Reviewing and revoking authorization of GitHub Apps](https://docs.github.com/en/apps/using-github-apps/reviewing-and-revoking-authorization-of-github-apps)

## Account management

To convert your account into an organization, select **Convert**.
For more information on converting your account, see
[Convert an account into an organization](../admin/organization/convert-account.md).

To deactivate your account, select **Deactivate**.
For information on deactivating your account, see
[Deactivating a user account](./deactivate-user-account.md).
