---
description: Company FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on
title: FAQs on organizations
---

### What if the Docker ID I want for my organization or company is taken?

This depends on the state of the namespace, if trademark claims exist for the organization or company Docker ID, a manual flow for legal review is required.

### What if I want to create more than 3 organizations?

You can create multiple organizations or multiple teams under a single company. SSO is available at the company level.

### How do I add an organization owner?

An existing owner can add additional team members as organization owners. All
they need to do is select the organization from the
[Organizations](https://hub.docker.com/orgs) page in Docker Hub, add the Docker ID/Email of the user, and then
select the owner role from the drop-down menu. See [Organization owner](manage-a-team.md#organization-owner).

### How do I know how many active users are part of my organization?

If your organization uses a Software Asset Management tool, you can use it to find out how many users have Docker Desktop installed. If your organization doesn't use this software, you can run an internal survey to find out who is using Docker Desktop. See [Identify your Docker users and their Docker accounts](../docker-hub/onboard.md#step-1-identify-your-docker-users-and-their-docker-accounts). With a Docker Business subscription, you can manage members in your identity provider and automatically provision them to your Docker organization with [SSO](../single-sign-on/_index.md) or [SCIM](../docker-hub/scim.md).

### Do users first need to authenticate with Docker before an owner can add them to an organization?

No. Organization owners can invite users through email and also choose a team for them to join within the invite.

### Can I force my organization's members to authenticate before using Docker Desktop and are there any benefits?

Yes. You can [enforce sign-in](../docker-hub/configure-sign-in.md) and some benefits are:

- Administrators can enforce features like [Image Access Management](../docker-hub/image-access-management.md) and [Registry Access Management](../docker-hub/registry-access-management.md).
 - Administrators can ensure compliance by blocking Docker Desktop usage for users who do not sign in as members of the organization.

### If a user has their personal email associated with a user account in Docker Hub, do they have to convert to using the org’s domain before they can be invited to join an organization?

Yes. When SSO is enabled for your organization, each user must sign in with the company’s domain. However, the user can retain their personal credentials and create a new Docker ID associated with their organization's domain.

### Can I convert my personal user account (Docker ID) to an organization account?

Yes. You can convert your user account to an organization account. Once you
convert a user account into an organization, it is not possible to
revert it to a personal user account. For prerequisites and instructions, see
[Convert an account into an organization](convert-account.md).

### Our users create Docker Hub accounts through self-service. How do we know when the total number of users for the requested licenses has been met? Is it possible to add more members to the organization than the total number of licenses?

There isn't any automatic notification when the total number of users for the requested licenses has been met. However, if the number of team
members exceed the number of licenses, you will receive an error informing you
to contact the administrator due to lack of seats.

### How can I merge organizations in Docker Hub?

Reach out to your Support contact if you need to consolidate organizations.

### Do organization invitees take up seats?

Yes. A user invited to an organization will take up one of the provisioned
seats, even if that user hasn’t accepted their invitation yet. Organization
owners can manage the list of invitees through the **Invitees** tab on the organization settings page in Docker Hub.

### Do organization owners take a seat?

Yes. Organization owners will take up a seat.

### What is the difference between user, invitee, seat, and member?

User may refer to a Docker user with a Docker ID.

An invitee is a user who has been invited to join an organization, but has not yet accepted their invitation.

Seats is the number of planned members within an organization.

Member may refer to a user that has received and accepted an invitation to join an organization. Member can also refer to a member of a team within an organization.


### If there are two organizations and a user belongs to both orgs, do they take up two seats?

Yes. In a scenario where a user belongs to two orgs, they take up one seat in each organization.

### Is it possible to set permissions for repositories within an organization?

Yes. You can configure repository access on a per-team basis. For example, you
can specify that all teams within an organization have **Read and Write** access
to repositories A and B, whereas only specific teams have **Admin** access. Org
owners have full administrative access to all repositories within the
organization. See [Configure repository permissions for a team](manage-a-team.md#configure-repository-permissions-for-a-team).

### Does my organization need to use Docker's registry?

A registry is a hosted service containing repositories of images that responds to the Registry API. Docker Hub is Docker's primary registry, but you can use Docker with other container image registries. You can access the default registry by browsing to [Docker Hub](https://hub.docker.com) or using the `docker search` command.