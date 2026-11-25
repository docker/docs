---
title: Docker Sandboxes API and SDK
linkTitle: Sandboxes API and SDK
description: Use the Docker Sandboxes API and TypeScript SDK to manage cloud sandboxes programmatically.
keywords: docker sandboxes API, cloud sandboxes API, sandbox SDK, TypeScript SDK, JavaScript
weight: 15
params:
  apiReferenceBanner: false
  sidebar:
    group: AI and agents
aliases:
  - /ai/sandboxes/api/
grid:
  - title: "SDK cookbook"
    description: "Follow examples for processes, files, storage, networking, and other sandbox operations."
    link: /ai/sandboxes-api/cookbook/
    icon: code
  - title: "API concepts"
    description: "Choose a kit or image, identify resources, and wait for actions to finish."
    link: /ai/sandboxes-api/concepts/
    icon: book-open
  - title: "Authentication and authorization"
    description: "Authenticate requests and understand which permissions your application needs."
    link: /ai/sandboxes-api/authentication/
    icon: settings
  - title: "Errors and retries"
    description: "Handle failures and retry requests without duplicating work."
    link: /ai/sandboxes-api/errors/
    icon: search
  - title: "Compute sizes and limits"
    description: "Choose resources and handle account quotas and request rate limits."
    link: /ai/sandboxes-api/limits/
    icon: cloud
---

Use the Docker Sandboxes API to create cloud sandboxes, run commands, and
transfer files from your applications and automated workflows. You can also
manage related resources, including images, snapshots, volumes, and secrets.

To try it, [run your first cloud sandbox](get-started.md) with the TypeScript
SDK. The tutorial uses a bundled kit that supplies a shell environment for
running commands.

## Activate cloud access

To use the API or SDK, [activate a Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access).
Use the same Docker account to [authenticate your application](authentication.md).

Cloud compute is billed separately from your Docker subscription. See
[Billing](/manuals/agentic-platform/signup.md#billing) for details.

## Choose an interface

Use the TypeScript SDK in your JavaScript or TypeScript application. The SDK
provides typed requests and responses, waits for sandboxes to start or stop,
and handles file transfers and interactive processes.
See [Install the SDK](install.md) for installation instructions.

You can also call the REST API directly from any language or HTTP tool.
See the [API reference](/reference/api/sandboxes/latest/) for operations,
request fields, responses, and the downloadable OpenAPI specification.

To run agents from your terminal, see [Docker Sandboxes](../sandboxes/_index.md).

## Develop your application

{{< grid >}}
