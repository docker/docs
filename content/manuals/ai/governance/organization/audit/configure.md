---
title: Configure audit delivery
linkTitle: Configure delivery
weight: 10
description: Configure local and cloud delivery for Docker AI Governance Audit Logs, and review who changed the delivery settings.
keywords: docker sandboxes, audit delivery, AI Governance, audit logs, retention, cloud delivery, AI Platform
aliases:
  - /ai/sandboxes/governance/audit/configure/
---

{{< summary-bar feature_name="AI Governance" >}}

Configure where Docker delivers your organization's audit records, and review
who changed those settings. You can deliver records to local disk on each host,
to Docker Cloud, or to both. For how the two delivery modes differ, see
[Delivery modes](_index.md#delivery-modes).

## Before you begin

Your organization needs:

- A Docker [AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md)
- An enforced organization governance policy
- Organization owner access, or a
  [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
  with AI Governance audit permissions

Only users who have an AI Governance license and are governed by the enforced
organization policy send Docker Sandboxes audit data.

## Configure delivery

To configure audit delivery:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **Audit delivery**.
1. Choose one or both delivery modes:
   - **Local disk** writes audit records to JSON Lines files on each host.
   - **Docker Cloud** stores audit records in Docker Cloud for hosted search,
     CSV export, and SIEM forwarding.
1. Save your changes.

Cloud delivery is on by default when AI Governance is enabled, except for
organizations that used local audit logging before hosted audit logs were
available. Those organizations start with cloud delivery off until an owner
turns it on. To keep records local to your hosts, turn off **Docker Cloud** and
keep **Local disk** on.

Docker Cloud keeps audit records searchable in the hosted audit log view for 90
days, and in archive storage for 90 days.

## Review delivery change history

Docker keeps a record of every change made to your organization's audit delivery
settings. Each entry shows who made the change, when it happened, and the before
and after values. Use this history to audit configuration changes and verify
when delivery modes were modified.

To review the history:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **History**.
