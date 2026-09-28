---
title: "Recover when the endpoint moves"
linkTitle: "Recover when the endpoint moves"
description: "Read a sandbox again to learn its current endpoint, then run one command through a connection built for that endpoint."
keywords: "cloud sandboxes, sandboxes api, recover when the endpoint moves"
weight: 704
params:
  sidebar:
    group: "Long-running work"
---

Refresh a sandbox handle after a restart or connection change. A saved address can become stale even though the sandbox still has the same resource name.

Start with the handle your application already holds. Keep its name and UID; do not replace them with another sandbox's identity.

## Compare the current endpoint {#1-compare-the-current-endpoint}

Refresh the handle, then compare the old and new connection information. The SDK checks the resource identity while constructing the refreshed handle.

The example reports whether the endpoint changed and returns a refreshed handle. Use that handle for later process and file calls.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const current = await held.refresh();
const before = held.core.endpoint;
const after = current.core.endpoint;
return {
  moved:
    before?.uri !== after?.uri || before?.protocol !== after?.protocol,
  current,
};
```

<details>
<summary>Complete TypeScript example: rebind/detect.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function endpointMoved(held: Sandbox) {
  const current = await held.refresh();
  const before = held.core.endpoint;
  const after = current.core.endpoint;
  return {
    moved:
      before?.uri !== after?.uri || before?.protocol !== after?.protocol,
    current,
  };
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Run through the refreshed handle {#2-run-through-the-refreshed-handle}

Refresh before opening a new process connection. The SDK obtains the appropriate sandbox credential for the current endpoint; your application does not need to pass its Docker account token there.

This example starts new work after refresh. It does not replay a command whose completion is unknown. If an earlier command may still be running, [find and reconnect to it](find-a-process-you-lost-track-of.md) first.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const current = await held.refresh();
return current.processes.run({ args }, { timeoutMs: 300_000 });
```

<details>
<summary>Complete TypeScript example: rebind/rebuild.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function runAndRebind(held: Sandbox, args: string[]) {
  const current = await held.refresh();
  return current.processes.run({ args }, { timeoutMs: 300_000 });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
