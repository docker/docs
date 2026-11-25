---
title: Run a cloud sandbox with API v2
description: Repeat the local sandbox workflow against Docker Cloud.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 20
linkTitle: Run a cloud sandbox
---

This guide will repeat the [local getting started workflow](local.md) with a
sandbox running in Docker Cloud. The task stays the same: create a shell
sandbox, run a command, inspect the result, and delete it.

Cloud requests use a service endpoint and Docker account credentials. Cloud
access and billing are managed through
[Docker Agentic Platform](/manuals/agentic-platform/signup.md).

## Workflow

1. Activate cloud access for your Docker account.
2. Authenticate your application and select the cloud endpoint.
3. Create a sandbox using a shell workload kit.
4. Run the same command as in the local guide.
5. Delete the sandbox and confirm cleanup.

The authentication instructions and runnable example are deferred. This guide
creates a separate cloud sandbox. Moving an existing sandbox between local
and cloud is a different workflow.

<!-- Complete after release validation: cloud entitlement and billing wording,
supported authentication, matching local/cloud examples, readiness, bounded
waiting, and cleanup. Do not assume the v1 SDK can target this endpoint. -->
