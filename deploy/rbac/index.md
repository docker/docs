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
- path: /deploy/rbac/usermgmt-create-subjects/
  title: Create and configure users and teams
- path: /deploy/rbac/usermgmt-define-roles/
  title: Define roles with authorized API operations
- path: /deploy/rbac/resources-group-resources/
  title: Group and isolate cluster resources
- path: /deploy/rbac/usermgmt-grant-permissions/
  title: Grant access to cluster resources
---

{% if include.ui %}

The Docker Enterprise Edition UI, or the Docker Univeral Control Plane (UCP),
lets you authorize users to view, edit, and use cluster resources by granting
role-based permissions. Resources can be grouped and isolated inline with an
organization's needs and users can be granted more than one role.


{% if include.version=="ucp-3.0" %}

To authorize access to cluster resources across your organization, Docker EE
administrators might take the following high-level steps:

- Add and configure **subjects** (users and teams).
- Define custom **roles** (or use defaults) by adding permitted operations per
  resource types.
- Group cluster **resources** into Swarm collections or Kubernetes namespaces.
- Create **grants** by marrying subject + role + resource.

For a simple example, see [Deploy a multi-tier application with RBAC](./howto-wordpress-multitier).

## Subjects

A subject represents a user, team, or organization. A subject can be granted a
role that defines permitted operations against one or more resource types.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level. A
  team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

For more, see:
- [Create and configure users and teams](./usermgmt-create-subjects.md)
- [Synchronize users and teams with LDAP](./usermgmt-sync-with-ldap.md)

## Roles

Roles define what operations can be done by whom. A role is a set of permitted
operations against a *resource type* (such as an image, container, volume) that
is assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to
view and schedule a node but not update. A custom **DBA** role might include
permissions to create, attach, view, and remove volumes.

Most organizations use multiple roles to fine-tune the approprate access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

For more, see: [Define roles with authorized API operations](./usermgmt-define-roles.md)

## Resources

Cluster resources are grouped into Swarm collections or Kubernetes namespaces.

A collection is a directory that holds Swarm resources. You can create
collections in the Docker EE UI by defining a path and moving resources. Or you
can create the path in the Docker EE UI and use *labels* in your YAML file to
assign application resources to that path. For an example, see [Deploy a multi-tier service with multiple roles](/deploy/rbac/howto-wordpress-multitier/).

> Resource types that can be placed into a Swarm collection include: Containers,
> Networks, Nodes, Services, Secrets, and Volumes.

A
[namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
is a logical area for a Kubernetes cluster. Kuberenetes comes with a "default"
namespace for your cluster objects (plus two more for system and public
resources).

> Resource types that can be placed into a Kubernetes namespace include: Pods,
> Deployments, NetworkPolcies, Nodes, Services, Secrets, and many more.

For more, see: [Group and isolate cluster resources](./resources-group-resources.md).

## Grants

A grant is made up of *subject*, *role*, and *resource group*.

Grants define which users can access what resources in what way. Grants are
effectively Access Control Lists (ACLs), and when grouped together, can
provide comprehensive access policies for an entire organization.

Only an administrator can manage grants, subjects, roles, and resources.

> Administrators are users who create subjects, group resources by moving them
> into directories or namespaces, define roles by selecting allowable operations,
> and apply grants to users and teams.

For more, see: [Grant access to cluster resources](./usermgmt-grant-permissions.md).


{% elsif include.version=="ucp-2.2" %}

To authorize access to cluster resources across your organization, Docker EE
administrators might take the following high-level steps:

- Add and configure **subjects** (users and teams)
- Define custom **roles** (or use defaults) by adding permissions to resource
  types
- Group cluster **resources** into Swarm collections
- Create **grants** by marrying subject + role + resource

## Subjects

A subject represents a user, team, or organization. A subject is granted a
role that defines permitted operations against one or more resource types.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level. A
  team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

For more, see:
- [Create and configure users and teams](./usermgmt-create-subjects.md)
- [Synchronize users and teams with LDAP](./usermgmt-sync-with-ldap.md)

## Roles

Roles define what operations can be done by whom. A role is a set of permitted
operations against a *resource type* (such as an image, container, volume) that
is assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to
view and schedule a node but not update. A custom **DBA** role might include
permissions to create, attach, view, and remove volumes.

Most organizations use different roles to fine-tune the approprate access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

For more, see: [Define roles with authorized API operations](./usermgmt-define-roles.md)

## Resources

Cluster resources are grouped into Swarm collections.

A collection is a directory that holds Swarm resources. You can create
collections in the Docker EE UI by defining a path and moving resources. Or you
can create the path in the Docker EE UI and use *labels* in your YAML file to
assign application resources to that path. For an example, see [Deploy a multi-tier service with multiple roles](./deploy/rbac/howto-wordpress-multitier/).

> Resource types that can be placed into a Swarm collection include: Containers,
> Networks, Nodes, Services, Secrets, and Volumes.

For more, see: [Group and isolate cluster resources](./resources-group-resources.md).

## Grants

A grant is made up of a *subject*, *resource group*, and *role*.

Grants define which users can access what resources in what way. Grants are
effectively Access Control Lists (ACLs), and when grouped together, can
provide comprehensive access policies for an entire organization.

Only an administrator can manage grants, subjects, roles, and resources.

> Administrators are users who create subjects, group resources by moving them
> into directories or namespaces, define roles by selecting allowable operations,
> and apply grants to users and teams.

For more, see: [Grant access to cluster resources](./usermgmt-grant-permissions.md).

## Transition from UCP 2.1 access control

- Your existing access labels and permissions are migrated automatically during
  an upgrade from UCP 2.1.x.
- Unlabeled "user-owned" resources are migrated into each user's private
  collection, in `/Shared/Private/<username>`.
- Old access control labels are migrated into `/Shared/Legacy/<labelname>`.
- When deploying a resource, choose a collection instead of an access label.
- Use grants for access control, instead of unlabeled permissions.

For hands-on tutorials, see:
- [Deploy a simple Wordpress service with RBAC](./howto-wordpress-view-only.md)
- [Deploy a multi-tier application with RBAC](./howto-wordpress-multitier.md)

{% endif %}
{% endif %}
