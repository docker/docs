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

With Docker Universal Control Plane (UCP), you can authorize how users view,
edit, and use cluster resources by granting role-based permissions. Resources
can be grouped according to an organization's needs and users can be granted
more than one role.

To authorize access to cluster resources across your organization, Docker
administrators might take the following high-level steps:

- Add and configure subjects (users and teams)
- Define customer roles (or use defaults) by adding permissions to resource types
- Group cluster resources into Swarm collections or Kubernetes namespaces
- Create grants by marrying subject + role + resource

## Subjects

A subject represents a user, team, or organization. A subject is granted a
role that defines permitted operations against one or more resources.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level. A
  team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

For more, see:
- [Create users and teams in UCP](./usermgmt-create-subjects.md)
- [Synchronize users and teams with LDAP](./usermgmt-sync-with-ldap.md)

## Roles

Roles define what operations can be done by whom. A role is a set of permitted
operations (view, edit, use, etc.) against a *resource type* (such as an image,
container, volume, etc.) that is assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to
view and schedule a node but not update it. A custom **DBA** role might include
permissions to create, attach, view, and remove volumes.

Most organizations use different roles to fine-tune the approprate access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

For more, see: [Create roles that define user access to resources](./usermgmt-define-roles.md)

## Resources

Cluster resources are grouped into Swarm collections or Kubernetes namespaces.

A collection is a directory that holds Swarm resources. You can define and build
collections in UCP by assinging resources to a collection. Or you can create the
path in UCP and use *labels* in your YAML file to assign application resources to
that path.

> Swarm resources types that can be placed into a collection include: Containers,
> Networks, Nodes, Services, Secrets, and Volumes.

A [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
holds Kubernetes resources.

> Kubernetes resources types that can be placed into a namespace include: Pods,
> Deployments, NetworkPolcies, Nodes, Services, Secrets, and many more.

For more, see: [Group resources into collections or namespaces](resources-group-resources.md).

## Grants

A grant is made up of a *subject*, *resource group*, and *role*.

Grants define which users can access what resources in what way. Grants are
effectively Access Control Lists (ACLs), and when grouped together, can
provide comprehensive access policies for an entire organization.

Only an administrator can manage grants, subjects, roles, and resources.

> Administrators are users who create subjects, group resources by moving them
> into directories or namespaces, define roles by selecting allowable operations,
> and apply grants to users and teams.

For more, see: [Create grants and authorize access to users and teams](usermgmt-grant-permissions.md).


{% elsif include.version=="ucp-2.2" %}

With [Docker Universal Control Plane (UCP)](https://docs.docker.com/datacenter/ucp/2.2/guides/),
you can authorize how users view, edit, and use cluster resources by granting
role-based permissions. Resources can be grouped according to an organization's
needs and users can be granted more than one role.

To authorize access to cluster resources across your organization, Docker
administrators might take the following high-level steps:

- Add and configure subjects (users and teams)
- Define customer roles (or use defaults) by adding permissions to resource types
- Group cluster resources into Swarm collections
- Create grants by marrying subject + role + resource

## Subjects

A subject represents a user, team, or organization. A subject is granted a
role that defines permitted operations against one or more resources.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level. A
  team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

For more, see:
- [Create users and teams in UCP](./usermgmt-create-subjects.md)
- [Synchronize users and teams with LDAP](./usermgmt-sync-with-ldap.md)

## Roles

Roles define what operations can be done by whom. A role is a set of permitted
operations (view, edit, use, etc.) against a *resource type* (such as an image,
container, volume, etc.) that is assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to
view and schedule a node but not update it. A custom **DBA** role might include
permissions to create, attach, view, and remove volumes.

Most organizations use different roles to fine-tune the approprate access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

For more, see [Create roles that define user access to resources](./usermgmt-define-roles.md)

## Resources

Cluster resources are grouped into collections.

A collection is a directory that holds Swarm resources. You can define and build
collections in UCP by assinging resources to a collection. Or you can create the
path in UCP and use *labels* in your YAML file to assign application resources to
that path.

> Swarm resources types that can be placed into a collection include: Containers,
> Networks, Nodes (physical or virtual), Services, Secrets, and Volumes.

For more, see: [Group resources into collections](resources-group-resources.md).

## Grants

A grant is made up of a *subject*, *resource collection*, and *role*.

Grants define which users can access what resources in what way. Grants are
effectively Access Control Lists (ACLs), which when grouped together, can
provide comprehensive access policies for an entire organization.

Only an administrator can manage grants, subjects, roles, and resources.

> Administrators are users who create subjects, group resources by moving them
> into directories or namespaces, define roles by selecting allowable operations,
> and apply grants to users and teams.

For more, see: [Create grants and authorize access to users and teams](usermgmt-grant-permissions.md).

## Transition from UCP 2.1 access control

- Your existing access labels and permissions are migrated automatically during
  an upgrade from UCP 2.1.x.
- Unlabeled "user-owned" resources are migrated into the user's private
  collection, in `/Shared/Private/<username>`.
- Old access control labels are migrated into `/Shared/Legacy/<labelname>`.
- When deploying a resource, choose a collection instead of an access label.
- Use grants for access control, instead of unlabeled permissions.

[See a deeper tutorial on how to design access control architectures.](access-control-design-ee-standard.md)

{% endif %}
{% endif %}
