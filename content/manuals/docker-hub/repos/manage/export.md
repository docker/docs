---
title: Export organization repositories to CSV
linkTitle: Export repositories
description: Learn how to export a complete list of your organization's Docker Hub repositories using the API.
keywords: docker hub, organization, repositories, export, csv, api, personal access token, pat
---

This guide shows you how to export a complete list of repositories from your
Docker Hub organization, including private repositories. You'll use a
Personal Access Token (PAT) from an administrator account to authenticate with
the Docker Hub API and export repository details to a CSV file for reporting or
analysis.

The exported data includes repository name, visibility status, last updated
date, pull count, and star count.

## Prerequisites

Before you begin, ensure you have:

- Administrator access to a Docker Hub organization
- `curl` installed for making API requests
- `jq` installed for JSON parsing
- A spreadsheet application to view the CSV

## Create a personal access token

[Create a personal access token](/security/access-tokens/) from
a user account that has access to the organization's repositories. When creating
the token, select at minimum **Read-only** access permissions to list
repositories.

> [!IMPORTANT]
>
> Use a PAT from a user account that is a member of the organization. Users
> with owner roles can export all organization repositories. Members can only
> export repositories they have permission to access.

## Authenticate with the Docker Hub API

Exchange your personal access token for a JWT bearer token that you'll use
for subsequent API requests.

1. Set your Docker Hub username, organization name, and personal access token as variables:

   ```bash
   USERNAME="<your-docker-username>"
   ORG="<org-name>"
   PAT="<your_personal_access_token>"
   ```

2. Call the authentication endpoint to get a JWT:

   ```bash
   TOKEN=$(
     curl -s https://hub.docker.com/v2/auth/token \
       -H 'Content-Type: application/json' \
       -d "{\"identifier\":\"$USERNAME\",\"secret\":\"$PAT\"}" \
     | jq -r '.access_token'
   )
   ```

3. Verify the token was retrieved successfully:

   ```console
   $ echo "Got JWT: ${#TOKEN} chars"
   ```

You'll use this JWT as a Bearer token in the `Authorization` header for all
subsequent API calls.

## Retrieve all repositories

The Docker Hub API paginates repository lists. This script retrieves all pages
and combines the results.

1. Set the page size and initial API endpoint:

   ```bash
   PAGE_SIZE=100
   URL="https://hub.docker.com/v2/namespaces/$ORG/repositories?page_size=$PAGE_SIZE"
   ```

2. Paginate through all results:

   ```bash
   ALL=$(
     while [ -n "$URL" ] && [ "$URL" != "null" ]; do
       RESP=$(curl -s "$URL" -H "Authorization: Bearer $TOKEN")
       echo "$RESP" | jq -c '.results[]'
       URL=$(echo "$RESP" | jq -r '.next')
     done | jq -s '.'
   )
   ```

3. Verify the number of repositories retrieved:

   ```console
   $ echo "$ALL" | jq 'length'
   ```

The script continues requesting the `next` URL from each response until
pagination is complete.

## Export to CSV

Generate a CSV file with repository details that you can open in
spreadsheet applications.

Run the following command to create `repos.csv`:

```bash
echo "$ALL" | jq -r '
  (["namespace","name","is_private","last_updated","pull_count","star_count"] | @csv),
  (.[] | [
    .namespace, .name, .is_private, .last_updated, (.pull_count//0), (.star_count//0)
  ] | @csv)
' > repos.csv
```

Verify the export completed:

```console
$ echo "Rows:" $(wc -l < repos.csv)
```

Open the `repos.csv` file in your preferred
spreadsheet application to view and analyze your repository data.

## Troubleshooting

### Only public repositories appear

The Docker Hub account associated with your personal access token may not have
access to private repositories in the organization.

To fix this:

1. Verify the account is a member of the organization
2. Check that the account has appropriate permissions (owner or member role)
3. Ensure the personal access token has sufficient access permissions
4. Regenerate the JWT and retry the export

### API returns 403 or missing fields

Ensure you're using the JWT from the `/v2/auth/token` endpoint as a
Bearer token in the `Authorization` header, not the personal access
token directly.

Verify your authentication:

```console
$ curl -s "https://hub.docker.com/v2/namespaces/$ORG/repositories?page_size=1" \
  -H "Authorization: Bearer $TOKEN" | jq
```

If this returns an error, re-run the authentication step to get a fresh JWT.
