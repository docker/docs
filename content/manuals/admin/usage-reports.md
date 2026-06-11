---
title: Usage reports
weight: 45
description: Access daily usage reports for your organization via the Docker API.
keywords: organization, usage reports, pull data, enterprise, export, API, CSV
---

{{< summary-bar feature_name="Usage reports" >}}

Usage reports give organization administrators programmatic access to daily
Docker Hub pull activity for their organization. Reports are generated
automatically and available for download as CSV files through the Docker API.

Key benefits include:

- **Chargeback and cost allocation**: Track pull activity by user, repository,
  and IP address to attribute usage to internal teams.
- **Adoption tracking**: Monitor which images and repositories your organization
  uses most.
- **Security and compliance**: Audit who pulled what, when, and from where.
- **Self-service**: Download reports directly through the API without contacting
  Docker support.

## Prerequisites

To access usage reports, you must meet the following requirements:

- [Docker Business subscription](https://www.docker.com/pricing?ref=Docs)
- Organization owner or a custom role with report access permissions

## Available report types

### Usage pulls

Daily aggregated pull activity for your organization. Each report covers one
calendar day (UTC) and includes all authenticated pulls attributed to your
organization's billing entity.

**Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `day` | date | Date of the pull activity (UTC) |
| `hub_username` | string | Docker Hub username that performed the pull |
| `image_repository` | string | Repository name, excluding namespace and tag |
| `namespace` | string | Docker Hub namespace (organization or personal account) |
| `tag` | string | Image tag that was pulled |
| `image_digest` | string | SHA256 content digest of the image manifest |
| `privacy` | string | Repository visibility at time of pull (`public` or `private`) |
| `billing_entity` | string | Customer billing entity the pull is attributed to |
| `auth_name` | string | Authenticated identity used for the pull (username or token name) |
| `ip_address` | string | Client IP address |
| `os` | string | Client operating system reported by the Docker client |
| `user_agent` | string | Docker client user agent string |
| `version_checks` | integer | Number of manifest check (HEAD) requests |
| `data_downloads` | integer | Number of layer download (GET) requests |
| `event_count` | integer | Total pull events (`version_checks` + `data_downloads`) |
| `egress_size_bytes` | integer | Total egress in bytes across all layer downloads |

> [!NOTE]
>
> Reports are generated daily at 16:00 UTC. Data for a given day becomes
> available the following day after the generation job completes. Column
> definitions may change over time; use the schema endpoint to get the current
> column definitions for any report.

## API reference

The usage reports API uses [ConnectRPC](https://connectrpc.com/) (compatible
with gRPC and JSON). All endpoints require authentication with a valid
organization access token or Hub session JWT.

Base URL: `https://marlin-api.docker.com`

### List reports

Returns metadata and download links for available reports.

**Endpoint:**
`POST /docker.marlin.query.v1.ReportService/ListReports`

**Request:**

```json
{
  "orgId": "your-org-id",
  "reportType": "usage_pulls",
  "cadence": "daily",
  "pageSize": 30,
  "pageToken": ""
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `orgId` | Yes | Your organization ID |
| `reportType` | Yes | Report type (e.g., `usage_pulls`) |
| `cadence` | Yes | Report cadence (`daily`) |
| `pageSize` | No | Number of reports per page (default: 30, max: 100) |
| `pageToken` | No | Cursor for the next page (from a previous response) |

**Response:**

```json
{
  "reports": [
    {
      "reportType": "usage_pulls",
      "cadence": "daily",
      "date": "2026-06-01",
      "sizeBytes": "15234",
      "downloadUrl": "https://...",
      "downloadUrlExpiresAt": "2026-06-02T12:15:00Z"
    }
  ],
  "nextPageToken": "..."
}
```

Each report includes a `downloadUrl` that can be used to download the CSV file
directly. The URL expires after 15 minutes. To get a fresh URL, call
`ListReports` again.

When `nextPageToken` is present, pass it as `pageToken` in the next request to
get the next page of results. Reports are returned most recent first.

### Get report schema

Returns the column definitions for a specific report. Use this to understand the
structure of the CSV before parsing it, or to detect when columns have been
added or changed.

**Endpoint:**
`POST /docker.marlin.query.v1.ReportService/GetReportSchema`

**Request:**

```json
{
  "orgId": "your-org-id",
  "reportType": "usage_pulls",
  "cadence": "daily",
  "date": "2026-06-01"
}
```

**Response:**

```json
{
  "category": "usage_pulls",
  "fields": [
    {
      "name": "day",
      "type": "date",
      "description": "Date of the pull activity (UTC)"
    },
    {
      "name": "hub_username",
      "type": "string",
      "description": "Docker Hub username that performed the pull"
    }
  ]
}
```

> [!NOTE]
>
> Schema is per-report-date. If columns are added or removed over time,
> historical reports retain the schema they were generated with. Always check
> the schema for the specific date you are downloading.

## Example: download a report

```bash
# 1. List available reports
curl -s https://marlin-api.docker.com/docker.marlin.query.v1.ReportService/ListReports \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"orgId":"your-org-id","reportType":"usage_pulls","cadence":"daily","pageSize":7}' \
  | jq '.reports[] | {date, sizeBytes}'

# 2. Download the most recent report
URL=$(curl -s https://marlin-api.docker.com/docker.marlin.query.v1.ReportService/ListReports \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"orgId":"your-org-id","reportType":"usage_pulls","cadence":"daily","pageSize":1}' \
  | jq -r '.reports[0].downloadUrl')

curl -o usage_pulls_latest.csv "$URL"
```

## Data retention

Reports are retained for 30 days. After 30 days, the report files are
permanently deleted and cannot be recovered. Download and archive reports
locally if you need to retain them longer.

## Frequently asked questions

### How often are reports generated?

Reports are generated daily at 16:00 UTC. Each report covers one calendar day
(midnight to midnight UTC).

### Can I get weekly or monthly reports?

Weekly and monthly cadences are planned for a future release. Currently, only
daily reports are available.

### Why are some fields missing from older reports?

Reports retain the schema they were generated with. If a new column is added,
it only appears in reports generated after the change. Use the schema endpoint
to check which columns are available for a specific date.

### Can I request additional columns?

Contact your Docker account team to discuss additional data fields. Some fields
available in the underlying data (such as IP geolocation and cloud provider)
can be added on request.
