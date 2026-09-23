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
import { oauth, SandboxesClient } from '@docker/sandboxes-api';

const client = new SandboxesClient({
  auth: oauth({
    onVerification({ verificationUriComplete, verificationUri, userCode }) {
      console.log(`Open ${verificationUriComplete ?? verificationUri}`);
      console.log(`Verification code: ${userCode}`);
    },
  }),
});

try {
  let sandbox = await client.kits.launch('shell', {
    resources: { cpus: 2, memoryMib: 4096 },
  });
  console.log('Sandbox:', sandbox.name);
  sandbox = await sandbox.waitUntilRunning();

  const result = await sandbox.processes.run({
    args: ['echo', 'Hello from Docker Sandboxes'],
  });
  console.log(result.stdout.trim());

  const deleting = await sandbox.delete({ force: true });
  await deleting?.waitUntilDeleted();
  console.log('Deleted', sandbox.name);
} finally {
  await client.close();
}
```

`kits.launch` creates the sandbox, and `waitUntilRunning` waits until it's
ready to run commands. `processes.run` runs the command and collects its output.
The program then deletes the sandbox and waits for deletion to finish. The
`force` option permits deletion while the sandbox is running.

The SDK handles sign-in, access tokens, and the connection to the sandbox.
Closing the client releases its local resources.

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

If the program fails before deleting the sandbox, use the printed sandbox
name to retrieve it with `client.get(name)` and delete it when you're finished.
Closing the client doesn't delete the sandbox.

## Next steps

- Install an [SDK](sdks.md) for Go or Python.
- Read [API concepts](concepts.md) to learn about resource names, lifecycle
  states, and supported Cloud options.
- Review [Errors and retries](errors.md) before adding recovery logic.
