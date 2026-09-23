---
title: Compute sizes and limits
description: Choose a Cloud Sandboxes compute size and plan for account quotas and request rate limits.
keywords: cloud sandboxes, compute sizes, CPU, memory, quotas, rate limits
weight: 50
---

Choose a compute size for each sandbox and keep your application's resource
usage within your account's quotas. Request rate limits also constrain how
quickly your application can send API requests.

## Compute sizes

Cloud Sandboxes supports these CPU and memory pairs:

| Size | CPUs | Memory | SDK memory value in MiB |
| --- | ---: | ---: | ---: |
| Micro | 1 | 2 GiB | 2048 |
| Small | 2 | 4 GiB | 4096 |
| Medium | 4 | 8 GiB | 8192 |
| Large | 8 | 16 GiB | 16384 |
| XL | 16 | 32 GiB | 32768 |

Set both `resources.cpus` and `resources.memoryMib` when launching a kit or
registry image. CPU and memory aren't independent settings: for example,
1 CPU with 1024 MiB is not a supported pair. When creating from an existing
image resource, omit resource settings because the image supplies them.

Your account access and available capacity determine whether a request can
be accepted. See [Billing](/manuals/agentic-platform/signup.md#billing) for
pricing and usage information.

## Account quotas

The default quotas apply across an account:

| Resource | Default limit |
| --- | ---: |
| Concurrent sandboxes | 10 |
| Stored sandboxes | 50 |
| Volumes | 100 |
| Secrets | 100 |
| Images being prepared at the same time | 3 |

Your account can have different quotas. Confirm your account's limits with
Docker before planning a workload that depends on a particular allowance.

Stopping an ordinary sandbox releases its concurrency slot, but the sandbox
still counts toward stored usage. Starting it again needs a concurrency slot.
An always-on sandbox retains its concurrency reservation while stopped.
Delete sandboxes you no longer need to release stored usage.

Listing sandboxes is not a reservation or an authoritative quota balance.
Other callers can consume capacity between your list and create requests.
Handle a quota refusal even if a preceding check showed available capacity.
Reduce concurrency or remove unused resources before retrying a request that
exceeds your quota.

## Request rate limits

Request rate limits are separate from resource quotas and can vary by
operation. Limit concurrent requests, honor the server's retry delay when
provided, and bound both retry attempts and elapsed time.

Retrying immediately doesn't resolve an exhausted resource quota. Inspect
the error details to determine whether to wait, reduce usage, or correct the
request. See [Errors and retries](errors.md) before retrying a mutation.
