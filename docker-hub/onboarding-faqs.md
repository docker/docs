---
description: Frequently asked questions
keywords: onboarding, docker, teams, orgs
toc_max: 2
---

### What is a Docker ID?

A Docker ID is a username to access Docker Hub repositories and hosted Docker
services. All you need is an email address to create a Docker ID. Your Docker ID must be between 4 and 30 characters long, and can only contain
numbers and lowercase letters. You cannot use any special characters or spaces.
For more information, see [Docker ID](../docker-id/index.md). If your admin enforces [Single sign-on (SSO)](../single-sign-on/index.md), a Docker ID is provisioned for new users.

Developers may have multiple Docker IDs in order to separate their Docker IDs  that are associated with an organization in Docker Business or Team, and their personal use Docker IDs.

### What if my Docker ID is taken?

All Docker IDs are first-come, first-served except for companies that have a US Trademark on a username. If you have a trademark for your namespace, [Docker Support](https://hub.docker.com/support/contact/){: target="_blank" rel="noopener"
class="_"} can retrieve the Docker ID for you.

### What’s an organization?

Docker users become members of an organization when they are assigned to at
least one team in the organization. When you first create an organization,
you’ll see that you have a team, the **Owners** (Admins) team, with a single
member. An organization owner is someone that is part of the owners team. They
can create new teams and add members to an existing team using their Docker ID
or email address and by selecting a team the user should be part of. An
organization owner can also add additional organization owners to help them
manage users, teams, and repositories in the organization. [Learn more](orgs.md).

### How many organizations can I create?

To begin, you should set up your first organization and contact the Customer Success team at customer-success@docker.com to add the additional organizations. Also, if you are enabling SSO, it is configured based on your domain, not your organization.

### What's an organization name or namespace?

The organization name, sometimes referred to as the organization namespace or the org ID, is the unique identifier of a Docker organization. The organization name cannot be the same as an existing Docker ID.

### What’s a team?

A **Team** is a group of Docker users that belong to an organization. An organization can have multiple teams. When you first create an organization, you’ll see that you have a team, the owners team, with a single member. An organization owner can then create new teams and add members to an existing team using Docker IDs or email address and by selecting a team the user should be part of. [Learn more](orgs.md#create-a-team).

### Who is an organization owner?

An organization owner is an administrator who is responsible to manage
repositories and add team members to the organization. They have full access to
private repositories, all teams, billing information, and organization settings.
An organization owner can also specify [permissions](orgs.md#configure-repository-permissions) for each team in the
organization. Only an organization owner can enable SSO for the organization.
When SSO is enabled for your organization, the organization owner can also
manage users.

Docker can auto-provision Docker IDs for new end-users or users who'd like to
have a separate Docker ID for company use through SSO enforcement.

The organization owner can also add additional owners to help them manage users, teams, and repositories in the organization.

### How do I add an organization owner?

An existing owner can add additional team members as organization owners. All
they need to do is select the organization from the
[Organizations](https://hub.docker.com/orgs){: target="_blank" rel="noopener"
class="_"} page in Docker Hub, add the Docker ID/Email of the user, and then
select the **Owners** team from the drop-down menu. [Learn more](orgs.md#the-owners-team).

### Do users first need to authenticate with Docker before an owner can add them to an organization?

No. Organization owners can invite users through email and also choose a team for them to join within the invite.

### Can I force my organization's members to authenticate before using Docker Desktop and are there any benefits?

Yes. You can [enforce sign in](../docker-hub/configure-sign-in.md) and there are benefits.

Some benefits of enforcing sign in are:
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

Currently, we don’t have a way to notify you. However, if the number of team
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
organization. [Learn more](orgs.md#configure-repository-permissions).

### Can I configure multiple SSO identity providers (IdPs) to authenticate users to a single org?

Docker SSO allows only one IdP configuration per organization. For more
information, see [Configure SSO](../single-sign-on/index.md) and [SSO FAQs](../single-sign-on/faqs.md).

### What is a service account?

A [service account](../docker-hub/service-accounts.md) is a Docker ID used for automated management of container images or containerized applications. Service accounts are typically used in automated workflows, and do not share Docker IDs with the members in the Team or Business plan. Common use cases for service accounts include mirroring content on Docker Hub, or tying in image pulls from your CI/CD process.

### Does my organization need to use Docker's registry?

A registry is a hosted service containing repositories of images which responds to the Registry API. Docker Hub is Docker's primary registry, but Docker can easily be used with other container image registries. The default registry can be accessed by browsing to [Docker Hub](https://hub.docker.com) or using the `docker search` command.

### What is included in my Docker Business or Team plan?

For a list of features available in each tier, see [Docker subscription overview](../subscription/index.md).