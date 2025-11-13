---
title: Troubleshoot provisioning
linkTitle: Troubleshoot provisioning
description: Troubleshoot common user provisioning issues with SCIM and Just-in-Time provisioning
keywords: SCIM troubleshooting, user provisioning, JIT provisioning, group mapping, attribute conflicts
tags: [Troubleshooting]
toc_max: 2
aliases:
 - /security/troubleshoot/troubleshoot-provisioning/
---

This page helps troubleshoot common user provisioning issues including user roles, attributes, and unexpected account behavior with SCIM and Just-in-Time (JIT) provisioning.

## SCIM attribute values are overwritten or ignored

### Error message

Typically, this scenario does not produce an error message in Docker or your
IdP. This issue usually surfaces as incorrect role or team assignment.

### Causes

- JIT provisioning is enabled, and Docker is using values from your IdP's
SSO login flow to provision the user, which overrides
SCIM-provided attributes.
- SCIM was enabled after the user was already provisioned via JIT, so SCIM
updates don't take effect.

### Affected environments

- Docker organizations using SCIM with SSO
- Users provisioned via JIT prior to SCIM setup

### Steps to replicate

1. Enable JIT and SSO for your Docker organization.
1. Sign in to Docker as a user via SSO.
1. Enable SCIM and set role/team attributes for that user.
1. SCIM attempts to update the user's attributes, but the role or team
assignment does not reflect changes.

### Solutions

#### Disable JIT provisioning (recommended)

1. Sign in to [Docker Home](https://app.docker.com/).
1. Select **Admin Console**, then **SSO and SCIM**.
1. Find the relevant SSO connection.
1. Select the **actions menu** and choose **Edit**.
1. Disable **Just-in-Time provisioning**.
1. Save your changes.

With JIT disabled, Docker uses SCIM as the source of truth for user creation
and role assignment.

**Keep JIT enabled and match attributes**

If you prefer to keep JIT enabled:

- Make sure your IdP's SSO attribute mappings match the values being sent
by SCIM.
- Avoid configuring SCIM to override attributes already set via JIT.

This option requires strict coordination between SSO and SCIM attributes
in your IdP configuration.

## SCIM updates don't apply to existing users

### Causes

User accounts were originally created manually or via JIT, and SCIM is not
linked to manage them.

### Solution

SCIM only manages users that it provisions. To allow SCIM to manage an
existing user:

1. Remove the user manually from the Docker [Admin Console](https://app.docker.com/admin).
1. Trigger provisioning from your IdP.
1. SCIM will re-create the user with correct attributes.

> [!WARNING]
>
> Deleting a user removes their resource ownership (e.g., repositories).
Transfer ownership before removing the user.
