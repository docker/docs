---
description: Group mapping in Docker Hub
keywords: Group Mapping, SCIM, Docker Hub
title: Group Mapping
---

With directory group-to-team provisioning from your IdP, user updates will automatically sync with your Docker organizations and teams. 

To correctly assign your users to Docker teams, you must create groups in your IDP following the naming pattern `organization:team`. For example, if you want to manage provisioning for the team "developers” in Docker, and your organization name is “moby,” you must create a group in your IdP with the name “moby:developers”. 

Once you enable group mappings in your connection, users assigned to that group in your IdP will automatically be added to the team “developers” in Docker.

>**Tip**
>
>Use the same names for the Docker teams as your group names in the IdP to prevent further configuration. When you sync groups, a group is created if it doesn’t already exist.
{: .tip}

To take advantage of group mapping, make sure you have [enabled SCIM](scim.md) and then follow the instructions provided by your IdP:

- [Okta](https://help.okta.com/en-us/Content/Topics/users-groups-profiles/usgp-enable-group-push.htm){: target="_blank" rel="noopener" class="_" }
- [Azure AD](https://learn.microsoft.com/en-us/azure/active-directory/app-provisioning/customize-application-attributes){: target="_blank" rel="noopener" class="_" }
- [OneLogin](https://developers.onelogin.com/scim/create-app){: target="_blank" rel="noopener" class="_" }

Once complete, a user who signs in to Docker through SSO is automatically added to the organizations and teams mapped in the IdP.