---
title: Access control model
description: Manage access to resources with role-based access control.
keywords: ucp, grant, role, permission, authentication, authorization, resource, namespace, Kubernetes
ui_tabs:
- version: ucp-3.0
  orlower: false
- version: ucp-2.2
  orlower: false
next_steps:
- path: create-users-and-teams-manually/
  title: Create and configure users and teams
- path: define-roles/
  title: Define roles with authorized API operations
- path: group-resources/
  title: Group and isolate cluster resources
- path: grant-permissions/
  title: Grant role-access to cluster resources
---
{% if include.version=="ucp-3.0" %}

[Docker Universal Control Plane (UCP)](../index.md),
the UI for [Docker EE](https://www.docker.com/enterprise-edition), lets you
authorize users to view, edit, and use cluster resources by granting role-based
permissions against resource types.

To authorize access to cluster resources across your organization, UCP
administrators might take the following high-level steps:

- Add and configure **subjects** (users and teams).
- Define custom **roles** (or use defaults) by adding permitted operations per
  resource types.
- Group cluster **resources** into Swarm collections or Kubernetes namespaces.
- Create **grants** by combining subject + role + resource group.

For an example, see [Deploy stateless app with RBAC](deploy-stateless-app.md).

## Subjects

A subject represents a user, team, or organization. A subject can be granted a
role that defines permitted operations against one or more resource types.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level. A
  team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.

Learn to [create and configure users and teams](create-users-and-teams-manually.md)

## Roles

Roles define what operations can be done by whom. A role is a set of permitted
operations against a *resource type*, like a container or volume, that's
assigned to a user or team with a grant.

For example, the built-in role, **Restricted Control**, includes permission to
view and schedule nodes but not to update nodes. A custom **DBA** role might
include permissions to `r-w-x` volumes and secrets.

Most organizations use multiple roles to fine-tune the appropriate access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

Learn to [define roles with authorized API operations](define-roles.md)

## Resource sets

To control user access, cluster resources are grouped into Docker Swarm *collections*
or Kubernetes *namespaces*.

- **Swarm collections**: A collection has a directory-like structure that holds
  Swarm resources. You can create collections in UCP by defining a directory path
  and moving resources into it. Also, you can create the path in UCP and use
  *labels* in your YAML file to assign application resources to the path.
  Resource types that users can access in a Swarm collection include containers,
  networks, nodes, services, secrets, and volumes.

- **Kubernetes namespaces**: A
[namespace](https://v1-8.docs.kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
  is a logical area for a Kubernetes cluster. Kubernetes comes with a `default`
  namespace for your cluster objects, plus two more namespaces for system and
  public resources. You can create custom namespaces, but unlike Swarm
  collections, namespaces _can't be nested_. Resource types that users can
  access in a Kubernetes namespace include pods, deployments, network policies,
  nodes, services, secrets, and many more.

Together, collections and namespaces are named *resource sets*. Learn to
[group and isolate cluster resources](group-resources.md).

## Grants

A grant is made up of *subject*, *role*, and *resource set*.

Grants define which users can access what resources in what way. Grants are
effectively Access Control Lists (ACLs), and when grouped together, they
provide comprehensive access policies for an entire organization.

Only an administrator can manage grants, subjects, roles, and access to
resources.

> About administrators
>
> An administrator is a user who creates subjects, groups resources by moving them
> into collections or namespaces, defines roles by selecting allowable operations,
> and applies grants to users and teams.

{% elsif include.version=="ucp-2.2" %}

Learn about [access control model in UCP](/datacenter/ucp/2.2/guides/access-control/index.md).

{% endif %}
