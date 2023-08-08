{% if include.product == "admin" %}
  {% assign product_link = "[Docker Admin](https://admin.docker.com)" %}
  {% if include.layer == "company" %}
    {% assign sso_link = "[configured SSO](/admin/company/settings/sso-configuration/)" %}
    {% assign sso_navigation="Select your company in the left navigation drop-down menu, and then select **SSO & SCIM.**" %}
  {% else %}
    {% assign sso_link = "[configured SSO](/admin/organization/security-settings/sso-configuration/)" %}
    {% assign sso_navigation="Select your organization in the left navigation drop-down menu, and then select **SSO & SCIM.**" %}
  {% endif %}
{% else %}
  {% assign product_link = "[Docker Hub](https://hub.docker.com)" %}
  {% assign sso_link = "[configured SSO](/single-sign-on/configure/)" %}
  {% assign sso_navigation="Navigate to the SSO settings page for your organization or company.
    - Organization: Select **Organizations**, your organization, **Settings**, and then **Security**.
    - Company: Select **Organizations**, your company, and then **Settings**." %}
{% endif %}

This section is for administrators who want to enable System for Cross-domain Identity Management (SCIM) 2.0 for their business. It is available for Docker Business customers.

SCIM provides automated user provisioning and de-provisioning for your Docker organization or company through your identity provider (IdP).  Once you enable SCIM in Docker and your IdP, any user assigned to the Docker application in the IdP is automatically provisioned in Docker and added to the organization or company.

Similarly, if a user gets unassigned from the Docker application in the IdP, the user is removed from the organization or company in Docker. SCIM also synchronizes changes made to a user's attributes in the IdP, for instance the user’s first name and last name.

The following provisioning features are supported:
 - Creating new users
 - Push user profile updates
 - Remove users
 - Deactivate users
 - Re-activate users
 - Group mapping

The following table lists the supported attributes. Note that your attribute mappings must match for SSO to prevent duplicating your members.

| Attribute    | Description
|:---------------------------------------------------------------|:-------------------------------------------------------------------------------------------|
| userName             | User's primary email address. This is used as the unique identifier of the user. |
| name.givenName | User’s first name |
| name.familyName | User’s surname |
| active | Indicates if a user is enabled or disabled. Can be set to false to de-provision the user. |

For additional details about supported attributes and SCIM, see [Docker Hub API SCIM reference](/docker-hub/api/latest/#tag/scim).

## Set up SCIM

You must make sure you have {{ sso_link }} before you enable SCIM. Enforcing SSO is not required.

### Step one: Enable SCIM in Docker

1. Sign in to {{ product_link }}{: target="_blank" rel="noopener" class="_" }.
2. {{ sso_navigation }}
3. In the SSO connections table, select the **Actions** icon and **Setup SCIM**.
4. Copy the **SCIM Base URL** and **API Token** and paste the values into your IdP.

### Step two: Enable SCIM in your IdP

Follow the instructions provided by your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm){: target="_blank" rel="noopener" class="_" }
- [Azure AD](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/user-provisioning){: target="_blank" rel="noopener" class="_" }
- [OneLogin](https://developers.onelogin.com/scim/create-app){: target="_blank" rel="noopener" class="_" }

## Disable SCIM

If SCIM is disabled, any user provisioned through SCIM will remain in the organization. Future changes for your users will not sync from your IdP. User de-provisioning is only possible when manually removing the user from the organization.

1. Sign in to {{ product_link }}{: target="_blank" rel="noopener" class="_" }.
2. {{ sso_navigation }}
3. In the SSO connections table, select the **Actions** icon.
4. Select **Disable SCIM**.

## Limitations

Administrators can assign [roles](/docker-hub/roles-and-permissions/) to organization members. However, SCIM doesn't support role management.
