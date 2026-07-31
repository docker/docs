---
title: View and export audit events
linkTitle: View and export
weight: 30
description: Search, filter, export, and stream Docker AI Governance audit events from the hosted audit log UI.
keywords: docker sandboxes, audit events, audit logs, AI Governance, CSV export, SIEM, Splunk, NDJSON
---

Cloud delivery stores AI Governance audit records in Docker Cloud and makes
them available in the hosted audit log UI. Use the hosted view to investigate
policy decisions, export events to CSV, or stream events to your SIEM.

## View audit events

To view audit events:

1. Sign in to [Docker Home](https://app.docker.com/).
1. Open your organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **Audit Events**.

The **Audit Events** view includes summary tiles for total events, allowed
events, denied events, and consent-required events. The event table includes:

| Column    | Description                                                   |
| --------- | ------------------------------------------------------------- |
| Time      | When Docker recorded the event.                               |
| Event     | The event type or policy action.                              |
| Principal | The Docker user associated with the event.                    |
| Resource  | The target resource, such as a domain, file path, or tool.    |
| Decision  | The governance decision, such as allow, deny, or consent.     |
| Agent     | The AI agent associated with the event, when Docker knows it. |

## Filter and search events

Use the audit log filters to narrow the event table by decision and time range.
Use search to find events by principal, resource, event type, or agent.

The event table uses cursor pagination for large result sets.

## Export events to CSV

Use CSV export when you need an offline copy of filtered audit events:

1. Open **Audit Events**.
1. Apply the filters and search terms for the events you want to export.
1. Select **Export**.
1. Download the generated CSV file from the link Docker provides.

CSV exports include up to 1 000 000 rows. Download links expire after 24 hours.

## SIEM forwarding

Docker can forward audit events to your security information and event
management (SIEM) system, letting you centralize Docker governance data
alongside other security signals in tools such as Datadog or Splunk.

To set up forwarding, you provide an HTTPS endpoint that will receive the audit
events, choose an output format, and supply the credentials Docker will use to
authenticate against your SIEM. Docker verifies the endpoint is reachable before
saving the configuration.

Once configured, forwarding can be enabled or disabled without losing your saved
settings.

To collect host-local files with your own log shipper instead, see
[Local audit logs](local.md).
