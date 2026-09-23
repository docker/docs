---
title: Docker Sandboxes SDKs
linkTitle: SDKs
description: Install and use the Go, TypeScript, and Python SDKs for the Docker Sandboxes API.
keywords: docker sandboxes SDK, Go SDK, TypeScript SDK, Python SDK, install sandbox SDK
weight: 10
---

Use the Docker Sandboxes SDKs to create and manage cloud sandboxes from Go,
TypeScript, or Python. The SDKs include methods for running commands,
transferring files, and waiting for a sandbox to start or stop.

## Install an SDK

{{< tabs >}}
{{< tab name="Go" >}}

Add the SDK to your Go module:

```console
$ go get github.com/docker/sandboxes-api
```

Import `github.com/docker/sandboxes-api/gen/go/sandboxes` and create a client
with `sandboxes.New`.

{{< /tab >}}
{{< tab name="TypeScript" >}}

With Node.js 20 or later, install the SDK in your project:

```console
$ npm install @docker/sandboxes-api
```

Import `SandboxesClient` from `@docker/sandboxes-api`.

{{< /tab >}}
{{< tab name="Python" >}}

With Python 3.11 or later and Git installed, install the SDK from its source
repository:

```console
$ python -m pip install "git+https://github.com/docker/sandboxes-api.git#subdirectory=gen/python/sandboxes"
```

Import `SandboxesClient` from `docker_sandboxes_api`.

{{< /tab >}}
{{< /tabs >}}

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

Use the SDKs to transfer files over HTTP and interact with running processes
over WebSocket. TypeScript also exports these helpers from
`@docker/sandboxes-api/streams`, and Python exposes them in
`docker_sandboxes_api.streams`.

See the [SDK repository](https://github.com/docker/sandboxes-api) for client
examples and source code.
