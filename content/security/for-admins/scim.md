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
>SSO uses Just-in-Time (JIT) provisioning by default. If you [enable SCIM](scim.md#set-up-scim), JIT values still overwrite the attribute values set by SCIM provisioning whenever users log in. To avoid conflicts, make sure your JIT values match your SCIM values. For more information, see [SSO attributes](../for-admins/single-sign-on/_index.md#sso-attributes).
{.important}

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
    - Display name

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
{{< tab name="OneLogin" >}}

1. Go to the OneLogin admin portal.
2. Go to the app that you created when you configured SSO.
3. Go to the **Access** tab and enable the app for the `ol_role` that you created when you configured SSO.
4. Go to the **Provisioning** tab and set up like the following:
    - Select **Enable provisioning**.
    - De-select **Create user**, **Delete user**, **Update user**.
    - When users are deleted in OneLogin, or the user's app access is removed, perform the below action: **Suspend**
    - When user accounts are suspended in OneLogin, perform the following action: **Suspend**
5. Select **Save**.

{{< /tab >}}
{{< /tabs >}}


See the documentation for your IdP for additional details:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm)
- [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/user-provisioning)
- [OneLogin](https://developers.onelogin.com/scim/create-app)

## Set up role mapping

You can assign [roles](/security/for-admins/roles-and-permissions/) to members in your organization in the IdP. To set up a role, you can use optional user-level attributes for the person you want to assign a role. In addition to roles, you can set an organization and team to override the default provisioning values set by the SSO connection.

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

For how to add these attributes, see the documentation for your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-add-custom-user-attributes.htm)
- [Entra ID (formerly Azure AD)](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes#provisioning-a-custom-extension-attribute-to-a-scim-compliant-application)
- [OneLogin](https://onelogin.service-now.com/support?id=kb_article&sys_id=742a000d4740f1909d8dfd1f536d435f&kb_category=566ffd6887332910695f0f66cebb3556#config-info-custom)

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