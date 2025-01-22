---
description: Learn about provisioning users for your SSO configuration.
keywords: provision users, provisioning, JIT, SCIM, group mapping, sso, docker hub, hub, docker admin, admin, security
title: Provision users
linkTitle: Provision
weight: 20
---

{{< summary-bar feature_name="SSO" >}}

Once you've configured your SSO connection, the next step is to provision users. This process ensures that users can access your organization.
This guide provides an overview of user provisioning and supported provisioning methods.

## What is provisioning?

Provisioning helps manage users by automating tasks like creating, updating, and deactivating users based
on data from your identity provider (IdP). There are three methods for user provisioning, with benefits for
different organization needs:

| Provisioning method | Description | Default setting in Docker | Recommended for |
| :--- | :--- | :------------- | :--- |
| Just-in-Time (JIT) | Automatically create and provisions user accounts when they first sign in via SSO | Enabled by default | Best for organizations who need minimal setup, who have smaller teams, or low-security environments |
| System for Cross-domain Identity Management (SCIM) | Continuously syncs user data between your IdP and Docker, ensuring user attributes remain updated without requiring manual updates | Disabled by default | Best for larger organizations or environments with frequent changes in user information or roles |
| Group mapping | Maps user groups from your IdP to specific roles and permissions within Docker, enabling fine-tuned access control based on group membership | Disabled by default | Best for organizations that require strict access control and for managing users based on their roles and permissions |

## Default provisioning setup

By default, Docker enables JIT provisioning when you configure an SSO connection. With JIT enabled, user accounts are automatically created the first time a user signs in using your SSO flow.

JIT provisioning may not provide the level of control or security some organizations need. In such cases, SCIM or group mapping can be configured to give administrators more control over user access and attributes.

## SSO attributes

When a user signs in through SSO, Docker obtains several attributes from your IdP to manage the user's identity and permissions. These attributes include:
- **Email address**: The unique identifier for the user
- **Full name**: The user's complete name
- **Groups**: Optional. Used for group-based access control
- **Docker Org**: Optional. Specifies the organization the user belongs to
- **Docker Team**: Optional. Defines the team the user belongs to within the organization
- **Docker Role**: Optional. Determines the user's permission within Docker

If your organization uses SAML for SSO, Docker retrieves these attributes from the SAML assertion message. Keep in mind that different IdPs may use different names for these attributes. The following reference table outlines possible SAML attributes used by Docker:

| SSO Attribute	| SAML Assertion Message Attributes |
| :--- | :--- |
| Email address |	`"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`, `email` |
| Full name	| `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"`, `name`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"` |
| Groups (optional) |	`"http://schemas.xmlsoap.org/claims/Group"`, `"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"`, `Groups`, `groups` |
| Docker Org (optional)	| `dockerOrg` |
| Docker Team (optional) |	`dockerTeam` |
| Docker Role (optional) |	`dockerRole` |

## What's next?

Review the provisioning method guides for steps on configuring provisioning methods:
- [JIT](/manuals/security/for-admins/provisioning/just-in-time.md)
- [SCIM](/manuals/security/for-admins/provisioning/scim.md)
- [Group mapping](/manuals/security/for-admins/provisioning/group-mapping.md)