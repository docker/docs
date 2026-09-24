---
title: "Keep a cloud sandbox running"
linkTitle: "Keep a cloud sandbox running"
description: "Choose a sandbox's lifetime and what happens when it ends, then renew the lifetime before the deadline passes."
keywords: "cloud sandboxes, sandboxes api, keep a cloud sandbox running"
weight: 503
params:
  sidebar:
    group: "Connect and configure"
---

Set a sandbox lifetime to prevent abandoned work from running indefinitely, and choose what happens when it expires.

You need an [authenticated client](connect-to-cloud-with-a-bearer-token.md) and a managed image. The same lifecycle options can be supplied when launching a kit.

## Set the lifetime at creation {#1-set-the-lifetime-at-creation}

Pass the lifetime in the units shown by your SDK and choose the expiry action. The example accepts seconds and converts to the SDK's duration representation.

Stopping preserves the sandbox for later use; deletion removes it. Check that your choice matches whether the application needs the sandbox's files after expiry. Account limits may restrict permitted durations and actions.

Accounts with always-on access can instead select the restart action. This preserves memory through automatic stop-and-resume cycles. It retains a concurrency reservation while stopped and is different from restarting in response to an incoming request.

The example waits for startup. Its wait deadline is separate from the sandbox lifetime: ending a client wait does not delete the sandbox.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const sandbox = await client.create(
  {
    displayName: name,
    image,
    lifecycle: {
      timeoutMs: lifeSeconds * 1_000,
      onTimeout: stopOnTimeout ? 'stop' : 'delete',
    },
  },
  { timeoutMs: 300_000, idempotencyKey: requestId },
);
return sandbox.waitUntilRunning();
```

<details>
<summary>Complete TypeScript example: timeout/deadline.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createWithDeadline(
  client: Sandboxes,
  name: string,
  image: string,
  lifeSeconds: number,
  stopOnTimeout: boolean,
  requestId: string,
) {
  const sandbox = await client.create(
    {
      displayName: name,
      image,
      lifecycle: {
        timeoutMs: lifeSeconds * 1_000,
        onTimeout: stopOnTimeout ? 'stop' : 'delete',
      },
    },
    { timeoutMs: 300_000, idempotencyKey: requestId },
  );
  return sandbox.waitUntilRunning();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Extend the remaining lifetime {#2-extend-the-remaining-lifetime}

Use `renewTimeout()` to extend the sandbox's remaining lifetime. Pass the desired duration from the time the service accepts the renewal, not an increment to add to the existing deadline. Renewal can only extend the lifetime; it cannot shorten it.

Renew a configured timeout before it expires. If the call fails, inspect the sandbox's current state instead of assuming that renewal took effect.

The service manages timeouts for automatically restarting sandboxes; do not renew those yourself.

For request-triggered restart after a stop, see [Resume on demand](let-a-stopped-sandbox-resume-on-demand.md).

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.renewTimeout({ timeoutMs: remainingSeconds * 1_000 });
```

<details>
<summary>Complete TypeScript example: timeout/renew.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function renewSandboxTimeout(
  client: Sandboxes,
  name: string,
  remainingSeconds: number,
) {
  const sandbox = await client.get(name);
  return sandbox.renewTimeout({ timeoutMs: remainingSeconds * 1_000 });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
