---
description: Use roles in Docker Hub to control who has access to content, registry, and organization management permissions.
keywords: members, teams, organizations, company, roles
title: Roles and permissions
---

Member roles let organization and company owners control who in an organization has access to products and features on Docker Hub. This section is for organization or company owners who want to learn about the defined roles and their permission scopes.

## Roles and access

When you invite users to your organization on Docker Hub, you assign a role. A role is a collection of permissions. Roles define access to perform actions like creating repositories, pulling images, publishing extensions, or managing repositories.

> Tip
>
> You don't need to create teams to assign a role to grant administrative access. You can still use teams to group members in your organization and assign [repository access permissions](/docker-hub/manage-a-team/#permissions-reference). Note that when you create new organizations, the owners team is no longer created automatically and only the owner role will grant a user owner permissions.
{: .tip }

The following roles are available to assign:

- **Member** - Non-administrative role. Members can view other members that are in the same organization.
- **Editor** - Partial administrative access to the organization. Editors can create, edit, and delete repositories. They can also edit an existing team's access permissions.
- **Organization owner** - Full organization administrative access. Organization owners can manage organization repositories, teams, members, settings, and billing.
- **Company owner** - In addition to the permissions of an organization owner, company owners can configure settings for their associated organizations.

## Permissions

The following lists the content and registry permissions for each role:

| Permission | Member | Editor | Organization owner | Company owner |
|:----------------------- |:------ |:-------|:------------------ |:----------- |
| Explore images and extensions | ✅ | ✅ | ✅ | ✅ |
| Star, favorite, vote, and comment on content | ✅ | ✅ | ✅ | ✅ |
| Pull images | ✅ | ✅ | ✅ | ✅ |
| Create and publish an extension | ✅ | ✅ | ✅ | ✅ |
| Become a Verified, Official, or Open Source publisher | ❌ | ❌ | ✅ | ✅ |
| Observe content engagement as a publisher | ❌ | ❌ | ✅ | ✅ |
| Create public and private repositories | ❌ | ✅ | ✅ | ✅ |
| Pull, push, edit, and delete a repository | ❌ | ✅ | ✅ | ✅ |
| Assign team permissions to repositories | ❌ | ✅ | ✅ | ✅ |
| Manage tags | ❌ | ✅ | ✅ | ✅ |
| View repository activity | ❌ | ❌ | ✅ | ✅ |
| Set up automated builds | ❌ | ✅ | ✅ | ✅ |
| Edit build settings | ❌ | ✅ | ✅ | ✅ |
| Set up vulnerability analysis with Docker Scout | ❌ | ✅ | ✅ | ✅ |
| View teams | ❌ | ✅ | ✅ | ✅ |
| Create teams | ❌ | ❌ | ✅ | ✅ |
| Manage, including delete, teams | ❌ | ❌ | ✅ | ✅ |

Note that editors and owners can give teams repository access permissions. See [Create and manage a team permissions reference](/docker-hub/manage-a-team/#permissions-reference).

The following lists the organization management permissions for each role:

| Permission | Member | Editor | Organization owner | Company owner |
|:----------------------- |:------ |:-------|:------------------ |:----------- |
| Configure the organization's settings (including linked services) | ❌ | ❌ | ✅ | ✅ |
| Add organizations to a company | ❌ | ❌ | ✅ | ✅ |
| Invite members | ❌ | ❌ | ✅ | ✅ |
| Manage members | ❌ | ❌ | ✅ | ✅ |
| Manage member permissions | ❌ | ❌ | ✅ | ✅ |
| View member activity | ❌ | ❌ | ✅ | ✅ |
| Export and reporting | ❌ | ❌ | ✅ | ✅ |
| Image access | ❌ | ❌ | ✅ | ✅ |
| Registry access | ❌ | ❌ | ✅ | ✅ |
| Setup Single Sign-On (SSO) and SCIM | ❌ | ❌ | ✅ * | ✅ |
| Require Desktop login | ❌ | ❌ | ✅ * | ✅ |
| Manage billing information (e.g. billing address) | ❌ | ❌ | ✅ | ✅ |
| Manage payment methods (e.g. credit card or invoice) | ❌ | ❌ | ✅ | ✅ |
| View billing history | ❌ | ❌ | ✅ | ✅ |
| Manage subscriptions | ❌ | ❌ | ✅ | ✅ |
| Manage seats | ❌ | ❌ | ✅ | ✅ |
| Upgrade and downgrade plans | ❌ | ❌ | ✅ | ✅ |

_* If not part of a company_
