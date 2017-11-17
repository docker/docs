---
title: Access control model
description: Manage access to resources with role-based access control.
keywords: ucp, grant, role, permission, authentication, authorization
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
next_steps:
- path: /deploy/access-control/
  title: Create and manage users
- path: /deploy/access-control/
  title: Create and manage teams
- path: /deploy/access-control/
  title: Deploy a service with view-only access across an organization
- path: /deploy/access-control/
  title: Isolate volumes between two different teams
- path: /deploy/access-control/
  title: Isolate swarm nodes between two different teams
---

{% if include.ui %}

With Docker Universal Control Plane (UCP), you can configure how users access
resources by assigning role-based permissions with grants.

{% if include.version=="ucp-3.0" %}

UCP administrators can control who views, edits, and uses Swarm and Kubernetes
resources. They can grant and manage permissions to enforce fine-grained access
control as needed.

## Grants

Grants define which users can access what resources. Grants are effectively
Access Control Lists (ACLs), which, when grouped together, can provide
comprehensive access policies for an entire organization.

A grant is made up of a *subject*, *namespace*, and *role*.

Administrators are users who create subjects, define namespaces by labelling
resources, define roles by selecting allowable operations, and apply grants to
users and teams.

> Only an administrator can manage grants, subjects, roles, and resources.

## Subjects

A subject represents a user, team, or organization. A subject is granted a
role that defines permitted operations against one or more resources.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level.
  A team exists only as part of an organization, and all of its members
  must be members of the organization. Team members share organization
  permissions. A team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

## Namespaces

A namespace is ...

## Roles

Roles define what operations can be done by whom against which cluster
resources. A role is a set of permitted operations against a resource that is
assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to
view and schedule a node (in a granted namespace) but not update it.

Most organizations use different roles to assign the right kind of access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

You can build custom roles to meet your organizational needs or use the
following built-in roles:

- **View Only** - Users can see all cluster resources but not edit or use them.
- **Restricted Control** - Users can view containers and run a shell inside a
  container process (with `docker exec`) but not view or edit other resources.
- **Full Control** - Users can perform all operations against granted resources.
- **Scheduler** - Users can view and schedule nodes.

[Learn more about roles and permissions](permission-levels.md).


{% elsif include.version=="ucp-2.2" %}

UCP administrators can control who views, edits, and uses resources such as
nodes, services, images, networks, and volumes. They can grant and manage
permissions to enforce fine-grained access control as needed.


## Grants
Grants define which users can access what resources. Grants are effectively
Access Control Lists (ACLs), which, when grouped together, can provide
comprehensive access policies for an entire organization.

A grant is made up of a *subject*, *resource collection*, and *role*.

Administrators are users who create subjects, define collections by labelling
resources, define roles by selecting allowable operations, and apply grants to
users and teams.

> Only an administrator can manage grants, subjects, roles, and resources.

![](../images/ucp-grant-model.svg)

## Subjects

A subject represents a user, team, or organization. A subject is granted a
role that defines permitted operations against one or more resources.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level.
  A team exists only as part of an organization, and all of its members
  must be members of the organization. Team members share organization
  permissions. A team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

## Collections

A collection is a group of resources that you define with labels and access by
specifying a directory-like path.

Swarm resources that can be placed into a collection include:

- Application configs
- Containers
- Networks
- Nodes (Physical or virtual)
- Services
- Secrets
- Volumes

## Roles

Roles define what operations can be done by whom against which cluster
resources. A role is a set of permitted operations against a resource that is
assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to view
and schedule a node (in a granted collection) but not update it.

Most organizations use different roles to assign the right kind of access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

You can build custom roles to meet your organizational needs or use the following
built-in roles:

- **View Only** - Users can see all cluster resources but not edit or use them.
- **Restricted Control** - Users can view containers and run a shell inside a
  container process (with `docker exec`) but not view or edit other resources.
- **Full Control** - Users can perform all operations against granted resources.
- **Scheduler** - Users can view and schedule nodes.

[Learn more about roles and permissions](permission-levels.md).


## Collection architecture

Before grants can be implemented, collections must group resources in a way that
makes sense for an organization.

For example, consider an organization with two application teams, Mobile and
Payments, which share cluster hardware resources but segregate access to their
individual applications.

```
orcabank (organization)
├── ops (team)
├── security (team)
├── mobile-dev (team)
└── payments-dev (team)
```

To define a potential access policy, the collection architecture should map to
the organizational structure. For a production UCP cluster, it might look like
this:

```
prod (collection)
├── mobile (sub-collection)
└── payments (sub-collection)
```

> A subject that has access to any level in a collection hierarchy has the
> same access to any collections below it.



## Transition from UCP 2.1 access control

- Your existing access labels and permissions are migrated automatically during
  an upgrade from UCP 2.1.x.
- Unlabeled "user-owned" resources are migrated into the user's private
  collection, in `/Shared/Private/<username>`.
- Old access control labels are migrated into `/Shared/Legacy/<labelname>`.
- When deploying a resource, choose a collection instead of an access label.
- Use grants for access control, instead of unlabeled permissions.

[See a deeper tutorial on how to design access control architectures.](access-control-design-ee-standard.md)

[Learn to manage access to resources by using collections.](manage-access-with-collections.md).

{% endif %}
{% endif %}
