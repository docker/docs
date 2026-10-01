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

Turn two-factor authentication (2FA) on or off for your Docker account
in **Account settings**. For how 2FA works, when Docker asks for the
code, and what the recovery code does, see
[Two-factor authentication][overview].

## Prerequisites

Before you turn on 2FA, you need:

- A time-based one-time password (TOTP) authenticator app on your phone or
  another device
- Your Docker account password
- A verified email address on your account

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

To move 2FA to a new phone or device,
[turn 2FA off](#disable-two-factor-authentication), then
[turn it on again](#enable-two-factor-authentication) from the new
device.

## Next steps

- [Recover your account][recover] if you lose your authenticator app or
  recovery code.
- Create a [personal access token][pat] to sign in from the Docker CLI,
  scripts, and CI.

[overview]: /manuals/security/authentication/2fa/_index.md
[pat]: /manuals/security/access-tokens/personal-access-tokens.md
[recover]: /manuals/security/authentication/2fa/recover-hub-account.md
