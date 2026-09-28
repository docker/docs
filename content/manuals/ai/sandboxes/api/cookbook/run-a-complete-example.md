---
title: "Run a complete example"
linkTitle: "Run a complete example"
description: "Sign in, launch the shell kit, print a greeting, and clean up with one runnable program."
keywords: "cloud sandboxes, sandboxes api, run a complete example"
weight: 102
params:
  sidebar:
    group: "Get started"
---

Run a small program that signs in to Docker, creates a sandbox from the `shell` kit, prints a greeting, and deletes the sandbox. You do not need a model-provider key for this example.

[Install your SDK](https://docs.docker.com/ai/sandboxes-api/install/) first. Your Docker account must have Cloud Sandboxes access and billing set up. This program creates a real sandbox with 2 CPUs and 4096 MiB of memory; compute usage is subject to your Docker billing terms.

## Run the program {#1-run-the-program}

Expand the complete example and copy the whole program into a file. Save it as `example.mts`, then run `npx tsx example.mts`.

The program prints a verification URL and code. Open the URL, enter the code, and approve sign-in. Your terminal then shows the sandbox's name and `Hello from Docker Sandboxes`. The program checks the command's exit status and attempts to delete its sandbox before closing the client.

The program waits for sign-in to finish or the verification code to expire. After sign-in, creating the sandbox and running the command share a five-minute time limit. Once the program receives the new sandbox's details, it allows another 30 seconds to delete that sandbox, even if a later step fails.

If sign-in fails, follow the steps in [Authenticate to Docker](connect-to-cloud-with-a-bearer-token.md). If sandbox creation fails, check account access and [resource limits](work-within-the-limits.md).

A creation timeout can occur after Docker creates the sandbox but before the program receives its details. In that case, the program cannot delete it. [List your sandboxes](page-through-and-filter-lists.md) to check for a sandbox created by this run and [delete it](delete-a-cloud-sandbox.md) if needed. If deletion fails or times out, use the sandbox name printed in your terminal to check whether it still exists.

This example deletes its sandbox because it is a one-off demonstration. Next, [keep a sandbox for later work](create-your-first-sandbox.md) or [run an agent kit](add-tools-with-kits.md).

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
const result = await sandbox.processes.run(
  {
    args: ['echo', 'Hello from Docker Sandboxes'],
  },
  operation,
);
console.log(result.stdout);
```

<details>
<summary>Complete TypeScript example: hello/main.ts</summary>

```typescript
import { pathToFileURL } from 'node:url';
import { oauth, Sandboxes, type Sandbox } from '@docker/sandboxes';

export async function main() {
  const auth = oauth({
    onVerification: ({ verificationUri, userCode }) => {
      console.log(`Open ${verificationUri} and enter ${userCode}`);
    },
  });
  const client = new Sandboxes({ auth });
  let sandbox: Sandbox | undefined;
  try {
    await auth.getAccessToken();
    const operation = {
      signal: AbortSignal.timeout(300_000),
      timeoutMs: 300_000,
    };
    sandbox = await client.kits.launch(
      'shell',
      {
        resources: { cpus: 2, memoryMib: 4096 },
      },
      operation,
    );
    console.log(`Sandbox: ${sandbox.name}`);
    sandbox = await sandbox.waitUntilRunning(operation);
    const result = await sandbox.processes.run(
      {
        args: ['echo', 'Hello from Docker Sandboxes'],
      },
      operation,
    );
    console.log(result.stdout);
    if (result.exitCode !== 0)
      throw new Error(`Command exited with status ${result.exitCode}`);
  } finally {
    try {
      if (sandbox) {
        const cleanup = { signal: AbortSignal.timeout(30_000) };
        const deleting = await sandbox.delete({ force: true }, cleanup);
        await deleting?.waitUntilDeleted(cleanup);
      }
    } finally {
      await client.close();
    }
  }
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch((error: unknown) => {
    console.error(error);
    process.exitCode = 1;
  });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}
