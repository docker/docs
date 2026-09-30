---
title: Recover your Docker account and two-factor recovery code
linkTitle: Recover your account
description: >-
  Sign in with a recovery code, generate a new recovery code, or contact
  Support when you lose your authenticator app.
keywords: account recovery, two-factor authentication, 2FA, recovery code,
  lost authenticator app, 2FA lockout, Docker account, generate recovery
  code
aliases:
  - /docker-hub/2fa/recover-hub-account/
  - /security/for-developers/2fa/recover-hub-account/
  - /security/2fa/new-recovery-code/
  - /security/2fa/recover-hub-account/
weight: 20
---

{{< summary-bar feature_name="2FA" >}}

Get back into your Docker account when part of your two-factor
authentication (2FA) setup is missing. What you do depends on what you
still have:

- You lost your recovery code and can still sign in.
  [Generate a new recovery code](#generate-a-new-recovery-code).
- You lost your authenticator app and still have your recovery code.
  [Sign in with your recovery code](#sign-in-with-your-recovery-code).
- You lost both your authenticator app and your recovery code.
  [Contact Docker Support](#contact-docker-support).

## Generate a new recovery code

If you lost your recovery code and can still sign in, generate a new one.
The new code replaces the previous code.

1. Sign in to your [Docker account](https://app.docker.com/login). Enter
   your password, then the code from your authenticator app.
1. Select your avatar in the top-right corner, then select **Account
   settings**.
1. Select **2FA**.
1. Enter your password, then select **Confirm**.
1. Select **Generate new code**.

Select the visibility icon to view the new code. Then select **Copy**,
**Download**, or **Print**, and store the code somewhere safe.

## Sign in with your recovery code

If you lost your authenticator app and still have your recovery code, use
the code to sign in.

> [!IMPORTANT]
>
> The recovery code works once. Using it signs you in and turns 2FA off.
> Turn 2FA on again from your new device as soon as you're signed in.

1. Sign in to your [Docker account](https://app.docker.com/login) with your
   username and password.
1. On the **Two-Factor Authentication** page, select **I've lost my
   authentication device**.
1. Enter your recovery code, then select **Verify**.

You're signed in and 2FA is off. To protect your account again, follow
[Turn on 2FA][enable].

## Contact Docker Support

If you lose both your authenticator app and your recovery code, contact
Docker Support to restore access.

Open the
[Contact Support](https://hub.docker.com/support/contact/?category=2fa-lockout).
The subject and description already describe a 2FA lockout. Enter the
email address on your Docker account, then follow the instructions from
Docker Support.

## Next steps

- [Turn on 2FA][enable] again after you recover your account.
- Create a [personal access token][pat] for the Docker CLI and automation.

[enable]: /manuals/security/authentication/2fa/manage.md
[pat]: /manuals/security/access-tokens/personal-access-tokens.md
