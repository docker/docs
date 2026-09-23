---
title: Docker Sandboxes API concepts
linkTitle: API concepts
description: Learn how to connect to cloud sandboxes, identify resources, wait for actions to finish, and retrieve paginated results.
keywords: docker sandboxes API concepts, cloud sandbox endpoint, API capabilities, API operations
weight: 20
---

An application uses the Docker Sandboxes API to create sandboxes, connect to
them, and track their state. These concepts explain where to send requests and
how to work with resources throughout their lifecycle.

## Kits and sandbox images

A kit packages an image and configuration for an agent or tool. Use a kit to
start with that environment, then use the SDK to run processes, transfer
files, and manage the sandbox's lifecycle.

The SDK includes a versioned catalog of bundled kits. In TypeScript,
`client.kits.list()` reads that catalog without making an API request.
`client.kits.launch('shell', options)` creates a sandbox from the shell kit.
Wait for the returned sandbox to reach the running state before using it.
Agent kits may also need a model-provider credential.

You can also create a sandbox from a registry image with `imageRef`, or from
an existing image resource with `image`. A named kit supplies its own image,
so don't also pass `image` or `imageRef` to `kits.launch`.

## Management and sandbox endpoints

Creating a sandbox and running a command inside it use different endpoints:

| Endpoint | Use it to |
| --- | --- |
| Management API at `https://connect.docker.com/sandboxes` | Create, inspect, and delete sandboxes and manage related resources. |
| Sandbox API at the returned `core.endpoint.uri` | Run processes and read or write files inside that sandbox. |

The SDKs build request URLs from these base URLs. If you make HTTP requests
directly, append the `/v1` route to the base URL, preserving any existing path.
For example, the management route `/v1/sandboxes` becomes
`https://connect.docker.com/sandboxes/v1/sandboxes`.

Each sandbox endpoint requires a token that grants access to that sandbox.
The SDK obtains this token when you use a sandbox's process or file methods. See
[Authentication and authorization](authentication.md) for details.

A sandbox's endpoint can change when its runtime changes. Read the sandbox
resource again before reconnecting to get its endpoint.

## Resource names

Use a resource's returned `name` to refer to it in later requests. For example,
a sandbox name has the form `sandboxes/<uid>`. The server assigns the ID, which
stays the same throughout the resource's lifetime. Pass the complete name,
including `sandboxes/`, when reading or deleting that sandbox.

The `displayName` field is a label for people to read. Changing this label
doesn't change the resource's `name`.

## Wait for an action to finish

Wait until a sandbox is running before sending commands to it. Creating a
sandbox takes time, so the API can return HTTP 202 with the sandbox still in a
pending state. Read the resource repeatedly until it reaches the state you
need. In TypeScript, call `waitUntilRunning()` on the sandbox returned by
`client.create()`, or use `client.withSandbox()` to wait, run code, and clean up.

Sandbox creation can continue after your client stops waiting. Read the
sandbox again to check its state, and inspect its `failure` field if it has
failed. See [Errors and retries](errors.md) for how to recover.

Deletion can also take time. The API returns HTTP 202 while the sandbox is
being deleted and HTTP 204 when deletion is complete. After deletion,
authorized reads return `notFound`.

## Read all results from a list

List requests return one page of results at a time. To retrieve the next page,
pass the response's `nextPageToken` as the next request's `pageToken`. Keep the
same page size, filter, and ordering. Continue until `nextPageToken` is empty,
even if a page contains fewer items than you requested.

The default page size is 25 for Cloud sandbox, image, snapshot, volume, and
secret lists. You can request up to 100 items per page.

## Choose supported Cloud options

Before adding optional features to a create request, check the
[Cloud support guide](https://github.com/docker/sandboxes-api/blob/main/CLOUD_SUPPORT.md).
It lists supported inputs, required permissions, and account requirements.
The API schema also describes backends other than Cloud, so some of its fields
aren't supported in Cloud requests.

Cloud supports kits, sandbox timeouts, stored secrets, and volume attachments,
subject to account permissions and feature availability. For example, volume
access must be enabled for your account. An SDK method's presence doesn't
guarantee that your account can use it.

Leave `parent` empty for Cloud requests. When launching a kit or registry
image, specify both CPU and memory from a [supported compute size](limits.md#compute-sizes).
An existing image resource supplies its own resources, so omit resource
settings when creating from `image`.
