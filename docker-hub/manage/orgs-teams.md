---
title: Create organizations and teams
description: Control access to repos with user organizations and teams
keywords: Docker Hub, registry, organizations, teams, resources, permissions
redirect_from:
- /docker-hub/orgs/
- /docker-cloud/orgs/
---

To share and restrict access to repositories in Docker Hub, create user
organizations and teams.

Members of an organization can only see the teams to which they belong and the
membership of those teams. Members of the `Owners` team can see and edit all
teams and all team membership lists.

Docker Hub users outside of an organization cannot see the organizations or
teams of other users.

## Create an organization

An organization is a group of teams, and a team is a group of users. You cannot
add users directly to an organization.

Organizations can have repositories and images associated with them. Paid
features such as private repositories are purchased with the billing information
associated with the organization.

To create an organization:

1.  Log in to Docker Hub.

2.  Select **Create Organization** from the user icon menu at the top right.

3.  Enter a name for your organization in the dialog.

4.  Enter billing information for the organization (for paid features, such as
    private repositories).

5.  Click **Save**.

    The Docker Hub interface changes to the new organization view. Use the menu
    at the top right to return to your individual user account.

When you create an organization, your user account is automatically added to the
`Owners` team of that organization, allowing you to manage the organization.

The `Owners` team must always have at least one member, and you can add other
members to it at any time.

### Convert user acount to organization

Individual user accounts can be converted to organizations if needed -- but they
canont be converted back so be careful. Also, create a new Docker ID before
converting.

> Account conversion cannot be undone
>
> Account conversion cannot be undone! Alos, after converting, you cannot log in
> to the _original account_. Email addresses, linked source repositories, and
> collaborators are removed. Automated builds are migrated.
{: .warning}

All existing automated builds are migrated to the first member of the `Owners`
team of the new organization (which you specify if the procedure below). This
person can configure the newly converted organization settings to grant access
to other users.

1. Log in to Docker Hub using the user account that you want to convert.

2. Click **Settings** in the user account menu in the top right corner.

3. Scroll down and click `Convert to organization`.

4. Read through the list of warnings and actions.

5. Enter the Docker ID of the user to be the first member of the Owners team.

6. Click **Save and Continue**.

The UI refreshes. Log in from the Docker ID you specified as the first Owner,
and then continue on to configure the organization as described below.

## Configure the Owners team

Each organization has an `Owners` team with members who manage the settings of
the organization. There must always be at least one member of the `Owners` team.

If you created the organization, you are automatically added to the `Owners`
team. You can add new members and also leave the team if you want to transfer
ownership.

Owners team members can:

- Create, edit, and delete teams.
* Configure and edit team access permissions.
* Manage billing information for the organization.
* Configure the organization settings (including linked services such as AWS and Github).
* Create, edit, and delete repositories associated with the organization.

> You cannot change the Owners team permission settings. Only add users to the
> Owners team who you are comfortable granting this level of access.

1.  Select an organization from the menu in the top right corner of the UI.

2.  Click **Teams** in the left navigation panel.

3.  Click **owners**.

4.  Click **Add user**.

5.  Enter the Docker ID of a user to add.

6.  Click **Create**.

7.  Repeat for each user who you want to add.

To transfer ownership of an organization, add the new owner to the `Owners`
team, then go to your Teams list and click **Leave** on the `Owners` team line.

> Email notifications for Owners
>
> Only members of the `Owners` team receive email notifications for events (such
> as automated builds) in the organization's resources. The email "notification
> level" setting for the organization affects only the `Owners` team.

## Create teams

You can create teams within an Organization to add users and manage access to
repositories.

Every organization contains an `Owners` team for users who manage the team
settings. You should create at least one team separate from the owners team so
that you can add members to your organization without giving them this level of
access.

1.  Select an organization from the menu in the top right corner of the UI.

2.  Click **Teams** in the left navigation panel.

3.  Click **Create** to create a new team.

4.  Give the new team a name and description, and click **Create**.

5.  On the screen that appears, click **Add User**.

6.  Enter the Docker ID of the user and click **Create**.

7.  Repeat this process for each user you want to add.

## Configure team permissions

You can give Teams within an organization different levels of access to
resources that the organization owns.

You can then assign individual users to a Team to grant them that level of
access. Team permissions are set by members of the `Owners` team.

> Additive permissions
>
> If a user is a member of multiple teams, access is conjunctive (inclusive or
> additive). For example, if Team A grants Alice `No access` to repositories,
> and Team B grants her `Read and Write` access, she has `Read and Write` access.

To set or edit Team permissions:

1.  From the Team detail view, click **Permissions**.

<!--
Select a default access level for all `Repositories`.
    This access level is applied to any repositories without specific settings in the section below.
Optionally, override the default access level for specific repositories. -->
2.  Grant the team access to one or more repositories in the **Repositories**
    section.

    a. Enter the name of the repository.

    b. Select an access level.

    c. Click the plus sign (`+`) icon. The change is saved immediately.

    d. Repeat this for each repository that the team needs access to.

    > Repo visibility
    >
    > An organization can have public repositories which are visible to **all**
    > users (including those outside the organization). Team members can view
    > public repositories even if you have not given them `View` permission. You
    > can use team permissions to grant write and admin access to public
    > repositories.

### Edit permissions for individual repos

You can also grant teams access to a repository from the repository's
**Permissions** page rather than from each team's permissions settings. You
might do this if you create repositories after you have already configured your
teams, and want to grant access to several teams at the same time.

If the organization's repository is private, you must explicitly grant any
access that your team members require. If the repository is public, all users
are granted read-only access by default.

Members of the organization's `Owners` team, and members of any team with
`admin` access to the repository can change the repository's access permissions.

To grant a team access to an organization's repository:

1.  Navigate to the organization's repository.

2.  Click the **Permissions** tab.

3.  Select the name of the team you want to add from the drop down menu.

4.  Choose the access level the team should have.

5.  Click the **plus sign** to add the selected team and permission setting. Your choice is saved immediately.

6.  Repeat this process for each team to which you want to grant access.

To edit a team's permission level, select a new setting in the **Permission**
drop down menu.

To remove a team's access to the repository, click the **trashcan** icon next to
the team's access permission line.

> Public vs private
>
> If the organization's repository is _public_, team members without explicit
> access permissions still have read-only access to the repository. If the
> repository is _private_, removing a team's access completely prevents the team
> members from seeing the repository.

### Permissions reference for teams

**Team access levels**:

* **No access**: no access at all. The resource is not visible to members of this team.
* **Read only**: users can view the resource and its configuration, but cannot perform actions on the resource.
* **Read and Write**: users can view _and change_ the resource and its configuration.
* **Admin**: users can view, and edit the resource and its configuration, and can create or delete new instances of the resource.

> Only users who are members of the `Owners` team can create _new_ repositories.

| Permission level  | Access                                                     |
| ----------------- | ---------------------------------------------------------- |
| **Repositories**  |                                                            |
| Read              | Pull                                                       |
| Read/Write        | Pull, push                                                 |
| Admin             | Pull, push, update description, create and delete          |
|                   |                                                            |
| **Build**         |                                                            |
| Read              | View build settings and timeline                           |
| Read/write        | View build settings and timeline, start/retry/cancel build |
| Admin             | View build settings and timeline, start/retry/cancel/change build configuration and source, create and delete |

## Machine user accounts in organizations

Your organization might find it useful to have a dedicated account for
programmatic or scripted access to your organization's resources using the
[Docker Hub APIs](/apidocs/docker-Hub/).

> These users may not be _created_ using scripts (even though these accounts are
> referred to as "robot" accounts or "bots").

To create a "robot" or machine account for your organization:

1.  Create a new Docker ID for the machine user. Verify the email address associated with the user.

2.  If necessary, create a new Team for the machine user, and grant that team access to the required resources.

    This method is recommended because it makes it easier for administrators to
    understand the machine user's access, and modify it without affecting other
    users' access.

3.  Add the machine user to the new Team.

## Modify a team

To modify an existing team, log in to Docker Hub and switch to your
organization, click **Teams** in the left navigation menu, then click the team
you want to modify.

You can manage team membership from the first page that appears when you select the team.

To change the team name or description, click **Settings**.

To manage team permissions for runtime resources (nodes and applications) and
repositories click **Permissions**.

## Manage resources for an organization

An organization can have its own resources including repositories, nodes and
node clusters, containers, services, and service stacks, just as if it was a
normal user account.

If you're a member of the `Owners` team, you can create these resources when
logged in as the Organization, and manage which Teams can view, edit, and create
and delete each resource.

### Link a service provider to an organization

1.  Log in to Docker Hub as a member of the `Owners` team.

2.  Switch to the Organization account by selecting it from the user icon menu at the top right.

3.  Click **Hub Settings** in the left navigation.

    From the Organization's Hub settings page, you can link to the organization source code repositories in [GitHub](../build/github/) or [Bitbucket](../build/bitbucket/).

    The steps are the same as when you perform these actions as an individual user.

### Create repositories

When a member of the `Owners` team creates a repository for an organization,
they can configure which teams within the organization can access the
repository. No access controls are configured by default on repository creation.
If the repository is _private_, this leaves it accessible only to members of the
`Owners` team until other teams are granted access.

> **Tip**:
>
> Members of the `Owners` team can configure this default from the
> **Default privacy** section of the organization's **Hub Settings** page.

See [Create new user repository](../manage/repos#create-new-user-repository)

1.  Log in to Docker Hub as a member of the `Owners` team.

2.  Switch to the Organization account by selecting it from the user icon menu at the top right.

3.  [Create the repository](../manage/repos#create-new-user-repository) as usual.

4.  Once the repository has been created, navigate to it and click **Permissions**.

5.  [Grant access](#configure-team-permissions) to any teams that require access to the repository.

### Manage organization settings

From the Organization's **Hub Settings** page you can also manage the
organization's Plan and billing account information, notifications, and API
keys.
