---
title: "Create your first sandbox"
linkTitle: "Create your first sandbox"
description: "Launch a kit, run a command in its sandbox, and inspect the result."
keywords: "cloud sandboxes, sandboxes api, create your first sandbox"
weight: 103
params:
  sidebar:
    group: "Get started"
---

Create a sandbox from a kit and run a command inside it. A kit bundles an image and configuration for a tool or agent. The `shell` kit is a useful first choice because a greeting command needs no model-provider credential.

[Install an SDK](../install.md) and [authenticate to Docker](connect-to-cloud-with-a-bearer-token.md) first. Pass the authenticated client to this example along with a kit name and command.

## Launch a kit and run a command {#1-launch-a-kit-and-run-a-command}

Choose `shell` as the kit name and pass `['echo', 'Hello from Docker Sandboxes']` as the command. The example requests 2 CPUs and 4096 MiB of memory. See [compute sizes and account limits](work-within-the-limits.md) when sizing other workloads.

The example launches the kit, waits until the sandbox is running, and runs the command through the sandbox's process collection. The SDK handles the connection and its sandbox-scoped credential.

A successful greeting returns `Hello from Docker Sandboxes` in `stdout` and exit code zero. See [Run your first command](run-your-first-command.md) for argument arrays, shell syntax, and interpreting command results.

Keep the returned sandbox name. This example leaves the sandbox available for the next guides. [Delete it](delete-a-cloud-sandbox.md) when you finish; closing the SDK client does not delete it. If a wait times out after creation was accepted, the sandbox can still exist.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const created = await client.kits.launch(
  kitName,
  { resources: { cpus: 2, memoryMib: 4096 } },
  { timeoutMs, signal },
);
const sandbox = await created.waitUntilRunning({ timeoutMs, signal });
const result = await sandbox.processes.run(
  { args },
  { timeoutMs, signal },
);
return { sandbox, result };
```

<details>
<summary>Complete TypeScript example: quickstart/run.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function createAndRun(
  client: Sandboxes,
  kitName: string,
  args: string[],
  timeoutMs = 300_000,
) {
  const signal = AbortSignal.timeout(timeoutMs);
  const created = await client.kits.launch(
    kitName,
    { resources: { cpus: 2, memoryMib: 4096 } },
    { timeoutMs, signal },
  );
  const sandbox = await created.waitUntilRunning({ timeoutMs, signal });
  const result = await sandbox.processes.run(
    { args },
    { timeoutMs, signal },
  );
  return { sandbox, result };
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Clean up a temporary sandbox automatically {#2-clean-up-a-temporary-sandbox-automatically}

For a single task, the scoped sandbox helper creates a sandbox, runs your callback, and attempts deletion afterward, including when the callback fails. This example explicitly stops the sandbox and waits for that transition before the helper attempts deletion. Supply the command arguments and creation options appropriate for the task; for a published image, set its reference and resources.

Cleanup has its own deadline. If it fails, inspect the workflow error and the retained sandbox identity so you can finish cleanup. Do not treat a client timeout as confirmation that the sandbox was deleted.

Next, [choose an agent kit](add-tools-with-kits.md) or [work with processes](run-your-first-command.md).

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.withSandbox(options, async (sandbox) => {
  const outcome = await sandbox.processes
    .run({ args }, { timeoutMs: 300_000 })
    .then(
      (value) => ({ value }),
      (error: unknown) => ({ error }),
    );
  try {
    const cleanup = { signal: AbortSignal.timeout(30_000) };
    const current = await client.get(sandbox.name, cleanup);
    if (current.uid !== sandbox.uid)
      throw new Error(
        'Sandbox identity changed; refusing to stop a replacement',
      );
    const stopped = await current.stop(cleanup);
    await stopped.waitUntilStopped(cleanup);
  } catch (error) {
    if ('error' in outcome)
      throw new AggregateError(
        [outcome.error, error],
        'Command and stop both failed',
      );
    throw error;
  }
  if ('error' in outcome) throw outcome.error;
  return outcome.value;
});
```

<details>
<summary>Complete TypeScript example: quickstart/scoped.ts</summary>

```typescript
import type { ClientCreateOptions, Sandboxes } from '@docker/sandboxes';

export async function runTemporary(
  client: Sandboxes,
  options: ClientCreateOptions,
  args: string[],
) {
  return client.withSandbox(options, async (sandbox) => {
    const outcome = await sandbox.processes
      .run({ args }, { timeoutMs: 300_000 })
      .then(
        (value) => ({ value }),
        (error: unknown) => ({ error }),
      );
    try {
      const cleanup = { signal: AbortSignal.timeout(30_000) };
      const current = await client.get(sandbox.name, cleanup);
      if (current.uid !== sandbox.uid)
        throw new Error(
          'Sandbox identity changed; refusing to stop a replacement',
        );
      const stopped = await current.stop(cleanup);
      await stopped.waitUntilStopped(cleanup);
    } catch (error) {
      if ('error' in outcome)
        throw new AggregateError(
          [outcome.error, error],
          'Command and stop both failed',
        );
      throw error;
    }
    if ('error' in outcome) throw outcome.error;
    return outcome.value;
  });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
