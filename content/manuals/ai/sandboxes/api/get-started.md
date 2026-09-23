---
title: Run your first cloud sandbox
linkTitle: Get started
description: Use the TypeScript SDK to create a cloud sandbox, run a command, and delete the sandbox.
keywords: Docker Sandboxes API tutorial, cloud sandbox TypeScript, create sandbox, sandbox SDK quickstart
weight: 5
---

Create a cloud sandbox, run a command inside it, and delete it using the Docker
Sandboxes TypeScript SDK. The example prints `Hello from Docker Sandboxes`
from inside the sandbox.

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

The example uses the public Alpine Linux image `docker.io/library/alpine:3.22`
with two CPUs and 4 GiB of memory on `linux/amd64`.

Create a file named `index.ts` with the following code. The program prompts
you to sign in, creates a sandbox, runs a command, and deletes the sandbox.

```typescript
import { randomUUID } from 'node:crypto';
import { oauth, SandboxesClient, WorkflowError } from '@docker/sandboxes-api';

const auth = oauth({
  onVerification({ verificationUriComplete, verificationUri, userCode }) {
    console.log(`Open ${verificationUriComplete ?? verificationUri}`);
    console.log(`Verification code: ${userCode}`);
  },
});
const client = new SandboxesClient({ auth });
const requestId = randomUUID();

try {
  await auth.getAccessToken();
  console.log('Create request ID:', requestId);
  const result = await client.withSandbox(
    {
      imageRef: 'docker.io/library/alpine:3.22',
      resources: { cpus: 2, memoryMib: '4096' },
      platform: { os: 'linux', architecture: 'amd64' },
    },
    async (sandbox) => {
      console.log('Sandbox:', sandbox.name);
      return sandbox.processes.run({
        args: ['sh', '-lc', 'echo "Hello from Docker Sandboxes"'],
      });
    },
    {
      idempotencyKey: requestId,
      timeoutMs: 300_000,
      cleanup: { timeoutMs: 30_000 },
    },
  );
  if (result.exitCode !== 0 || result.incomplete) {
    throw new Error(`Command failed or output was incomplete: ${result.stderr}`);
  }
  console.log(result.stdout.trim());
} catch (error) {
  if (error instanceof WorkflowError) {
    console.error('Failed during:', error.phase);
    console.error('Sandbox:', error.resource?.name);
    console.error('Cleanup:', error.cleanup);
  }
  throw error;
} finally {
  await client.close();
}
```

`withSandbox` waits for the sandbox to run before calling your function. It
attempts deletion after the function finishes, including when the function
fails. Creation, waiting, and command execution share a five-minute deadline;
cleanup has a separate 30-second deadline.

The SDK handles access tokens and the connection to the sandbox's endpoint.
Closing the client releases local connections and credentials.

## Run the program

Run the program with `tsx`:

```console
$ npx tsx index.ts
```

Open the printed verification URL and sign in with the Docker account that
has your cloud subscription. The program continues after sign-in. On success,
it prints `Hello from Docker Sandboxes` after deleting the sandbox.

For CI jobs and other unattended applications, use
[PAT authentication](authentication.md#authenticate-automation-with-a-pat).

## If the program fails

If the program reports a workflow error, check its cleanup status. If cleanup
failed, use the printed sandbox name to inspect it with `client.get(name)`
and delete it when you no longer need it. Closing the client doesn't delete
remote sandboxes.

If a lost response leaves you without a sandbox name, use the printed request
ID as the idempotency key when retrying the same create request. See
[Errors and retries](errors.md#retry-without-duplicating-work) for how to retry
without duplicating work.

## Next steps

- Install an [SDK](sdks.md) for Go or Python.
- Read [API concepts](concepts.md) to learn about resource names, lifecycle
  states, and supported Cloud options.
- Review [Errors and retries](errors.md) before adding recovery logic.
