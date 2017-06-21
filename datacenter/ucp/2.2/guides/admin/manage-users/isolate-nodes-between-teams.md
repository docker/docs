---
title: Isolate swarm nodes between two different teams
description: Create grants that limit access to nodes to specific teams.
keywords: ucp, grant, role, permission, authentication
---

With Docker EE Advanced, you can enable physical isolation of resources
by organizing nodes into collections and granting `Scheduler` access for
different users.

In this example, two teams get access to two different node collections, 
and UCP access control ensures that the teams can't view or use each other's
container resources. You need at least two worker nodes to complete this
example.

1.  Create an `Ops` team and a `QA` team.
2.  Create `/Prod` and `/Staging` collections for the two teams.
3.  Assign worker nodes to one collection or the other.
4.  Grant the `Ops` and `QA` teams access against their
    corresponding collections.

## Create two teams

Click the **Organizations** link in the web UI to create two teams in your
organization, named "Ops" and "QA". For more info, see
[Create and manage teams](create-and-manage-teams.md). 

## Create resource collections

In this example, the Ops and QA teams use two different node groups,
which they access through corresponding resource collections.

1.  In the left pane, click **Collections** to show all of the resource
    collections in the swarm.
2.  Click **Create collection**, and in the **Collection Name** textbox, enter
    "Prod".
3.  Click **Create** to create the collection.
4.  Find **Prod** in the list, and click **View collection**.
5.  Click **Create collection**, and in the **Collection Name** textbox, enter
    "ApplicationA". This creates a sub-collection for access control. 
6.  Navigate to the collections list by clicking **Collections** in the left pane
    or at the top of the page.
7.  Click **Create collection** again, and in the **Collection Name** textbox, enter
    "Staging". Also, create a sub-collection named "ApplicationA".

You've created four new collections. The `/Prod` and `/Staging` collections
are for the worker nodes, and the `/Prod/ApplicationA` and `/Staging/ApplicationA`
sub-collections are for access control to an application that will be deployed on the corresponding worker nodes.

## Move worker nodes to collections 

By default, worker nodes are located in the `/Shared` collection. To control
access to nodes, move them to dedicated collections where you can grant 
access to specific users, teams, and organizations.

Move worker nodes by changing the value of the access label key,
 `com.docker.ucp.access.label`, to a different collection.

1.  In the left pane, click **Nodes** to view all of the nodes in the swarm.
2.  Click a worker node, and in the details pane on the right, click **Edit**. 
3.  In the **Labels** section, find the access label with the value `/Shared` and
    change it to `/Prod`. 
4.  Click **Save** to move the node to the `/Prod` collection.
5.  Repeat the previous steps for another worker node, and move it to the 
    `/Staging` collection. 

> Note: If you're not running Docker EE Advanced, you'll get the following 
> error message when you try to change the access label: 
> Nodes must be in either the shared or system collection without an advanced license.

## Grant access for specific teams

You'll need four grants to control access to nodes and container resources:

-  Grant the `Ops` team the `Scheduler` role against the `/Prod` nodes.
-  Grant the `Ops` team the `Restricted Control` role against the `/Prod/ApplicationA` resources.
-  Grant the `QA` team the `Scheduler` role against the `/Staging` nodes.
-  Grant the `QA` team the `Restricted Control` role against the `/Staging/ApplicationA` resources.

These are the steps for creating the grants for the resource collections.  

1.  Navigate to **User Management > Manage Grants** and click **Create grant**.
2.  In the left pane, click **Collections**, navigate to **/Prod/ApplicationA**,
    and click **Select**.
3.  Click **Roles**, and select **Restricted Control** in the dropdown list.
4.  Click **Subjects**, and under **Select subject type**, click **Organizations**.
    select **Ops** from the **Team** dropdown. 
5.  Click **Create** to grant permissions to the Ops team.
6.  Click **Create grant** and repeat the previous steps for the **/Staging/ApplicationA**
    collection and the QA team.

The same workflow applies for creating the grants against the node collections. 
Apply the `Scheduler` role to the `/Prod` and `/Staging` collections.

With these four grants in place, members of the Staging team won't be able
to view or use the `/Prod` nodes, and members of the Ops team won't be able
to view or use the `/Staging` nodes.

## Access control in action

You can see access control in action with the following two scenarios.

### Create production workloads

Users on the Prod team have permissions to create workloads on the `/Prod`
nodes. 

1.  Log in as a user on the Prod team.
2.  Change the user's default collection to `/Prod/ApplicationA`.
3.  Run `docker stack deploy` with any compose/stack file.
4.  All resources are deployed under `/Prod/ApplicationA`, and the
    containers are scheduled only on the nodes under `/Prod`.

### New users can't inspect isolated nodes and container resources

1.  Create a new user.
2.  Log in as the new user.
3.  Ensure that the `/Shared` collection has at least one worker node.
4.  As the new user, run `docker stack deploy <stack-name>`.
5.  The new workload is deployed on the nodes under `/Shared` and under
    the user's private collection.
6.  The new user can't view any of the nodes under `/Prod` or `/Shared`.

