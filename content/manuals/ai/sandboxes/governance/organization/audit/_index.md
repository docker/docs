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

{{< summary-bar feature_name="AI Governance" >}}

Your organization's policies decide what sandboxes can do. AI Governance Audit
Logs show you what happened: each time a policy decides a request from a
governed developer's sandbox, a record captures who made the request, what it
was, and what the policy decided. You can review the records in Docker Home,
export them, forward them to your SIEM, or collect them from each machine.

This page explains what audit logs record, what you need to generate them,
where records go, and what you can do with them.

> [!NOTE]
> Audit logs cover local sandboxes, the sandboxes that run on developers'
> machines. Docker Cloud delivery stores their audit records in the cloud. It
> doesn't add cloud sandbox coverage. For cloud network decisions, see
> [Cloud policy logs](../../../cloud/network-policy.md#inspect-network-policy).

## What audit logs record

AI Governance Audit Logs record Docker AI Governance activity for your
organization. Each record captures the principal, action, target, decision, and
time for a governance event. Records contain metadata only. They don't contain
prompt content, agent output, or parameter values.

Audit logs cover Docker Sandboxes policy decisions and sandbox session events.
Other Docker AI sources can emit records through the same schema as they
become available. For every field a record can contain, see the
[Audit record reference](../../reference/audit-records.md).

## What you need

Audit logs are exposed when AI Governance is enabled for your organization.
To use them, your organization needs:

- A Docker [AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md)
- An enforced organization governance policy
- An organization owner, or a user with a
  [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
  that includes AI Governance audit permissions, to configure delivery and view
  hosted events

Records come only from developers whose activity your organization governs.
Docker Sandboxes send audit records only for signed-in users who have an AI
Governance license and are governed by an enforced centralized
[organization policy](../_index.md). Users without both don't generate audit
data, and they don't appear in audit events or SIEM forwarding output.

## Where records go

Docker supports two delivery modes for audit records:

- **Local disk**: the sandbox daemon writes JSON Lines (`.jsonl`)
  files on each host. Use this mode for host-local retention, air-gapped
  collection, or collection through your own log shipper.
- **Docker Cloud**: Docker stores audit records in Docker Cloud. Cloud
  delivery powers the hosted audit log view, CSV export, and SIEM streaming from
  app.docker.com.

To choose which delivery modes your organization uses, see
[Configure audit delivery](configure.md).

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

## What you can do with records

How you use audit records depends on where they go:

- With Docker Cloud delivery, [view and export events](view-export.md) in
  Docker Home, and [forward them to your SIEM](siem.md). The hosted view, CSV
  export, and SIEM forwarding all require Docker Cloud delivery. Local delivery
  alone doesn't power them.
- With local delivery, collect the `.jsonl` files from each machine with your
  own log shipper. For more details, see
  [Collect local audit files with a log shipper](siem.md#collect-local-audit-files-with-a-log-shipper).

To set up delivery, see [Configure audit delivery](configure.md).

## Data handling

When Docker Cloud delivery is enabled, Docker stores audit records in Docker
Cloud for 90 days. For legal and privacy terms that govern Docker services, see
Docker's
[Terms of Service](https://www.docker.com/legal/docker-terms-service/) and
[Privacy Policy](https://www.docker.com/legal/privacy/).
