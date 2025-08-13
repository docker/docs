---
title: Group mapping
description: Automate team membership by syncing identity provider groups with Docker teams
keywords: Group Mapping, SCIM, Docker Admin, admin, security, team management, user provisioning, identity provider
aliases:
- /admin/company/settings/group-mapping/
- /admin/organization/security-settings/group-mapping/
- /docker-hub/group-mapping/
- /security/for-admins/group-mapping/
- /security/for-admins/provisioning/group-mapping/
weight: 30
---

{{< summary-bar feature_name="SSO" >}}

Group mapping automatically synchronizes user groups from your identity provider (IdP) with teams in your Docker organization. For example, when you add a developer to the "backend-team" group in your IdP, they're automatically added to the corresponding team in Docker

This page explains how group mapping works, and how to set up group mapping.

> [!TIP]
>
> Group mapping is ideal for adding users to multiple organizations or multiple teams within one organization. If you don't need to set up multi-organization or multi-team assignment, SCIM [user-level attributes](scim.md#set-up-role-mapping) may be a better fit for your needs.

## Prerequisites

Before you being, you must have:

- SSO configured for your organization
- Administrator access to Docker Home and your identity provider

## How group mapping works

Group mapping keeps your Docker teams synchronized with your IdP groups through these key components:

- Authentication flow: When users sign in through SSO, your IdP shares user attributes with Docker including email, name, and group memberships.
- Automatic updates: Docker uses these attributes to create or update user profiles and manage team assignments based on IdP group changes.
- Unique identification: Docker uses email addresses as unique identifiers, so each Docker account must have a unique email address.
- Team synchronization: Users' team memberships in Docker automatically reflect changes made in your IdP groups.

## Set up group mapping

Group mapping setup involves configuring your identity provider to share group
information with Docker. This requires:

- Creating groups in your IdP using Docker's naming format
- Configuring attributes so your IdP sends group data during authentication
- Adding users to the appropriate groups
- Testing the connection to ensure groups sync properly

You can use group mapping with SSO only, or with both SSO and SCIM for enhanced
user lifecycle management.

### Group naming format

Create groups in your IdP using the format: `organization:team`.

For example:

- For the "developers" team in the "moby" organization: `mobdy:developers`
- For multi-organization access: `moby:backend` and `whale:desktop`

Docker creates teams automatically if they don't already exist when groups sync.

### Supported attributes

| Attribute | Description |
|:--------- | :---------- |
| `id` | Unique ID of the group in UUID format. This attribute is read-only. |
| `displayName` | Name of the group following the group mapping format: `organization:team`. |
| `members` | A list of users that are members of this group. |
| `members(x).value` | Unique ID of the user that is a member of this group. Members are referenced by ID. |

## Configure group mapping with SSO

Use group mapping with SSO connections that use the SAML authentication method.

> [!NOTE]
>
> Group mapping with SSO isn't supported with the Azure AD (OIDC) authentication method. SCIM isn't required for these configurations.

{{< tabs >}}
{{< tab name="Okta" >}}

The user interface for your IdP may differ slightly from the following steps. Refer to the [Okta documentation](https://help.okta.com/oie/en-us/content/topics/apps/define-group-attribute-statements.htm) to verify.

To set up group mapping:

1. Sign in to Okta and open your application.
1. Navigate to the **SAML Settings** page for your application.
1. In the **Group Attribute Statements (optional)** section, configure like the following:
   - **Name**: `groups`
   - **Name format**: `Unspecified`
   - **Filter**: `Starts with` + `organization:` where `organization` is the name of your organization
   The filter option will filter out the groups that aren't affiliated with your Docker organization.
1. Create your groups by selecting **Directory**, then **Groups**.
1. Add your groups using the format `organization:team` that matches the names of your organization(s) and team(s) in Docker.
1. Assign users to the group(s) that you create.

The next time you sync your groups with Docker, your users will map to the Docker groups you defined.

{{< /tab >}}
{{< tab name="Entra ID" >}}

The user interface for your IdP may differ slightly from the following steps. Refer to the [Entra ID documentation](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes) to verify.

To set up group mapping:

1. Sign in to Entra ID and open your application.
1. Select **Manage**, then **Single sign-on**.
1. Select **Add a group claim**.
1. In the Group Claims section, select **Groups assigned to the application** with the source attribute **Cloud-only group display names (Preview)**.
1. Select **Advanced options**, then the **Filter groups** option.
1. Configure the attribute like the following:
   - **Attribute to match**: `Display name`
   - **Match with**: `Contains`
   - **String**: `:`
1. Select **Save**.
1. Select **Groups**, **All groups**, then **New group** to create your group(s).
1. Assign users to the group(s) that you create.

The next time you sync your groups with Docker, your users will map to the Docker groups you defined.

{{< /tab >}}
{{< /tabs >}}

## Configure group mapping with SCIM

Use group mapping with SCIM for more advanced user lifecycle management. Before you begin, make sure you [set up SCIM](./scim.md#enable-scim) first.

{{< tabs >}}
{{< tab name="Okta" >}}

The user interface for your IdP may differ slightly from the following steps. Refer to the [Okta documentation](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-enable-group-push.htm) to verify.

To set up your groups:

1. Sign in to Okta and open your application.
1. Select **Applications**, then **Provisioning**, and **Integration**.
1. Select **Edit** to enable groups on your connection, then select **Push groups**.
1. Select **Save**. Saving this configuration will add the **Push Groups** tab to your application.
1. Create your groups by navigating to **Directory** and selecting **Groups**.
1. Add your groups using the format `organization:team` that matches the names of your organization(s) and team(s) in Docker.
1. Assign users to the group(s) that you create.
1. Return to the **Integration** page, then select the **Push Groups** tab to open the view where you can control and manage how groups are provisioned.
1. Select **Push Groups**, then **Find groups by rule**.
1. Configure the groups by rule like the following:
    - Enter a rule name, for example `Sync groups with Docker Hub`
    - Match group by name, for example starts with `docker:` or contains `:` for multi-organization
    - If you enable **Immediately push groups by rule**, sync will happen as soon as there's a change to the group or group assignments. Enable this if you don't want to manually push groups.

Find your new rule under **By rule** in the **Pushed Groups** column. The groups that match that rule are listed in the groups table on the right-hand side.

To push the groups from this table:

1. Select **Group in Okta**.
1. Select the **Push Status** drop-down.
1. Select **Push Now**.

{{< /tab >}}
{{< tab name="Entra ID" >}}

The user interface for your IdP may differ slightly from the following steps. Refer to the [Entra ID documentation](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes) to verify.

Complete the following before configuring group mapping:

1. Sign in to Entra ID and go to your application.
1. In your application, select **Provisioning**, then **Mappings**.
1. Select **Provision Microsoft Entra ID Groups**.
1. Select **Show advanced options**, then **Edit attribute list**.
1. Update the `externalId` type to `reference`, then select the **Multi-Value** checkbox and choose the referenced object attribute `urn:ietf:params:scim:schemas:core:2.0:Group`.
1. Select **Save**, then **Yes** to confirm.
1. Go to **Provisioning**.
1. Toggle **Provision Status** to **On**, then select **Save**.

Next, set up group mapping:

1. Go to the application overview page.
1. Under **Provision user accounts**, select **Get started**.
1. Select **Add user/group**.
1. Create your group(s) using the `organization:team` format.
1. Assign the group to the provisioning group.
1. Select **Start provisioning** to start the sync.

To verify, select **Monitor**, then **Provisioning logs** to see that your groups were provisioned successfully. In your Docker organization, you can check that the groups were correctly provisioned and the members were added to the appropriate teams.

{{< /tab >}}
{{< /tabs >}}

Once complete, a user who signs in to Docker through SSO is automatically added to the organizations and teams mapped in the IdP.

> [!TIP]
>
> [Enable SCIM](scim.md) to take advantage of automatic user provisioning and de-provisioning. If you don't enable SCIM users are only automatically provisioned. You have to de-provision them manually.
