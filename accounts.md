+++
title = "Account and repository management"
description = "Account and repository management"
keywords = ["docker, documentation, about, technology, hub, user interface, account management, accounts, teams, enterprise"]
[menu.main]
parent="workw_dtr"
weight=10
+++

# Account and repository management introduction

Administrators assign permissions to control users level of access to the
Trusted Registry. To access repositories, these users are grouped into teams and
organizations. Users can use the Trusted Registry UI to view their teams,
organizations, and the repositories they belong to.

This document describes:

* How account management is organized, and
* Provides examples of using the UI to organize and view users, teams, organizations, and repositories.

>**Note**: You can manage accounts and repositories from a Trusted Registry API. Refer to the API documentation located under the user drop down menu in the Trusted Registry UI to learn about it.

## Account types

Docker defines two types of accounts: users and organizations.

* **User** accounts correspond to individuals able to authenticate to the Trusted Registry.
* **Organization** accounts consists of subsets of users grouped into teams. Organizations can correspond to separate divisions, business units, or product teams within in your corporation.

Both user and organization accounts are uniquely identified by an account name.
These names also serve to organize repositories in the Trusted Registry into
namespaces corresponding to each account using the format
`<account-name>/<repository-name>`.

### User accounts

Users may be members of multiple teams and multiple organizations. As such, they may have differing permissions that are particular to that team. For example, team A can have admin permissions to manage teams while Team B can have only read permissions.

* Users, through their own namespace, can manage and share their own repositories.
* Users can be granted permission to repositories through team membership in an organization.
* Users can be an "owner" member of an organization. This grants them permission to manage all teams and repositories in that organization.

The following table depicts user permissions and roles:  

| Permissions/Roles                       | sys admin | org admin | org member | team member | user |
|:----------------------------------------|:---------:|:---------:|:----------:|:-----------:|:----:|
| orgs: create, edit, delete              |     x     |           |            |             |      |
| orgs: view public repos, teams, members |     x     |     x     |     x      |             |      |
| orgs: view public repos only            |     x     |     x     |     x      |      x      |  x   |
| teams: create, edit, delete             |     x     |     x     |            |             |      |
| teams: view public  repos, members      |     x     |     x     |     x      |      x      |      |
| teams: set repo permissions             |     x     |     x     |            |             |      |

### Organization accounts

Organizations are defined by its own namespace of repositories. They can consist
of one or more teams with each team having its own set of permissions. It can
also be managed by anyone in an initial "owners" team, which is created by
default.

* **Organization owners** have the highest level of permissions within the
organization. They can manage all teams and repositories and create, modify,
or delete teams.

* **Organization members** must be a member of one or more teams within the
organization.  

All organization members can see teams and their members. However, they are not
visible to users outside that organization.


### Teams in an organization

Teams are a convenient grouping of users. Organization owners can create a team and control team membership using the following authentication methods:

* Importing a managed list of users.
* Setting up LDAP integration and configuring team membership sync with an LDAP group.

Organization owners, other than the system administrator, are the only people
who can create, modify, or delete those teams that belong to that organization.

The following table depicts teams permissions to their repositories:

| Repository access      | read | read-write | admin |
|:-----------------------|:----:|:----------:|:-----:|
| view/ browse           |  x   |     x      |   x   |
| pull                   |  x   |     x      |   x   |
| push                   |      |     x      |   x   |
| edit/delete tags       |      |     x      |   x   |
| edit the description   |      |            |   x   |
| make public or private |      |            |   x   |
| manage user access     |      |            |   x   |

**Team permissions are additive**. This means that the highest level of permissions is always granted and can’t be reduced. For example, a user belongs to two teams. Team A grants its members "read-write" access to a repository and Team B grants "read-only" access to the same repository. That user will have "read-write" access to the repository because it is the higher permission level of the two.

### Repository behavior

Repositories are identified by a namespace value. A namespace has the format `account-name/repository-name`. The `account-name` can be either a user or organization account. Upon creation, you can determine whether the repository has either public or private visibility.

Public repositories are visible to all accounts. But it can only be written to by accounts granted explicit write access. However private repositories can’t be discovered by any account type unless it has explicit read access to it.

You must first create a repository before pushing an image to it. Otherwise you will see the following error message:

```
% docker push my.dtr.host/user1/myimage
The push refers to a repository [my.dtr.host/user1/myimage] (len: 1) 1d073211c498: Image push failed
unauthorized: access to the requested resource is not authorized
```

## Manage repositories, organizations, and teams

This section provides workflows for you to manage your users using the Trusted Registry’s repositories.


### Create an organization

1. From the Trusted Registry dashboard, click the Organizations submenu.

2. Click New organization. The Organization details screen displays.

3. Enter a unique name for your organization and save.


### Add teams to your organization

1. From the Trusted Registry dashboard, click the Organizations submenu.

2. Find your organization and select it. The Organization details screen displays.

3. Select the submenu Teams.

4. By default, the `owners` team box displays where you can add members who will have full admin access to that repository.

5. Click New team and enter the required fields.

6. Click Add members to select members to the team. Save your work.

At this point, you have created an organization and populated it with at least
one team. Next, you will either create or associate a repository to that
organization.


### Create a new repository for the team or organization

1. From the Organization details screen, click the desired organization.

2. If you click New repository, follow the steps to create a new repository that is associated to the organization.

3. To associate that repo to a team, click the Teams subtab, then click the targeted team.

4. Click Add repository and select a permission set from the drop down menu.  

5. You can either create a new repository or find an existing repository to associate to the team.


### Create a new repository

1. From the Trusted Registry dashboard, click the Repositories submenu.

2. Click New repository. The Repositories details screen displays.

3. Select an account type and enter a repository name.

4. Determine visibility. By default, the repository is public.

5. (Optional) Enter a description.

6. Save your work.

From the Repository submenu, you can:

* View, search, and filter the list of your repositories.
* Create either public or private repositories.
* Select a repository and edit it.
* Drill down to see details and teams that are associated with it.  

### View repository details

1. From the Trusted Registry dashboard, navigate to the Repositories menu.

2. Find a repository that contains images in it.

3. Click the submenu to see either details, tags, or settings.

The **Details** screen contains a brief description, a longer README, and the permissions associated with it.

The **Tag** screen contains the list of image tags. If you wanted to delete an image for garbage collection, click the garbage can icon beside it.

![Repositories page</repositories>](images/accounts-long-tag.png)

The **Settings** screen is where you edit the details screen.

## See also

* To configure for your environment, see the
[Configuration instructions](../configure/configuration.md).
* To administer the Trusted Registry, see the [Admin guide ](adminguide.md).
* To use Docker Trusted Registry, see the [User guide](userguide.md).
* To upgrade, see the [Upgrade guide](../install/upgrade.md).
* To see previous changes, see the [release notes](release-notes.md).
