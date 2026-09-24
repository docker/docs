---
title: "Run your first command"
linkTitle: "Run your first command"
description: "Connect to a running sandbox, run a command in it, and tell a command that exited unsuccessfully apart from a call that failed."
keywords: "cloud sandboxes, sandboxes api, run your first command"
weight: 105
params:
  sidebar:
    group: "Get started"
---

Run a command in an existing sandbox and collect its output. Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md) and the sandbox name returned by [your first launch](create-your-first-sandbox.md).

For most commands, use the process collection's run helper. It starts the process and waits for its result.

## Open the sandbox {#1-open-the-sandbox}

Get the sandbox by its resource name. The returned handle exposes processes and files and manages their authenticated connections. You do not need to copy the sandbox address or Docker token into a second client.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.get(name);
```

<details>
<summary>Complete TypeScript example: exec/endpoint.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function processesAt(client: Sandboxes, name: string) {
  return client.get(name);
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Use the process run helper {#2-use-the-process-run-helper}

Pass an argument array: the program followed by its arguments, such as `['echo', 'Hello from Docker Sandboxes']`. Arguments are not shell syntax. For pipes, redirection, or variable expansion, run `sh -c` with your script as the next argument.

The run helper returns the process result after collecting output. A successful SDK call does not mean the command succeeded: inspect its exit code. Standard error is diagnostic output and can be nonempty even when the exit code is zero.

An SDK error means the request, wait, or connection failed. Keep the process handle carried by an error when available so you can inspect it instead of starting duplicate work. [Error handling](handle-errors-and-degradation.md) explains the distinction.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.processes.run({ args }, { timeoutMs: 300_000 });
```

<details>
<summary>Complete TypeScript example: exec/outcome.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function classifyOutcome(sandbox: Sandbox, args: string[]) {
  return sandbox.processes.run({ args }, { timeoutMs: 300_000 });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Optional: execute a command directly {#3-optional-execute-a-command-directly}

For a short command whose output you only need after it exits, direct execution is simpler: one request waits for the command and returns its captured output and exit code, without a separate process-creation and output-reading sequence.

Check the incomplete flag before treating captured output as complete. Choose [streaming output](run-something-that-produces-real-output.md) when you need output while the command runs.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const result = await endpoint.api.exec({ body: { cmd: args } });
```

<details>
<summary>Complete TypeScript example: exec/exec.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function runCommand(sandbox: Sandbox, args: string[]) {
  const endpoint = await sandbox.endpointClient(['sandboxesExec'], 'exec');
  try {
    const result = await endpoint.api.exec({ body: { cmd: args } });
    return {
      exitCode: result.exitCode ?? 0,
      stdout: new TextDecoder().decode(result.stdout),
      stderr: new TextDecoder().decode(result.stderr),
      incomplete: result.incomplete ?? false,
    };
  } finally {
    await endpoint.close();
  }
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
