---
title: Usage reports
linkTitle: Usage reports
description: Learn how to retrieve enterprise usage reports for your Docker organization using the Reports API.
keywords: docker, enterprise, reports, usage, pulls, api, csv, organization access token, oat
toc_max: 4
weight: 50
---

{{< summary-bar feature_name="Usage reports" >}}

Docker generates daily pull activity reports for your organization as CSV downloads, accessible through the Reports API.

## Prerequisites

- A [Docker Business subscription](https://www.docker.com/pricing/)
- Organization owner role, or a [custom role](/manuals/enterprise/security/roles-and-permissions/custom-roles.md)
  with the report read permission

> [!NOTE]
> Personal Access Tokens
> (PATs) are unsupported for usage reports.

## Retrieve and download usage reports

These procedures walk you through fetching and downloading usage reports for your organization. 

### Step 1. Create an OAT for the Reports API

OATs are
org-scoped tokens designed for machine-to-machine access, making them
suitable for automated report retrieval workflows.

1. Go to [Docker Home](https://app.docker.com/) to [create your OAT](/enterprise/security/access-tokens/#create-an-organization-access-token). To add the report read scope:
    - Go to **Resources** and select the **Organization scopes** drop-down.
    - Select **Report read** from the drop-down.
2. Set your variables so the report read OAT is associated with the correct organization:

   ```bash
   ORG="<your-org-name>"
   OAT="<your-organization-access-token>"
   ```
3. Exchange the OAT for a JWT bearer token:

   ```bash
   TOKEN=$(curl -s https://hub.docker.com/v2/auth/token \
     -H 'Content-Type: application/json' \
     -d "{\"identifier\":\"$ORG\",\"secret\":\"$OAT\"}" \
     | jq -r '.access_token')
   ```

4. Validate the JWT bearer token against the list reports endpoint:

   ```console
   $ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports" \
     -H "Authorization: Bearer $TOKEN" | jq .
   ```

You use this `TOKEN` value in the `Authorization: Bearer` header for all
subsequent API calls. 

> [!IMPORTANT]
> The JWT token expires after a period, 
> so re-run the OAT exchange step to
> refresh the JWT token.

### Step 2. List usage reports types

#### List available reports 

Fetch the available report types and cadences for your organization:

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

#### List reports for a type and cadence

List reports for a specific type and cadence. Results are in reverse chronological order:

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

#### Set pagination for queries

Results are paginated with a default page size of 30 and a maximum of 100.

Use the `page_size` and `page_token` query parameters to control pagination: 

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily?page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

When `next_page_token` is non-empty, pass it as `page_token` to fetch the next page:

```console
$ curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily?page_size=10&page_token=NEXT_TOKEN" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Step 3. Download a report

#### Download a specific report

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

#### Download latest report

Combine the list and download steps to always fetch the most recent report:

```bash
DATE=$(
  curl -s "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily?page_size=1" \
    -H "Authorization: Bearer $TOKEN" \
  | jq -r '.reports[0].Date // empty'
)

if [ -z "$DATE" ]; then
  echo "No reports available."
  exit 1
fi

curl -L -o "usage_pulls_${DATE}.csv" \
  "https://api.docker.com/enterprise-data/v1/orgs/$ORG/reports/usage_pulls/daily/${DATE}/download" \
  -H "Authorization: Bearer $TOKEN"
```

## Get report schema

Fetch the schema for a specific report date to discover available columns and their types:

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

## API reference

| Endpoint                                                                      | Description                      |
| ----------------------------------------------------------------------------- | -------------------------------- |
| `GET /enterprise-data/v1/orgs/{org}/reports`                                  | List available report types      |
| `GET /enterprise-data/v1/orgs/{org}/reports/{type}/{cadence}`                 | List reports with pagination     |
| `GET /enterprise-data/v1/orgs/{org}/reports/{type}/{cadence}/{date}/download` | Download a report (302 redirect) |
| `GET /enterprise-data/v1/orgs/{org}/reports/{type}/{cadence}/{date}/schema`   | Get report column schema         |

For the full API specification, see the
[Enterprise Data API reference](/reference/api/enterprise-data/latest/).

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
