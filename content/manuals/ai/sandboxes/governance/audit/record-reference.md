---
title: Audit record reference
linkTitle: Record reference
weight: 50
description: Reference fields, categories, decisions, and payload types for Docker AI Governance audit records.
keywords: docker sandboxes, audit record, audit schema, AI Governance, policy decision, action_type, jsonl
---

Docker AI Governance audit records use one schema across delivery modes. Local
JSON Lines files and cloud-delivered records contain the same metadata fields.

Records capture metadata only. They don't contain prompt content, agent output,
or parameter values. Parameter keys may appear when they help identify an
action.

## Common fields

| Field              | Description                                                                                                      |
| ------------------ | ---------------------------------------------------------------------------------------------------------------- |
| `audit_event_id`   | Unique ID for the audit event.                                                                                   |
| `timestamp`        | UTC time when Docker recorded the event.                                                                         |
| `schema_version`   | Version of the record schema. Pin SIEM field mappings to this value.                                             |
| `category`         | Event category, such as `AUDIT_CATEGORY_MANAGEMENT`, `AUDIT_CATEGORY_EVALUATION`, or `AUDIT_CATEGORY_EXECUTION`. |
| `decision`         | Governance decision for evaluation records.                                                                      |
| `username`         | The signed-in Docker user's Docker Hub username.                                                                 |
| `user_email`       | The signed-in Docker user's email address.                                                                       |
| `org_id`           | ID of the organization whose governance policy is in effect.                                                     |
| `org_name`         | Display name of the organization whose governance policy is in effect.                                           |
| `audit_session_id` | Identifies the daemon session that produced the record.                                                          |
| `resource_id`      | Target of the evaluation, such as a host and port, file path, or tool.                                           |
| `os`               | Operating system that produced the record.                                                                       |
| `app_version`      | Version of the Docker component that produced the record.                                                        |
| `client_name`      | Source component, such as `sbx` for Docker Sandboxes.                                                            |
| `hostname`         | Hostname of the machine that produced the record.                                                                |
| `deny_reason`      | Why a denied request was blocked. Present on deny decisions.                                                     |
| `enforcement_mode` | Governance mode in effect when Docker evaluated the request. See [Enforcement modes](#enforcement-modes).        |
| `approval`         | Approval details. Present when the decision involves an approval. See [Approval fields](#approval-fields).       |
| `action_type`      | Payload discriminator that identifies the action-specific object in the record.                                  |
| `agent`            | AI agent associated with the event, when Docker knows it.                                                        |

## Categories

| Category                    | Description                                                   |
| --------------------------- | ------------------------------------------------------------- |
| `AUDIT_CATEGORY_MANAGEMENT` | Session lifecycle, policy sync, and configuration events.     |
| `AUDIT_CATEGORY_EVALUATION` | Governance policy decisions, such as allow, deny, or consent. |
| `AUDIT_CATEGORY_EXECUTION`  | Outcomes after an evaluated action runs.                      |

## Decisions

| Decision                           | Description                                                |
| ---------------------------------- | ---------------------------------------------------------- |
| `AUDIT_DECISION_ALLOW`             | Docker allowed the action.                                 |
| `AUDIT_DECISION_DENY`              | Docker denied the action.                                  |
| `AUDIT_DECISION_APPROVAL_REQUIRED` | Docker held the action until a user approves or denies it. |
| `AUDIT_DECISION_APPROVAL_ALLOW`    | A user approved a pending request.                         |
| `AUDIT_DECISION_APPROVAL_DENY`     | A user denied a pending request.                           |

## Enforcement modes

The `enforcement_mode` field shows how governance applied to the request when
Docker evaluated it.

| Mode      | Description                                                                     |
| --------- | ------------------------------------------------------------------------------- |
| `enforce` | Your organization's governance policy is enforced.                              |
| `audit`   | Docker records the decision for observability but doesn't enforce policy.       |
| `off`     | Governance isn't active. Docker still writes the record to the local audit log. |

Approval resolution records don't include `enforcement_mode`.

## Approval fields

When a policy requires approval, Docker writes an
`AUDIT_DECISION_APPROVAL_REQUIRED` record while the request waits. When a user
approves or denies the request, Docker writes a second record with
`AUDIT_DECISION_APPROVAL_ALLOW` or `AUDIT_DECISION_APPROVAL_DENY`. Both records
carry the same `approval.approval_request_id`, so you can join them to
reconstruct the full approval.

| Field                          | Description                                                                       |
| ------------------------------ | --------------------------------------------------------------------------------- |
| `approval.approval_request_id` | ID that links an approval required record to the record that resolves it.         |
| `approval.grant_scope`         | Scope the user accepted. Present only on `AUDIT_DECISION_APPROVAL_ALLOW` records. |

Docker Sandboxes records an approved request's `approval.grant_scope` as
`AUDIT_APPROVAL_GRANT_SCOPE_PERSISTENT`, which means that the approval persists
beyond the request that was approved.

If a user approves a request but Docker can't apply the grant, the record has
an `AUDIT_DECISION_DENY` decision and includes the `approval` object.

## Action types

The `action_type` field identifies the action-specific payload in the record.

| Action type                  | Description                                                                 |
| ---------------------------- | --------------------------------------------------------------------------- |
| `session`                    | Sandbox daemon session lifecycle event.                                     |
| `network_egress`             | Network access evaluation.                                                  |
| `http_request`               | HTTP request evaluation. See [HTTP request payload](#http-request-payload). |
| `filesystem_mount`           | Filesystem mount or path access evaluation.                                 |
| `tool_invocation`            | Tool invocation evaluation.                                                 |
| `resource_read`              | Resource read evaluation.                                                   |
| `server_registration`        | Server registration event.                                                  |
| `prompt`                     | Prompt-related metadata event.                                              |
| `network_execution`          | Outcome of a network action.                                                |
| `filesystem_execution`       | Outcome of a filesystem action.                                             |
| `tool_execution`             | Outcome of a tool invocation.                                               |
| `resource_execution`         | Outcome of a resource read.                                                 |
| `policy_sync`                | Policy synchronization event.                                               |
| `pii_detection`              | Metadata event for a data detection result.                                 |
| `c_score_report`             | Metadata event for a C-score report.                                        |
| `policy_action`              | Policy configuration or policy action event.                                |

## HTTP request payload

When a network policy requires each HTTP request to be authorized, Docker
evaluates the request method and path in addition to the network connection.
These evaluations produce records with an `action_type` of `http_request`,
separate from the `network_egress` record for the connection.

| Field                 | Description                                                                 |
| --------------------- | --------------------------------------------------------------------------- |
| `http_request.method` | HTTP method, such as `HTTP_METHOD_GET` or `HTTP_METHOD_POST`.               |
| `http_request.host`   | Destination hostname or IP address. IPv6 addresses appear without brackets. |
| `http_request.port`   | Destination port.                                                           |
| `http_request.path`   | URL path, without the query string or fragment.                             |

## Sample record

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

The following sample shows a user approving an HTTP request that required
approval:

```json
{
  "audit_event_id": "3c1f6a2e-8b47-4d0a-9e55-2f7b1d6c9a04",
  "timestamp": "2026-09-25T17:42:08.114502Z",
  "schema_version": "1.216.0",
  "category": "AUDIT_CATEGORY_EVALUATION",
  "decision": "AUDIT_DECISION_APPROVAL_ALLOW",
  "username": "jordandoe",
  "user_email": "jordandoe@example.com",
  "org_id": "9f8e7d6c-5b4a-3210-fedc-ba9876543210",
  "org_name": "Acme Inc",
  "audit_session_id": "8a3bc076-79d0-4502-baf3-cc6ad35fb578",
  "resource_id": "api.example.com:443",
  "os": "macos",
  "app_version": "v0.46.0",
  "client_name": "sbx",
  "hostname": "host-machine",
  "approval": {
    "approval_request_id": "b7e2c9d1-4f60-4a8e-93c2-5d1e8f7a6b30",
    "grant_scope": "AUDIT_APPROVAL_GRANT_SCOPE_PERSISTENT"
  },
  "action_type": "http_request",
  "http_request": {
    "method": "HTTP_METHOD_POST",
    "host": "api.example.com",
    "port": "443",
    "path": "/v1/items"
  },
  "agent": "claude"
}
```
