---
title: Enforce sign-in for Docker Desktop
linkTitle: Enforce sign-in
description: Require users to sign in to Docker Desktop to access organization benefits and security features
toc_max: 2
keywords: authentication, registry.json, configure, enforce sign-in, docker desktop, security, .plist, registry key, mac, windows, organization
tags: [admin]
aliases:
 - /security/for-admins/configure-sign-in/
 - /security/for-admins/enforce-sign-in/
 - /enterprise/security/enforce-sign-in/
weight: 30
---

{{< summary-bar feature_name="Enforce sign-in" >}}

By default, users can access Docker Desktop without signing in to your organization.
When users don't sign in as organization members, they miss out on subscription benefits and can bypass security features configured for your organization.

You can enforce sign-in using several methods, depending on your setup:

- [Registry key method (Windows only)](methods.md#windows-registry-key-method)
- [Configuration profiles method (Mac only)](methods.md#mac-configuration-profiles-method-recommended)
- [`.plist` method (Mac only)](methods.md#mac-plist-file-method)
- [`registry.json` method (all platforms)](methods.md#all-platforms-registryjson-method)

Deploying a `admin-settings.json` file with [Settings Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/_index.md) also enforces sign-in. See [Settings Management and sign-in enforcement](methods.md#settings-management-and-sign-in-enforcement).

This page provides an overview of how sign-in enforcement works.

## How sign-in enforcement works

When Docker Desktop detects a registry key, configuration profile, `.plist` file, or
`registry.json` file:

- A **Sign in using your work email address** prompt appears, requiring users to
  sign in as organization members to use Docker Desktop. The prompt states which
  organizations are required and which method enforces it.
- If users sign in with accounts that aren't organization members, they're
  automatically signed out and can't use Docker Desktop. The prompt changes to
  **You have been signed out** and explains why. They can sign in again with a
  different account.
- When users sign in with organization member accounts, they can use Docker
  Desktop normally.
- When users sign out, the sign-in prompt reappears and they can
  no longer use Docker Desktop unless they sign back in.

### Impact on the Docker CLI

Sign-in enforcement also blocks the Docker CLI. While the sign-in prompt is
showing, Docker Desktop's API proxy rejects almost every request with an
explanation at the terminal, for example:

```text
Sign in to continue using Docker Desktop. Membership in the [myorg] organization
is required. Sign in enforced by your administrators (via registry.json).
```

- `docker run`, `docker pull`, `docker build`, `docker ps`, and other commands
  that reach the engine fail until the user signs in.
- `docker version`, `docker info`, and `docker login` continue to work, so users
  can sign in from the CLI.

> [!IMPORTANT]
>
> Make sure you plan for blocking the Docker CLI before you roll out enforcement. Any scripted or CI use of the
> Docker CLI on an enforced machine stops working until that machine's user signs
> in as an organization member.

Sign-in enforcement is separate from [SSO enforcement](#enforcing-sign-in-versus-enforcing-single-sign-on-sso), which governs how users authenticate as opposed to whether they must authenticate.

### Impact on already-signed-in users

When enforcement is first deployed, users who are already running Docker Desktop are not immediately affected. Docker Desktop re-evaluates enforcement when it starts, and when a user signs in or out. It doesn't poll for new configuration while running, so a newly deployed registry key, configuration profile, `.plist`, or `registry.json` file takes effect on the next restart.

On the next Docker Desktop restart:

- Users signed in with an organization member account are automatically re-authenticated and continue working uninterrupted.
- Users signed in with a non-member account are immediately signed out on startup and see the **Sign in required!** prompt. 

## Enforcing sign-in versus enforcing single sign-on (SSO)

Enforcing Docker Desktop sign-in and [enforcing SSO](/manuals/security/authentication/single-sign-on/connect.md#enforce-sso) are different features that serve different purposes:

| Enforcement                       | Description                                                     | Benefits                                                                                                                                                                                                                                                   |
|:----------------------------------|:----------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Enforce sign-in only              | Users must sign in before using Docker Desktop                 | Ensures users receive the benefits of your subscription and ensures security features are applied. In addition, you gain insights into users’ activity.                                                                                                    |
| Enforce single sign-on (SSO) only | If users sign in, they must sign in using SSO                  | Centralizes authentication and enforces unified policies set by the identity provider.                                                                                                                                                                     |
| Enforce both                      | Users must sign in using SSO before using Docker Desktop       | Ensures users receive the benefits of your subscription and ensures security features are applied. In addition, you gain insights into users’ activity. It also centralizes authentication and enforces unified policies set by the identity provider. |
| Enforce neither                   | If users sign in, they can use SSO or their Docker credentials | Lets users access Docker Desktop without barriers, at the cost of reduced security and insights.                                                                                                                                                  |

## Next steps

- To set up sign-in enforcement, see [Configure sign-in enforcement](methods.md).
- To configure SSO enforcement, see [Enforce SSO](/manuals/security/authentication/single-sign-on/connect.md).
