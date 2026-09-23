---
title: Run your first cloud sandbox
linkTitle: Get started
description: Use the TypeScript SDK to create a cloud sandbox, run a command, and delete the sandbox.
keywords: Docker Sandboxes API tutorial, cloud sandbox TypeScript, create sandbox, sandbox SDK quickstart
weight: 5
---

> [!NOTE]
> The Docker Sandboxes API and SDKs are experimental. Features and behavior may change.

Create a cloud sandbox, run a command inside it, and delete it using the Docker
Sandboxes TypeScript SDK. The example prints `Hello from Docker Sandboxes`
from inside the sandbox. It uses the `shell` kit, so you don't need a
model-provider credential.

## Prerequisites

To follow this tutorial, you need:

- A Docker account with an active
  [Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access)
- Node.js 20 or later and npm

Your credentials need permission to create and read sandboxes, run commands,
and delete sandboxes. The required API permissions are `sandboxesCreate`,
`sandboxesRead`, `sandboxesCredential`, `sandboxesExec`, and `sandboxesDelete`.

## Create a project

Create a directory and initialize a Node.js project:

```console
$ mkdir sandboxes-api-tutorial
$ cd sandboxes-api-tutorial
$ npm init --yes
$ npm pkg set type=module
```

Install the SDK and TypeScript tooling:

```console
$ npm install @docker/sandboxes-api
$ npm install --save-dev tsx typescript @types/node
```

## Create and use a sandbox

A kit supplies the sandbox's image and configuration for an agent or tool.
This example launches the bundled `shell` kit with two CPUs and 4 GiB of
memory. See [Compute sizes and limits](limits.md) for other sizes.
Cloud compute is billed to your subscription.

Create a file named `index.ts` with the following code. The program prompts
you to sign in, creates a sandbox, runs a command, and deletes the sandbox.

```typescript
import { randomUUID } from 'node:crypto';
import { oauth, SandboxesClient, type Sandbox } from '@docker/sandboxes-api';

const auth = oauth({
  onVerification({ verificationUriComplete, verificationUri, userCode }) {
    console.log(`Open ${verificationUriComplete ?? verificationUri}`);
    console.log(`Verification code: ${userCode}`);
  },
});
const client = new SandboxesClient({ auth });
const requestId = randomUUID();
let sandbox: Sandbox | undefined;

try {
  await auth.getAccessToken();
  console.log('Create request ID:', requestId);
  const operation = { signal: AbortSignal.timeout(300_000) };
  sandbox = await client.kits.launch(
    'shell',
    { resources: { cpus: 2, memoryMib: 4096 } },
    { ...operation, idempotencyKey: requestId },
  );
  console.log('Sandbox:', sandbox.name);
  sandbox = await sandbox.waitUntilRunning(operation);
  const result = await sandbox.processes.run(
    { args: ['echo', 'Hello from Docker Sandboxes'] },
    operation,
  );
  if (result.exitCode !== 0 || result.incomplete) {
    throw new Error(`Command failed or output was incomplete: ${result.stderr}`);
  }
  console.log(result.stdout.trim());
} finally {
  try {
    if (sandbox) {
      const cleanup = { signal: AbortSignal.timeout(30_000) };
      const deleting = await sandbox.delete({ force: true }, cleanup);
      await deleting?.waitUntilDeleted(cleanup);
      console.log('Deleted', sandbox.name);
    }
  } catch (error) {
    console.error('Cleanup failed; inspect sandbox:', sandbox?.name);
    throw error;
  } finally {
    await client.close();
  }
}
```

`kits.launch` returns after the API accepts creation. `waitUntilRunning`
waits until you can run commands. Creation, waiting, and command execution
share a five-minute deadline, starting after browser sign-in.

The `finally` block attempts deletion even if waiting or command execution
fails. Cleanup has a separate 30-second deadline. The `force` option permits
deletion of a running sandbox.

The SDK handles access tokens and the connection to the sandbox's endpoint.
Closing the client releases its local resources. It doesn't revoke your
Docker sign-in or delete remote sandboxes.

## Run the program

Run the program with `tsx`:

```console
$ npx tsx index.ts
```

Open the printed verification URL and sign in with the Docker account that
has your cloud subscription. The program continues after sign-in. On success,
it prints `Hello from Docker Sandboxes`, then confirms sandbox deletion.

For CI jobs and other unattended applications, use
[PAT authentication](authentication.md#authenticate-automation-with-a-pat).

## If the program fails

If the program reports a cleanup failure, use the printed sandbox name to
inspect it with `client.get(name)` and delete it when you no longer need it.
A timeout does not prove that creation failed or deletion succeeded.

If a lost response leaves you without a sandbox name, use the printed request
ID as the idempotency key when retrying the same create request. See
[Errors and retries](errors.md#retry-without-duplicating-work) for how to retry
without duplicating work.

## Next steps

- Install an [SDK](sdks.md) for Go or Python.
- Read [API concepts](concepts.md) to learn about resource names, lifecycle
  states, and supported Cloud options.
- Review [Errors and retries](errors.md) before adding recovery logic.
