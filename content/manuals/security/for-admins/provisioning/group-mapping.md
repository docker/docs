---
description: Group mapping for administrators
keywords: Group Mapping, SCIM, Docker Hub, Docker Admin, admin, security
title: Group mapping
aliases:
- /admin/company/settings/group-mapping/
- /admin/organization/security-settings/group-mapping/
- /docker-hub/group-mapping/
- /security/for-admins/group-mapping/
weight: 40
---

{{< summary-bar feature_name="SSO" >}}

Group mapping lets you sync user groups from your identity provider (IdP) with teams in your Docker organization. This automates team membership management, keeping your Docker teams up to date based on changes in your IdP. You can use group mapping once you have configured [single sign-on (SSO)](../single-sign-on/_index.md).

> [!TIP]
>
> Group mapping is ideal for adding users to multiple organizations or multiple teams within one organization. If you don't need to set up multi-organization or multi-team assignment, you can use SCIM [user-level attributes](scim.md#set-up-role-mapping).

## How group mapping works

With group mapping enabled, when a user authenticates through SSO, your IdP shares key attributes with Docker, such as the user's email address, name, and groups. Docker uses these attributes to create or update the user's profile, as well as to manage their team and organization assignments. With group mapping, users’ team memberships in Docker automatically reflect changes made in your IdP groups.

It's important to note that Docker uses the user's email address as a unique identifier. Each Docker account must always have a unique email address.

## Use group mapping

To assign users to Docker teams through your IdP, you must create groups in your IdP following the naming pattern: `organization:team`. For example, if your organization is called "moby" and you want to manage the "developers" team, the group name in your IdP should be `moby:developers`. In this example, any user added to this group in your IdP is automatically assigned to the "developers" team in Docker.

You can also use this format to assign users to multiple organizations. For example, to add a user to the "backend" team in the "moby" organization and the "desktop" team in the "whale" organization, the group names would be `moby:backend` and `whale:desktop`.

> [!TIP]
>
> Match the group names in your IdP with your Docker teams. When groups are synced, Docker creates a team if it doesn’t already exist.

The following lists the supported group mapping attributes:

| Attribute | Description |
|:--------- | :---------- |
| id | Unique ID of the group in UUID format. This attribute is read-only. |
| displayName | Name of the group following the group mapping format: `organization:team`. |
| members | A list of users that are members of this group. |
| members(x).value | Unique ID of the user that is a member of this group. Members are referenced by ID. |

The general steps to use group mapping are:

1. In your IdP, create groups with the `organization:team` format.
2. Add users to the group.
3. Add the Docker application that you created in your IdP to the group.
4. Add attributes in the IdP.
5. Push groups to Docker.

The exact configuration may vary depending on your IdP. You can use [group mapping with SSO](#use-group-mapping-with-sso), or with SSO and [SCIM enabled](#use-group-mapping-with-scim).

### Use group mapping with SSO

The following steps describe how to set up and use group mapping with SSO
connections that use the SAML authentication method. Note that group mapping
with SSO isn't supported with the Azure AD (OIDC) authentication method.
Additionally, SCIM isn't required for these configurations.

{{< tabs >}}
{{< tab name="Okta" >}}

The user interface for your IdP may differ slightly from the following steps. You can refer to the [Okta documentation](https://help.okta.com/oie/en-us/content/topics/apps/define-group-attribute-statements.htm) to verify.

To set up group mapping:

1. Sign in to Okta and open your application.
2. Navigate to the **SAML Settings** page for your application.
3. In the **Group Attribute Statements (optional)** section, configure like the following:
   - **Name**: `groups`
   - **Name format**: `Unspecified`
   - **Filter**: `Starts with` + `organization:` where `organization` is the name of your organization
   The filter option will filter out the groups that aren't affiliated with your Docker organization.
4. Create your groups by selecting **Directory**, then **Groups**.
5. Add your groups using the format `organization:team` that matches the names of your organization(s) and team(s) in Docker.
6. Assign users to the group(s) that you create.

The next time you sync your groups with Docker, your users will map to the Docker groups you defined.

{{< /tab >}}
{{< tab name="Entra ID" >}}

The user interface for your IdP may differ slightly from the following steps. You can refer to the [Entra ID documentation](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes) to verify.

To set up group mapping:

1. Sign in to Entra ID and open your application.
2. Select **Manage**, then **Single sign-on**.
3. Select **Add a group claim**.
4. In the Group Claims section, select **Groups assigned to the application** with the source attribute **Cloud-only group display names (Preview)**.
5. Select **Advanced options**, then the **Filter groups** option.
6. Configure the attribute like the following:
   - **Attribute to match**: `Display name`
   - **Match with**: `Contains`
   - **String**: `:`
7. Select **Save**.
8. Select **Groups**, **All groups**, then **New group** to create your group(s).
9. Assign users to the group(s) that you create.

The next time you sync your groups with Docker, your users will map to the Docker groups you defined.

{{< /tab >}}
{{< /tabs >}}

### Use group mapping with SCIM

The following steps describe how to set up and use group mapping with SCIM. Before you begin, make sure you [set up SCIM](./scim.md#enable-scim) first.

{{< tabs >}}
{{< tab name="Okta" >}}

The user interface for your IdP may differ slightly from the following steps. You can refer to the [Okta documentation](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-enable-group-push.htm) to verify.

To set up your groups:

1. Sign in to Okta and open your application.
2. Select **Applications**, then **Provisioning**, and **Integration**.
3. Select **Edit** to enable groups on your connection, then select **Push groups**.
4. Select **Save**. Saving this configuration will add the **Push Groups** tab to your application.
5. Create your groups by navigating to **Directory** and selecting **Groups**.
6. Add your groups using the format `organization:team` that matches the names of your organization(s) and team(s) in Docker.
7. Assign users to the group(s) that you create.
8. Return to the **Integration** page, then select the **Push Groups** tab to open the view where you can control and manage how groups are provisioned.
9. Select **Push Groups**, then **Find groups by rule**.
10. Configure the groups by rule like the following:
    - Enter a rule name, for example `Sync groups with Docker Hub`
    - Match group by name, for example starts with `docker:` or contains `:` for multi-organization
    - If you enable **Immediately push groups by rule**, sync will happen as soon as there's a change to the group or group assignments. Enable this if you don't want to manually push groups.

Find your new rule under **By rule** in the **Pushed Groups** column. The groups that match that rule are listed in the groups table on the right-hand side.

To push the groups from this table:

1. Select **Group in Okta**.
2. Select the **Push Status** drop-down.
3. Select **Push Now**.

{{< /tab >}}
{{< tab name="Entra ID" >}}

The user interface for your IdP may differ slightly from the following steps. You can refer to the [Entra ID documentation](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes) to verify.

Complete the following before configuring group mapping:

1. Sign in to Entra ID and go to your application.
2. In your application, select **Provisioning**, then **Mappings**.
3. Select **Provision Microsoft Entra ID Groups**.
4. Select **Show advanced options**, then **Edit attribute list**.
5. Update the `externalId` type to `reference`, then select the **Multi-Value** checkbox and choose the referenced object attribute `urn:ietf:params:scim:schemas:core:2.0:Group`.
6. Select **Save**, then **Yes** to confirm.
7. Go to **Provisioning**.
8. Toggle **Provision Status** to **On**, then select **Save**.

Next, set up group mapping:

1. Go to the application overview page.
2. Under **Provision user accounts**, select **Get started**.
3. Select **Add user/group**.
4. Create your group(s) using the `organization:team` format.
5. Assign the group to the provisioning group.
6. Select **Start provisioning** to start the sync.

To verify, select **Monitor**, then **Provisioning logs** to see that your groups were provisioned successfully. In your Docker organization, you can check that the groups were correctly provisioned and the members were added to the appropriate teams.

{{< /tab >}}
{{< /tabs >}}

Once complete, a user who signs in to Docker through SSO is automatically added to the organizations and teams mapped in the IdP.

> [!TIP]
>
> [Enable SCIM](scim.md) to take advantage of automatic user provisioning and de-provisioning. If you don't enable SCIM users are only automatically provisioned. You have to de-provision them manually.

## More resources

The following videos demonstrate how to use group mapping with your IdP with SCIM enabled:

- [Video: Group mapping with Okta](https://youtu.be/c56YECO4YP4?feature=shared&t=3023)
- [Video: Attribute and group mapping with Entra ID (Azure)](https://youtu.be/bGquA8qR9jU?feature=shared&t=2039)
