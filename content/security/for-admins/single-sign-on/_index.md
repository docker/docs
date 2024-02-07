---
description: Overview of Single Sign-On
keywords: Single Sign-On, SSO, sign-on, admin, docker hub, admin console, security
title: Single Sign-On overview
aliases:
- /single-sign-on/
- /admin/company/settings/sso/
- /admin/organization/security-settings/sso-management/
---

SSO allows users to authenticate using their identity providers (IdPs) to access Docker. SSO is available for a whole company, and all associated organizations, or an individual organization that has a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

## How it works

When you enable SSO, this redirects to your IdP's authentication page to sign in. They can't authenticate using their Docker login credentials (Docker ID and password). Docker supports Service Provider Initiated SSO flow. Your users must sign in to Docker Hub or Docker Desktop to initiate the SSO authentication process.

The following diagram shows how SSO operates and is managed in Docker Hub and Docker Desktop. In addition, it provides information on how to authenticate between your IdP.

![SSO architecture](images/security/for-admins/single-sign-on/SSO.png)

## How to set it up

Before enabling SSO in Docker, administrators must first configure their IdP to work with Docker. Docker provides the Assertion Consumer Service (ACS) URL and the Entity ID. Administrators use this information to establish a connection between their IdP server and Docker Hub.

After establishing the connection between the IdP server and Docker, administrators sign in to Docker Hub or the Docker Admin Console and complete SSO enablement.

When you enable SSO for your company, a first-time user can sign in to Docker Hub using their company's domain email address. They're then added to your company, assigned to an organization, and optionally assigned to a team.

Administrators can then choose to enforce SSO sign-in and manage SSO connections for their individual company.

### SSO attributes

When a user signs in using SSO, Docker obtains the following attributes from the IdP:

- **Email address** - unique identifier of the user
- **Full name** - name of the user
- **Groups (optional)** - list of groups to which the user belongs

If you use SAML for your SSO connection, Docker obtains these attributes from the SAML assertion message. Your IdP may use different naming for SAML attributes than those in the previous list. The following table lists the possible SAML attributes that can be present in order for your SSO connection to work. 

You can also configure attributes to override default values, such as default team or organization. See [role mapping](../scim.md#set-up-role-mapping).

| SSO attribute | SAML assertion message attributes |
| ---------------- | ------------------------- |
| Email address    | `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`, `email`                           |
| Full name        | `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"`, `name`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"`  |
| Groups (optional) | `"http://schemas.xmlsoap.org/claims/Group"`, `"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"`, `Groups`, `groups` |
| Docker Org (optional)        | `dockerOrg`   |
| Docker Team (optional)     | `dockerTeam`  |
| Docker Role (optional)      | `dockerRole`  |

> **Important**
>
> If none of the email address attributes listed in the previous table are found, SSO returns an error.
{ .important}

## Prerequisites

* You must first notify your company about the new SSO login procedures.
* Verify that your members have Docker Desktop version 4.4.2, or later, installed on their machines.
* If your organization uses the Docker Hub CLI, new org members must [create a Personal Access Token (PAT)](/docker-hub/access-tokens/) to sign in to the CLI. There is a grace period for existing users, which will expire in the near future. Before the grace period ends, your users can sign in from Docker Desktop CLI using their previous credentials until PATs are mandatory.
In addition, you should add all email addresses to your IdP.
* Confirm that all CI/CD pipelines have replaced their passwords with PATs.
* For your service accounts, add your additional domains or enable it in your IdP.

## What's next?

- Start [configuring SSO](configure/_index.md)
- Explore the [FAQs](../../../faq/security/single-sign-on/faqs.md)
