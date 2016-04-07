<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/accounts/"]
title = "Permission levels"
description = "Learn about the permission levels available on Docker Trusted Registry."
keywords = ["docker, registry, security, permissions, users"]
[menu.main]
parent="dtr_menu_user_management"
identifier="dtr_permission_levels"
+++
<![end-metadata]-->

# DTR permission levels

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
