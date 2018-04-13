---
title: Create users and teams manually
description: Learn how to add users and define teams in Docker Universal Control Plane.
keywords: rbac, authorize, authentication, users, teams, UCP
redirect_from:
  - /datacenter/ucp/3.0/guides/authorization/create-users-and-teams-manually/
---

Users, teams, and organizations are referred to as subjects in Docker EE.

Individual users can belong to one or more teams but each team can only be in
one organization. At the fictional startup, Acme Company, all teams in the
organization are necessarily unique but the user, Alex, is on two teams:

```
acme-datacenter
├── dba
│   └── Alex*
├── dev
│   └── Bett
└── ops
    ├── Alex*
    └── Chad
```

## Authentication

All users are authenticated on the backend. Docker EE provides built-in
authentication and also integrates with LDAP directory services.

To use Docker EE's built-in authentication, you must [create users manually](#create-users-manually).

> To enable LDAP and authenticate and synchronize UCP users and teams with your
> organization's LDAP directory, see:
> - [Synchronize users and teams with LDAP in the UI](create-teams-with-ldap.md)
> - [Integrate with an LDAP Directory](../admin/configure/external-auth/index.md).

## Build an organization architecture

The general flow of designing an organization with teams in UCP is:

1. Create an organization.
2. Add users or enable LDAD (for syncing users).
3. Create teams under the organization.
4. Add users to teams manually or sync with LDAP.

### Create an organization with teams

To create an organization in UCP:

1. Click **Organization & Teams** under **User Management**.
2. Click **Create Organization**.
3. Input the organization name.
4. Click **Create**.

To create teams in the organization:

1. Click through the organization name.
2. Click **Create Team**.
3. Input a team name (and description).
4. Click **Create**.
5. Add existing users to the team. To sync LDAP users, see: [Integrate with an LDAP Directory](../admin/configure/external-auth/index.md).
   - Click the team name and select **Actions** > **Add Users**.
   - Check the users to include and click **Add Users**.

> **Note**: To sync teams with groups in an LDAP server, see [Sync Teams with LDAP](create-teams-with-ldap.md).

### Create users manually

New users are assigned a default permission level so that they can access the
cluster. To extend a user's default permissions, add them to a team and [create grants](grant-permissions.md). You can optionally grant them Docker EE
administrator permissions.

To manually create users in UCP:

1. Click **Users** under **User Management**.
2. Click **Create User**.
3. Input username, password, and full name.
4. Click **Create**.
5. Optionally, check "Is a Docker EE Admin" to give the user administrator 
   privileges.

> A `Docker EE Admin` can grant users permission to change the cluster
> configuration and manage grants, roles, and resource sets.

![](../images/ucp_usermgmt_users_create01.png){: .with-border}
![](../images/ucp_usermgmt_users_create02.png){: .with-border}

## Where to go next

- [Synchronize teams with LDAP](create-teams-with-ldap.md)
- [Define roles with authorized API operations](define-roles.md)
- [Group and isolate cluster resources](group-resources.md)
- [Grant role-access to cluster resources](grant-permissions.md)