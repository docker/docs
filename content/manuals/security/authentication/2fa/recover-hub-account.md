---
title: Recover your Docker account
linkTitle: Recover your account
description: >-
  Recover your Docker account and manage two-factor authentication recovery
  codes.
keywords: account recovery, two-factor authentication, 2FA, recovery code,
  docker hub security, lost authenticator app, 2FA lockout
aliases:
  - /docker-hub/2fa/recover-hub-account/
  - /security/for-developers/2fa/recover-hub-account/
  - /security/2fa/new-recovery-code/
  - /security/2fa/recover-hub-account/
weight: 20
---

{{< summary-bar feature_name="2FA" >}}

This page explains how to get back into your Docker account when part of
your two-factor authentication (2FA) setup is missing. What you do depends
on what you still have:

- Lost your recovery code but can still sign in:
  [Generate a new recovery code](#generate-a-new-recovery-code).
- Lost your authenticator app but have your recovery code:
  [Sign in with your recovery code](#sign-in-with-your-recovery-code).
- Lost both:
  [Contact Docker Support](#recover-your-account-without-access).

## Generate a new recovery code

If you lost your recovery code but can still sign in, generate a new one.
The new code replaces the old one, which stops working.

1. Sign in to your [Docker account](https://app.docker.com/login). Enter
   your password, then the code from your authenticator app.
1. Select your avatar and from the drop-down menu, select **Account
   settings**.
1. Select **2FA**.
1. Enter your password, then select **Confirm**.
1. Select **Generate new code**.

Select the visibility icon to view the new code, then **Copy**, **Download**,
or **Print** it. Store it somewhere safe.

## Sign in with your recovery code

If you lost your authenticator app but still have your recovery code, use
the code to sign in.

> [!IMPORTANT]
>
> The recovery code works once. Using it signs you in and turns 2FA off. Turn
> 2FA on again from your new device as soon as you're signed in.

1. Sign in to your [Docker account](https://app.docker.com/login) with your
   username and password.
1. On the **Two-Factor Authentication** page, select **I've lost my
   authentication device**.
1. Enter your recovery code, then select **Verify**.

You're signed in and 2FA is off. To protect your account again, follow
[Enable two-factor authentication](/manuals/security/authentication/2fa/_index.md#enable-two-factor-authentication).

## Recover your account without access

If you lost both your authenticator app and your recovery code, you can't
complete sign-in on your own.

Open the
[Contact Support form](https://hub.docker.com/support/contact/?category=2fa-lockout).
The form is prefilled for a 2FA lockout. Enter the email address on your
Docker account and follow the instructions from Docker Support.

## Next steps

- [Enable two-factor authentication](/manuals/security/authentication/2fa/_index.md)
  again after you recover your account.
- [Create a personal access token](/manuals/security/access-tokens/personal-access-tokens.md)
  for CLI sign-in and automation.
