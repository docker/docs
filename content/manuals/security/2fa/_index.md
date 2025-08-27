---
title: Enable two-factor authentication for your Docker account
linkTitle: Two-factor authentication
description: Enable or disable two-factor authentication on your Docker account for enhanced security and account protection
keywords: two-factor authentication, 2FA, docker hub security, account security, TOTP, authenticator app, disable 2FA
weight: 20
aliases:
 - /docker-hub/2fa/
 - /security/2fa/disable-2fa/
 - /security/for-developers/2fa/
 - /security/for-developers/2fa/disable-2fa/
---

Two-factor authentication (2FA) adds an essential security layer to your Docker account by requiring a unique security code in addition to your password when signing in. This prevents unauthorized access even if your password is compromised.

When you turn on two-factor authentication, Docker provides a unique recovery code specific to your account. Store this code securely as it lets you recover your account if you lose access to your authenticator app.

## Key benefits

Two-factor authentication significantly improves your account security:

- Protection against password breaches: Even if your password is stolen or leaked, attackers can't access your account without your second factor.
- Secure CLI access: Required for Docker CLI authentication when 2FA is turned on, ensuring automated tools use personal access tokens instead of passwords.
- Compliance requirements: Many organizations require 2FA for accessing development and production resources.
- Peace of mind: Know that your Docker repositories, images, and account settings are protected by industry-standard security practices.

## Prerequisites

Before turning on two-factor authentication, you need:

- A smartphone or device with a Time-based One-time password (TOTP) authenticator app installed
- Access to your Docker account password

## Enable two-factor authentication

To turn on 2FA for your Docker account:

1. Sign in to your [Docker account](https://app.docker.com/login).
1. Select your avatar and then from the drop-down menu, select **Account settings**.
1. Select **2FA**.
1. Enter your account password, then select **Confirm**.
1. Save your recovery code and store it somewhere safe. You can use your recovery code to recover your account in the event you lose access to your authenticator app.
1. Use a TOTP mobile app to scan the QR code or enter the text code.
1. Once you've linked your authenticator app, enter the six-digit code in the text-field.
1. Select **Enable 2FA**.

Two-factor authentication is now active on your account. You'll need to enter a security code from your authenticator app each time you sign in.

## Disable two-factor authentication

> [!WARNING]
>
> Disabling two-factor authentication results in decreased security for your Docker account.

1. Sign in to your [Docker account](https://app.docker.com/login).
2. Select your avatar and then from the drop-down menu, select **Account settings**.
3. Select **2FA**.
4. Enter your password, then select **Confirm**.
5. Select **Disable 2FA**.
