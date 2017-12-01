---
title: Group and isolate cluster resources
description: Learn how to group resources into collections or namespaces to control access.
keywords: rbac, ucp, grant, role, permission, authentication, resource collection
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
next_steps:
- path: /deploy/rbac/basics-create-subjects/
  title: Create and configure users and teams
- path: /deploy/rbac/basics-define-roles/
  title: Create roles to authorize access
- path: /deploy/rbac/basics-grant-permissions/
  title: Grant access to cluster resources
---

{% if include.ui %}

{% if include.version=="ucp-3.0" %}

## Kubernetes namespace

A
[namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
is a logical area for a Kubernetes cluster. Kuberenetes comes with a "default"
namespace for your cluster objects (plus two more for system and public
resources). You can create custom namespaces, but unlike Swarm collections,
namespaces _cannot be nested_.

> Resource types that can be placed into a Kubernetes namespace include: Pods,
> Deployments, NetworkPolcies, Nodes, Services, Secrets, and many more.

## Swarm collection

A collection is a directory of grouped resources, such as services, containers,
volumes, networks, and secrets. To authorize access, administrators create
grants against directory branches.

![](../images/collections-and-resources.svg){: .with-border}

## Access label

Access to a collection is granted with a path defined in an access label.

For example, each user has a private collection (with default permisions) and
the path is `/Shared/Private/<username>`. The private collection for user "hans"
would have the following access label:

```
com.docker.ucp.access.label = /Shared/Private/hans
```

## Nested collections

You can nest collections. If a user has a grant against a collection, the grant
applies to all of its child collections.

For a child collection, or for a user who belongs to more than one team, the
system concatenates permissions from multiple roles into an "effective role" for
the user, which specifies the operations that are allowed against the target.

> **Note**: Permissions are concatenated from multiple roles into an "effective
> role".

## Built-in collections

Docker EE provides a number of built-in collections.

| Default collection | Description |
| ------------------ | --------------------------------------------------------------------------------------- |
| `/`                | Path to all resources in the Swarm cluster. Resources not in a collection are put here. |
| `/System`          | Path to UCP managers, DTR nodes, and UCP/DTR system services. By default, only admins have access, but this is configurable. |
| `/Shared`          | Default path to all worker nodes for scheduling. In Docker EE Standard, all worker nodes are located here. In [Docker EE Advanced](https://www.docker.com/enterprise-edition), worker nodes can be moved and [isolated](./howto-isolate-nodes/). |
| `/Shared/Private/` | Path to a user's private collection. |
| `/Shared/Legacy`   | Path to the access control labels of legacy versions (UCP 2.1 and lower). |


This diagram shows the `/System` and `/Shared` collections created by Docker EE.
User private collections are children of the `/Shared/private` collection. The
Docker EE  administrator user created a `/prod` collection and a child
collection, `/webserver`.

![](../images/collections-diagram.svg){: .with-border}

## Default collections

Each user has a default collection which can be changed in UCP preferences.

Users can't deploy a resource without a collection. When a user deploys a
resource in the CLI without an access label, Docker EE automatically places the
resource in the user's default collection.

[Learn how to add labels to nodes](../../datacenter/ucp/2.2/guides/admin/configure/add-labels-to-cluster-nodes/).

With Docker Compose, the system applies default collection labels across all
resources in the stack unless `com.docker.ucp.access.label` has been explicitly
set.

> Default collections and collection labels
>
> Default collections are good for users who ony work on a well-defined slice of
> the system, as well as users who deploy stacks and don't want to edit the
> contents of their compose files. A user with more versatile roles in the
> system, such as an adminitrator, might find it better to set custom labels for
> each resource.

## Collections and labels

Resources are marked as being in a collection by using labels. Some resource
types don't have editable labels, so you can't move resources like this across
collections. You can't modify collections after resource creation for
containers, networks, and volumes, but you can update labels for services,
nodes, secrets, and configs.

For editable resources, like services, secrets, nodes, and configs, you can
change the `com.docker.ucp.access.label` to move resources to different
collections. With the CLI, you can use this label to deploy resources to a
collection other than your default collection. Omitting this label on the CLI
deploys a resource on the user's default resource collection.

The system uses the additional labels, `com.docker.ucp.collection.*`, to enable
efficient resource lookups. By default, nodes have the
`com.docker.ucp.collection.root`, `com.docker.ucp.collection.shared`, and
`com.docker.ucp.collection.swarm` labels set to `true`. UCP
automatically controls these labels, and you don't need to manage them.

Collections get generic default names, but you can give them meaningful names,
like "Dev", "Test", and "Prod".

A *stack* is a group of resources identified by a label. You can place the
stack's resources in multiple collections. Resources are placed in the user's
default collection unless you specify an explicit `com.docker.ucp.access.label`
within the stack/compose file.

## Control access to nodes

The Docker EE Advanced license enables access control on worker nodes. Admin
users can move worker nodes from the default `/Shared` collection into other
collections and create corresponding grants for scheduling tasks.

In this example, an administrator has moved worker nodes to a `/prod`
collection:

![](../images/containers-and-nodes-diagram.svg)

When you deploy a resource with a collection, Docker EE sets a constraint
implicitly based on what nodes the collection, and any ancestor collections, can
access. The `Scheduler` role allows users to deploy resources on a node. By
default, all users have the `Scheduler` role against the `/Shared` collection.

When deploying a resource that isn't global, like local volumes, bridge
networks, containers, and services, the system identifies a set of "schedulable
nodes" for the user. The system identifies the target collection of the
resource, like `/Shared/Private/hans`, and it tries to find the parent that's
closest to the root that the user has the `Node Schedule` permission on.

For example, when a user with a default configuration runs `docker container run
nginx`, the system interprets this to mean, "Create an NGINX container under the
user's default collection, which is at `/Shared/Private/hans`, and deploy it on
one of the nodes under `/Shared`.

If you want to isolate nodes against other teams, place these nodes in new
collections, and assign the `Scheduler` role, which contains the `Node Schedule`
permission, to the team. [Isolate swarm nodes to a specific team](howto-isolate-notes.md).


{% elsif include.version=="ucp-2.2" %}

## Swarm collection

A collection is a directory of grouped resources, such as services, containers,
volumes, networks, and secrets. To authorize access, administrators create
grants against directory branches.

![](../images/collections-and-resources.svg){: .with-border}

## Access label

Access to a collection is granted with a path defined in an access label.

For example, each user has a private collection (with default permisions) and
the path is `/Shared/Private/<username>`. The private collection for user "hans"
would have the following access label:

```
com.docker.ucp.access.label = /Shared/Private/hans
```

## Nested collections

You can nest collections. If a user has a grant against a collection, the grant
applies to all of its child collections.

For a child collection, or for a user who belongs to more than one team, the
system concatenates permissions from multiple roles into an "effective role" for
the user, which specifies the operations that are allowed against the target.

> **Note**: Permissions are concatenated from multiple roles into an "effective
> role".

## Built-in collections

Docker EE provides a number of built-in collections.

| Default collection | Description |
| ------------------ | --------------------------------------------------------------------------------------- |
| `/`                | Path to all resources in the Swarm cluster. Resources not in a collection are put here. |
| `/System`          | Path to UCP managers, DTR nodes, and UCP/DTR system services. By default, only admins have access, but this is configurable. |
| `/Shared`          | Default path to all worker nodes for scheduling. In Docker EE Standard, all worker nodes are located here. In [Docker EE Advanced](https://www.docker.com/enterprise-edition), worker nodes can be moved and [isolated](./howto-isolate-nodes/). |
| `/Shared/Private/` | Path to a user's private collection. |
| `/Shared/Legacy`   | Path to the access control labels of legacy versions (UCP 2.1 and lower). |


This diagram shows the `/System` and `/Shared` collections created by Docker EE.
User private collections are children of the `/Shared/private` collection. The
Docker EE  administrator user created a `/prod` collection and a child
collection, `/webserver`.

![](../images/collections-diagram.svg){: .with-border}

## Default collections

Each user has a default collection which can be changed in UCP preferences.

Users can't deploy a resource without a collection. When a user deploys a
resource in the CLI without an access label, Docker EE automatically places the
resource in the user's default collection.

[Learn how to add labels to nodes](../../datacenter/ucp/2.2/guides/admin/configure/add-labels-to-cluster-nodes/).

With Docker Compose, the system applies default collection labels across all
resources in the stack unless `com.docker.ucp.access.label` has been explicitly
set.

> Default collections and collection labels
>
> Default collections are good for users who ony work on a well-defined slice of
> the system, as well as users who deploy stacks and don't want to edit the
> contents of their compose files. A user with more versatile roles in the
> system, such as an adminitrator, might find it better to set custom labels for
> each resource.

## Collections and labels

Resources are marked as being in a collection by using labels. Some resource
types don't have editable labels, so you can't move resources like this across
collections. You can't modify collections after resource creation for
containers, networks, and volumes, but you can update labels for services,
nodes, secrets, and configs.

For editable resources, like services, secrets, nodes, and configs, you can
change the `com.docker.ucp.access.label` to move resources to different
collections. With the CLI, you can use this label to deploy resources to a
collection other than your default collection. Omitting this label on the CLI
deploys a resource on the user's default resource collection.

The system uses the additional labels, `com.docker.ucp.collection.*`, to enable
efficient resource lookups. By default, nodes have the
`com.docker.ucp.collection.root`, `com.docker.ucp.collection.shared`, and
`com.docker.ucp.collection.swarm` labels set to `true`. UCP
automatically controls these labels, and you don't need to manage them.

Collections get generic default names, but you can give them meaningful names,
like "Dev", "Test", and "Prod".

A *stack* is a group of resources identified by a label. You can place the
stack's resources in multiple collections. Resources are placed in the user's
default collection unless you specify an explicit `com.docker.ucp.access.label`
within the stack/compose file.

## Control access to nodes

The Docker EE Advanced license enables access control on worker nodes. Admin
users can move worker nodes from the default `/Shared` collection into other
collections and create corresponding grants for scheduling tasks.

In this example, an administrator has moved worker nodes to a `/prod`
collection:

![](../images/containers-and-nodes-diagram.svg)

When you deploy a resource with a collection, Docker EE sets a constraint
implicitly based on what nodes the collection, and any ancestor collections, can
access. The `Scheduler` role allows users to deploy resources on a node. By
default, all users have the `Scheduler` role against the `/Shared` collection.

When deploying a resource that isn't global, like local volumes, bridge
networks, containers, and services, the system identifies a set of "schedulable
nodes" for the user. The system identifies the target collection of the
resource, like `/Shared/Private/hans`, and it tries to find the parent that's
closest to the root that the user has the `Node Schedule` permission on.

For example, when a user with a default configuration runs `docker container run
nginx`, the system interprets this to mean, "Create an NGINX container under the
user's default collection, which is at `/Shared/Private/hans`, and deploy it on
one of the nodes under `/Shared`.

If you want to isolate nodes against other teams, place these nodes in new
collections, and assign the `Scheduler` role, which contains the `Node Schedule`
permission, to the team. [Isolate swarm nodes to a specific team](howto-isolate-notes.md).

{% endif %}
{% endif %}
