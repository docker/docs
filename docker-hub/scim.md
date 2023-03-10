---
description: System for Cross-domain Identity Management
keywords: SCIM, SSO
title: SCIM
direct_from: 
- /docker-hub/company-scim/
- /docker-hub/group-mapping/
---

This section is for administrators who want to enable System for Cross-domain Identity Management (SCIM) 2.0 for their business. It is available for Docker Business customers. 

SCIM provides automated user provisioning and de-provisioning for your Docker organization or company through your identity provider (IdP).  Once you enable SCIM in Docker Hub and your IdP, any user assigned to the Docker application in the IdP is automatically provisioned in Docker Hub and added to the organization or company. 

Similarly, if a user gets unassigned from the Docker application in the IdP, the user is removed from the organization or company in Docker Hub. SCIM also synchronizes changes made to a user's attributes in the IdP, for instance the user’s first name and last name.

The following provisioning features are supported:
 - Creating new users
 - Push user profile updates
 - Remove users
 - Deactivate users 
 - Re-activate users
 - Group mapping
 
The table below lists the supported attributes. Note that your attribute mappings must match for SSO to prevent duplicating your members.

| Attribute    | Description
|:---------------------------------------------------------------|:-------------------------------------------------------------------------------------------|
| username             | Unique identifier of the user (email)                                   |
| givenName                            | User’s first name |
| familyName |User’s surname                                              |

## Set up SCIM

You must make sure you have [configured SSO](../single-sign-on/index.md) before you enable SCIM. Enforcing SSO is not required.

### Step one: Enable SCIM in Docker Hub

1. Sign in to Docker Hub, navigate to the **Organizations** page and select your organization or company.
2. Select **Settings**. If you are setting up SCIM for an organization you then need to select **Security**. 
3. n the **Single Sign-On Connection** table, select the **Actions** icon and **Setup SCIM**.
4. Copy the **SCIM Base URL** and **API Token** and paste the values into your IdP.

### Step two: Enable SCIM in your IdP

Follow the instructions provided by your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm){: target="_blank" rel="noopener" class="_" }
- [Azure AD](https://learn.microsoft.com/en-us/azure/databricks/administration-guide/users-groups/scim/aad#step-2-configure-the-enterprise-application){: target="_blank" rel="noopener" class="_" }
- [OneLogin](https://developers.onelogin.com/scim/create-app){: target="_blank" rel="noopener" class="_" }

### Optional step 
You also have the option to use group mapping within your IdP. To take advantage of group mapping, follow the instructions provided by your IdP:
- [Okta](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-about-group-push.htm){: target="_blank" rel="noopener" class="_" }
- [Azure AD](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes){: target="_blank" rel="noopener" class="_" }
- [OneLogin](https://developers.onelogin.com/scim/create-app){: target="_blank" rel="noopener" class="_" }

Once complete, a user who signs in to Docker through SSO is automatically added to the organizations and teams mapped in the IdP.

## Disable SCIM

If SCIM is disabled, any user provisioned through SCIM will remain in the organization. Future changes for your users will not sync from your IdP. User de-provisioning is only possible when manually removing the user from the organization.

1. In the **Single Sign-On Connection** table, select the **Actions** icon
2. Select **Disable SCIM**.
