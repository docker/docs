---
title: Isolate volumes between two different teams
description: Create grants that limit access to volumes to specific teams.
keywords: ucp, grant, role, permission, authentication
---

In this example, two teams are granted access to volumes in two different
resource collections. UCP access control prevents the teams from viewing and
accessing each other's volumes, even though they may be located in the same
nodes.

The procedure has the following steps.

1.  Create two teams
2.  Create resource collections
3.  Create grants
4.  Team member creates volumes

## Create two teams

Use the **Organizations** web UI to create two teams in your organization, 
named "Dev" and "Prod". 
[Learn how to create and manage teams](create-and-manage-teams.md). 

## Create resource collections

In this example, the Dev and Prod teams use two different volumes, which they 
access through two corresponding resource collections. The collections are
placed under the `/Shared` collection.

1.  In the left pane, click **Collections** to show all of the resource
    collections in the swarm.
2.  Find the **/Shared** collection and click **View collection**.
2.  Click **Create collection**, and in the **Collection Name** input, enter
    "dev-volumes".
3.  Click **Create** to create the collection.
4.  Click **Create collection** again, and in the **Collection Name** input, enter
    "prod-volumes", and click **Create**.

## Create grants for controlling access to the new volumes

1.  Navigate to **User Management > Manage Grants** and click **Create grant**.
2.  In the left pane, click **Collections**, navigate to **/Shared/dev-volumes**,
    and click **Select**.
3.  Click **Roles**, and select **Restricted Control** in the dropdown list.
4.  Click **Subjects**, and under **Select subject type**, click **Organizations**.
    In the dropdown, pick the organization that you want to associate with this grant.
    Also, pick **Dev** from the **Team** dropdown. 
5.  Click **Create** to grant permissions to the Dev team.
6.  Click **Create grant** and repeat the previous steps for the **/Shared/prod-volumes**
    collection and the Prod team.

## Create a volume as a team member

Team members have permission to create volumes in their assigned collection.

1.  Log in as one of the users on the Dev team.
2.  In the left pane, click **Volumes** to show all of the 
    volumes in the swarm that the user can access.
2.  Click **Create volume** and name the new volume "dev-data".
3.  In the left pane, click **Collections**. The default collection appears.
    At the top of the page, click **Shared**, find the **dev-volumes**
    collection in the list, and click **Select**.
4.  Click **Create** to add the "dev-data" volume to the collection. 


