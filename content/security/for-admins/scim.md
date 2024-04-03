---
keywords: SCIM, SSO, user provisioning, de-provisioning, role mapping, assign users
title: SCIM overview
description: Learn how System for Cross-domain Identity Management works and how to set it up.

direct_from:
- /docker-hub/company-scim/
- /docker-hub/scim/
- /admin/company/settings/scim/
- /admin/organization/security-settings/scim/
---

This section is for administrators who want to enable System for Cross-domain Identity Management (SCIM) 2.0 for their business. It is available for Docker Business customers.

SCIM provides automated user provisioning and de-provisioning for your Docker organization or company through your identity provider (IdP).  Once you enable SCIM in Docker and your IdP, any user assigned to the Docker application in the IdP is automatically provisioned in Docker and added to the organization or company.

Similarly, if a user gets unassigned from the Docker application in the IdP, this removes the user from the organization or company in Docker. SCIM also synchronizes changes made to a user's attributes in the IdP, for example the user’s first name and last name.

The following lists the supported provisioning features:
 - Creating new users
 - Push user profile updates
 - Remove users
 - Deactivate users
 - Re-activate users
 - Group mapping

## Supported attributes

The following table lists the supported attributes. Note that your attribute mappings must match for SSO to prevent duplicating your members.

| Attribute    | Description |
|:---------------------------------------------------------------|:-------------------------------------------------------------------------------------------|
| userName             | User's primary email address. This is the unique identifier of the user. |
| name.givenName | User’s first name |
| name.familyName | User’s surname |
| active | Indicates if a user is enabled or disabled. Can be set to false to de-provision the user. |

For additional details about supported attributes and SCIM, see [Docker Hub API SCIM reference](/docker-hub/api/latest/#tag/scim).

> **Important**
>
> SSO uses Just-in-Time (JIT) provisioning by default. If you [enable SCIM](scim.md#set-up-scim), JIT values still overwrite the attribute values set by SCIM provisioning whenever users log in. To avoid conflicts, make sure your JIT values match your SCIM values. For more information, see [SSO attributes](../for-admins/single-sign-on/configure/configure-idp.md#sso-attributes).
{.important}

> **Beta feature**
>
> Optional Just-in-Time (JIT) provisioning is available in Private Beta when you use the Admin Console. If you're participating in this program, you can avoid conflicts between SCIM and JIT by disabling JIT provisioning in your SSO connection.
{ .experimental }

## Enable SCIM in Docker

You must make sure you have [configured SSO](single-sign-on/configure/_index.md) before you enable SCIM. Enforcing SSO isn't required.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-scim %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-scim product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## Enable SCIM in your IdP

The user interface for your IdP may differ slightly from the following steps. You can refer to the documentation for your IdP to verify.

{{< tabs >}}
{{< tab name="Okta" >}}

### Enable SCIM

1. Go to the Okta admin portal.
2. Go to the app you created when you configured your SSO connection.
3. On the app page, go to the **General** tab and select **Edit App Settings**.
4. Enable SCIM provisioning, then select **Save**.
5. Now you can access the **Provisioning** tab. Navigate to this tab, then select **Edit SCIM Connection**.
6. To configure SCIM in Okta, set up your connection like the following:
    - SCIM Base URL: SCIM connector base URL (copied from Docker Hub)
    - Unique identifier field for users: `email`
    - Supported provisioning actions: **Push New Users** and **Push Profile Updates**
    - Authentication Mode: HTTP Header
    - SCIM Bearer Token: HTTP Header Authorization Bearer Token (copied from Docker Hub)
7. Select **Test Connector Configuration**.
8. Review the test results.
9. Select **Save**.

### Enable synchronization

1. Go to **Provisioning > To App > Edit**.
2. Enable **Create Users**, **Update User Attributes**, and **Deactivate Users**.
3. Select **Save**.
4. Remove unnecessary mappings. The necessary mappings are:
    - Username
    - Given name
    - Family name
    - Email

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. In the Azure admin portal, go to **Enterprise Applications**, then select the **Docker** application you created when you set up your SSO connection.
2. Go to **Provisioning** and select **Get Started**.
3. Select **Automatic** provisioning mode.
4. Enter the SCIM Base URL and API Token from Docker Hub into the **Admin Credentials** form.
5. Test the connection, then select **Save**.
6. Go to  **Mappings** , then select **Provision Azure Active Directory Groups**.
7. Set the **Enabled** value to **No**.
8. Select **Provision Azure Active Directory Users**.
9. Remove all unsupported attributes.
10. Select **Save**.
11. Set the provisioning status to **On**.

{{< /tab >}}
{{< /tabs >}}


See the documentation for your IdP for additional details:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm)
- [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/user-provisioning)

## Set up role mapping

You can assign [roles](/security/for-admins/roles-and-permissions/) to members in your organization in the IdP. To set up a role, you can use optional user-level attributes for the person you want to assign a role. In addition to roles, you can set an organization or team to override the default provisioning values set by the SSO connection.

> **Note**
>
> These mappings are supported for both SCIM and JIT provisioning. With JIT provisioning, role mapping only applies when a user is initially provisioned to the organization.

The following table lists the supported optional user-level attributes.

| Attribute | Possible values    | Considerations |
| --------- | ------------------ | -------------- |
| `dockerRole` | `member`, `editor`, or `owner`. For a list of permissions for each role, see [Roles and permissions](/security/for-admins/roles-and-permissions/). | If you don't assign a role in the IdP, the value of the `dockerRole` attribute defaults to `member`. When you set the attribute, this overrides the default value. |
| `dockerOrg` | `organizationName`. For example, an organization named "moby" would be `moby`. | Setting this attribute overrides the default organization configured by the SSO connection. Also, this won't add the user to the default team. If this attribute isn't set, the user is provisioned to the default organization and the default team. If set and `dockerTeam` is also set, this provisions the user to the team within that org. |
| `dockerTeam` | `teamName`. For example, a team named "developers" would be `developers`. | Setting this attribute provisions the user to the default org and to the specified team, instead of the SSO connection's default team. This also creates the team if it doesn't exist. You can still use group mapping to provision users to teams in multiple orgs. See [Group mapping](/security/for-admins/group-mapping/). |

After you set the role in the IdP, you need to sync to push the changes to Docker.

The external namespace to use to set up these attributes is `urn:ietf:params:scim:schemas:extension:docker:2.0:User`.

{{< tabs >}}
{{< tab name="Okta" >}}

### Set up

1. Setup [SSO](./single-sign-on/configure/_index.md) and SCIM first.
2. In the Okta admin portal, go to **Directory > Profile Editor** and select **User (Default)**.
3. Select **Add Attribute** and configure the values for the role, org, or team you want to add. Exact naming isn't required.
4. Return to the **Profile Editor** and select your application.
5. Select **Add Attribute** and enter the required values. The **External Name** and **External Namespace** must be exact. The external name values for org/team/role mapping are `dockerOrg`, `dockerTeam`, and `dockerRole` respectively, as listed in the previous table. The external namespace is the same for all of them: `urn:ietf:params:scim:schemas:extension:docker:2.0:User`.
6. After creating the attributes, go to the top and select **Mappings > Okta User to YOUR APP**.
7. Go to the newly created attributes and map the variable names selected above to the external names, then select **Save Mappings**. If you’re using JIT provisioning, continue to the following step.
8. Go to **Applications > YOUR APP > General > SAML Settings > Edit > Step 2** and configure the mapping from the user attribute to the docker variables.

### Assign roles by user

1. Go to **Directory > People > YOUR USER > Profile**, then select **Edit** on **Attributes**.
2. Update the attributes to the desired values.

### Assign roles by group

1. Go to **Directory > People > YOUR GROUP > Applications > YOUR APPLICATION**, then select the **Edit** icon.
2. Update the attributes to the desired values.

If a user doesn't already have attributes set up, users who are added to the group will inherit these attributes upon provsioning.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

### Set up

1. Setup [SSO](./single-sign-on/configure/_index.md) and SCIM first.
2. In the Azure AD admin portal, go to**Enterprise Apps > YOUR APP > Provisioning > Mappings > Provision Azure Active Directory Users**.
3. To set up the new mapping, check **Show advanced options**, then select **Edit attribute options**.
4. Create new entries with the desired mapping for role, org, or group (for example, `urn:ietf:params:scim:schemas:extension:docker:2.0:User:dockerRole`) as a string type.
5. Go back to **Attribute Mapping** for users and click **Add new mapping**.

### Expression mapping

This implementation works best for roles, but can't be used along with organization and team mapping using the same method. With this approach, you can assign attributes at a group level, which members can inherit. This is the recommended approach for role mapping.

1. In the **Edit Attribute** view, select the **Expression** mapping type.
2. If you can create app roles named as the role directly (for example, `owner` or `editor`), in the **Expression** field, you can use `SingleAppRoleAssignment([appRoleAssignments])`.

   Alternatively, if you’re restricted to using app roles you have already defined (for example, `My Corp Administrators`) you’ll need to setup a switch for these roles. For example:

    ```text
    Switch(SingleAppRoleAssignment([appRoleAssignments]), "member", "My Corp Administrator", "owner", "My Corp Editor", "editor")`
    ```
3. Set the following fields:
    - **Target attribute**: `urn:ietf:params:scim:schemas:extension:docker:2.0:User:dockerRole`.
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

If you used expression mapping in the previous step, go to **App registrations > YOUR APP > App Roles** and create an app role for each Docker role. If possible, create it with a display name that is directly equivalent to the role in Docker, for example, `owner` instead of `Owner`. If set up this way, then you can use expression mapping to `SingleAppRoleAssignment([appRoleAssignments])`. Otherwise, a custom switch will have to be used. See [Expression mapping](#expression-mapping).

To add a user:
1. Go to **YOUR APP > Users and groups**. Select **Add user/group**.
2. Select the user you want to add, then **Select** their desired role.

To add a group:
1. Go to **YOUR APP > Users and groups**. Select **Add user/group**.
2. Select the group you want to add, then **Select** the desired role for the users in that group.

If you used direct mapping in the previous step, go to **Microsoft Graph Explorer** and sign in to your tenant. You need to be a tenant admin to use this feature. Use the Microsoft Graph API to assign the extension attribute to the user with the value that corresponds to what the attribute was mapped to. See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/extensibility-overview?tabs=http) on adding or updating data in extension attributes.

{{< /tab >}}
{{< /tabs >}}

See the documentation for your IdP for additional details:

- [Okta](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-add-custom-user-attributes.htm)
- [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes#provisioning-a-custom-extension-attribute-to-a-scim-compliant-application)

## Disable SCIM

If SCIM is disabled, any user provisioned through SCIM will remain in the organization. Future changes for your users will not sync from your IdP. User de-provisioning is only possible when manually removing the user from the organization.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-scim-disable %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-scim-disable product="admin" %}}

{{< /tab >}}
{{< /tabs >}}