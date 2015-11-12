+++
title = "Account and repository management"
description = "Account and repository management"
keywords = ["docker, documentation, about, technology, hub, user interface, account management, accounts, enterprise"]
[menu.main]
parent="smn_dhe"
weight=10
+++

# Account and repository management introduction

This document describes the various actions you can perform surrounding account
and repository management through the Docker Trusted Registry (Trusted Registry)
user interface (UI). What you see and do depends on your role and
permissions which is outlined in this document. Use the UI to quickly see, and
with the correct permissions do the following:

* manage repositories
* manage organizations
* create and view teams
* assign members to teams

**Example**: as a organization owner in the organization-owned repository
Animals repo, you might create multiple teams in your organization. You manage
the teams access. You give Team Alpha permissions so they can modify,
delete, and push tags. Team Beta can only view, browse, and pull tags.

**Note**: You can perform many functions in the UI as you can from the
command line or API. To see what you can do from a command line or API, refer to
the API documentation accessed from the Trusted Registry UI.

There are three scopes with which you can manage permissions:
* Organizations
* Teams
* Repositories

This document explains the relationship between users, organizations, teams,
and repositories and gives examples of potential workflows you might use.

## Accounts

Docker defines two types of accounts: users and organizations.

* Users are individuals. They are sometimes called members in the context of a team or organization.
* An organization is a group of members.

Both accounts are defined by a namespace containing two component names in the
form of account-name/repository-name. There is no limit to how many you can have of each type of account.

## Repositories

In the Trusted Registry, repositories can be either public or private. Anyone can  create them, but who sees and accesses them are determined by permissions as well.

**Public**:

* visible to all accounts in the system
* can only be written to by accounts granted explicit write access

**Private**:

* cannot be discovered by any account unless having explicit read access to it
* can be created by users and organizations

## User accounts

A user (member):
* Can belong to an organization
* Be a part of a team that is a part of that organization.
* Belong to more than one team, in more than one organization and have differing roles within those teams.

**Example**:
In team A, they can have admin permissions so they can help manage their group,
while in team B, those users only have read permission.  

User can also create repositories under their own name and share those
repositories with other users. They confer permissions to other users on a
per-repository basis. The following table depicts the combination of users and possible permissions:  

| Permissions/Roles                       | sys admin | org admin | org member | team member | user |
|-----------------------------------------|:---------:|:---------:|:----------:|:-----------:|:----:|
| orgs: create, edit, delete              |     x     |           |            |             |      |
| orgs: view public repos, teams, members |     x     |     x     |      x     |             |      |
| orgs: view public repos only            |     x     |     x     |      x     |      x      |   x  |
| teams: create, edit, delete             |     x     |     x     |            |             |      |
| teams: view public  repos, members      |     x     |     x     |      x     |      x      |      |
| teams: set repo permissions             |     x     |     x     |            |             |      |

## Organization accounts

System administrators can also create an organization account, with its own
namespace of repositories. Comprised of one or more teams, they can be managed
by anyone in an initial "owners" team, which is created by default.

* **Organization owners** have the highest level of permissions within the
organization. They can manage all teams and repositories and create, modify,
or delete teams.

* **Organization members** must be a member of one or more teams within the
organization.  

All organization members can see teams and their members. However, they are not
visible to users outside that organization.

## Teams

Teams are configured in two ways:

* As a list of users managed by an organization owner, or
* Through LDAP system integration which can then be periodically synced

The organization owner, other than the system administrator, is the only person
who can create, modify, or delete those teams that belong to that organization.

Teams, like users, can also be granted permissions to their repositories as seen in the  following table:

| Repository access      | read | read-write | admin |
|------------------------|:----:|:----------:|:-----:|
| view/ browse           |   x  |      x     |   x   |
| pull                   |   x  |      x     |   x   |
| push                   |      |      x     |   x   |
| edit/delete tags       |      |      x     |   x   |
| edit the description   |      |            |   x   |
| make public or private |      |            |   x   |
| manage user access     |      |            |   x   |

  **Note**: These permissions are additive. This means you cannot override a team level permission to prevent access to a specific repository. If a team has read-write access to the entire namespace of repositories, then granting that team 'read-only' access to a specific repository will not reduce their access to that repository, as the team will still have read-write access to that repository through its namespace access.

### Working with repositories, organizations, and teams

From the Trusted Registry dashboard, click the Repositories submenu.

From the Repository submenu, you can:

* View, search, and filter the list of your repositories
* Create either public or private repositories
* Select a repository and edit it
* Drill down to see details, teams that are associated with it, and settings.  

There are submenus which you can see additional information:

* **Details** screen: see any permissions or tags

* **Teams** screen: see teams and associated members. Select a team to see what repositories they are associated with and what permissions they have. This takes you to the Organization screen where you set a team's permissions on that page or select the team to delete it.

From the Organizations submenu, you can:

* Create a new organization
* View, delete, or edit an existing organization
* Add teams to it
* View  and add members to the team

## See also

* To configure for your environment, see the
[Configuration instructions](configuration.md).
* To administer the Trusted Registry, see [the Admin guide ](adminguide.md).
* To use Docker Trusted Registry, see the [User guide](userguide.md).
* To upgrade, see the [Upgrade guide](install/upgrade.md).
* To see previous changes, see the [release notes](release-notes.md).



<!---
\\Todo:
--->
