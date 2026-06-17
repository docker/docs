---
title: Usage reports
description: Learn how to retrieve enterprise usage reports for your Docker organization using the Reports API.
keywords: docker, enterprise, reports, usage, pulls, api, csv, personal access token, pat, organization access token, oat
weight: 20
params:
  sidebar:
    group: Enterprise
---

Docker provides daily usage reports for organizations with a Docker Business
subscription. These reports contain pull activity data for your organization and
are available as CSV downloads through the Reports API.

Reports are generated automatically. You use the API to list what is available
and download the files you need.

## Prerequisites

Before you begin, ensure you have:

- A [Docker Business subscription](/subscription/core-subscription/details/)
- One of the following:
  - Organization owner role
  - A custom role that includes the **report-read** permission
- `curl` installed for making API requests
- `jq` installed for JSON parsing (optional, for formatting responses)

## Authentication

The Reports API accepts two authentication methods: Personal Access Tokens (PATs)
and Organization Access Tokens (OATs). Both are sent directly as Bearer tokens.

### Option A: Personal Access Token

Use a PAT from a user account that has the required role in the organization.

1. [Create a personal access token](/security/access-tokens/) with
   **Read-only** access or broader.

2. Set your variables:

   ```bash
   ORG="<your-org-name>"
   TOKEN="<your-personal-access-token>"
   ```

3. Test the token:

   ```console
   $ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports" \
     -H "Authorization: Bearer $TOKEN" | jq .
   ```

### Option B: Organization Access Token

Use an OAT issued to the organization. OATs do not require a specific user
account.

1. [Create an organization access token](/security/oat/) with
   the appropriate permissions.

2. Set your variables:

   ```bash
   ORG="<your-org-name>"
   TOKEN="<your-organization-access-token>"
   ```

3. Test the token:

   ```console
   $ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports" \
     -H "Authorization: Bearer $TOKEN" | jq .
   ```

You use this `TOKEN` value in the `Authorization: Bearer` header for all
subsequent API calls.

## List available report types

Discover which report types and cadences are available for your organization.

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Example response:

```json
{
  "report_types": [
    {
      "ReportType": "usage_pulls",
      "Cadence": "daily"
    }
  ]
}
```

Each entry represents a distinct combination of report type and cadence. Use
these values in subsequent calls.

## List reports

List available reports for a given type and cadence. Reports are returned in
reverse chronological order (most recent first).

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Example response:

```json
{
  "reports": [
    {
      "ReportType": "usage_pulls",
      "Cadence": "daily",
      "Date": "2026-06-16",
      "SizeBytes": 48210,
      "Key": "my-org/usage_pulls/daily/2026-06-16.csv"
    },
    {
      "ReportType": "usage_pulls",
      "Cadence": "daily",
      "Date": "2026-06-15",
      "SizeBytes": 51003,
      "Key": "my-org/usage_pulls/daily/2026-06-15.csv"
    }
  ],
  "next_page_token": ""
}
```

### Pagination

Results are paginated with a default page size of 30 and a maximum of 100.
Use the `page_size` and `page_token` query parameters to control pagination.

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily?page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

When `next_page_token` is non-empty, pass it as `page_token` to fetch the next
page:

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily?page_size=10&page_token=NEXT_TOKEN" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## Download a report

Download the CSV file for a specific date. The API responds with a `302`
redirect to a pre-signed URL. With `curl -L`, the redirect is followed
automatically and the file is saved locally.

```console
$ curl -L -o "usage_pulls_2026-06-16.csv" \
  "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily/2026-06-16/download" \
  -H "Authorization: Bearer $TOKEN"
```

The pre-signed download URL expires after 15 minutes. If the link expires,
call the endpoint again to get a fresh URL.

### Download the latest report

Combine the list and download steps to always fetch the most recent report:

```bash
DATE=$(
  curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily?page_size=1" \
    -H "Authorization: Bearer $TOKEN" \
  | jq -r '.reports[0].Date'
)

curl -L -o "usage_pulls_${DATE}.csv" \
  "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily/${DATE}/download" \
  -H "Authorization: Bearer $TOKEN"
```

## Get report schema

Retrieve the schema for a specific report date. The schema describes the columns
in the CSV file.

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily/2026-06-16/schema" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Example response:

```json
{
  "category": "usage_pulls",
  "fields": [
    {
      "name": "date",
      "type": "string",
      "description": "The date of the pull event (YYYY-MM-DD)."
    },
    {
      "name": "repository",
      "type": "string",
      "description": "The repository that was pulled."
    },
    {
      "name": "pull_count",
      "type": "integer",
      "description": "Number of pulls for the repository on this date."
    }
  ]
}
```

Use the schema endpoint to programmatically discover column names and types
before processing a report.

## API reference

| Endpoint | Description |
|---|---|
| `GET /enterprise-data/v1/orgs/{org}/reports` | List available report types |
| `GET /enterprise-data/v1/orgs/{org}/reports/{type}/{cadence}` | List reports with pagination |
| `GET /enterprise-data/v1/orgs/{org}/reports/{type}/{cadence}/{date}/download` | Download a report (302 redirect) |
| `GET /enterprise-data/v1/orgs/{org}/reports/{type}/{cadence}/{date}/schema` | Get report column schema |

For the full API specification, see the
[Enterprise Data API reference](/reference/api/enterprise-data/).

## Troubleshooting

### 401 Unauthorized

Your token is missing or invalid. Verify that you are passing the token as a
Bearer token in the `Authorization` header and that the token has not expired.

### 403 Forbidden

The authenticated user or token does not have permission to access reports for
this organization. Verify that:

- The user has the organization owner role or a custom role with the
  **report-read** permission.
- The organization has an active Docker Business subscription.

### 404 Not Found

The requested report type, cadence, or date does not exist. Use the list
endpoints to discover available reports before attempting a download.
