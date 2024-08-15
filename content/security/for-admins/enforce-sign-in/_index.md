---
description: Understand what happens when you force users to sign in to Docker Desktop
toc_max: 2
keywords: authentication, registry.json, configure, enforce sign-in, docker desktop, security,
title: Enforce sign-in for Docker Desktop
aliases:
- /security/for-admins/configure-sign-in/
---

By default, members of your organization can use Docker Desktop without signing
in. When users don’t sign in as a member of your organization, they don’t
receive the [benefits of your organization’s
subscription](../../../subscription/core-subscription/details.md) and they can circumvent [Docker’s
security features](/security/for-admins/hardened-desktop/_index.md) for your organization.

There are multiple ways you can enforce sign-in, depending on your companies' set up and preferences:
- [Registry key method (Windows only)](methods.md#registry-key-method-windows-only){{< badge color=violet text="Early Access" >}}
- [`.plist` method (Mac only)](methods.md#plist-method-mac-only){{< badge color=violet text="Early Access" >}}
- [`registry.json` method (All)](methods.md#registryjson-method-all)

## How is sign-in enforced?

When Docker Desktop starts and it detects a registry key, a `.plist` file or `registry.json` file, the
following occurs:

- A **Sign in required!** prompt appears requiring the user to sign
  in as a member of your organization to use Docker Desktop. ![Enforce Sign-in
  Prompt](../../images/enforce-sign-in.png?w=400)
- When a user signs in to an account that isn’t a member of your organization,
  they are automatically signed out and can’t use Docker Desktop. The user
  can select **Sign in** and try again.
- When a user signs in to an account that is a member of your organization, they
 can use Docker Desktop.
- When a user signs out, the **Sign in required!** prompt appears and they can
  no longer use Docker Desktop.

> **Enforce sign-in versus enforce SSO**
>
> Enforcing sign-in ensures that users are required to sign in to use Docker Desktop.
> If your organization is also using single sign-on (SSO), you can optionally enforce SSO.
> This means that your users must use SSO to sign in, instead of a username and password.
> When you enforce sign-in and enforce SSO, your users must sign in and must use SSO to do so.
> See [Enforce SSO](/security/for-admins/single-sign-on/connect#optional-enforce-sso) for details on how to enable this for your SSO connection.
