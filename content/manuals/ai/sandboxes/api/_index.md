---
title: Docker Sandboxes API and SDK
linkTitle: API and SDK
description: Use the Docker Sandboxes API and TypeScript SDK to manage cloud sandboxes programmatically.
keywords: docker sandboxes API, cloud sandboxes API, sandbox SDK, TypeScript SDK, JavaScript
weight: 60
params:
  apiReferenceBanner: false
  sidebar:
    badge:
      color: violet
      text: Experimental
aliases:
  - /ai/sandboxes-api/
---

> [!NOTE]
> The Docker Sandboxes API and SDK are experimental. Features, interfaces,
> and behavior may change.

Use the Docker Sandboxes API to create cloud sandboxes, run commands, and
transfer files from your applications and automated workflows. You can also
manage related resources, including images, snapshots, volumes, and secrets.

To try it, [run your first cloud sandbox](/manuals/ai/sandboxes/api/get-started.md) with the TypeScript
SDK. The tutorial uses a bundled kit that supplies a shell environment for
running commands.

## Activate cloud access

To use the API or SDK, [activate a Docker Agentic Platform subscription](/manuals/ai/sandboxes/cloud-access/_index.md#activate-cloud-access).
Use the same Docker account to [authenticate your application](/manuals/ai/sandboxes/api/authentication.md).

Cloud compute is billed separately from your Docker subscription. See
[Billing](/manuals/ai/sandboxes/cloud-access/_index.md#billing) for details.

## Choose an interface

Use the TypeScript SDK in your JavaScript or TypeScript application. The SDK
provides typed requests and responses, waits for sandboxes to start or stop,
and handles file transfers and interactive processes.
See [Install the SDK](/manuals/ai/sandboxes/api/install.md) for installation instructions.

You can also call the REST API directly from any language or HTTP tool.
See the [API reference](/reference/api/sandboxes/latest/) for operations,
request fields, responses, and the downloadable OpenAPI specification.

To run agents from your terminal, see [Docker Sandboxes](/manuals/ai/sandboxes/_index.md).

## Develop your application

- [SDK cookbook](/manuals/ai/sandboxes/api/cookbook/_index.md): follow examples for processes, files,
  storage, networking, and other sandbox operations.
- [API concepts](/manuals/ai/sandboxes/api/concepts.md): choose a kit or image, identify resources, and
  wait for actions to finish.
- [Authentication and authorization](/manuals/ai/sandboxes/api/authentication.md): authenticate requests
  and understand which permissions your application needs.
- [Errors and retries](/manuals/ai/sandboxes/api/errors.md): handle failures and retry requests without
  duplicating work.
- [Compute sizes and limits](/manuals/ai/sandboxes/api/limits.md): choose resources and handle account
  quotas and request rate limits.
