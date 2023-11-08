---
description: System for Cross-domain Identity Management
keywords: SCIM, SSO
title: SCIM
direct_from:
- /docker-hub/company-scim/
- /docker-hub/scim/
- /admin/company/settings/scim/
-/admin/organization/security-settings/scim/
---

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

You must make sure you have [configured SSO](single-sign-on/configure/_index.md) before you enable SCIM. Enforcing SSO is not required.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-scim %}}

{{< /tab >}}
{{< tab name="Docker Admin" >}}

{{< include "admin-early-access.md" >}}

{{% admin-scim product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

