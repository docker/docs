---
title: Manage domains
description: Add, verify, and manage domains to control user access and enable auto-provisioning in Docker organizations
keywords: domain management, domain verification, auto-provisioning, user management, DNS, TXT record, Admin Console
weight: 55
aliases:
 - /security/for-admins/domain-management/
---

{{< summary-bar feature_name="Domain management" >}}

Domain management lets you add and verify domains for your organization, then enable auto-provisioning to automatically add users when they sign in with email addresses that match your verified domains.

This approach simplifies user management, ensures consistent security settings, and reduces the risk of unmanaged users accessing Docker without visibility or control.

## Add and verify a domain

Adding a domain requires verification to confirm ownership. The verification process uses DNS records to prove you control the domain.

### Add a domain

1. Sign in to [Docker Home](https://app.docker.com) and select
your organization. If your organization is part of a company, select the company
and configure the domain for the organization at the company level.
1. Select **Admin Console**, then **Domain management**.
1. Select **Add a domain**.
1. Enter your domain and select **Add domain**.
1. In the pop-up modal, copy the **TXT Record Value** to verify your domain.

### Verify a domain

Verification confirms that you own the domain by adding a TXT record to your Domain Name System (DNS) host.

It can take up to 72 hours for the DNS change to propagate. Docker automatically
checks for the record and confirms ownership once the change is recognized.

> [!TIP]
>
> The record name field determines where the TXT record is added in your domain
(root or subdomain). In general, refer to the following tips for
adding a record name:
>
> - Use `@` or leave the record name empty for root domains like `example.com`,
depending on your provider.
> - Don't enter values like `docker`, `docker-verification`, `www`, or your
domain name. These values may direct to the wrong place.
>
> Check your DNS provider's documentation to verify record name requirements.

Follow the steps for your DNS provider to add the **TXT Record Value**. If
your provider isn't listed, use the steps for "Other providers":

{{< tabs >}}
{{< tab name="AWS Route 53" >}}

1. Add your TXT record to AWS by following [Creating records by using the Amazon Route 53 console](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-creating.html).
1. Wait up to 72 hours for TXT record verification.
1. Return to the **Domain management** page of the
[Admin Console](https://app.docker.com/admin) and select **Verify** next to
your domain name.

{{< /tab >}}
{{< tab name="Google Cloud DNS" >}}

1. Add your TXT record to Google Cloud DNS by following [Verifying your domain with a TXT record](https://cloud.google.com/identity/docs/verify-domain-txt).
1. Wait up to 72 hours for TXT record verification.
1. Return to the **Domain management** page of the
[Admin Console](https://app.docker.com/admin) and select **Verify** next to
your domain name.

{{< /tab >}}
{{< tab name="GoDaddy" >}}

1. Add your TXT record to GoDaddy by following [Add a TXT record](https://www.godaddy.com/help/add-a-txt-record-19232).
1. Wait up to 72 hours for TXT record verification.
1. Return to the **Domain management** page of the
[Admin Console](https://app.docker.com/admin) and select **Verify** next to
your domain name.

{{< /tab >}}
{{< tab name="Other providers" >}}

1. Sign in to your domain host.
1. Add a TXT record to your DNS settings using the **TXT Record Value** from Docker.
1. Wait up to 72 hours for TXT record verification.
1. Return to the **Domain management** page of the
[Admin Console](https://app.docker.com/admin) and select **Verify** next to
your domain name.

{{< /tab >}}
{{< /tabs >}}

## Configure auto-provisioning

Auto-provisioning automatically adds users to your organization when they sign in with email addresses that match your verified domains. You must verify a domain before enabling auto-provisioning.

> [!IMPORTANT]
>
> For domains that are part of an SSO connection, Just-in-Time (JIT) provisioning takes precedence over auto-provisioning when adding users to an organization.

### How auto-provisioning works

When auto-provisioning is enabled for a verified domain:

- Users who sign in to Docker with matching email addresses are automatically added to your organization.
- Auto-provisioning only adds existing Docker users to your organization, it doesn't create new accounts.
- Users experience no changes to their sign-in process.
- Company and organization owners receive email notifications when new users are added.
- You may need to [manage seats](/manuals/subscription/manage-seats.md) to accomodate new users.

### Enable auto-provisioning

Auto-provisioning is configured per domain. To enable it:

1. Sign in to [Docker Home](https://app.docker.com) and select
your organization. If your organization is part of a company, select the company
and configure the domain for the organization at the company level.
1. Select **Admin Console**, then **Domain management**.
1. Select the **Actions menu** next to the domain you want to enable
auto-provisioning for.
1. Select **Enable auto-provisioning**.
1. Optional. If enabling auto-provisioning at the company level, select an
organization.
1. Select **Enable** to confirm.

The **Auto-provisioning** column will update to **Enabled** for the domain.

### Disable auto-provisioning

To disable auto-provisioning for a user:

1. Sign in to [Docker Home](https://app.docker.com) and select
your organization. If your organization is part of a company, select the company
and configure the domain for the organization at the company level.
1. Select **Admin Console**, then **Domain management**.
1. Select the **Actions menu** next to your domain.
1. Select **Disable auto-provisioning**.
1. Select **Disable** to confirm.

## Audit domains for uncaptured users

{{< summary-bar feature_name="Domain audit" >}}

Domain audit identifies uncaptured users. Uncaptured users are Docker users who have authenticated using an email address associated with your verified domains but aren't members of your Docker organization.

### Limitations

Domain audit can't identify:

- Users who access Docker Desktop without authenticating
- Users who authenticate using an account that doesn't have an
email address associated with one of your verified domains

To prevent unidentifiable users from accessing Docker Desktop, [enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).

### Run a domain audit

1. Sign in to [Docker Home](https://app.docker.com) and choose your
company.
1. Select **Admin Console**, then **Domain management**.
1. In **Domain audit**, select **Export Users** to export a CSV file
of uncaptured users.

The CSV file contains the following columns:

    - Name: Docker user's display name
    - Username: Docker ID of the user
    - Email: Email address of the user

### Invite uncaptured users

You can bulk invite uncaptured users to your organization using the exported
CSV file. For more information on bulk inviting users, see
[Manage organization members](/manuals/admin/organization/members.md).

## Delete a domain

Deleting a domain removes its TXT record value and disables any associated auto-provisioning.

>[!WARNING]
>
> Deleting a domain will disable auto-provisioning for that domain and remove verification. This action cannot be undone.

To delete a domain:

1. Sign in to [Docker Home](https://app.docker.com) and select
your organization. If your organization is part of a company, select the company
and configure the domain for the organization at the company level.
1. Select **Admin Console**, then **Domain management**.
1. For the domain you want to delete, section the **Actions** menu, then
**Delete domain**.
1. To confirm, select **Delete domain** in the pop-up modal.
