---
title: Two-factor authentication for your Docker account
linkTitle: Two-factor authentication
description: >-
  Learn how two-factor authentication protects a Docker account, when Docker
  asks for the code, and what the recovery code does.
keywords: two-factor authentication, 2FA, how 2FA works, Docker account,
  TOTP, authenticator app, authentication code, recovery code, personal
  access token, docker login, account security, Docker Hub
weight: 20
aliases:
  - /docker-hub/2fa/
  - /security/for-developers/2fa/
  - /security/2fa/
grid:
  - title: Turn 2FA on or off
    description: >-
      Set up an authenticator app, save the recovery code, or turn 2FA off.
    icon: device-phone-mobile
    link: /security/authentication/2fa/manage/
  - title: Recover your account
    description: >-
      Sign in with a recovery code, generate a new one, or contact Support.
    icon: key
    link: /security/authentication/2fa/recover-hub-account/
---

{{< summary-bar feature_name="2FA" >}}

Two-factor authentication (2FA) adds a code from an authenticator app to
your password when you sign in to your Docker account. Someone who knows
your password still needs the code from your device to sign in.

2FA is a setting on individual Docker accounts. You turn it on and off
yourself in **Account settings**. Organization and company settings do not
include 2FA. To control how members sign in across an organization, use
[single sign-on](/manuals/security/authentication/single-sign-on/_index.md).

## How two-factor authentication works

When you turn on 2FA, you pair a time-based one-time password (TOTP)
authenticator app with your account by scanning a QR code or entering a
text code. Any authenticator app that supports TOTP works. The app then
shows a six-digit code that changes every 30 seconds.

After you enter your password, Docker shows the **Two-Factor
Authentication** page. Enter the code from your app in
**Authentication Code** and select **Verify**. If you no longer have the
device, select **I've lost my authentication device** to sign in with
your recovery code instead.

A wrong code fails with `Invalid one-time password`. After repeated
failures, Docker returns `Too many failed login attempts` and blocks
further attempts for a short time.

Docker keeps one authenticator per account. To move 2FA to a new device,
turn it off and turn it on again from the new device.

## When Docker asks for the code

| Sign-in | What Docker asks for |
| --- | --- |
| Browser sign-in to Docker Home or Docker Hub | Your password, then the code from your authenticator app |
| `docker login` with no username | The same browser sign-in, if the browser is not already signed in |
| `docker login -u`, scripts, and CI | A [personal access token](/manuals/security/access-tokens/personal-access-tokens.md) in the password prompt. Password sign-in from the CLI is not supported when 2FA is on |
| Enforced single sign-on | No code. Your identity provider handles sign-in, and the **2FA** page tells you to contact your administrator |

## Recovery code

Docker gives you one recovery code when you turn on 2FA. The code signs
you in if you lose your authenticator app, so copy, download, or print it
and store it somewhere safe.

> [!IMPORTANT]
>
> The recovery code works once. Using it on the **Lost Authentication
> Device** page signs you in, turns 2FA off, and deletes the code. Turn
> 2FA on again from your new device as soon as you are signed in.

- **Generate new code** replaces the previous code. Only the latest code
  works.
- Docker asks for your password before it shows the code or lets you
  change 2FA settings.
- If you lose both the authenticator app and the recovery code, contact
  Docker Support to recover your account.

Docker emails the verified address on your account when you turn 2FA on
or off, when a recovery code is generated, and when a recovery code is
used to sign in. The email does not contain the code.

## Next steps

{{< grid >}}
