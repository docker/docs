---
title: Create teams with LDAP
description: Learn how to enable LDAP and sync users and teams in Docker Universal Control Plane.
keywords: authorize, authentication, users, teams, UCP, LDAP
---

To enable LDAP in UCP and sync to your LDAP directory:

1. Click **Admin Settings** under your username drop down.
2. Click **Authentication & Authorization**.
3. Scroll down and click `Yes` by **LDAP Enabled**. A list of LDAP settings displays.
4. Input values to match your LDAP server installation.
5. Test your configuration in UCP.
6. Manually create teams in UCP to mirror those in LDAP.
6. Click **Sync Now**.

If Docker EE is configured to sync users with your organization's LDAP directory
server, you can enable syncing the new team's members when creating a new team
or when modifying settings of an existing team.

For more, see: [Integrate with an LDAP Directory](../admin/configure/external-auth/index.md).

![](../images/create-and-manage-teams-5.png){: .with-border}

## Binding to the LDAP server

There are two methods for matching group members from an LDAP directory, direct
bind and search bind.

Select **Immediately Sync Team Members** to run an LDAP sync operation
immediately after saving the configuration for the team. It may take a moment
before the members of the team are fully synced.

### Match Group Members (Direct Bind)

This option specifies that team members should be synced directly with members
of a group in your organization's LDAP directory. The team's membership will by
synced to match the membership of the group.

- **Group DN**: The distinguished name of the group from which to select users.
- **Group Member Attribute**: The value of this group attribute corresponds to
  the distinguished names of the members of the group.


### Match Search Results (Search Bind)

This option specifies that team members should be synced using a search query
against your organization's LDAP directory. The team's membership will be
synced to match the users in the search results.

- **Search Base DN**: Distinguished name of the node in the directory tree where
  the search should start looking for users.
- **Search Filter**: Filter to find users. If null, existing users in the search
  scope are added as members of the team.
- **Search subtree**: Defines search through the full LDAP tree, not just one
  level, starting at the Base DN.
