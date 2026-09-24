---
title: Errors and retries
description: Handle Docker Sandboxes API errors, recover from failed waits, and retry requests without duplicating work.
keywords: docker sandboxes API errors, REST errors, API retries, idempotency keys, rate limits
weight: 40
---

> [!NOTE]
> The Docker Sandboxes API and SDK are experimental. Features, interfaces,
> and behavior may change.

Before retrying a failed request, check whether the service already started
the work. For example, a create request can succeed even if your application
loses the response. Retrying without checking can create a second sandbox.

An API request can succeed even when the command it runs fails. Check the
command result's exit code and output separately from request and wait errors.
A nonzero exit code reports a command failure.

## Read an error response

API errors contain a `code`, a `message`, and optional typed `details`. Use the
code to decide how to respond. Several codes share an HTTP status, so the
status alone might not explain the failure.

| Code | What to do |
| --- | --- |
| `invalidArgument` | Correct the malformed request or unsupported value before retrying. |
| `unauthenticated` | Obtain a valid credential to replace the missing, invalid, or expired one. |
| `permissionDenied` | Check that your credentials have permission for the action. |
| `notFound` | Check the resource name and request URL. Cloud also returns this code for routes it doesn't serve. |
| `failedPrecondition` | Check the resource state, required features, and any `If-Match` header. |
| `resourceExhausted` | Check the error details for a quota, rate limit, or request-size limit. |
| `unimplemented` | Check whether the backend supports the requested feature. |
| `unavailable` | Retry after a delay, once you know the retry won't duplicate work. |

For more detail, inspect `google.rpc.ErrorInfo` when it is present. Its `reason`
and `domain` fields identify documented causes. Base your error handling on
these fields and the error code rather than free-form message text. Your
client should also handle responses containing detail types it doesn't
recognize.

## Recover from a failed wait

If the API returns HTTP 202, the work is still in progress. Keep reading the
resource until it reaches the state you need or fails. The resource's `failure`
field describes a failure that occurs after the initial request succeeds.

If a wait times out or is canceled, the action can still finish. Inspect the
resource before trying again or deleting it. TypeScript wait helpers report
`WaitError`, which includes the last resource the client received and
any failure details. If the client never received a resource, use the original
request and idempotency key to recover it as described in
[Retry without duplicating work](#retry-without-duplicating-work).

When using `withSandbox`, inspect `WorkflowError.phase`, `resource`, and
`cleanup` to find the failed step and whether the SDK deleted the sandbox.

## Handle concurrent changes

To avoid changing a resource that someone else has modified, send its `etag`
in the `If-Match` header. An etag identifies the version of the resource you
read. Operations such as sandbox updates and deletion require this header.
SDK resource handles send the etag from the resource you used to create the
handle.

Pass the etag exactly as returned, including its quotes. A missing required
header returns HTTP 428, and a stale etag returns HTTP 412. If the etag is
stale, read the resource again and decide whether your change is still
appropriate before submitting another request. In the SDK, `refresh()` returns
a separate handle. Use that returned handle for the next operation. The
original handle still has the old etag. If you supply an idempotency key, use
a different key for the request with the updated etag.

## Retry without duplicating work

To retry a request such as sandbox creation without creating another sandbox,
include an `Idempotency-Key` header in the original request. Keep that key and
the exact request, including any `If-Match` value, for retries. For supported
operations, the service keeps accepted results for at least 24 hours and
returns the original response when you repeat the request with the same key.
Read the resource afterward to check its latest state.

Use the same key only when repeating the same request. Changing the request
under that key causes an error, and using a different key submits another
action. Send the header only for operations that support it:

| Operations | `Idempotency-Key` |
| --- | --- |
| Create a sandbox, image, process, port, snapshot, secret, or volume | Optional |
| Restore a snapshot | Optional |
| Start, stop, or delete a sandbox | Optional |
| Update a secret | Optional |
| Update a sandbox | Required |

Other operations don't accept an idempotency key. The SDK generates a key for
each supported mutation unless you supply one.

Once you know a create request succeeded, poll the returned resource to wait
for completion. Retry only when you need to recover from a failure or a lost
response. For temporary failures, wait between retries and limit the number
of attempts.

Process creation supports an idempotency key. Keep the returned process name
and find that process after a lost response before starting another one.
This does not make process input, signals, or file writes safe to replay.
Check the outcome before repeating those actions.

## Account for SDK retries

The SDK makes up to two additional attempts for eligible transient failures.
Use its retry policy, or disable automatic retries when your application owns
the retry loop. Set `maxRetries: 0` for that call. Combining both policies can
produce more attempts than you intended.

Keep the idempotency key across application-level retries. Automatic retries
within a call reuse its key, but a separate create call can generate a new
one. Use a deadline and a bounded attempt count, and honor server retry delays.
See [Request rate limits](limits.md#request-rate-limits) for how rate limits
differ from resource quotas.

## Set a timeout for commands

`processes.run()` waits for a command to finish without a default overall
timeout. Individual requests to create the process and read its output have
a 30-second timeout. To bound the whole run, pass a timeout in the second
argument, for example `sandbox.processes.run(input, { timeoutMs: 300_000 })`.
Timing out or canceling the call stops local waiting. It doesn't kill the
process in the sandbox.
