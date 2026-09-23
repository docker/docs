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

- A Docker account with access to Cloud Sandboxes
- An image available to your account that includes `sh`
- A [personal access token](/manuals/security/access-tokens/personal-access-tokens.md)
- Node.js 20 or later and npm
- `curl` and `jq`

Your credentials need permission to list images, create and read sandboxes,
run commands, and delete sandboxes. The required API permissions are
`imagesRead`, `sandboxesCreate`, `sandboxesRead`, `sandboxesCredential`,
`sandboxesExec`, and `sandboxesDelete`.

## Create an access token

To authenticate API requests, exchange your Docker ID and personal access
token (PAT) for a short-lived Docker Hub access token:

```console
$ export DOCKER_ID=<your-docker-id>
$ export DOCKER_PAT=<your-personal-access-token>
$ export DOCKER_ACCESS_TOKEN=$(curl --silent --show-error --fail --request POST \
  --url https://hub.docker.com/v2/auth/token \
  --header "Content-Type: application/json" \
  --data "$(jq -n --arg identifier "$DOCKER_ID" --arg secret "$DOCKER_PAT" \
    '{identifier: $identifier, secret: $secret}')" \
  | jq -er '.access_token')
$ unset DOCKER_PAT
```

The program reads `DOCKER_ACCESS_TOKEN` from your environment and uses it to
authenticate with the management API.

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

## Select an image

Choose an image for the sandbox. It must be available to your account and
include `sh` to run the example command. List your images:

```console
$ curl --silent --show-error --fail \
  --url https://connect.docker.com/sandboxes/v1/images \
  --header "Authorization: Bearer $DOCKER_ACCESS_TOKEN" \
  | jq '.images[] | {name, displayName}'
```

Copy the image's complete `name`, such as `images/<uid>`, and save it in
`SBX_IMAGE`:

```console
$ export SBX_IMAGE=<IMAGE_RESOURCE_NAME>
```

## Create and use a sandbox

Create a file named `index.ts` with the following code. The program creates a
sandbox and waits until it's running before sending the command. It then
attempts to delete the sandbox, even if the command fails.

```typescript
import { randomUUID } from 'node:crypto';
import { ResourceWaitError, SandboxesClient } from '@docker/sandboxes-api';

const token = process.env.DOCKER_ACCESS_TOKEN;
const image = process.env.SBX_IMAGE;
if (!token || !image) {
  throw new Error('Set DOCKER_ACCESS_TOKEN and SBX_IMAGE');
}

const client = new SandboxesClient({
  baseUrl: 'https://connect.docker.com/sandboxes',
  headers: { Authorization: `Bearer ${token}` },
});

const requestId = randomUUID();
console.log('Create request ID:', requestId);
const sandbox = await client.createSandboxAndWait(
  { body: { image }, headers: { 'Idempotency-Key': requestId } },
  { timeoutMs: 300_000 },
).catch((error) => {
  if (error instanceof ResourceWaitError && error.resource) {
    console.error('Inspect accepted sandbox:', error.resource.name);
  }
  throw error;
});
console.log('Created', sandbox.name);

try {
  const endpoint = sandbox.core?.endpoint;
  if (!endpoint) {
    throw new Error('Sandbox returned no endpoint');
  }

  const sandboxClient = await client.forEndpoint(endpoint, ['sandboxesExec'], {
    signal: AbortSignal.timeout(30_000),
  });
  const result = await sandboxClient.run(
    ['sh', '-lc', 'echo "Hello from Docker Sandboxes"'],
    {},
    { signal: AbortSignal.timeout(30_000) },
  );
  if (result.exitCode !== 0 || result.incomplete) {
    throw new Error(`Command failed or output was incomplete: ${result.stderr}`);
  }
  console.log(result.stdout.trim());
} finally {
  try {
    const current = await client.getSandbox({ name: sandbox.name }, {
      signal: AbortSignal.timeout(30_000),
    });
    if (!current) {
      throw new Error('Cleanup read returned no sandbox');
    }
    await client.sandbox(current).deleteAndWait({}, { timeoutMs: 300_000 });
    console.log('Deleted', sandbox.name);
  } catch (error) {
    console.error('Cleanup failed; inspect sandbox:', sandbox.name, error);
    process.exitCode = 1;
  }
}
```

The program connects to two API endpoints: the management API to create and
delete the sandbox, and the sandbox's own API to run the command. Your Docker
Hub token authenticates the management requests. For command execution,
`forEndpoint` obtains a token limited to that sandbox and creates the client
that sends the command.

## Run the program

Run the program with `tsx`:

```console
$ npx tsx index.ts
```

On success, the program prints `Hello from Docker Sandboxes` and deletes the
sandbox. It also prints the create request ID and sandbox resource name.

## If the program fails

Check the printed sandbox name before running the program again. A sandbox
can remain if creation times out or deletion fails, and another run can create
another sandbox. Read the resource with `getSandbox` and delete it when you
no longer need it.

If a lost response leaves you without a sandbox name, use the printed request
ID as the idempotency key when retrying the same create request. See
[Errors and retries](errors.md#retry-without-duplicating-work) for how to retry
without duplicating work.

## Next steps

- Install an [SDK](sdks.md) for Go or Python.
- Read [API concepts](concepts.md) to learn about resource names, lifecycle
  states, and supported Cloud options.
- Review [Errors and retries](errors.md) before adding recovery logic.
