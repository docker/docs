---
title: Deploy a service with view-only access across an organization
description: Create a grant to control access to a service.
keywords: ucp, grant, role, permission, authentication
---

In this example, your organization is granted access to a new resource
collection that contains one service. If you don't have an organization 
already, create one by navigating to **User Management > Organizations** 
and clicking **Create organization**. 

1.  In the left pane, click **Collections** to show all of the resource
    collections in the swarm.
2.  Find the **Shared** collection and click **View collection**.
2.  Click **Create collection**, and in the **Collection Name** textbox, enter
    a name that describes the resources that you want to group. In this example,
    name the collection "View-only services".
3.  Click **Create** to create the collection.

Currently, the new collection is empty. To populate it, deploy a new service
and add it to the collection.

1.  In the left pane, click **Services** to show all of the services running 
    in the swarm.
2.  Click **Create service**, and in the **Name** textbox, enter "WordPress".
3.  In the **Image** textbox, enter "wordpress". This identifies the latest
    `wordpress` image in the Docker Store.
4.  In the left pane, click **Collections**. The user's default collection
    appears.
    Click **Selected** to list all of the collections. Click **Shared**,
    find the **View-only services** collection in the list, and click
    **Select**.
5.  Click **Deploy** to add the "WordPress" service to the collection and
    deploy it.

You're ready to create a grant for controlling access to the "HelloWorld" service.

1.  Navigate to **User Management > Manage Grants** and click **Create grant**.
2.  In the left pane, click **Collections**, navigate to **/Shared/View-only services**,
    and click **Select**.
3.  Click **Roles**, and select **View Only** in the dropdown list.
4.  Click **Subjects**, and under **Select subject type**, click **Organizations**.
    In the dropdown, pick the organization that you want to associate with this grant. 
5.  Click **Create** to grant permissions to the organization.
