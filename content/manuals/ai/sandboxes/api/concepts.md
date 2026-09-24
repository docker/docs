---
title: Docker Sandboxes API concepts
linkTitle: API concepts
description: Learn how kits and images define a cloud sandbox, how to connect to it, and how to manage resources throughout their lifecycle.
keywords: docker sandboxes API concepts, sandbox kits, sandbox images, cloud sandbox endpoint, API operations
weight: 20
---

> [!NOTE]
> The Docker Sandboxes API and SDK are experimental. Features, interfaces,
> and behavior may change.

An application uses the Docker Sandboxes API to create sandboxes, connect to
them, and track their state. Choose the environment for your sandbox, then
learn how to work with its resources throughout their lifecycle.

## Kits and sandbox images

A sandbox kit defines an environment for an agent or tool: its container
image, setup, network rules, and credential requirements. You can use a kit
bundled with the SDK or one published separately in an OCI registry.

You can also create a sandbox from a container image. Choose the source based
on how much of the environment you want to configure yourself:

| Source | What it provides | How to create a sandbox |
| --- | --- | --- |
| [Bundled kit](#bundled-kits) | An image reference and configuration included in the SDK's catalog | `client.kits.launch('shell')` |
| Kit from another source | An environment defined by a kit you obtain separately, such as a published Hermes kit | `client.create()` with [kit artifact inputs](#supply-kit-artifacts) |
| Registry image (`imageRef`) | A container image to use with your own sandbox settings | `client.create({ imageRef: 'ubuntu:24.04', resources: 'small' })` |
| Image resource (`image`) | An image already prepared for Cloud Sandboxes, including its compute settings | `client.create({ image: 'images/<uid>' })` |

The `imageRef` value is an image name in a registry. The `image` value is a
resource name returned by the Sandboxes API. When you use `image`, omit
`resources` because the image resource supplies its compute settings.

A sandbox kit supplies its own image reference. Cloud Sandboxes pulls that
container image as needed, separately from loading the kit's configuration.
Bundling a kit with the SDK doesn't bundle its container image.

Both `create()` and `kits.launch()` return after creation is accepted. Call
`waitUntilRunning()` on the returned sandbox before running commands. For a
bundled kit, `kits.launchAndWait()` combines creation and waiting in one call.

### Bundled kits

The npm package includes the following kit definitions and their supporting
files. These copies are tied to the SDK release. Launch a bundled kit by its
short name, such as `shell`, without downloading the kit from a registry:

| Kit name | Environment |
| --- | --- |
| `shell` | A shell environment for running your own commands |
| `claude` | Claude Code |
| `codex` | Codex |
| `cursor` | Cursor |
| `devin` | Devin |
| `docker-agent` | Docker Agent |
| `gemini` | Gemini CLI |
| `opencode` | OpenCode |

For example, `client.kits.launchAndWait('shell')` creates a shell sandbox and
waits until it's running. The kit launch helpers default to Small compute,
with two CPUs and 4 GiB of memory. To see the catalog bundled with your
installed SDK version, call `client.kits.list()`.

The launch helpers prepare the bundled kit's content and pass it to
`client.create()` through the `kits` field. Omit `image` and `imageRef` when
using these helpers because the kit supplies the image reference.

To run an AI agent, also provide credentials for the service that supplies its
models. For example, Claude Code can use an Anthropic API key, and Codex can
use an OpenAI API key. Signing in to Docker gives you access to sandboxes.
The provider key gives the agent access to its models. See
[Authenticate agents](authentication.md#authenticate-agents) for how to supply
these credentials. The `shell` kit needs no provider key to run commands.

### Supply kit artifacts

You can also use kits published by Docker, your organization, or other
authors. For example, the
[Hermes agent kit](https://hub.docker.com/r/sbx/hermes-agent-kit) is distributed
as `docker.io/sbx/hermes-agent-kit:latest`. This reference identifies a kit
artifact in a registry. The kit, in turn, names the container image to run.

To supply a kit outside the bundled catalog, pass its content to
`client.create()` in the `kits` array. Each entry contains a kit reference and
the resolved kit artifact serialized as JSON bytes. A sandbox kit defines
the environment, and mixin kits add configuration to it. See
[Kits](../customize/kits.md) for how these kinds of kits work together.

The `kits` field is a low-level API input. It expects the kit definition and
supporting files in the serialized
[v2 artifact format](https://github.com/docker/sbx-kits-contrib/blob/v0.17.0/spec/types.go).
The API doesn't accept v3 kit descriptors directly. The npm SDK has no
helper to load registry kits or convert them into this input, so using them
through the SDK requires loading and conversion code outside the SDK.
The `kits.launch()` and `kits.launchAndWait()` helpers accept only names
from the bundled catalog.

This function shows the expanded request structure for a prepared sandbox
kit. It takes the kit's source reference and serialized artifact as inputs:

```typescript
function createFromKit(reference: string, artifactBytes: Uint8Array) {
  return client.create({
    resources: 'small',
    kits: [
      {
        artifact: {
          ref: { ref: reference, kind: 'sandbox' },
          inline: artifactBytes,
        },
      },
    ],
  });
}
```

The `ref` identifies the kit's source, such as the Hermes registry reference.
It doesn't trigger a registry pull. The `inline` value carries the prepared
JSON artifact, including the kit's file content, as a `Uint8Array`. Passing
raw `spec.yaml`, a kit ZIP, or an OCI manifest in this field isn't supported.

To launch a public kit directly by its registry reference, you can use the
[Docker Agentic Platform Console](/manuals/agentic-platform/kits.md#run-a-kit-by-reference),
which loads the kit for you. Its backend can load v2 and v3 kits and converts
compatible v3 kits into the v2 artifact format that Cloud Sandboxes accepts.
This conversion doesn't support every v3 capability.

## Management and sandbox endpoints

Creating a sandbox and running a command inside it use different endpoints:

| Endpoint | Use it to |
| --- | --- |
| Management API at `https://connect.docker.com/sandboxes` | Create, inspect, and delete sandboxes and manage related resources. |
| Sandbox API at the returned `core.endpoint.uri` | Run processes and read or write files inside that sandbox. |

The SDK builds request URLs from these base URLs. If you make HTTP requests
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
need. For kits, `client.kits.launchAndWait()` creates the sandbox and waits
until it's running. If you use `client.create()` or `client.kits.launch()`,
call `waitUntilRunning()` on the returned sandbox before running commands.

A running sandbox can still be completing setup specified by its kit, such
as cloning a repository. Wait for any files or services your workload needs
before starting that work.

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

Cloud supports kits, sandbox timeouts, stored secrets, and volume attachments,
subject to account permissions and feature availability. For example, volume
access must be enabled for your account. An SDK method's presence doesn't
guarantee that your account can use it.

Leave `parent` empty for Cloud requests. The kit launch helpers default to
Small compute. When calling `client.create()` with kit artifacts or a
registry image, select a [compute size](limits.md#compute-sizes) with
`resources`, such as `resources: 'small'`.
