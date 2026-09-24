---
title: "Find a process you lost track of"
linkTitle: "Find a process you lost track of"
description: "Look up a running process by the session tag you gave it when you created it, then pick up its output from the resume point its own report publishes."
keywords: "cloud sandboxes, sandboxes api, find a process you lost track of"
weight: 701
params:
  sidebar:
    group: "Long-running work"
---

Reconnect to work after losing a process connection. Starting a new command can duplicate effects, so first look for the process you already started.

Connections opened through TypeScript process handles recover from temporary disconnects while you consume output. They reconnect to the same process and resume after the last delivered chunk, with a 30-second budget per recovery episode. They never restart the command or replay terminal input. Raw streams require explicit reconnection. Use this guide when automatic recovery stops, you close a connection, or your application restarts.

You need a current sandbox handle and a session tag that your application assigned when it started the process. A session tag is searchable metadata, not the process's resource name.

## Find the running session {#1-find-the-running-session}

List processes filtered by session and running state. The example follows pagination. If multiple processes match, choose the intended resource by its saved name or metadata rather than attaching to an arbitrary result.

An empty list means no running process matched. It does not prove that the original command never started: it may have exited.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.processes
  .all({ filter: `session=${session},state=running` })
  .collect();
```

<details>
<summary>Complete TypeScript example: reconnect/list.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function findSession(sandbox: Sandbox, session: string) {
  return sandbox.processes
    .all({ filter: `session=${session},state=running` })
    .collect();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Reconnect to its output {#2-reconnect-to-its-output}

Connect to the selected process handle and consume output events. The example resumes from the last sequence reported by that handle.

For delivery to another system, keep your own cursor for the last chunk that system actually handled. The latest sequence reported by the process may be ahead of that cursor. Use [explicit output resumption](run-an-interactive-shell-in-a-cloud-sandbox.md) when you need to replay from your saved position.

Close the connection when finished. Keep the process name if the connection drops again; do not replace reconnection with a second process start.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const connection = await process.connect({
  resumeFrom: process.lastStreamSequence,
});
let lastSequence = BigInt(process.lastStreamSequence ?? '0');
try {
  for await (const event of connection) {
    if (event.type === 'chunk')
      write(
        event.data ?? new Uint8Array(),
        (lastSequence = BigInt(event.streamSequence ?? '0')),
      );
  }
  return lastSequence;
} finally {
  await connection.close();
}
```

<details>
<summary>Complete TypeScript example: reconnect/attach.ts</summary>

```typescript
import type { Process } from '@docker/sandboxes';

export async function rejoinProcess(
  process: Process,
  write: (bytes: Uint8Array, sequence: bigint) => void,
) {
  const connection = await process.connect({
    resumeFrom: process.lastStreamSequence,
  });
  let lastSequence = BigInt(process.lastStreamSequence ?? '0');
  try {
    for await (const event of connection) {
      if (event.type === 'chunk')
        write(
          event.data ?? new Uint8Array(),
          (lastSequence = BigInt(event.streamSequence ?? '0')),
        );
    }
    return lastSequence;
  } finally {
    await connection.close();
  }
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
