---
title: Forward audit events to a SIEM
linkTitle: Forward to a SIEM
weight: 30
description: Forward Docker AI Governance audit events to Splunk, Dynatrace, Datadog, Sumo Logic, or CrowdStrike Falcon Next-Gen SIEM.
keywords: docker sandboxes, SIEM, audit logs, Splunk, Dynatrace, Datadog, Sumo Logic, CrowdStrike, Falcon Next-Gen SIEM, AI Governance, forwarding, NDJSON
aliases:
  - /ai/sandboxes/governance/audit/siem/
---

{{< summary-bar feature_name="AI Governance" >}}

Send your Docker AI Governance audit events to your security information and
event management (SIEM) system to centralize Docker governance data alongside
other security signals. How you do it depends on where your audit records are
delivered:

- With Docker Cloud delivery, Docker forwards events directly to a supported
  SIEM destination.
- With local disk delivery, your own log shipper collects the audit files from
  each machine.

## Forward events from Docker Cloud

Docker forwards audit events from Docker Cloud to the SIEM destinations you
add in Docker Home.

### Supported destinations

| Destination                                                                                          | Connects through           | What you need                                                                            |
| ---------------------------------------------------------------------------------------------------- | -------------------------- | ---------------------------------------------------------------------------------------- |
| [Splunk Cloud](https://docs.splunk.com/)                                                             | HTTP Event Collector (HEC) | HEC ingest URL and an HEC token. Optionally, a Splunk index name.                        |
| [Dynatrace](https://docs.dynatrace.com/)                                                             | Log Ingest API             | Log Ingest API URL and an API token with the `logs.ingest` scope.                        |
| [Datadog](https://docs.datadoghq.com/)                                                               | HTTP log intake API        | Logs intake URL for your Datadog site and an API key.                                    |
| [Sumo Logic](https://www.sumologic.com/help/)                                                        | HTTP Source                | HTTP Source URL and an auth token from an HTTP Logs & Metrics source.                    |
| [CrowdStrike Falcon Next-Gen SIEM](https://library.humio.com/logscale-api/log-shippers-hec-raw.html) | HEC / HTTP Event Connector | API URL and an API key from the connector. Append `/raw` to the API URL when you add it. |

### Before you begin

You need organization owner access, or a
[custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
with AI Governance audit permissions.

SIEM forwarding requires Docker Cloud delivery to be enabled for your
organization. If you haven't already, enable it under **AI Platform** > **Audit
logs** > **Audit delivery** before configuring a SIEM destination. For more
details, see [Configure audit delivery](configure.md).

Gather the endpoint URL and credential your destination needs, listed in
Supported destinations. Some destinations require you
to create a connector or source first.

### Add a SIEM destination

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Go to **AI Platform** > **Audit logs**.
1. Open **Export & Connectors**.
1. Select **Add destination**.
1. Select your destination and complete the form.
1. Select **Save**.

Docker verifies the endpoint is reachable with the supplied credential before
saving. If verification fails, check that the URL and credential are correct
and that the endpoint is accessible from the internet.

### Edit or delete a destination

From the **SIEM forwarding** list, select the menu next to a destination to
edit or delete it. The edit form lets you update credentials and toggle
forwarding on or off for that destination. Deleting a destination permanently
removes the endpoint and its stored credential and can't be undone.

## Collect local audit files with a log shipper

With local disk delivery, each machine writes its audit records to `.jsonl`
files in its own audit directory. For the directory on each operating system
and how files are named, see
[Where records are stored](_index.md#where-records-are-stored).

Point your log shipper at the audit directory and configure it to collect
`.jsonl` files only. Tools such as the Splunk Universal Forwarder,
Filebeat, and CrowdStrike Falcon LogScale read the directory and forward each
line as an event. Because in-progress records live in `.tmp` files until they
are finalized, collectors never see partial records.
