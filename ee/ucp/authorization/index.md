---
title: Access control model
description: Manage access to resources with role-based access control.
keywords: ucp, grant, role, permission, authentication, authorization, resource, namespace, Kubernetes
---

[Docker Universal Control Plane (UCP)](../index.md),
the UI for [Docker EE](https://www.docker.com/enterprise-edition), lets you
authorize users to view, edit, and use cluster resources by granting role-based
permissions against resource sets.

To authorize access to cluster resources across your organization, UCP
administrators might take the following high-level steps:

- Add and configure **subjects** (users, teams, and service accounts).
- Define custom **roles** (or use defaults) by adding permitted operations per
  type of resource.
- Group cluster **resources** into resource sets of Swarm collections or
  Kubernetes namespaces.
- Create **grants** by combining subject + role + resource set.

For an example, see [Deploy stateless app with RBAC](deploy-stateless-app.md).

## Subjects

A subject represents a user, team, organization, or a service account. A subject
can be granted a role that defines permitted operations against one or more
resource sets.

- **User**: A person authenticated by the authentication backend. Users can
  belong to one or more teams and one or more organizations.
- **Team**: A group of users that share permissions defined at the team level. A
  team can be in one organization only.
- **Organization**: A group of teams that share a specific set of permissions,
  defined by the roles of the organization.
- **Service account**: A Kubernetes object that enables a workload to access
  cluster resources which are assigned to a namespace.

Learn to [create and configure users and teams](create-users-and-teams-manually.md).

## Roles

Roles define what operations can be done by whom. A role is a set of permitted
operations against a type of resource, like a container or volume, which is
assigned to a user or a team with a grant.

For example, the built-in role, **Restricted Control**, includes permissions to
view and schedule nodes but not to update nodes. A custom **DBA** role might
include permissions to `r-w-x` (read, write, and execute) volumes and secrets.

Most organizations use multiple roles to fine-tune the appropriate access. A
given team or user may have different roles provided to them depending on what
resource they are accessing.

Learn to [define roles with authorized API operations](define-roles.md).

## Resource sets

To control user access, cluster resources are grouped into Docker Swarm
*collections* or Kubernetes *namespaces*.

- **Swarm collections**: A collection has a directory-like structure that holds
  Swarm resources. You can create collections in UCP by defining a directory path
  and moving resources into it. Also, you can create the path in UCP and use
  *labels* in your YAML file to assign application resources to the path.
  Resource types that users can access in a Swarm collection include containers,
  networks, nodes, services, secrets, and volumes.

- **Kubernetes namespaces**: A
[namespace](https://v1-11.docs.kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
  is a logical area for a Kubernetes cluster. Kubernetes comes with a `default`
  namespace for your cluster objects, plus two more namespaces for system and
  public resources. You can create custom namespaces, but unlike Swarm
  collections, namespaces _cannot be nested_. Resource types that users can
  access in a Kubernetes namespace include pods, deployments, network policies,
  nodes, services, secrets, and many more.

Together, collections and namespaces are named *resource sets*. Learn to
[group and isolate cluster resources](group-resources.md).

## Grants

A grant is made up of a *subject*, a *role*, and a *resource set*.

Grants define which users can access what resources in what way. Grants are
effectively **Access Control Lists** (ACLs) which
provide comprehensive access policies for an entire organization when grouped
together.

Only an administrator can manage grants, subjects, roles, and access to
resources.

> About administrators
>
> An administrator is a user who creates subjects, groups resources by moving them
> into collections or namespaces, defines roles by selecting allowable operations,
> and applies grants to users and teams.
{: .important}

## Secure Kubernetes defaults

For cluster security, only users and service accounts granted the `cluster-admin` ClusterRole for 
all Kubernetes namespaces via a ClusterRoleBinding can deploy pods with privileged options. This prevents a
platform user from being able to bypass the Universal Control Plane Security Model.

These privileged options include: 

  - `PodSpec.hostIPC` - Prevents a user from deploying a pod in the host's IPC
    Namespace.
  - `PodSpec.hostNetwork` - Prevents a user from deploying a pod in the host's
    Network Namespace.
  - `PodSpec.hostPID` - Prevents a user from deploying a pod in the host's PID
    Namespace.
  - `SecurityContext.allowPrivilegeEscalation` - Prevents a child process
    of a container from gaining more privileges than its parent.
  - `SecurityContext.capabilities` - Prevents additional  [Linux
    Capabilities](https://docs.docker.com/engine/security/security/#linux-kernel-capabilities)
    from being added to a pod.
  - `SecurityContext.privileged` - Prevents a user from deploying a [Privileged
    Container](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities).
  - `Volume.hostPath` - Prevents a user from mounting a path from the host into
    the container. This could be a file, a directory, or even the Docker Socket.

If a user without a cluster admin role tries to deploy a pod with any of these
privileged options, an error similar to the following example is displayed:

```bash
Error from server (Forbidden): error when creating "pod.yaml": pods "mypod" is forbidden: user "<user-id>" is not an admin and does not have permissions to use privileged mode for resource
```

## Where to go next

- [Create and configure users and teams](create-users-and-teams-manually.md)
- [Define roles with authorized API operations](define-roles.md)
- [Grant role-access to cluster resources](grant-permissions.md)
- [Group and isolate cluster resources](group-resources.md)
