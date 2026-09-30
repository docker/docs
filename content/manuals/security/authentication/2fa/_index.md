---
title: Enable two-factor authentication for your Docker account
linkTitle: Two-factor authentication
description: >-
  Enable or disable two-factor authentication on your Docker account for
  enhanced security and account protection.
keywords: two-factor authentication, 2FA, docker hub security,
  account security, TOTP, authenticator app, disable 2FA, recovery code
weight: 20
aliases:
  - /docker-hub/2fa/
  - /security/2fa/disable-2fa/
  - /security/for-developers/2fa/
  - /security/for-developers/2fa/disable-2fa/
  - /security/2fa/
---

{{< summary-bar feature_name="2FA" >}}

Two-factor authentication (2FA) adds a second step to password sign-in. After
you enter your password, Docker asks for a code from an authenticator app. If
someone learns your password, they still can't sign in without the code.

When you turn on 2FA, Docker gives you a recovery code for your account. The
recovery code is how you get back in if you lose your authenticator app, so
store it somewhere safe.

> [!IMPORTANT]
>
> The recovery code works once. Using it signs you in and turns 2FA off. If
> you lose both your authenticator app and your recovery code, you need to
> contact Docker Support to recover your account.

## Key benefits

- Protection against stolen passwords: An attacker who has your password
  still needs the code from your authenticator app.
- Secure CLI access: When 2FA is on, the Docker CLI needs a personal access
  token instead of your password, so scripts and CI never hold your password.
- Compliance: Many organizations require 2FA for access to development and
  production resources.

## Prerequisites

Before you turn on 2FA, you need:

- A time-based one-time password (TOTP) authenticator app on your phone or
  another device
- Your Docker account password
- A verified email address on your account. If your email isn't verified,
  Docker asks you to verify it before you can open the 2FA settings.

> [!NOTE]
>
> If your organization enforces single sign-on (SSO), the **2FA** page shows
> a message to contact your administrator instead. Your identity provider
> manages authentication for your account.

## Enable two-factor authentication

To turn on 2FA for your Docker account:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar and then from the drop-down menu, select **Account
   settings**.
1. Select **2FA**.
1. Enter your account password, then select **Confirm**.
1. Save your recovery code. Select **Copy**, or open the menu next to it to
   **Download** or **Print** the code. Store it somewhere safe.
1. Open your authenticator app and either scan the code on the **QR Code**
   tab or enter the code from the **Text Code** tab.
1. Enter the six-digit code from your authenticator app in the
   **Authentication code** field.
1. Select **Enable 2FA**.

Two-factor authentication is on. From now on, when you sign in with your
password, Docker asks for a code from your authenticator app.

To sign in from the Docker CLI, use a
[personal access token](/manuals/security/access-tokens/personal-access-tokens.md)
in place of your password.

## Disable two-factor authentication

> [!WARNING]
>
> Turning off 2FA leaves your account protected by your password alone.

To switch to a new device, turn off 2FA and then turn it on again from the
new device.

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar and then from the drop-down menu, select **Account
   settings**.
1. Select **2FA**.
1. Enter your password, then select **Confirm**.
1. Select **Disable 2FA**.

## Next steps

- [Recover your account](/manuals/security/authentication/2fa/recover-hub-account.md)
  if you lose your authenticator app or recovery code.
- [Create a personal access token](/manuals/security/access-tokens/personal-access-tokens.md)
  for CLI sign-in and automation.
