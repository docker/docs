---
title: Manage a Docker account
linkTitle: Manage an account
weight: 20
description: Learn how to manage your Docker account.
keywords: accounts, docker ID, account settings, account management, docker home
---

You can centrally manage your Docker account using Docker Home, including
adminstrative and security settings.

> [!TIP]
>
> If your account is associated with an organization that enforces single
> sign-on (SSO), you may not have permissions to update your account settings.
> You must contact your administrator to update your settings.

## Update account information

Account information is visible on your **Account settings** page. You can
update the following account information:

- Full name
- Company
- Location
- Website
- Gravatar email

To add or update your avatar using Gravatar:

1. Create a [Gravatar account](https://gravatar.com/).
1. Create your avatar.
1. Add your Gravatar email to your Docker account settings.

It may take some time for your avatar to update in Docker.

## Update email address

To update your email address:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner and select **Account settings**.
1. Select **Email**.
1. Enter your new email address and your password to confirm the change.
1. Select **Send verification email**. Docker sends a verification
link to your new email.

Your new email address will appear as unverified until you complete the
verification process. You can:

- Resend the verification email if needed.
- Remove the unverified email address at any time before verification.

To verify your email, open your email client and follow the instructions
in the Docker verification email.

> [!NOTE]
>
> Docker accounts only support one verified email address at a time, which
is used for account notifications and security-related communications. You
can't add multiple verified email addresses to your account.

## Change your password

You can change your password by initiating a password reset via email. To change your password:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner and select **Account settings**.
1. Select **Password**, then **Reset password**.
1. Docker will send you a password reset email with instructions to reset
your password.

## Manage two-factor authentication

To update your two-factor authentication (2FA) settings:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner and select **Account settings**.
1. Select **2FA**.

For more information, see
[Enable two-factor authentication](../security/2fa/_index.md).

## Manage personal access tokens

To manage personal access tokens:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner and select **Account settings**.
1. Select **Personal access tokens**.

For more information, see
[Create and manage access tokens](../security/access-tokens.md).

## Manage connected accounts

You can unlink connected Google or GitHub accounts:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner and select **Account settings**.
1. Select **Connected accounts**.
1. Select **Disconnect** on your connected account.

To fully unlink your Docker account, you must also unlink Docker from Google
or GitHub. See Google or GitHub's documentation for more information:

- [Manage connections between your Google Account and third-parties](https://support.google.com/accounts/answer/13533235?hl=en)
- [Reviewing and revoking authorization of GitHub Apps](https://docs.github.com/en/apps/using-github-apps/reviewing-and-revoking-authorization-of-github-apps)

## Convert your account

For information on converting your account into an organization, see
[Convert an account into an organization](../admin/organization/convert-account.md).

## Deactivate your account

For information on deactivating your account, see
[Deactivating a user account](./deactivate-user-account.md).
