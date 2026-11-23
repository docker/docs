---
title: Docker Sandboxes API v2
description: Manage local and cloud sandbox environments through the redesigned Docker Sandboxes API.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 10
linkTitle: v2
params:
  apiReferenceBanner: false
---

The Docker Sandboxes API v2 is a programmatic interface for managing sandbox
environments locally and in Docker Cloud. Use it to automate development tasks,
manage the environments you work in, or build applications that use sandboxes.

> [!NOTE]
> This documentation is in development. Guides are outlines, and runnable
> examples, SDK installation instructions, and the API reference are deferred
> until the interface is stable.

## Start locally, then use cloud

The getting started sequence begins with a local sandbox: connect to the daemon,
create a shell environment, run a command, and delete the sandbox. The cloud
guide follows the same sequence with cloud authentication and a cloud endpoint.

- [Run a local sandbox](get-started/local.md)
- [Run a cloud sandbox](get-started/cloud.md)

The API shares operations across local and cloud environments. Available
features and configuration options differ between them. See
[Local and cloud capabilities](capabilities.md).

## Compose your environment

A sandbox can combine a workload kit with mixin kits. The workload supplies
the base environment and launch command. Mixins contribute tools,
configuration, or runtime behavior. The API prepares the environment from
published kit references.

This section uses v3 kits. API versions and kit format versions are independent:
API v2 does not mean v2 kits. See
[Customize an environment](recipes/customize-an-environment.md) and the
[kit documentation](/manuals/ai/sandboxes/customize/_index.md).

## Automate sandbox workflows

Start with the [recipes](recipes/_index.md) for managing environments,
working with project files, and connecting agents to services. Each recipe
identifies the execution environment it covers.

Use [Connect and authenticate](authentication.md) to understand local and cloud
connections. Read [Handle failures and cleanup](errors.md) before building
unattended workflows.

## API and SDK versions

V2 has a different interface from [API v1](../v1/_index.md). The v1 SDK and
cookbook examples are not v2 instructions. See [SDKs](sdks.md) for the scope of
the planned SDK documentation.

<!-- Before publication: confirm the release name and wire API version. The
implementation reviewed for this draft exposes /v0, not /v2. Do not derive
request paths or SDK package names from this documentation's v2 label. -->
