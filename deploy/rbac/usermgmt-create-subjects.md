---
title: Create and configure users and teams
description: Learn how to add users and define teams in Docker Universal Control Plane.
keywords: rbac, authorize, authentication, users, teams, UCP, Docker
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
next_steps:
- path: /deploy/rbac/usermgmt-sync-with-ldap/
  title: Synchronize teams with LDAP
- path: /datacenter/ucp/2.2/guides/admin/configure/external-auth/
  title: Integrate with an LDAP Directory
- path: /deploy/rbac/usermgmt-define-roles/
  title: Create roles to authorize access
- path: /deploy/rbac/usermgmt-grant-permissions/
  title: Grant access to cluster resources
---

{% if include.ui %}

Users, teams, and organizations are referred to as subjects in Docker Enterprise
Edition.

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

To use Docker EE's built-in authentication, you must [create users manually](#Create-users-manually).

> To enable LDAP and authenticate and synchronize UCP users and teams with your
> organization's LDAP directory, see:
> - [Synchronize users and teams with LDAP in the UI](usermgmt-sync-with-ldap.md)
> - [Integrate with an LDAP Directory](/datacenter/ucp/2.2/guides/admin/configure/external-auth/index.md).

## Build an organization architecture

The general flow of designing an organization with teams in the Docker EE UI is:

1. Create an organization.
2. Add users or enable LDAP.
3. Create teams under the organization.
4. Add users to teams manually or sync with LDAP.

### Create an organization with teams

To create an organzation in the Docker EE UI:

1. Click **Organization & Teams** under **User Management**.
2. Click **Create Organization**.
3. Input the organization name.
4. Click **Create**.

To create teams in the organization:

1. Click through the organization name.
2. Click **Create Team**.
3. Input a team name (and description).
4. Click **Create**.
5. Add existing users to the team. If they don't exist, see [Integrate with an LDAP Directory](../../datacenter/ucp/2.2/guides/admin/configure/external-auth/index.md).
   - Click the team name and select **Actions** > **Add Users**.
   - Check the users to include and click **Add Users**.

> **Note**: To sync teams with groups in an LDAP server, see [Sync Teams with LDAP](./usermgmt-sync-with-ldap).


{% if include.version=="ucp-3.0" %}

### Create users manually

New users are assigned a default permission level so that they can access the
cluster. To extend a user's default permissions, add them to a team and [create grants](./usermgmt-grant-permissions/). You can optionally grant them Docker EE
administrator permissions.

To manally create users in the Docker EE UI:

1. Click **Users** under **User Management**.
2. Click **Create User**.
3. Input username, password, and full name.
4. Click **Create**.
5. [optional] Check "Is a Docker EE Admin".

> A `Docker EE Admin` can grant users permission to change the cluster
> configuration and manage grants, roles, and collections.

![](../images/ucp_usermgmt_users_create01.png){: .with-border}
![](../images/ucp_usermgmt_users_create02.png){: .with-border}


{% elsif include.version=="ucp-2.2" %}

### Create users manuallly

New users are assigned a default permission level so that they can access the
cluster. To extend a user's default permissions, add them to a team and [create grants](/deploy/rbac/usermgmt-grant-permissions/). You can optionally grant them Docker EE
administrator permissions.

To manally create users in the Docker EE UI:

1. Navigate to the **Users** page.
2. Click **Create User**.
3. Input username, password, and full name.
4. Click **Create**.
5. [optional] Check "Is a UCP admin?".

> A `UCP Admin` can grant users permission to change the cluster configuration
> and manage grants, roles, and collections.

![](../images/create-users-1.png){: .with-border}
![](../images/create-users-2.png){: .with-border}

{% endif %}
{% endif %}
