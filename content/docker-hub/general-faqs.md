---
title: General FAQs for Docker Hub
description: Frequently asked administration questions
keywords: onboarding, docker, teams, orgs
redirect:
- /docker-hub/onboarding-faqs/
---

### What is a Docker ID?

A Docker ID is a username for your Docker account that lets you access Docker products. All you need is an email address to create a Docker ID, or you can sign up with your Google or GitHub account. Your Docker ID must be between 4 and 30 characters long, and can only contain
numbers and lowercase letters. You cannot use any special characters or spaces.

For more information, see [Docker ID](../docker-id/index.md). If your admin enforces [Single sign-on (SSO)](../security/for-admins/single-sign-on/index.md), a Docker ID is provisioned for new users.

Developers may have multiple Docker IDs in order to separate their Docker IDs that are associated with an organization in Docker Business or Team, and their personal use Docker IDs.

### What if my Docker ID is taken?

All Docker IDs are first-come, first-served except for companies that have a US Trademark on a username. If you have a trademark for your namespace, [Docker Support](https://hub.docker.com/support/contact/) can retrieve the Docker ID for you.

### What’s an organization?

Docker users become members of an organization when they're associated with the organization by an organization owner. An organization owner is someone assigned the owner role. They
can create new teams and add members to an existing team using their Docker ID
or email address and by selecting a team the user should be part of. An
organization owner can also add additional organization owners to help them
manage users, teams, and repositories in the organization. See [Create your organization](orgs.md).

### What's an organization name or namespace?

The organization name, sometimes referred to as the organization namespace or the org ID, is the unique identifier of a Docker organization. The organization name cannot be the same as an existing Docker ID.

### What are roles?

A role is a collection of permissions granted to members. Roles define access to perform actions in Docker Hub like creating repositories, managing tags, or viewing teams. See [Roles and permissions](roles-and-permissions.md).

### What’s a team?

A **Team** is a group of Docker users that belong to an organization. An organization can have multiple teams. An organization owner can then create new teams and add members to an existing team using Docker IDs or email address and by selecting a team the user should be part of. See [Create and manage a team](manage-a-team.md).

### What's a company?

A **Company** is a management layer that centralizes administration of multiple organizations. Administrators can add organizations with a Docker Business subscription to a company and configure settings for all organizations under the company. See [Set up your company](creating-companies.md).

### Who is an organization owner?

An organization owner is an administrator who is responsible to manage
repositories and add team members to the organization. They have full access to
private repositories, all teams, billing information, and organization settings.
An organization owner can also specify [permissions](manage-a-team.md#configure-repository-permissions-for-a-team) for each team in the
organization. Only an organization owner can enable SSO for the organization.
When SSO is enabled for your organization, the organization owner can also
manage users.

Docker can auto-provision Docker IDs for new end-users or users who'd like to
have a separate Docker ID for company use through SSO enforcement.

The organization owner can also add additional owners to help them manage users, teams, and repositories in the organization.

### Can I configure multiple SSO identity providers (IdPs) to authenticate users to a single org?

Docker SSO allows only one IdP configuration per organization. For more
information, see [Configure SSO](../security/for-admins/single-sign-on/configure/_index.md) and [SSO FAQs](../faq/security/single-sign-on/faqs.md).

### What is a service account?

A [service account](../docker-hub/service-accounts.md) is a Docker ID used for automated management of container images or containerized applications. Service accounts are typically used in automated workflows, and do not share Docker IDs with the members in the Team or Business plan. Common use cases for service accounts include mirroring content on Docker Hub, or tying in image pulls from your CI/CD process.

### Can I delete or deactivate a Docker account for another user?

Only someone with access to the Docker account can deactivate the account. For more details, see [Deactivating an account](../docker-hub/deactivate-account.md/).

If the user is a member of your organization, you can remove the user from your organization. For more details, see [Remove a member or invitee](/docker-hub/members/#remove-a-member-or-invitee).