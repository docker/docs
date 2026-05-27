---
title: Manage license assignment
linkTitle: License assignment
description: Manage product licenses for your organization, including assignment, revocation, and automatic assignment.
keywords: licenses, organization, members, Docker Offload, AI governance, license assignment, admin console
weight: 30
---

You can manage licenses on the Members page from Docker Home. The Members page gives admins a bird’s eye view of pending invites, roles, and membership. You can use this page to manually or automatically assign licenses, or revoke licenses on a per-member basis.

## Prerequisites

You can manage licenses for Docker Offload or AI Governance. Each product has its own requirements to get started:

- To purchase Docker Offload licenses, you must have a Docker Team or Docker Business subscription.
- AI Governance licenses are currently available by contacting sales.

You can purchase licenses for both products by [contacting sales](https://www.docker.com/pricing/contact-sales/).

## License Management

You manage licenses from the [Members page in the Admin Console](/manuals/admin/organization/manage/members.md). After purchasing a license-based product, each product appears on the Members page with information about the number of total licenses and licenses used by members.

- Inviting a member may consume a seat, but doesn’t auto-assign product licenses by default.
- Only admin owners can view details about organization licenses.

## Assign and revoke licenses

1. Sign in to [Docker Home](https://app.docker.com), then select your organization.
2. Go to **Members** to view all your organization members.
3. Assign or revoke licenses to members.

- You can search individual members from the members table or filtering for their username, email, or full name.
- You can bulk invite members by selecting multiple members.

## Automated license assignment

Automatic assignment gives members a product license from your pool the first time they use a supported product, instead of assigning licenses manually on **Members**.

### Prerequisites

Before the automatic assignment toggle can assign licenses, members must already be in your organization through **SSO and user provisioning** (Just-in-Time, SCIM, or your configured provisioning flow). Automatic assignment does not invite users or create accounts—it only assigns licenses to existing members.

- **AI Governance** — SSO and provisioning capabilities are included when you purchase AI Governance. Set up [SSO](/manuals/enterprise/security/single-sign-on/connect.md) and [provisioning](/manuals/enterprise/security/provisioning/_index.md) before relying on automatic assignment for `sbx`.
- **Docker Core** — On Team and Business subscriptions, SSO and provisioning are already part of your subscription. Complete the same setup before automatic Docker Core assignment applies on Desktop or CLI sign-in.

Automatic assignment applies to **AI Governance** and **Docker Core** only, not Docker Offload.

### Configure

On **Members**, each license pool tile has an automatic assignment toggle (**on** by default). When an org member without a license uses:

- **AI Governance** — Docker Sandboxes (`sbx`) (policy check)
- **Docker Core** — Docker Desktop or the CLI (entitlement check)

Docker assigns a seat from the pool if one is available. If the pool is exhausted, no license is issued and org owners see a notice on **Members**.

Turning the toggle off stops new automatic assignments; it does not revoke existing licenses.

## What’s next
