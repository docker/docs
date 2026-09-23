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

To use the API or an SDK, [activate a Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access).
Use the same Docker account to [create your API access token](authentication.md#create-an-access-token).

Cloud compute is billed separately from your Docker subscription. See
[Billing](/manuals/agentic-platform/signup.md#billing) for details.

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
