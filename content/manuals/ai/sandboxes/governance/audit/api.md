---
title: Retrieve audit events with the API
linkTitle: Audit Logs API
weight: 45
description: Retrieve and filter Docker AI Governance audit events with the Audit Logs API.
keywords: docker sandboxes, AI Governance, audit logs API, audit events, organization access token, OAT
---

{{< summary-bar feature_name="AI Governance Audit Logs" >}}

The Audit Logs API retrieves cloud-delivered Docker AI Governance events for
automation, investigation, and reporting. The API returns the same event data
shown in the hosted audit log view.

## Prerequisites

Before you use the API, you need:

- A Docker [AI Governance plan](/manuals/subscription/plans/ai-governance.md)
- Docker Cloud delivery turned on in [audit delivery settings](configure.md)
- An Organization Access Token (OAT) with the **View audit logs** permission

Personal Access Tokens and Docker account session tokens aren't supported.

## Create an access token

Create an OAT for the organization whose events you want to retrieve:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Identity & auth**, then **Access tokens**.
1. Select **Generate access token**.
1. Under the AI Governance organization permissions, select **View audit
   logs**.
1. Select **Generate token**, then store the token in a credential manager.

For token expiration, rotation, and management instructions, see
[Organization access tokens](/manuals/enterprise/security/access-tokens.md).

Set variables for the examples:

```bash
ORG="<YOUR_ORGANIZATION_NAME>"
OAT="<YOUR_ORGANIZATION_ACCESS_TOKEN>"
```

Pass the raw OAT in the bearer header. Don't exchange it for another token.

## List audit events

Request the first 25 events for your organization:

```console
$ curl --get "https://api.docker.com/v2/auditlogs/governance/$ORG" \
  --header "Authorization: Bearer $OAT" \
  --data-urlencode "page_size=25" | jq .
```

Events are returned in reverse chronological order. The response contains an
array of events and an opaque cursor for the next page:

```json
{
  "events": [
    {
      "audit_event_id": "95e7257f-93c9-4f29-bde7-88830e2dae80",
      "org_id": "9f8e7d6c-5b4a-3210-fedc-ba9876543210",
      "event_type": "docker.marlin.audit.event.v1.AuditRecord",
      "category": "AUDIT_CATEGORY_EVALUATION",
      "decision": "AUDIT_DECISION_DENY",
      "action_type": "network_egress",
      "enforcement_mode": "enforce",
      "audit_session_id": "8a3bc076-79d0-4502-baf3-cc6ad35fb578",
      "username": "jordandoe",
      "created_at": "2026-08-25T18:42:31.123Z",
      "ingested_at": "2026-08-25T18:42:31.456Z",
      "payload": {
        "resource_id": "example.com:443",
        "agent": "claude",
        "network_egress": { "protocol": "tcp" }
      }
    }
  ],
  "next_page_token": "eyJ0cyI6IjIwMjYtMDgtMjVUMTg6NDI6MzEuMTIzWiIsImVpZCI6Ijk1ZTcyNTdmLTkzYzktNGYyOS1iZGU3LTg4ODMwZTJkYWU4MCJ9"
}
```

See the [audit record reference](record-reference.md) for category, decision,
action, and payload fields.

## Filter events

Use query parameters to narrow the result set:

| Parameter | Match behavior |
| --- | --- |
| `from`, `to` | Events inside an inclusive RFC 3339 time range |
| `audit_session_id` | Exact session ID |
| `username` | Exact Docker username |
| `decision` | Exact audit decision |
| `enforcement_mode` | Exact `audit`, `warn`, or `enforce` mode |
| `action_type` | Exact action type |
| `query` | Case-insensitive text search across username, action type, decision, and event type |
| `resource_id` | Case-insensitive resource ID prefix |
| `agent` | Exact agent value |

For example, retrieve denied network events from a time range:

```console
$ curl --get "https://api.docker.com/v2/auditlogs/governance/$ORG" \
  --header "Authorization: Bearer $OAT" \
  --data-urlencode "decision=AUDIT_DECISION_DENY" \
  --data-urlencode "action_type=network_egress" \
  --data-urlencode "from=2026-08-01T00:00:00Z" \
  --data-urlencode "to=2026-08-25T23:59:59Z" | jq .
```

A `query` search requires both `from` and `to`. The search term must contain at
least three characters, and the time range must not exceed 30 days.

## Retrieve the next page

The API returns up to 25 events by default and accepts a `page_size` up to 500.
When `next_page_token` isn't empty, pass it unchanged as `page_token`:

```console
$ curl --get "https://api.docker.com/v2/auditlogs/governance/$ORG" \
  --header "Authorization: Bearer $OAT" \
  --data-urlencode "page_size=25" \
  --data-urlencode "page_token=$NEXT_PAGE_TOKEN" | jq .
```

Keep the same filters while paging through one result set.

## Request limits

The API applies all of the following request limits:

- 100 requests per minute for each organization
- 600 requests per minute for each source IP address
- 1000 requests per hour for each organization

Requests that exceed a limit receive `429 Too Many Requests` with a JSON error
response. The minute limit controls bursts, while the hourly limit controls
sustained traffic.

## API reference

See the [Audit Logs API reference](/reference/api/ai-governance-audit-logs/)
for the complete request and response schema.
