---
description: Docker Cloud for Organizations and Teams
keywords: organizations, teams, Docker Cloud, resources, permissions
title: Organizations and Teams in Docker Cloud
---

You can create Organizations in Docker Cloud to share repositories, and infrastructure and applications with coworkers and collaborators.

Members of an organization can see only the teams to which they belong, and
their membership. Members of the `Owners` team can see and edit
all of the teams and all of the team membership lists. Docker Cloud users
outside an organization cannot see the Organizations or teams another user
belongs to.

## Create an Organization

An Organization in Docker Cloud contains Teams, and each Team contains users.
You cannot add users directly to an Organization. Organizations can also have
repositories, applications (services and containers), and infrastructure (nodes
and node clusters) associated with them. Paid features such as private
repositories and extra nodes are paid for using the billing information
associated with the Organization.

To create an organization:

1. Log in to Docker Cloud.
2. Select **Create Organization** from the user icon menu at the top right.
2. In the dialog that appears, enter a name for your organization.
3. Enter billing information for the organization.

    This is used for paid features used by the Organization account, including private repositories and additional nodes.

4. Click **Save**.

    The Docker Cloud interface switches you to the new organization view. You
    can return to your individual user account from the menu at the top right
    corner.

When you create an Organization, your user account is automatically added to the
Organization's `Owners` team, which allows you to manage the Organization. This
team must always have at least one member, and you can add other members to it
at any time.

### Convert a user to an Organization

Individual user accounts can be converted to organizations if needed. You can
no longer log in to the account; email addresses, linked source
repositories and collaborators are removed. Automated builds are
migrated. **Account conversion cannot be undone.**

You need another valid Docker ID (not the account you are converting) for
the user who is the first member of the `Owners` team. All existing
automated builds are migrated to this user, and they can configure
the newly converted organization's settings to grant access to other users.

1. Log in to Docker Cloud using the user account that you want to convert.
2. Click **Settings** in the user account menu in the top right corner.
3. Scroll down and click `Convert to organization`.
4. Read through the list of warnings and actions.
5. Enter the Docker ID of the user to be the first member of the Owners team.
6. Click **Save and Continue**.

The UI refreshes. Log in from the Docker ID you specified as the first Owner, and then continue on to configure the organization as described below.

#### What's next?

Once you've created an organization:

* Add users to [the Owners team](orgs.md#configure-the-owners-team) to help you manage the organization
* [Create teams](orgs.md#create-teams)
* [Set team permissions](orgs.md#set-team-permissions)
* Set up [linked providers](orgs.md#link-a-service-provider-to-an-organization), and [manage resources](orgs.md#manage-resources-for-an-organization) for the organization

## Configure the Owners team

Each organization has an `Owners` team which contains the users who manage the
organization's settings. If you created the organization, you are automatically
added to the `Owners` team. You can add new users to the `Owners` team and then
leave the team if you want to transfer ownership. There must always be at least
one member of the `Owners` team.

Owners team members can:

* create, change, and delete teams
* set and change team access permissions
* manage the organization's billing information
* configure the organization's settings (including linked services such as AWS and Github)
* view, change, create and delete repositories, services, and node clusters associated with the organization

> **Note**: You cannot change the Owners team permission settings. Only add users to the Owners team who you are comfortable granting this level of access.

1. While logged in to Docker Cloud, use the menu in the top right corner to switch to the organization you want to work on.
2. Click **Teams** in the lower left corner.
3. Click **owners**.
4. Click **Add user**.
5. Enter the Docker ID of a user to add.
6. Click **Create**.
6. Repeat for each user who you want to add.

To transfer ownership of an organization, add the new owner to the `Owners`
team, then go to your Teams list and click **Leave** on the `Owners` team line.

> **Note**: At this time, only members of the `Owners` team receive email
notifications for events (such as builds and container redeploys) in the
organization's resources. The email "notification level" setting for the
organization affects only the `Owners` team.

## Create teams

You can create Teams within an Organization to add users and manage access to infrastructure, applications, and repositories.

Every organization contains an `Owners` team for users who manage the team
settings. You should create at least one team separate from the owners team so
that you can add members to your organization without giving them this level of
access.

1. While logged in to Docker Cloud, switch to the organization you want to work on from the menu in the upper right corner.
2. Click **Teams** in the lower left corner of the navigation bar.
3. Click **Create** to create a new team.
4. Give the new team a name and description, and click **Create**.
4. On the screen that appears, click **Add User**.
5. Enter the Docker ID of the user and click **Create**.
6. Repeat this process for each user you want to add.

## Set team permissions

You can give Teams within an organization different levels of access to
resources that the organization owns. You can then assign individual users to a
Team to grant them that level of access. Team permissions are set by members of
the `Owners` team.

> **Note**: If a user is a member of multiple teams, their access settings are conjunctive (sometimes called inclusive or additive). For example, if a user is a member of Team A that grants them `No access` to repositories, and they're also a member of Team B that grants them `Read and Write` access to repositories, the user has `Read and Write` access.

To set or edit Team permissions:

1. From the Team detail view, click **Permissions**.
2. Select an access level for `Runtime` resources.
    Runtime resources include both infrastructure and applications.
<!--
Select a default access level for all `Repositories`.
    This access level is applied to any repositories without specific settings in the section below.
Optionally, override the default access level for specific repositories. -->
3. Optionally, grant the team access to one or more repositories in the **Repositories** section.
    1. Enter the name of the repository.
    2. Select an access level.
    3. Click the plus sign (`+`) icon. The change is saved immediately.
    4. Repeat this for each repository that the team needs access to.

    > **Note**: An organization can have public repositories which are visible to **all** users (including those outside the organization). Team members can view public repositories even if you have not given them `View` permission. You can use team permissions to grant write and admin access to public repositories.

### Change team permissions for an individual repository

You can also grant teams access to a repository from the repository's
**Permissions** page rather than from each team's permissions settings. You
might do this if you create repositories after you have already configured your
teams, and want to grant access to several teams at the same time.

If the organization's repository is private, you must explicitly grant any access that your team members require. If the repository is public, all users are granted read-only access by default.

Members of the organization's `Owners` team, and members of any team with `admin` access to the repository can change the repository's access permissions.

To grant a team access to an organization's repository:

1. Navigate to the organization's repository.
2. Click the **Permissions** tab.
3. Select the name of the team you want to add from the drop down menu.
5. Choose the access level the team should have.
6. Click the **plus sign** to add the selected team and permission setting.

     Your choice is saved immediately.

7. Repeat this process for each team to which you want to grant access.

To edit a team's permission level, select a new setting in the **Permission** drop down menu.

To remove a team's access to the repository, click the **trashcan** icon next to the team's access permission line.

> **Note**: If the organization's repository is _public_, team members without explicit access permissions still have read-only access to the repository. If the repository is _private_, removing a team's access completely prevents the team members from seeing the repository.

### Docker Cloud team permission reference

**General access levels**:

* **No access**: no access at all. The resource is not visible to members of this team.
* **Read only**: users can view the resource and its configuration, but cannot perform actions on the resource.
* **Read and Write**: users can view *and change* the resource and its configuration.
* **Admin**: users can view, and edit the resource and its configuration, and can create or delete new instances of the resource.

> **Note**: Only users who are members of the `Owners` team can create _new_ repositories.

| Permission level | Access |
| ------------- | ------------- |
| **Swarms**  (Beta)| |
| Admin | View swarms, manage swarms, add users |
| **Repositories** | |
| Read | Pull |
| Read/Write | Pull, push |
| Admin | All of the above, plus update description, create and delete |
| **Build** | |
| Read | View basic build settings and Timeline |
| Read/write | All of the above plus start, retry, or cancel build |
| Admin | All of the above, plus view and change build configuration, change build source, create and delete |
| **Nodes** | |
| Read | View |
| Read/write  | View, scale, check node health |
| Admin | All of the above plus terminate, upgrade daemon, get certificate, create BYON token, update, deploy, and create |
| **Applications** | |
| Read | View, get logs, export stackfile |
| Read/write  | All of the above, plus start, stop, redeploy, and scale  |
| Admin | All of the above plus, open a terminal window, terminate, update, and create |

## Machine user accounts in organizations

Your organization might find it useful to have a dedicated account that is used for programmatic or scripted access to your organization's resources using the [Docker Cloud APIs](/apidocs/docker-cloud/).

> **Note**: While these accounts are sometimes called "robot" accounts or "bots", these users may not be _created_ using scripts.

To create a "robot" or machine account for your organization:

1. Create a new Docker ID for the machine user. Verify the email address associated with the user.
2. If necessary, create a new Team for the machine user, and grant that team access to the required resources.

    This method is recommended because it makes it easier for administrators to
    understand the machine user's access, and modify it without affecting other
    users' access.

3. Add the machine user to the new Team.

## Modify a team

To modify an existing team, log in to Docker Cloud and switch to your
organization, click **Teams** in the left navigation menu, then click the team
you want to modify.

You can manage team membership from the first page that appears when you select the team.

To change the team name or description, click **Settings**.

To manage team permissions for runtime resources (nodes and applications) and
repositories click **Permissions**.

## Manage resources for an Organization

An organization can have its own resources including repositories, nodes and
node clusters, containers, services, and service stacks, just as if it was a
normal user account.

If you're a member of the `Owners` team, you can create these resources when
logged in as the Organization, and manage which Teams can view, edit, and create
and delete each resource.

#### Link a service provider to an Organization

1. Log in to Docker Cloud as a member of the `Owners` team.

2. Switch to the Organization account by selecting it from the user icon menu at the top right.

3. Click **Cloud Settings** in the left navigation.

    From the Organization's Cloud settings page, you can [link to the organization's source code repositories](builds/link-source.md), [link to infrastructure hosts](infrastructure/index.md) such as a cloud service providers.

    The steps are the same as when you perform these actions as an individual user.

#### Create repositories

When a member of the `Owners` team creates a repository for an organization,
they can configure which teams within the organization can access the
repository. No access controls are configured by default on repository creation.
If the repository is _private_, this leaves it accessible only to members of the
`Owners` team until other teams are granted access.

> **Tip**: Members of the `Owners` team can configure this default from the **Default privacy** section of the organization's **Cloud Settings** page.

1. Log in to Docker Cloud as a member of the `Owners` team.

2. Switch to the Organization account by selecting it from the user icon menu at the top right.

3. [Create the repository](builds/repos.md#create-a-new-repository-in-docker-cloud) as usual.

4. Once the repository has been created, navigate to it and click **Permissions**.

5. [Grant access](#change-team-permissions-for-an-individual-repository) to any teams that require access to the repository.

#### Manage organization settings

From the Organization's **Cloud Settings** page you can also manage the
Organization's Plan and billing account information, notifications, and API
keys.

#### Create organization resources

To create resources for an Organization such as services and node clusters, log
in to Docker Cloud and switch to the Organization account. Create the
repositories, services, stacks, or node clusters as you would for any other
account.
