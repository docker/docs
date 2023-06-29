---
description: System for Cross-domain Identity Management
keywords: SCIM, SSO
title: SCIM
---

{% include admin-early-access.md %}

{% include admin-scim.md %}

## Set up SCIM

You must make sure you have [configured SSO](sso.md) before you enable SCIM. Enforcing SSO is not required.

### Step one: Enable SCIM in Docker Admin

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Security**.
4. In the **Single Sign-On Connection** table, select the **Actions** icon and **Setup SCIM**.
5. Copy the **SCIM Base URL** and **API Token** and paste the values into your IdP.

### Step two: Enable SCIM in your IdP

Follow the instructions provided by your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/Apps/Apps_App_Integration_Wizard_SCIM.htm){: target="_blank" rel="noopener" class="_" }
- [Azure AD](https://learn.microsoft.com/en-us/azure/databricks/administration-guide/users-groups/scim/aad#step-2-configure-the-enterprise-application){: target="_blank" rel="noopener" class="_" }
- [OneLogin](https://developers.onelogin.com/scim/create-app){: target="_blank" rel="noopener" class="_" }

## Disable SCIM

If SCIM is disabled, any user provisioned through SCIM will remain in the organization. Future changes for your users will not sync from your IdP. User de-provisioning is only possible when manually removing the user from the organization.

1. In the **Single Sign-On Connection** table, select the **Actions** icon.
2. Select **Disable SCIM**.