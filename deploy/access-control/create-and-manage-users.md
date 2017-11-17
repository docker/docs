---
title: Create and manage users
description: Learn how to administer user permissions in Docker Universal Control Plane.
keywords: authorize, authentication, users, teams, UCP, Docker
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
---

{% if include.ui %}
Docker Universal Control Plane provides built-in authentication and also
integrates with LDAP directory services.

> To enable LDAP and manage users and groups from your organization's directory,
> go to Admin > Admin Settings > Authentication and Authorization.
> [Learn how to integrate with an LDAP directory](../configure/external-auth/index.md).

{% if include.version=="ucp-3.0" %}
To use UCP built-in authentication, you must manually create users. New users
are assigned a default permission level so that they can access the cluster.
You can optionally grant them UCP administrator permissions.

To create a new user in the UCP web UI:
1. Navigate to the **Users** page.
2. Click **Create User**.
3. Input username, password, and full name.
4. Click **Create**.
5. [optional] Check "Is a Docker EE Admin".

> Check `Is a Docker EE Admin` to grant a user permission to change the cluster
> configuration and manage grants, roles, and collections.

![](../images/ucp_usermgmt_users_create01.png){: .with-border}

![](../images/ucp_usermgmt_users_create02.png){: .with-border}

{% elsif include.version=="ucp-2.2" %}
To use UCP built-in authentication, you must manually create users. New users
are assigned a default permission level so that they can access the swarm.
You can optionally grant them UCP administrator permissions.

To create a new user in the UCP web UI:
1. Navigate to the **Users** page.
2. Click **Create User**.
3. Input username, password, and full name.
4. Click **Create**.
5. [optional] Check "Is a UCP admin?".

![](../images/create-users-1.png){: .with-border}

> Check `Is a UCP admin?` to grant a user permission to change the swarm
> configuration and manage grants, roles, and collections.

![](../images/create-users-2.png){: .with-border}
{% endif %}

## Where to go next

* [Create and manage teams](create-and-manage-teams.md)
* [UCP permission levels](permission-levels.md)
{% endif %}
