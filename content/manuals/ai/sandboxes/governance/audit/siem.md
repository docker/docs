---
title: SIEM forwarding
linkTitle: SIEM forwarding
weight: 35
description: Forward Docker AI Governance audit events to Splunk, Dynatrace, or a custom HTTPS endpoint.
keywords: docker sandboxes, SIEM, audit logs, Splunk, Dynatrace, AI Governance, forwarding, NDJSON
---

Docker can forward audit events to your security information and event
management (SIEM) system, letting you centralize Docker governance data
alongside other security signals. Events are forwarded in NDJSON format.
Docker verifies the endpoint is reachable with the supplied credential before
saving.

## Supported destinations

| Destination                      | Description                                                     |
| -------------------------------- | --------------------------------------------------------------- |
| Splunk Cloud (HEC)               | Hosted Splunk using the HTTP Event Collector                    |
| Splunk Enterprise (self-hosted)  | Self-hosted Splunk using the HTTP Event Collector               |
| Dynatrace                        | Dynatrace Log Management using the Log Ingest API               |
| Custom HTTPS endpoint (advanced) | Any SIEM that accepts HTTPS with a custom authentication header |

## Before you begin

SIEM forwarding requires Docker Cloud delivery to be enabled for your
organization. If you haven't already, enable it under **AI Platform** >
**Audit logs** > **Audit delivery** before configuring a SIEM destination. See
[Configure audit delivery](configure.md).

Gather credentials from your SIEM before configuring forwarding:

- **Splunk Cloud**: HEC ingest URL and an HEC token. Optionally, a Splunk index
  name. See [Splunk documentation](https://docs.splunk.com/).
- **Splunk Enterprise**: HEC endpoint URL (typically port 8088) and an HEC
  token. The endpoint must present a publicly-trusted TLS certificate.
  Optionally, a Splunk index name. See [Splunk documentation](https://docs.splunk.com/).
- **Dynatrace**: Log Ingest API URL and an API token with the `logs.ingest`
  scope. See [Dynatrace documentation](https://docs.dynatrace.com/).
- **Custom HTTPS endpoint**: Your endpoint URL, authentication header name, and
  full header value including any scheme (for example, `Bearer <token>`).

## Add a SIEM destination

1. Sign in to [Docker Home](https://app.docker.com/).
1. Open your organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **SIEM forwarding**.
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

To collect host-local files with your own log shipper instead, see
[Local audit logs](local.md).
