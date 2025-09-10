---
description: Learn about provisioning users for your SSO configuration.
keywords: provision users, provisioning, JIT, SCIM, group mapping, sso, docker admin, admin, security
title: Provision users
linkTitle: Provision
weight: 20
aliases:
 - /security/for-admins/provisioning/
grid:
  - title: "Just-in-Time (JIT) provisioning"
    description: "Set up automatic user creation on first sign-in. Ideal for smaller teams with minimal setup requirements."
    icon: "schedule"
    link: "just-in-time/"
  - title: "SCIM provisioning"
    description: "Enable continuous user data synchronization between your IdP and Docker. Best for larger organizations."
    icon: "sync"
    link: "scim/"
  - title: "Group mapping"
    description: "Configure role-based access control using IdP groups. Perfect for strict access control requirements."
    icon: "group"
    link: "group-mapping/"
---

{{< summary-bar feature_name="SSO" >}}

After configuring your SSO connection, the next step is to provision users. This process ensures that users can access your organization through automated user management.

This page provides an overview of user provisioning and the supported provisioning methods.

## What is provisioning?

Provisioning helps manage users by automating tasks like account creation, updates, and deactivation based on data from your identity provider (IdP). There are three methods for user provisioning, each offering benefits for different organizational needs:

| Provisioning method | Description | Default setting in Docker | Recommended for |
| :--- | :--- | :------------- | :--- |
| Just-in-Time (JIT) | Automatically creates and provisions user accounts when they first sign in via SSO | Enabled by default | Organizations needing minimal setup, smaller teams, or low-security environments |
| System for Cross-domain Identity Management (SCIM) | Continuously syncs user data between your IdP and Docker, ensuring user attributes remain updated without manual intervention | Disabled by default | Larger organizations or environments with frequent changes in user information or roles |
| Group mapping | Maps user groups from your IdP to specific roles and permissions within Docker, enabling fine-grained access control based on group membership | Disabled by default | Organizations requiring strict access control and role-based user management |

## Default provisioning setup

By default, Docker enables JIT provisioning when you configure an SSO connection. With JIT enabled, user accounts are automatically created the first time a user signs in using your SSO flow.

JIT provisioning may not provide sufficient control or security for some organizations. In such cases, SCIM or group mapping can be configured to give administrators more control over user access and attributes.

## SSO attributes

When a user signs in through SSO, Docker obtains several attributes from your IdP to manage the user's identity and permissions. These attributes include:

- Email address: The unique identifier for the user
- Full name: The user's complete name
- Groups: Optional. Used for group-based access control
- Docker Org: Optional. Specifies the organization the user belongs to
- Docker Team: Optional. Defines the team the user belongs to within the organization
- Docker Role: Optional. Determines the user's permissions within Docker
- Docker session minutes: Optional. Sets the session duration before users must re-authenticate with their IdP. Must be a positive integer greater than 0. If not provided, default session timeouts apply

> [!NOTE]
>
> Default session timeouts apply when Docker session minutes is not specified. Docker Desktop sessions expire after 90 days or 30 days of inactivity. Docker Hub and Docker Home sessions expire after 24 hours.

## SAML attribute mapping

If your organization uses SAML for SSO, Docker retrieves these attributes from the SAML assertion message. Different IdPs may use different names for these attributes.

| SSO Attribute	| SAML Assertion Message Attributes |
| :--- | :--- |
| Email address |	`"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`, `email` |
| Full name	| `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"`, `name`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"` |
| Groups (optional) |	`"http://schemas.xmlsoap.org/claims/Group"`, `"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"`, `Groups`, `groups` |
| Docker Org (optional)	| `dockerOrg` |
| Docker Team (optional) |	`dockerTeam` |
| Docker Role (optional) |	`dockerRole` |
| Docker session minutes (optional) | `dockerSessionMinutes`, must be a positive integer > 0 |

## Next steps

Choose the provisioning method that best fits your organization's needs:

{{< grid >}}