---
title: Manage existing sandboxes with API v2
description: Find sandbox environments and manage their lifecycle from code.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 20
linkTitle: Manage existing sandboxes
---

Use the API to find sandbox environments and manage their lifecycle from code.
Listing, inspection, starting, stopping, and deletion apply to both local and
cloud sandboxes.

Read a sandbox's identity and state before acting on it. Keep the returned
identity for later requests rather than using a display name as a unique
identifier. Stopping and deleting serve different tasks: stop an environment
you intend to use again, and delete one you no longer need.

Resources visible to your application depend on the endpoint, credentials, and
sandbox compatibility. Do not assume that every sandbox shown by the CLI or
Console is visible through a particular API endpoint.

<!-- Add list/filter/get and stop/start/delete examples. Verify which CLI- and
Console-created sandboxes are visible, local v3 compatibility, state transitions,
concurrency controls, pagination, and ownership before claiming interoperability.
Use a selection criterion that cannot accidentally target unrelated sandboxes. -->
