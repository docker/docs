---
title: Grant permissions to users based on roles
description: Grant access to swarm resources by using role-based access control.
keywords: ucp, grant, role, permission, authentication, authorization
---

If you're a UCP administrator, you can create *grants* to control how users 
and organizations access swarm resources.

A grant is made up of a *subject*, a *role*, and a *resource collection*.
A grant defines who (subject) has how much access (role) 
to a set of resources (collection). Each grant is a 1:1:1 mapping of 
subject, role, collection. For example, you can grant the "Prod Team" 
"View Only" permissions for the "/Production" collection.

The usual workflow for creating grants has four steps.

1.  Set up your users and teams. For example, you might want three teams,
    Dev, QA, and Prod.
2.  Organize swarm resources into separate collections that each team uses.
3.  Optionally, create custom roles for specific permissions to the Docker API.
4.  Grant role-based access to collections for your teams.

## Create a grant

When you have your users, collections, and roles set up, you can create
grants. Administrators create grants on the **Manage Grants** page.

1.  Click **Create Grant**. The default collection, usually `/Swarm` is listed.
2.  Click **Selected** to list all of the collections.
3.  Click **Select** on the collection you want to grant access to.
4.  In the left pane, click **Roles** and select a role from the dropdown list.
5.  In the left pane, click **Subjects**. Click **All Users** to create a grant
    for a specific user, or click **Organizations** to create a grant for an
    organization or a team.
6.  Select a user, team, or organization and click **Create**. 

