---
title: Handle failures and cleanup in API v2
description: Plan for request failures, asynchronous operations, and sandbox cleanup.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 50
linkTitle: Handle failures and cleanup
---

A request failure, an unfinished operation, and a command's nonzero exit code
need different handling. Keep enough information to inspect accepted work
before starting it again.

## Wait for accepted work

An accepted request can leave work in progress. Wait for the state your task
needs before running dependent operations. A client deadline can stop your
application from waiting while the server continues the work.

## Retry deliberately

Before retrying a write, determine whether the original request was accepted.
Use idempotency only on operations that support it, and retain the same request
identity when repeating the same work.

## Plan cleanup

Closing a connection and deleting a sandbox are separate actions. Make sandbox
ownership and cleanup explicit in your application. An SDK can provide cleanup
helpers, but their behavior must be checked for the SDK and version you use.

Check a command's exit code separately from the success of the API request.
A successful request can report a failed command.

<!-- Add verified error shapes, operation waiting, retry rules, deadlines, and
cleanup behavior once the API and SDK contracts are stable. Do not carry v1
error types, retry defaults, or resource-handle behavior into these examples. -->
