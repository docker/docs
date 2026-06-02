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
  "schema_version": "1.80.0",
  "category": "AUDIT_CATEGORY_EVALUATION",
  "audit_session_id": "8a3bc076-79d0-4502-baf3-cc6ad35fb578",
  "username": "5f3e3556-ed49-4431-bdd8-24958cdc4340",
  "user_email": "jordan@example.com",
  "action_type": "network_egress",
  "resource_id": "example.com",
  "decision": "AUDIT_DECISION_DENY",
  "deny_reason": [
    "no applicable policies for op(action=net:connect:tcp, resource=net:domain:example.com)"
  ],
  "network_egress": { "protocol": "tcp" },
  "os": "macos",
  "hostname": "host-machine",
  "client_name": "sbx",
  "app_version": "v0.30.0"
}
```

Common fields include:

| Field              | Description                                                                                                  |
| ------------------ | ------------------------------------------------------------------------------------------------------------ |
| `timestamp`        | UTC time of the decision.                                                                                    |
| `schema_version`   | Version of the record schema. Pin your SIEM field mappings to it, as the format is a stable contract.        |
| `category`         | `AUDIT_CATEGORY_EVALUATION` for policy decisions, `AUDIT_CATEGORY_MANAGEMENT` for session lifecycle records. |
| `audit_session_id` | Identifies the daemon run that produced the record.                                                          |
| `username`         | The signed-in Docker user's account UUID.                                                                    |
| `user_email`       | The signed-in Docker user's email address.                                                                   |
| `action_type`      | The kind of access evaluated, such as `network_egress`.                                                      |
| `resource_id`      | The target of the evaluation, such as a hostname.                                                            |
| `decision`         | `AUDIT_DECISION_ALLOW` or `AUDIT_DECISION_DENY`.                                                             |
| `deny_reason`      | Why a denied request was blocked. Present on deny decisions.                                                 |

Identity is resolved from the signed-in Docker user at daemon startup. If no
user is signed in when the daemon starts, records still ship, but with empty
`username` and `user_email` fields.

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

## Override the storage location

Two environment variables change where records are written. Set them on the
daemon process, not the CLI.

| Variable                             | Effect                                                                                                                                          |
| ------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `SANDBOXES_STORAGE_ROOT=<dir>`       | Override the base storage directory. Audit records move under `<dir>/logs/`, keeping the same platform-specific namespace as the default paths. |
| `DOCKER_SANDBOXES_APP_NAME=<suffix>` | Append a suffix to the app name (`sandboxes` becomes `sandboxes-<suffix>`). Useful for running multiple daemon instances side by side.          |

The CLI starts the daemon automatically when none is running, so exporting
either variable in your shell propagates to the daemon it spawns. If a daemon
is already running, stop it first with `sbx daemon stop` so the next `sbx`
command starts a daemon that picks up the new value.
