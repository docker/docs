---
title: "Run work across many sandboxes"
linkTitle: "Run work across many sandboxes"
description: "Run one command on every running sandbox you hold, a bounded number at a time, and read one answer per sandbox even when some of them fail."
keywords: "cloud sandboxes, sandboxes api, run work across many sandboxes"
weight: 703
params:
  sidebar:
    group: "Long-running work"
---

Run a command across several existing sandboxes while keeping concurrency bounded. This can collect diagnostics or apply the same task to a group of independent workspaces.

Use an [authenticated client](connect-to-cloud-with-a-bearer-token.md) and an argument array. Only run the command in sandboxes your application owns for this task.

## Select running sandboxes {#1-select-running-sandboxes}

Walk the sandbox collection and retain running handles. The example selects all running sandboxes visible to the client; narrow that selection to your application's jobs before executing a command with side effects.

Listing and execution are separate requests. A sandbox can change state between them, so the execution step still needs to handle failures.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const running: Sandbox[] = [];
for await (const sandbox of client.all()) {
  if (sandbox.status === 'running') running.push(sandbox);
}
return running;
```

<details>
<summary>Complete TypeScript example: fanout/spread.ts</summary>

```typescript
import type { Sandbox, Sandboxes } from '@docker/sandboxes';

export async function runningSandboxes(client: Sandboxes) {
  const running: Sandbox[] = [];
  for await (const sandbox of client.all()) {
    if (sandbox.status === 'running') running.push(sandbox);
  }
  return running;
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Collect individual results {#2-collect-individual-results}

Choose a concurrency bound and run the command on each selected sandbox. The example keeps each result or failure alongside the sandbox name, so one failed request does not discard every other result.

Inspect both request failures and process exit codes. A returned process result can contain a nonzero exit code.

Keep concurrency below your account's limits and your application's memory budget. For large output, stream results instead of collecting them all. The example does not delete the existing sandboxes.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const outcomes: Outcome[] = [];
const size = Number.isFinite(bound) ? Math.max(1, Math.trunc(bound)) : 1;
for (let start = 0; start < sandboxes.length; start += size) {
  const group = sandboxes.slice(start, start + size);
  const settled = await Promise.allSettled(
    group.map((sandbox) => sandbox.processes.run({ args })),
  );
  settled.forEach((result, index) => {
    outcomes.push(
      result.status === 'fulfilled'
        ? { sandboxName: group[index].name, response: result.value }
        : { sandboxName: group[index].name, failure: result.reason },
    );
  });
}
return outcomes;
```

<details>
<summary>Complete TypeScript example: fanout/gather.ts</summary>

```typescript
import type { RunResult, Sandbox } from '@docker/sandboxes';

export type Outcome = {
  sandboxName: string;
  response?: RunResult;
  failure?: unknown;
};

export async function runOnEach(
  sandboxes: Sandbox[],
  args: string[],
  bound: number,
) {
  const outcomes: Outcome[] = [];
  const size = Number.isFinite(bound) ? Math.max(1, Math.trunc(bound)) : 1;
  for (let start = 0; start < sandboxes.length; start += size) {
    const group = sandboxes.slice(start, start + size);
    const settled = await Promise.allSettled(
      group.map((sandbox) => sandbox.processes.run({ args })),
    );
    settled.forEach((result, index) => {
      outcomes.push(
        result.status === 'fulfilled'
          ? { sandboxName: group[index].name, response: result.value }
          : { sandboxName: group[index].name, failure: result.reason },
      );
    });
  }
  return outcomes;
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
