---
description: Understand what happens when you force users to sign in to Docker Desktop
toc_max: 2
keywords: authentication, registry.json, configure, enforce sign-in, docker desktop, security, .plist, registry key, mac, windows
title: Enforce sign-in for Docker Desktop
linkTitle: Enforce sign-in
tags: [admin]
aliases:
 - /security/for-admins/configure-sign-in/
 - /docker-hub/configure-sign-in/
weight: 30
---

{{< summary-bar feature_name="Enforce sign-in" >}}

By default, members of your organization can use Docker Desktop without signing
in. When users don’t sign in as a member of your organization, they don’t
receive the [benefits of your organization’s
subscription](../../../subscription/details.md) and they can circumvent [Docker’s
security features](/manuals/security/for-admins/hardened-desktop/_index.md) for your organization.

There are multiple methods for enforcing sign-in, depending on your companies' set up and preferences:
- [Registry key method (Windows only)](methods.md#registry-key-method-windows-only){{< badge color=green text="New" >}}
- [Configuration profiles method (Mac only)](methods.md#configuration-profiles-method-mac-only){{< badge color=green text="New" >}}
- [`.plist` method (Mac only)](methods.md#plist-method-mac-only){{< badge color=green text="New" >}}
- [`registry.json` method (All)](methods.md#registryjson-method-all)

## How is sign-in enforced?

When Docker Desktop starts and it detects a registry key, `.plist` file, or `registry.json` file, the following occurs:

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

> [!NOTE]
>
> Enforcing sign-in for Docker Desktop does not impact accessing the Docker CLI. CLI access is only impacted for organizations that enforce single sign-on.

## Enforcing sign-in versus enforcing single sign-on (SSO)

[Enforcing SSO](/manuals/security/for-admins/single-sign-on/connect.md#optional-enforce-sso) and enforcing sign-in are different features. The following table provides a
description and benefits when using each feature.

| Enforcement                       | Description                                                     | Benefits                                                                                                                                                                                                                                                   |
|:----------------------------------|:----------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Enforce sign-in only              | Users must sign in before using Docker Desktop.                 | Ensures users receive the benefits of your subscription and ensures security features are applied. In addition, you gain insights into users’ activity.                                                                                                    |
| Enforce single sign-on (SSO) only | If users sign in, they must sign in using SSO.                  | Centralizes authentication and enforces unified policies set by the identity provider.                                                                                                                                                                     |
| Enforce both                      | Users must sign in using SSO before using Docker Desktop.       | Ensures users receive the benefits of your subscription and ensures security features are applied. In addition, you gain insights into users’ activity. Finally, it centralizes authentication and enforces unified policies set by the identity provider. |
| Enforce neither                   | If users sign in, they can use SSO or their Docker credentials. | Lets users access Docker Desktop without barriers, but at the cost of reduced security and insights.                                                                                                                                                  |

## What's next?

- To enforce sign-in, review the [Methods](/manuals/security/for-admins/enforce-sign-in/methods.md) guide.
- To enforce SSO, review the [Enforce SSO](/manuals/security/for-admins/single-sign-on/connect.md) steps.