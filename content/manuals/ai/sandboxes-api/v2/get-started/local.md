---
title: Run a local sandbox with API v2
description: Create a local shell sandbox, run a command, and clean up through the API.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 10
linkTitle: Run a local sandbox
---

This guide will introduce the API using a sandbox on your machine. A shell
workload keeps the first task focused on creating an environment and running
a command.

The local daemon manages the sandbox. Your application connects through the
daemon's local connection rather than a cloud service endpoint. See
[Connect and authenticate](../authentication.md).

## Workflow

1. Install a compatible version of `sbx` and start its daemon.
2. Connect to the local API and check the operations it supports.
3. Create a sandbox using a shell workload kit.
4. Run a command and inspect its output and exit code.
5. Delete the sandbox when the work is complete.

The connection commands and runnable example are deferred. After this guide,
continue with [Run a cloud sandbox](cloud.md).

<!-- Complete after release validation: supported sbx versions and platforms,
feature flags, sign-in prerequisites, IPC discovery, create readiness, command
execution, cleanup on failure, and expected output. Use synchronous execution,
which is shared by the local and cloud implementations. -->
