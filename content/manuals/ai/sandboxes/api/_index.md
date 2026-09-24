---
title: Docker Sandboxes API
linkTitle: API and SDK
description: Use the Docker Sandboxes API and TypeScript SDK to manage cloud sandboxes programmatically.
keywords: docker sandboxes API, cloud sandboxes API, sandbox SDK, TypeScript SDK, JavaScript
weight: 35
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

> [!NOTE]
> The Docker Sandboxes API and SDK are experimental. Features, interfaces,
> and behavior may change.

Use the Docker Sandboxes API to create cloud sandboxes, run commands, and
transfer files from your applications and automated workflows. You can also
manage related resources, including images, snapshots, volumes, and secrets.

Start from a kit that packages an agent or tool environment. The SDK manages
the resulting sandbox and its processes, files, and lifecycle.

To try it, [run your first cloud sandbox](get-started.md) with the TypeScript SDK.

## Activate cloud access

To use the API or SDK, [activate a Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access).
Use the same Docker account to [authenticate your application](authentication.md).

Cloud compute is billed separately from your Docker subscription. See
[Billing](/manuals/agentic-platform/signup.md#billing) for details.

## Choose an interface

Use the TypeScript SDK in your JavaScript or TypeScript application. The SDK
provides typed requests and responses, waits for sandboxes to start or stop,
and handles file transfers and interactive processes.
See [Install the SDK](sdks.md) for installation instructions.

You can also call the REST API directly from any language or HTTP tool.

## Develop your application

- [API concepts](concepts.md): connect to the right endpoint, identify
  resources, and wait for actions to finish.
- [Authentication and authorization](authentication.md): authenticate requests
  and understand which permissions your application needs.
- [Errors and retries](errors.md): handle failures and retry requests without
  duplicating work.
- [Compute sizes and limits](limits.md): choose resources and handle account
  quotas and request rate limits.
