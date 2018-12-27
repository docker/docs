---
title: Isolate swarm nodes to a specific team
description: Create grants that limit access to nodes to specific teams.
keywords: ucp, grant, role, permission, authentication
---

With Docker EE Advanced, you can enable physical isolation of resources
by organizing nodes into collections and granting `Scheduler` access for
different users. To control access to nodes, move them to dedicated collections
where you can grant access to specific users, teams, and organizations.

In this example, a team gets access to a node collection and a resource
collection, and UCP access control ensures that the team members can't view
or use swarm resources that aren't in their collection.

You need a Docker EE Advanced license and at least two worker nodes to
complete this example.

1.  Create an `Ops` team and assign a user to it.
2.  Create a `/Prod` collection for the team's node.
3.  Assign a worker node to the `/Prod` collection.
4.  Grant the `Ops` teams access to its collection.

![](../images/isolate-nodes-diagram.svg){: .with-border}

## Create a team

In the web UI, navigate to the **Organizations & Teams** page to create a team
named "Ops" in your organization. Add a user who isn't a UCP administrator to
the team.
[Learn to create and manage teams](create-and-manage-teams.md).

## Create a node collection and a resource collection

In this example, the Ops team uses an assigned group of nodes, which it
accesses through a collection. Also, the team has a separate collection
for its resources.

Create two collections: one for the team's worker nodes and another for the
team's resources.

1.  Navigate to the **Collections** page to view all of the resource
    collections in the swarm.
2.  Click **Create collection** and name the new collection "Prod".
3.  Click **Create** to create the collection.
4.  Find **Prod** in the list, and click **View children**.
5.  Click **Create collection**, and name the child collection
    "Webserver". This creates a sub-collection for access control.

You've created two new collections. The `/Prod` collection is for the worker
nodes, and the `/Prod/Webserver` sub-collection is for access control to
an application that you'll deploy on the corresponding worker nodes.

## Move a worker node to a collection

By default, worker nodes are located in the `/Shared` collection.
Worker nodes that are running DTR are assigned to the `/System` collection.
To control access to the team's nodes, move them to a dedicated collection.

Move a worker node by changing the value of its access label key,
`com.docker.ucp.access.label`, to a different collection.

1.  Navigate to the **Nodes** page to view all of the nodes in the swarm.
2.  Click a worker node, and in the details pane, find its **Collection**.
    If it's in the `/System` collection, click another worker node,
    because you can't move nodes that are in the `/System` collection. By
    default, worker nodes are assigned to the `/Shared` collection.  
3.  When you've found an available node, in the details pane, click
    **Configure**.
3.  In the **Labels** section, find `com.docker.ucp.access.label` and change
    its value from `/Shared` to `/Prod`.
4.  Click **Save** to move the node to the `/Prod` collection.

> Docker EE Advanced required
>
> If you don't have a Docker EE Advanced license, you see the following
> error message when you try to change the access label:
> **Nodes must be in either the shared or system collection without an advanced license.**
> [Get a Docker EE Advanced license](https://www.docker.com/pricing).

![](../images/isolate-nodes-1.png){: .with-border}

## Grant access for a team

You need two grants to control access to nodes and container resources:

-  Grant the `Ops` team the `Restricted Control` role for the `/Prod/Webserver`
   resources.
-  Grant the `Ops` team the `Scheduler` role against the nodes in the `/Prod`
   collection.

Create two grants for team access to the two collections:

1.  Navigate to the **Grants** page and click **Create Grant**.
2.  In the left pane, click **Collections**, and in the **Swarm** collection,
    click **View Children**.
3.  In the **Prod** collection, click **View Children**.
4.  In the **Webserver** collection, click **Select Collection**.  
5.  In the left pane, click **Roles**, and select **Restricted Control**
    in the dropdown.
6.  Click **Subjects**, and under **Select subject type**, click **Organizations**.
7.  Select your organization, and in the **Team** dropdown, select **Ops**.
8.  Click **Create** to grant the Ops team access to the `/Prod/Webserver`
    collection.

The same steps apply for the nodes in the `/Prod` collection.

1.  Navigate to the **Grants** page and click **Create Grant**.
2.  In the left pane, click **Collections**, and in the **Swarm** collection,
    click **View Children**.
3.  In the **Prod** collection, click **Select Collection**.
4.  In the left pane, click **Roles**, and in the dropdown, select **Scheduler**.
5.  In the left pane, click **Subjects**, and under **Select subject type**, click
    **Organizations**.
6.  Select your organization, and in the **Team** dropdown, select **Ops** .
7.  Click **Create** to grant the Ops team `Scheduler` access to the nodes in the
    `/Prod` collection.

![](../images/isolate-nodes-2.png){: .with-border}

## Deploy a service as a team member

Your swarm is ready to show role-based access control in action. When a user
deploys a service, UCP assigns its resources to the user's default collection.
From the target collection of a resource, UCP walks up the ancestor collections
until it finds nodes that the user has `Scheduler` access to. In this example,
UCP assigns the user's service to the `/Prod/Webserver` collection and schedules
tasks on nodes in the `/Prod` collection.

As a user on the Ops team, set your default collection to `/Prod/Webserver`.

1.  Log in as a user on the Ops team.
2.  Navigate to the **Collections** page, and in the **Prod** collection,
    click **View Children**.
3.  In the **Webserver** collection, click the **More Options** icon and
    select **Set to default**.

Deploy a service automatically to worker nodes in the `/Prod` collection.
All resources are deployed under the user's default collection,
`/Prod/Webserver`, and the containers are scheduled only on the nodes under
`/Prod`.

1.  Navigate to the **Services** page, and click **Create Service**.
2.  Name the service "NGINX", use the "nginx:latest" image, and click
    **Create**.
3.  When the **nginx** service status is green, click the service. In the
    details view, click **Inspect Resource**, and in the dropdown, select
    **Containers**.
4.  Click the **NGINX** container, and in the details pane, confirm that its
    **Collection** is **/Prod/Webserver**.

    ![](../images/isolate-nodes-3.png){: .with-border}

5.  Click **Inspect Resource**, and in the dropdown, select **Nodes**.
6.  Click the node, and in the details pane, confirm that its **Collection**
    is **/Prod**.

    ![](../images/isolate-nodes-4.png){: .with-border}

## Alternative: Use a grant instead of the default collection

Another approach is to use a grant instead of changing the user's default
collection. An administrator can create a grant for a role that has the
`Service Create` permission against the `/Prod/Webserver` collection or a child
collection. In this case, the user sets the value of the service's access label,
`com.docker.ucp.access.label`, to the new collection or one of its children
that has a `Service Create` grant for the user.

## Where to go next

- [Node access control in Docker EE Advanced](access-control-node.md)
- [Isolate volumes between two different teams](isolate-volumes-between-teams.md)
- [Deploy a service with view-only access across an organization](deploy-view-only-service.md)

