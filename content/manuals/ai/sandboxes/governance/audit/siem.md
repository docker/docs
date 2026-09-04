---
title: SIEM forwarding
linkTitle: SIEM forwarding
weight: 40
description: Forward Docker AI Governance audit events to Splunk, Dynatrace, or Datadog.
keywords: docker sandboxes, SIEM, audit logs, Splunk, Dynatrace, Datadog, AI Governance, forwarding, NDJSON
---

{{< summary-bar feature_name="AI Governance Audit Logs" >}}

Docker can forward audit events to your security information and event
management (SIEM) system, letting you centralize Docker governance data
alongside other security signals. Docker verifies the endpoint is reachable
with the supplied credential before saving.

## Supported destinations

| Destination                      | Description                                                     |
| -------------------------------- | --------------------------------------------------------------- |
| Splunk Cloud (HEC)               | Hosted Splunk using the HTTP Event Collector                    |
| Dynatrace                        | Dynatrace Log Management using the Log Ingest API               |
| Datadog                          | Datadog Logs using the HTTP log intake API                      |

## Before you begin

SIEM forwarding requires Docker Sandboxes
[0.39.0](/manuals/ai/sandboxes/release-notes.md) or later. Earlier versions
don't deliver audit records to a SIEM destination, even when forwarding is
configured. Update Docker Sandboxes before enabling a new destination.

SIEM forwarding requires Docker Cloud delivery to be enabled for your
organization. If you haven't already, enable it under **AI Platform** >
**Audit logs** > **Audit delivery** before configuring a SIEM destination. See
[Configure audit delivery](configure.md).

Gather credentials from your SIEM before configuring forwarding:

- **Splunk Cloud**: HEC ingest URL and an HEC token. Optionally, a Splunk index
  name. See [Splunk documentation](https://docs.splunk.com/).
- **Dynatrace**: Log Ingest API URL and an API token with the `logs.ingest`
  scope. See [Dynatrace documentation](https://docs.dynatrace.com/).
- **Datadog**: Logs intake URL for your Datadog site and an API key. See
  [Datadog documentation](https://docs.datadoghq.com/).

## Add a SIEM destination

1. Sign in to [Docker Home](https://app.docker.com/).
1. Open your organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **Export & Connectors**.
1. Select **Add destination**.
1. Select your destination and complete the form.
1. Select **Save**.

If verification fails, check that the URL and credential are correct and that
the endpoint is accessible from the internet.

## Manage destinations

From the **SIEM forwarding** list, select the menu next to a destination to
edit or delete it. The edit form lets you update credentials and toggle
forwarding on or off for that destination. Deleting a destination permanently
removes the endpoint and its stored credential and cannot be undone.
