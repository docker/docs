---
description: Learn how single sign-on works, how to set it up, and the required SSO attributes.
keywords: Single Sign-On, SSO, sign-on, admin, docker hub, admin console, security
title: Single sign-on overview
linkTitle: Single sign-on
aliases:
- /single-sign-on/
- /admin/company/settings/sso/
- /admin/organization/security-settings/sso-management/
weight: 10
---

{{< summary-bar feature_name="SSO" >}}

Single sign-on (SSO) lets users access Docker by authenticating using their identity providers (IdPs). SSO is available for a whole company, and all associated organizations within that company, or an individual organization that has a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

## How SSO works

When you enable SSO, Docker supports a non-IdP-initiated SSO flow for user login. Instead of users authenticating using their Docker username and password, they are redirected to your identity provider's authentication page to sign in. Users must sign in to Docker Hub or Docker Desktop to initiate the SSO authentication process.

The following diagram shows how SSO operates and is managed in Docker Hub and Docker Desktop. In addition, it provides information on how to authenticate between your IdP.

![SSO architecture](images/SSO.png)

## How to set it up

SSO is configured using the following steps:
1. [Configure SSO](../single-sign-on/configure.md) by creating and verifying a domain in Docker.
2. [Create your SSO connection](../single-sign-on/connect.md) in Docker and your IdP.
3. Cross-connect Docker and your IdP.
4. Test your connection.
5. Provision users.
6. Optional. [Enforce sign-in](../enforce-sign-in/_index.md).
7. [Manage your SSO configuration](../single-sign-on/manage.md).

Once your SSO configuration is complete, a first-time user can sign in to Docker Hub or Docker Desktop using their company's domain email address. Once they sign in, they are added to your company, assigned to an organization, and if necessary, assigned to a team.

## Prerequisites

Before configuring SSO, ensure you meet the following prerequisites:
* Notify your company about the new SSO sign in procedures.
* Verify that all users have Docker Desktop version 4.4.2 or later installed.
* If your organization is planning to [enforce SSO](/manuals/security/for-admins/single-sign-on/connect.md#optional-enforce-sso), members using the Docker CLI are required to [create a Personal Access Token (PAT)](/docker-hub/access-tokens/). The PAT will be used instead of their username and password. Docker plans to deprecate signing in to the CLI with a password in the future, so using a PAT will be required to prevent issues with authentication. For more details see the [security announcement](/security/security-announcements/#deprecation-of-password-logins-on-cli-when-sso-enforced).
* Ensure all your Docker users have a valid user on your IdP with the same email address as their Unique Primary Identifier (UPN).
* Confirm that all CI/CD pipelines have replaced their passwords with PATs.
* For your service accounts, add your additional domains or enable it in your IdP.

## What's next?

- Start [configuring SSO](../../for-admins/single-sign-on/configure.md) in Docker
- Explore the [FAQs](../../../security/faqs/single-sign-on/_index.md)
