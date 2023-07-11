---
description: System for Cross-domain Identity Management
keywords: SCIM, SSO
title: SCIM
direct_from: 
- /docker-hub/company-scim/
---

{% include admin-scim.md %}

## Set up SCIM

You must make sure you have [configured SSO](../single-sign-on/index.md) before you enable SCIM. Enforcing SSO is not required.

### Step one: Enable SCIM in Docker Hub

1. Sign in to Docker Hub, navigate to the **Organizations** page and select your organization or company.
2. Select **Settings**. If you are setting up SCIM for an organization you then need to select **Security**. 
3. In the **Single Sign-On Connection** table, select the **Actions** icon and **Setup SCIM**.
4. Copy the **SCIM Base URL** and **API Token** and paste the values into your IdP.

### Step two: Enable SCIM in your IdP

Follow the instructions provided by your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm){: target="_blank" rel="noopener" class="_" }
- [Azure AD](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/user-provisioning){: target="_blank" rel="noopener" class="_" }
- [OneLogin](https://developers.onelogin.com/scim/create-app){: target="_blank" rel="noopener" class="_" }

## Disable SCIM

If SCIM is disabled, any user provisioned through SCIM will remain in the organization. Future changes for your users will not sync from your IdP. User de-provisioning is only possible when manually removing the user from the organization.

1. In the **Single Sign-On Connection** table, select the **Actions** icon.
2. Select **Disable SCIM**.