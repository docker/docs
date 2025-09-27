---
description: Organization FAQs
linkTitle: Organization
weight: 20
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, organizations, administration, Admin Console, members, organization management, manage orgs
title: FAQs on organizations
tags: [FAQ]
aliases:
- /docker-hub/organization-faqs/
- /faq/admin/organization-faqs/
---

### How can I see how many active users are in my organization?

If your organization uses a Software Asset Management tool, you can use it to
find out how many users have Docker Desktop installed. If your organization
doesn't use this software, you can run an internal survey
to find out who is using Docker Desktop.

For more information, see [Identify your Docker users and their Docker accounts](../../admin/organization/onboard.md#step-1-identify-your-docker-users-and-their-docker-accounts).

### Do users need to authenticate with Docker before an owner can add them to an organization?

No. Organization owners can invite users with their email addresses, and also
assign them to a team during the invite process.

### Can I force my organization's members to authenticate before using Docker Desktop and are there any benefits?

Yes. You can
[enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).

Some benefits of enforcing sign-in are:

- Administrators can enforce features like [Image Access Management](/manuals/enterprise/security/hardened-desktop/image-access-management.md) and [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md).
 - Administrators can ensure compliance by blocking Docker Desktop usage for
 users who don't sign in as members of the organization.

### Can I convert my personal Docker ID to an organization account?

Yes. You can convert your user account to an organization account. Once you
convert a user account into an organization, it's not possible to
revert it to a personal user account.

For prerequisites and instructions, see
[Convert an account into an organization](convert-account.md).

### Do organization invitees take up seats?

Yes. A user invited to an organization will take up one of the provisioned
seats, even if that user hasn’t accepted their invitation yet.

To manage invites, see [Manage organization members](/manuals/admin/organization/members.md).

### Do organization owners take a seat?

Yes. Organization owners occupy a seat.

### What is the difference between user, invitee, seat, and member?

- User: Docker user with a Docker ID.
- Invitee: A user that an administrator has invited to join an organization but
has not yet accepted their invitation.
- Seats: The number of purchased seats in an organization.
- Member: A user who has received and accepted an invitation to join an
organization. Member can also refer to a member of a team within an
organization.

### If I have two organizations and a user belongs to both organizations, do they take up two seats?

Yes. In a scenario where a user belongs to two organizations, they take up one
seat in each organization.
