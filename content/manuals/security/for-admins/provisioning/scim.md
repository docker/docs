---
keywords: SCIM, SSO, user provisioning, de-provisioning, role mapping, assign users
title: SCIM provisioning
linkTitle: SCIM
description: Learn how System for Cross-domain Identity Management works and how to set it up.
aliases:
  - /security/for-admins/scim/
  - /docker-hub/scim/
weight: 30
---

{{< summary-bar feature_name="SSO" >}}

System for Cross-domain Identity Management (SCIM) is available for Docker Business customers. This guide provides an overview of SCIM provisioning.

## How SCIM works

SCIM offers automated user provisioning and de-provisioning for Docker through your identity provider (IdP). Once SCIM is enabled, users assigned to the Docker application in your IdP are automatically provisioned and added to your Docker organization. If a user is unassigned, they are removed from Docker.

SCIM also syncs user profile updates, such as name changes, made in your IdP. SCIM can be used with Docker’s default Just-in-Time (JIT) provisioning configuration, or on its own with JIT disabled.

SCIM supports the automation of:
- Creating users
- Updating user profiles
- Removing and deactivating users
- Re-activating users
- Group mapping

## Supported attributes

> [!IMPORTANT]
>
> Docker uses JIT provisioning by default for SSO configurations. If you enable SCIM, JIT values still overwrite the attribute
values set by SCIM provisioning. To avoid conflicts, your JIT attribute values must match your SCIM attribute values. To avoid conflicts between SCIM and JIT, you can also disable JIT provisioning. See [Just-in-Time](/manuals/security/for-admins/provisioning/just-in-time.md) for more information.

Attributes are pieces of user information, such as name and email, that are synchronized between your IdP and Docker when using SCIM. Proper mapping of these attributes is essential for seamless user provisioning and to prevent duplicate entries when using SSO.

The following table lists the supported attributes for SCIM:

| Attribute    | Description |
|:---------------------------------------------------------------|:-------------------------------------------------------------------------------------------|
| userName             | User’s primary email address, used as the unique identifier |
| name.givenName | User’s first name |
| name.familyName | User’s surname |
| active | Indicates if a user is enabled or disabled, set to “false” to de-provision a user |

For additional details about supported attributes and SCIM, see [Docker Hub API SCIM reference](/reference/api/hub/latest/#tag/scim).

## Enable SCIM in Docker

{{< summary-bar feature_name="Admin console early access" >}}

You must [configure SSO](../single-sign-on/configure/_index.md) before you enable SCIM. Enforcing SSO isn't required to use SCIM.

{{< tabs >}}
{{< tab name="Admin Console" >}}

{{% admin-scim product="admin" %}}

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% admin-scim %}}

{{< /tab >}}
{{< /tabs >}}

## Enable SCIM in your IdP

The user interface for your IdP may differ slightly from the following steps. You can refer to the documentation for your IdP to verify. For additional details, see the documentation for your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm)
- [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/user-provisioning)

{{< tabs >}}
{{< tab name="Okta" >}}

### Enable SCIM

1. Sign in to Okta and select **Admin** to open the admin portal.
2. Open the application you created when you configured your SSO connection.
3. On the application page, select the **General** tab, then **Edit App Settings**.
4. Enable SCIM provisioning, then select **Save**.
5. Now you can access the **Provisioning** tab in Okta. Navigate to this tab, then select **Edit SCIM Connection**.
6. To configure SCIM in Okta, set up your connection using the following values and settings:
    - SCIM Base URL: SCIM connector base URL (copied from Docker Hub)
    - Unique identifier field for users: `email`
    - Supported provisioning actions: **Push New Users** and **Push Profile Updates**
    - Authentication Mode: HTTP Header
    - SCIM Bearer Token: HTTP Header Authorization Bearer Token (copied from Docker Hub)
7. Select **Test Connector Configuration**.
8. Review the test results and select **Save**.

### Enable synchronization

1. In Okta, select **Provisioning**.
2. Select **To App**, then **Edit**.
3. Enable **Create Users**, **Update User Attributes**, and **Deactivate Users**.
4. Select **Save**.
5. Remove unnecessary mappings. The necessary mappings are:
    - Username
    - Given name
    - Family name
    - Email

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. In the Azure admin portal, go to **Enterprise Applications**, then select the **Docker** application you created when you set up your SSO connection.
2. Select **Provisioning**, then **Get Started**.
3. Select **Automatic** provisioning mode.
4. Enter the **SCIM Base URL** and **API Token** from Docker into the **Admin Credentials** form.
5. Test the connection, then select **Save**.
6. Go to  **Mappings**, then select **Provision Azure Active Directory Groups**.
7. Set the **Enabled** value to **No**.
8. Select **Provision Azure Active Directory Users**.
9. Remove all unsupported attributes.
10. Select **Save**.
11. Set the provisioning status to **On**.

{{< /tab >}}
{{< /tabs >}}

## Set up role mapping

You can assign [roles](/security/for-admins/roles-and-permissions/) to members in your organization in your IdP. To set up a role, you can use optional user-level attributes for the person you want to assign a role. In addition to roles, you can set an organization or team to override the default provisioning values set by the SSO connection.

> [!NOTE]
>
> Role mappings are supported for both SCIM and JIT provisioning. With JIT provisioning, role mapping only applies when a user is initially provisioned to the organization.

The following table lists the supported optional user-level attributes.

| Attribute | Possible values    | Considerations |
| --------- | ------------------ | -------------- |
| `dockerRole` | `member`, `editor`, or `owner`, for a list of permissions for each role, see [Roles and permissions](/security/for-admins/roles-and-permissions/) | If you don't assign a role in the IdP, the value of the `dockerRole` attribute defaults to `member`. When you set the attribute, this overrides the default value. |
| `dockerOrg` | `organizationName`, for example, an organization named "moby" would be `moby` | Setting this attribute overrides the default organization configured by the SSO connection. Also, this won't add the user to the default team. If this attribute isn't set, the user is provisioned to the default organization and the default team. If set and `dockerTeam` is also set, this provisions the user to the team within that organization. |
| `dockerTeam` | `teamName`, for example, a team named "developers" would be `developers` | Setting this attribute provisions the user to the default organization and to the specified team, instead of the SSO connection's default team. This also creates the team if it doesn't exist. You can still use group mapping to provision users to teams in multiple organizations. See [Group mapping](/security/for-admins/provisioning/group-mapping/) for more details. |

After you set the role in the IdP, you must initiate a sync in your IdP to push the changes to Docker.

The external namespace to use to set up these attributes is `urn:ietf:params:scim:schemas:extension:docker:2.0:User`.

{{< tabs >}}
{{< tab name="Okta" >}}

### Set up role mapping in Okta

1. Setup [SSO](../single-sign-on/configure/_index.md) and SCIM first.
2. In the Okta admin portal, go to **Directory**, select **Profile Editor**, and then **User (Default)**.
3. Select **Add Attribute** and configure the values for the role, organization, or team you want to add. Exact naming isn't required.
4. Return to the **Profile Editor** and select your application.
5. Select **Add Attribute** and enter the required values. The **External Name** and **External Namespace** must be exact. The external name values for organization/team/role mapping are `dockerOrg`, `dockerTeam`, and `dockerRole` respectively, as listed in the previous table. The external namespace is the same for all of them: `urn:ietf:params:scim:schemas:extension:docker:2.0:User`.
6. After creating the attributes, navigate to the top of the page and select **Mappings**, then **Okta User to YOUR APP**.
7. Go to the newly created attributes and map the variable names to the external names, then select **Save Mappings**. If you’re using JIT provisioning, continue to the following steps.
8. Navigate to **Applications** and select **YOUR APP**.
9. Select **General**, then **SAML Settings**, and **Edit**.
10. Select **Step 2** and configure the mapping from the user attribute to the Docker variables.

### Assign roles by user

1. In the Okta admin portal, select **Directory**, then **People**.
2. Select **Profile**, then **Edit**.
3. Select **Attributes** and update the attributes to the desired values.

### Assign roles by group

1. In the Okta admin portal, select **Directory**, then **People**.
2. Select **YOUR GROUP**, then **Applications**.
3. Open **YOUR APPLICATION** and select the **Edit** icon.
4. Update the attributes to the desired values.

If a user doesn't already have attributes set up, users who are added to the group will inherit these attributes upon provisioning.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

### Set up role mapping in Azure AD

1. Setup [SSO](../single-sign-on/configure/_index.md) and SCIM first.
2. In the Azure AD admin portal, open **Enterprise Apps** and select **YOUR APP**.
3. Select **Provisioning**, then **Mappings**, and **Provision Azure Active Directory Users**.
4. To set up the new mapping, check **Show advanced options**, then select **Edit attribute options**.
5. Create new entries with the desired mapping for role, organization, or group (for example, `urn:ietf:params:scim:schemas:extension:docker:2.0:User:dockerRole`) as a string type.
6. Navigate back to **Attribute Mapping** for users and select **Add new mapping**.

### Expression mapping

This implementation works best for roles, but can't be used along with organization and team mapping using the same method. With this approach, you can assign attributes at a group level, which members can inherit. This is the recommended approach for role mapping.

1. In the **Edit Attribute** view, select the **Expression** mapping type.
2. If you can create app roles named as the role directly (for example, `owner` or `editor`), in the **Expression** field, you can use `SingleAppRoleAssignment([appRoleAssignments])`.

   Alternatively, if you’re restricted to using app roles you have already defined (for example, `My Corp Administrators`) you’ll need to setup a switch for these roles. For example:

    ```text
    Switch(SingleAppRoleAssignment([appRoleAssignments]), "member", "My Corp Administrator", "owner", "My Corp Editor", "editor")`
    ```
3. Set the following fields:
    - **Target attribute**: `urn:ietf:params:scim:schemas:extension:docker:2.0:User:dockerRole`
    - **Match objects using this attribute**: No
    - **Apply this mapping**: Always
4. Save your configuration.

### Direct mapping

Direct mapping is an alternative to expression mapping. This implementation works for all three mapping types at the same time. In order to assign users, you'll need to use the Microsoft Graph API.

1. In the **Edit Attribute** view, select the **Direct** mapping type.
2. Set the following fields:
    - **Source attribute**: choose one of the allowed extension attributes in Entra (for example, `extensionAttribute1`)
    - **Target attribute**: `urn:ietf:params:scim:schemas:extension:docker:2.0:User:dockerRole`
    - **Match objects using this attribute**: No
    - **Apply this mapping**: Always

    If you're setting more than one attribute, for example role and organization, you need to choose a different extension attribute for each one.
3. Save your configuration.

### Assign users

If you used expression mapping in the previous step, navigate to **App registrations**, select **YOUR APP**, and **App Roles**. Create an app role for each Docker role. If possible, create it with a display name that is directly equivalent to the role in Docker, for example, `owner` instead of `Owner`. If set up this way, then you can use expression mapping to `SingleAppRoleAssignment([appRoleAssignments])`. Otherwise, a custom switch will have to be used. See [Expression mapping](#expression-mapping).

To add a user:
1. Select **YOUR APP**, then **Users and groups**.
2. Select **Add user/groups**, select the user you want to add, then **Select** their desired role.

To add a group:
1. Select **YOUR APP**, then **Users and groups**.
2. Select **Add user/groups**, select the user you want to add, then **Select** their desired role.

If you used direct mapping in the previous step, go to **Microsoft Graph Explorer** and sign in to your tenant. You need to be a tenant admin to use this feature. Use the Microsoft Graph API to assign the extension attribute to the user with the value that corresponds to what the attribute was mapped to. See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/extensibility-overview?tabs=http) on adding or updating data in extension attributes.

{{< /tab >}}
{{< /tabs >}}

See the documentation for your IdP for additional details:

- [Okta](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-add-custom-user-attributes.htm)
- [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes#provisioning-a-custom-extension-attribute-to-a-scim-compliant-application)

## Disable SCIM

{{< summary-bar feature_name="Admin console early access" >}}

If SCIM is disabled, any user provisioned through SCIM will remain in the organization. Future changes for your users will not sync from your IdP. User de-provisioning is only possible when manually removing the user from the organization.

{{< tabs >}}
{{< tab name="Admin Console" >}}

{{% admin-scim-disable product="admin" %}}

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% admin-scim-disable %}}

{{< /tab >}}
{{< /tabs >}}

## More resources

The following videos demonstrate how to configure SCIM for your IdP:

- [Video: Configure SCIM with Okta](https://youtu.be/c56YECO4YP4?feature=shared&t=1314)
- [Video: Attribute mapping with Okta](https://youtu.be/c56YECO4YP4?feature=shared&t=1998)
- [Video: Configure SCIM with Entra ID (Azure)](https://youtu.be/bGquA8qR9jU?feature=shared&t=1668)
- [Video: Attribute and group mapping with Entra ID (Azure)](https://youtu.be/bGquA8qR9jU?feature=shared&t=2039)
