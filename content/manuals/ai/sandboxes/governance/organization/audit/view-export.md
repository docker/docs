---
title: View and export audit events
linkTitle: View and export
weight: 20
description: Search, filter, and export Docker AI Governance audit events from the hosted audit log UI.
keywords: docker sandboxes, audit events, audit logs, AI Governance, CSV export
aliases:
  - /ai/sandboxes/governance/audit/view-export/
---

{{< summary-bar feature_name="AI Governance" >}}

Cloud delivery stores AI Governance audit records in Docker Cloud and makes
them available in the hosted audit log UI. Use the hosted view to investigate
policy decisions, such as why a developer's request was blocked or what
sandboxes reached over the past week. Open the view to see recent events,
narrow it to the ones you care about, then export them to CSV if you need an
offline copy.

## Before you begin

- [Docker Cloud delivery](configure.md) must be turned on for your
  organization. The hosted view shows only records delivered to Docker Cloud,
  not records written to local disk.
- You need organization owner access, or a
  [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
  with AI Governance audit permissions.

## View audit events

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **Audit events**.

The **Audit events** view includes summary tiles for total events, allowed
events, denied events, and consent-required events. The event table includes:

| Column    | Description                                                         |
| --------- | ------------------------------------------------------------------- |
| Time      | When Docker recorded the event.                                     |
| Event     | The event type or policy action.                                    |
| Principal | The Docker user associated with the event.                          |
| Resource  | The target resource, such as a domain, file path, or tool.          |
| Decision  | The governance decision, such as Allow, Deny, or Approval required. |
| Agent     | The AI agent associated with the event, when Docker knows it.       |

For every field in a record and what its values mean, see
[Audit record reference](../../reference/audit-records.md).

## Filter and search events

The event table can hold a large number of events, so narrow it to the ones
you're investigating. Use the filters above the table:

- **Date range**: Today, Last 7 days, Last 14 days, Last 30 days, or a custom
  range. The table shows the last 7 days by default.
- **Decision**: Allow, Deny, Approval required, or Approval unavailable.
- **Event type**: the kind of action, such as a network or filesystem
  evaluation. For the full list, see
  [Action types](../../reference/audit-records.md#action-types).
- **Agent**: the AI agent associated with the event.
- **Principal**: a Docker username.
- **Session ID**: an audit session ID. This filter shows the whole session,
  whatever the date range.

To search across events, principals, resources, and decisions, enter at least
three characters in the search box, then press Enter.

For example, set **Decision** to Deny and **Principal** to a developer's
username to see what their sandboxes were blocked from.

The table shows 25 events per page by default. You can show 100 per page
instead.

## Export events to CSV

When you need an offline copy of the events you found, export them. The export
uses your current filters and search terms:

1. Open **Audit events**.
1. Apply the filters and search terms for the events you want to export.
1. Select **Export**.
1. Download the generated CSV file from the link Docker provides.

CSV exports include up to 1,000,000 rows. An export stays available to download
for 24 hours after Docker generates it.

To monitor audit events continuously instead of searching them by hand,
[forward them to your SIEM](siem.md).
