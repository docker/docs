---
title: Enable two-factor authentication for your Docker account
linkTitle: Two-factor authentication
description: >-
  Turn on two-factor authentication for your Docker account, save the
  recovery code, or turn 2FA off.
keywords: two-factor authentication, 2FA, Docker account, TOTP,
  authenticator app, recovery code, QR code, personal access token,
  disable 2FA, Docker Hub
weight: 20
aliases:
  - /docker-hub/2fa/
  - /security/2fa/disable-2fa/
  - /security/for-developers/2fa/
  - /security/for-developers/2fa/disable-2fa/
  - /security/2fa/
---

{{< summary-bar feature_name="2FA" >}}

Two-factor authentication (2FA) adds a code from an authenticator app after
you sign in with your password. Someone who knows your password still needs
that code to sign in.

When you turn on 2FA, Docker gives you a recovery code. Keep it somewhere
safe. You use it to get back in if you lose your authenticator app.

> [!IMPORTANT]
>
> The recovery code works once. Using it signs you in and turns 2FA off.
> If you lose both your authenticator app and your recovery code, contact
> Docker Support to recover your account.

## Prerequisites

Before you turn on 2FA, you need:

- A time-based one-time password (TOTP) authenticator app on your phone or
  another device
- Your Docker account password
- A verified email address on your account

Docker opens the 2FA settings after your email address is verified.

> [!NOTE]
>
> If your organization enforces single sign-on (SSO), the **2FA** page
> tells you to contact your administrator. Your identity provider manages
> sign-in for your account.

## Enable two-factor authentication

To turn on 2FA for your Docker account:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner, then select **Account
   settings**.
1. Select **2FA**.
1. Enter your account password, then select **Confirm**.
1. Save your recovery code. Select **Copy**, or open the menu next to
   **Copy** and select **Download** or **Print**.
1. Open your authenticator app. Scan the code on the **QR Code** tab, or
   enter the code from the **Text Code** tab.
1. Enter the six-digit code from your authenticator app in
   **Authentication code**.
1. Select **Enable 2FA**.

Two-factor authentication is on. When you sign in with your password,
Docker asks for a code from your authenticator app.

To sign in from the Docker CLI, use a [personal access token][pat] in
place of your password.

## Disable two-factor authentication

> [!WARNING]
>
> Turning off 2FA leaves your account protected by your password alone.

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar in the top-right corner, then select **Account
   settings**.
1. Select **2FA**.
1. Enter your password, then select **Confirm**.
1. Select **Disable 2FA**.

To move 2FA to a new device, turn it off, then turn it on again from that
device.

## Next steps

- [Recover your account][recover] if you lose your authenticator app or
  recovery code.
- Create a [personal access token][pat] for the Docker CLI and automation.

[pat]: /manuals/security/access-tokens/personal-access-tokens.md
[recover]: /manuals/security/authentication/2fa/recover-hub-account.md
