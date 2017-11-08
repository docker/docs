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
{% if include.version=="ucp-3.0" %}
This topic is under construction.
With Docker Universal Control Plane, you can control who creates and edits
resources, such as nodes, services, images, networks, and volumes. You can grant
and manage permissions to enforce fine-grained access control as needed.

## Grant access to Swarm resources
This topic is under construction.

## Grant access to Kubernetes resources
This topic is under construction.

## Transition from UCP 2.2 access control
This topic is under construction.


{% elsif include.version=="ucp-2.2" %}
## Grant access to Swarm resources

UCP administrators control how subjects (users, teams, organizations) access
resources (collections) by assigning role-based permissions with *grants*.

A grant is made up of a *subject*, *role*, and *resource collection*.

Administrators are users who create subjects, define collections by labelling
resources, define roles by selecting allowable operations, and apply grants to
users and teams.

> Only an administrator can manage grants, subjects, roles, and resources.

![](../images/ucp-grant-model.svg)

## Subjects

A subject represents a user, team, or organization. A subject is granted a
role for a collection of resources.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level.
  A team exists only as part of an organization, and all of its members
  must be members of the organization. Team members share organization
  permissions. A team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

## Roles

A role is a set of permitted API operations that you can assign to a specific
subject and collection by using a grant. UCP administrators view and manage
roles by navigating to the **Roles** page.
[Learn more about roles and permissions](permission-levels.md).

## Resource collections

Docker EE allows you to control access to cluster resources with *collections*.
A collection is a group of resources that you access by specifying a
directory-like path.

Resources that can be placed into a collection include:

- Physical or virtual nodes
- Containers
- Services
- Networks
- Volumes
- Secrets
- Application configs

## Collection architecture

Grants define which users can access what resources. Grants are effectively
Access Control Lists (ACLs), which, when grouped together, can provide
comprehensive access policies for an entire organization.

Before grants can be implemented, collections must be defined and group
resources in a way that makes sense for an organization.

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

## Role composition

Roles define what operations can be done against cluster resources. Most
organizations use different roles to assign the right kind of access. A given
team or user may have different roles provided to them depending on what
resource they are accessing.

UCP provides default roles and lets you build custom roles. For example, here,
three different roles are used:

- **Full Control** (default role) - Allows users to perform all operations
  against cluster resources.
- **View Only** (default role) - Allows users to see all cluster resources but
  not edit or delete them.
- **Dev** (custom role) - Allows users to view containers and run a shell inside
  a container process (with `docker exec`) but not to view or edit any other
  cluster resources.

## Grant composition

The following four grants define the access policy for the entire organization
for this cluster. They tie together the collections that were created, the
default and custom roles, and also teams of users that are in UCP.

![image](../images/access-control-grant-composition.png){: .with-border}

## Access architecture

The resulting access architecture defined by these grants is depicted below.

![image](../images/access-control-collection-architecture.png){: .with-border}

There are four teams that are given access to cluster resources:

- The `ops` team has `Full Control` against the entire `/prod` collection. It
  can deploy, view, edit, and remove applications and application resources.
- The `security` team has the `View Only` role. They can see, but not edit, all
  resources in the `/prod` collection.
- The `mobile` team has the `Dev` role against the `/prod/mobile` collection
  only. This team can see and `exec` into their own applications, but not the
  `payments` applications.
- The `payments` team has the `Dev` role for the `/prod/payments` collection
  only.

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
