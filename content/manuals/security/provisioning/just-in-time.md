---
description: Learn how Just-in-Time provisioning creates and updates accounts during SSO sign-in, and how to turn it off.
keywords: user provisioning, just-in-time provisioning, JIT, autoprovision, SSO, single sign-on, SCIM, disable JIT, Docker Admin, admin, security
title: Just-in-Time provisioning
linkTitle: Just-in-Time
weight: 20
aliases:
  - /security/for-admins/provisioning/just-in-time/
  - /enterprise/security/provisioning/just-in-time/
---

{{< summary-bar feature_name="SSO" >}}

Just-in-Time (JIT) provisioning creates and updates user accounts during SSO
authentication. JIT verifies that users belong to the organization and assigns
them to teams based on your identity provider (IdP) configuration. JIT doesn't
deprovision users.

## Overview

JIT gives users access the first time they sign in with SSO, with no
invitations to send or directory sync to set up. It fits organizations that
want the quickest path to access and can remove users manually. 

The diagram shows how JIT provisioning works:

```mermaid {title="SSO authentication with and without JIT provisioning" caption="With JIT, Docker adds users to organizations and teams at sign-in. Without it, only members can sign in. Highlighted steps happen only with JIT."}
---
config:
  themeVariables:
    fontSize: 14px
  flowchart:
    wrappingWidth: 80
    nodeSpacing: 14
    rankSpacing: 22
    padding: 5
    curve: linear
    subGraphTitleMargin:
      top: 4
      bottom: 12
---
flowchart LR
  start(["Sign in with SSO"]) --> invite["Accept pending invitation"]
  invite --> jit{"JIT enabled?"}
  jit -->|Yes| map{"IdP group mappings?"}
  map -->|Yes| assign
  map -->|No| default
  subgraph jitbox["JIT provisioning"]
    assign["Add to orgs and teams from IdP"]
    default["Add to default org and team"]
  end
  assign --> ok(["SSO successful"])
  default --> ok
  jit -->|No| member{"Member of org?"}
  member -->|Yes| ok
  member -->|No| denied(["Access denied"])

  classDef success fill:#22c55e22,stroke:#22c55e
  classDef failure fill:#ef444422,stroke:#ef4444
  class ok success
  class denied failure
  style jitbox fill:#3b82f622,stroke:#3b82f6
```

When a user signs in with SSO, Docker looks for a Docker account with the email address your IdP sent. Docker creates one if none exists, and generates a unique username from the user's email and name, and random numbers. If one does exist, Docker updates the user's full name, if needed.

If JIT is disabled, sign-in succeeds for a member of the organization,
including a user whose invitation Docker accepted at sign-in. Everyone else sees an
`Access denied` error and must ask an administrator for an invitation.

Disable JIT provisioning if you want SCIM to be the source of truth for
provisioning. For example, you have multiple organizations or want to restrict
access based on your organization's security configuration.

## Enable

When you create an SSO connection, Docker turns on JIT provisioning by
default. To enable JIT provisioning:

1. Set up an [SSO connection](/manuals/security/authentication/single-sign-on/connect.md)
   for your organization. Docker turns on JIT when you create the connection.
1. Optional. Set the default organization and team that Docker uses for users
   who sign in without IdP group mappings. You set these when you create the
    connection.
1. Optional. Set up
   [group mapping](/manuals/security/provisioning/scim/group-mapping.md#configure-group-mapping-with-sso)
   to assign users to organizations and teams from your IdP.
1. Have your users sign in with SSO.

Setting up SSO requires that at least one provisioning method stay enabled: JIT,
SCIM, or both. Before adding SCIM, review
[how SCIM works with JIT](/manuals/security/provisioning/scim/_index.md#choose-how-scim-works-with-jit).

## Disable

Disabling JIT provisioning may disrupt your users' access and workflows,
because only existing organization members can then sign in with SSO. To
disable it:

1. Go to [Docker Home](https://app.docker.com/) and select your organization
   from the top-left account drop-down.
1. Select **Identity & auth**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Actions** icon, then select
   **Disable JIT provisioning**.
1. Select **Disable** to confirm.

You can disable JIT only while
[SCIM is enabled](/manuals/security/provisioning/scim/provision-scim.md#enable-scim-in-docker),
and you can disable SCIM only while JIT is enabled.

## Re-enable

If you disabled JIT provisioning, you can turn it back on:

1. Go to [Docker Home](https://app.docker.com/) and select your organization
   from the top-left account drop-down.
1. Select **Identity & auth**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Actions** icon, then select
   **Enable JIT provisioning**.
1. Select **Enable** to confirm.

## Next steps

- Set up [group mapping](/manuals/security/provisioning/scim/group-mapping.md)
  to automatically assign users to teams.
- Review [Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md).
