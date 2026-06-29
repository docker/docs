---
title: Audit logging
linkTitle: Audit logs
weight: 28
description: Capture a structured, durable record of every sandbox policy decision for SIEM ingestion and compliance.
keywords: docker sandboxes, audit log, audit logging, policy decision, SIEM, compliance, jsonl, splunk, filebeat
---

The sandbox daemon records a structured audit event for every policy decision
it makes. Each record captures who triggered the evaluation, when it happened,
which rule matched, and whether the resource was allowed or denied. Records are
written to disk as JSON Lines (`.jsonl`) so existing SIEM and log-shipping
tools can collect them. The records stay on the machine that produced them.
Docker doesn't collect or ingest audit data.

> [!NOTE]
> Audit logging is part of Docker AI Governance and requires a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

Audit logging is active only while your organization enforces a centralized
governance policy. The subscription alone doesn't produce records. If your
organization hasn't configured and enforced an [organization policy](org.md),
the daemon writes no audit logs. To confirm governance is active, run `sbx
policy ls` — the output begins with a `Policy rules` header listing a
`Governance  Managed by <org>` line when an organization policy is in effect.

Audit logging complements [monitoring](monitoring.md). Monitoring with `sbx
policy ls` and `sbx policy log` is for live, interactive debugging. Audit
logging produces a durable trail for security review and compliance.

## What gets recorded

The daemon writes two categories of record:

- Evaluation records capture each policy decision: the resource, the
  action, the verdict, and the reason for a denial.
- Session lifecycle records mark the start and end of each daemon run.
  Evaluation records share the run's `audit_session_id`, so you can correlate
  every decision back to a single daemon session.

A network evaluation record looks like this:

```json
{
  "audit_event_id": "95e7257f-93c9-4f29-bde7-88830e2dae80",
  "timestamp": "2026-05-28T19:15:00.728933Z",
  "schema_version": "1.82.0",
  "category": "AUDIT_CATEGORY_EVALUATION",
  "decision": "AUDIT_DECISION_DENY",
  "username": "jordandoe",
  "user_email": "jordandoe@example.com",
  "org_id": "9f8e7d6c-5b4a-3210-fedc-ba9876543210",
  "org_name": "Acme Inc",
  "audit_session_id": "8a3bc076-79d0-4502-baf3-cc6ad35fb578",
  "resource_id": "example.com:443",
  "os": "macos",
  "app_version": "v0.31.0",
  "client_name": "sbx",
  "hostname": "host-machine",
  "deny_reason": [
    "no applicable policies for op(action=net:connect:tcp, resource=net:domain:example.com:443)"
  ],
  "action_type": "network_egress",
  "network_egress": { "protocol": "tcp" },
  "agent": "claude"
}
```

Common fields include:

| Field              | Description                                                                                                  |
| ------------------ | ------------------------------------------------------------------------------------------------------------ |
| `timestamp`        | UTC time of the decision.                                                                                    |
| `schema_version`   | Version of the record schema. Pin your SIEM field mappings to it, as the format is a stable contract.        |
| `category`         | `AUDIT_CATEGORY_EVALUATION` for policy decisions, `AUDIT_CATEGORY_MANAGEMENT` for session lifecycle records. |
| `audit_session_id` | Identifies the daemon run that produced the record.                                                          |
| `username`         | The signed-in Docker user's Docker Hub username.                                                             |
| `user_email`       | The signed-in Docker user's email address.                                                                   |
| `org_id`           | ID of the organization whose governance policy is in effect.                                                 |
| `org_name`         | Display name of the organization whose governance policy is in effect.                                       |
| `action_type`      | The kind of access evaluated, such as `network_egress`.                                                      |
| `resource_id`      | The target of the evaluation, such as a host and port.                                                       |
| `decision`         | `AUDIT_DECISION_ALLOW` or `AUDIT_DECISION_DENY`.                                                             |
| `deny_reason`      | Why a denied request was blocked. Present on deny decisions.                                                 |
| `agent`            | The AI agent driving the sandbox (for example, `claude`, `codex`). Omitted when the agent is unknown.       |

Each record is attributed to the signed-in Docker user and the organization
whose governance policy is in effect.

## Where records are stored

The daemon writes audit records, not the CLI. Running a command such as `sbx
create` sends a request to the daemon, and the daemon emits the resulting
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

The daemon writes in-progress records to a temporary `.tmp` file and seals it
into a final `.jsonl` file by atomic rename. Sealing happens at a rotation
threshold (by default 5 minutes, 1000 events, or 50 MiB, whichever comes
first) or when the daemon shuts down cleanly. Only sealed `.jsonl` files are
complete. Treat `.tmp` files as incomplete and don't collect them.

Sandboxes never delete sealed files. Retention and cleanup are the
responsibility of your log shipper or your own housekeeping.

## Collect records with a SIEM

Point your log shipper at the audit directory and configure it to collect
sealed `.jsonl` files only. Tools such as the Splunk Universal Forwarder,
Filebeat, and CrowdStrike Falcon LogScale read the directory and forward each
line as an event. Because in-progress records live in `.tmp` files until they
are sealed, collectors never see partial records.
