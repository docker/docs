+++
title = "Docker Trusted Registry API Design document"
description = "Docker Trusted Registry 1.3 API Design document"
draft = true
keywords = ["API, Docker, index, REST, documentation, Docker Trusted Registry, registry"]
[menu.main]
parent = "smn_dtrapi"
+++

# Docker Trusted Registry 1.3

## Organization Accounts and Repository Access Control

We need an API for fine-grained role based access control in Docker Trusted
Registry.

### Goals

1. An API for managing accounts

    This includes creating an account, listing existing accounts, creating a
    team within an organization, listing teamns in an organization, getting a
    specific team, listing members of a team, adding and removing members
    from a team (if managed whitelist) or editing ldap sync config.

2. LDAP Syncing

    We need to be able to automatically sync members of a team according to
    some LDAP group filter configured by an org admin.

3. An API for repositories and access control.

    This includes creating a repository, listing repositories for an account,
    adding collaborators to a repository, setting namespace-level access for
    teams, etc.

4. Minimal or No Data Migration.

    If we can help it, we’d like to not require any data migration in this
    release, only schema migration.

### Non-Goals

No UI. This release will be API only.

### Design

There are two types of accounts: Users and Organizations. We will be creating a
sense of account-owned namespaces in Docker Trusted Registry much like in Docker Hub Registry.
Account-owned (non-global) repositories will have two-component names of the
form `namespace/reponame`.

Users will be able to create a repository under their own namespace and be able
to control what other users have read-only, read-write, or admin access to any
of their repositories.

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

Repositories can be set as either public or private. Public repositories can be
read by any account in the system and only written to by accounts with explicit
write access granted to them. Private repositories will not be discoverable by
any account that does not have at least explicit read access to that
repository.

User owned repositories can only be given access to the owner and other
individual user accounts, i.e., you cannot grant access to a user-owned
repository to a team of users in an organization. If this level of control is
needed on a repository then consider moving that repository within an
organization namespace.

Organization owned repositories can only be given access to the teams within
that organization, i.e., you cannot grant access to an organization-owned
repository to an individual user account or to a team in another organization.
If this level of control is needed on a repository then consider adding those
individual users to a team within the owning organization or add users in the
other organization’s team to a team within the owning organization.

### Notable Differences from Docker Hub Registry

- Repositories must be explicitly created using the API. A `docker push` will
  not create a repository if it does not exist. This prevents a typo from
  creating an unwanted repository in Docker Trusted Registry. This policy
  will be globally enforced in Docker Trusted Registry 1.3.

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

### API

#### Authentication

Clients may authenticate API requests by providing either Basic Auth
credentials (i.e., username and password) in an "Authorization" header for each
request.

**Accounts**

- [Accounts API Doc](https://gist.github.com/jlhawn/1e7e35c7eed536f4bdf2)
- [Teams API Doc](https://gist.github.com/jlhawn/9a1f24d44e9fce541cfb)

**Repositories**

- [Repositories API Doc](https://gist.github.com/jlhawn/dbdeb81a2724e913f036)
- [User-Owned Repository Access API Doc](https://gist.github.com/jlhawn/90e5442e1b1503b1f970)
- [Organization-Owned Repository Access API Doc](https://gist.github.com/jlhawn/45c8d6587f632975a345)
- [Organization-Owned Repository Namespace Access API Doc](https://gist.github.com/jlhawn/68f78af0e6b0325facd3)
