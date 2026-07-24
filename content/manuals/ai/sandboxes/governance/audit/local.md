---
title: Local audit logs
linkTitle: Local audit logs
weight: 10
description: Collect Docker AI Governance audit records from local JSON Lines files written by the sandbox daemon.
keywords: docker sandboxes, audit log, local audit logs, audit logging, jsonl, SIEM, splunk, filebeat
---

The sandbox daemon writes local audit records as JSON Lines (`.jsonl`) files.
Local audit logs stay on the host that produced them and can be collected by
your own log shipper. Local delivery is independent of
[cloud delivery](configure.md), so you can use local files with cloud delivery
on or off.

The sandbox daemon writes local audit records only for signed-in users who have
an AI Governance license and are governed by an enforced centralized
[organization policy](../org.md). Docker Sandboxes users without both do not
send audit data to audit logs. To confirm governance is active, run
`sbx policy ls`. The output includes a `Governance: Managed by <org>` line
when an organization policy is in effect.

## What gets recorded

The daemon writes two categories of local record:

- Evaluation records capture each policy decision: the resource, the action,
  the verdict, and the reason for a denial.
- Session lifecycle records mark the start and end of each daemon run.
  Evaluation records share the run's `audit_session_id`, so you can correlate
  every decision back to a daemon session.

Records contain metadata only. They don't contain prompt content, agent output,
or parameter values. For field details, see
[Audit record reference](record-reference.md).

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

## Where records are stored

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

## Collect records with a SIEM

Point your log shipper at the audit directory and configure it to collect
`.jsonl` files only. Tools such as the Splunk Universal Forwarder,
Filebeat, and Crowdstrike Falcon LogScale read the directory and forward each
line as an event. Because in-progress records live in `.tmp` files until they
are finalized, collectors never see partial records.
