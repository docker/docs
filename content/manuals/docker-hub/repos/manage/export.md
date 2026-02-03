---
title: Export organization repositories to CSV
linkTitle: Export repositories
description: Learn how to export a complete list of your organization's Docker Hub repositories using the API.
keywords: docker hub, organization, repositories, export, csv, api, access token
---

This guide shows you how to export a complete list of repositories from your
Docker Hub organization, including private repositories. You'll use an
Organization Access Token (OAT) to authenticate with the Docker Hub API and
export repository details to a CSV file for reporting or analysis.

The exported data includes repository name, visibility status, last updated
date, pull count, and star count.

## Prerequisites

Before you begin, ensure you have:

- Administrator access to a Docker Hub organization
- `curl` installed for making API requests
- `jq` installed for JSON parsing
- A spreadsheet application to view the CSV

## Create an organization access token

Organization access tokens let you authenticate API requests without
interactive login steps.

1. Navigate to your organization in [Docker Home](https://app.docker.com) and select **Admin Console**.

2. Select **Access tokens** from the sidebar.

3. Select **Generate access token**.

4. Configure the token permissions:

   - Under **Repository permissions**, add every repository you want the
     token to access
   - Assign at least **Image Pull** (read) access to each repository
   - You can add up to 50 repositories per token

5. Copy the generated token and store it securely.

> [!IMPORTANT]
>
> If you only enable **Read public repositories**, the API will only return
> public repositories. To include private repositories in your export, you must
> explicitly add them to the token's repository permissions.

## Authenticate with the Docker Hub API

Exchange your organization access token for a JWT bearer token that you'll use
for subsequent API requests.

1. Set your organization name and access token as variables:

   ```bash
   ORG="<your-org>"
   OAT="<your_org_access_token>"
   ```

2. Call the authentication endpoint to get a JWT:

   ```bash
   TOKEN=$(
     curl -s https://hub.docker.com/v2/users/login \
       -H 'Content-Type: application/json' \
       -d "{\"username\":\"$ORG\",\"password\":\"$OAT\"}" \
     | jq -r '.token'
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

Your organization access token may only have **Read public repositories**
enabled, or it lacks permissions for specific private repositories.

To fix this:

1. Navigate to your organization's access tokens in Docker Hub
2. Select the token you created
3. Add private repositories to the token's permissions with at
   least **Image Pull** access
4. Regenerate the JWT and retry the export

### API returns 403 or missing fields

Ensure you're using the JWT from the `/v2/users/login` endpoint as a
Bearer token in the `Authorization` header, not the organization access
token directly.

Verify your authentication:

```console
$ curl -s "https://hub.docker.com/v2/namespaces/$ORG/repositories?page_size=1" \
  -H "Authorization: Bearer $TOKEN" | jq
```

If this returns an error, re-run the authentication step to get a fresh JWT.

### Need access to all repositories

Organization access tokens are scoped to specific repositories you select
during token creation. To export all repositories, you have two options:

1. Add all repositories to the organization access token (up to 50 repositories)
2. Use a Personal Access Token (PAT) from an administrator account that has
   access across the entire organization

The choice between these approaches depends on your organization's security
policies.
