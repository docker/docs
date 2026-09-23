---
title: Docker Sandboxes API
linkTitle: API and SDKs
description: Use the Docker Sandboxes API and SDKs to manage cloud sandboxes programmatically.
keywords: docker sandboxes API, cloud sandboxes API, sandbox SDK, Go SDK, Python SDK, TypeScript SDK
weight: 35
---

Use the Docker Sandboxes API to create cloud sandboxes, run commands, and
transfer files from your applications and automated workflows. You can also
manage related resources, including images, snapshots, volumes, and secrets.

To try it, [run your first cloud sandbox](get-started.md) with the TypeScript SDK.

## Activate cloud access

To use the API or an SDK, you need a Docker Personal or Docker Pro account with
an active cloud sandbox subscription. Cloud compute is billed separately from
your Docker plan on a pay-as-you-go basis.

Signup and billing use Docker Agentic Platform, so you'll see that name in the
web console, checkout, and your active plans. To activate access:

1. Open the [Docker Agentic Platform console](https://agentic-platform.docker.com/)
   and sign in with the Docker account you want to use for cloud sandboxes.
2. If your account doesn't have access, follow the redirect to Docker Billing
   and complete checkout for the Docker Agentic Platform pay-as-you-go plan.
   If you see **Docker Agentic Platform access required** instead of a redirect,
   select **Go to Docker Billing**.
3. Return to the console after checkout. If your account already has an active
   subscription, you can skip checkout.

Use the same Docker account to [create your API access token](authentication.md#create-an-access-token).

## Choose an interface

Use an SDK for Go, TypeScript, or Python to work with the API in your
application. The SDKs provide typed requests and responses, wait for sandboxes
to start or stop, and handle file transfers and interactive processes.
See [SDKs](sdks.md) for installation instructions.

You can also call the REST API directly from any language or HTTP tool. The
[OpenAPI specification](https://github.com/docker/sandboxes-api/blob/main/openapi/sandboxes-v1.openapi.yaml)
describes its requests and responses.

## Develop your application

- [API concepts](concepts.md): connect to the right endpoint, identify
  resources, and wait for actions to finish.
- [Authentication and authorization](authentication.md): authenticate requests
  and understand which permissions your application needs.
- [Errors and retries](errors.md): handle failures and retry requests without
  duplicating work.
