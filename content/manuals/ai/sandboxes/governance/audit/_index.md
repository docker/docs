---
title: AI Governance Audit Logs
linkTitle: Audit logs
weight: 28
description: Capture, view, export, and collect structured audit records for Docker AI Governance policy decisions.
keywords: docker sandboxes, audit log, audit logging, AI Governance, policy decision, SIEM, compliance, jsonl
---

{{< summary-bar feature_name="AI Governance Audit Logs" >}}

AI Governance Audit Logs record Docker AI Governance activity for your
organization. Each record captures the principal, action, target, decision, and
time for a governance event. Records contain metadata only. They don't contain
prompt content, agent output, or parameter values.

Audit logs are exposed when AI Governance is enabled for your organization.
Docker Sandboxes send audit records only for signed-in users who have an AI
Governance license and are governed by an enforced centralized [organization
policy](../access-controls/organization.md). Docker Sandboxes users without both
don't send audit data to audit logs.

> [!NOTE]
> AI Governance Audit Logs are part of Docker AI Governance and require a
> separate paid subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Requirements

To use AI Governance Audit Logs, your organization needs:

- A Docker [AI Governance plan](/manuals/subscription/plans/ai-governance.md)
- An enforced organization governance policy
- A Docker organization account
- An organization owner, or a user with a [custom role](/manuals/enterprise/security/roles-and-permissions/custom-roles/_index.md) that includes AI Governance audit permissions, to configure delivery and view hosted events

> [!NOTE]
> Other Docker subscriptions are not sufficient on their own to use AI Governance
> Audit Logs. Users without an AI Governance license and an enforced organization
> policy will not generate audit data and will not appear in audit events or SIEM
> forwarding output. Personal accounts are not supported.

## Coverage

AI Governance Audit Logs cover Docker Sandboxes policy decisions and sandbox
session events. Other Docker AI sources can emit records through the same schema
as they become available.

## Delivery modes

Docker supports two delivery modes for audit records:

- **Local disk**: the sandbox daemon writes JSON Lines (`.jsonl`)
  files on each host. Use this mode for host-local retention, air-gapped
  collection, or collection through your own log shipper.
- **Docker Cloud**: Docker stores audit records in Docker Cloud. Cloud
  delivery powers the hosted audit log view, CSV export, and SIEM streaming from
  app.docker.com. Cloud delivery is on by default when AI Governance is enabled.
  Organization owners can disable it in [audit delivery settings](configure.md).

Organization owners and users with a [custom role](/manuals/enterprise/security/roles-and-permissions/custom-roles/_index.md) that includes AI Governance audit permissions can configure local disk, Docker Cloud, or both.

The hosted audit log view, CSV export, and SIEM forwarding all require Docker Cloud delivery to be enabled. Local delivery alone does not power these features.

Organizations that used local audit logging before hosted audit logs were
available start with cloud delivery off until an owner opts in from
[audit delivery settings](configure.md).

## Data handling

When Docker Cloud delivery is enabled, Docker stores audit records in Docker
Cloud for the retention window configured by your organization. For legal and
privacy terms that govern Docker services, see Docker's [Terms of
Service](https://www.docker.com/legal/docker-terms-service/) and [Privacy
Policy](https://www.docker.com/legal/privacy/).

## Learn more

- [Local audit logs](local.md)
- [Configure audit delivery](configure.md)
- [View and export audit events](view-export.md)
- [SIEM forwarding](siem.md)
- [Audit record reference](record-reference.md)
