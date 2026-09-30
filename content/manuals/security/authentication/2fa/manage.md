---
title: Manage two-factor authentication for your Docker account
linkTitle: Manage
description: >-
  Turn on two-factor authentication for your Docker account, save the
  recovery code, move 2FA to a new device, or turn 2FA off.
keywords: enable 2FA, disable 2FA, turn on 2FA, turn off 2FA, two-factor
  authentication, Docker account, TOTP, authenticator app, QR code,
  recovery code, new device, personal access token, Docker Hub
weight: 10
aliases:
  - /security/2fa/disable-2fa/
  - /security/for-developers/2fa/disable-2fa/
---

{{< summary-bar feature_name="2FA" >}}

Turn on two-factor authentication (2FA) to require a code from your
authenticator app when you sign in with your password. Turn it off to
sign in with your password alone, or to move 2FA to a new device. For
how 2FA works and what the recovery code does, see
[Two-factor authentication][overview].

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
Docker asks for a code from your authenticator app. Docker also emails
you a reminder to save your recovery code.

> [!IMPORTANT]
>
> The recovery code works once. Using it signs you in and turns 2FA off.
> Keep it somewhere safe. If you lose both your authenticator app and your
> recovery code, contact Docker Support to recover your account.

To sign in with `docker login -u`, or from scripts and CI, use a
[personal access token][pat] in place of your password.

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

Two-factor authentication is off. Docker emails you to confirm the
change.

## Move 2FA to a new device

Docker keeps one authenticator per account. To move 2FA to a new phone or
device, [turn 2FA off](#disable-two-factor-authentication), then
[turn it on again](#enable-two-factor-authentication) from the new
device. Setup is not available while 2FA is on.

## Next steps

- [Recover your account][recover] if you lose your authenticator app or
  recovery code.
- Create a [personal access token][pat] for the Docker CLI and automation.

[overview]: /manuals/security/authentication/2fa/_index.md
[pat]: /manuals/security/access-tokens/personal-access-tokens.md
[recover]: /manuals/security/authentication/2fa/recover-hub-account.md
