+++
title = "DTR Accounts & Repos API: Intro & Overview"
description = "Overview of the structure and design of the DTR Accounts & Repos API"
keywords = ["API, Docker, index, REST, documentation, Docker Trusted Registry, registry"]
[menu.main]
parent = "smn_dtrapi"
+++

# Docker Trusted Registry 1.3: Accounts & Repos API

## Introduction

The Accounts & Repos API lets you integrate Docker Trusted Registry (DTR) with your enterprise's organizational structure by providing fine-grained, role-based access control for your repositories. Specifically, this API provides:

* An API for account management, including creating an account, listing existing accounts, creating a team within an organization, listing teamns in an organization, getting a specific team, listing members of a team, adding and removing members from a team (if using a managed whitelist), or editing LDAP syncing configuration.

* Methods for syncing members of a team in DTR with an LDAP group filter configured by an admin.

* An API for repository management and access control, including creating a repository, listing repositories for an account, adding collaborators to a repository, setting namespace-level access for teams, etc.

The API is designed so that minimal data migration is required, only schema migration. There is no UI accompanying this API.

## Design overview

This API defines two types of accounts that can own repositories: Users and Organizations. Account-owned (i.e., non-global) repos define a namespace similar to that of the Docker Hub, with two component names in the form `namespace/reponame`. 

Repositories can be either public or private. Public repositories can be
read by any account in the system, but can only be written to by accounts granted explicit write access. Private repositories cannot be discovered by
any account that does not have at least explicit read access to that
repository.

### User accounts

DTR users can create a repository under their own namespace and can control which other users have read-only, read-write, or admin access to any
of their repositories.

User owned repositories can only be accessed by the owner and other
individual user accounts, i.e., you cannot grant access to a user-owned
repository to a team of users in an organization. If a repository requires this level of control, consider moving it within an organization namespace.

### Organization accounts

System administrators can also create an Organization account, with it’s own
namespace of repositories. Organization accounts will have one or more teams
which can be managed by anyone in the initial ‘owners’ team. Teams can either
be created with a managed whitelist of users known to the system or with an
LDAP group search filter which will be periodically synced.

Any member of an organization’s owners team can create new repositories under
the organization’s namespace and even create and edit other teams. Each team
can be given read-only or read-write access to all repositories in the
organization’s namespace and/or be granted separate levels of access on a
per-repository basis. Permissions are additive, so there is no way to override
a team level permission to prevent access to a specific repository.

Teams within an organization can also be granted read-only, read-write, or
admin level access to all repositories in the organizations namespace. This
allows a team to pull, push, and/or manage repositories for an organization
but *not* manage the teams themselves.

Organization-owned repositories can only be given access to the teams within
that organization, i.e., you cannot grant access to an organization-owned
repository to an individual user account or to a team in another organization.
If this level of control is needed on a repository then consider adding those
individual users to a team within the owning organization or add users in the
other organization’s team to a team within the owning organization.

### Notable differences from Docker Hub

- Repositories must be explicitly created using the API. A `docker push` will
  not create a repository if it does not exist. This prevents a typo from
  creating an unwanted repository in DTR. This policy
  will be globally enforced in DTR 1.3.

- Organizations can only be created by system admins. This should prevent the
  proliferation of unwanted organization accounts.

- Collaborators on user-owned repositories can be given more granularity of
  access. Docker Hub Registry offers only read-write access. Docker Trusted
  Registry plans to offer read-only, read-write, and admin access per
  user-owned repository.

- Teams can be granted access to all repositories in an organization's
  namespace. Docker Hub Registry offers team access control on a
  per-repository level only, and only an organization's 'owners' team can
  manage access and create new repositories. Docker Trusted Registry plans to
  also offer the ability to grant a team access to (and/or the ability to
  manage) all repositories under a namespace.

- Teams within an organization will be visible to all members of the
  organization. In Docker Hub Registry, users are 'blind' to teams that they
  themselves are not a member of. In Docker Trusted Registry, teams will be
  visible to the members of the organization, but will not be able to see a
  teams's members unless they are also a member of that team.
  
### Authentication

Clients may authenticate API requests by providing either Basic Auth
credentials (i.e., username and password) in an "Authorization" header for each
request.

### API Documentation

The following documents detail the API:

- [Accounts API Doc]({{< relref "dtr_1_3_accounts.md" >}})
- [Teams API Doc]({{< relref "dtr_1_3_teams.md" >}})
- [Repositories API Doc]({{< relref "dtr_1_3_repositories.md" >}})
- [User-Owned Repository Access API Doc]({{< relref "dtr_1_3_user_repo_access.md" >}})
- [Organization-Owned Repository Access API Doc]({{< relref "dtr_1_3_team_repo_access.md" >}})
- [Organization-Owned Repository Namespace Access API Doc]({{< relref "dtr_1_3_team_repo_namespace_access.md" >}})

