---
title: Docker Sandboxes SDK cookbook
linkTitle: Cookbook
description: Build agent workflows with the Docker Sandboxes SDK for TypeScript.
keywords: cloud sandboxes, sandboxes sdk, typescript, agent kits
weight: 60
params:
  recipeCatalog: true
  sidebar:
    activeChildrenOnly: true
---

Use these recipes to add file transfers, processes, storage, and other sandbox
operations to your application. Each recipe shows the relevant SDK calls and
an expandable complete TypeScript example.

Start with [Get started](../get-started.md) to run your first sandbox, or
[install the SDK](../sdks.md) to use these examples in an existing project.
For HTTP operations and request fields, see the
[API reference](/reference/api/sandboxes/latest/).

## Get started

Authenticate and launch a kit before exploring individual SDK operations.

{{< recipe-list group="Get started" >}}

## Working in a sandbox

Use a running sandbox for commands, project files, and web applications.

{{< recipe-list group="Working in a sandbox" >}}

## Packaging and state

Keep an environment or its state for later work.

{{< recipe-list group="Packaging and state" >}}

## Security and policy

Give agents the access they need without embedding credentials in application code.

{{< recipe-list group="Security and policy" >}}

## Connect and configure

Choose a different image, control lifetime, or connect additional storage and clients.

{{< recipe-list group="Connect and configure" >}}

## Requests and responses

Handle request options and failures deliberately.

{{< recipe-list group="Requests and responses" >}}

## Long-running work

Continue work across lost connections or coordinate several sandboxes.

{{< recipe-list group="Long-running work" >}}
