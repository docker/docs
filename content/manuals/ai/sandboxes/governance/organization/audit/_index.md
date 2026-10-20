---
title: AI Governance Audit Logs
linkTitle: Audit logs
weight: 40
description: Capture, view, export, and collect structured audit records for Docker AI Governance policy decisions.
keywords: docker sandboxes, audit log, audit logging, AI Governance, policy decision, SIEM, compliance, jsonl
aliases:
  - /ai/sandboxes/governance/audit/
  - /ai/sandboxes/governance/audit/local/
---

{{< summary-bar feature_name="AI Governance Audit Logs" >}}

The Docker Sandboxes coverage on this page applies to local sandboxes.
Docker Cloud delivery stores their audit records in the cloud; it does not add
cloud sandbox coverage. For cloud network decisions, see
[Cloud policy logs](../../../cloud/network-policy.md#inspect-network-policy).

AI Governance Audit Logs record Docker AI Governance activity for your
organization. Each record captures the principal, action, target, decision, and
time for a governance event. Records contain metadata only. They don't contain
prompt content, agent output, or parameter values.

Audit logs are exposed when AI Governance is enabled for your organization.
Docker Sandboxes send audit records only for signed-in users who have an AI
Governance license and are governed by an enforced centralized [organization
policy](../_index.md). Docker Sandboxes users without both
don't send audit data to audit logs.

> [!NOTE]
> AI Governance Audit Logs are part of Docker AI Governance and require a
> separate paid subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Requirements

To use AI Governance Audit Logs, your organization needs:

- A Docker [AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md)
- An enforced organization governance policy
- A Docker organization account
- An organization owner, or a user with a [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md) that includes AI Governance audit permissions, to configure delivery and view hosted events

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

Organization owners and users with a [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md) that includes AI Governance audit permissions can configure local disk, Docker Cloud, or both.

The hosted audit log view, CSV export, and SIEM forwarding all require Docker Cloud delivery to be enabled. Local delivery alone does not power these features.

Organizations that used local audit logging before hosted audit logs were
available start with cloud delivery off until an owner opts in from
[audit delivery settings](configure.md).

### Where records are stored

Local audit files stay on the host that produced them. Local delivery is
independent of [cloud delivery](configure.md), so you can use local files with
cloud delivery on or off. To confirm that a machine is governed and writes
records, run `sbx policy ls`. The output includes a `Governance: Managed by
<org>` line when an organization policy is in effect.

The daemon writes audit records, not the CLI. Running a command such as
`sbx create` sends a request to the daemon, and the daemon emits the resulting
record to its own audit directory.

The default location depends on your operating system:

| OS      | Default path                                                      |
| ------- | ----------------------------------------------------------------- |
| macOS   | `~/Library/Logs/com.docker.sandboxes/sandboxes/auditkit/`         |
| Linux   | `${XDG_STATE_HOME:-~/.local/state}/sandboxes/sandboxes/auditkit/` |
| Windows | `%LOCALAPPDATA%\DockerSandboxes\sandboxes\logs\auditkit\`         |

The directory layout differs by platform because each operating system places
application logs in its own conventional location.

Files are named `audit-<utc-timestamp>-<process-uuid>-<seq>.jsonl`.

The daemon writes in-progress records to a temporary `.tmp` file and
finalizes it into a `.jsonl` file by atomic rename. Finalization happens at a
rotation threshold: by default, 5 minutes, 1000 events, or 50 MiB, whichever
comes first. Finalization also happens when the daemon shuts down cleanly. Only
`.jsonl` files are complete. Treat `.tmp` files as incomplete and don't collect
them.

Sandboxes never delete `.jsonl` files. Retention and cleanup are the
responsibility of your log shipper or your own housekeeping.

## Data handling

When Docker Cloud delivery is enabled, Docker stores audit records in Docker
Cloud for the retention window configured by your organization. For legal and
privacy terms that govern Docker services, see Docker's [Terms of
Service](https://www.docker.com/legal/docker-terms-service/) and [Privacy
Policy](https://www.docker.com/legal/privacy/).

## Learn more

- [Configure audit delivery](configure.md)
- [View and export audit events](view-export.md)
- [SIEM forwarding](siem.md)
- [Retrieve audit events with the API](api.md)
- [Audit record reference](../../reference/audit-records.md)
