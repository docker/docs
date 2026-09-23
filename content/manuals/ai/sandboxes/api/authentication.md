---
title: Authentication and authorization
description: Authenticate requests to the Docker Sandboxes API with access tokens and learn which permissions your application needs.
keywords: docker sandboxes API authentication, sandbox authorization, bearer token, sandbox permissions, owner scope
weight: 30
---

Authenticate with a Docker Hub access token to create and manage cloud
sandboxes. To run commands or transfer files inside a sandbox, use a separate
token that grants access to that sandbox. The SDK can obtain this second token
for you.

Your account and credential type must be enabled for Cloud Sandboxes. Each
request also requires permission for the action you want to perform.

## Create an access token

Exchange your Docker ID and personal access token (PAT) for a short-lived
access token using the
[Docker Hub authentication API](/reference/api/hub/latest/operations/AuthCreateAccessToken/):

```console
$ ACCESS_TOKEN=$(curl --silent --show-error --fail --request POST \
  --url https://hub.docker.com/v2/auth/token \
  --header "Content-Type: application/json" \
  --data '{"identifier":"<DOCKER_ID>","secret":"<PERSONAL_ACCESS_TOKEN>"}' \
  | jq -er '.access_token')
```

Use the returned `access_token` to authenticate management requests. The
Sandboxes API doesn't accept a PAT, organization access token (OAT), or
password directly.

## Authenticate management requests

Send the access token in the `Authorization` header when calling
`https://connect.docker.com/sandboxes`:

```http
Authorization: Bearer <access_token>
```

Supply this header when you configure your SDK client. Your application must
also obtain a replacement token when it expires. The SDK doesn't refresh
tokens or load credentials saved by `sbx login`. If a request returns HTTP 401
`unauthenticated`, obtain a valid token before retrying.

## Authenticate sandbox requests

To run a command or transfer a file, your application needs a token for that
sandbox with permission for the action. For example, `sandboxesExec` permits
command execution and `sandboxesFilesWrite` permits writing files.

In TypeScript, call `client.forEndpoint(endpoint, permissions)` on your
management client. Pass the sandbox's `core.endpoint` and the permissions you
need. The SDK requests a short-lived token, checks that it is valid for the
sandbox and endpoint, and creates a client that uses it. The caller needs
`sandboxesCredential` to obtain this token, as well as the permissions it
requests.

Use the returned client for process and file requests, including streams.
Keep its token private and out of logs. The Docker Hub token used by the
management client must not be sent to a sandbox endpoint.

## Resource access and permissions

Your credentials determine which account's resources you can access. The API
calls this the caller's owner scope. Cloud resolves this scope from your
credentials, so leave the optional `parent` field empty in requests.

Each request also checks whether you have permission for the action on the
target resource. For example, creating a sandbox requires `sandboxesCreate`,
reading it requires `sandboxesRead`, and deleting it requires
`sandboxesDelete`.

See the [Cloud support guide](https://github.com/docker/sandboxes-api/blob/main/CLOUD_SUPPORT.md)
for the permissions and service configuration that each feature requires.
