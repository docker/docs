---
title: Install the Docker Sandboxes SDK
linkTitle: Install the SDK
description: Install and use the Docker Sandboxes SDK in your JavaScript or TypeScript application.
keywords: docker sandboxes SDK, TypeScript SDK, JavaScript, npm, install sandbox SDK
weight: 10
---

> [!NOTE]
> The Docker Sandboxes API and SDK are experimental. Features, interfaces,
> and behavior may change.

Use the Docker Sandboxes SDK to create and manage cloud sandboxes from
JavaScript or TypeScript. The SDK includes methods for running commands,
transferring files, and waiting for a sandbox to start or stop.

You don't need the Docker CLI to use the SDK.

## Install the SDK

With Node.js 20 or later, install the SDK in your project:

```console
$ npm install @docker/sandboxes
```

Import `Sandboxes` from `@docker/sandboxes`.

## Connect to Cloud Sandboxes

Before connecting, [activate a Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access)
for your Docker account.

Configure your client with browser sign-in for interactive use or a personal
access token for automation. The SDK supplies the service URL and manages
access tokens. See [Authentication and authorization](authentication.md) for
setup instructions.

Use the returned sandbox's methods to run commands and transfer files. The
SDK connects to its endpoint and obtains the required sandbox credentials.

Follow [Run your first cloud sandbox](get-started.md) for a complete TypeScript
example that creates a sandbox, runs a command, and deletes it.

## File transfers and interactive processes

Use the SDK to transfer files over HTTP and interact with running processes
over WebSocket. The package also exports stream helpers from
`@docker/sandboxes/streams`.
