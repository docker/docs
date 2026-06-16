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
- An [organization access token (OAT)](/manuals/security/for-admins/access-tokens/)

## Available report types

### Usage pulls

Daily aggregated pull activity for your organization. Each report covers one
calendar day (UTC) and includes all authenticated pulls attributed to your
organization.

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

All endpoints require authentication with an
[organization access token (OAT)](/manuals/security/for-admins/access-tokens/)
in the `Authorization` header.

Base URL: `https://api.docker.com`

### Discover available reports

Returns the report types and cadences available for your organization.

```
GET /v2/orgs/{org_name}/reports
```

**Response:**

```json
{
  "report_types": [
    {"ReportType": "usage_pulls", "Cadence": "daily"}
  ]
}
```

### List reports

Returns metadata for available reports of a given type and cadence.

```
GET /v2/orgs/{org_name}/reports/{type}/{cadence}
```

**Query parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `page_token` | No | Cursor for the next page (from a previous response) |

**Response:**

```json
{
  "reports": [
    {
      "ReportType": "usage_pulls",
      "Cadence": "daily",
      "Date": "2026-06-01",
      "SizeBytes": 15234
    }
  ],
  "next_page_token": "..."
}
```

Reports are returned most recent first. When `next_page_token` is present,
pass it as `page_token` in the next request to get more results.

### Download a report

Downloads a specific report CSV. The response is a redirect (HTTP 302) to
a short-lived download URL. Most HTTP clients follow the redirect
automatically.

```
GET /v2/orgs/{org_name}/reports/{type}/{cadence}/{date}/download
```

The download URL expires after 15 minutes.

### Get report schema

Returns the column definitions for a specific report date. Use this to
understand the structure of the CSV before parsing it, or to detect when
columns have been added or changed.

```
GET /v2/orgs/{org_name}/reports/{type}/{cadence}/{date}/schema
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
# Set your org access token
export TOKEN="dckr_oat_..."

# 1. Discover what report types are available
curl -s -H "Authorization: Bearer $TOKEN" \
  https://api.docker.com/v2/orgs/your-org/reports | jq .

# 2. List available daily usage reports
curl -s -H "Authorization: Bearer $TOKEN" \
  https://api.docker.com/v2/orgs/your-org/reports/usage_pulls/daily | jq .

# 3. Download the most recent report
curl -L -o usage_pulls_2026-06-01.csv -H "Authorization: Bearer $TOKEN" \
  https://api.docker.com/v2/orgs/your-org/reports/usage_pulls/daily/2026-06-01/download

# 4. Check the column schema
curl -s -H "Authorization: Bearer $TOKEN" \
  https://api.docker.com/v2/orgs/your-org/reports/usage_pulls/daily/2026-06-01/schema | jq .
```

> [!TIP]
>
> Use `curl -L` for the download endpoint. The API returns a redirect (HTTP 302)
> to the file, and `-L` tells curl to follow it.

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
